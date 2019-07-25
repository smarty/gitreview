# gitreview
--
gitreview scans entries found in the CDPATH environment variable looking for git
repositories that are messy or behind and opens a git GUI (sublime merge by
default) for each to facilitate a review. It only runs `git status` and `git
fetch`, which should be safe. After all reviews are complete it prints (to
stdout) a concatenated report of all 'git fetch' output for repos that were
behind their origin.
