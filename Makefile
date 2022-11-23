# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build_all:
	${MAKE} prepare
	${GOCMD} build -ldflags '-extldflags "-lstdc++"' ./cmd/sibyl

prepare:
	cd ./pkg/server; swag init -g app.go --parseDepth 1 --parseDependency;

test:
	$(GOTEST) ./...
