PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

.PHONY: install
install: ## make install project={name of your application} project_id={project ID} builder_token={build token} output={absolute path to install}
	go run ./tools/install.go $(project) $(project_id) $(builder_token) $(output)
