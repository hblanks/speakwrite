#!/bin/sh -e

DIR=$(cd $(dirname $0)/..; pwd)

THEME_DIR=$DIR/theme \
    CONTENT_DIR=$DIR/content \
    LISTEN_ADDR=localhost:8080 \
    exec $DIR/build/intervald "$@"
