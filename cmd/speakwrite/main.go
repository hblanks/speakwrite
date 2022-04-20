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
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage: %s [serve|render]\n", os.Args[0])
		fmt.Fprintf(out,
			`Required environment variables:
	CONTENT_DIR		= Path to site content/ dir
	THEME_DIR		= Path to theme/ dir
	PUBLIC_URL		= Public URL for the site

Optional environment variables:
	LISTEN_ADDR		= For "serve": listen address (default: localhost:8080)
	OUTPUT_DIR      = For "render": where to write output (default: speakwrite-out)
`)
	}
	flag.Parse()

	contentDir := os.Getenv("CONTENT_DIR")
	themeDir := os.Getenv("THEME_DIR")
	publicURL := os.Getenv("PUBLIC_URL")
	if contentDir == "" || themeDir == "" || publicURL == "" {
		flag.Usage()
		os.Exit(1)
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
		outputDir := os.Getenv("OUTPUT_DIR")
		if outputDir == "" {
			outputDir = "speakwrite-out"
		}
		if err := render.WriteURLs(server, outputDir); err != nil {
			log.Fatalf("Write error: %v", err)
		}

	case "serve":
		listenAddr := os.Getenv("LISTEN_ADDR")
		if listenAddr == "" {
			listenAddr = "localhost:8080"
		}

		log.Fatalf("Server listen error: %v",
			http.ListenAndServe(listenAddr, server))

	default:
		flag.Usage()
		os.Exit(1)

	}
}
