VERSION := 0.1.0
DIST_DIR := dist
WASM_SRC := $(shell find . -name '*.go')
BUILD_FLAGS := -s -w -X main.version=$(VERSION)

.PHONY: run
run:
	@go run -tags=ebitenginedebug -ldflags "$(BUILD_FLAGS)" .

.PHONY: start
start: build
	@deno run -NR=$(DIST_DIR) jsr:@std/http/file-server $(DIST_DIR)

.PHONY: fmt
fmt:
	@go fmt ./...
	@deno fmt .

.PHONY: build
build: $(DIST_DIR)/index.html $(DIST_DIR)/west.wasm $(DIST_DIR)/wasm_exec.js

.PHONY: clean
clean:
	rm -rf $(DIST_DIR)

$(DIST_DIR)/wasm_exec.js:
	cp $(shell go env GOROOT)/lib/wasm/wasm_exec.js $@

$(DIST_DIR)/index.html: main.html
	@mkdir -p $(DIST_DIR)
	cp $< $@

%.wasm: $(WASM_SRC)
	GOOS=js GOARCH=wasm go build -ldflags "$(BUILD_FLAGS)" -o $@ .
