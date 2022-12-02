
APP_NAME=music

VERSION=$(shell git describe --abbrev=0 --tags 2>/dev/null || echo "v0.0.0")
COMMIT=$(shell git rev-parse --short HEAD || echo "HEAD")
MODULE=$(shell grep ^module go.mod | awk '{print $$2;}')
LD_FLAGS="-w -X $(MODULE)/goapp.Version=$(VERSION) -X $(MODULE)/goapp.Commit=$(COMMIT)"
MAIN="goapp/service/main/main.go"

.PHONY: model

run: binary
	# ADDRESS=:9021
	@DEV=1 ADDRESS=:9021 ./temp/$(APP_NAME)

binary: wasm
	@mkdir -p temp
	@echo "ldflags=$(LD_FLAGS)"
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