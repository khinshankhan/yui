# Set .SHELLFLAGS to use "bash strict mode". This will fail on any error, fail
# on an error in a pipeline (usually it just returns the value of the last
# command in a pipeline), and will fail if using any undefined variables.
SHELL := bash
# Set .ONESHELL config. Runs the whole make recipe in one shell session.
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
# Set .DELETE_ON_ERROR. Deletes any files generated on error.
.DELETE_ON_ERROR:
# Set makeflags to --warn-on-unused-variables (probably an error) and to avoid
# the built-in rules (this removes a lot of magic that is related to yacc and
# other tools).
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# VARIABLES
GIT_COMMIT = $(shell git rev-parse --verify HEAD)
BUILD_DATE = $(shell date +%Y.%m.%d.%H%M%S)

GOSERVICES = $(sort $(notdir $(realpath $(dir $(wildcard ./cmd/*/main.go)))))

.PHONY: $(GOSERVICES)

.PHONY: build
build: $(GOSERVICES)

# creates a TARGET per go service
$(GOSERVICES): % : ./cmd/%/main.go
	go build -ldflags="-s -w -X main.CommitHash=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)" -o bin/$@ cmd/$@/*.go

# list available go services
.PHONY: services
services:
	@echo $(GOSERVICES)
