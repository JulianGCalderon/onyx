package wikilink_test

import (
	"juliangcalderon/onyx/internal/wikilink"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"
)

func TestWikilink(t *testing.T) {
	cwd := "dir"
	files := map[string]struct{}{
		"file1.md":         {},
		"file1.png":        {},
		"dir/file2.md":     {},
		"dir/file3.md":     {},
		"dir/sub/file4.md": {},
	}

	m := goldmark.New(
		goldmark.WithExtensions(
			wikilink.NewWikilink(cwd, files),
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
		{
			No:          8,
			Description: "Absolute Wikilink",
			Markdown:    `[[file1]]`,
			Expected:    `<p><a href="../file1" title="file1">file1</a></p>`,
		},
		{
			No:          9,
			Description: "Unique Wikilink",
			Markdown:    `[[file4]]`,
			Expected:    `<p><a href="sub/file4" title="file4">file4</a></p>`,
		},
		{
			No:          10,
			Description: "Relative Wikilink in Subdir",
			Markdown:    `[[sub/file4]]`,
			Expected:    `<p><a href="sub/file4" title="sub/file4">sub/file4</a></p>`,
		},
		{
			No:          11,
			Description: "Wikilink with Title",
			Markdown:    `[[sub/file4|custom title]]`,
			Expected:    `<p><a href="sub/file4" title="custom title">custom title</a></p>`,
		},
		{
			No:          12,
			Description: "Wikilink with Hash",
			Markdown:    `[[sub/file4#heading1]]`,
			Expected:    `<p><a href="sub/file4#heading1" title="sub/file4#heading1">sub/file4#heading1</a></p>`,
		},
		{
			No:          13,
			Description: "Wikilink with Hash and Title",
			Markdown:    `[[sub/file4#heading1|custom title]]`,
			Expected:    `<p><a href="sub/file4#heading1" title="custom title">custom title</a></p>`,
		},
		{
			No:          14,
			Description: "Wikilink with non markdown link",
			Markdown:    `[[file1.png]]`,
			Expected:    `<p><a href="../file1.png" title="file1.png">file1.png</a></p>`,
		},
		{
			No:          15,
			Description: "Wikilink with embed",
			Markdown:    `![[file1.png|title]]`,
			Expected:    `<p><img src="../file1.png" alt="title" title="title"></p>`,
		},
	}, t)
}
