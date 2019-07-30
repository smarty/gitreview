# gitreview

gitreview scans path entries found in the an environment
variable looking for git repositories that have uncommitted
changes or are behind their remote and opens a git GUI for
each to facilitate a review.

On each repository it runs `git status` and `git fetch`,
both of which should be safe. After all reviews are complete
it prints (to stdout) a concatenated report of all `git fetch`
output for repositories that were behind their origin.

Installation:

`go get -u github.com/mdwhatcott/gitreview`

CLI Flags:

```
  -gui string
    	The external git GUI application to use for reviews. (default "smerge")
  -roots string
    	The name of the environment variable containing colon-separated path values to scan. (default "CDPATH")
```
