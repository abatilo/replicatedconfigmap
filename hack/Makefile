SHELL := /bin/bash

GO_VERSION = 1.13.4
TOOLS_CONTAINER = abatilo/rcm-tools

.PHONY: help
help: ## View help information
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build_tools
build_tools:
	docker build -t $(TOOLS_CONTAINER) -f hack/Dockerfile.tools .

.PHONY: tools
tools: build_tools ## Start a shell with tools installed and the repo directory mounted
	docker run -it --rm -v$(HOME)/.kube:/root/.kube:ro -v$(PWD):/src -w /src $(TOOLS_CONTAINER) bash
	# Take ownership of any new files that came from mounting
	sudo chown -R $(shell whoami):$(shell whoami) ../
