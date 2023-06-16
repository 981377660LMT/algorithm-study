// golang实现一个Promise.all(
// 本质上和js一样

package main

import (
	"fmt"
	"sync"
	"time"
)

// Promise是一个容量为1的channel.
type Promise = chan interface{}

// Promise.all
func PromiseAll(tasks ...Promise) Promise {
	var wg sync.WaitGroup
	res := make([]string, len(tasks))
	promise := make(Promise, 1)

	go func() {
		for taskId, task := range tasks {
			wg.Add(1)
			go func(taskId int, task Promise) {
				awaited := <-task // await
				switch v := awaited.(type) {
				case error:
					promise <- v // reject
				case string:
					res[taskId] = v
				}
				wg.Done()
			}(taskId, task)
		}

		wg.Wait() // count == 0
		promise <- res
		close(promise)
	}()

	return promise
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		awaited := <-asyncTask("foo")
		switch v := awaited.(type) {
		case string:
			fmt.Println(v)
		case error:
			fmt.Println(v)
		}
		wg.Done()
	}()

	go func() {
		awaited := <-PromiseAll(
			asyncTask("A"),
			asyncTask("B"),
			asyncTask("C"),
		)
		switch v := awaited.(type) {
		case []string:
			fmt.Println(v)
		case error:
			fmt.Println(v)
		}
		wg.Done()
	}()

	wg.Wait()
}

func asyncTask(value string) Promise {
	promise := make(Promise, 1)
	go func() {
		time.Sleep(1 * time.Second)
		promise <- "resolved: " + value
		close(promise)
	}()
	return promise
}
