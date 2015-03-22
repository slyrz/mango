package main

type Builder struct {
	r Reader
	w Writer
}

func NewBuilder(r Reader, w Writer) *Builder {
	return &Builder{r, w}
}

func (b *Builder) write(n *Node) {
	switch n.Type {
	case BlockNode:
		b.w.Block(n.Text)
	case ListNode:
		b.w.List(n.Text)
	case SectionNode:
		b.w.Section(n.Text)
	case TextNode:
		b.w.Text(n.Text)
	case TextBoldNode:
		b.w.TextBold(n.Text)
	case TextUnderlineNode:
		b.w.TextUnderline(n.Text)
	case BreakNode:
		b.w.Break(n.Text)
	}
	for _, c := range n.Childs {
		b.write(c)
	}
}

func (b *Builder) Build(f *File) (string, error) {
	root, err := b.r.Read(f.Doc)
	if err != nil {
		return "", err
	}

	b.w.Section("Name")
	for _, node := range root.Childs {
		if node.Type == SectionNode {
			break
		}
		b.write(node)
	}

	b.w.Section("Synopsis")
	b.w.Text(f.Name)
	b.w.TextUnderline("[options...]")
	b.w.TextUnderline("[argument...]")

	b.w.Section("Options")
	for _, opt := range f.Flags {
		var doc string

		b.w.Flag(opt.Name, opt.Short, opt.Param)
		if opt.Doc != "" {
			doc = opt.Doc
		} else {
			doc = opt.Usage
		}
		node, err := b.r.Read(doc)
		if err != nil {
			panic(err)
		}
		b.write(node)
	}

	for _, node := range root.Childs {
		if node.Type == SectionNode {
			b.write(node)
		}
	}

	return b.w.Done(), nil
}
