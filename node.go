package html

import (
	"fmt"
	html_template "html/template"
	"io"
	"strings"
)

//go:generate ./generate-elems.sh elems.go

type Node interface {
	RenderHTML(ctx *Context, w io.Writer) error
}

func Textf(text string, args ...interface{}) TextNode {
	return Text(fmt.Sprintf(text, args...))
}

func Text(text string) TextNode {
	return TextNode(text)
}

type TextNode string

func (node TextNode) RenderHTML(ctx *Context, w io.Writer) error {
	_, _ = fmt.Fprint(w, strings.Repeat("\t", ctx.depth))
	html_template.HTMLEscape(w, []byte(string(node)))
	_, _ = w.Write([]byte("\n"))
	return nil
}

func Attr(k, v string) AttrNode {
	return AttrNode{k, v}
}

type AttrNode struct {
	Name  string
	Value string
}

func (node AttrNode) RenderHTML(ctx *Context, w io.Writer) error {
	_, _ = fmt.Fprintf(
		w, `%s="%s"`,
		html_template.HTMLEscapeString(node.Name),
		html_template.HTMLEscapeString(node.Value))
	return nil
}

func Func(f func() Node) FuncNode {
	return FuncNode{f}
}

type FuncNode struct {
	Func func() Node
}

func (node FuncNode) RenderHTML(ctx *Context, w io.Writer) error {
	// rendered via ElemNode
	return nil
}

type Children []Node

func (node Children) RenderHTML(ctx *Context, w io.Writer) error {
	for _, child := range node {
		if err := child.RenderHTML(ctx, w); err != nil {
			return err
		}
	}

	return nil
}

func Elem(name string, children ...Node) ElemNode {
	node := ElemNode{
		Name:       name,
		Attributes: nil,
		Children:   nil,
	}

	for _, child := range children {
		node.Add(child)
	}

	return node
}

type ElemNode struct {
	Name       string
	Attributes []AttrNode
	Children   Children
}

// todo: move to generate-elems.sh
var _voidElems = []string{"area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr"}

func (node ElemNode) RenderHTML(ctx *Context, w io.Writer) error {
	_, _ = fmt.Fprintf(
		w, "%s<%s",
		strings.Repeat("\t", ctx.depth),
		node.Name)

	if node.Attributes != nil {
		for _, attr := range node.Attributes {
			_, _ = fmt.Fprintf(w, " ")
			if err := attr.RenderHTML(ctx, w); err != nil {
				return err
			}
		}
	}

	isVoidElem := false
	for _, name := range _voidElems {
		if node.Name == name {
			isVoidElem = true
			break
		}
	}

	switch {
	case isVoidElem:
		_, _ = fmt.Fprintf(w, " />\n")
	default:
		_, _ = fmt.Fprintf(w, ">\n")
		ctx.depth++
		if err := node.Children.RenderHTML(ctx, w); err != nil {
			return err
		}
		ctx.depth--
		_, _ = fmt.Fprintf(w, "%s</%s>\n", strings.Repeat("\t", ctx.depth), node.Name)
	}

	return nil
}

func (node *ElemNode) Add(child Node) {
	switch child := child.(type) {
	case AttrNode:
		node.Attributes = append(node.Attributes, child)
	case *AttrNode:
		node.Attributes = append(node.Attributes, *child)
	case FuncNode:
		if y := child.Func(); y != nil {
			node.Add(y)
		}
	case Children:
		for _, subchild := range child {
			node.Add(subchild)
		}
	default:
		node.Children = append(node.Children, child)
	}
}
