BINARY := gcp-sales-api
GOPATH := $(GOPATH)
GOARCH := amd64

# DEF: VERSION/BUILDTIME
# Need to follow Semantic Versioning.
# VERSION=`git describe --abbrev=0 --tags`
VERSION := `git describe --tags --abbrev=0`
COMMIT := `git rev-parse HEAD`
BRANCH := `git rev-parse --abbrev-ref HEAD`

# Package target
BUILD_DIR := ${GOPATH}/src/service/cmd/${BINARY}

# Build Binaries setting BuildInfo vars
LDFLAGS :=-ldflags "-d -s -w -X main.version=${VERSION} \
	-X main.commit=${COMMIT} \
	-X main.branch=${BRANCH}"

# Build the project
all: clean linux darwin

linux: 
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} CGO_ENABLED=0 go build -v ${LDFLAGS} -a -tags netgo -installsuffix netgo -o ${BINARY}-linux-${GOARCH} . ; \
	cd - >/dev/null

darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} . ; \
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

clean:
	-rm -f ${BINARY}-*


.PHONY: linux darwin fmt clean
