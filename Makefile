# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/
# use go module since version 1.15
# go mod init github.com/jgrelet/oceano2oceansites
# go test

BINARY = oceano2oceansites
GOARCH = amd64

VERSION = 0.2.6
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

# set env
export OCEANO2OCEANSITES_CFG = ${CONFIG}
export ROSCOP_CSV = ${ROSCOP}

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -buildmode=exe -ldflags  "-X main.Version=${VERSION}  \
-X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'`"

# Build the project
all: dep build install test demo

dep:
	@go get -v -t ./...
	@go vet  ./...

build: 
	@go build ${LDFLAGS}  .

buildall: 
	@go build ${LDFLAGS} -a -v .

install:
	@go install 

mod:
	@go mod tidy  
	$(MAKE) test 

test:
	@go test -v ./...  

fmt:
	@go fmt $$(go list ./... | grep -v /vendor/) 

demo:
	${BINARY} -e -c ${CONFIG} -r ${ROSCOP} --files=${DRIVE}/${PREFIX}*.cnv 
	${BINARY} -v
	${BINARY} --files=${DRIVE}/${PREFIX}*.cnv 

ncdump:
	ncdump -v PROFILE,LATITUDE,LONGITUDE,BATH netcdf/OS_${CRUISE}_CTD.nc

clean:
	-rm -f ${BINARY}-*

.PHONY: build buildall install test demo ncdump fmt  mod clean
