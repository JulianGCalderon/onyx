package wikilink

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

type wikilink struct {
	cwd   string
	files map[string]struct{}
}

func NewWikilink(cwd string, files map[string]struct{}) goldmark.Extender {
	return &wikilink{
		cwd:   cwd,
		files: files,
	}
}

func (e *wikilink) Extend(m goldmark.Markdown) {
	resolver := resolver{
		cwd:   e.cwd,
		files: e.files,
	}

	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(parser.InlineParser(&wikilinkParser{resolver: resolver}), 199),
		),
	)
}
