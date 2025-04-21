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
