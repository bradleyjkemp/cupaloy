.PHONY: install
install: get_dependencies install_linters

.PHONY: get_dependencies
get_dependencies:
	go get github.com/mattn/goveralls

.PHONY: install_linters
install_linters:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.12.5

.PHONY: lint
lint:
	$(GOPATH)/bin/golangci-lint run

.PHONY: test
test: lint
	go test ./...


.PHONY: test-ci
test-ci: coverage lint

.PHONY: coverage
coverage:
	go test -v -coverpkg ./... -coverprofile coverage.out ./...

.PHONY: clean
clean:
	rm -rf examples/ignored*
