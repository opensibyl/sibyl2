# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build:
	${GOCMD} build ./cmd/sibyl

test:
	$(GOTEST) ./...
