package handlers

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	maxWorkers int
	jobs       chan func()
	wg         sync.WaitGroup
	shutdown   chan struct{}
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	pool := &WorkerPool{
		maxWorkers: maxWorkers,
		jobs:       make(chan func(), maxWorkers*2), // Buffer for 2x workers
		shutdown:   make(chan struct{}),
	}

	for i := 0; i < maxWorkers; i++ {
		go pool.worker()
	}

	return pool
}

func (p *WorkerPool) worker() {
	for {
		select {
		case job := <-p.jobs:
			// Recover from panics in job execution
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Recovered from panic in worker: %v\n", r)
					}
					p.wg.Done()
				}()
				job()
			}()
		case <-p.shutdown:
			return
		}
	}
}

func (p *WorkerPool) Submit(job func()) {
	p.wg.Add(1)
	p.jobs <- job
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

func (p *WorkerPool) Shutdown() {
	close(p.shutdown)
	close(p.jobs)
}
