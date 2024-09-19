build:
	@go build -o bin/my-blogs cmd/web/*.go

run: build
	@./bin/my-blogs

test:
	@go test -v ./...
