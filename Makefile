PB_FILES = $(shell find . -name "*.pb.go")
GENERATED_FILES += $(PB_FILES)
GENERATED_FILES += www/dist/main.js

-include artifacts/make/go.mk

artifacts/make/%.mk:
	curl -sf https://dogmatiq.io/makefiles/fetch | bash /dev/stdin $*

JS_PB_DIR=www/src/pb
%.pb.go: %.proto
	@mkdir -p $(@D)
	protoc --go_out=paths=source_relative,plugins=grpc:. $(@D)/*.proto
	@mkdir -p $(JS_PB_DIR)
	protoc 	--proto_path=$(@D) \
		--js_out=import_style=commonjs:$(JS_PB_DIR)\
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:$(JS_PB_DIR)\
		$(@D)/*.proto

www/node_modules:
	cd www; npm install

www/dist/main.js: www/src/index.js www/node_modules $(PB_FILES)
	cd www; npx webpack --mode=development $(notdir $(<D))/$(<F)
