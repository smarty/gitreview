#!/usr/bin/make -f

docs:
	-go run *.go -help 2>&1 >/dev/null | grep -v 'exit status 2' > README.md

install:
	go install

package:
	go build -trimpath -o gitreview
	zip gitreview.zip gitreview README.md LICENSE.md
	rm gitreview
