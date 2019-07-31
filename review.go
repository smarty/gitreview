package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type GitReviewer struct {
	config    *Config
	repoPaths []string

	erred   map[string]string
	messy   map[string]string
	ahead   map[string]string
	behind  map[string]string
	fetched map[string]string
	journal map[string]string
}

func NewGitReviewer(config *Config) *GitReviewer {
	return &GitReviewer{
		config: config,
		repoPaths: append(
			collectGitRepositories(config.GitRepositoryRoots),
			filterGitRepositories(config.GitRepositoryPaths)...
		),
		erred:   make(map[string]string),
		messy:   make(map[string]string),
		ahead:   make(map[string]string),
		behind:  make(map[string]string),
		fetched: make(map[string]string),
		journal: make(map[string]string),
	}
}

func (this *GitReviewer) GitAnalyzeAll() {
	log.Printf("Analyzing %d git repositories...", len(this.repoPaths))
	reports := NewAnalyzer(workerCount).AnalyzeAll(this.repoPaths)
	for _, report := range reports {
		if len(report.StatusError) > 0 {
			this.erred[report.RepoPath] += report.StatusError
			log.Println("[ERROR]", report.StatusError)
		}
		if len(report.FetchError) > 0 {
			this.erred[report.RepoPath] += report.FetchError
			log.Println("[ERROR]", report.FetchError)
		}
		if len(report.RevListError) > 0 {
			this.erred[report.RepoPath] += report.RevListError
			log.Println("[ERROR]", report.RevListError)
		}

		if len(report.StatusOutput) > 0 {
			this.messy[report.RepoPath] += report.StatusOutput
		}
		if len(report.RevListAhead) > 0 {
			this.ahead[report.RepoPath] += report.RevListAhead
		}
		if len(report.RevListBehind) > 0 {
			this.behind[report.RepoPath] += report.RevListBehind
		}

		if len(report.FetchOutput) > 0 {
			this.fetched[report.RepoPath] += report.FetchOutput + report.RevListOutput
			this.journal[report.RepoPath] += report.FetchOutput + report.RevListOutput
		}
	}
}

func (this *GitReviewer) ReviewAll() {
	for path := range this.journal {
		if !strings.Contains(strings.ToLower(path), "smartystreets") {
			delete(this.journal, path) // Don't include external code in review log.
		}
	}

	reviewable := sortUniqueKeys(this.erred, this.messy, this.ahead, this.behind, this.fetched, this.journal)
	if len(reviewable) == 0 {
		log.Println("Nothing to review at this time.")
		return
	}

	printMapKeys(this.erred, "Repositories with git errors: %d")
	printMapKeys(this.messy, "Repositories with uncommitted changes: %d")
	printMapKeys(this.ahead, "Repositories ahead of origin master: %d")
	printMapKeys(this.behind, "Repositories behind origin master: %d")
	printMapKeys(this.fetched, "Repositories with new content since the last review: %d")
	printStrings(reviewable, "Repositories to be reviewed: %d")

	prompt(fmt.Sprintf("Press <ENTER> to initiate the review process (will open %d review windows)...", len(reviewable)))

	for _, path := range reviewable {
		log.Printf("Opening %s at %s", this.config.GitGUILauncher, path)
		err := exec.Command(this.config.GitGUILauncher, path).Run()
		if err != nil {
			log.Println("Failed to open git GUI:", err)
		}
	}
}

func (this *GitReviewer) PrintCodeReviewLogEntry() {
	printMapKeys(this.journal, "Repositories to be included in the final report: %d")

	if len(this.journal) == 0 {
		return
	}

	writer := this.config.OpenOutputWriter()
	defer func() { _ = writer.Close() }()

	prompt("Press <ENTER> to conclude review process and print code review log entry...")

	fmt.Fprintln(writer)
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "##", time.Now().Format("2006-01-02"))
	fmt.Fprintln(writer)
	for _, review := range this.journal {
		fmt.Fprintln(writer, review)
	}
}

const workerCount = 16
