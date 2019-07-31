package main

func main() {
	config := ReadConfig()
	reviewer := NewGitReviewer(
		config.GitRepositoryRoots,
		config.GitRepositoryPaths,
		config.GitGUILauncher,
	)
	reviewer.GitAnalyzeAll()
	reviewer.ReviewAll()
	reviewer.PrintCodeReviewLogEntry(config.OpenOutputWriter)
}
