CGO_ENABLED = 1

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.PHONY: graph
graph: artifacts/graph.png
	open "$<"

artifacts/graph.png: $(shell find . -name '*.go')
	@mkdir -p "$(@D)"
	go run cmd/graph/main.go | dot -Tpng -o "$@"

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

artifacts/plugin/%: $(GO_SOURCE_FILES) $(GENERATED_FILES)
	$(eval PARTS := $(subst /, ,$*))
	$(eval BUILD := $(word 1,$(PARTS)))
	$(eval OS    := $(word 2,$(PARTS)))
	$(eval ARCH  := $(patsubst arm%,arm,$(word 3,$(PARTS))))
	$(eval GOARM := $(patsubst arm%,%,$(filter arm%,$(word 3,$(PARTS)))))
	$(eval BIN   := $(word 4,$(PARTS)))
	$(eval PKG   := $(basename $(BIN)))
	$(eval ARGS  := $(if $(findstring debug,$(BUILD)),$(GO_DEBUG_ARGS),$(GO_RELEASE_ARGS)))

	CGO_ENABLED=$(CGO_ENABLED) GOOS="$(OS)" GOARCH="$(ARCH)" GOARM="$(GOARM)" go build $(ARGS) -o "$@" "./cmd/$(PKG)"
