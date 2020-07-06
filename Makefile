#!/usr/bin/make -f

test: fmt
	go test ./...

fmt:
	go mod tidy
	go fmt ./...

docs:
	-go run *.go -help 2>&1 >/dev/null | grep -v 'exit status 2' > README.md

install:
	go install

package:
	go build -trimpath -o gitreview
	zip gitreview.zip gitreview README.md LICENSE.md
	# TODO: use 'hub' to upload artifacts w/ release
	rm gitreview


.PHONY: test fmt docs install package
