package wikilink

import (
	"log"
	"path/filepath"
	"strings"
)

type resolver struct {
	cwd   string
	files map[string]struct{}
}

func (r resolver) resolve(target, fragment string) (destination string) {
	destination = r.resolveTarget(target)
	if len(fragment) > 0 {
		destination += "#" + normalizeFragment(fragment)
	}
	return
}

// Resolves the wikilink target returning a path relative to the current directory.
// If the target has no extension, it assumes that it has extension ".md". If
// the target file is a markdown file, the final path will have no extension.
//
// The note resolution has the following precedence rules:
//
// 1. Absolute path from vault root.
// 2. Relative path from current directory.
// 3. Target note has an unique basename.
func (r resolver) resolveTarget(target string) string {
	if target == "" {
		return ""
	}
	if filepath.Ext(target) == "" {
		target += ".md"
	}

	// resolve absolute target
	_, ok := r.files[target]
	if ok {
		return r.buildPath(target)
	}

	// resolve relative target
	targetAsRelative := filepath.Join(r.cwd, target)
	_, ok = r.files[targetAsRelative]
	if ok {
		return r.buildPath(targetAsRelative)
	}

	// resolve path to unique note
	for note := range r.files {
		if filepath.Base(note) == target {
			return r.buildPath(note)
		}
	}

	return r.buildPath(target)
}

// Builds a path to `destination`, relative to the current directory.
//
// If the target is a markdown file, then the final path will have no extension.
func (r resolver) buildPath(destination string) string {
	destination, err := filepath.Rel(r.cwd, destination)
	if err != nil {
		log.Panicf("%v and %v should share root", r.cwd, destination)
	}

	if filepath.Ext(destination) == ".md" {
		destination = strings.TrimSuffix(destination, ".md")
	}

	return destination

}

func normalizeFragment(hash string) string {
	hash = strings.ToLower(hash)
	hash = strings.ReplaceAll(hash, " ", "-")
	return hash
}
