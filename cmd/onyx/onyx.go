package main

import (
	"bytes"
	"io"
	"io/fs"
	"juliangcalderon/onyx/internal"
	myExtension "juliangcalderon/onyx/internal/extension"
	"juliangcalderon/onyx/internal/utils"
	"os"
	"path/filepath"
	"text/template"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
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
		if d.IsDir() {
			return nil
		}

		path, err = filepath.Rel(ContentPath, path)
		utils.AssertNil(err)
		files = append(files, path)

		return nil
	})
	utils.AssertNil(err)

	filesSet := make(map[string]struct{})
	for _, file := range files {
		filesSet[file] = struct{}{}
	}

	templates, err := template.ParseGlob(filepath.Join(TemplatesPath, "*"))
	utils.AssertNil(err)

	for _, file := range files {
		srcPath := filepath.Join(ContentPath, file)
		dstPath := internal.BuildDstPath(file, PublicPath)
		fileDir := filepath.Dir(file)

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
		markdown := goldmark.New(

			goldmark.WithExtensions(extension.GFM, extension.DefinitionList, extension.Footnote, extension.Typographer, myExtension.NewWikilink(fileDir, filesSet), mathjax.MathJax),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
		)
		err = markdown.Convert(content, &html)
		utils.AssertNil(err)

		ctx := internal.PageContext{
			Dir:     fileDir,
			Content: html.String(),
		}

		err = templates.ExecuteTemplate(dstFile, "root.gotmpl", ctx)
		utils.AssertNil(err)
	}
}
