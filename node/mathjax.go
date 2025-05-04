package node

import (
	"github.com/dop251/goja"
)

type Mathjax func() MathjaxInstance

type MathjaxInstance struct {
	Css    func() string
	Render func(code string, ty string) string
}

func NewMathjax(runtime *goja.Runtime) (Mathjax, error) {
	module, err := runtime.RunString(MathjaxBundle)
	if err != nil {
		return nil, err
	}

	var mathjax Mathjax
	err = runtime.ExportTo(module, &mathjax)
	if err != nil {
		return nil, err
	}

	return mathjax, nil
}
