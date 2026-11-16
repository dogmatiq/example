DOCKER_REPO = ghcr.io/dogmatiq/example

DOCKER_PLATFORMS += linux/amd64
DOCKER_PLATFORMS += linux/arm64

GO_EMBEDDED_FILES += $(wildcard ui/projections/*.sql)
GO_EMBEDDED_FILES += $(wildcard ui/templates/*.html ui/templates/*.css)

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile
-include .makefiles/pkg/docker/v1/Makefile
-include .makefiles/pkg/apm/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
