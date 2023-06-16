package main

import (
	"fmt"
	"time"
)

func main() {
	cb := func() {
		fmt.Println("hello")
	}
	clear := SetInterval(cb, 500*time.Millisecond)

	time.Sleep(1 * time.Second)
	clear()
	time.Sleep(1 * time.Second)
	fmt.Println("done")
}

func SetInterval(callback func(), interval time.Duration) (ClearInterval func()) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			callback()
		}
	}()

	return func() {
		ticker.Stop()
	}
}
