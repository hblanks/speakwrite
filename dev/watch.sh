#!/bin/sh -e

die() {
	echo $* > /dev/stderr
	exit 1
}

DIR=$(cd $(dirname $0)/..; pwd)

test -d "$CONTENT_DIR" || die "CONTENT_DIR must point to a directory"
test -d "$THEME_DIR" || die "THEME_DIR must point to a directory"

WATCH_DIRS=$(find $DIR -maxdepth 1 -mindepth 1 -type d \
    -not -name .git \
    -not -name build \
    -not -name docs)

# Make gomon run after startup since it lacks an option
# to do that.
(sleep 0.5 ; touch $DIR/cmd/speakwrite/main.go) &

make -C $DIR build/gomon
$DIR/build/gomon -d -R -m='\.(html|json|go|md)$$' \
    $WATCH_DIRS $CONTENT_DIR $THEME_DIR -- \
    sh -c "make && exec $DIR/build/speakwrite serve"
