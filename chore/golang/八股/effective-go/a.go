package main

import "fmt"

func main() {
	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d ", i)
	}
}

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)
