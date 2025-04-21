package internal

import (
	"juliangcalderon/onyx/internal/utils"
	"path/filepath"
)

func BuildDstPath(path, root string) string {
	if utils.IsMarkdown(path) {
		path = utils.SetExt(path, ".html")
	}

	return filepath.Join(root, path)
}

type PageContext struct {
	Dir     string
	Content string
}

func (c PageContext) Root() string {
	return utils.Must(filepath.Rel(c.Dir, "."))
}
