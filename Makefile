#!/usr/bin/make -f

VERSION := $(shell git describe)

test: fmt
	GORACE="atexit_sleep_ms=50" go test ./...

fmt:
	go mod tidy && go fmt ./...

docs:
	-go run *.go -help 2>&1 >/dev/null | grep -v 'exit status 2' > README.md

install:
	go install -ldflags="-X 'main.Version=$(VERSION)'"

package:
	go build -ldflags="-X 'main.Version=$(VERSION)'" -trimpath -o gitreview
	zip gitreview.zip gitreview README.md LICENSE.md
	# TODO: use 'hub' to upload artifacts w/ release
	rm gitreview


.PHONY: test fmt docs install package
