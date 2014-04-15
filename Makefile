default: build

.PHONY: build

package = gosearch

test:
	GOPATH=`pwd` go test $(package)/...

build:test
	GOPATH=`pwd` go build $(package)/index
	GOPATH=`pwd` go build $(package)/search

