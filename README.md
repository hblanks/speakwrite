# speakwrite

A minimal static blog renderer and dev server.

## Setup

To build the server, run:

```
make
```

## How it works: overview

speakwrite expects two directories: a `CONTENT_DIR` and a `THEME_DIR`.

`CONTENT_DIR` contains pages and posts:

```
content/
  <!-- TODO pages/
    {PAGE_NAME}.md              Pages. Mapped to /{PAGE_NAME}/ -->
  posts/
    {ISO_8601}-{POST_NAME}/     Articles. Mapped to /posts/{SERIES_NAME}/{POST_NAME}/
      index.md                  Article text.
    {SERIES_NAME}/
      {ISO_8601}-{POST_NAME}/   Articles in a series. Mapped to /posts/{SERIES_NAME}/{POST_NAME}/
        index.md                Article text.
```

`THEME_DIR` contains static assets and templates:

```
theme/
  static/                       Static assets. Mapped to /static
    v/                          Versioned static assets. Cached forever.
  template/                         
    base.html                   Base HTML template.
    root.html                   Template for /
    <!-- TODO page.html                   Template for pages not under posts/ -->
    post.html                   Template for articles under posts/
```

## What to put in a post

## What to put in a template

Template variables

Base (all pages):

Root:

Post:





## Development

To run the server locally with auto-reloading for content
and code changes:

```
CONTENT_DIR=/path/to/content THEME_DIR=/path/to/theme ./dev/watch.sh
```
