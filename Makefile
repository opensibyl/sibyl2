# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build:
	${GOCMD} build -ldflags '-extldflags "-lstdc++"' ./cmd/sibyl

test:
	$(GOTEST) ./...
