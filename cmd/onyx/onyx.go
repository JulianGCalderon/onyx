package main

import (
	"bytes"
	"io"
	"io/fs"
	"juliangcalderon/onyx"
	"juliangcalderon/onyx/extension/mathjax"
	"juliangcalderon/onyx/extension/wikilink"
	"juliangcalderon/onyx/utils"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/anchor"
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

		source, err := os.ReadFile(srcPath)
		utils.AssertNil(err)

		dstParent := filepath.Dir(dstPath)
		err = os.MkdirAll(dstParent, 0o775)
		utils.AssertNil(err)

		dstFile, err := os.Create(dstPath)
		utils.AssertNil(err)
		defer dstFile.Close()

		if !utils.IsMarkdown(file) {
			io.Copy(dstFile, bytes.NewReader(source))
			continue
		}

		mathjax, err := mathjax.NewIsolatedMathjax()
		utils.AssertNil(err)

		var htmlContent bytes.Buffer
		markdown := goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.DefinitionList,
				extension.Footnote,
				extension.Typographer,
				mathjax,
				&anchor.Extender{
					Texter: anchor.Text("#"),
				},
				meta.New(meta.WithStoresInDocument()),
				wikilink.NewWikilink(fileDir, filesSet),
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		)

		document := markdown.Parser().Parse(text.NewReader(source))

		meta := document.OwnerDocument().Meta()

		if title, ok := meta["title"].(string); ok {
			newHeading := ast.NewHeading(1)
			newHeading.AppendChild(newHeading, ast.NewString([]byte(title)))
			document.InsertBefore(document, document.FirstChild(), newHeading)
		}
		if heading, ok := document.FirstChild().(*ast.Heading); !ok || heading.Level > 1 {
			title := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
			newHeading := ast.NewHeading(1)
			newHeading.AppendChild(newHeading, ast.NewString([]byte(title)))
			document.InsertBefore(document, document.FirstChild(), newHeading)
		}

		err = markdown.Renderer().Render(&htmlContent, source, document)
		utils.AssertNil(err)

		ctx := internal.PageContext{
			Dir:     fileDir,
			Content: htmlContent.String(),
			Style:   mathjax.CSS(),
		}

		err = templates.ExecuteTemplate(dstFile, "root.gotmpl", ctx)
		utils.AssertNil(err)
	}
}
