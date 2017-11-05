# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = oceano2oceansites
VET_REPORT = vet.report
TEST_REPORT = tests.x
GOARCH = amd64

VERSION = 0.2.5
COMMIT = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=jgrelet
CURRENT_DIR=$(shell pwd)

# demo
DRIVE = data/CTD/test
CRUISE = TEST
CONFIG = oceano2oceansites.toml
ROSCOP = roscop/code_roscop.csv
PREFIX = csp

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.Version=${VERSION}  \
-X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'`"

# Build the project
all: clean build test vet demo

build: 
	go get .
	go build ${LDFLAGS}  .

buildall: 
	go get .
	go build ${LDFLAGS} -a -v .

install:
	go install 

test:
	go test -v ./...  

demo:
	${BINARY} -v
	${BINARY} -c ${CONFIG} -r ${ROSCOP} -e --files=${DRIVE}/${PREFIX}*.cnv 
	ncdump -v PROFILE,LATITUDE,LONGITUDE,BATH netcdf/OS_${CRUISE}_CTD.nc

fmt:
	go fmt $$(go list ./... | grep -v /vendor/) 

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BINARY}-*

.PHONY: build install test vet fmt clean