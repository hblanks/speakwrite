# speakwrite

A minimal static blog renderer and dev server.

You won't want this. But, it's what I use for
[confidentialinterval.com](https://confidentialinterval.com/).

## To run the server locally

To run the renderer as a dev server with auto-reload on code,
content, or template changes:

```
CONTENT_DIR=/path/to/content THEME_DIR=/path/to/theme ./dev/watch.sh
```

## To render into a directory

```
CONTENT_DIR=/path/to/content \
  THEME_DIR=/path/to/theme \
  PUBLIC_URL=http://website.url \
  OUTPUT_DIR=/path/to/html/output \
  build/speakwrite render
```

## Development

Build and run tests with:

```
make
```

## How it works: overview

speakwrite expects two directories: a `CONTENT_DIR` and a `THEME_DIR`.

`CONTENT_DIR` contains pages and posts:

```
content/
  posts/
    metadata.json               (Optional) metadata about the entire post series.
    {ISO_8601}-{POST_NAME}/     Articles. Mapped to /posts/{SERIES_NAME}/{POST_NAME}/
      index.md                  Article text.
      metadata.json             (Optional) metadata about this post.
    {SERIES_NAME}/
      metadata.json             (Optional) metadata about this named series.
      {ISO_8601}-{POST_NAME}/   Articles in a series. Mapped to /posts/{SERIES_NAME}/{POST_NAME}/
        index.md                Article text.
        metadata.json           (Optional) metadata about this post.
```

`THEME_DIR` contains static assets and templates:

```
theme/
  static/                       Static assets. Mapped to /static
    v/                          Versioned static assets. Cached forever.
  template/                         
    base.html                   Base HTML template.
    root.html                   Template for /
    post.html                   Template for articles under posts/
```

## What to put in a post

Well. It's markdown, with support for footnotes and some other stuff
you'll find in the code. But if you put a pandoc title block at the top
that will get parsed and used. I.e.

```markdown
% This is the title of the post

This is the content.
```

## What to put in a template

To paraphrase the resource fork of an old Marathon binary:

> Wouldn't you know it's all in the source code.
