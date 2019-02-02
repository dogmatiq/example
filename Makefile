GENERATED_FILES += $(shell find . -name "*.pb.go")

-include artifacts/make/go.mk

artifacts/make/%.mk:
	curl -sf https://dogmatiq.io/makefiles/fetch | bash /dev/stdin $*


%.pb.go: %.proto
	protoc --go_out=paths=source_relative,plugins=grpc:. $(@D)/*.proto
