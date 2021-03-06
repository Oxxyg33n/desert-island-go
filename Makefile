mod-vendor:
	go mod tidy && go mod vendor

clean:
	rm -Rf ./output/images/* && rm -Rf ./output/traits/*

build:
	@go build

test:
	@go list ./... | xargs go test

run: build
	@./desert-island-go

validations:
	@golangci-lint run
