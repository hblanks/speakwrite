package content

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/gomarkdown/markdown"
)

//
// Extra metadata about a post.
//
type PostMetadata struct {
	Tags []string `json:"tags"`
	Deck string   `json:"deck"` // the "deck" or "drop line" of the post
}

//
// A Post is a dated page.
//
type Post struct {
	Date        time.Time
	ContentPath string
	Metadata    PostMetadata
	Name        string
	Series *Series
	Title       string
}

func NewPost(dateStr, name, contentPath, metadataPath string, series *Series) (*Post, error) {
	t, err := time.Parse(IsoDateFormat, dateStr)
	if err != nil {
		return nil, err
	}

	// Parse content, but only so we can get the title.
	doc, err := parseMarkdown(contentPath)
	if err != nil {
		return nil, err
	}
	title := getTitle(doc)
	if title == "" {
		return nil, fmt.Errorf("No title found for %s", contentPath)
	}

	post := &Post{
		Date:        t,
		Name:        name,
		ContentPath: contentPath,
		Title:       title,
		Series:      series,
	}

	if metadataPath != "" {
		b, err := ioutil.ReadFile(metadataPath)
		if err != nil {
			return nil, fmt.Errorf("NewPost error: %w", err)
		}
		err = json.Unmarshal(b, &post.Metadata)
		if err != nil {
			return nil, err
			return nil, fmt.Errorf("NewPost error: %w", err)
		}
	}

	return post, nil
}

func (p *Post) HTML() (template.HTML, error) {
	doc, err := parseMarkdown(p.ContentPath)
	if err != nil {
		return template.HTML(""), err
	}
	if doc == nil {
		panic("wat")
	}
	output := markdown.Render(doc, mdrenderer)
	if len(output) == 0 {
		return template.HTML(""), errors.New("Failed to render document")
	}
	return template.HTML(output), nil
}

func (p *Post) RelativeURL() string {
	if p.Series.Name == "" {
		return path.Join("/posts", p.Name) + "/"
	} else {
		return path.Join("/posts", p.Series.Name, p.Name) + "/"
	}
}

// Sorts a slice of posts by descending (date, name).
func sortPosts(posts []*Post) {
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Date == posts[j].Date {
			return posts[i].Name > posts[j].Name
		}
		return posts[i].Date.After(posts[j].Date)
	})
}

//
// Series
//

type SeriesMetadata struct {
	Title string `json:"title"`
}

type Series struct {
	Name  string
	Title string
	Posts []*Post
}

func NewSeries(name, metadataPath string) (*Series, error) {
	s := &Series{Name: name}
	if metadataPath != "" {
		sm := &SeriesMetadata{}
		b, err := ioutil.ReadFile(metadataPath)
		if err != nil {
			return nil, fmt.Errorf("NewSeries error: %w", err)
		}
		err = json.Unmarshal(b, sm)
		if err != nil {
			return nil, fmt.Errorf("NewSeries error: %w", err)
		}
		s.Title = sm.Title
	}
	return s, nil
}

// Sorts a slice of Series by ascending name with "" first.
func sortSeries(s []*Series) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Name < s[j].Name
	})
}

//
// PostIndex
//

type PostIndex struct {
	postMap   map[string]map[string]*Post
	seriesMap map[string]*Series
	Posts     []*Post
	Series    []*Series
}

func (p *PostIndex) Get(series, name string) *Post {
	// log.Printf("PostIndex.Get: series=%q name=%q", series, name)
	s := p.postMap[series]
	if s != nil {
		return s[name]
	}
	return nil
}

// Returns the most recent post of all series.
func (p *PostIndex) GetLatest() *Post {
	if len(p.Posts) == 0 {
		return nil
	}
	return p.Posts[0]
}

