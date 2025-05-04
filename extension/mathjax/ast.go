package mathjax

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

type InlineMath struct {
	ast.BaseInline
}

func (n *InlineMath) Inline() {
}

func (n *InlineMath) IsBlank(source []byte) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		text := c.(*ast.Text).Segment
		if !util.IsBlank(text.Value(source)) {
			return false
		}
	}
	return true
}

func (n *InlineMath) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindInlineMath = ast.NewNodeKind("InlineMath")

func (n *InlineMath) Kind() ast.NodeKind {
	return KindInlineMath
}

func NewInlineMath() *InlineMath {
	return &InlineMath{
		BaseInline: ast.BaseInline{},
	}
}

type MathBlock struct {
	ast.BaseBlock
	Info *ast.Text

	language []byte
}

func (n *MathBlock) Language(source []byte) []byte {
	if n.language == nil && n.Info != nil {
		segment := n.Info.Segment
		info := segment.Value(source)
		i := 0
		for ; i < len(info); i++ {
			if info[i] == ' ' {
				break
			}
		}
		n.language = info[:i]
	}
	return n.language
}

func (n *MathBlock) IsRaw() bool {
	return true
}

func (n *MathBlock) Dump(source []byte, level int) {
	m := map[string]string{}
	if n.Info != nil {
		m["Info"] = fmt.Sprintf("\"%s\"", n.Info.Text(source))
	}
	ast.DumpHelper(n, source, level, m, nil)
}

var KindMathBlock = ast.NewNodeKind("MathCodeBlock")

func (n *MathBlock) Kind() ast.NodeKind {
	return KindMathBlock
}

func (n *MathBlock) Text(source []byte) []byte {
	return n.Lines().Value(source)
}

func NewMathBlock(info *ast.Text) *MathBlock {
	return &MathBlock{
		BaseBlock: ast.BaseBlock{},
		Info:      info,
	}
}
