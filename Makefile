SHELL = /bin/bash

# Definitions
#############

src_dir := internal
domain_services_src_dir := $(src_dir)/domain/service
concrete_services_src_dir := $(src_dir)/services

protos_source_dir := protos
protos_target_dir := internal/protos
# list of .api files in its source directory
protos_ff := ${shell find ${protos_source_dir} -maxdepth 1 -type f -print -name *.proto}
services_dir := internal/domain/service

concrete_services_dir = ./internal/services

services_ff := ${shell find ${services_dir} -maxdepth 3 -type f -print -name *.proto}

version ?= $(shell git describe --tags --always --match=v* 2> /dev/null)

generated_from_openapi_dir := build/generated
open_api_files_dir := docs/openapi2protofiles
yaml_text := .yaml
proto_text := .proto

test_suffix := _test.go

sources = $(services_dir)/$(wildcard *.proto) $(wildcard */*.proto)

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

# parse arguments for changelog-finalize target
ifeq (generate-test,$(firstword $(MAKECMDGOALS)))
  generate_test_args := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(generate_test_args):;@:)
endif

dirs := $(shell find $(concrete_services_dir) -type d -name "api")
srcs := $(foreach d,$(dirs),$(wildcard $(d)/*.proto))
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
	@echo $(version)

.PHONY: t-test #              -- Not a tests realy!
t-test: $(sources)
	@echo $^

.PHONY: generate #            -- Generates gRPC API from .api file
generate: $(protos_ff)
	for FILE in $(protos_ff); do \
  		echo "API for $${FILE} generated"; \
  		protoc \
            --proto_path=$(protos_source_dir)/ \
            --go_out=$(protos_target_dir)/$$(basename $${FILE%.*}) \
            --go_opt=paths=source_relative \
            --go-grpc_out=$(protos_target_dir)/$$(basename $${FILE%.*}) \
            --go-grpc_opt=paths=source_relative \
            --grpc-gateway_out $(protos_target_dir)/$$(basename $${FILE%.*}) \
            --grpc-gateway_opt logtostderr=true \
            --grpc-gateway_opt paths=source_relative \
            $$FILE; \
  	done

.PHONY: clean #               -- removes protoc code generation artifacts
clean:
	@rm -R gen

test: $(services_ff)
	for FILE in $(services_ff); do \
		echo $$(dirname $${FILE}); \
	done

.PHONY: convert #             -- Converts .api files from OpenApi documentation
convert: $(open_api_files_dir)/*
	@echo 'File $^ will be converted to .proto file format'
	@openapi2proto \
	-spec $^ \
	-out $(subst $(yaml_text),$(proto_text),$(subst $(open_api_files_dir),$(generated_from_openapi_dir), $^)) \
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


.PHONY: changelog-init #      -- Initialize changelog file for project. Syntax: make changelog-init [version]
changelog-init: dummy-changelog-init
	@changelog init --author "Alexey 'hIMEI' Matveev" --email "himei@tuta.io" --since $(changelog_init_args)

.PHONY: changelog-prepare #   -- Prepares changelog for release.
changelog-prepare:
	@changelog prepare --author "Alexey 'hIMEI' Matveev" --email "himei@tuta.io"

.PHONY: changelog-finalize #  -- Finalizes changelog with given version. Syntax: make changelog-finalize [version]
changelog-finalize: dummy-changelog-finalize
	@changelog finalize --version=$(changelog_finalize_args)

.PHONY: changelog-out #       -- Creates changelog file in md format
changelog-out:
	@changelog md --out=CHANGELOG.md

.PHONY: gotests #             -- Checks if gotests installed
gotests:
	@ if ! which gotests > /dev/null; then \
		echo "error: gotests not installed" >&2; \
		echo "to install it run <make gotests_install>" >&2; \
		exit 1; \
	fi

.PHONY: generate-test #       -- Generates tests for given fikle
generate-test: dummy-generate-test
	@gotests -all $(generate_test_args) >> $(basename $(generate_test_args))$(test_suffix)

.PHONY: traverse
traverse: $(dirs)
	for d in $+; do \
		echo "$$d"; \
	done