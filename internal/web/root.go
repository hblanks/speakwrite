package web

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/hblanks/speakwrite/internal/content"
)

type RootData struct {
	BaseData
	Posts       *content.PostIndex
	LatestURL   string // *absolute* URL
	RelativeURL string
	Title       string
}

// All prior posts, oldest ones last.
// func (rd *RootData) PriorPosts() []*content.Post {
// 	count := len(rd.Posts)
// 	switch count {
// 	case 0, 1:
// 		return nil
// 	default:
// 		result := make([]*content.Post, count-1, count-1)
// 		for i, p := range rd.Posts[0 : count-1] {
// 			result[count-2-i] = p
// 		}
// 		return result
// 	}
// }

// Serve GET / requests.
func (s *Server) getRoot(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := RootData{
		BaseData: BaseData{
			Now: time.Now(),
		},
		Posts:       s.Posts,
		RelativeURL: "/",
		Title:       "",
	}

	if post := s.Posts.GetLatest(); post != nil {
		data.LatestURL = post.RelativeURL()
	}

	t := s.GetTemplate(w, "root.html")
	if t == nil {
		return
	}
	if err := t.Execute(w, &data); err != nil {
		log.Printf("getRoot: error %v", err)
	}
}
