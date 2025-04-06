package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const TASKS_NUM = 100_000
const ITEMS_NUM = 10_000
const TASKS_IN_BUNCH = 10
const TIME_BETWEEN_BUNCHES_MS = 1

type SomeData struct {
	Name string
	Age  uint32
}

var workers *cpuWorkers[float64] = NewCpuWorkers[float64](0)

// var pool *Pool[SomeData] = new(Pool[SomeData])

var pool = sync.Pool{
	New: func() any {
		return new(SomeData)
	},
}

func doWork() float64 {
	start := time.Now()
	dataMap := make(map[string]SomeData)
	var sum uint64 = 0

	for j := uint32(0); j < ITEMS_NUM; j++ {
		name := fmt.Sprintf("name-%d", j)

		dataMap[name] = SomeData{
			Name: name,
			Age:  j,
		}

		val, exists := dataMap[name]
		if exists && val.Name == name {
			sum += uint64(val.Age)
		}
	}
	return time.Since(start).Seconds()
}

func doWorkWithPool() float64 {
	start := time.Now()
	dataMap := make(map[string]*SomeData)
	var sum uint64 = 0

	for j := uint32(0); j < ITEMS_NUM; j++ {
		name := fmt.Sprintf("name-%d", j)

		var data = pool.Get().(*SomeData)
		data.Name = name
		data.Age = j
		dataMap[name] = data

		val, exists := dataMap[name]
		if exists && val.Name == name {
			sum += uint64(val.Age)
		}
	}

	for k := range dataMap {
		pool.Put(dataMap[k])
	}
	return time.Since(start).Seconds()
}

func goroutinesOnly(ch chan float64) {
	go func() {
		ch <- doWork()
	}()
}

func goroutineWithCpuWorkers(ch chan float64) {
	go func() {
		workers.DoWork(ch, doWork)
	}()
}

func cpuWorkersOnly(ch chan float64) {
	workers.DoWork(ch, doWork)
}

func cpuWorkersAndPool(ch chan float64) {
	go func() {
		workers.DoWork(ch, doWorkWithPool)
	}()
}

func runTest(name string, testFn func(ch chan float64)) {
	fmt.Println(name)

	start := time.Now()
	ch := make(chan float64, 128)

	// don't block main thread when running the test as it needs to start receving from channel asap
	go func() {
		for taskIdx := range TASKS_NUM {
			testFn(ch)
			if taskIdx%TASKS_IN_BUNCH == 0 {
				// simulate requests coming sequentially not all at once
				time.Sleep(TIME_BETWEEN_BUNCHES_MS * time.Millisecond)
			}
		}
	}()

	taskSum := float64(0)
	taskMin := math.MaxFloat64
	taskMax := -math.MaxFloat64

	for range TASKS_NUM {
		taskTime := <-ch
		taskSum += taskTime

		if taskMin > taskTime {
			taskMin = taskTime
		}
		if taskMax < taskTime {
			taskMax = taskTime
		}
	}
	total := time.Since(start).Seconds()
	taskAvg := taskSum / TASKS_NUM

	fmt.Printf(" - finished in %.4fs, task avg %.4fs, min %.4fs, max %.4fs\n", total, taskAvg, taskMin, taskMax)
}

func main() {
	// runTest("With pure goroutines.", goroutineOnly)

	// test2Name := fmt.Sprintf("With CPU workers: %d workers.", workers.Num)
	// runTest(test2Name, goroutineWithCpuWorkers)

	test3Name := fmt.Sprintf("With CPU workers only: %d workers.", workers.Num)
	runTest(test3Name, cpuWorkersOnly)

	// test4Name := fmt.Sprintf("With CPU workers and pool: %d workers.", cpu.NumWorkers)
	// runTest(test4Name, cpuWorkersAndPool)
}
