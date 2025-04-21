package extension

import (
	"log"
	"path/filepath"
	"slices"
	"strings"

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

func (p *wikilinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()

	// wikilink should start with '[['
	if len(line) < 5 {
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

	block.Advance(closingIndex + 2)
	var (
		wikilinkStart int = 2
		wikilinkEnd   int = closingIndex
	)

	var (
		titleStart int
		targetEnd  int
	)

	verticalIndex := slices.Index(line, '|')
	if verticalIndex >= 0 {
		targetEnd = verticalIndex
		titleStart = verticalIndex + 1
	} else {
		targetEnd = wikilinkEnd
		titleStart = wikilinkStart
	}

	title := line[titleStart:wikilinkEnd]
	target := line[wikilinkStart:targetEnd]
	text := text.Segment{
		Start:        segment.Start + titleStart,
		Stop:         segment.Start + wikilinkEnd,
		Padding:      segment.Padding,
		ForceNewline: segment.ForceNewline,
	}

	destination := []byte(p.resolveTarget(string(target)))

	link := ast.NewLink()
	link.Title = title
	link.Destination = destination
	link.AppendChild(link, ast.NewTextSegment(text))

	return link
}

// Resolves the wikilink target returning a path relative to the current directory.
//
// If the target has no extension, it assumes that it has extension ".md". If
// the target file is a markdown file, the final path will have no extension.
//
// The note resolution has the following precedence rules:
//
// 1. Absolute path from vault root.
// 2. Relative path from current directory.
// 3. Target note has an unique basename.
func (p *wikilinkParser) resolveTarget(target string) string {
	if filepath.Ext(target) == "" {
		target += ".md"
	}

	// resolve absolute target
	_, ok := p.files[target]
	if ok {
		return p.buildPath(target)
	}

	// resolve relative target
	targetAsRelative := filepath.Join(p.cwd, target)
	_, ok = p.files[targetAsRelative]
	if ok {
		return p.buildPath(targetAsRelative)
	}

	// resolve path to unique note
	for note := range p.files {
		if filepath.Base(note) == target {
			return p.buildPath(note)
		}
	}

	return p.buildPath(target)
}

// Builds a path to `destination`, relative to the current directory.
//
// If the target is a markdown file, then the final path will have no extension.
func (p *wikilinkParser) buildPath(destination string) string {
	destination, err := filepath.Rel(p.cwd, destination)
	if err != nil {
		log.Panicf("%v and %v should share root", p.cwd, destination)
	}

	if filepath.Ext(destination) == ".md" {
		destination = strings.TrimSuffix(destination, ".md")
	}

	return destination

}
