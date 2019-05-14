-include artifacts/make/go.mk

.PHONY: graph
graph: artifacts/graph.png
	open "$<"

artifacts/graph.png: $(shell find . -name '*.go')
	go run cmd/graph/main.go | dot -Tpng -o "$@"

artifacts/make/%.mk:
	curl -sf https://dogmatiq.io/makefiles/fetch | bash /dev/stdin $*
