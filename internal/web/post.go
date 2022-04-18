package web

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/hblanks/speakwrite/internal/content"
	"github.com/julienschmidt/httprouter"
)

type PostData struct {
	BaseData
	*content.Post
	Content template.HTML
}

func (s *Server) identifyPost(filepath string) (*content.Post, string) {
	filepath = strings.TrimPrefix(filepath, "/")
	part0, filepath, _:= strings.Cut(filepath, "/")
	var part1 string
	defer log.Printf("identifyPost: part0=%q part1=%q filepath=%q",
		part0, part1, filepath)
	
	// Posts without a series are most common.
	post := s.Posts.Get("", part0)
	if post != nil {
		return post, filepath
	}

	// Else, look for post in a series.
	part1, filepath, _ = strings.Cut(filepath, "/")
	post = s.Posts.Get(part0, part1)
	if post != nil {
		return post, filepath
	}

	return nil, ""
}

// Serve post and associated files.
func (s *Server) getPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	post, extra := s.identifyPost(ps.ByName("filepath"))
	if post == nil {
		http.NotFound(w, r)
		return
	}

	if extra != "" {
		log.Printf("getPost: post=%v filepath=%s", post.RelativeURL(), extra)
		fs := ContentDir(path.Dir(post.ContentPath))
		fileServer := http.FileServer(fs)
		r.URL.Path = extra
		fileServer.ServeHTTP(w, r)
		return
	}

	log.Printf("getPost: post=%v", post.RelativeURL())

	t := s.GetTemplate(w, "post.html")
	if t == nil {
		return
	}
	postContent, err := post.HTML()
	if err != nil {
		log.Printf("getPost: error %v", err)
		sendError(w, http.StatusInternalServerError)
	}

	data := PostData{
		BaseData: BaseData{
			Now: time.Now(),
		},
		Post:    post,
		Content: postContent,
	}
	if err := t.Execute(w, &data); err != nil {
		log.Printf("getPost: name=%s error %v", post.RelativeURL(), err)
	}
}
