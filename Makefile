PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

.PHONY: init
init: ## make init project_id={project ID} builder_token={build token} output={absolute path to install}
	go run ./tools/install.go $(project_id) $(builder_token) $(output)

.PHONY: puffer
puffer:
	./puffer/puffer_gen -lang go -output ./src/cmds/custom/ -definitions ./puffer/json_definitions -tpl ./puffer/tpl
