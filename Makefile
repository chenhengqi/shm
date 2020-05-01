.PHONY: test
test: lint dep
	go test -v -cover -coverprofile=coverage.out ./...

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...
	golint `go list ./... | grep -v /vendor/`

.PHONY: dep
dep:
	go mod tidy
