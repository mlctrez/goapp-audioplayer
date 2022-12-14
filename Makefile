
APP_NAME=music

VERSION=$(shell git describe --abbrev=0 --tags 2>/dev/null || echo "v0.0.0")
COMMIT=$(shell git rev-parse --short HEAD || echo "HEAD")
MOD=$(shell grep ^module go.mod | awk '{print $$2;}')
DIR=$(shell pwd)
LD_FLAGS="-w -X $(MOD)/goapp.Version=$(VERSION) -X $(MOD)/goapp.Commit=$(COMMIT) -X $(MOD)/goapp.BuildDir=$(DIR) -X $(MOD)/goapp.Module=$(MOD)"
MAIN="goapp/service/main/main.go"

.PHONY: model

run: binary
	@DEV=1 ./temp/$(APP_NAME)

binary: wasm
	@mkdir -p temp
	@#echo "ldflags=$(LD_FLAGS)"
	@go build -o temp/$(APP_NAME) -ldflags $(LD_FLAGS) $(MAIN)

wasm: model
	@rm -f goapp/web/app.wasm
	@GOARCH=wasm GOOS=js go build -o goapp/web/app.wasm -ldflags $(LD_FLAGS) $(MAIN)

model:
	@go run model/generate/generate.go

clean:
	@rm -rf temp
	@rm -f goapp/web/app.wasm

# used only to build the create binary in github.com/mlctrez/goappcreate
create:
	@rm -f goapp/web/app.wasm
	@go build -o temp/create .

deploy: binary
	sudo temp/music -action deploy