// FastRandom

package main

import (
	"fmt"
	"time"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	r := NewRandom()
	r.Shuffle(nums)
	fmt.Println(r.RandRange(1, 10, 2))
}

type Random struct {
	seed     int
	hashBase int
}

func NewRandom() *Random                 { return &Random{seed: int(time.Now().UnixNano()/2 + 1)} }
func NewRandomWithSeed(seed int) *Random { return &Random{seed: seed} }

func (r *Random) Rng() int {
	r.seed ^= r.seed << 7
	r.seed ^= r.seed >> 9
	return r.seed
}

func (r *Random) Next() int { return r.Rng() }

func (r *Random) RngWithMod(mod int) int { return r.Rng() % mod }

// [left, right]
func (r *Random) RandInt(min, max int) int { return min + r.Rng()%(max-min+1) }

// [start:stop:step]
func (r *Random) RandRange(start, stop int, step int) int {
	return start + r.Rng()%(stop-start)/step*step
}

// FastShuffle
func (r *Random) Shuffle(nums []int) {
	for i := range nums {
		rand := r.RandInt(0, i)
		nums[i], nums[rand] = nums[rand], nums[i]
	}
}

func (r *Random) Sample(nums []int, k int) []int {
	nums = append(nums[:0:0], nums...)
	r.Shuffle(nums)
	return nums[:k]
}

// 元组哈希
func (r *Random) HashPair(a, b int) int {
	if r.hashBase == 0 {
		r.hashBase = r.Rng()
	}
	return a*r.hashBase + b
}
