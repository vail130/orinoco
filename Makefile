PACKAGES ?= $$(go list ./... | grep -v '/vendor/')
FILES ?= $$(find . -type f -name "*.go" | grep -v '/vendor/')

default: clean configure fmt lint test build

all: clean configure build

clean:
	rm -rf bin

configure:
	go get github.com/tools/godep
	go get github.com/axw/gocov/...
	go get github.com/golang/lint/golint
	go get github.com/jteeuwen/go-bindata/...
	go get gopkg.in/check.v1
	go get
	go install

build:
	go build -o bin/orinoco main.go

fmt:
	gofmt -l -s -w $(FILES)

.PHONY: test
test:
	TEST=1 go test -v $(PACKAGES)

lint:
	go tool vet ./*/*.go \
		&& find . -name "*.go" -maxdepth 2 | grep -v "templates.go" | xargs -n1 golint

.PHONY: vendor
vendor:
	godep save
