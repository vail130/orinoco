default: build

clean:
	rm -rf bin

build:
	go build -o bin/orinoco orinoco.go

deps:
	go get

test:
	docker run --rm -e "TEST_PKG=$${pkg}" -v `pwd`:/go/src/github.com/vail130/orinoco vail130/orinoco-test
