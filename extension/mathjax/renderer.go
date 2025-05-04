package mathjax

import (
	"github.com/yuin/goldmark/renderer"
)

type mathjaxRenderer struct{}

func NewMathjaxRenderer() renderer.NodeRenderer {
	return &mathjaxRenderer{}
}

func (m *mathjaxRenderer) RegisterFuncs(renderer.NodeRendererFuncRegisterer) {
	panic("unimplemented")
}
