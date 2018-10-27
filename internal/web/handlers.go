package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/hblanks/confint/internal/content"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) getRoot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	log.Printf("getRoot: name=%q", name)
}

func (s *Server) getPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	post := s.posts[name]
	log.Printf("getPost: name=%q post=%v", name, post)

	t := s.getTemplate(w, "post.html")
	if t == nil {
		return
	}
	postContent, err := post.HTML()
	if err != nil {
		log.Printf("getPost: error %v", err)
		sendError(w, http.StatusInternalServerError)
	}

	data := struct {
		*content.Post
		Content template.HTML
	}{
		Post:    post,
		Content: postContent,
	}
	err = t.Execute(w, &data)
	if err != nil {
		log.Printf("getPost: name=%q error %v", err)
	}
}

func (s *Server) addHandlers(staticDir string) {
	s.router.GET("/", s.getRoot)
	// s.router.GET("/about", s.getRoot)
	s.router.GET("/posts/:name", s.getPost)
	s.router.ServeFiles("/static/*filepath", http.Dir(staticDir))
}
