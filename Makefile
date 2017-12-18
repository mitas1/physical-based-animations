COMMIT := $(shell git rev-parse HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

CC_WINDOWS := i686-w64-mingw32-gcc
CC_LINUX := x86_64-pc-linux-gcc
CC_DARWIN := o64-clang

ARCH_LINUX := amd64
ARCH_WINDOWS := 386
ARCH_DARWIN := amd64

TARGET := physical-based-animations
BINDATA := bindata.go
.DEFAULT_GOAL: $(TARGET)

LINUX_TARGET := $(TARGET)-linux-$(ARCH_LINUX)
WINDOWS_TARGET := $(TARGET)-windows-$(ARCH_WINDOWS).exe
DARWIN_TARGET := $(TARGET)-darwin-$(ARCH_DARWIN)

VERSION := 0.1

GITHUB_USERNAME := mitas1
BUILD_DIR := ${GOPATH}/src/github.com/${GITHUB_USERNAME}/${TARGET}/bin
CURRENT_DIR := $(shell pwd)

LDFLAGS=-ldflags '-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}'

all: install goinstall

$(BINDATA):
	@go get -v github.com/jteeuwen/go-bindata/...
	@go-bindata assets/...

install: $(BINDATA)
	@go get -v ./...

$(BUILD_DIR):
	@mkdir -p ${BUILD_DIR}

$(LINUX_TARGET): install $(BUILD_DIR)
	@./build.sh "linux" "${VERSION}" "${BUILD_DIR}/${LINUX_TARGET}" "${CC_LINUX}" "${ARCH_LINUX}"

$(WINDOWS_TARGET): install $(BUILD_DIR)
	@./build.sh "windows" "${VERSION}" "${BUILD_DIR}/${WINDOWS_TARGET}" "${CC_WINDOWS}" "${ARCH_WINDOWS}"

$(DARWIN_TARGET): install $(BUILD_DIR)
	@./build.sh "darwin" "${VERSION}" "${BUILD_DIR}/${DARWIN_TARGET}" "${CC_DARWIN}" "${ARCH_DARWIN}"

deps:
	@go get -v github.com/tools/godep
	@godep save
	@-rm -rf vendor/github.com/go-gl/glfw
	@git clone https://github.com/go-gl/glfw.git vendor/github.com/go-gl/glfw
	@ls -d ${CURRENT_DIR}/vendor/github.com/go-gl/glfw/** | grep -P "^.+[^(v3.2)]$$" | xargs -d"\n" rm -rf
	@rm -f vendor/github.com/go-gl/glfw/.travis.yml

linux: $(LINUX_TARGET)

darwin: $(DARWIN_TARGET)

windows: $(WINDOWS_TARGET)

build: $(LINUX_TARGET) $(DARWIN_TARGET) $(WINDOWS_TARGET)

clean:
	@rm -fr ${BUILD_DIR}
	@rm -fr vendor
	@rm -f $(BINDATA)

goinstall: install
	@go install $(LDFLAGS)

uninstall: clean
	@rm $$(which ${TARGET})

run: goinstall
	@${TARGET} &

test:
	go test -v ./...

.PHONY: all install build clean goinstall uninstall run
