package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type GitReviewer struct {
	gitRoots []string
	gitGUI   string
	problems map[string]string
	messes   map[string]string
	reviews  map[string]string
}

func NewGitReviewer(gitRoots []string, gitGUI string) *GitReviewer {
	return &GitReviewer{
		gitRoots: gitRoots,
		gitGUI:   gitGUI,
		problems: make(map[string]string),
		messes:   make(map[string]string),
		reviews:  make(map[string]string),
	}
}

func (this *GitReviewer) FetchAllRepositories() {
	for _, root := range this.gitRoots {
		if root == "." {
			continue
		}
		listing, err := ioutil.ReadDir(root)
		if err != nil {
			log.Panicln(err)
		}
		for i, item := range listing {
			path := filepath.Join(root, item.Name())
			if !item.IsDir() {
				continue
			}
			git := filepath.Join(path, ".git")
			_, err := os.Stat(git)
			if os.IsNotExist(err) {
				continue
			}

			out, err := execute(path, gitStatusCommand)
			if err != nil {
				this.problems[path] = fmt.Sprintln("[ERROR] Could not ascertain repo status:", err)
				continue
			}

			if len(strings.TrimSpace(string(out))) > 0 {
				this.messes[path] = string(out)
			}

			log.Printf("Fetching %d/%d: %s", i+1, len(listing), path)
			out, err = execute(path, gitFetchCommand)
			if err != nil {
				this.problems[path] = fmt.Sprintln("[ERROR] Could not fetch:", err)
				continue
			}

			if strings.Contains(string(out), pendingReviewIndicator) {
				this.reviews[path] = string(out)
			}
		}
	}
}

func (this *GitReviewer) notableRepositoryCount() int {
	return len(this.problems) + len(this.messes) + len(this.reviews)
}

func (this *GitReviewer) ReviewAllNotableRepositories() {
	if this.notableRepositoryCount() == 0 {
		log.Println("Nothing to review today.")
		return
	}

	printMap(this.problems, "The following %d repositories experienced errors:")
	printMap(this.messes, "The following %d repositories have uncommitted changes:")
	printMap(this.reviews, "The following %d repositories have been updated:")

	keys := sortUniqueKeys(this.problems, this.messes, this.reviews)
	log.Printf("Now beginning review of %d total repositories...", len(keys))

	for _, path := range keys {
		if containsKey(this.problems, path) {
			log.Println(path, this.problems[path])
		}
		if containsKey(this.messes, path) {
			log.Println(path, this.messes[path])
		}
		if containsKey(this.reviews, path) {
			log.Println(this.reviews[path])
		}
		log.Printf("Press <ENTER> to open git GUI...")
		bufio.NewScanner(os.Stdin).Scan()
		err := exec.Command(this.gitGUI, path).Run()
		if err != nil {
			log.Println("Failed to open git GUI:", err)
		}
	}
}

func (this *GitReviewer) PrintCodeReviewLogEntry() {
	if this.notableRepositoryCount() == 0 {
		return
	}

	log.Println("--------------------------------------------")
	log.Println("Copy the following into the code review log:")
	log.Println("--------------------------------------------")
	log.Println()

	fmt.Printf("## %s\n\n", time.Now().Format("2006-01-02"))
	for _, fetch := range this.reviews {
		if !strings.Contains(strings.ToLower(fetch), "smartystreets") {
			continue // Don't include external code in review log.
		}
		fmt.Println(fetch)
	}
}

func sortUniqueKeys(maps ...map[string]string) (unique []string) {
	combined := make(map[string]struct{})
	for _, m := range maps {
		for key := range m {
			combined[key] = struct{}{}
		}
	}
	for key := range combined {
		unique = append(unique, key)
	}
	sort.Strings(unique)
	return unique
}

func containsKey(m map[string]string, key string) bool {
	_, found := m[key]
	return found
}

func printMap(m map[string]string, preamble string) {
	if len(m) == 0 {
		return
	}
	log.Printf(preamble, len(m))
	log.Println()
	for path := range m {
		log.Println(path)
	}
	log.Println()
}

func execute(dir, command string) (string, error) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

const (
	gitStatusCommand       = "git status --porcelain -uall"
	gitFetchCommand        = "git fetch"
	pendingReviewIndicator = ".." // ie. 7761a97..1bbecb6  master     -> origin/master
)
