set dotenv-load

# Disable CGo globally
export CGO_ENABLED := "0"

DB_URL := "postgres://postgres:postgres@localhost:5432/knowhere"

@default:
	just --list --unsorted

@deps:
	go mod tidy
	go mod download

@build: deps
	go build -o knowhere-cafe .

@lint: deps
	go vet

@test: deps
	go test .

@run: deps
	go run . --dev {{DB_URL}}

watch:
	ls **/*.go | entr -c just run

alias fmt := format
@format:
	# try to use golines, fall back to go fmt
	(which golines >> /dev/null && golines -w -m 80 --ignore-generated **/*.go) || go fmt

@migrate:
	go run . --migrate {{DB_URL}}

alias dbg := debug
@debug:
	dlv debug -- --dev {{DB_URL}}


VENDOR_FOLDER := "static/vendor"

@vendor-unpoly:
	mkdir -p {{VENDOR_FOLDER}}
	curl  -o {{VENDOR_FOLDER}}/unpoly.min.js https://cdn.jsdelivr.net/npm/unpoly@3.9.2/unpoly.min.js
	curl  -o {{VENDOR_FOLDER}}/unpoly.min.css https://cdn.jsdelivr.net/npm/unpoly@3.9.2/unpoly.min.css
