default: build

clean:
	rm -rf bin

build:
	go build -o bin/orinoco main.go

deps:
	go get

test:
	docker run --rm \
		-e "TEST_PKG=$${pkg}" \
		-v `pwd`:/go/src/github.com/vail130/orinoco \
		vail130/orinoco-test

test-s3:
	docker run --rm \
		-e "TEST_PKG=$${pkg}" \
		-e "AWS_ACCESS_KEY_ID=$${AWS_ACCESS_KEY_ID}" \
		-e "AWS_SECRET_ACCESS_KEY=$${AWS_SECRET_ACCESS_KEY}" \
		-v `pwd`:/go/src/github.com/vail130/orinoco \
		vail130/orinoco-test-s3

build-docker-images:
	docker build -t vail130/orinoco-base -f docker/base.Dockerfile docker
	docker build -t vail130/orinoco-test -f docker/test.Dockerfile docker
	docker build -t vail130/orinoco-test-s3 -f docker/test-s3.Dockerfile docker
	docker build -t vail130/orinoco-orinoco -f docker/orinoco.Dockerfile docker
	
push-docker-images:
	docker push vail130/orinoco-base
	docker push vail130/orinoco-test
	docker push vail130/orinoco-test-s3
	docker push vail130/orinoco-orinoco