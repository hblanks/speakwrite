.DEFAULT_GOAL := all

BUILD := build

#
# gomon
#

GOMON := $(BUILD)/gomon

$(GOMON): # go.sum
	go mod download
	@mkdir -p $(BUILD)
	GOBIN=$$(cd $(BUILD); pwd) go install github.com/c9s/gomon

#
# go
#

GO_SOURCES := \
	go.sum \
	$(shell find . -path ./$(BUILD) -prune -o -name \*.go -print)

GO_OUTPUTS := \
	$(BUILD)/speakwrite

$(GO_OUTPUTS): $(GO_SOURCES)
	@mkdir -p $(BUILD)
	go build -o $(BUILD)/ ./cmd/...

.PHONY: test
test:
	go test ./...

#
# all
#

.PHONY: all
all: $(GO_OUTPUTS) $(GOMON) test
