# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = oceano2oceansites
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION = 0.2.5
COMMIT = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=jgrelet
CURRENT_DIR=$(shell pwd)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.Binary=${BINARY} -X main.Version=${VERSION} -X main.COMMIT=${COMMIT} \
-X main.BRANCH=${BRANCH} \
-X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'`"

# Build the project
all: clean test vet linux darwin windows

linux: 
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} . ; \

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} . ; \

windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; \

test:
	if ! hash go2xunit 2>/dev/null; then go install github.com/tebeka/go2xunit; fi
	godep go test -v ./... 2>&1 | go2xunit -output ${TEST_REPORT} ; \

vet:
	godep go vet ./... > ${VET_REPORT} 2>&1 ; \

fmt:
	go fmt $$(go list ./... | grep -v /vendor/) ; \

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BINARY}-*

.PHONY: link linux darwin windows test vet fmt clean