func (p *PostIndex) GetPriorPosts(seriesName string) []*Post {
	latest := p.GetLatest()
	series, ok := p.seriesMap[seriesName]
	switch {
	case !ok:
		return nil

	case len(series.Posts) == 0:
		return nil

	case series.Posts[0] == latest:
		return series.Posts[1:]

	default:
		return series.Posts
	}
}

// Pattern for a post directory: ${ISO_8601}-${NAME}
var postRegexp = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})-(.*)`)

//
// Reads posts within a posts/ directory including down one layer for
// named series.
//
func readPosts(postsDir string, series *Series) ([]*Post, []*Series, error) {
	d, err := os.Open(postsDir)
	if err != nil {
		return nil, nil, err
	}
	defer d.Close()

	infos, err := d.Readdir(-1)
	if err != nil {
		// NB: IO errors are ignored by Glob!
		return nil, nil, err
	}

	posts := make([]*Post, 0)
	allSeries := make([]*Series, 0)

	if err != nil {
		return nil, nil, err
	}

	for _, info := range infos {
		// Skip files.
		if !info.IsDir() {
			continue
		}

		// Both posts and series can have metadata. Look for taht first.
		metadataPath := filepath.Join(d.Name(), info.Name(), "metadata.json")
		if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
			metadataPath = ""
		} else if err != nil {
			log.Printf("Failed to stat %s: %v", metadataPath, err)
			return nil, nil, err
		}

		contentPath := filepath.Join(d.Name(), info.Name(), "index.md")
		switch _, err := os.Stat(contentPath); {
		case err == nil:
			// index.md exists, so this is a post directory.
			basename := filepath.Base(info.Name())
			if m := postRegexp.FindStringSubmatch(basename); m != nil {
				post, err := NewPost(m[1], m[2], contentPath, metadataPath, series)
				if err != nil {
					return nil, nil, err
				}
				posts = append(posts, post)
				series.Posts = append(series.Posts, post)
			}

		case os.IsNotExist(err) && series.Name == "":
			// No index, but we're top-level. Treat this directory
			// as a series of posts.
			seriesName := info.Name()
			newSeries, err := NewSeries(info.Name(), metadataPath)
			if err != nil {
				return nil, nil, err
			}

			seriesPosts, _, err := readPosts(
				filepath.Join(d.Name(), seriesName), newSeries)
			if err != nil {
				return nil, nil, err
			}
			if len(posts) == 0 {
				return nil, nil, fmt.Errorf(
					"Series %s contained no posts", seriesName)
			}
			sortPosts(newSeries.Posts)
			posts = append(posts, seriesPosts...)
			allSeries = append(allSeries, newSeries)

		case os.IsNotExist(err) && series.Name != "":
			log.Printf("Expected index.md at %s, but not found", contentPath)
			return nil, nil, err

		default:
			log.Printf("Failed to stat %s: %v", contentPath, err)
			return nil, nil, err
		}
	}

	if len(series.Posts) > 0 {
		sortPosts(series.Posts)
		allSeries = append(allSeries, series)
	}

	return posts, allSeries, nil
}

// Loads posts from a directory into a PostIndex.
func NewPostIndex(contentDir string) (*PostIndex, error) {
	posts, series, err := readPosts(filepath.Join(contentDir, "posts"), &Series{})
	if err != nil {
		return nil, err
	}

	// Sort all posts.
	sortPosts(posts)

	// Sort all series
	sortSeries(series)

	pi := &PostIndex{
		postMap:   make(map[string]map[string]*Post),
		seriesMap: make(map[string]*Series),
		Posts:     posts,
		Series:    series,
	}

	// Index series by name
	for _, s := range series {
		pi.seriesMap[s.Name] = s
	}

	// Index posts by series
	for _, post := range posts {
		if _, ok := pi.postMap[post.Series.Name]; !ok {
			pi.postMap[post.Series.Name] = make(map[string]*Post)
		}
		pi.postMap[post.Series.Name][post.Name] = post
	}

	return pi, nil
}
