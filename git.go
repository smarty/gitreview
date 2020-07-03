package main

import (
	"fmt"
	"strings"
)

var (
	gitRemoteCommand      = "git remote -v"                                    // ie. [origin	git@github.com:mdwhatcott/gitreview.git (fetch)]
	gitStatusCommand      = "git status --porcelain -uall"                     // parse-able output, including untracked
	gitFetchCommand       = "git fetch"                                        // --dry-run"  // for debugging
	gitFetchPendingReview = ".."                                               // ie. [7761a97..1bbecb6  master     -> origin/master]
	gitRevListCommand     = "git rev-list --left-right master...origin/master" // 1 line per commit w/ prefix '<' (ahead) or '>' (behind)
	gitErrorTemplate      = "[ERROR] Could not execute [%s]: %v" + "\n"
	gitOmitCommand        = "git config --get review.omit"
	gitSkipCommand        = "git config --get review.skip"
)

type GitReport struct {
	RepoPath string

	RemoteError  string
	StatusError  string
	FetchError   string
	RevListError string
	SkipError    string

	RemoteOutput  string
	StatusOutput  string
	FetchOutput   string
	RevListOutput string
	OmitOutput    string
	SkipOutput    string

	RevListAhead  string
	RevListBehind string
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
func (this *GitReport) GitSkipStatus() bool {
	out, err := execute(this.RepoPath, gitSkipCommand)
	if err != nil && err.Error() != "exit status 1" {
		this.SkipError = fmt.Sprintf(gitErrorTemplate, gitSkipCommand, err)
	}
	if strings.Contains(out, "true") {
		this.SkipOutput = out
		return true
	}
	return false
}
func (this *GitReport) GitOmitStatus() bool {
	out, err := execute(this.RepoPath, gitOmitCommand)
	if err != nil && err.Error() != "exit status 1" {
		this.SkipError = fmt.Sprintf(gitErrorTemplate, gitOmitCommand, err)
	}
	if strings.Contains(out, "true") {
		this.OmitOutput = out
		return true
	}
	return false
}
func (this *GitReport) GitFetch() {
	out, err := execute(this.RepoPath, gitFetchCommand)
	if err != nil {
		this.FetchError = fmt.Sprintf(gitErrorTemplate, gitFetchCommand, err)
	}
	if output := string(out); strings.Contains(output, gitFetchPendingReview) {
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
		this.RevListAhead = fmt.Sprintf("The default branch is %d commits ahead of origin.\n", ahead)
	}
	if behind > 0 {
		this.RevListBehind = fmt.Sprintf("The default branch is %d commits behind origin.\n", behind)
	}
}

func (this *GitReport) Progress() string {
	status := ""
	if len(this.StatusError+this.FetchError+this.RemoteError+this.RevListError) > 0 {
		status += "!"
	} else {
		status += " "
	}
	if len(this.StatusOutput) > 0 {
		status += "M"
	} else {
		status += " "
	}
	if len(this.RevListAhead) > 0 {
		status += "A"
	} else {
		status += " "
	}
	if len(this.RevListBehind) > 0 {
		status += "B"
	} else {
		status += " "
	}
	if len(this.FetchOutput) > 0 {
		status += "F"
	} else {
		status += " "
	}
	if len(this.OmitOutput) > 0 {
		status += "O"
	} else {
		status += " "
	}
	if len(this.SkipOutput) > 0 {
		status += "S"
	} else {
		status += " "
	}
	return fmt.Sprintf("[%-7s] %s", status, this.RepoPath)
}
