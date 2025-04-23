package main

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"time"
)

const TASKS_NUM = 100_000
const ITEMS_NUM = 10_000
const TASKS_IN_BUNCH = 10

// const TIME_BETWEEN_BUNCHES_MS = 1

type SomeData struct {
	Name string
	Age  uint32
}

func doWork() float64 {
	start := time.Now()
	dataMap := make(map[string]SomeData, ITEMS_NUM)
	var sum uint64 = 0

	for j := uint32(0); j < ITEMS_NUM; j++ {
		name := strconv.Itoa(int(j))

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

func Test(numWorkers int, recycle bool) {
	start := time.Now()
	workers := NewCpuWorkers[float64](numWorkers, recycle)

	fmt.Printf("%d CPU workers (recycle: %t)...\n", workers.Num, recycle)

	// don't block main thread when running the test as it needs to start receving from channel asap
	go func() {
		for range TASKS_NUM {
			workers.DoWork(doWork)
			// if taskIdx%TASKS_IN_BUNCH == 0 {
			// 	// simulate requests coming sequentially not all at once
			// 	time.Sleep(TIME_BETWEEN_BUNCHES_MS * time.Millisecond)
			// }
		}
	}()

	taskSum := float64(0)
	taskMin := math.MaxFloat64
	taskMax := -math.MaxFloat64

	for range TASKS_NUM {
		taskTime := <-workers.results
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
	Test(runtime.NumCPU()*2, true)
	Test(runtime.NumCPU()*2, false)
	Test(runtime.NumCPU()*5, true)
	Test(runtime.NumCPU()*5, false)
}
