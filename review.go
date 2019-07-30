package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

type GitReviewer struct {
	gitGUI       string
	repoPaths    []string
	pathToRemote map[string]string
	problems     map[string]string
	messes       map[string]string
	reviews      map[string]string
}

func NewGitReviewer(gitRoots, gitPaths []string, gitGUI string) *GitReviewer {
	return &GitReviewer{
		repoPaths: append(
			collectGitRepositories(gitRoots),
			filterGitRepositories(gitPaths)...
		),
		gitGUI:       gitGUI,
		pathToRemote: make(map[string]string),
		problems:     make(map[string]string),
		messes:       make(map[string]string),
		reviews:      make(map[string]string),
	}
}

func (this *GitReviewer) GitAnalyzeAll() {
	log.Printf("Analyzing %d git repositories...", len(this.repoPaths))
	reports := NewAnalyzer(16).AnalyzeAll(this.repoPaths)
	for _, report := range reports {
		this.pathToRemote[report.RepoPath] = report.RemoteOutput

		if len(report.StatusError) > 0 {
			this.problems[report.RepoPath] += report.StatusError
			log.Println("[WARN]", report.StatusError)
		}
		if len(report.FetchError) > 0 {
			this.problems[report.RepoPath] += report.FetchError
			log.Println("[WARN]", report.FetchError)
		}
		if len(report.RevListError) > 0 {
			this.problems[report.RepoPath] += report.RevListError
			log.Println("[WARN]", report.RevListError)
		}

		if len(report.StatusOutput) > 0 {
			this.messes[report.RepoPath] += report.StatusOutput
		}
		if len(report.RevListAhead) > 0 {
			this.messes[report.RepoPath] += report.RevListAhead
		}

		if len(report.FetchOutput) > 0 {
			this.reviews[report.RepoPath] += report.FetchOutput
		}

		if len(report.FetchOutput) > 0 && len(report.RevListBehind) > 0 {
			this.reviews[report.RepoPath] += report.RevListOutput
		}
	}
}

func (this *GitReviewer) ReviewAll() {
	if len(this.problems)+len(this.messes)+len(this.reviews) == 0 {
		log.Println("Nothing to review at this time.")
		return
	}

	printMap(this.problems, "The following %d repositories experienced errors:")
	printMap(this.messes, "The following %d repositories have uncommitted changes or are ahead of origin master:")
	printMap(this.reviews, "The following %d repositories have new content or are behind origin master:")

	keys := sortUniqueKeys(this.problems, this.messes, this.reviews)
	log.Printf("A total of %d repositories need to be reviewed.", len(keys))
	prompt(fmt.Sprintf("Press <ENTER> to initiate review (will open %d review windows)...", len(keys)))

	for _, path := range keys {
		log.Printf("Opening %s at %s", this.gitGUI, path)
		err := exec.Command(this.gitGUI, path).Run()
		if err != nil {
			log.Println("Failed to open git GUI:", err)
		}
	}
}

func (this *GitReviewer) PrintCodeReviewLogEntry(output io.Writer) {
	if len(this.reviews) == 0 {
		log.Println("Nothing to report at this time.")
		return
	}

	prompt("Press <ENTER> to conclude review process and print code review log entry...")

	fmt.Fprintln(output)
	fmt.Fprintln(output)
	fmt.Fprintf(output, "## %s\n\n", time.Now().Format("2006-01-02"))
	for path, review := range this.reviews {
		if !strings.Contains(strings.ToLower(path), "smartystreets") {
			continue // Don't include external code in review log.
		}
		fmt.Fprintln(output, review)
	}
}
