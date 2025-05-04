package node_test

import (
	"juliangcalderon/onyx/node"
	"testing"

	"github.com/dop251/goja"
)

func TestInit(t *testing.T) {
	runtime := goja.New()

	mathjax, err := node.NewMathjax(runtime)
	if err != nil {
		t.Fatalf("%v", err)
	}
	instance := mathjax()

	html := instance.Render("a^2 + b^2 = c^2", "display")
	if len(html) == 0 {
		t.Fatalf("render output should not be empty")
	}

	css := instance.Css()
	if len(css) == 0 {
		t.Fatalf("css output should not be empty")
	}
}
