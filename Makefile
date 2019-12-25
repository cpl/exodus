# add all your cmd/<things> in here
TARGETS = executable


CMD_DIR := ./cmd
PKG_DIR := ./pkg
OUT_DIR := ./out

COV_FILE := cover.out

GO111MODULE := on


GO_TEST_FLAGS := -v -count=1 -race -coverprofile=$(OUT_DIR)/$(COV_FILE) -covermode=atomic

.PHONY: $(OUT_DIR) clean build test mod cover test-deps fmt vet purge bench

all: clean mod fmt vet test build

$(OUT_DIR):
	@mkdir -p $(OUT_DIR)

clean:
	@rm -rf $(OUT_DIR)

purge: clean
	go mod tidy
	go clean -cache
	go clean -testcache
	go clean -modcache

build: $(OUT_DIR)
	$(foreach target,$(TARGETS),go build -o $(OUT_DIR)/$(target) $(CMD_DIR)/$(target)/*.go;)

test: $(OUT_DIR)
	go test $(GO_TEST_FLAGS) ./...

mod:
	go mod tidy
	go mod verify

cover:
	go tool cover -html=$(OUT_DIR)/$(COV_FILE)

test-deps:
	go test all

fmt:
	go fmt ./...

vet:
	go vet ./...

bench:
	go test -bench=. -benchmem -benchtime=10s ./...
