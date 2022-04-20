package feed

import (
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/hblanks/speakwrite/internal/content"
)

const maxItems = 50

func toTitle(series *content.Series, postTitle string) string {
	if series.Name == "" || series.Title == "" {
		return postTitle
	}
	return series.Title + ": " + postTitle
}

func toDescription(m *content.PostMetadata) string {
	tags := ""
	if len(m.Tags) > 0 {
		tags = "[" + strings.Join(m.Tags, ", ") + "]"
	}
	if m.Deck == "" {
		return tags
	}
	if tags == "" {
		return m.Deck
	}
	return m.Deck + " " + tags
}

func join(u url.URL, relPath string) string {
	u.Path = path.Join(u.Path, relPath) // Safe b/c u is not a pointer.
	return u.String()
}

// Takes the base (unnamed) series and an ordered slice of posts.
// Constructs an RSS feed and writes the most recent posts out
// to it.
func WriteRSS(publicURL *url.URL, md *content.SeriesMetadata, posts []*content.Post, w io.Writer) error {
	feed := &feeds.Feed{
		Title:       md.Title,
		Link:        &feeds.Link{Href: publicURL.String()},
		Description: md.Description,
		Author: &feeds.Author{
			Name:  md.Author.Name,
			Email: md.Author.Email,
		},
		Created: md.Created,
	}

	feed.Items = make([]*feeds.Item, 0, len(posts))
	for i, p := range posts {
		if i > maxItems {
			break
		}
		item := &feeds.Item{
			Title:       toTitle(p.Series, p.Title),
			Link:        &feeds.Link{Href: join(*publicURL, p.RelativeURL())},
			Description: toDescription(&p.Metadata),
			Created:     p.Date,
		}
		feed.Items = append(feed.Items, item)
	}

	return feed.WriteRss(w)
}
