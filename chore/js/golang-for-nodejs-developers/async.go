// https://github.com/miguelmota/golang-for-nodejs-developers/blob/26bfd24d0caef71beeabbfb3061f70604f4a1646/examples/async_await.go#L1
// !Promise是容量为1的channel(hh)

package main

import (
	"errors"
	"fmt"
	"time"
)

func Async(name string) chan interface{} {
	promise := make(chan interface{}, 1) // !阻塞, 直到产生结果
	go func() {                          // new Promise
		time.Sleep(1 * time.Second) // async task
		if name == "fail" {
			promise <- errors.New("failed") // reject
		} else {
			promise <- "hello " + name // resolve
		}
	}()

	return promise
}

func main() {
	awaited := <-Async("bob") // await
	switch v := awaited.(type) {
	case string:
		fmt.Println(v)
	case error:
		fmt.Println(v) // log.Error(v)
	}

	awaited = <-Async("fail")
	switch v := awaited.(type) {
	case string:
		fmt.Println(v)
	case error:
		fmt.Println(v)
	}
}
