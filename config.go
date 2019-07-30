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
	GitRoots []string
	GitGUI   string
}

func ReadConfig() (config Config) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	flag.CommandLine.SetOutput(io.MultiWriter(os.Stdout, os.Stderr))
	flag.Usage = func() {
		fmt.Println(doc)
		fmt.Println()
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}
	gitRoots := flag.String("roots", "CDPATH",
		"The name of the environment variable containing colon-separated path values to scan.")
	flag.StringVar(&config.GitGUI, "gui", "smerge",
		"The external git GUI application to use for reviews.")
	flag.Parse()

	config.GitRoots = strings.Split(os.Getenv(*gitRoots), ":")
	return config
}

var doc = strings.Join([]string{
	"gitreview scans path entries found in the an environment",
	"variable looking for git repositories that have uncommitted",
	"changes or are behind their remote and opens a git GUI for",
	"each to facilitate a review.",
	"",
	"On each repository it runs `git status` and `git fetch`,",
	"both of which should be safe. After all reviews are complete",
	"it prints (to stdout) a concatenated report of all `git fetch`",
	"output for repositories that were behind their origin.",
	"",
	"Installation:",
	"",
	"`go get -u github.com/mdwhatcott/gitreview`",
}, "\n")
