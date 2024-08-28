PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

DIARKIS_VERSION = $(shell go list -m github.com/Diarkis/diarkis | cut -d " " -f 2)
DIARKIS_CLI = ./diarkis-cli/os/linux/bin/diarkis-cli
BUILD_CONFIG = ./build/linux-build.yml
AWS_ACCOUNT_NUM = __AWS_ACCOUNT_NUM__
GCP_PROJECT_ID = __GCP_PROJECT_ID__
ACR_DOMAIN = __ACR_DOMAIN__

ifeq ($(shell uname), Darwin)
	DIARKIS_CLI = ./diarkis-cli/os/mac/bin/diarkis-cli
	BUILD_CONFIG = ./build/mac-build.yml
endif

PHONY: init
init: ## Initialize project
	./script/change_diarkis_version.sh $(DIARKIS_VERSION)

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
run-docker: build-linux ## Run Diarkis locally with docker compose
	docker compose up -d

.PHONY: stop-docker
stop-docker: ## Stop Diarkis
	docker compose down

.PHONY: setup-aws
setup-aws:
	./script/setup_aws.sh

.PHONY: auth-docker-aws
auth-docker-aws:
	aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $(AWS_ACCOUNT_NUM).dkr.ecr.ap-northeast-1.amazonaws.com

.PHONY: build-container-aws
build-container-aws: build-linux ## Build container for AWS
	docker build --platform=linux/amd64 -f docker/http/Dockerfile remote_bin -t $(AWS_ACCOUNT_NUM).dkr.ecr.ap-northeast-1.amazonaws.com/http
	docker build --platform=linux/amd64 -f docker/udp/Dockerfile remote_bin -t $(AWS_ACCOUNT_NUM).dkr.ecr.ap-northeast-1.amazonaws.com/udp
	docker build --platform=linux/amd64 -f docker/tcp/Dockerfile remote_bin -t $(AWS_ACCOUNT_NUM).dkr.ecr.ap-northeast-1.amazonaws.com/tcp
	docker build --platform=linux/amd64 -f docker/mars/Dockerfile remote_bin -t $(AWS_ACCOUNT_NUM).dkr.ecr.ap-northeast-1.amazonaws.com/mars

.PHONY: push-container-aws
push-container-aws: auth-docker-aws ## Push container to AWS
	docker push $(AWS_ACCOUNT_NUM).dkr.ecr.ap-northeast-1.amazonaws.com/http
	docker push $(AWS_ACCOUNT_NUM).dkr.ecr.ap-northeast-1.amazonaws.com/udp
	docker push $(AWS_ACCOUNT_NUM).dkr.ecr.ap-northeast-1.amazonaws.com/tcp
	docker push $(AWS_ACCOUNT_NUM).dkr.ecr.ap-northeast-1.amazonaws.com/mars

.PHONY: setup-gcp
setup-gcp: ## Set up GCP
	./script/setup_gcp.sh

.PHONY: auth-docker-gcp
auth-docker-gcp: ## Authenticate docker to push images to GCP
	gcloud auth configure-docker asia-northeast1-docker.pkg.dev --quiet

.PHONY: build-container-gcp
build-container-gcp: build-linux ## Build container for GCP
	docker build --platform=linux/amd64 -f docker/http/Dockerfile remote_bin -t asia-northeast1-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis/http
	docker build --platform=linux/amd64 -f docker/udp/Dockerfile remote_bin -t  asia-northeast1-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis/udp
	docker build --platform=linux/amd64 -f docker/tcp/Dockerfile remote_bin -t  asia-northeast1-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis/tcp
	docker build --platform=linux/amd64 -f docker/mars/Dockerfile remote_bin -t  asia-northeast1-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis/mars

.PHONY: push-container-gcp
push-container-gcp: auth-docker-gcp ## Push container to GCP
	docker push  asia-northeast1-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis/http
	docker push  asia-northeast1-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis/udp
	docker push  asia-northeast1-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis/tcp
	docker push  asia-northeast1-docker.pkg.dev/$(GCP_PROJECT_ID)/diarkis/mars

.PHONY: setup-azure
setup-azure: ## Set up Azure
	./script/setup_azure.sh

.PHONY: build-container-azure
build-container-azure: build-linux ## Build container for azure
	docker build --platform=linux/amd64 -f docker/http/Dockerfile remote_bin -t $(ACR_DOMAIN)/http
	docker build --platform=linux/amd64 -f docker/udp/Dockerfile remote_bin -t  $(ACR_DOMAIN)/udp
	docker build --platform=linux/amd64 -f docker/tcp/Dockerfile remote_bin -t  $(ACR_DOMAIN)/tcp
	docker build --platform=linux/amd64 -f docker/mars/Dockerfile remote_bin -t $(ACR_DOMAIN)/mars

.PHONY: push-container-azure
push-container-azure: ## Push container to azure
	docker push $(ACR_DOMAIN)/http
	docker push $(ACR_DOMAIN)/udp
	docker push $(ACR_DOMAIN)/tcp
	docker push $(ACR_DOMAIN)/mars
