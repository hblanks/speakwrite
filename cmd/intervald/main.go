package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hblanks/confint/internal/web"
)

func main() {
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

	log.Fatalf("Server listen error: %v",
		http.ListenAndServe(listenAddr, server))
}
