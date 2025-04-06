package main

import (
	"runtime"
)

type cpuWorkers[T any] struct {
	sem        chan struct{}
	NumWorkers int
}

func NewCpuWorkers[T any]() *cpuWorkers[T] {
	workersNum := runtime.NumCPU()

	w := &cpuWorkers[T]{
		sem:        make(chan struct{}, workersNum),
		NumWorkers: workersNum,
	}
	return w
}

// will block until free CPU thread is available to execute the workFn()
func (w *cpuWorkers[T]) DoWork(workFn func() T) T {
	w.sem <- struct{}{} // acquire slot
	ch := make(chan T, 1)

	go func() {
		defer func() { <-w.sem }() // release slot
		ch <- workFn()
	}()
	return <-ch
}
