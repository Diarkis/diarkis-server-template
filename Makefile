PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Set default command of make to help, so that running make will output help texts
.DEFAULT_GOAL := help

COPYRIGHT := // © 2019-2024 Diarkis Inc. All rights reserved.

.PHONY: init
init: ## make init project_id={project ID} builder_token={build token} output={absolute path to install}
	./init.sh $(project_id) $(builder_token) $(output)

.PHONY: fmt
fmt: add-license
	gofmt -w src/
	npx prettier --write "**/*.{yml,yaml,json,md}"

.PHONY: add-license
add-license: $(shell find . -type f -name '*.go')
	for f in $^; do \
		head -n 1 "$$f" | grep -q '^$(COPYRIGHT)' && continue; \
		echo $(COPYRIGHT) | cat - "$$f" > temp && mv temp "$$f"; \
	done
