package extension

import (
	"slices"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(newWikilinkParser(e.cwd, e.files), 199),
		),
	)
}

type wikilinkParser struct {
	cwd   string
	files map[string]struct{}
}

func newWikilinkParser(cwd string, files map[string]struct{}) parser.InlineParser {
	return &wikilinkParser{
		cwd:   cwd,
		files: files,
	}
}

func (w *wikilinkParser) Trigger() []byte {
	return []byte{'['}
}

func (w *wikilinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()

	// wikilink should start with '[['
	if len(line) < 2 {
		return nil
	}
	if !slices.Equal(line[:2], []byte{'[', '['}) {
		return nil
	}

	// wikilink should end with ']]'
	closingIndex := slices.Index(line, ']')
	if closingIndex < 0 || len(line) < closingIndex+1 {
		return nil
	}
	if !slices.Equal(line[closingIndex:closingIndex+2], []byte{']', ']'}) {
		return nil
	}

	// wikilink should not be empty
	line = line[2:closingIndex]
	if len(line) == 0 {
		return nil
	}

	block.Advance(closingIndex + 2)

	segment.Start += 2
	segment.Stop = segment.Start + len(line)

	link := ast.NewLink()
	link.Title = line
	link.Destination = line
	link.AppendChild(link, ast.NewTextSegment(segment))

	return link
}
