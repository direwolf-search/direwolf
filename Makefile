# Makefile of DireWolf project
SHELL = /bin/bash

# Prerequisites
###############

# dirs
SRC_DIR := internal
BIN_DIR := bin
DW_MIGRATIONS_BIN_DIR := migrate
DW_MIGRATIONS_SRC_DIR := cmd/migrate

DW_MIGRATIONS_BUILD_TARGET := dw-migrations

# migrate subcommands
DB_COMMAND := db
INIT_COMMAND := init
MIGRATE_COMMAND := migrate
ROLLBACK_COMMAND := rollback
LOCK_COMMAND := lock
UNLOCK_COMMAND := unlock
CREATE_GO_COMMAND := create_go
CREATE_SQL_COMMAND :=create_sql
STATUS_COMMAND := status
MARK_APPLIED := mark_applied

# version
REV_LIST := $(shell git rev-list --tags --max-count=1)
VERSION ?= $(shell git describe --tags $(REV_LIST)) # git describe --tags $(git rev-list --tags --max-count=1)

# replacements
TEST_SUFFIX := _test.go

# binaries
CHANGELOG := changelog
GOTESTS := gotests
GIT := git

# go compiler
GO := go build


# parse arguments for changelog-init target
ifeq (changelog-init,$(firstword $(MAKECMDGOALS)))
  CHANGELOG_INIT_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(CHANGELOG_INIT_ARGS):;@:)
endif

# parse arguments for push-bitbucket target
ifeq (push-bitbucket,$(firstword $(MAKECMDGOALS)))
  PUSH_BITBUCKET_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(PUSH_BITBUCKET_ARGS):;@:)
endif

# parse arguments for push-origin target
ifeq (push-origin,$(firstword $(MAKECMDGOALS)))
  PUSH_ORIGIN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(PUSH_ORIGIN_ARGS):;@:)
endif

# parse arguments for changelog-finalize target
ifeq (changelog-finalize,$(firstword $(MAKECMDGOALS)))
  CHANGELOG_FINALIZE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(CHANGELOG_FINALIZE_ARGS):;@:)
endif

# parse arguments for gotests-generate target
ifeq (gotests-generate,$(firstword $(MAKECMDGOALS)))
  GOTESTS_GENERATE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(GOTESTS_GENERATE_ARGS):;@:)
endif

# parse arguments for migrate-create-go target
ifeq (migrate-create-go,$(firstword $(MAKECMDGOALS)))
  MIGRATE_CREATE_GO_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(MIGRATE_CREATE_GO_ARGS):;@:)
endif

# Targets
##########

default: help

dummy-push-origin:
	# ...

dummy-push-bitbucket:
	# ...

dummy-generate-test:
	# ...

dummy-changelog-init:
	# ...

dummy-changelog-finalize:
	# ...

dummy-migrate-create-go:
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

.PHONY: clean #                     -- removes protoc code generation artifacts
clean:
	@rm -R gen

# changelog commands

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
	@changelog init --author "Alexey 'hIMEI' Matveev" --email "himei@tuta.io" --since $(CHANGELOG_INIT_ARGS)

.PHONY: changelog-prepare #         -- Prepares changelog for release.
changelog-prepare:
	@$(CHANGELOG) prepare --author "Alexey 'hIMEI' Matveev" --email "himei@tuta.io"

.PHONY: changelog-finalize #        -- Finalizes changelog with given version. Syntax: make changelog-finalize [version]
changelog-finalize: dummy-changelog-finalize
	@$(CHANGELOG) finalize --version=$(CHANGELOG_FINALIZE_ARGS)

.PHONY: changelog-out #             -- Creates changelog file in md format
changelog-out:
	@$(CHANGELOG) md --out=CHANGELOG.md

# gotests commands

.PHONY: gotests-check #             -- Checks if gotests installed
gotests-check:
	@ if ! which $(GOTESTS) > /dev/null; then \
		echo "error: gotests not installed" >&2; \
		echo "to install it run <make gotests_install>" >&2; \
		exit 1; \
	fi

.PHONY: gotests-generate #          -- Generates tests for given file
gotests-generate: dummy-generate-test
	@$(GOTESTS) -all -template testify $(SRC_DIR)/$(GOTESTS_GENERATE_ARGS) >> $(SRC_DIR)/$(basename $(GOTESTS_GENERATE_ARGS))$(TEST_SUFFIX)

# manage remotes commands

.PHONY: push-bitbucket #            -- Pushes to bitbucket remote
push-bitbucket: dummy-push-bitbucket
	@$(GIT) push -u bitbucket $(PUSH_BITBUCKET_ARGS)

.PHONY: push-origin #               -- Pushes to bitbucket remote
push-origin: dummy-push-origin
	@$(GIT) push -u origin $(PUSH_ORIGIN_ARGS)

# migrate commands

.PHONY: migrate-build
migrate-build:
	@$(GO) -o $(BIN_DIR)/$(DW_MIGRATIONS_BIN_DIR)/$(DW_MIGRATIONS_BUILD_TARGET) $(DW_MIGRATIONS_SRC_DIR)/main.go

.PHONY: migrate-create-go
migrate-create-go:
	@$(BIN_DIR)/$(DW_MIGRATIONS_BIN_DIR)/$(DW_MIGRATIONS_BUILD_TARGET) $(DB_COMMAND) $(CREATE_GO_COMMAND) $(MIGRATE_CREATE_GO_ARGS)

.PHONY: migrate-init
migrate-init:
	@$(BIN_DIR)/$(DW_MIGRATIONS_BIN_DIR)/$(DW_MIGRATIONS_BUILD_TARGET) $(DB_COMMAND) $(INIT_COMMAND)

.PHONY: migrate-migrate
migrate-migrate:
	@$(BIN_DIR)/$(DW_MIGRATIONS_BIN_DIR)/$(DW_MIGRATIONS_BUILD_TARGET) $(DB_COMMAND) $(MIGRATE_COMMAND)

.PHONY: migrate-rollback
migrate-rollback:
	@$(BIN_DIR)/$(DW_MIGRATIONS_BIN_DIR)/$(DW_MIGRATIONS_BUILD_TARGET) $(DB_COMMAND) $(ROLLBACK_COMMAND)

.PHONY: migrate-lock
migrate-lock:
	@$(BIN_DIR)/$(DW_MIGRATIONS_BIN_DIR)/$(DW_MIGRATIONS_BUILD_TARGET) $(DB_COMMAND) $(LOCK_COMMAND)

.PHONY: migrate-unlock
migrate-unlock:
	@$(BIN_DIR)/$(DW_MIGRATIONS_BIN_DIR)/$(DW_MIGRATIONS_BUILD_TARGET) $(DB_COMMAND) $(UNLOCK_COMMAND)

.PHONY: migrate-status
migrate-status:
	@$(BIN_DIR)/$(DW_MIGRATIONS_BIN_DIR)/$(DW_MIGRATIONS_BUILD_TARGET) $(DB_COMMAND) $(STATUS_COMMAND)

.PHONY: migrate-mark-applied
migrate-mark-applied:
	@$(BIN_DIR)/$(DW_MIGRATIONS_BIN_DIR)/$(DW_MIGRATIONS_BUILD_TARGET) $(DB_COMMAND) $(MARK_APPLIED)
