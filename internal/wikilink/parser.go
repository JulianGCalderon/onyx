package wikilink

import (
	"fmt"
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type wikilinkParser struct {
	resolver resolver
}

func (w *wikilinkParser) Trigger() []byte {
	return []byte{'[', '!'}
}

func (p *wikilinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	reTarget := `([^|#]*?)`
	reHash := `(#[^|]+?)?`
	reTitle := `(\|.+?)?`
	reFull := fmt.Sprintf(`^!?\[\[%v%v%v\]\]`, reTarget, reHash, reTitle)

	_, segment := block.Position()
	re := regexp.MustCompile(reFull)
	matches := block.FindSubMatch(re)
	if len(matches) != 4 {
		return nil
	}
	target := matches[1]
	hash := matches[2]
	title := matches[3]

	if len(target) == 0 && len(hash) <= 1 && len(title) <= 1 {
		return nil
	}

	embed := matches[0][0] == '!'
	if embed {
		segment.Start += len("!")
	}

	if len(title) == 0 {
		segment.Start += len("[[")
		title = append(target, hash...)
	} else {
		segment.Start += len("[[") + len(target) + len(hash) + len("|")
		title = title[1:]
	}

	if len(hash) > 0 {
		hash = hash[1:]
	}

	segment.Stop = segment.Start + len(title)

	hash = []byte(normalizeFragment(string(hash)))

	destination := []byte(p.resolver.resolve(string(target), string(hash)))

	link := ast.NewLink()
	link.Title = title
	link.Destination = destination
	link.AppendChild(link, ast.NewTextSegment(segment))

	if embed {
		return ast.NewImage(link)
	} else {
		return link
	}
}
