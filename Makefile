#!/usr/bin/make -f

all: format lint test

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify
	go mod tidy
	@echo "--> Download go modules to local cache"
	go mod download

.PHONY: go.sum

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@echo "--> Running linter"
	@which golangci-lint > /dev/null || echo "\033[91m install golangci-lint ...\033[0m" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run -v --out-format=tab

format: format-goimports
	golangci-lint run --fix --out-format=tab --issues-exit-code=0

format-goimports:
	@go install github.com/incu6us/goimports-reviser/v3@latest
	@find . -name '*.go' -type f -not -name '*.pb.go' -not -name '*.pb.gw.go' -exec goimports-reviser -use-cache -rm-unused {} \;

.PHONY: format lint format-goimports

###############################################################################
###                                 Tests                                   ###
###############################################################################

test:
	@echo "--> Running tests"
	go test $(BUILD_FLAGS) -mod=readonly ./...

test-count:
	go test $(BUILD_FLAGS) -mod=readonly -cpu 1 -count 1 -cover ./...

.PHONY: test

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.13.5
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=docker run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint ./proto

proto-gen:
	@$(protoImage) sh ./proto/gen.sh

.PHONY: proto-all proto-format proto-lint proto-gen