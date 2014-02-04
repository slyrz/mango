package markup

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	TOKEN_EOL = iota
	TOKEN_SECTION
	TOKEN_TEXT
	TOKEN_INDENT
	TOKEN_BLOCKITEM
	TOKEN_LISTITEM
	TOKEN_UNDERLINE
	TOKEN_STAR
)

var (
	reBlockItem = regexp.MustCompile(`^\s*>\s*`)
	reListItem  = regexp.MustCompile(`^(\w+|\*)\)\s*`)
	reSectionGo = regexp.MustCompile(`^([A-Z][a-z]*\s*)+:$`)
	reSectionMd = regexp.MustCompile(`^(-+|=+)$`)
)

type Token struct {
	Kind int
	Text string
}

func NewToken(kind int) *Token {
	return NewTokenWithText(kind, "")
}

func NewTokenWithText(kind int, text string) *Token {
	return &Token{kind, text}
}

// Return true if token is of given kind.
func (t *Token) Is(kind int) bool {
	return t.Kind == kind
}

// A sequence of tokens. Offers the function Are() which checks if the tokens
// match a given order.
type Tokens []*Token

// Returns true if the tokens match kinds.
func (tokens Tokens) Are(kinds ...int) bool {
	if len(kinds) > len(tokens) {
		return false
	}
	for i, kind := range kinds {
		if !tokens[i].Is(kind) {
			return false
		}
	}
	return true
}

func kindToString(kind int) string {
	switch kind {
	case TOKEN_EOL:
		return "Eol"
	case TOKEN_TEXT:
		return "Text"
	case TOKEN_INDENT:
		return "Indent"
	case TOKEN_UNDERLINE:
		return "Underline"
	case TOKEN_STAR:
		return "Star"
	case TOKEN_SECTION:
		return "Section"
	case TOKEN_LISTITEM:
		return "Listitem"
	case TOKEN_BLOCKITEM:
		return "Blockitem"
	}
	return "Unknown"
}

func (t *Token) String() string {
	return fmt.Sprintf("Token{%s, '%s'}", kindToString(t.Kind), t.Text)
}

type Tokenizer struct {
	text   string
	skip   int
	tokens []*Token
}

func NewTokenizer() *Tokenizer {
	return new(Tokenizer)
}

func (t *Tokenizer) TokenizeString(text string) ([]*Token, error) {
	return t.TokenizeLines(strings.Split(text, "\n"))
}

func (t *Tokenizer) TokenizeLines(lines []string) ([]*Token, error) {
	// We always start with an end of line token, so we don't need edge cases
	// to detect sections at the start.
	t.tokens = []*Token{NewToken(TOKEN_EOL)}
	for _, line := range lines {
		t.text = line
		t.skip = 0
		t.parseLine()
	}
	return t.tokens, nil
}

// Add a new token of specified kind with an optional text to the token list.
func (t *Tokenizer) addToken(kind int, text ...string) {
	var token *Token = nil
	if len(text) == 0 {
		token = NewToken(kind)
	} else {
		token = NewTokenWithText(kind, text[0])
	}
	t.tokens = append(t.tokens, token)
}

// Skip the next n characters of the current line.
func (t *Tokenizer) skipChars(n int) {
	if n < len(t.text) {
		t.text = t.text[n:]
	} else {
		t.skipLine()
	}
}

// Skip the current line.
func (t *Tokenizer) skipLine() {
	t.text = ""
}

// Parse and skip the indentation.
func (t *Tokenizer) parseIndentation() {
	skipped := 0 // Number of chars consumed
	spaces := 0  // Number of spaces detected
SkipLoop:
	for _, char := range t.text {
		switch char {
		case ' ':
			spaces += 1
		case '\t':
			spaces += 4
		default:
			break SkipLoop
		}
		if spaces >= 4 {
			t.addToken(TOKEN_INDENT)
			spaces = 0
		}
		skipped++
	}
	t.skipChars(skipped)
}

// Try to parse a list item declaration.
func (t *Tokenizer) parseListItem() {
	if match := reListItem.FindStringSubmatch(t.text); match != nil {
		matchTotal := match[0] // everything including ')' and trailing whitespace
		matchItem := match[1]  // the word or '*' before ')'
		t.addToken(TOKEN_LISTITEM, matchItem)
		t.skipChars(len(matchTotal))
	}
}

// Return a slice containing the n last tokens.
func (t *Tokenizer) lastTokens(n int) Tokens {
	i := 0
	if n > 0 && n < len(t.tokens) {
		i = len(t.tokens) - n
	}
	return Tokens(t.tokens[i:])
}

// Parse Markdown-style section declaration: A line of text followed by a line
// of dashes or equal signs.
//
//		This Is A Title  or  This Is A Title
//		---------------      ===============
//
func (t *Tokenizer) parseSectionMd() bool {
	if len(t.tokens) < 3 || !reSectionMd.MatchString(t.text) {
		return false
	}
	lastTokens := t.lastTokens(3)
	if !lastTokens.Are(TOKEN_EOL, TOKEN_TEXT, TOKEN_EOL) {
		return false
	}
	lastTokens[1].Kind = TOKEN_SECTION
	t.skipLine()
	return true
}

// Parse Go-style section declaration: a line in title case, ending with a
// colon and followed by an empty line.
//
//		This Is A Title:
//
//
func (t *Tokenizer) parseSectionGo() bool {
	if len(t.tokens) < 3 {
		return false
	}
	lastTokens := t.lastTokens(3)
	if !lastTokens.Are(TOKEN_EOL, TOKEN_TEXT, TOKEN_EOL) {
		return false
	}
	if !reSectionGo.MatchString(lastTokens[1].Text) {
		return false
	}

	lastTokens[1].Text = strings.TrimRight(lastTokens[1].Text, ":")
	lastTokens[1].Kind = TOKEN_SECTION
	t.skipLine()
	return true
}

func (t *Tokenizer) parseSection() bool {
	if t.parseSectionMd() {
		return true
	}
	if t.parseSectionGo() {
		return true
	}
	return false
}

func (t *Tokenizer) parseBlockItem() {
	match := reBlockItem.FindString(t.text)
	if len(match) > 0 {
		t.skipChars(len(match))
		t.addToken(TOKEN_BLOCKITEM, t.text)
		t.skipLine()
	}
}

func (t *Tokenizer) parseLine() {
	if t.parseSection() {
		return
	}
	t.parseIndentation()
	t.parseBlockItem()
	t.parseListItem()

	lastPos := 0
	for currPos, char := range t.text {
		kind := 0
		switch char {
		case '*':
			kind = TOKEN_STAR
		case '_':
			kind = TOKEN_UNDERLINE
		}
		if kind != 0 {
			// Everything between the current token and the last token is text.
			if lastPos < currPos {
				t.addToken(TOKEN_TEXT, t.text[lastPos:currPos])
			}
			t.addToken(kind)
			lastPos = currPos + 1
		}
	}
	if lastPos < len(t.text) {
		t.addToken(TOKEN_TEXT, t.text[lastPos:])
	}
	t.addToken(TOKEN_EOL)
}
