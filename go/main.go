package main

import (
	"fmt"
	"sync"
	"time"
)

const THREADS_NUM = 100_000
const VALUES_NUM = 10_000

type SomeData struct {
	Name string
	Age  int32
}

var (
	pool      []*SomeData = make([]*SomeData, 0)
	poolMutex sync.Mutex
)

func getSomeData() *SomeData {
	poolMutex.Lock()

	poolLen := len(pool)
	var data *SomeData

	if poolLen == 0 {
		data = &SomeData{}
		pool = append(pool, data)
	} else {
		data = pool[poolLen-1]
		pool = pool[:poolLen-1]
	}
	poolMutex.Unlock()
	return data
}

func releaseSomeData(d *SomeData) {
	poolMutex.Lock()
	pool = append(pool, d)
	poolMutex.Unlock()
}

func testWithStruct() {
	var wg sync.WaitGroup

	for i := range THREADS_NUM {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			dataMap := make(map[string]SomeData)

			for j := int32(0); j < VALUES_NUM; j++ {
				name := fmt.Sprintf("name-%d", j)

				dataMap[name] = SomeData{
					Name: name,
					Age:  j,
				}

				_, exists := dataMap[name]
				if exists {
					//
				}
			}
		}(i)
	}

	wg.Wait()
}

func testWithStructPtr() {
	var wg sync.WaitGroup

	for i := range THREADS_NUM {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			dataMap := make(map[string]*SomeData)

			for j := int32(0); j < VALUES_NUM; j++ {
				name := fmt.Sprintf("name-%d", j)

				dataMap[name] = &SomeData{
					Name: name,
					Age:  j,
				}

				_, exists := dataMap[name]
				if exists {
					//
				}
			}
		}(i)
	}

	wg.Wait()
}

func testWithPool() {
	var wg sync.WaitGroup

	for i := range THREADS_NUM {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			dataMap := make(map[string]*SomeData)

			for j := int32(0); j < VALUES_NUM; j++ {
				name := fmt.Sprintf("name-%d", j)

				var data = getSomeData()
				data.Name = name
				data.Age = j
				dataMap[name] = data

				_, exists := dataMap[name]
				if exists {
					//
				}
			}

			for k := range dataMap {
				releaseSomeData(dataMap[k])
			}
		}(i)
	}

	wg.Wait()
}

func main() {
	start := time.Now()
	testWithStruct()
	//testWithStructPtr()
	//testWithPool()
	fmt.Printf("%d threads finished %d iterrations each in %.2f seconds\n", THREADS_NUM, VALUES_NUM, time.Since(start).Seconds())
}
