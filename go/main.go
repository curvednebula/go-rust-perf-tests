package main

import (
	"fmt"
	"math"
	"time"
)

const TASKS_NUM = 100_000
const VALUES_NUM = 10_000

type SomeData struct {
	Name string
	Age  uint32
}

func testNoPool() (total float64, min float64, max float64, avg float64) {
	start := time.Now()
	ch := make(chan float64)

	for range TASKS_NUM {
		go func() {
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
			ch <- time.Since(start).Seconds()
		}()
	}

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
	duration := time.Since(start)

	return duration.Seconds(), taskMin, taskMax, taskSum / TASKS_NUM
}

// func testWithPool() {
// 	var wg sync.WaitGroup
//  pool := Pool{}

// 	for i := range TASKS_NUM {
// 		wg.Add(1)

// 		go func(i int) {
// 			defer wg.Done()
// 			dataMap := make(map[string]*SomeData)

// 			for j := uint32(0); j < VALUES_NUM; j++ {
// 				name := fmt.Sprintf("name-%d", j)

// 				var data = pool.Get()
// 				data.Name = name
// 				data.Age = j
// 				dataMap[name] = data

// 				_, exists := dataMap[name]
// 				if exists {
// 					//
// 				}
// 			}

// 			for k := range dataMap {
// 				pool.Release(dataMap[k])
// 			}
// 		}(i)
// 	}

// 	wg.Wait()
// }

func main() {
	total, min, max, avg := testNoPool()
	//testWithPool()
	fmt.Printf("%d tasks, %d iterrations in each: finished in %.2fs, one task avg %.2fs, min %.2fs, max %.2fs\n", TASKS_NUM, VALUES_NUM, total, avg, min, max)
}
