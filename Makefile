-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

.PHONY: run
run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/bank
	$<
