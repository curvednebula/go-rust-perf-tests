package main

import (
	"runtime"
)

type cpuWorkers[T any] struct {
	sem chan struct{}
	Num int
}

// numWorkers = 0 -> num workers as CPU threads
func NewCpuWorkers[T any](numWorkers int) *cpuWorkers[T] {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}
	w := &cpuWorkers[T]{
		sem: make(chan struct{}, numWorkers),
		Num: numWorkers,
	}
	return w
}

// will execute workFn() when free CPU thread is available
func (w *cpuWorkers[T]) DoWork(resultCh chan<- T, workFn func() T) {
	w.sem <- struct{}{} // acquire slot

	go func() {
		defer func() { <-w.sem }() // release slot
		resultCh <- workFn()
	}()
}
