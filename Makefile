# Common Makefile

# This file is included in Makefiles for subdirectories, so should only contain
# directives shared throughout the project.

KINNEY=github.com/CamusEnergy/kinney

# Base go inside of module, to avoid collisions with other projects.
ifndef GOPATH
	GOMOD=$(shell go env GOMOD)
  GOPATH=$(dir ${GOMOD})go
  export GOPATH
endif

ifndef GOBIN
  GOBIN=$(GOPATH)/bin
	export GOBIN
endif

env:
	go get -v -u \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	pip install pipenv
	pipenv install
	@if [ -z "${PIPENV_ACTIVE}" ]; \
	then \
		echo ====================================; \
		echo WARNING: Currently outside of pipenv; \
		echo Run to enter: \`pipenv shell\`; \
		echo ====================================; \
	fi

.PHONY: env
