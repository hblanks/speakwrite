package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"

	"github.com/hblanks/confint/internal/content"
)

func (s *Server) getRoot(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// name := ps.ByName("name")
	// log.Printf("getRoot: name=%q", name)

	data := struct {
		Title       string
		LatestURL   string
		LatestTitle string
	}{Title: "The static redirect"}

	if post := s.Posts.GetLatest(); post != nil {
		data.LatestURL = postRelativeURL(post)
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

// Serve post and associated files.
func (s *Server) getPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	post := s.Posts.Get(name)
	if post == nil {
		http.NotFound(w, r)
		return
	}

	if fpath := ps.ByName("filepath"); fpath != "/" {
		log.Printf("getPost: name=%q post=%v filepath=%s", name, post, fpath)
		fs := ContentDir(filepath.Dir(post.ContentPath))
		fileServer := http.FileServer(fs)
		r.URL.Path = fpath
		fileServer.ServeHTTP(w, r)
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
