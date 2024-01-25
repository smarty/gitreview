package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	GitFetch           bool
	GitRepositoryPaths []string
	GitRepositoryRoots []string
	GitGUILauncher     string
	OutputFilePath     string
	ReviewAhead        bool
	ReviewBehind       bool
	ReviewError        bool
	ReviewFetched      bool
	ReviewJournal      bool
	ReviewMessy        bool
}

func ReadConfig(version string) *Config {
	log.SetFlags(log.Ltime | log.Lshortfile)

	config := new(Config)

	flags := flag.NewFlagSet(fmt.Sprintf("%s @ %s", filepath.Base(os.Args[0]), Version), flag.ExitOnError)

	flags.Usage = func() {
		_, _ = fmt.Fprintf(flags.Output(), "Usage of %s:\n\n", flags.Name())
		_, _ = fmt.Fprintf(flags.Output(), "%s\n\n```\n", doc)
		flags.PrintDefaults()
		_, _ = fmt.Fprintln(flags.Output(), "```")
	}

	flags.StringVar(&config.GitGUILauncher,
		"gui", "smerge", ""+
			"The external git GUI application to use for visual reviews.\n"+
			"-->",
	)

	flags.StringVar(&config.OutputFilePath,
		"outfile", "SMARTY_REVIEW_LOG", ""+
			"The path or name of the environment variable containing the\n"+
			"path to your pre-existing code review file. If the file exists\n"+
			"the final log entry will be appended to that file instead of stdout.\n"+
			"-->",
	)

	flags.BoolVar(&config.GitFetch,
		"fetch", true, ""+
			"When false, suppress all git fetch operations via --dry-run.\n"+
			"Repositories with updates will still be included in the review.\n"+
			"-->",
	)

	gitRoots := flags.String(
		"roots", "CDPATH", ""+
			"The name of the environment variable containing colon-separated\n"+
			"path values to scan for any git repositories contained therein.\n"+
			"Scanning is NOT recursive.\n"+
			"NOTE: this flag will be ignored in the case that non-flag command\n"+
			"line arguments representing paths to git repositories are provided.\n"+
			"-->",
	)

	repoList := flags.String(
		"roots-file", "", ""+
			"A colon-separated list of file paths, where each file contains a\n"+
			"list of repositories to examine, with one repository on a line.\n"+
			"-->",
	)

	review := flags.String(
		"review", "abejm", ""+
			"Letter code of repository statuses to review; where (a) is ahead,\n"+
			"origin/master (b) is behind origin/master, (e) has git errors,\n"+
			"(f) has new fetched contents, and (m) is messy with uncommitted\n"+
			"changes. (j) is like (f) except only 'smarty' repositories\n"+
			"are considered\n"+
			"-->",
	)

	_ = flags.Parse(os.Args[1:])

	config.ReviewAhead = strings.ContainsAny(*review, "aA")
	config.ReviewBehind = strings.ContainsAny(*review, "bB")
	config.ReviewError = strings.ContainsAny(*review, "eE")
	config.ReviewFetched = strings.ContainsAny(*review, "fF")
	config.ReviewJournal = strings.ContainsAny(*review, "jJ")
	config.ReviewMessy = strings.ContainsAny(*review, "mM")

	config.GitRepositoryPaths = flags.Args()
	roots := strings.Split(os.Getenv(*gitRoots), ":")

	if len(*repoList) > 0 {
		list := strings.Split(*repoList, ";")
		for _, l := range list {
			config.handleRepoFile(config.tryPaths(l, roots), roots)
		}
	}

	if len(config.GitRepositoryPaths) == 0 {
		config.GitRepositoryRoots = roots
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
		log.Println("Final report will be written to stdout.")
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
		file, err2 := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, stat.Mode())
		if err2 == nil {
			log.Println("Final report will be appended to", path)
			return file
		} else {
			log.Printf("Could not open file for appending: [%s] Error: %v", this.OutputFilePath, err2)
		}
	}

	log.Println("Final report will be written to stdout.")
	return os.Stdout
}

func (this *Config) handleRepoFile(path string, prefixes []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Path for roots-file cannot be opened: %s: %s", path, err)
	}

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 1 && line[0] != '#' {
			line = this.tryPaths(line, prefixes)
			this.GitRepositoryPaths = append(this.GitRepositoryPaths, line)
			i++
		}
	}

	_ = file.Close()
	if err = scanner.Err(); err != nil {
		log.Printf("Error reading roots-file: %s: %s", path, err)
	}

	log.Printf("Added %d repositories from file: %s", i, path)
}

func (this *Config) tryPaths(path string, prefixes []string) string {
	path = strings.TrimSpace(path)
	cnt := len(path)
	if cnt == 0 {
		return ""
	}

	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}

	if !filepath.IsAbs(path) {
		for _, p := range prefixes {
			test := filepath.Join(p, path)
			if _, err := os.Stat(test); err == nil {
				return test
			}
		}
	}

	return path
}

const rawDoc = `

#### SMARTY DISCLAIMER:

Subject to the terms of the associated license agreement, this
software is freely available for your use. This software is
FREE, AS IN PUPPIES, and is a gift. Enjoy your new
responsibility. This means that while we may consider
enhancement requests, we may or may not choose to entertain
requests at our sole and absolute discretion.

# gitreview

gitreview facilitates visual inspection (code review) of git
repositories that meet any of the following criteria:

1. New content was fetched
2. Behind origin/<default-branch>
3. Ahead of origin/<default-branch>
4. Messy (have uncommitted state)
5. Throw errors for the required git operations (listed below)

We use variants of the following commands to ascertain the
status of each repository:

- ''git remote''           (shows remote address)
- ''git status''           (shows uncommitted files)
- ''git fetch''            (finds new commits/tags/branches)
- ''git rev-list''         (lists commits behind/ahead-of <default-branch>)
- ''git config --get ...'' (show config parameters of a repo)

...all of which should be safe enough. 

Each repository that meets any criteria above will be
presented for review. After all reviews are complete a
concatenated report of all output from ''git fetch'' for
repositories that were behind their origin is printed to
stdout. Only repositories with "smarty" in their
path are included in this report.

Repositories are identified for consideration from path values
supplied as non-flag command line arguments or via the roots
flag (see details below).

Installation:

    go get -u github.com/smarty/gitreview


Skipping Repositories:

If you have repositories in your list that you would rather not review,
you can mark them to be skipped by adding a config variable to the
repository. The following command will produce this result:

    git config --add review.skip true


Omitting Repositories:

If you have repositories in your list that you would still like to audit
but aren't responsible to sign off (it's code from another team), you can 
mark them to be omitted from the final report by adding a config variable
to the repository. The following command will produce this result:

    git config --add review.omit true


Specifying the ''default'' branch:

This tool assumes that the default branch of all repositories is ''master''.
If a repository uses a non-standard default branch (ie. ''main'', ''trunk'')
and you want this tool to focus  reviews on commits pushed to that branch
instead, run the following command:

	git config --add review.branch <branch-name>


CLI Flags:
`

var doc = strings.ReplaceAll(strings.TrimSpace(rawDoc), "''", "`")
