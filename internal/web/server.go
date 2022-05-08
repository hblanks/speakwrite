package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/hblanks/speakwrite/internal/content"
)

type Server struct {
	PublicURL *url.URL

	router *httprouter.Router

	contentDir string
	Posts      *content.PostIndex
	// Pages      content.PageIndex

	templates map[string]*template.Template

	staticDir string
}

// Creates (but does not run!) a server. Steps include:
//	- load all templates
//	- load all content
//  - set up all routes
func NewServer(publicURL, contentDir, themeDir string) (*Server, error) {
	s := &Server{
		router:     httprouter.New(),
		contentDir: contentDir,
		templates:  make(map[string]*template.Template),
	}

	u, err := url.Parse(publicURL)
	if err != nil {
		return nil, err
	}
	s.PublicURL = u

	if err := s.loadTemplates(filepath.Join(themeDir, "templates")); err != nil {
		return nil, fmt.Errorf("loadTemplates error: %w", err)
	}
	if err := s.loadContent(contentDir); err != nil {
		return nil, fmt.Errorf("loadContent error: %w", err)
	}

	s.staticDir = filepath.Join(themeDir, "static")
	if _, err := ioutil.ReadDir(s.staticDir); err != nil {
		return nil, err
	}
	s.addHandlers()
	return s, nil
}

func (s *Server) loadContent(contentDir string) error {
	postIndex, err := content.NewPostIndex(s.contentDir)
	if err != nil {
		return err
	}
	s.Posts = postIndex

	// pageIndex, err := content.LoadPages(s.contentDir)
	// if err != nil {
	// 	return err
	// }
	// s.Pages = pageIndex

	log.Printf("Server.loadContent: posts=%d", len(postIndex.Posts))
	return nil
}

func (s *Server) loadTemplates(templatesDir string) error {
	paths, err := filepath.Glob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return fmt.Errorf("failed to read templates dir %s: %w",
			templatesDir, err)
	}
	basePath := filepath.Join(templatesDir, "base.html")
	for _, p := range paths {
		t := template.New("base.html")
		t.Funcs(map[string]interface{}{
			"isoDate": func(t *time.Time) string {
				return t.Format("2006-01-02")
			},
			"englishDate": func(t *time.Time) string {
				return t.Format("January 2, 2006")
			},
		})
		if p == basePath {
			_, err = t.ParseFiles(p)
		} else {
			_, err = t.ParseFiles(basePath, p)
		}
		if err != nil {
			return err
		}
		s.templates[filepath.Base(p)] = t
	}
	return nil
}

func sendError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (s *Server) GetTemplate(w http.ResponseWriter, name string) *template.Template {
	t := s.templates[name]
	if t == nil {
		log.Printf("getTemplate: %s not found!", name)
		sendError(w, http.StatusInternalServerError)
	}
	return t
}

func joinURL(base *url.URL, relpath string) string {
	var u url.URL = *base
	u.Path = path.Join(u.Path, relpath)
	return u.String()
}

func (s *Server) staticURL(fpath string) string {
	rel, err := filepath.Rel(s.staticDir, fpath)
	if err != nil {
		panic(err)
	}
	return joinURL(s.PublicURL, path.Join("/static", rel))
}

func (s *Server) postURL(post *content.Post) string {
	return joinURL(s.PublicURL, post.RelativeURL()) + "/"
}

func (s *Server) addHandlers() {
	s.router.GET("/", s.getRoot)
	s.router.GET("/rss.xml", s.getRSS)
	s.router.GET("/posts/*filepath", s.getPost)
	s.router.ServeFiles("/static/*filepath", http.Dir(s.staticDir))
}

func (s *Server) GetURLs() ([]string, error) {
	urls := make([]string, 0)

	publicURL := s.PublicURL.String()
	if !strings.HasSuffix(publicURL, "/") {
		publicURL += "/"
	}
	urls = append(urls, publicURL)

	// Find all posts and related files
	for _, post := range s.Posts.Posts {
		u := s.postURL(post)
		urls = append(urls, u)
		postDir := filepath.Dir(post.ContentPath)
		err := filepath.Walk(postDir,
			func(path string, info os.FileInfo, err error) error {
				switch {
				case err != nil:
					return err
				case info.IsDir() && filepath.Base(path) == "exclude":
					return filepath.SkipDir
				case info.IsDir():
					return nil
				case path == post.ContentPath:
					return nil
				case path == filepath.Join(postDir, "metadata.json"):
					return nil
				}
				rel, err := filepath.Rel(postDir, path)

				parsedURL, err := url.Parse(u)
				if err != nil {
					return err
				}
				urls = append(urls, joinURL(parsedURL, rel))
				return nil
			})
		if err != nil {
			return nil, err
		}
	}

	// Add RSS feed.
	urls = append(urls, publicURL+"rss.xml")

	// Find all static assets
	err := filepath.Walk(s.staticDir,
		func(path string, info os.FileInfo, err error) error {
			switch {
			case err != nil:
				return err
			case info.IsDir():
				return nil
			}
			urls = append(urls, s.staticURL(path))
			return nil
		})
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
