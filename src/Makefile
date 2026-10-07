.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

DIARKIS_VERSION = $(shell go mod edit -json | jq -r '.Require[] | select(.Path | contains("github.com/Diarkis/diarkis")).Version')
DIARKIS_CLI = ./diarkis-cli/os/linux/bin/diarkis-cli
BUILD_CONFIG = ./build/linux-build.yml
BUILD_CONFIG_AMD64 = ./build/linux-build-amd64.yml
BUILD_CONFIG_ARM64 = ./build/linux-build-arm64.yml
DIARKIS_CLI_HOST = v3.builder.diarkis.io
AWS_ACCOUNT_NUM = __AWS_ACCOUNT_NUM__
AWS_REGION = ap-northeast-1
AWS_CONTAINER_REGISTRY = $(AWS_ACCOUNT_NUM).dkr.ecr.$(AWS_REGION).amazonaws.com
GCP_PROJECT_ID = __GCP_PROJECT_ID__
GCP_REGION = asia-northeast1
GCP_CONTAINER_REGISTRY = $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis
ACR_DOMAIN = __ACR_DOMAIN__

ifeq ($(shell uname), Darwin)
	DIARKIS_CLI = ./diarkis-cli/os/mac/bin/diarkis-cli
	BUILD_CONFIG = ./build/mac-build.yml
endif

.PHONY: init
init: ## Initialize project
	./script/change_diarkis_version.sh $(DIARKIS_VERSION)

.PHONY: build-local
build-local: gen ## Build server binary for local use
	echo "Build server binaries for local use"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c $(BUILD_CONFIG) --host $(DIARKIS_CLI_HOST)
	chmod -R 700 remote_bin/

.PHONY: build-linux
build-linux: gen ## Build server binary for linux
	echo "Build server binaries"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c ./build/linux-build.yml --host $(DIARKIS_CLI_HOST)
	chmod -R 700 remote_bin/

.PHONY: build-linux-amd64
build-linux-amd64: gen ## Build server binary for linux/amd64 container images (output: remote_bin/linux/amd64/)
	echo "Build server binaries for linux/amd64"
	rm -rfv remote_bin/linux/amd64
	$(DIARKIS_CLI) build -c $(BUILD_CONFIG_AMD64) --host $(DIARKIS_CLI_HOST)
	chmod -R 700 remote_bin/linux/amd64/

.PHONY: build-linux-arm64
build-linux-arm64: gen ## Build server binary for linux/arm64 container images (output: remote_bin/linux/arm64/)
	echo "Build server binaries for linux/arm64"
	rm -rfv remote_bin/linux/arm64
	$(DIARKIS_CLI) build -c $(BUILD_CONFIG_ARM64) --host $(DIARKIS_CLI_HOST)
	chmod -R 700 remote_bin/linux/arm64/

.PHONY: build-mac
build-mac: gen ## Build server binary for mac use
	echo "Build server binaries"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c ./build/mac-build.yml --host $(DIARKIS_CLI_HOST)
	chmod -R 700 remote_bin/

.PHONY: server
server: ## Start a server locally: [ target=mars ] [ target=http ] [ target=udp ] [ target=tcp ]
	echo "Starting $(target) server..."

ifeq ($(target), mars)
	./remote_bin/$(target) ./configs/mars/main.json
else
	./remote_bin/$(target)
endif

.PHONY: change-diarkis-version
change-diarkis-version: ## Change diarkis version
	./script/change_diarkis_version.sh

.PHONY: gen
gen: ## Generate go, cpp, and cs code files using puffer (Diarkis packet gen module) from packet definition written in json
	make -C puffer gen

.PHONY: clean
clean: ## Deletes all puffer generated code files
	make -C puffer clean

.PHONY: go-cli
go-cli: ## Starts Go test client: host=<HTTP address> uid=<client user ID> clientKey=<client key> puffer=<true/false>
	./remote_bin/testcli --host=$(host) --uid=$(uid) --clientKey=$(clientKey) --puffer=$(puffer)

.PHONY: run-docker
run-docker: build-linux-amd64 ## Run Diarkis locally with docker compose
	docker compose up -d --build

.PHONY: stop-docker
stop-docker: ## Stop Diarkis
	docker compose down

