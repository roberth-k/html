package html

import (
	"context"
	"io"
	"strings"
)

type Context struct {
	context.Context
	depth int
}

type PartialFunc func(ctx *Context) Node

func RenderString(node Node) (string, error) {
	var buf strings.Builder

	if err := Render(&buf, node); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Render(w io.Writer, node Node) error {
	ctx := &Context{nil, 0}
	return node.RenderHTML(ctx, w)
}
