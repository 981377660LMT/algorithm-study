package main

import (
	"fmt"
	"math/rand"
)

func main() {
	R := RandomHash(0, 1e18)
	fmt.Println(R(1))
	fmt.Println(R(1))
	fmt.Println(R(2))
}

type Value = int

// 生成[min,max]范围内的随机数,并保证同一个value对应的随机数是固定的.
func RandomHash(min, max uint64) func(value Value) uint64 {
	pool := make(map[Value]uint64)
	f := func(value Value) uint64 {
		if hash, ok := pool[value]; ok {
			return hash
		}
		rand := rand.Uint64()%(max-min+1) + min
		pool[value] = rand
		return rand
	}
	return f
}
