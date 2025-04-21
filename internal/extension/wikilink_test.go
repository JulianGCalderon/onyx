package extension_test

import (
	"juliangcalderon/onyx/internal/extension"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"
)

func TestWikilink(t *testing.T) {
	cwd := "dir"
	files := map[string]struct{}{
		"file1.md":         {},
		"dir/file2.md":     {},
		"dir/file3.md":     {},
		"dir/sub/file4.md": {},
	}

	m := goldmark.New(
		goldmark.WithExtensions(
			extension.NewWikilink(cwd, files),
		),
	)
	testutil.DoTestCases(m, []testutil.MarkdownTestCase{
		{
			No:          1,
			Description: "Relative Wikilink",
			Markdown:    `[[file3]]`,
			Expected:    `<p><a href="file3" title="file3">file3</a></p>`,
		},
		{
			No:          2,
			Description: "No Wikilink",
			Markdown:    `[`,
			Expected:    `<p>[</p>`,
		},
		{
			No:          3,
			Description: "No Wikilink",
			Markdown:    `]`,
			Expected:    `<p>]</p>`,
		},
		{
			No:          4,
			Description: "No Wikilink",
			Markdown:    `[]`,
			Expected:    `<p>[]</p>`,
		},
		{
			No:          5,
			Description: "No Wikilink",
			Markdown:    `[[]`,
			Expected:    `<p>[[]</p>`,
		},
		{
			No:          6,
			Description: "No Wikilink",
			Markdown:    `[]]`,
			Expected:    `<p>[]]</p>`,
		},
		{
			No:          7,
			Description: "No Wikilink",
			Markdown:    `[[]]`,
			Expected:    `<p>[[]]</p>`,
		},
	}, t)
}
