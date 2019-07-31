package main

import (
	"fmt"
	"strings"
)

const (
	gitRemoteCommand       = "git remote -v"                                    // ie. [origin	git@github.com:mdwhatcott/gitreview.git (fetch)]
	gitStatusCommand       = "git status --porcelain -uall"                     // parse-able output, including untracked
	gitFetchCommand        = "git fetch "                                       // --dry-run"  // for debugging TODO
	gitRevListCommand      = "git rev-list --left-right master...origin/master" // 1 line per commit w/ prefix '<' (ahead) or '>' (behind)
	pendingReviewIndicator = ".."                                               // ie. [7761a97..1bbecb6  master     -> origin/master]
	gitErrorTemplate       = "[ERROR] Could not execute [%s]: %v" + "\n"
)

type GitReport struct {
	RepoPath      string
	RemoteOutput  string
	RemoteError   string
	StatusOutput  string
	StatusError   string
	FetchOutput   string
	FetchError    string
	RevListError  string
	RevListAhead  string
	RevListBehind string
	RevListOutput string
}

func (this *GitReport) GitRemote() {
	out, err := execute(this.RepoPath, gitRemoteCommand)
	if err != nil {
		this.RemoteError = fmt.Sprintf(gitErrorTemplate, gitRemoteCommand, err)
		this.RemoteOutput = this.RepoPath
		return
	}
	fields := strings.Fields(string(out))
	if len(fields) < 2 {
		return
	}
	this.RemoteOutput = fields[1]
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
	behind, ahead := 0, 0
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, ">") {
			this.RevListOutput += "  " + line + "\n"
			behind++
		} else if strings.HasPrefix(line, "<") {
			ahead++
		}
	}
	if ahead > 0 {
		this.RevListAhead = fmt.Sprintf("The master branch is %d commits ahead of origin/master.\n", ahead)
	}
	if behind > 0 {
		this.RevListBehind = fmt.Sprintf("The master branch is %d commits behind origin/master.\n", behind)
	}
}
