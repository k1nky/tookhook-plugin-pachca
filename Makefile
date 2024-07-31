SHELL:=/bin/bash
STATICCHECK=$(shell which staticcheck)

.DEFAULT_GOAL := build

test:
	go test -cover ./...

vet:
	go vet ./...
	$(STATICCHECK) ./...

generate:
	go generate ./...

gvt: generate vet test

cover:
	go test -cover ./... -coverprofile cover.out
	go tool cover -html cover.out -o cover.html

build: gvt 
	CGO_ENABLED=0 go build -o build/pachca cmd/*.go

plugin-dev:
	go build -o dev/pachca cmd/*.go

prepare:
	go mod tidy
	go install go.uber.org/mock/mockgen@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@latest