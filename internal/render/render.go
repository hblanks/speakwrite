package render

//
// Methods for rendering all content to disk.
//

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hblanks/speakwrite/internal/web"
)

type responseWriter struct {
	*httptest.ResponseRecorder
	f *os.File
}

func (r *responseWriter) Write(buf []byte) (int, error) {
	return r.f.Write(buf)
}

func writeURL(s *web.Server, u, outputRoot string) error {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	relpath := parsedURL.Path
	if relpath == "" {
		relpath = "/index.html"
	}
	if strings.HasSuffix(relpath, "/") {
		relpath += "index.html"
	}

	p := filepath.Join(outputRoot, relpath)
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := ioutil.TempFile(dir, "")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := os.Chmod(f.Name(), 0644); err != nil {
		return err
	}

	r := httptest.NewRequest("GET", u, nil)
	w := &responseWriter{httptest.NewRecorder(), f}
	s.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 {
		os.Remove(f.Name())
		return fmt.Errorf("URL %s returned %d, not 200!",
			u, w.Result().StatusCode)
	}
	os.Rename(f.Name(), p)
	return nil
}

func WriteURLs(s *web.Server, outputRoot string) error {
	urls, err := s.GetURLs()
	if err != nil {
		return err
	}
	for _, u := range urls {
		writeURL(s, u, outputRoot)
	}
	return nil
}
