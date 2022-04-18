package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hblanks/speakwrite/internal/render"
	"github.com/hblanks/speakwrite/internal/web"
)

func main() {
	outputPath := flag.String("output", "build/html", "Output path when rendering")

	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage: %s [serve|render]\n", os.Args[0])
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

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	server, err := web.NewServer(publicURL, contentDir, themeDir)
	if err != nil {
		log.Fatalf("Server init error: %v", err)
	}

	switch args[0] {
	case "render":
		if err := render.WriteURLs(server, *outputPath); err != nil {
			log.Fatalf("Write error: %v", err)
		}

	case "serve":
		log.Fatalf("Server listen error: %v",
			http.ListenAndServe(listenAddr, server))

	default:
		flag.Usage()
		os.Exit(1)

	}
}
