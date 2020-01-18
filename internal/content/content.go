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
	"sort"
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
	Flags:          html.CommonFlags, // | html.FootnoteReturnLinks,
	RenderNodeHook: html.RenderNodeFunc(nodeHook),
})

const IsoDateFormat = "2006-01-02"

func parse(path string) (ast.Node, error) {
	mdparser := parser.NewWithExtensions(
		parser.CommonExtensions | parser.Footnotes |
			parser.MathJax | parser.AutoHeadingIDs | parser.Titleblock,
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

//
// Page
//
//
// type Page struct {
// 	Title       string
// 	Name        string
// 	contentPath string
// }
//
//
// PageIndex
//
//
// type PageIndex map[string]*Page
//
// func (p PageIndex) Routes() ([]string, error) {
// 	return []string{}, nil
// }
//
// func LoadPages(contentDir string) (PageIndex, error) {
// 	root := filepath.Join(contentDir, "pages")
// 	paths := make([]string, 0)
// 	filepath.Walk(root,
// 		func(path string, info os.FileInfo, err error) error {
// 			if err != nil || info.IsDir() {
// 				return nil
// 			}
// 			if strings.HasSuffix(path, ".md") {
// 				paths = append(paths, path)
// 			}
// 			return nil
// 		})
//
// 	log.Printf("LoadPages: paths=%v", paths)
//
// 	return nil, nil
// }

//
// Post
//

type Post struct {
	Date        time.Time
	Title       string
	Name        string
	ContentPath string
	index       int
}

func NewPost(dateStr, name, ContentPath string) (*Post, error) {
	t, err := time.Parse(IsoDateFormat, dateStr)
	if err != nil {
		return nil, err
	}

	doc, err := parse(ContentPath)
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
		ContentPath: ContentPath,
		Title:       title,
	}, nil
}

func (p *Post) HTML() (template.HTML, error) {
	doc, err := parse(p.ContentPath)
	if err != nil {
		return template.HTML(""), err
	}
	if doc == nil {
		panic("wat")
	}
	output := markdown.Render(doc, mdrenderer)
	if len(output) == 0 {
		return template.HTML(""), errors.New("Failed to render document")
	}
	return template.HTML(output), nil
}

//
// PostIndex
//

type PostIndex struct {
	names map[string]*Post
	Posts []*Post
}

func (p *PostIndex) Get(name string) *Post {
	return p.names[name]
}

func (p *PostIndex) GetLatest() *Post {
	if len(p.Posts) == 0 {
		return nil
	}
	return p.Posts[len(p.Posts)-1]
}

var postRegexp = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})-(.*)`)

func readPosts(contentDir string) ([]*Post, error) {
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

	posts := make([]*Post, 0)
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
			posts = append(posts, post)
		}
	}
	return posts, nil
}

// Loads posts from a directory into a PostIndex.
func LoadPosts(contentDir string) (*PostIndex, error) {
	posts, err := readPosts(contentDir)
	if err != nil {
		return nil, err
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Name < posts[j].Name
	})

	result := &PostIndex{
		names: make(map[string]*Post),
		Posts: posts,
	}
	for i, post := range posts {
		post.index = i
		result.names[post.Name] = post
	}
	return result, nil
}
