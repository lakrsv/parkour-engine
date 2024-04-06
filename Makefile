BINARY_NAME=game
BUILD_PATH=build

define package
	mkdir -p "./${BUILD_PATH}/$(1)/$(2)"
    mv ${BINARY_NAME}-$(1)* ./${BUILD_PATH}/$(1)/$(2)
    cp -r ./game/assets ./${BUILD_PATH}/$(1)/$(2)/assets
    cp -r ./binaries/$(1)/$(2)/. ./${BUILD_PATH}/$(1)/$(2)
endef

.PHONY: clean run build

clean:
	rm -rf ${BUILD_PATH}

run:
	cd ./game && go run .


build: build-win64 build-linux

build-win64: ${BUILD_PATH}/windows/amd64
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
