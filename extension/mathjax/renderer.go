package mathjax

import (
	"juliangcalderon/onyx/node"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type mathjaxRenderer struct {
	js node.MathjaxInstance
}

func NewMathjaxRenderer(js node.MathjaxInstance) renderer.NodeRenderer {
	return &mathjaxRenderer{
		js: js,
	}
}

func (m *mathjaxRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindInlineMath, m.renderInlineMath)
	reg.Register(KindMathBlock, m.renderMathBlock)
}

func (r *mathjaxRenderer) renderInlineMath(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		content := make([]byte, 0)
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			value := segment.Value(source)
			content = append(content, value...)
		}
		html := r.js.Render(string(content), "inline")

		_, err := w.WriteString(html)
		if err != nil {
			return ast.WalkStop, err
		}

		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *mathjaxRenderer) renderMathBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {

		content := make([]byte, 0)
		l := n.Lines().Len()
		for i := 0; i < l; i++ {
			line := n.Lines().At(i)
			content = append(content, line.Value(source)...)
		}

		html := r.js.Render(string(content), "display")

		_, err := w.WriteString(html)
		if err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}
