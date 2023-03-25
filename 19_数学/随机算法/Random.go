// FastRandom

package main

import (
	"fmt"
	"time"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	r := NewRandom()
	r.ShuffleInt(nums)
	fmt.Println(nums)
}

type Random struct {
	seed     uint64
	hashBase uint64
}

func NewRandom() *Random {
	return &Random{seed: uint64(time.Now().UnixNano()/2 + 1)}
}

func (r *Random) Rng() uint64 {
	r.seed ^= r.seed << 7
	r.seed ^= r.seed >> 9
	return r.seed
}

func (r *Random) Next() uint64 { return r.Rng() }

func (r *Random) RngWithMod(mod uint64) uint64 { return r.Rng() % mod }

// [left, right]
func (r *Random) RandInt(left, right uint64) uint64 { return left + r.Rng()%(right-left+1) }

// FastShuffle
func (r *Random) ShuffleInt(nums []int) {
	for i := range nums {
		rand := r.RandInt(0, uint64(i))
		nums[i], nums[rand] = nums[rand], nums[i]
	}
}

// 元组哈希
func (r *Random) HashPair(a, b uint64) uint64 {
	if r.hashBase == 0 {
		r.hashBase = r.Rng()
	}
	return a*r.hashBase + b
}
