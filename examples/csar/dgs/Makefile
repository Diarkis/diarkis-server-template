PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

DIARKIS_VERSION = $(shell go mod edit -json | jq -r '.Require[] | select(.Path | contains("github.com/Diarkis/diarkis")).Version')
DIARKIS_CLI = ../../diarkis-cli/os/linux/bin/diarkis-cli
BUILD_CONFIG = ./build/linux-build.yml
AWS_ACCOUNT_NUM = __AWS_ACCOUNT_NUM__
AWS_REGION = ap-northeast-1
AWS_CONTAINER_REGISTRY = $(AWS_ACCOUNT_NUM).dkr.ecr.$(AWS_REGION).amazonaws.com
GCP_PROJECT_ID = __GCP_PROJECT_ID__
GCP_REGION = asia-northeast1
GCP_CONTAINER_REGISTRY = $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis
ACR_DOMAIN = __ACR_DOMAIN__

ifeq ($(shell uname), Darwin)
	DIARKIS_CLI = ../../diarkis-cli/os/mac/bin/diarkis-cli
	BUILD_CONFIG = ./build/mac-build.yml
endif

.PHONY: build-local
build-local: ## Build server binary for local use
	echo "Build server binaries for local use"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c $(BUILD_CONFIG) --host v3.builder.diarkis.io
	chmod -R 700 remote_bin/

.PHONY: build-linux
build-linux: ## Build server binary for linux or container environment
	echo "Build server binaries"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c ./build/linux-build.yml --host v3.builder.diarkis.io
	chmod -R 700 remote_bin/

.PHONY: build-mac
build-mac: ## Build server binary for mac use
	echo "Build server binaries"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c ./build/mac-build.yml --host v3.builder.diarkis.io
	chmod -R 700 remote_bin/

.PHONY: server
server: ## Start a server locally: [ target=mars ] [ target=http ] [ target=udp ] [ target=dgs ]
	echo "Starting $(target) server..."

ifeq ($(target), mars)
	./remote_bin/$(target) ./configs/mars/main.json
else
	./remote_bin/$(target)
endif

GOCLI_ARGS =
ifdef host
$(eval GOCLI_ARGS += --host=$(host))
else
$(eval GOCLI_ARGS += --host=127.0.0.1:7000)
endif
ifdef uid
$(eval GOCLI_ARGS += --uid=$(uid))
else
$(eval GOCLI_ARGS += --uid=$(shell echo $$RANDOM))
endif
ifdef clientKey
$(eval GOCLI_ARGS += --clientKey=$(clientKey))
endif
ifdef puffer
$(eval GOCLI_ARGS += --puffer=$(puffer))
endif

.PHONY: go-cli
go-cli: ## Starts Go test client: host=<HTTP address> uid=<client user ID> clientKey=<client key> puffer=<true/false>
	./remote_bin/cli $(GOCLI_ARGS)

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

.PHONY: build-container-azure
build-container-azure: ## Build container for azure
	make .build-container registry=$(ACR_DOMAIN) tag=dev0

.PHONY: push-container-azure
push-container-azure: ## Push container to azure
	make .push-container registry=$(ACR_DOMAIN) tag=dev0

################################################################################
## Build and push container internal
################################################################################

.PHONY: .build-container
.build-container: build-linux ## Build container internal
	docker build --platform=linux/amd64 -f docker/http/Dockerfile remote_bin -t $(registry)/http:$(tag)
	docker build --platform=linux/amd64 -f docker/udp/Dockerfile remote_bin -t $(registry)/udp:$(tag)
	docker build --platform=linux/amd64 -f docker/mars/Dockerfile remote_bin -t $(registry)/mars:$(tag)
	docker build --platform=linux/amd64 -f docker/dgs/Dockerfile . -t $(registry)/dgs:$(tag)

.PHONY: .push-container
.push-container: ## Push container internal
	docker push $(registry)/http:$(tag)
	docker push $(registry)/udp:$(tag)
	docker push $(registry)/mars:$(tag)
	docker push $(registry)/dgs:$(tag)
