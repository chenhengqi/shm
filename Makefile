GO_VERSION=1.14.2-buster

.PHONY: test-local
test-local: lint dep
	go test -v -cover -coverprofile=coverage.out ./...

.PHONY: test
test: lint dep
	docker run -w /home -v $(PWD):/home golang:$(GO_VERSION) sh -c 'go test -v -cover -coverprofile=coverage.out ./...'

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...
	golint `go list ./... | grep -v /vendor/`

.PHONY: dep
dep:
	go mod tidy
