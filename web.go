package speakwrite

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hblanks/speakwrite/internal/web"
)

type ResponseWriter struct {
	*httptest.ResponseRecorder
	f *os.File
}

func (r *ResponseWriter) Write(buf []byte) (int, error) {
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
	w := &ResponseWriter{httptest.NewRecorder(), f}
	s.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 {
		os.Remove(f.Name())
		return fmt.Errorf("URL %s returned %d, not 200!",
			u, w.Result().StatusCode)
	}
	os.Rename(f.Name(), p)
	return nil
}

func writeURLs(s *web.Server, outputRoot string) error {
	urls, err := s.GetURLs()
	if err != nil {
		return err
	}
	for _, u := range urls {
		writeURL(s, u, outputRoot)
	}
	return nil
}

func Main() {
	outputPath := flag.String("-output", "build/html", "Output path when rendering")

	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage: %s\n", os.Args[0])
		fmt.Fprintf(out,
			`Environment variables:
	CONTENT_DIR		= Path to site content/ dir
	LISTEN_ADDR		= Listen address (default: localhost:8080)
	PUBLIC_URL		= Public URL for the site (default: http://localhost:8080)
	THEME_DIR		= Path to theme/ dir.
`)
	}
	flag.Parse()

	contentDir := os.Getenv("CONTENT_DIR")
	themeDir := os.Getenv("THEME_DIR")
	if contentDir == "" || themeDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		publicURL = "http://localhost:8080"
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = "localhost:8080"
	}

	server, err := web.NewServer(publicURL, contentDir, themeDir)
	if err != nil {
		log.Fatalf("Server init error: %v", err)
	}

	if args := flag.Args(); len(args) > 0 && args[0] == "serve" {
		log.Fatalf("Server listen error: %v",
			http.ListenAndServe(listenAddr, server))
	} else {
		if err := writeURLs(server, *outputPath); err != nil {
			log.Fatalf("Write error: %v", err)
		}
	}
}
