package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"

	"github.com/hblanks/confint/internal/content"
)

type Server struct {
	router     *httprouter.Router
	contentDir string

	posts     content.PostIndex
	templates map[string]*template.Template
}

func NewServer(contentDir, themeDir string) (*Server, error) {
	s := &Server{
		router:     httprouter.New(),
		contentDir: contentDir,
		posts:      make(content.PostIndex),
		templates:  make(map[string]*template.Template),
	}
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
	postIndex, err := content.ReadPosts(s.contentDir)
	if err != nil {
		return err
	}
	log.Printf("Server.loadContent: posts=%d", len(postIndex))
	s.posts = postIndex
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
		var t *template.Template
		if p == basePath {
			t, err = template.ParseFiles(p)
		} else {
			t, err = template.ParseFiles(p, basePath)
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

func (s *Server) getTemplate(w http.ResponseWriter, name string) *template.Template {
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
