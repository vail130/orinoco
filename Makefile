default: build

clean:
	rm -rf bin

build:
	go build -o bin/orinoco orinoco.go

deps:
	go get

test:
	make build
	./scripts/run-tests.sh