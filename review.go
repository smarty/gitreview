package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type GitReviewer struct {
	gitGUI    string
	repoPaths []string
	problems  map[string]string
	messes    map[string]string
	reviews   map[string]string
}

func NewGitReviewer(gitRoots []string, gitGUI string) *GitReviewer {
	return &GitReviewer{
		repoPaths: collectGitRepositoryPaths(gitRoots),
		gitGUI:    gitGUI,
		problems:  make(map[string]string),
		messes:    make(map[string]string),
		reviews:   make(map[string]string),
	}
}

func collectGitRepositoryPaths(gitRoots []string) (paths []string) {
	for _, root := range gitRoots {
		if root == "." {
			continue
		}
		if strings.TrimSpace(root) == "" {
			continue
		}
		listing, err := ioutil.ReadDir(root)
		if err != nil {
			log.Println("Couldn't resolve path (skipping):", err)
			continue
		}
		for _, item := range listing {
			path := filepath.Join(root, item.Name())
			if !item.IsDir() {
				continue
			}
			git := filepath.Join(path, ".git")
			_, err := os.Stat(git)
			if os.IsNotExist(err) {
				continue
			}

			paths = append(paths, path)
		}
	}

	return paths
}

func (this *GitReviewer) GitFetchAll() {
	log.Printf("Running `git status` and `get fetch` for %d repos...", len(this.repoPaths))
	for _, fetch := range NewGitClient(16).ScanAll(this.repoPaths) {
		if len(fetch.StatusError) > 0 {
			this.problems[fetch.RepoPath] += fetch.StatusError
		}
		if len(fetch.StatusOutput) > 0 {
			this.messes[fetch.RepoPath] = fetch.StatusOutput
		}
		if len(fetch.FetchError) > 0 {
			this.problems[fetch.RepoPath] += fetch.FetchError
		}
		if len(fetch.FetchOutput) > 0 {
			this.reviews[fetch.RepoPath] = fetch.FetchOutput
		}
	}
}

func (this *GitReviewer) formatFetchProgress(index int) string {
	progress := strings.TrimSpace(fmt.Sprintf("%3d / %-3d", index+1, len(this.repoPaths)))
	progress = "(" + progress + ")"
	for len(progress) < len("(999 / 999)") {
		progress = " " + progress
	}
	return progress
}

func (this *GitReviewer) reviewIsPending() bool {
	return len(this.problems)+len(this.messes)+len(this.reviews) > 0
}

func (this *GitReviewer) ReviewAll() {
	if !this.reviewIsPending() {
		log.Println("Nothing to review today.")
		return
	}

	printMap(this.problems, "The following %d repositories experienced errors:")
	printMap(this.messes, "The following %d repositories have uncommitted changes:")
	printMap(this.reviews, "The following %d repositories have been updated:")

	keys := sortUniqueKeys(this.problems, this.messes, this.reviews)
	log.Printf("A total of %d repositories need to be reviewed.", len(keys))
	prompt(fmt.Sprintf("Press <ENTER> to initiate review (will open %d review windows)...", len(keys)))

	for _, path := range keys {
		err := exec.Command(this.gitGUI, path).Run()
		if err != nil {
			log.Println("Failed to open git GUI:", err)
		}
	}
}

func (this *GitReviewer) PrintCodeReviewLogEntry() {
	if !this.reviewIsPending() {
		return
	}

	prompt("Press <ENTER> to conclude review process and print code review log entry...")

	fmt.Println()
	fmt.Println()
	fmt.Printf("## %s\n\n", time.Now().Format("2006-01-02"))
	for _, fetch := range this.reviews {
		if !strings.Contains(strings.ToLower(fetch), "smartystreets") {
			continue // Don't include external code in review log.
		}
		fmt.Println(fetch)
	}
}
