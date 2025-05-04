package mathjax

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type mathjax struct {
}

func NewMathjax() goldmark.Extender {
	return &mathjax{}
}

func (e *mathjax) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewMathjaxInlineParser(), 501),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewMathjaxRenderer(), 1001),
	))
}
