TARGET := physical-based-animations
GOPATH := $$(go env GOPATH)
BINDATA := bindata.go
.DEFAULT_GOAL: $(TARGET)

VERSION := 0.1
BUILD := `git rev-parse HEAD`

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: all install build clean goinstall uninstall check run

all: check install goinstall

$(BINDATA):
	@go get github.com/jteeuwen/go-bindata
	@$(GOPATH)/bin/go-bindata -debug assets/...

$(TARGET): *.go
	@go build $(LDFLAGS) -o $(TARGET)

install: $(BINDATA)
	@go get ./...

build: install $(TARGET)
	@true

clean:
	@rm -f $(BINDATA)
	@rm -f $(TARGET)

goinstall: install
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which $(GOPATH)/bin/${TARGET})

check:
	@test -z $(gofmt -l physical-based-animations.go | tee /dev/stderr) || \
	echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do $(GOPATH)/bin/golint $${d}; done
	@go tool vet *.go

run: goinstall
	@$(GOPATH)/bin/$(TARGET) &
