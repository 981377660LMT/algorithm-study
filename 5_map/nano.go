package main

import (
	"fmt"
	"time"
)

func main() {
	time1 := time.Now()
	R := NewRandom()
	upper := int(1e6)
	for i := 0; i < upper; i++ {
		R.Rng()
	}
	fmt.Println(time.Since(time1))
}

type Random struct {
	seed     uint64
	hashBase uint64
}

func NewRandom() *Random                 { return &Random{seed: uint64(time.Now().UnixNano()/2 + 1)} }
func NewRandomWithSeed(seed int) *Random { return &Random{seed: uint64(seed)} }

func (r *Random) Rng() uint64 {
	r.seed ^= r.seed << 7
	r.seed ^= r.seed >> 9
	return r.seed
}

func (r *Random) Rng61() uint64 { return r.Rng() & ((1 << 61) - 1) }
