// https://go.dev/play/p/0bw2icHuKRa
// js/python:无栈协程
// golang:有栈协程
// 不过无栈协程也可以有多线程的调度器
// 丢到一个goroutine里去，把cpu密集的工作，异步化
// 有栈协程更接近多线程，实际也可以多线程调度执行

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const MOD int = 1e11 + 7

func main() {

	time1 := time.Now()

	w := 10
	ModSumAsync(0, 5e8, MOD, &w)

	time2 := time.Now()
	fmt.Println(time2.Sub(time1)) // 692.8902ms
}

// 计算`sum(range(start,end))%mod`.
func ModSumAsync(start int, end int, mod int, worker *int) int {
	if start >= end {
		return 0
	}

	realWorker := runtime.NumCPU()
	if worker != nil {
		realWorker = *worker
	}

	count := end - start
	base := count / realWorker
	remain := count % realWorker
	run := func(workerId int) Promise {
		res := make(Promise, 1)

		go func() {
			more := min(workerId, remain)
			normal := max(0, workerId-remain)
			curStart := start + more*(base+1) + normal*base
			curEnd := curStart + base
			if more < remain {
				curEnd++
			}
			sum := 0
			for i := curStart; i < curEnd; i++ {
				sum = (sum + i) % mod
			}
			res <- sum
			close(res)
		}()
		return res
	}

	tasks := make([]Promise, realWorker)
	for i := 0; i < realWorker; i++ {
		tasks[i] = run(i)
	}
	res := PromiseAllInt(tasks...)
	awaited := <-res
	switch v := awaited.(type) {
	case error:
		panic(v)
	case []int:
		sum := 0
		for _, v := range v {
			sum = (sum + v) % mod
		}
		return sum
	}

	panic("unreachable")
}

type Promise = chan interface{}

func PromiseAllInt(tasks ...Promise) Promise {
	var wg sync.WaitGroup
	res := make([]int, len(tasks))
	promise := make(Promise, 1)

	for taskId, task := range tasks {
		wg.Add(1)
		go func(taskId int, task Promise) {
			awaited := <-task // await
			switch v := awaited.(type) {
			case error:
				promise <- v // reject
			case int:
				res[taskId] = v
			}
			wg.Done()
		}(taskId, task)
	}

	wg.Wait() // count == 0
	promise <- res
	close(promise)

	return promise
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
