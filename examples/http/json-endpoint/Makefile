CURRENT_MAKEFILE_DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))

PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

ROOT_DIR=${CURRENT_MAKEFILE_DIR}../../../
DIARKIS_CLI =${ROOT_DIR}/diarkis-cli/os/linux/bin/diarkis-cli
BUILD_CONFIG = linux-build.yml
ifeq ($(shell uname), Darwin)
	DIARKIS_CLI = ${ROOT_DIR}/diarkis-cli/os/mac/bin/diarkis-cli
	BUILD_CONFIG = mac-build.yml
endif

.PHONY: build-local
build-local: ## Build server binary for local use
	echo "Build server binaries for local use"
	rm -rfv remote_bin/*
	$(DIARKIS_CLI) build -c build/$(BUILD_CONFIG) --host v3.builder.diarkis.io
	chmod -R 700 remote_bin/
	echo "Build other diarkis dependencies"
	make -C $(ROOT_DIR) build-local

.PHONY: server
server: ## Start a server locally: [ target=mars ] [ target=http ]
	echo "Starting $(target) server..."
ifeq ($(target), mars)
	$(ROOT_DIR)/remote_bin/$(target) $(ROOT_DIR)/configs/mars/main.json
else
	DIARKIS_CONFIG_ROOT_PATH=../../.. $(CURRENT_MAKEFILE_DIR)/remote_bin/$(target)
endif
