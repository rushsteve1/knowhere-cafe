set dotenv-load

default: build

deps:
    go mod tidy

build: deps
    CGO_ENABLED=0 go build -o knowhere-cafe .

lint: deps
    go vet

test: deps
    go test .

watch: deps
    CGO_ENABLED=0 go run . -- postgres://postgres:postgres@localhost:5432/knowhere