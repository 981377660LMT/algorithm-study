package main

import "time"

func main() {
	defer func() {
		if r := recover(); r != nil {
			println("Recovered in main", r)
		}
	}()

	go panicFunc()

	time.Sleep(2 * time.Second)
}

func panicFunc() {
	panic("I'm panicking!")
}
