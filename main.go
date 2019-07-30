package main

import "io"

func main() {
	config := ReadConfig()
	writer := config.OpenOutputWriter()
	defer close_(writer)

	reviewer := NewGitReviewer(
		config.GitRepositoryRoots,
		config.GitRepositoryPaths,
		config.GitGUILauncher,
	)
	reviewer.GitAnalyzeAll()
	reviewer.ReviewAll()
	reviewer.PrintCodeReviewLogEntry(writer)
}

func close_(closer io.Closer) { _ = closer.Close() }
