package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	GitRepositoryPaths []string
	GitRepositoryRoots []string
	GitGUILauncher     string
	OutputFilePath     string
}

func ReadConfig() (config Config) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, doc)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "```")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "```")
	}

	flag.StringVar(&config.GitGUILauncher,
		"gui", "smerge", ""+
			"The external git GUI application to use for visual reviews."+"\n"+
			"-->",
	)

	outfile := flag.String(
		"outfile", "SMARTY_REVIEW_LOG", ""+
			"The path or name of the environment variable containing the"+"\n"+
			"path to your pre-existing code review file. If the file exists"+"\n"+
			"the final log entry will be appended to that file instead of stdout."+"\n"+
			"-->",
	)

	gitRoots := flag.String(
		"roots", "CDPATH", ""+
			"The name of the environment variable containing colon-separated"+"\n"+
			"path values to scan for any git repositories contained therein."+"\n"+
			"Scanning is NOT recursive."+"\n"+
			"NOTE: this flag will be ignored in the case that non-flag command"+"\n"+
			"line arguments representing paths to git repositories are provided."+"\n"+
			"-->",
	)

	flag.Parse()

	config.OutputFilePath = os.Getenv(*outfile)
	config.GitRepositoryRoots = strings.Split(os.Getenv(*gitRoots), ":")
	config.GitRepositoryPaths = flag.Args()
	return config
}

const rawDoc = `# gitreview

gitreview facilitates visual inspection (code review) of git
repositories that meet any of the following criteria:

2. Are behind the 'origin' remote,
3. Are ahead of their 'origin' remote,
1. Have uncommitted changes (including untracked files).

To ascertain the status of a repository we run variants of:

1. 'git status'
2. 'git fetch'
3. 'git rev-list'

...all of which should be safe enough. Each repository
that meets the criteria above will be presented for review.
After all reviews are complete a concatenated report of all
output from 'git fetch' for repositories that were behind 
their origin is printed to stdout. Only repositories with 
"smartystreets" in their path are included in this report.

Repositories are identified for consideration from path values
supplied as non-flag command line arguments or via the roots
flag (see details below).

Installation:

    go get -u github.com/mdwhatcott/gitreview

CLI Flags:
`

var doc = strings.ReplaceAll(rawDoc, "'", "`")
