package mathjax

import (
	"testing"

	_ "embed"

	v8 "rogchap.com/v8go"
)

//go:embed javascript/bundle.js
var mathjax string

func TestInteroperate(t *testing.T) {
	javascript := v8.NewContext()
	javascript.RunScript(mathjax, "bundle.js")

	mathjax, err := javascript.Global().Get("mathjax")
	if err != nil {
		t.Fatalf("%v", mathjax)
	}

	math := `a^2 + b^2 = c^2 \impliedby \text{"Pitagoras"}`
	javascript.Global().Set("math", math)
	javascript.Global().Set("type", "display")

	rendered, err := javascript.RunScript(
		"mathjax.render(math, type)",
		"mathjax_test_2.go",
	)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if !rendered.IsString() {
		t.Fatalf("render output should be a strign: %v", rendered)
	}

	css, err := javascript.RunScript(
		"mathjax.css()",
		"mathjax_test_2.go",
	)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if !css.IsString() {
		t.Fatalf("css output should be a strign: %v", css)
	}
}