################################################################################
## Setup, build and push container for AWS
################################################################################

.PHONY: setup-aws
setup-aws: ## Setup AWS
	./script/setup_aws.sh

.PHONY: auth-docker-aws
auth-docker-aws: ## Authenticate docker to push images to AWS
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_CONTAINER_REGISTRY)

.PHONY: build-container-aws
build-container-aws: ## Build container for AWS
	make .build-container registry=$(AWS_CONTAINER_REGISTRY) tag=dev0

.PHONY: push-container-aws
push-container-aws: auth-docker-aws ## Push container for AWS
	make .push-container registry=$(AWS_CONTAINER_REGISTRY) tag=dev0

################################################################################
## Setup, build and push container for GCP
################################################################################

.PHONY: setup-gcp
setup-gcp: ## Set up GCP
	./script/setup_gcp.sh

.PHONY: auth-docker-gcp
auth-docker-gcp: ## Authenticate docker to push images to GCP
	gcloud auth configure-docker $(GCP_REGION)-docker.pkg.dev --quiet

.PHONY: build-container-gcp
build-container-gcp: ## Build container for GCP
	make .build-container registry=$(GCP_CONTAINER_REGISTRY) tag=dev0

.PHONY: push-container-gcp
push-container-gcp: auth-docker-gcp ## Push container to GCP
	make .push-container registry=$(GCP_CONTAINER_REGISTRY) tag=dev0

################################################################################
## Setup, build and push container for Azure
################################################################################

.PHONY: setup-azure
setup-azure: ## Set up Azure
	./script/setup_azure.sh

.PHONY: auth-docker-azure
auth-docker-azure: ## Authenticate docker to push images to Azure
	az acr login --name $(ACR_DOMAIN)

.PHONY: build-container-azure
build-container-azure: ## Build container for azure
	make .build-container registry=$(ACR_DOMAIN) tag=dev0

.PHONY: push-container-azure
push-container-azure: auth-docker-azure ## Push container to azure
	make .push-container registry=$(ACR_DOMAIN) tag=dev0

################################################################################
## Build and push container internal
################################################################################

CONTAINER_PLATFORMS = linux/amd64,linux/arm64
BUILDX_BUILDER = diarkis-multiarch-builder
BUILDX_ARGS = --builder $(BUILDX_BUILDER) --platform $(CONTAINER_PLATFORMS) --provenance=false

.PHONY: setup-buildx
setup-buildx: ## Create the buildx builder for multi-platform container images
	docker buildx inspect $(BUILDX_BUILDER) >/dev/null 2>&1 || \
		docker buildx create --name $(BUILDX_BUILDER) --driver docker-container

.PHONY: .build-container
.build-container: setup-buildx build-linux-amd64 build-linux-arm64 ## Build multi-platform container images into the buildx cache
	docker buildx build $(BUILDX_ARGS) --output type=cacheonly -f docker/http/Dockerfile remote_bin/linux -t $(registry)/http:$(tag)
	docker buildx build $(BUILDX_ARGS) --output type=cacheonly -f docker/udp/Dockerfile remote_bin/linux -t $(registry)/udp:$(tag)
	docker buildx build $(BUILDX_ARGS) --output type=cacheonly -f docker/tcp/Dockerfile remote_bin/linux -t $(registry)/tcp:$(tag)
	docker buildx build $(BUILDX_ARGS) --output type=cacheonly -f docker/mars/Dockerfile remote_bin/linux -t $(registry)/mars:$(tag)

.PHONY: .push-container
.push-container: setup-buildx ## Push multi-platform container images built by .build-container
	docker buildx build $(BUILDX_ARGS) --push -f docker/http/Dockerfile remote_bin/linux -t $(registry)/http:$(tag)
	docker buildx build $(BUILDX_ARGS) --push -f docker/udp/Dockerfile remote_bin/linux -t $(registry)/udp:$(tag)
	docker buildx build $(BUILDX_ARGS) --push -f docker/tcp/Dockerfile remote_bin/linux -t $(registry)/tcp:$(tag)
	docker buildx build $(BUILDX_ARGS) --push -f docker/mars/Dockerfile remote_bin/linux -t $(registry)/mars:$(tag)
