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
content/
    root/                         Static pages, mapped to /
    posts/
        ${ISO_8601}-{PATH_NAME}/  Article. Mapped to /posts/
            index.md              Article text.

theme/
    static/                       Static assets. Mapped to /static
    template/                     HTML templates.
```

Server commandline:

```
CONTENT_DIR=content THEME=theme interval
```
