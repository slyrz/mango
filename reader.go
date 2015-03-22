package main

import (
	"bytes"
	"regexp"
	"strings"
)

type Reader interface {
	Read(string) (*Node, error)
}

func skip(text string, n int) string {
	if n < 0 || n >= len(text) {
		return ""
	}
	return text[n:]
}

func nextLine(text string) (string, string) {
	end := strings.IndexRune(text, '\n')
	if end == -1 {
		return text, ""
	}
	return text[:end], text[end+1:]
}

type MarkupReader struct {
	// empty
}

func NewMarkupReader() *MarkupReader {
	return &MarkupReader{}
}

var markupEmph = map[rune]NodeType{
	'*': TextBoldNode,
	'_': TextUnderlineNode,
	'`': TextNode,
}

// parseText transforms text into TextNode, TextBoldNode and TextUnderlineNode.
func (m *MarkupReader) parseText(text string, callback func(NodeType, string)) {
	var (
		buff bytes.Buffer
		skip = -1
	)
	for i, r := range text {
		if i <= skip {
			continue
		}
		if k, ok := markupEmph[r]; ok {
			if end := i + 1 + strings.IndexRune(text[i+1:], r); end > i {
				if buff.Len() > 0 {
					callback(TextNode, buff.String())
					buff.Reset()
				}
				callback(k, text[i+1:end])
				skip = end
				continue
			}
		}
		buff.WriteRune(r)
	}
	if buff.Len() > 0 {
		callback(TextNode, buff.String())
	}
}

var markupMatchers = []struct {
	Regex     *regexp.Regexp
	Multiline bool
	Type      NodeType
}{
	{regexp.MustCompile(`^((\w+\s*)+)\n(?:-+|=+)\n`), true, SectionNode},
	{regexp.MustCompile(`^(([A-Z]\w+\s*)+):\n\n`), true, SectionNode},
	{regexp.MustCompile(`^(\w+|\*)\)\s*`), false, ListNode},
	{regexp.MustCompile(`^(?:\>|\ {4}|\t)\s?(.*)\n`), true, BlockNode},
	{regexp.MustCompile(`^(\n+)`), true, BreakNode},
}

func (m *MarkupReader) parse(text string, callback func(NodeType, string)) {
	var line string
OuterLoop:
	for len(text) > 0 {
	InnerLoop:
		for _, matcher := range markupMatchers {
			match := matcher.Regex.FindStringSubmatch(text)
			if match == nil {
				continue
			}
			all := match[0]
			val := match[1]
			callback(matcher.Type, val)
			// Skip the matching part.
			text = skip(text, len(all))
			// If the matcher consumes full lines, text is at the beginning
			// of a new line and we go back to the OuterLoop.
			if matcher.Multiline {
				continue OuterLoop
			} else {
				break InnerLoop
			}
		}
		line, text = nextLine(text)
		m.parseText(line, callback)
	}
}

func (m *MarkupReader) Read(text string) (*Node, error) {
	root := &Node{Type: DocumentNode}
	curr := root
	m.parse(text, func(kind NodeType, text string) {
		node := &Node{Type: kind, Text: text}
		if kind == SectionNode {
			root.Childs = append(root.Childs, node)
			curr = node
		} else {
			curr.Childs = append(curr.Childs, node)
		}
	})
	return root, nil
}

type PlainReader struct {
	// empty
}

func NewPlainReader() *PlainReader {
	return &PlainReader{}
}

func (p *PlainReader) Read(text string) (*Node, error) {
	root := &Node{Type: DocumentNode}
	empty := 0
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			if empty > 0 {
				root.Childs = append(root.Childs, &Node{Type: BreakNode})
			}
			root.Childs = append(root.Childs, &Node{Type: TextNode, Text: line})
			empty = 0
		} else {
			empty++
		}
	}
	return root, nil
}
