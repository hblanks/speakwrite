package web

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/hblanks/confint/internal/content"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) postURL(post *content.Post) string {
	var u url.URL = *s.PublicURL
	u.Path = path.Join("/posts", post.Name)
	return u.String()
}

func (s *Server) getRoot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	log.Printf("getRoot: name=%q", name)

	data := struct {
		Title       string
		LatestURL   string
		LatestTitle string
	}{Title: "The static redirect"}

	if post := s.Posts.GetLatest(); post != nil {
		data.LatestURL = s.postURL(post)
		data.LatestTitle = post.Title
	}

	t := s.GetTemplate(w, "root.html")
	if t == nil {
		return
	}
	if err := t.Execute(w, &data); err != nil {
		log.Printf("getRoot: name=%q error %v", err)
	}
}

func (s *Server) getPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	post := s.Posts.Get(name)
	if post == nil {
		http.NotFound(w, r)
		return
	}
	log.Printf("getPost: name=%q post=%v", name, post)

	t := s.GetTemplate(w, "post.html")
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
	if err := t.Execute(w, &data); err != nil {
		log.Printf("getPost: name=%q error %v", err)
	}
}

func (s *Server) addHandlers(staticDir string) {
	s.router.GET("/", s.getRoot)
	// s.router.GET("/about", s.getRoot)
	s.router.GET("/posts/:name", s.getPost)
	s.router.ServeFiles("/static/*filepath", http.Dir(staticDir))
}
