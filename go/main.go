package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const TASKS_NUM = 100_000
const VALUES_NUM = 10_000

type SomeData struct {
	Name string
	Age  uint32
}

var cpu *cpuWorkers[float64] = NewCpuWorkers[float64]()

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

	for j := uint32(0); j < VALUES_NUM; j++ {
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

	for j := uint32(0); j < VALUES_NUM; j++ {
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

func testOnlyGoroutines(ch chan float64) {
	for range TASKS_NUM {
		go func() {
			ch <- doWork()
		}()
	}
}

func testWithCpuWorkers(ch chan float64) {
	for range TASKS_NUM {
		go func() {
			ch <- cpu.DoWork(doWork)
		}()
	}
}

func testWithCpuWorkersOnly(ch chan float64) {
	for range TASKS_NUM {
		ch <- cpu.DoWork(doWork)
	}
}

func testWithPoolAndCpuWorkers(ch chan float64) {
	for range TASKS_NUM {
		go func() {
			ch <- cpu.DoWork(doWorkWithPool)
		}()
	}
}

func runTest(name string, testFn func(ch chan float64)) {
	fmt.Println(name)

	start := time.Now()
	ch := make(chan float64, 128)

	// don't interrupt main thread when running the test as it need to start receving from channel asap
	go func() {
		testFn(ch)
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
	// runTest("With pure goroutines.", testOnlyGoroutines)

	test2Name := fmt.Sprintf("With CPU workers: %d workers.", cpu.NumWorkers)
	runTest(test2Name, testWithCpuWorkers)

	test3Name := fmt.Sprintf("With CPU workers only: %d workers.", cpu.NumWorkers)
	runTest(test3Name, testWithCpuWorkersOnly)

	test4Name := fmt.Sprintf("With CPU workers and pool: %d workers.", cpu.NumWorkers)
	runTest(test4Name, testWithPoolAndCpuWorkers)
}
