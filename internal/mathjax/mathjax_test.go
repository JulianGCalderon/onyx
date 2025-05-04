package mathjax

import (
	"juliangcalderon/onyx/javascript"
	"testing"

	_ "embed"

	"github.com/dop251/goja"
)

func TestInteroperate(t *testing.T) {
	runtime := goja.New()
	mathjax, err := runtime.RunString(javascript.Mathjax)
	if err != nil {
		t.Fatalf("%v", err)
	}

	mathjaxF, callable := goja.AssertFunction(mathjax)
	if !callable {
		t.Fatalf("%v should be callable", mathjax)
	}

	builder, err := mathjaxF(goja.Undefined())
	if err != nil {
		t.Fatalf("%v", mathjax)
	}

	cssV := builder.ToObject(runtime).Get("css")
	cssF, callable := goja.AssertFunction(cssV)
	if !callable {
		t.Fatalf("%v should be callable", cssV)
	}

	renderV := builder.ToObject(runtime).Get("render")
	renderF, callable := goja.AssertFunction(renderV)
	if !callable {
		t.Fatalf("%v should be callable", renderV)
	}

	math := runtime.ToValue(`a^2 + b^2 = c^2`)
	display := runtime.ToValue("display")

	pageHTML, err := renderF(goja.Undefined(), math, display)
	if err != nil {
		t.Fatalf("%v", err)
	}

	pageCSS, err := cssF(goja.Undefined())
	if err != nil {
		t.Fatalf("%v", err)
	}

	_ = pageHTML
	_ = pageCSS
}
