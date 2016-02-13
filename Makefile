default: build

clean:
	rm -rf bin

build:
	go build -o bin/orinoco main.go

deps:
	go get

# make pkg=PKG svc=SVC
test:
	docker run --rm \
		-e "TEST_PKG=$${pkg}" \
		-e "TEST_SVC=$${svc:-log}" \
		-v `pwd`:/go/src/github.com/vail130/orinoco \
		vail130/orinoco-test
		
test-stdout:
	make test svc=stdout

test-log:
	make test svc=log

test-http:
	make test svc=http

test-s3:
	docker run --rm \
		-e "TEST_PKG=$${pkg}" \
		-e "TEST_SVC=s3" \
		-e "AWS_ACCESS_KEY_ID=$${AWS_ACCESS_KEY_ID}" \
		-e "AWS_SECRET_ACCESS_KEY=$${AWS_SECRET_ACCESS_KEY}" \
		-v `pwd`:/go/src/github.com/vail130/orinoco \
		vail130/orinoco-test

build-docker-images:
	docker build -t vail130/orinoco-base -f docker/base.Dockerfile docker
	docker build -t vail130/orinoco-test -f docker/test.Dockerfile docker
	docker build -t vail130/orinoco-orinoco -f docker/orinoco.Dockerfile docker
	
push-docker-images:
	docker push vail130/orinoco-base
	docker push vail130/orinoco-test
	docker push vail130/orinoco-orinoco
