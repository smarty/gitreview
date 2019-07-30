package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

const (
	gitStatusCommand       = "git status --porcelain -uall"                             // parse-able output, including untracked
	gitFetchCommand        = "git fetch"                                                // --dry-run" // for debugging
	gitRevListCommand      = "git rev-list --left-right --count master...origin/master" // [4	6]    // (master is 4 ahead, 6 behind)
	pendingReviewIndicator = ".."                                                       // ie. [7761a97..1bbecb6  master     -> origin/master]
	gitErrorTemplate       = "[ERROR] Could not execute [%s]: %v" + "\n"
)

type GitClient struct {
	workerCount int
	workerInput chan string
}

func NewGitClient(workerCount int) *GitClient {
	return &GitClient{
		workerCount: workerCount,
		workerInput: make(chan string),
	}
}

func (this *GitClient) ScanAll(paths []string) (fetches []*GitReport) {
	go this.loadInputs(paths)
	outputs := this.startWorkers()
	for fetch := range merge(outputs...) {
		fetches = append(fetches, fetch)
	}
	return fetches
}

func (this *GitClient) startWorkers() (outputs []chan *GitReport) {
	for x := 0; x < this.workerCount; x++ {
		output := make(chan *GitReport)
		outputs = append(outputs, output)
		go NewGitWorker(x, this.workerInput, output).Start()
	}
	return outputs
}

func (this *GitClient) loadInputs(paths []string) {
	for _, path := range paths {
		this.workerInput <- path
	}
	close(this.workerInput)
}

type GitWorker struct {
	id  int
	in  chan string
	out chan *GitReport
}

func NewGitWorker(id int, in chan string, out chan *GitReport) *GitWorker {
	return &GitWorker{
		id:  id,
		in:  in,
		out: out,
	}
}

func (this *GitWorker) Start() {
	for path := range this.in {
		this.out <- this.git(path)
	}
	close(this.out)
}

func (this *GitWorker) git(path string) *GitReport {
	log.Println(path)
	report := &GitReport{RepoPath: path}
	report.GitStatus()
	report.GitFetch()
	report.GitRevList()
	return report
}

type GitReport struct {
	RepoPath     string
	StatusOutput string
	StatusError  string
	FetchOutput  string
	FetchError   string
	RevListError string
	RevListAhead string
}

func (this *GitReport) GitStatus() {
	out, err := execute(this.RepoPath, gitStatusCommand)
	if err != nil {
		this.StatusError = fmt.Sprintf(gitErrorTemplate, gitStatusCommand, err)
		return
	}
	if output := string(out); len(strings.TrimSpace(output)) > 0 {
		this.StatusOutput = output
	}
}
func (this *GitReport) GitFetch() {
	out, err := execute(this.RepoPath, gitFetchCommand)
	if err != nil {
		this.FetchError = fmt.Sprintf(gitErrorTemplate, gitFetchCommand, err)
	}
	if output := string(out); strings.Contains(output, pendingReviewIndicator) {
		this.FetchOutput = output
	}
}
func (this *GitReport) GitRevList() {
	out, err := execute(this.RepoPath, gitRevListCommand)
	if err != nil {
		this.RevListError = fmt.Sprintf(gitErrorTemplate, gitRevListCommand, err)
	}
	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) < 2 {
		return
	}
	ahead, _ := strconv.Atoi(fields[0])
	if ahead > 0 {
		this.RevListAhead = fmt.Sprintf("The master branch is %d commits ahead of origin/master.\n", ahead)
	}
}

func merge(fannedOut ...chan *GitReport) chan *GitReport {
	var waiter sync.WaitGroup
	waiter.Add(len(fannedOut))

	fannedIn := make(chan *GitReport)

	output := func(c <-chan *GitReport) {
		for n := range c {
			fannedIn <- n
		}
		waiter.Done()
	}

	for _, c := range fannedOut {
		go output(c)
	}

	go func() {
		waiter.Wait()
		close(fannedIn)
	}()

	return fannedIn
}
