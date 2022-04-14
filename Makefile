SHELL = /bin/bash

# Definitions
#############

SRC_DIR := internal
DOMAIN_SERVICES_SRC_DIR := $(SRC_DIR)/domain/service
CONCRETE_SERVICES_SRC_DIR := $(SRC_DIR)/services

PROTOS_SOURCE_DIR := protos
PROTOS_TARGET_DIR := internal/protos
# list of .proto files in its source directory
PROTOS_FF := ${shell find ${PROTOS_SOURCE_DIR} -maxdepth 1 -type f -print -name *.proto}
SERVICES_DIR := internal/domain/service

SERVICES_FF := ${shell find ${SERVICES_DIR} -maxdepth 3 -type f -print -name *.proto}

VERSION ?= $(shell git describe --tags --always --match=v* 2> /dev/null)

GENERATED_FROM_OPENAPI_DIR := build/generated
OPENAPI_FILES_DIR := docs/openapi2protofiles
YAML_EXT := .yaml
PROTO_EXT := .proto

SOURCES = $(SERVICES_DIR)/$(wildcard *.proto) $(wildcard */*.proto)

# parse arguments for changelog-init target
ifeq (changelog-init,$(firstword $(MAKECMDGOALS)))
  CHANGELOG_INIT_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(CHANGELOG_INIT_ARGS):;@:)
endif

# parse arguments for changelog-finalize target
ifeq (changelog-finalize,$(firstword $(MAKECMDGOALS)))
  CHANGELOG_FINALIZE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(CHANGELOG_FINALIZE_ARGS):;@:)
endif

# Targets
##########

default: help

dummy-changelog-init:
	# ...

dummy-changelog-finalize:
	# ...

.PHONY: help #                -- Shows help message
help:
	@echo ''
	@echo '========================='
	@echo 'usage: make [target] ...'
	@echo '========================='
	@echo ''
	@echo 'Targets: '
	@echo '--------'
	@echo ''
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1 \2/' | expand -t20
	@echo ''

.PHONY: version #             -- Prints current application version, v0 if not found
version:
	@echo $(VERSION)

.PHONY: generate-test #       -- Not a tests realy!
generate-test: $(SOURCES)
	@echo $^

.PHONY: generate #            -- Generates gRPC API from .proto file
generate: $(PROTOS_FF)
	for FILE in $(PROTOS_FF); do \
  		echo "API for $${FILE} generated"; \
  		protoc \
            --proto_path=$(PROTOS_SOURCE_DIR)/ \
            --go_out=$(PROTOS_TARGET_DIR)/$$(basename $${FILE%.*}) \
            --go_opt=paths=source_relative \
            --go-grpc_out=$(PROTOS_TARGET_DIR)/$$(basename $${FILE%.*}) \
            --go-grpc_opt=paths=source_relative \
            --grpc-gateway_out $(PROTOS_TARGET_DIR)/$$(basename $${FILE%.*}) \
            --grpc-gateway_opt logtostderr=true \
            --grpc-gateway_opt paths=source_relative \
            $$FILE; \
  	done

.PHONY: clean #               -- removes protoc code generation artifacts
clean:
	@rm -R gen

test: $(SERVICES_FF)
	for FILE in $(SERVICES_FF); do \
		echo $$(dirname $${FILE}); \
	done

.PHONY: convert #             -- Converts .proto files from OpenApi documentation
convert: $(OPENAPI_FILES_DIR)/*
	@echo 'File $^ will be converted to .proto file format'
	@openapi2proto \
	-spec $^ \
	-out $(subst $(YAML_EXT),$(PROTO_EXT),$(subst $(OPENAPI_FILES_DIR),$(GENERATED_FROM_OPENAPI_DIR), $^)) \
	-annotate \

.PHONY: changelog #           -- Checks if changelog installed
changelog:
	@ if ! which changelog > /dev/null; then \
		echo "error: changelog not installed" >&2; \
		echo "to install it run <make changelog_install>" >&2; \
		exit 1; \
	fi

.PHONY: changelog-install #   -- Installs changelog
changelog-install:
	@go install github.com/mh-cbon/changelog
	@go get github.com/mh-cbon/changelog
	@changelog
	@echo "mh-cbon/changelog installed if you see its usage message"


.PHONY: changelog-init #      -- Initialize changelog file for project. Syntax: make changelog-init [VERSION]
changelog-init: dummy-changelog-init
	@changelog init --author "Alexey 'hIMEI' Matveev" --email "himei@tuta.io" --since $(CHANGELOG_INIT_ARGS)

.PHONY: changelog-prepare #   -- Prepares changelog for release.
changelog-prepare:
	@changelog prepare --author "Alexey 'hIMEI' Matveev" --email "himei@tuta.io"

.PHONY: changelog-finalize #  -- Finalizes changelog with given version. Syntax: make changelog-finalize [VERSION]
changelog-finalize: dummy-changelog-finalize
	@changelog finalize --version=$(CHANGELOG_FINALIZE_ARGS)

.PHONY: changelog-out #       -- Creates changelog file in md format
changelog-out:
	@changelog md --out=CHANGELOG.md