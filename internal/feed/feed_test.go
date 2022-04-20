package feed

import (
	"bytes"
	"testing"
	"time"

	"github.com/hblanks/speakwrite/internal/content"
)

func TestWriteRSS(t *testing.T) {
	series := &content.Series{
		Name: "",
		SeriesMetadata: content.SeriesMetadata{
			Title:       "Here's a series",
			Description: "It's a thing",
			URL:         "http://speakwrite.blog/",
			Author: content.SeriesAuthor{
				Name:  "Pierce Inverarity",
				Email: "pierce@waste.org",
			},
		},
	}

	posts := []*content.Post{
		&content.Post{
			Name:  "first-post",
			Title: "First post",
			Metadata: content.PostMetadata{
				Deck: "yup",
			},
			Date:   time.Time{},
			Series: series,
		},
		&content.Post{
			Name:  "second-post",
			Title: "Second post",
			Metadata: content.PostMetadata{
				Deck: "yar",
			},
			Date:   time.Time{},
			Series: series,
		},
	}

	buf := &bytes.Buffer{}
	err := WriteRSS(&series.SeriesMetadata, posts, buf)
	if err != nil {
		t.Fatalf("WriteRSS() returned unexpected error: %v", err)
	}
	if len(buf.Bytes()) == 0 {
		t.Errorf("WriteRSS() wrote no bytes")
	}
}
