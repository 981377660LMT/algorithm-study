package main

import (
	"fmt"
	"time"
)

var seed uint64 = uint64(time.Now().UnixNano()/2 + 1)

func nextRand() uint64 {
	seed ^= seed << 7
	seed ^= seed >> 9
	return seed & (0xFFFFFFFF)
}

func main() {
	fmt.Println(nextRand())
}
