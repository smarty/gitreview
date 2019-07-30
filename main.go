package main

func main() {
	config := ReadConfig()
	reviewer := NewGitReviewer(config.GitRepositoryRoots, config.GitGUILauncher)
	reviewer.GitFetchAll()
	reviewer.ReviewAll()
	reviewer.PrintCodeReviewLogEntry()
}
