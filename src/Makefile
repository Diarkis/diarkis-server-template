PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

DIARKIS_CLI = ./diarkis-cli/os/linux/bin/diarkis-cli

ifeq ($(shell uname), Darwin)
	DIARKIS_CLI = ./diarkis-cli/os/mac/bin/diarkis-cli
endif

.PHONY: build-local
build-local: gen ## Build server binary for local use
	echo "Build server binaries"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c local-build.yml --host v3.builder.diarkis.io
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
gen: # Generate go, cpp, css code from packet definition written in json
	make -C puffer gen

.PHONY: run-docker
run-docker: ## Run Diarkis locally with docker compose
	docker compose up -d

.PHONY: stop-docker
stop-docker: ## Stop Diarkis
	docker compose down

.PHONY: setup-aws
setup-aws:
	./script/setup_aws.sh

.PHONY: build-container-aws
build-container-aws: ## Build container for AWS
	$(DIARKIS_CLI) build -c container-build.yaml --host v3.builder.diarkis.io
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
setup-gcp: 
	./script/setup_gcp.sh

.PHONY: build-container-gcp
build-container-gcp: ## Build container for GCP
	$(DIARKIS_CLI) build -c container-build.yaml --host v3.builder.diarkis.io
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
setup-azure:
	./script/setup_azure.sh


.PHONY: build-container-azure
build-container-azure: ## Build container for azure
	$(DIARKIS_CLI) build -c container-build.yaml --host v3.builder.diarkis.io
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
