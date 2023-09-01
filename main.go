package main

var Version = "dev"

func main() {
	config := ReadConfig(Version)
	reviewer := NewGitReviewer(config)
	reviewer.GitAnalyzeAll()
	reviewer.ReviewAll()
	reviewer.PrintCodeReviewLogEntry()
}
