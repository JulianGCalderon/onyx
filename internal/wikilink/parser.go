package wikilink

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type wikilinkParser struct {
	resolver resolver
}

var (
	tokenEmbed    []byte = []byte("![[")
	tokenLink     []byte = []byte("[[")
	tokenClose    []byte = []byte("]]")
	tokenFragment []byte = []byte("#")
	tokenTitle    []byte = []byte("|")
)

func (w *wikilinkParser) Trigger() []byte {
	return []byte{'[', '!'}
}

func (p *wikilinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	source, _ := block.PeekLine()

	var (
		indexStart int
		indexClose int
		embed      bool
	)

	// Find start of wikilink
	switch {
	case bytes.HasPrefix(source, tokenEmbed):
		embed = true
		indexStart = len(tokenEmbed)
	case bytes.HasPrefix(source, tokenLink):
		embed = false
		indexStart = len(tokenLink)
	default:
		return nil
	}

	// Find end of wikilink
	indexClose = bytes.Index(source, tokenClose)
	if indexClose <= 0 {
		return nil
	}

	// Obtain full wikilink content
	content := source[indexStart:indexClose]
	if len(content) == 0 {
		// Content cannot be empty
		return nil
	}

	// Separate destination from title: [[destination|title]]
	destination, title, cut := bytes.Cut(content, tokenTitle)
	if cut && len(title) == 0 {
		// Title cannot be empty
		return nil
	}
	if !cut {
		// If there is no specific title, fall back to the full destination
		title = destination
	}

	// Separate target from fragment: [[target#fragment]]
	target, fragment, cut := bytes.Cut(destination, tokenFragment)
	if cut && len(fragment) == 0 {
		// Fragment cannot be empty
		return nil
	}

	// Resolve destination using the resolver
	destination = []byte(p.resolver.resolve(string(target), string(fragment)))

	link := ast.NewLink()
	link.Title = title
	link.Destination = destination
	link.AppendChild(link, ast.NewString(title))

	block.Advance(indexClose + len(tokenClose))

	if embed {
		return ast.NewImage(link)
	} else {
		return link
	}
}
