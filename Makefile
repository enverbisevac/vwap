GOLANGCI_LINT_VERSION=v1.44.0

ifndef GOPATH
	GOPATH := $(shell go env GOPATH)
endif

format: $(GOPATH)/bin/goimports
	@echo "Checking that go fmt does not make any changes..."
	@test -z $$(go fmt $(go list ./...)) || (echo "go fmt would make a change. Please verify and commit the proposed changes"; exit 1)
	@echo "Checking go fmt complete"
	@echo "Running goimports"
	@test -z $$(goimports -w ./..) || (echo "goimports would make a change. Please verify and commit the proposed changes"; exit 1)

lint: format $(GOPATH)/bin/golangci-lint
	@echo "Linting $(1)"
	@golangci-lint run \
		-E asciicheck \
		-E bodyclose \
		-E exhaustive \
		-E exportloopref \
		-E gofmt \
		-E goimports \
		-E gosec \
		-E noctx \
		-E nolintlint \
		-E exportloopref \
		-E stylecheck \
		-E unconvert \
		-E unparam
	@echo "Lint-free"

test:
	@go test -race ./...

test-unit:
	@go test -race ./internal/vwap

test-integration:
	@go test -race ./internal/websocket

upgrade:
	@echo "Upgrading dependencies..."
	@go get -u
	@go mod tidy
	
run: $(GOPATH)/bin/easyjson
	@go generate ./...
	@go run main.go

build:
	@go build -o vwap main.go

clean:
	@rm -rf wvap

#
# Install Tools
#
sec: $(GOPATH)/bin/gosec
	@echo "Checking for security problems ..."
	@gosec -quiet -confidence high -severity medium ./...
	@echo "No problems found"; \

$(GOPATH)/bin/golangci-lint:
	@echo "🔘 Installing golangci-lint... (`date '+%H:%M:%S'`)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin

$(GOPATH)/bin/goimports:
	@echo "🔘 Installing goimports ... (`date '+%H:%M:%S'`)"
	@go install -u golang.org/x/tools/cmd/goimports

$(GOPATH)/bin/easyjson:
	@echo "🔘 Installing easyjson ... (`date '+%H:%M:%S'`)"
	@go install github.com/mailru/easyjson/...@latest

$(GOPATH)/bin/gosec:
	@echo "🔘 Installing gosec ... (`date '+%H:%M:%S'`)"
	@curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(GOPATH)/bin