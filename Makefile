mod-vendor:
	go mod tidy && go mod vendor

build:
	@go build

test:
	@go list ./... | xargs go test

run: build
	@./desert-island-go