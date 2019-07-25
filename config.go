package main

import (
	"flag"
	"log"
	"strings"
)

type Config struct {
	GitRoots []string
	GitGUI   string
}

func ReadConfig() (config Config) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gitRoots := flag.String("roots", "CDPATH",
		"The name of the environment variable containing colon-separated path values to scan.")
	flag.StringVar(&config.GitGUI, "gui", "smerge",
		"The external git GUI application to use for reviews.")
	flag.Parse()

	config.GitRoots = strings.Split(*gitRoots, ":")
	return config
}

