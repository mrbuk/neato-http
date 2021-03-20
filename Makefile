name = neato-http
version = 0.5
user = mrbuk

all: build
.PHONY : all

GOPATH := $(shell go env GOPATH)

dep:
	go get -u github.com/onsi/ginkgo/ginkgo
	go mod tidy

run-integrationtest:
	#$(GOPATH)/bin/ginkgo -r

test:
	#$(GOPATH)/bin/ginkgo -r -skipPackage integrationtest
	
lint:
	golint .
	
build: test lint
	mkdir -p build
	go build -o ./build/${name}

docker-build:
	docker build . -t ${user}/${name}:$(version) -t ${user}/${name}:latest

docker-push: docker-build
	docker push ${user}/${name}:$(version)
	docker push ${user}/${name}:latest

docker-save:
	mkdir -p images
	docker save ${user}/${name}:latest -o images/${user}_${name}_latest.tgz

install-service:
	./scripts/install-service.sh ${name}.service

