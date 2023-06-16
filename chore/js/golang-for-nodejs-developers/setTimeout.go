package main

import (
	"time"
)

func main() {
	clear := SetTimeout(func() {
		println("Hello, World!")
	}, 1000*time.Millisecond)

	// clear the timeout
	clear()

}

func SetTimeout(callback func(), duration time.Duration) (clearTimeout func()) {
	canceled := false
	time.AfterFunc(duration, func() {
		if canceled {
			return
		}
		callback()
	})

	return func() {
		canceled = true
	}
}
