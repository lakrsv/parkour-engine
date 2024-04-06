BINARY_NAME=game
BUILD_PATH=build

define package
	mkdir -p "./${BUILD_PATH}/$(1)/$(2)"
    mv ./game/${BINARY_NAME}* ./${BUILD_PATH}/$(1)/$(2)
    cp -r ./game/assets ./${BUILD_PATH}/$(1)/$(2)/assets
    cp -r ./binaries/$(1)/$(2)/. ./${BUILD_PATH}/$(1)/$(2)
endef

clean:
	rm -rf ${BUILD_PATH}


build: ${BUILD_PATH}/windows/amd64

${BUILD_PATH}/windows/amd64:
	# Build Windows Binary
	CGO_ENABLED="1" \
    CC="x86_64-w64-mingw32-gcc" \
	GOOS="windows" \
	GOARCH="amd64" \
	go build -o ${BINARY_NAME} ./game/
	$(call package,windows,amd64)