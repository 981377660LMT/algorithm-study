package main

import (
	"fmt"
	"time"
)

func main() {
	arr := make([]int, int(1e7))

	time1 := time.Now()
	for i := 0; i < len(arr); i++ {
		_ = arr[i]
	}
	fmt.Println("init:", time.Now().Sub(time1))

	time1 = time.Now()
	for i := len(arr) - 1; i >= 0; i-- {
		_ = arr[i]
	}
	fmt.Println("reverse:", time.Now().Sub(time1))

}
