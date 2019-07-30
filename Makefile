#!/usr/bin/make -f

docs:
	echo "# gitreview" > README.md
	echo >> README.md
	go run *.go -help >> README.md

install:
	go install

package:
	go build -o gitreview
	zip gitreview.zip gitreview README.md LICENSE.md
