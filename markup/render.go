package markup

type Renderer interface {
	Block(n *Node)
	Group(n *Node)
	List(n *Node)
	ListItem(n *Node)
	Section(text string)
	Text(text string)
	TextBold(text string)
	TextUnderline(text string)
	Space()
	Break()
}

func Render(r Renderer, node *Node) {
	switch node.Kind {
	case NODE_BLOCK:
		r.Block(node)
	case NODE_GROUP:
		r.Group(node)
	case NODE_LIST:
		r.List(node)
	case NODE_LISTITEM:
		r.ListItem(node)
	case NODE_SECTION:
		r.Section(node.Text)
	case NODE_SPACE:
		r.Space()
	case NODE_BREAK:
		r.Break()
	case NODE_TEXT:
		r.Text(node.Text)
	case NODE_TEXTBOLD:
		r.TextBold(node.Text)
	case NODE_TEXTUNDERLINE:
		r.TextUnderline(node.Text)
	}
}

func RenderChilds(r Renderer, node *Node) {
	for _, child := range node.Childs {
		Render(r, child)
	}
}
