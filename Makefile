# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build_all:
	# create swagger docs too
	cd ./pkg/server; swag init -g app.go --parseDepth 1 --parseDependency;
	${GOCMD} build -ldflags '-extldflags "-lstdc++"' ./cmd/sibyl

test:
	$(GOTEST) ./...
