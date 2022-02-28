.DEFAULT_GOAL := all

BUILD := build

#
# go
#

GO_SOURCES := \
	go.sum \
	$(shell find . -path ./$(BUILD) -prune -o -name \*.go -print)

GO_OUTPUTS := \
	$(BUILD)/blog-httpd

$(GO_OUTPUTS): $(GO_SOURCES)
	@mkdir -p $(BUILD)
	go build -o $(BUILD)/ ./cmd/...

.PHONY: all
all: $(GO_OUTPUTS)


#
# docker
#


.PHONY: docker
docker: all
	docker build -t speakwrite:latest .


#
# watch
#

# Because gomon doesn't like watching the whole repo..
WATCH_DIRS = $(shell \
	find . -maxdepth 1 -mindepth 1 -type d \
	-not -name .git \
	-not -name build \
	-not -name content \
	-not -name docs)


.PHONY: watch
watch:
ifeq ($(THEME_DIR),)
	$(error THEME_DIR must be set for watch target)
endif
ifeq ($(CONTAINER),)
	$(error CONTAINER must be set for watch target)
endif
	gomon -d -R -m='\.(go|html)$$' $(WATCH_DIRS) $(THEME_DIR) \
			-- sh -c "make && docker restart $(CONTAINER)"

.PHONY: html
html: all
	rm -rf $(BUILD)/html
	dev/blog-httpd.sh

.PHONY: deploy
deploy: html
	rsync -av --delete --checksum \
		--exclude="*.swp" --delete-excluded \
		build/html/ fex:/state/home/web/confidentialinterval.com/html/
