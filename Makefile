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
	mkdir -p binaries/macos/amd64

clean: prepare
	rm -rf ${BUILD_PATH}

run:
	cd ./game && go run .


BUILD_NATIVE = build-linux
ifeq ($(shell uname -s), Darwin)
	BUILD_NATIVE = build-macos
endif

build: generate build-windows ${BUILD_NATIVE}

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

build-macos: ${BUILD_PATH}/macos/amd64
${BUILD_PATH}/macos/amd64:
	GOOS="darwin" \
	GOARCH="amd64" \
	go build -o ${BINARY_NAME}-macos -trimpath ./game/
	$(call package,macos,amd64)

generate:
	go generate ./engine/...