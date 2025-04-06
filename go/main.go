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

func testNoPool(ch chan float64) {
	for range TASKS_NUM {
		go func() {
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
			// for k := range dataMap {
			// 	delete(dataMap, k)
			// }
			ch <- time.Since(start).Seconds()
		}()
	}
}

func testWithPool(ch chan float64) {
	pool := new(Pool[SomeData])

	for range TASKS_NUM {
		go func() {
			start := time.Now()
			dataMap := make(map[string]*SomeData)
			var sum uint64 = 0

			for j := uint32(0); j < VALUES_NUM; j++ {
				name := fmt.Sprintf("name-%d", j)

				var data = pool.Get()
				data.Name = name
				data.Age = j
				dataMap[name] = data

				val, exists := dataMap[name]
				if exists && val.Name == name {
					sum += uint64(val.Age)
				}
			}

			for k := range dataMap {
				pool.Release(dataMap[k])
			}
			ch <- time.Since(start).Seconds()
		}()
	}
}

func main() {
	start := time.Now()
	ch := make(chan float64)

	testNoPool(ch)
	//testWithPool(ch)

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

	fmt.Printf("%d tasks, %d items: finished in %.2fs, one task avg %.2fs, min %.2fs, max %.2fs\n", TASKS_NUM, VALUES_NUM, total, taskAvg, taskMin, taskMax)
}
