package main

import (
	"fmt"
	"io/fs"
	"juliangcalderon/onyx/internal/utils"
	"path/filepath"
)

const ContentPath = "content"

func main() {
	files := make([]string, 0)

	err := filepath.WalkDir(ContentPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !utils.IsMarkdown(path) {
			return nil
		}

		path, err = filepath.Rel(ContentPath, path)
		utils.AssertNil(err)
		files = append(files, path)

		return nil
	})
	utils.AssertNil(err)

	for _, file := range files {
		fmt.Println(file)
	}
}
