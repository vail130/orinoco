default: build

clean:
	rm -rf bin

build:
	go build -o bin/orinoco main.go

deps:
	go get

test:
	docker run --rm -e "TEST_PKG=$${pkg}" -v `pwd`:/go/src/github.com/vail130/orinoco vail130/orinoco-test

build-docker-images:
	docker build -t vail130/orinoco-base -f docker/base.Dockerfile docker
	docker build -t vail130/orinoco-test -f docker/test.Dockerfile docker
	docker build -t vail130/orinoco-sieve -f docker/sieve.Dockerfile docker
	docker build -t vail130/orinoco-tap -f docker/tap.Dockerfile docker
	docker build -t vail130/orinoco-pump -f docker/pump.Dockerfile docker
	docker build -t vail130/orinoco-litmus -f docker/litmus.Dockerfile docker
	docker build -t vail130/orinoco-orinoco -f docker/orinoco.Dockerfile docker
	
push-docker-images:
	docker push vail130/orinoco-base
	docker push vail130/orinoco-test
	docker push vail130/orinoco-sieve
	docker push vail130/orinoco-tap
	docker push vail130/orinoco-pump
	docker push vail130/orinoco-litmus
	docker push vail130/orinoco-orinoco