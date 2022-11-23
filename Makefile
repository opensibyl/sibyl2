# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build_all:
	$(MAKE) build_client
	$(MAKE) build_server

build_client:
	${GOCMD} build -ldflags '-extldflags "-lstdc++"' ./cmd/sibyl

build_server:
	# create swagger docs too
	cd ./pkg/server; swag init -g app.go --parseDepth 1 --parseDependency;
	${GOCMD} build -ldflags '-extldflags "-lstdc++"' ./cmd/sibyl_server

test:
	$(GOTEST) ./...
