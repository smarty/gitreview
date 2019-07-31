package main

import "io"

func main() {
	config := ReadConfig()
	reviewer := NewGitReviewer(
		config.GitRepositoryRoots,
		config.GitRepositoryPaths,
		config.GitGUILauncher,
	)
	reviewer.GitAnalyzeAll()
	reviewer.ReviewAll()
	reviewer.PrintCodeReviewLogEntry(config.OpenOutputWriter())
}

func close_(closer io.Closer) { _ = closer.Close() }
