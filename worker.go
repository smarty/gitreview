package main

import (
	"log"
	"path/filepath"
)

type Worker struct {
	id  int
	in  chan string
	out chan *GitReport
}

func NewWorker(id int, in chan string, out chan *GitReport) *Worker {
	return &Worker{id: id, in: in, out: out}
}

func (this *Worker) Start() {
	for path := range this.in {
		this.out <- this.git(path)
	}
	close(this.out)
}

func (this *Worker) git(path string) *GitReport {
	path, _ = filepath.Abs(path)
	report := &GitReport{RepoPath: path}
	if !report.GitSkipStatus() {
		report.GitOmitStatus()
		report.GitRemote()
		report.GitStatus()
		report.GitFetch()
		report.GitRevList()
	}
	log.Println(report.Progress())
	return report
}
