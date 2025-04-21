package extension

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
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
	reTarget := `([^|#]*)`
	reHash := `(#[^|]+)?`
	reTitle := `(\|.+)?`
	reFull := fmt.Sprintf(`^\[\[%v%v%v\]\]`, reTarget, reHash, reTitle)

	_, segment := block.Position()
	re := regexp.MustCompile(reFull)
	matches := block.FindSubMatch(re)
	if len(matches) != 4 {
		return nil
	}
	target := matches[1]
	hash := matches[2]
	title := matches[3]

	if len(title) == 0 {
		segment.Start += len("[[")
		title = append(target, hash...)
	} else {
		segment.Start += len("[[") + len(target) + len(hash) + len("|")
		title = title[1:]
	}
	segment.Stop = segment.Start + len(title)

	destination := []byte(p.resolveTarget(string(target)))
	destination = append(destination, hash...)

	link := ast.NewLink()
	link.Title = title
	link.Destination = destination
	link.AppendChild(link, ast.NewTextSegment(segment))

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
	if target == "" {
		return ""
	}
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
