package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/hblanks/confint/internal/content"
)

type Server struct {
	PublicURL *url.URL

	router     *httprouter.Router
	contentDir string

	Posts     *content.PostIndex
	templates map[string]*template.Template
}

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
		return nil, err
	}
	if err := s.loadContent(contentDir); err != nil {
		return nil, err
	}

	staticDir := filepath.Join(themeDir, "static")
	if _, err := ioutil.ReadDir(staticDir); err != nil {
		return nil, err
	}
	s.addHandlers(staticDir)
	return s, nil
}

func (s *Server) loadContent(contentDir string) error {
	postIndex, err := content.LoadPosts(s.contentDir)
	if err != nil {
		return err
	}
	log.Printf("Server.loadContent: posts=%d", len(postIndex.Posts))
	s.Posts = postIndex
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
				return t.Format("January 2nd, 2006")
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
