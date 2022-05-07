mod-vendor:
	go mod tidy && go mod vendor

build:
	@go build

test:
	@go list ./... | xargs go test

run: build
	@./metadata_updater

validations:
	@golangci-lint run
