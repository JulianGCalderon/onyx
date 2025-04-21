package main

import (
	"bytes"
	"io"
	"io/fs"
	"juliangcalderon/onyx/internal"
	"juliangcalderon/onyx/internal/utils"
	"os"
	"path/filepath"
	"text/template"

	"github.com/yuin/goldmark"
)

const ContentPath = "content"
const PublicPath = "public"
const TemplatesPath = "templates"

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

	templates, err := template.ParseGlob(filepath.Join(TemplatesPath, "*"))
	utils.AssertNil(err)

	for _, file := range files {
		srcPath := filepath.Join(ContentPath, file)
		dstPath := internal.BuildDstPath(file, PublicPath)

		content, err := os.ReadFile(srcPath)
		utils.AssertNil(err)

		dstParent := filepath.Dir(dstPath)
		err = os.MkdirAll(dstParent, 0o775)
		utils.AssertNil(err)

		dstFile, err := os.Create(dstPath)
		utils.AssertNil(err)
		defer dstFile.Close()

		if !utils.IsMarkdown(file) {
			io.Copy(dstFile, bytes.NewReader(content))
			continue
		}

		var html bytes.Buffer
		markdown := goldmark.New()
		err = markdown.Convert(content, &html)
		utils.AssertNil(err)

		ctx := internal.PageContext{
			Dir:     filepath.Dir(file),
			Content: html.String(),
		}

		err = templates.ExecuteTemplate(dstFile, "root.gotmpl", ctx)
		utils.AssertNil(err)
	}
}
