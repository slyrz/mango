package main

type NodeType uint32

const (
	DocumentNode NodeType = iota
	BlockNode
	BreakNode
	ListNode
	SectionNode
	TextBoldNode
	TextNode
	TextUnderlineNode
)

type Node struct {
	Type   NodeType
	Text   string
	Childs []*Node
}
