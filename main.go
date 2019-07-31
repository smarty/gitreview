package main

func main() {
	config := ReadConfig()
	reviewer := NewGitReviewer(config)
	reviewer.GitAnalyzeAll()
	reviewer.ReviewAll()
	reviewer.PrintCodeReviewLogEntry()
}
