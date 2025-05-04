package mathjax

import (
	"juliangcalderon/onyx/node"

	"github.com/dop251/goja"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type mathjax struct {
	js node.MathjaxInstance
}

func NewIsolatedMathjax() (*mathjax, error) {
	runtime := goja.New()
	module, err := node.NewMathjax(runtime)
	if err != nil {
		return nil, err
	}
	js := module()

	return &mathjax{
		js: js,
	}, nil
}

func (e *mathjax) CSS() string {
	return e.js.Css()
}

func (e *mathjax) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewMathjaxInlineParser(), 501),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewMathjaxRenderer(e.js), 1001),
	))
}
