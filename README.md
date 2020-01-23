# confidential interval

## Setup

To build the server (TODO and validate markdown?), run:

```
make
```

## Development

The devserver runs on http://localhost:8080/ with:

```
make up watch
```

## Deployment

```
make deploy
```

## How it works

Directory layout:

```
cmd/                              Server & renderer entrypoint.
content/
    root/                         Static pages, mapped to /
    posts/
        ${ISO_8601}-{PATH_NAME}/  Article. Mapped to /posts/
            index.md              Article text.
dev/                              Development scripts
theme/
    static/                       Static assets. Mapped to /static
      v/                          Versioned static assets. Cached forever.
    template/                     HTML templates.
    vendor/                       Vendored dependencies (CSS, JS)
internal/                         All go code.
```

Server commandline:

```
CONTENT_DIR=content THEME=theme interval
```
