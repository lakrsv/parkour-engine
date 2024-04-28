BINARY_NAME=game
BUILD_PATH=build

define package
	mkdir -p "./${BUILD_PATH}/$(1)/$(2)"
    mv ${BINARY_NAME}-$(1)* ./${BUILD_PATH}/$(1)/$(2)
    cp -r ./binaries/$(1)/$(2)/. ./${BUILD_PATH}/$(1)/$(2)
endef

.PHONY: clean prepare run build generate lint

lint:
	golangci-lint run ./engine/...
	golangci-lint run ./game/...

prepare:
	mkdir -p binaries/linux/amd64
	mkdir -p binaries/windows/amd64
	mkdir -p binaries/darwin/amd64

clean: prepare
	rm -rf ${BUILD_PATH}

run:
	cd ./game && go run .


build: generate build-windows build-linux build-darwin

build-windows: ${BUILD_PATH}/windows/amd64
${BUILD_PATH}/windows/amd64:
	# Build Windows Binary
	CGO_ENABLED="1" \
    CC="x86_64-w64-mingw32-gcc" \
	GOOS="windows" \
	GOARCH="amd64" \
	go build -o ${BINARY_NAME}-windows.exe ./game/
	$(call package,windows,amd64)

build-linux: ${BUILD_PATH}/linux/amd64

${BUILD_PATH}/linux/amd64:
	GOOS="linux" \
	GOARCH="amd64" \
	go build -o ${BINARY_NAME}-linux ./game/
	$(call package,linux,amd64)

build-darwin: ${BUILD_PATH}/darwin/amd64
${BUILD_PATH}/darwin/amd64:
	GOOS="darwin" \
	GOARCH="amd64" \
	go build -o ${BINARY_NAME}-darwin -trimpath ./game/
	$(call package,darwin,amd64)

generate:
	go generate ./engine/...