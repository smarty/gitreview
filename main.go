package main

func main() {
	config := ReadConfig()
	reviewer := NewGitReviewer(config.GitRoots, config.GitGUI)
	reviewer.GitFetchAll()
	reviewer.ReviewAll()
	reviewer.PrintCodeReviewLogEntry()
}
