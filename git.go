package main

import (
	"fmt"
	"strings"
)

var (
	gitRemoteCommand         = "git remote -v"                            // ie. [origin	git@github.com:mdwhatcott/gitreview.git (fetch)]
	gitStatusCommand         = "git status --porcelain -uall"             // parse-able output, including untracked
	gitFetchCommand          = "git fetch"                                // --dry-run"  // for debugging
	gitFetchPendingReview    = ".."                                       // ie. [7761a97..1bbecb6  master     -> origin/master]
	gitRevListCommand        = "git rev-list --left-right %s...origin/%s" // 1 line per commit w/ prefix '<' (ahead) or '>' (behind)
	gitErrorTemplate         = "[ERROR] Could not execute [%s]: %v" + "\n"
	gitOmitCommand           = "git config --get review.omit"
	gitSkipCommand           = "git config --get review.skip"
	gitDefaultBranchCommand  = "git config --get review.branch"
	gitStandardDefaultBranch = "master"
)

func GitRevListCommand(branch string) string {
	return fmt.Sprintf(gitRevListCommand, branch, branch)
}

type GitReport struct {
	RepoPath string

	RemoteError  string
	StatusError  string
	FetchError   string
	RevListError string

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
	fields := strings.Fields(out)
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
	if len(strings.TrimSpace(out)) > 0 {
		this.StatusOutput = out
	}
}
func (this *GitReport) GitSkipStatus() bool {
	out, _ := execute(this.RepoPath, gitSkipCommand)
	this.SkipOutput = out
	return strings.Contains(out, "true")
}
func (this *GitReport) GitOmitStatus() bool {
	out, _ := execute(this.RepoPath, gitOmitCommand)
	this.OmitOutput = out
	return strings.Contains(out, "true")
}
func (this *GitReport) GitDefaultBranch() string {
	out, _ := execute(this.RepoPath, gitDefaultBranchCommand)
	branch := strings.TrimSpace(out)
	if branch == "" {
		return gitStandardDefaultBranch
	}
	return branch
}
func (this *GitReport) GitFetch() {
	out, err := execute(this.RepoPath, gitFetchCommand)
	if err != nil {
		this.FetchError = fmt.Sprintf(gitErrorTemplate, gitFetchCommand, err)
	}
	if strings.Contains(out, gitFetchPendingReview) {
		this.FetchOutput = out
	}
}
func (this *GitReport) GitRevList() {
	branch := this.GitDefaultBranch()
	command := GitRevListCommand(branch)
	out, err := execute(this.RepoPath, command)
	if err != nil {
		this.RevListError = fmt.Sprintf(gitErrorTemplate, command, err)
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
		this.RevListAhead = fmt.Sprintf("The %s branch is %d commits ahead of origin/%s.\n", branch, ahead, branch)
	}
	if behind > 0 {
		this.RevListBehind = fmt.Sprintf("The %s branch is %d commits behind origin/%s.\n", branch, behind, branch)
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
