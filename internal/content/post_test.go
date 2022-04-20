package content

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func createContent(t *testing.T, series []*Series) string {
	contentRoot, err := os.MkdirTemp("", "speakwrite-content-test*")
	if err != nil {
		t.Fatalf("createContent error: %v", err)
		return ""
	}

	for _, s := range series {
		seriesPath := filepath.Join(contentRoot, "posts", s.Name)
		if err := os.MkdirAll(seriesPath, 0755); err != nil {
			t.Fatalf("createContent error: %v", err)
			return ""
		}

		if s.Title != "" {
			metadataPath := filepath.Join(seriesPath, "metadata.json")
			data, err := json.Marshal(&SeriesMetadata{Title: s.Title})
			if err != nil {
				t.Fatalf("createContent error: %v", err)
				return ""
			}
			if err := ioutil.WriteFile(metadataPath, data, 0644); err != nil {
				t.Fatalf("createContent error: %v", err)
				return ""
			}
		}

		for _, p := range s.Posts {
			postPath := filepath.Join(seriesPath, p.Date.Format("2006-01-02-")+p.Name)
			if err := os.Mkdir(postPath, 0755); err != nil {
				t.Fatalf("createContent error: %v", err)
				return ""
			}
			metadataPath := filepath.Join(postPath, "metadata.json")
			data, err := json.Marshal(&(p.Metadata))
			if err != nil {
				t.Fatalf("createContent error: %v", err)
				return ""
			}
			if err := ioutil.WriteFile(metadataPath, data, 0644); err != nil {
				t.Fatalf("createContent error: %v", err)
				return ""
			}

			// Create index.md with the the title. No other content
			// needed.
			contentPath := filepath.Join(postPath, "index.md")
			data = []byte("% " + p.Title + "\n")
			if err := ioutil.WriteFile(contentPath, data, 0644); err != nil {
				t.Fatalf("createContent error: %v", err)
				return ""
			}
		}
	}

	return contentRoot
}

func deleteContent(contentRoot string) {
	if contentRoot != "" { // Sanity check
		os.RemoveAll(contentRoot)
	}
}

var baseSeries = &Series{}
var seriesA = &Series{
	Name: "A",
	SeriesMetadata: SeriesMetadata{
		Title: "A was a series",
	},
}
var allSeries = []*Series{}

func init() {
	baseSeries.Posts = []*Post{
		&Post{
			Name:   "blah-blah-blah",
			Title:  "Blah Blah Blah",
			Date:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Series: baseSeries,
			Metadata: PostMetadata{
				Deck: "blabbbb",
				Tags: []string{"blah", "blab"},
			},
		},
	}
	allSeries = append(allSeries, baseSeries)

	seriesA.Posts = []*Post{
		&Post{
			Name:   "first post",
			Title:  "First Post",
			Date:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			Series: seriesA,
			Metadata: PostMetadata{
				Deck: "fooooo",
				Tags: []string{"foo", "fie"},
			},
		},
		&Post{
			Name:   "second post",
			Title:  "Second Post",
			Date:   time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			Series: seriesA,
			Metadata: PostMetadata{
				Deck: "waaaa",
				Tags: []string{"foo", "fie"},
			},
		},
	}
	allSeries = append(allSeries, seriesA)
}

func postsEqual(p1, p2 *Post) bool {
	return p1.Name == p2.Name &&
		p1.Title == p2.Title &&
		p1.Date.Equal(p2.Date) &&
		p1.Series.Name == p2.Series.Name &&
		reflect.DeepEqual(p1.Metadata, p2.Metadata)

}

// (Regrettably) tests everything in one spot.
func TestPostIndex(t *testing.T) {
	contentRoot := createContent(t, allSeries)
	defer deleteContent(contentRoot)

	pi, err := NewPostIndex(contentRoot)
	if err != nil {
		t.Fatalf("NewPostIndex failed: %v", err)
	}

	t.Run("GetLatest", func(t *testing.T) {
		expected := seriesA.Posts[1]
		pi.GetLatest()
		if actual := pi.GetLatest(); !postsEqual(expected, actual) {
			t.Errorf("expected %s != actual %s (%#v != %#v)",
				expected.Name, actual.Name, expected, actual)
		}
	})

	t.Run("GetPriorPosts: base series", func(t *testing.T) {
		expected := baseSeries.Posts[0]
		priors := pi.GetPriorPosts("")
		if len(priors) != 1 {
			t.Fatalf("returned %d posts, not 1", len(priors))
		}
		actual := priors[0]
		if !postsEqual(expected, actual) {
			t.Errorf("expected %s != actual %s", expected.Name, actual.Name)
		}
	})

	t.Run("GetPriorPosts: series A", func(t *testing.T) {
		expected := seriesA.Posts[0]
		priors := pi.GetPriorPosts("A")
		if len(priors) != 1 {
			t.Fatalf("returned %d posts, not 1", len(priors))
		}
		actual := priors[0]
		if !postsEqual(expected, actual) {
			t.Errorf("expected %s != actual %s", expected.Name, actual.Name)
		}
	})
}
