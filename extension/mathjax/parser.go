package mathjax

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type mathjaxInlineParser struct{}

func NewMathjaxInlineParser() parser.InlineParser {
	return &mathjaxInlineParser{}
}

func (w *mathjaxInlineParser) Trigger() []byte {
	panic("unimplemented")
}

func (p *mathjaxInlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	panic("unimplemented")
}
