package main

import "sync"

type Analyzer struct {
	workerCount int
	workerInput chan string
}

func NewAnalyzer(workerCount int) *Analyzer {
	return &Analyzer{
		workerCount: workerCount,
		workerInput: make(chan string),
	}
}

func (this *Analyzer) AnalyzeAll(paths []string) (fetches []*GitReport) {
	go this.loadInputs(paths)
	outputs := this.startWorkers()
	for fetch := range merge(outputs...) {
		fetches = append(fetches, fetch)
	}
	return fetches
}

func (this *Analyzer) loadInputs(paths []string) {
	for _, path := range paths {
		this.workerInput <- path
	}
	close(this.workerInput)
}

func (this *Analyzer) startWorkers() (outputs []chan *GitReport) {
	for x := 0; x < this.workerCount; x++ {
		output := make(chan *GitReport)
		outputs = append(outputs, output)
		go NewWorker(x, this.workerInput, output).Start()
	}
	return outputs
}

func merge(fannedOut ...chan *GitReport) chan *GitReport {
	var waiter sync.WaitGroup
	waiter.Add(len(fannedOut))

	fannedIn := make(chan *GitReport)

	output := func(c <-chan *GitReport) {
		for n := range c {
			fannedIn <- n
		}
		waiter.Done()
	}

	for _, c := range fannedOut {
		go output(c)
	}

	go func() {
		waiter.Wait()
		close(fannedIn)
	}()

	return fannedIn
}
