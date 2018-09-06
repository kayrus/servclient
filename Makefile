APP_NAME=servclient
DOCKER_REPO=kayrus
TAG=v1

build:
	CGO_ENABLED=0 GOPATH=$(PWD) go build -o siclient client
	CGO_ENABLED=0 GOPATH=$(PWD) go build -o siserver server

buildd:
	docker build -t $(DOCKER_REPO)/$(APP_NAME):$(TAG) -f Dockerfile .

buildd-nc:
	docker build --no-cache -t $(DOCKER_REPO)/$(APP_NAME):$(TAG) -f Dockerfile .

push:
	@echo 'push $(TAG) to $(DOCKER_REPO)/$(APP_NAME)'
	docker push $(DOCKER_REPO)/$(APP_NAME):$(TAG)

fmt:
	find . -name '*.go' -exec go fmt {} \;

all: build buildd
