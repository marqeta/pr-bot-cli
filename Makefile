GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

.PHONY: ci
ci: fmt vet lint build test

.PHONY: clean
clean:
	rm bin/* || true

.PHONY: vet
vet:
	go vet --mod=vendor ./...

.PHONY: lint
lint:
	golangci-lint run -v -c .golangci.yaml ./...

.PHONY: test
test:
	go test --mod=vendor -v -coverpkg=./... -coverprofile=coverage.out ./...

.PHONY: build
build: clean
	go build --mod=vendor -o ./bin/pr-bot-cli cmd/pr-bot-cli/main.go

.PHONY: fmt
fmt:
	go fmt --mod=vendor ./...
	goimports -l -w $(GOFMT_FILES)

.PHONY: run
run: build
	./bin/pr-bot -config ./config/local.yaml

.PHONY: mocks
mocks:
	go generate ./...

.PHONY: static
static:
	templ generate -path ui/templates
