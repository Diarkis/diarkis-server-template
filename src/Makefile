PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

DIARKIS_CLI = ./diarkis-cli/os/linux/bin/diarkis-cli
BUILD_CONFIG = ./build/linux-build.yml

ifeq ($(shell uname), Darwin)
	DIARKIS_CLI = ./diarkis-cli/os/mac/bin/diarkis-cli
	BUILD_CONFIG = ./build/mac-build.yml
endif

.PHONY: build-local
build-local: gen ## Build server binary for local use
	echo "Build server binaries for local use"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c $(BUILD_CONFIG) --host v3.builder.diarkis.io
	chmod -R 700 remote_bin/

.PHONY: build-linux
build-linux: gen ## Build server binary for linux or container environment
	echo "Build server binaries"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c ./build/linux-build.yml --host v3.builder.diarkis.io
	chmod -R 700 remote_bin/

.PHONY: build-mac
build-mac: gen ## Build server binary for mac use
	echo "Build server binaries"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c ./build/mac-build.yml --host v3.builder.diarkis.io
	chmod -R 700 remote_bin/


.PHONY: server
server: ## Start a server locally: [ target=mars ] [ target=http ] [ target=udp ] [ target=tcp ]
	echo "Starting $(target) server..."

ifeq ($(target), mars)
	./remote_bin/$(target) ./configs/mars/main.json
else
	./remote_bin/$(target)
endif

.PHONY: gen
gen: ## Generate go, cpp, and cs code files using puffer (Diarkis packet gen module) from packet definition written in json
	make -C puffer gen

.PHONY: clean
clean: ## Deletes all puffer generated code files
	make -C puffer clean

.PHONY: run-docker
run-docker: build-linux ## Run Diarkis locally with docker compose
	docker compose up -d

.PHONY: stop-docker
stop-docker: ## Stop Diarkis
	docker compose down

.PHONY: setup-aws
setup-aws:
	./script/setup_aws.sh

.PHONY: build-container-aws
build-container-aws: build-linux ## Build container for AWS
	docker build --platform=linux/amd64 -f docker/http/Dockerfile remote_bin -t __AWS_ACCOUNT_NUM__.dkr.ecr.ap-northeast-1.amazonaws.com/http
	docker build --platform=linux/amd64 -f docker/udp/Dockerfile remote_bin -t __AWS_ACCOUNT_NUM__.dkr.ecr.ap-northeast-1.amazonaws.com/udp
	docker build --platform=linux/amd64 -f docker/tcp/Dockerfile remote_bin -t __AWS_ACCOUNT_NUM__.dkr.ecr.ap-northeast-1.amazonaws.com/tcp
	docker build --platform=linux/amd64 -f docker/mars/Dockerfile remote_bin -t __AWS_ACCOUNT_NUM__.dkr.ecr.ap-northeast-1.amazonaws.com/mars

.PHONY: push-container-aws
push-container-aws: ## Push container to AWS
	docker push __AWS_ACCOUNT_NUM__.dkr.ecr.ap-northeast-1.amazonaws.com/http
	docker push __AWS_ACCOUNT_NUM__.dkr.ecr.ap-northeast-1.amazonaws.com/udp
	docker push __AWS_ACCOUNT_NUM__.dkr.ecr.ap-northeast-1.amazonaws.com/tcp
	docker push __AWS_ACCOUNT_NUM__.dkr.ecr.ap-northeast-1.amazonaws.com/mars

.PHONY: setup-gcp
setup-gcp: ## Set up GCP
	./script/setup_gcp.sh

.PHONY: build-container-gcp
build-container-gcp: build-linux ## Build container for GCP
	docker build --platform=linux/amd64 -f docker/http/Dockerfile remote_bin -t asia.gcr.io/__GCP_PROJECT_ID__/http
	docker build --platform=linux/amd64 -f docker/udp/Dockerfile remote_bin -t asia.gcr.io/__GCP_PROJECT_ID__/udp
	docker build --platform=linux/amd64 -f docker/tcp/Dockerfile remote_bin -t asia.gcr.io/__GCP_PROJECT_ID__/tcp
	docker build --platform=linux/amd64 -f docker/mars/Dockerfile remote_bin -t asia.gcr.io/__GCP_PROJECT_ID__/mars

.PHONY: push-container-gcp
push-container-gcp: ## Push container to GCP
	docker push asia.gcr.io/__GCP_PROJECT_ID__/http
	docker push asia.gcr.io/__GCP_PROJECT_ID__/udp
	docker push asia.gcr.io/__GCP_PROJECT_ID__/tcp
	docker push asia.gcr.io/__GCP_PROJECT_ID__/mars

.PHONY: setup-azure
setup-azure: ## Set up Azure
	./script/setup_azure.sh


.PHONY: build-container-azure
build-container-azure: build-linux ## Build container for azure
	docker build --platform=linux/amd64 -f docker/http/Dockerfile remote_bin -t __ACR_DOMAIN__/http
	docker build --platform=linux/amd64 -f docker/udp/Dockerfile remote_bin -t  __ACR_DOMAIN__/udp
	docker build --platform=linux/amd64 -f docker/tcp/Dockerfile remote_bin -t  __ACR_DOMAIN__/tcp
	docker build --platform=linux/amd64 -f docker/mars/Dockerfile remote_bin -t __ACR_DOMAIN__/mars

.PHONY: push-container-azure
push-container-azure: ## Push container to azure
	docker push __ACR_DOMAIN__/http
	docker push __ACR_DOMAIN__/udp
	docker push __ACR_DOMAIN__/tcp
	docker push __ACR_DOMAIN__/mars
