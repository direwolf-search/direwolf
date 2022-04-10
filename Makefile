SHELL = /bin/bash

# Definitions
#############

# protos generates gRPC API from .proto files

PROTOS_SOURCE_DIR := protos
PROTOS_TARGET_DIR := internal/protos
PROTOS_FF := ${shell find ${PROTOS_SOURCE_DIR} -maxdepth 1 -type f -print -name *.proto}

# openapi2protos target generates .proto files from OpenApi documentation

# parse arguments for openapi2protos target
ifeq (openapi2protos,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "openapi2protos"
  OPENAPI2PROTOS_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(OPENAPI2PROTOS_ARGS):;@:)
endif

GENERATED_FROM_OPENAPI_DIR := build/generated
OPENAPI_FILE_FULL_PATH := docs/openapi2protofiles/$(OPENAPI2PROTOS_ARGS)
YAML := .yaml
PROTO := .proto
GENERATED_FILE_NAME := $(subst $(YAML),$(PROTO),$(OPENAPI2PROTOS_ARGS))


# Targets
##########

openapi2protofiles:
	#...

.PHONY: protos

# generates .proto files from OpenApi documentation
openapi2protos: openapi2protofiles
	@openapi2proto -spec $(OPENAPI_FILE_FULL_PATH) -out $(GENERATED_FROM_OPENAPI_DIR)/$(GENERATED_FILE_NAME) -annotate

# generates gRPC API from .proto files
protos: $(PROTOS_FF)
	@echo $(PROTOS_FF)
	for FILE in $(PROTOS_FF); do \
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