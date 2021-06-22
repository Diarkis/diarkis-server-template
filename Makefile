PHONY: install
install:
	go run ./tools/install.go $(project) $(output)

.DEFAULT_GOAL := install
