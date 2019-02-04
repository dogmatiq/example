GENERATED_FILES += $(shell find . -name "*.pb.go")
GENERATED_FILES += web/assets/js/dist/main.js

-include artifacts/make/go.mk

artifacts/make/%.mk:
	curl -sf https://dogmatiq.io/makefiles/fetch | bash /dev/stdin $*

JS_PB_OUTPUT_DIR=web/assets/js
%.pb.go: %.proto
	protoc --go_out=paths=source_relative,plugins=grpc:. $(@D)/*.proto
	@mkdir -p $(JS_PB_OUTPUT_DIR)
	protoc 	--proto_path=$(@D) \
		--js_out=import_style=commonjs:$(JS_PB_OUTPUT_DIR)\
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:$(JS_PB_OUTPUT_DIR)\
		$(@D)/*.proto


web/assets/js/dist/main.js: web/assets/js/client.js
	cd $(<D); npx webpack --mode=development $(<F)
