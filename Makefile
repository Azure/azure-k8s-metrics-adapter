REGISTRY?=jsturtevant
IMAGE?=azure-k8-metrics-adapter
TEMP_DIR:=$(shell mktemp -d)
ARCH?=amd64
OUT_DIR?=./_output
SEMVER=minor

VERSION?=latest
GOIMAGE=golang:1.10

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

.PHONY: all build test verify build-container version save

all: build
build-local: vendor
	CGO_ENABLED=0 GOARCH=$(ARCH) go build -a -tags netgo -o $(OUT_DIR)/$(ARCH)/adapter github.com/Azure/azure-k8s-metrics-adapter

build:
	docker build -t $(REGISTRY)/$(IMAGE)-$(ARCH):$(VERSION) .

save:
	docker save -o app.tar $(REGISTRY)/$(IMAGE)-$(ARCH):$(VERSION)

version:	
ifeq ("$(BRANCH)", "master")
	@echo "versioning on master"
	go get -u github.com/Clever/gitsem
	gitsem $(SEMVER)
else
	@echo "must be on clean master branch"
endif	
	
push:
	@docker login -u $(DOCKER_USER) -p '$(DOCKER_PASS)'    
	docker push $(REGISTRY)/$(IMAGE)-$(ARCH):$(VERSION)

vendor: 
	dep ensure

test: vendor
	CGO_ENABLED=0 go test ./pkg/...

dev:
	skaffold dev

deploy:
	skaffold run




