package main

import (
	"runtime"
)

type cpuWorkersPool[T any] struct {
	requests chan func() T
	results  chan T
	Num      int
}

// numWorkers = 0 -> num workers as CPU threads
func NewCpuWorkersPool[T any](numWorkers int) *cpuWorkersPool[T] {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}
	w := &cpuWorkersPool[T]{
		requests: make(chan func() T, numWorkers),
		results:  make(chan T, numWorkers),
		Num:      numWorkers,
	}

	for range numWorkers {
		go func() {
			for {
				fn := <-w.requests
				w.results <- fn()
			}
		}()
	}
	return w
}

// will block until free CPU thread is available to execute workFn()
func (w *cpuWorkersPool[T]) DoWork(fn func() T) {
	w.requests <- fn
}
