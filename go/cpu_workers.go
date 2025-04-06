package main

import (
	"runtime"
)

type cpuWorkers[T any] struct {
	sem        chan struct{}
	NumWorkers int
}

func NewCpuWorkers[T any]() *cpuWorkers[T] {
	w := &cpuWorkers[T]{
		sem:        make(chan struct{}, runtime.NumCPU()),
		NumWorkers: runtime.NumCPU(),
	}
	return w
}

func (w *cpuWorkers[T]) DoWork(workFn func() T) T {
	w.sem <- struct{}{} // acquire slot
	ch := make(chan T, 1)

	go func() {
		defer func() { <-w.sem }() // release slot
		ch <- workFn()
	}()
	return <-ch
}
