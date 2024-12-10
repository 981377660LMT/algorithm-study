// Tip #1 一行代码测量函数的执行时间

package main

import (
	"fmt"
	"time"
)

func main() {
	defer TrackTime(time.Now()) // <- 在函数开头加入这一行，即可测量此函数的执行时间

	time.Sleep(500 * time.Millisecond)
}

func TrackTime(pre time.Time) time.Duration {
	elapsed := time.Since(pre)
	fmt.Println("Time elapsed: ", elapsed)
	return elapsed
}
