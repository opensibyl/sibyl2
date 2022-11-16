# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build:
	${GOCMD} build -ldflags '-extldflags "-lstdc++"' ./cmd/sibyl
	${GOCMD} build -ldflags '-extldflags "-lstdc++"' ./cmd/sibyl_server

test:
	$(GOTEST) ./...
