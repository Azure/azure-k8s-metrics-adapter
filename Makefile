# Following gets combined into: REGISTRY/REGISTRY_PATH/IMAGE:VERSION
REGISTRY?=csemcr.azurecr.io
IMAGE?=test/k8s/metrics/adapter
VERSION?=latest

ifneq ("$(REGISTRY)", "")
	FULL_IMAGE=$(REGISTRY)/$(IMAGE)
else
	FULL_IMAGE=$(IMAGE)
endif

OUT_DIR?=./_output
SEMVER=""
PUSH_LATEST=true

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

.PHONY: all build-local build vendor test version push \
		verify-deploy gen-deploy dev save tag-ci

all: build
build-local: test
	CGO_ENABLED=0 go build -a -tags netgo -o $(OUT_DIR)/adapter github.com/Azure/azure-k8s-metrics-adapter

build: vendor verify-deploy verify-apis
	docker build -t $(FULL_IMAGE):$(VERSION) .

vendor: 
	dep ensure -v

test: vendor
	hack/run-tests.sh

version: build	
ifeq ("$(SEMVER)", "")
	@echo "Please set sem version bump: can be 'major', 'minor', or 'patch'"
	exit
endif
ifeq ("$(BRANCH)", "master")
	@echo "versioning on master"
	go get -u github.com/jsturtevant/gitsem
	gitsem $(SEMVER)
else
	@echo "must be on clean master branch"
endif	

push:
ifdef DOCKER_PASS
	# non interactive login
	@echo $(DOCKER_PASS) | docker login -u $(DOCKER_USER) --password-stdin $(REGISTRY) 
else
	# interactive login (needed for WSL)
	docker login -u $(DOCKER_USER) $(REGISTRY)
endif
	docker push $(FULL_IMAGE):$(VERSION)
ifeq ("$(PUSH_LATEST)", "true")
	@echo "pushing to latest"
	docker tag $(FULL_IMAGE):$(VERSION) $(FULL_IMAGE):latest
	docker push $(FULL_IMAGE):latest
endif		

# dev setup
dev:
	skaffold dev

teste2e:
	hack/run-e2e.sh

# CI specific commands used during CI build
save:
	docker save -o app.tar $(FULL_IMAGE):$(VERSION)

tag-ci:
	docker tag $(FULL_IMAGE):$(CIRCLE_WORKFLOW_ID) $(FULL_IMAGE):$(VERSION)

# Code gen helpers
gen-apis: codegen-get
	hack/update-codegen.sh

verify-apis: codegen-get
	hack/verify-codegen.sh

codegen-get:
	go get -u k8s.io/code-generator/...
	hack/codegen-repo-fix.sh

# Helm deploy generator helpers
verify-deploy:
	hack/verify-deploy.sh

gen-deploy:
	hack/gen-deploy.sh

gen-all: gen-apis gen-deploy
