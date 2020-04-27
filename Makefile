GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

all: fmt build start

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o build/vault/plugins/drago cmd/drago/main.go

start:
	./dev/vault server -dev -dev-root-token-id=root -dev-plugin-dir=build/vault/plugins

enable:
	VAULT_ADDR='http://127.0.0.1:8200' ./dev/vault secrets enable -description="drago secrets plugin" drago

clean:
	rm -f ./build/vault/plugins/drago

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start enable