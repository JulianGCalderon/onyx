package utils

import (
	"path/filepath"
	"strings"
)

func Must[T any](obj T, err error) T {
	AssertNil(err)
	return obj
}
func AssertNil(err error) {
	if err != nil {
		panic(err)
	}
}

func IsMarkdown(path string) bool {
	return filepath.Ext(path) == ".md"
}

func SetExt(path, ext string) string {
	return strings.TrimSuffix(path, filepath.Ext(path)) + ext
}
