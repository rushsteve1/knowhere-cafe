set dotenv-load

@default:
    just --list --unsorted

@deps:
    go mod tidy

@build: deps
    CGO_ENABLED=0 go build -o knowhere-cafe .

@lint: deps
    go vet

@test: deps
    go test .

@run: deps
    CGO_ENABLED=0 go run . postgres://postgres:postgres@localhost:5432/knowhere

@watch:
    ls -d **/*.{go,html,js} | entr just run

@fmt:
    (which golines >> /dev/null && golines -w **/*.go) || go fmt
