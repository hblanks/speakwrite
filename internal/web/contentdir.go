package web

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// Type for serving static files within an individual post's directory.
type ContentDir string

// Returns files to serve over HTTP, but only if those files aren't the raw
// markdown and aren't in a directory called "exclude."
func (c ContentDir) Open(name string) (http.File, error) {
	dir := string(c)

	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	if fullName == filepath.Join(dir, "index.md") {
		return nil, os.ErrNotExist
	}
	relpath, _ := filepath.Rel(dir, fullName)
	for _, p := range filepath.SplitList(relpath) {
		if p == "exclude" {
			return nil, os.ErrNotExist
		}
	}

	f, err := http.Dir(dir).Open(name)
	if err != nil {
		return f, err
	}
	if st, err := f.Stat(); err != nil {
		return nil, err
	} else if st.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}
