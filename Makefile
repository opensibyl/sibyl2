# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build-all:
	${GOCMD} build ./cmd/sibyl

	cd ./extras/casedoctor && ${GOCMD} build ./cmd/casedoctor && cd ../..
	cd ./extras/storytrack && ${GOCMD} build ./cmd/storytrack && cd ../..

test:
	$(GOTEST) ./...

	cd ./extras/casedoctor && ${GOCMD} test ./... && cd ../..
	cd ./extras/storytrack && ${GOCMD} test ./... && cd ../..

