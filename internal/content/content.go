package content

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gomarkdown/markdown"
	_ "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

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
	Flags:          html.CommonFlags | html.FootnoteReturnLinks,
	RenderNodeHook: html.RenderNodeFunc(nodeHook),
})

var postRegexp = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})-(.*)`)

type Page struct {
	Title       string
	Name        string
	contentPath string
}

type Post struct {
	Date        time.Time
	Title       string
	Name        string
	contentPath string
}

const IsoDateFormat = "2006-01-02"

func parse(path string) (ast.Node, error) {
	mdparser := parser.NewWithExtensions(
		parser.CommonExtensions | parser.AutoHeadingIDs | parser.Titleblock,
	)

	log.Printf("content.parse: %s", path)
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

func NewPost(dateStr, name, contentPath string) (*Post, error) {
	t, err := time.Parse(IsoDateFormat, dateStr)
	if err != nil {
		return nil, err
	}

	doc, err := parse(contentPath)
	if err != nil {
		return nil, err
	}
	title := getTitle(doc)
	if title == "" {
		return nil, errors.New("No title found")
	}

	return &Post{
		Date:        t,
		Name:        name,
		contentPath: contentPath,
		Title:       title,
	}, nil
}

func (p *Post) HTML() (template.HTML, error) {
	doc, err := parse(p.contentPath)
	if err != nil {
		return template.HTML(""), err
	}
	if doc == nil {
		panic("wat")
	}
	ast.Print(os.Stderr, doc)
	output := markdown.Render(doc, mdrenderer)
	if len(output) == 0 {
		return template.HTML(""), errors.New("Failed to render document")
	}
	return template.HTML(output), nil
}

type PageIndex map[string]*Page
type PostIndex map[string]*Post

func ReadPosts(contentDir string) (PostIndex, error) {
	d, err := os.Open(filepath.Join(contentDir, "posts"))
	if err != nil {
		return nil, err
	}
	defer d.Close()

	infos, err := d.Readdir(-1)
	if err != nil {
		// NB: IO errors are ignored by Glob!
		return nil, err
	}
	result := make(PostIndex)
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		contentPath := filepath.Join(d.Name(), info.Name(), "index.md")
		if _, err := os.Stat(contentPath); err != nil {
			log.Printf("Failed to stat %s: %v", contentPath, err)
			return nil, err
		}

		basename := filepath.Base(info.Name())
		if m := postRegexp.FindStringSubmatch(basename); m != nil {

			post, err := NewPost(m[1], m[2], contentPath)
			if err != nil {
				return nil, err
			}
			result[post.Name] = post
		}
	}
	return result, nil
}

func ReadPages(contentDir string) (PageIndex, error) {
	return nil, nil
}
