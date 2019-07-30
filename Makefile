#!/usr/bin/make -f

docs:
	-go run *.go -help 2> README.txt
	echo "# gitreview" > README.md
	echo >> README.md
	cat README.txt | grep -v 'exit status 2' >> README.md
	rm README.txt

install:
	go install

package:
	go build -o gitreview
	zip gitreview.zip gitreview README.md LICENSE.md
	rm gitreview
