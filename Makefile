REGISTRY?=csemcr.azurecr.io
REGISTRY_PATH?=test/k8s/metrics
IMAGE?=adapter
TEMP_DIR:=$(shell mktemp -d)
ARCH?=amd64
OUT_DIR?=./_output
SEMVER=minor
PUSH_LATEST=true

VERSION?=latest
GOIMAGE=golang:1.10

ifeq ("$(REGISTRY_PATH)", "")
	FULL_IMAGE=$(REGISTRY)/$(IMAGE)
else
	FULL_IMAGE=$(REGISTRY)/$(REGISTRY_PATH)/$(IMAGE)
endif	

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

.PHONY: all build test verify build-container version save

all: build
build-local: test
	CGO_ENABLED=0 GOARCH=$(ARCH) go build -a -tags netgo -o $(OUT_DIR)/$(ARCH)/adapter github.com/Azure/azure-k8s-metrics-adapter

build: vendor verify-deploy
	docker build -t $(FULL_IMAGE):$(VERSION) .

save:
	docker save -o app.tar $(FULL_IMAGE):$(VERSION)

version: build	
ifeq ("$(BRANCH)", "master")
	@echo "versioning on master"
	go get -u github.com/jsturtevant/gitsem
	gitsem $(SEMVER)
else
	@echo "must be on clean master branch"
endif	

tag-ci:
	docker tag $(FULL_IMAGE):$(CIRCLE_WORKFLOW_ID) $(FULL_IMAGE):$(VERSION)
	
push:
	@echo $(DOCKER_PASS) | docker login -u $(DOCKER_USER) --password-stdin csemcr.azurecr.io 
	docker push $(FULL_IMAGE):$(VERSION)
ifeq ("$(PUSH_LATEST)", "true")
	@echo "pushing to latest"
	docker tag $(FULL_IMAGE):$(VERSION) $(FULL_IMAGE):latest
	docker push $(FULL_IMAGE):latest
endif		

vendor: 
	dep ensure

test: vendor
	CGO_ENABLED=0 go test ./pkg/...

verify-deploy:
	hack/verify-deploy.sh

gen-deploy:
	hack/gen-deploy.sh

dev:
	skaffold dev




