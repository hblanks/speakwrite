package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/hblanks/speakwrite/internal/feed"
)

func (s *Server) getRSS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	series := s.Posts.GetBaseSeries()
	if series == nil {
		log.Printf("getRSS: base series not found")
		sendError(w, http.StatusNotFound)
		return
	}

	err := feed.WriteRSS(&series.SeriesMetadata, s.Posts.Posts, w)
	if err != nil {
		sendError(w, http.StatusInternalServerError)
	}
}
