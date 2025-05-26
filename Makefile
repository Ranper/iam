# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: help

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/Ranper/iam
VERSION_PACKAGE=github.com/marmotedu/component-base/pkg/version

# ==============================================================================
# Includes

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
include scripts/make-rules/image.mk


## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
