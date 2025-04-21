# builds onyx
build:
    go build ./...

# runs checks
check:
    go fmt ./...
    go vet ./...

# tests onyx
test:
    go test ./...

# runs onyx
run:
    go run ./cmd/onyx

# starts server
serve:
    static-web-server -p 8080
