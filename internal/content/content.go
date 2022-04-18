package content

// Reads and parses post content, notably Markdown.

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	_ "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Tweak Markdown rendering so that:
//	- the pandoc-style title block is not rendered
//	- comments are not included
func nodeHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch v := node.(type) {
	case *ast.Heading:
		if v.IsTitleblock {
			return ast.SkipChildren, true
		}
	case *ast.HTMLBlock:
		if bytes.HasPrefix(v.Literal, []byte("<!--")) {
			return ast.SkipChildren, true
		}
	case *ast.HTMLSpan:
		if bytes.HasPrefix(v.Literal, []byte("<!--")) {
			return ast.SkipChildren, true
		}
	}
	return ast.GoToNext, false
}

var mdrenderer = html.NewRenderer(html.RendererOptions{
	Title:          "A custom title",
	Flags:          html.CommonFlags, // | html.FootnoteReturnLinks,
	RenderNodeHook: html.RenderNodeFunc(nodeHook),
})

const IsoDateFormat = "2006-01-02"

// Parses a markdown file at a given path.
func parseMarkdown(path string) (ast.Node, error) {
	mdparser := parser.NewWithExtensions(
		parser.CommonExtensions | parser.Footnotes |
			parser.MathJax | parser.AutoHeadingIDs | parser.Titleblock,
	)

	// log.Printf("content.parse: %s", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	md, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return mdparser.Parse(md), nil
}

// Walks the AST and returns the title.
func getTitle(doc ast.Node) string {
	var title string
	var inTitle bool
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if node, ok := node.(*ast.Heading); ok && node.IsTitleblock {
			if entering {
				inTitle = true
				return ast.GoToNext
			}
		}
		if inTitle {
			if node, ok := node.(*ast.Text); ok {
				title = string(node.Literal)
				return ast.Terminate
			}
		}
		return ast.GoToNext
	})
	return title
}
