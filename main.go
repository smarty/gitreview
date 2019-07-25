// gitreview scans entries found in the `CDPATH` environment variable
// looking for git repositories that are messy or behind and opens
// a git GUI (Sublime Merge by default) for each to facilitate a review.
// It only runs `git status` and `git fetch`, which should be safe.
// After all reviews are complete it prints (to `stdout`) a concatenated
// report of all `git fetch` output for repos that were behind their origin.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)

	var gitGUI string
	flag.StringVar(&gitGUI, "gui", "smerge", "The external git GUI application to use for reviews.")
	flag.Parse()

	messes := make(map[string]string)
	reviews := make(map[string]string)

	for _, root := range strings.Split(os.Getenv("CDPATH"), ":") {
		if root == "." {
			continue
		}
		listing, err := ioutil.ReadDir(root)
		if err != nil {
			log.Panicln(err)
		}
		for i, item := range listing {
			path := filepath.Join(root, item.Name())
			if !item.IsDir() {
				continue
			}
			git := filepath.Join(path, ".git")
			_, err := os.Stat(git)
			if os.IsNotExist(err) {
				continue
			}

			messy := exec.Command("git", "status", "--porcelain", "-uall")
			messy.Dir = path
			out, err := messy.CombinedOutput()
			if err != nil {
				log.Printf("[ERROR] Could not ascertain repo status for %s: %v", path, err)
				continue
			}

			if len(strings.TrimSpace(string(out))) > 0 {
				log.Printf("[MESSY] %d/%d: %s", i+1, len(listing), path)
				messes[path] = string(out)
			}

			log.Printf("Fetching %d/%d: %s", i+1, len(listing), path)
			fetch := exec.Command("git", "fetch")
			fetch.Dir = path
			out, err = fetch.CombinedOutput()
			if err != nil {
				log.Printf("[ERROR] Could not fetch %s: %v", path, err)
				continue
			}

			if strings.Contains(string(out), "..") {
				reviews[path] = string(out)
			}
		}
	}

	if len(reviews) + len(messes) == 0 {
		log.Println("Nothing to review today.")
		return
	}

	log.Printf("There are %d messy repositories and %d updated repositories to review.", len(messes), len(reviews))

	for path, mess := range messes {
		log.Println(path)
		log.Println(mess)
		log.Printf("Press <ENTER> to open git GUI...")
		err := exec.Command(gitGUI, path).Run()
		if err != nil {
			log.Println(err)
		}
		bufio.NewScanner(os.Stdin).Scan()
	}

	for path, fetch := range reviews {
		if _, found := messes[path]; found {
			continue // already reviewed
		}
		log.Println(fetch)
		log.Printf("Press <ENTER> to open git GUI...")
		err := exec.Command(gitGUI, path).Run()
		if err != nil {
			log.Println(err)
		}
		bufio.NewScanner(os.Stdin).Scan()
	}

	log.Println("---------------------------------------------")
	log.Println("Copy the following into your code review log:")
	log.Println("---------------------------------------------")
	log.Println()

	fmt.Printf("## %s\n\n", time.Now().Format("2006-01-02"))
	for _, fetch := range reviews {
		fmt.Println(fetch)
	}
}
