package main

import (
	"fmt"
	"sync"
)

func main() {
	wgExample()
}

// !数据聚合
func wgExample() {
	var wg sync.WaitGroup
	recChan := make(chan int)

	// producer
	numTask := 10
	for i := 0; i < numTask; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			recChan <- i
		}(i)
	}

	// closer
	go func() {
		wg.Wait()
		close(recChan)
	}()

	// consumer
	for v := range recChan {
		fmt.Println(v)
	}
}
