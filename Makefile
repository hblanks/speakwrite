.DEFAULT_GOAL := all

BUILD := build

#
# go
#

GO_SOURCES := \
	go.sum \
	$(shell find . -path ./$(BUILD) -prune -o -name \*.go -print)

GO_OUTPUTS := \
	$(BUILD)/intervald

$(GO_OUTPUTS): $(GO_SOURCES)
	@mkdir -p $(BUILD)
	go build -o $(BUILD)/ ./cmd/...

.PHONY: all
all: $(GO_OUTPUTS)


#
# docker-compose
#

.PHONY: up
up: all
	@docker-compose down
	@docker-compose up -d
	@sleep 0.5
	@docker-compose exec web true


# Because gomon doesn't like watching the whole repo..
WATCH_DIRS = $(shell \
	find . -maxdepth 1 -mindepth 1 -type d \
	-not -name .git \
	-not -name build \
	-not -name content \
	-not -name docs)


.PHONY: watch
watch:
	gomon -d -R -m='\.(go|html)$$' $(WATCH_DIRS) \
			-- sh -c "make && docker-compose restart web"

.PHONY: deploy
deploy: all
	rm -rf $(BUILD)/html
	dev/intervald.sh
	rsync -av --delete \
		build/html/ fex:/state/home/web/confidentialinterval.com/html/
