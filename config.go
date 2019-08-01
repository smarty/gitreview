package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Config struct {
	GitFetch           bool
	GitRepositoryPaths []string
	GitRepositoryRoots []string
	GitGUILauncher     string
	OutputFilePath     string
}

func ReadConfig() *Config {
	log.SetFlags(log.Ltime | log.Lshortfile)

	config := new(Config)

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

	flag.StringVar(&config.OutputFilePath,
		"outfile", "SMARTY_REVIEW_LOG", ""+
			"The path or name of the environment variable containing the"+"\n"+
			"path to your pre-existing code review file. If the file exists"+"\n"+
			"the final log entry will be appended to that file instead of stdout."+"\n"+
			"-->",
	)

	flag.BoolVar(&config.GitFetch,
		"fetch", true, ""+
			"When false, suppress all git fetch operations via --dry-run."+"\n"+
			"Repositories with updates will still be included in the review."+"\n"+
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

	config.GitRepositoryPaths = flag.Args()
	if len(config.GitRepositoryPaths) == 0 {
		config.GitRepositoryRoots = strings.Split(os.Getenv(*gitRoots), ":")
	}
	if !config.GitFetch {
		log.Println("Running git fetch with --dry-run (updated repositories will not be reviewed).")
		gitFetchCommand += " --dry-run"
	}
	return config
}

func (this *Config) OpenOutputWriter() io.WriteCloser {
	this.OutputFilePath = strings.TrimSpace(this.OutputFilePath)
	if this.OutputFilePath == "" {
		return os.Stdout
	}

	path, found := os.LookupEnv(this.OutputFilePath)
	if found {
		log.Printf("Found output path in environment variable: %s=%s", this.OutputFilePath, path)
	} else {
		path = this.OutputFilePath
	}

	stat, err := os.Stat(path)
	if err == nil && err != os.ErrNotExist {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, stat.Mode())
		if err == nil {
			log.Println("Final report will be appended to", path)
			return file
		} else {
			log.Printf("Could not open file for appending: [%s] Error: %v", this.OutputFilePath, err)
		}
	}

	log.Println("Final report will appear in stdout.")
	return os.Stdout
}

const rawDoc = `# gitreview

gitreview facilitates visual inspection (code review) of git
repositories that meet any of the following criteria:

1. New content was fetched
2. Behind origin/master
3. Ahead of origin/master
4. Messy (have uncommitted state)
5. Throw errors for the required git operations (listed below)

We use variants of the following commands to ascertain the
status of each repository:

- 'git remote'    (shows remote address)
- 'git status'    (shows uncommitted files)
- 'git fetch'     (finds new commits/tags/branches)
- 'git rev-list'  (lists commits behind/ahead of master)

...all of which should be safe enough. 

Each repository that meets any criteria above will be
presented for review. After all reviews are complete a
concatenated report of all output from 'git fetch' for
repositories that were behind their origin is printed to
stdout. Only repositories with "smartystreets" in their
path are included in this report.

Repositories are identified for consideration from path values
supplied as non-flag command line arguments or via the roots
flag (see details below).

Installation:

    go get -u github.com/mdwhatcott/gitreview

CLI Flags:
`

var doc = strings.ReplaceAll(rawDoc, "'", "`")
