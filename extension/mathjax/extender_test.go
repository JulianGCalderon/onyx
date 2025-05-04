package mathjax_test

import (
	"bytes"
	"juliangcalderon/onyx/extension/mathjax"
	"juliangcalderon/onyx/node"
	"testing"

	"github.com/dop251/goja"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"
)

func TestMathjax(t *testing.T) {
	mathjax, err := mathjax.NewIsolatedMathjax()
	if err != nil {
		t.Fatalf("failed to build mathjax extension: %v", err)
	}

	module, err := node.NewMathjax(goja.New())
	if err != nil {
		t.Fatalf("failed to build mathjax instance: %v", err)
	}
	instance := module()

	m := goldmark.New(
		goldmark.WithExtensions(
			mathjax,
		),
	)
	testutil.DoTestCases(m, []testutil.MarkdownTestCase{
		{
			No:          1,
			Description: "Inline Pitagoras",
			Markdown:    `$a^2 + b^2 + c^2$`,
			Expected:    "<p>" + instance.Render("a^2 + b^2 + c^2", "inline") + "</p>",
		},
		{
			No:          2,
			Description: "Block Pitagoras",
			Markdown: "$$\n" +
				"a^2 + b^2 + c^2\n" +
				"$$\n",
			Expected: instance.Render("a^2 + b^2 + c^2", "display"),
		},
	}, t)

	if t.Failed() {
		return
	}

	actualCSS := []byte(mathjax.CSS())
	expectedCSS := []byte(instance.Css())

	ok := bytes.Equal(actualCSS, expectedCSS)
	if !ok {
		t.Fatalf("CSS:\n%s\n", testutil.DiffPretty(expectedCSS, actualCSS))
	}
}
