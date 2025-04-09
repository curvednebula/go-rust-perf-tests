package main

import (
	"runtime"
)

type cpuWorkers[T any] struct {
	sem      chan struct{}
	requests chan func() T
	results  chan T
	recycle  bool
	Num      int
}

// numWorkers = 0 -> num workers as CPU threads
func NewCpuWorkers[T any](numWorkers int, recycle bool) *cpuWorkers[T] {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	var w *cpuWorkers[T]

	if recycle {
		w = &cpuWorkers[T]{
			requests: make(chan func() T, numWorkers),
		}

		for range numWorkers {
			go func() {
				for {
					fn := <-w.requests
					w.results <- fn()
				}
			}()
		}
	} else {
		w = &cpuWorkers[T]{
			sem: make(chan struct{}, numWorkers),
		}
	}

	w.results = make(chan T, numWorkers)
	w.recycle = recycle
	w.Num = numWorkers

	return w
}

// will block until free CPU thread is available to execute workFn()
func (w *cpuWorkers[T]) DoWork(fn func() T) {
	if w.recycle {
		// send request to already running goroutines
		w.requests <- fn
		return
	}

	w.sem <- struct{}{} // acquire slot

	go func() {
		defer func() { <-w.sem }() // release slot
		w.results <- fn()
	}()
}
