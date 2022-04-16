# Makefile of DireWolf project
SHELL = /bin/bash

# prerequisites
#############

# dirs
SRC_DIR := internal
DOMAIN_SERVICES_SRC_DIR := $(SRC_DIR)/domain/service
CONCRETE_SERVICES_SRC_DIR := $(SRC_DIR)/services
API_DIR = $(SRC_DIR)/api
PROTO_PACKAGES_DIR := protos
GENERATED_FROM_OPENAPI_DIR := build/generated
OPENAPI_FILES_DIR := docs/openapi2protofiles

# version
VERSION ?= $(shell git describe --tags --always --match=v* 2> /dev/null)

# replacements
YAML_SUFFIX := .yaml
PROTO_SUFFIX := *.proto
TEST_SUFFIX := _test.go
GRPC_SUFFIX := grpc

# binaries
PROTOC := protoc
CHANGELOG := changelog
GOTESTS := gotests
CONVERTER := openapi2proto

# parse arguments for changelog-init target
ifeq (changelog-init,$(firstword $(MAKECMDGOALS)))
  changelog_init_args := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(changelog_init_args):;@:)
endif

# parse arguments for changelog-finalize target
ifeq (changelog-finalize,$(firstword $(MAKECMDGOALS)))
  changelog_finalize_args := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(changelog_finalize_args):;@:)
endif

# parse arguments for generate-test target
ifeq (generate-test,$(firstword $(MAKECMDGOALS)))
  generate_test_args := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(generate_test_args):;@:)
endif

SERVICERS_API_DIR := $(shell find $(CONCRETE_SERVICES_SRC_DIR) -name "*.proto")
srcs := $(foreach d,$(SERVICERS_API_DIR),$(wildcard $(d)/*.proto))
objs := $(srcs:.cpp=.o)

# Targets
##########


default: help

dummy-generate-test:
	# ...

dummy-changelog-init:
	# ...

dummy-changelog-finalize:
	# ...

.PHONY: help #                      -- Shows help message
help:
	@echo ''
	@echo '==============================='
	@echo 'usage: make <TARGET> [ARGUMENT]'
	@echo '==============================='
	@echo ''
	@echo 'Targets: '
	@echo '--------'
	@echo ''
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1 \2/' | expand -t20
	@echo ''

.PHONY: version #                   -- Prints current application version, v0 if not found
version:
	@echo $(VERSION)

.PHONY: generate #                  -- Generates gRPC API from .proto file
generate: $(SERVICERS_API_DIR)
	for file in $^ ; do \
    	$(PROTOC) -I$(PROTO_PACKAGES_DIR) \
        	--proto_path=$$(dirname $${file%.*})/ \
        	--go_out=$(api_dir)/$$(basename $${file%.*}) \
			--go_opt=paths=source_relative \
        	--go-grpc_out=$(api_dir)/$$(basename $${file%.*}) \
        	--go-grpc_opt=paths=source_relative \
        	--grpc-gateway_out $(api_dir)/$$(basename $${file%.*}) \
        	--grpc-gateway_opt logtostderr=true \
        	--grpc-gateway_opt paths=source_relative \
        	$$file; \
    done

.PHONY: clean #                     -- removes protoc code generation artifacts
clean:
	@rm -R gen

.PHONY: convert #                   -- Converts OpenApi documentation in .yaml to .proto files
convert: $(OPENAPI_FILES_DIR)/*
	@echo 'File $^ will be converted to .proto file format'
	@$(CONVERTER) \
	-spec $^ \
	-out $(subst $(YAML_SUFFIX),$(PROTO_SUFFIX),$(subst $(OPENAPI_FILES_DIR),$(GENERATED_FROM_OPENAPI_DIR), $^)) \
	-annotate \

.PHONY: changelog-check #           -- Checks if changelog installed
changelog-check:
	@ if ! which $(CHANGELOG) > /dev/null; then \
		echo "error: changelog not installed" >&2; \
		echo "to install it run <make changelog_install>" >&2; \
		exit 1; \
	fi

.PHONY: changelog-install #         -- Installs changelog
changelog-install:
	@go get github.com/mh-cbon/changelog
	@go install github.com/mh-cbon/changelog
	@$(CHANGELOG)
	@echo "mh-cbon/changelog installed if you see its usage message"


.PHONY: changelog-init #            -- Initialize changelog file for project. Syntax: make changelog-init [version]
changelog-init: dummy-changelog-init
	@changelog init --author "Alexey 'hIMEI' Matveev" --email "himei@tuta.io" --since $(changelog_init_args)

.PHONY: changelog-prepare #         -- Prepares changelog for release.
changelog-prepare:
	@$(CHANGELOG) prepare --author "Alexey 'hIMEI' Matveev" --email "himei@tuta.io"

.PHONY: changelog-finalize #        -- Finalizes changelog with given version. Syntax: make changelog-finalize [version]
changelog-finalize: dummy-changelog-finalize
	@$(CHANGELOG) finalize --version=$(changelog_finalize_args)

.PHONY: changelog-out #             -- Creates changelog file in md format
changelog-out:
	@$(CHANGELOG) md --out=CHANGELOG.md

.PHONY: gotests-check #             -- Checks if gotests installed
gotests-check:
	@ if ! which $(GOTESTS) > /dev/null; then \
		echo "error: gotests not installed" >&2; \
		echo "to install it run <make gotests_install>" >&2; \
		exit 1; \
	fi

.PHONY: generate-test #             -- Generates tests for given file
generate-test: dummy-generate-test
	@$(GOTESTS) -all $(generate_test_args) >> $(basename $(generate_test_args))$(TEST_SUFFIX)


