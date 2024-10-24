PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

COPYRIGHT := Diarkis Inc. All rights reserved.

.PHONY: init
init: ## make init project_id={project ID} builder_token={build token} output={absolute path to install} module_name={go module name}
	go run ./tools/init.go $(project_id) $(builder_token) $(output) $(module_name)

.PHONY: fmt
fmt: add-license
	gofmt -w src/
	npx prettier --write "**/*.{yml,yaml,json,md}"

.PHONY: add-license ## add license header to all go files
add-license: $(shell find . -type f -name '*.go')
	for f in $^; do \
		head -n 1 "$$f" | grep -q '$(COPYRIGHT)' && tail -n +2 "$$f" > temp && mv temp "$$f"; \
		head -n 1 "$$f" | grep -q '^$$' && tail -n +2 "$$f" > temp && mv temp "$$f"; \
		echo "// © 2019-$(shell date +%Y) $(COPYRIGHT)\n" | cat - "$$f" > temp && mv temp "$$f"; \
	done
