# gitreview

gitreview facilitates visual inspection (code review) of git
repositories that meet any of the following criteria:

1. Have uncommitted changes
2. Are behind the 'origin' remote
3. Are ahead of their 'origin' remote

For each considered repository we run variants of:

1. `git status`
2. `git fetch`
3. `git rev-list`

...all of which should be safe enough. After all reviews are
complete a concatenated report of all `git fetch` output for
repositories that were behind their origin is printed to stdout.
Only repositories with "smartystreets" in their path are
included in this final report.

Repositories are identified for consideration from path values
supplied as non-flag command line arguments or via the roots
flag (see details below).

Installation:

    go get -u github.com/mdwhatcott/gitreview

CLI Flags:


```
  -gui string
    	The external git GUI application to use for visual reviews.
    	--> (default "smerge")
  -outfile string
    	The path or name of the environment variable containing the
    	path to your pre-existing code review file. If the file exists
    	the final log entry will be appended to that file instead of stdout.
    	--> (default "SMARTY_REVIEW_LOG")
  -roots string
    	The name of the environment variable containing colon-separated
    	path values to scan for any git repositories contained therein.
    	Scanning is NOT recursive.
    	NOTE: this flag will be ignored in the case that non-flag command
    	line arguments representing paths to git repositories are provided.
    	--> (default "CDPATH")
```
