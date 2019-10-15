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
