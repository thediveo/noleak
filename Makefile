.PHONY: clean coverage deploy undeploy help install test report buildapp startapp docsify scan

help: ## list available targets
	@# Shamelessly stolen from Gomega's Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

clean: ## clean up build and coverage artefacts
	rm -f coverage.html coverage.out

test: ## runs all tests
	ginkgo -r --randomize-all
