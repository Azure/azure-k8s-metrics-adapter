REGISTRY?=jsturtevant
IMAGE?=azure-k8-metrics-adapter
TEMP_DIR:=$(shell mktemp -d)
ARCH?=amd64
OUT_DIR?=./_output

VERSION?=latest
GOIMAGE=golang:1.10

.PHONY: all build test verify build-container

all: build
build: vendor
	CGO_ENABLED=0 GOARCH=$(ARCH) go build -a -tags netgo -o $(OUT_DIR)/$(ARCH)/adapter github.com/jsturtevant/azure-k8-metrics-adapter

vendor: 
	dep ensure

test: vendor
	CGO_ENABLED=0 go test ./pkg/...

build-container: build
	cp deploy/Dockerfile $(TEMP_DIR)
	cp $(OUT_DIR)/$(ARCH)/adapter $(TEMP_DIR)/adapter
	cd $(TEMP_DIR) && sed -i "s|BASEIMAGE|scratch|g" Dockerfile
	sed -i 's|REGISTRY|'${REGISTRY}'|g' deploy/manifests/custom-metrics-apiserver-deployment.yaml
	docker build -t $(REGISTRY)/$(IMAGE)-$(ARCH):$(VERSION) $(TEMP_DIR)
	rm -rf $(TEMP_DIR)

container-push: 
	docker push $(REGISTRY)/$(IMAGE)-$(ARCH):$(VERSION)
