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
	counter := make(map[uint64]int)
	for i := 0; i < 1000; i++ {
		counter[r.RandRange(1, 10, 2)]++
	}
	fmt.Println(counter)
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

func (r *Random) Next() uint64 { return r.Rng() }

func (r *Random) RngWithMod(mod int) uint64 { return r.Rng() % uint64(mod) }

// [left, right]
func (r *Random) RandInt(min, max int) uint64 { return uint64(min) + r.Rng()%(uint64(max-min+1)) }

// [start:stop:step]
func (r *Random) RandRange(start, stop int, step int) uint64 {
	width := stop - start
	// Fast path.
	if step == 1 {
		return uint64(start) + r.Rng()%uint64(width)
	}
	var n uint64
	if step > 0 {
		n = uint64((width + step - 1) / step)
	} else {
		n = uint64((width + step + 1) / step)
	}
	return uint64(start) + uint64(step)*(r.Rng()%n)
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
func (r *Random) HashPair(a, b int) uint64 {
	if r.hashBase == 0 {
		r.hashBase = r.Rng()
	}
	return uint64(a)*r.hashBase + uint64(b)
}

func (r *Random) GetHashBase1D(nums []int) []uint64 {
	hashBase := make([]uint64, len(nums))
	for i := range hashBase {
		hashBase[i] = r.Rng()
	}
	return hashBase
}

func (r *Random) GetHashBase2D(nums [][]int) [][]uint64 {
	hashBase := make([][]uint64, len(nums))
	for i := range hashBase {
		hashBase[i] = make([]uint64, len(nums[i]))
		for j := range hashBase[i] {
			hashBase[i][j] = r.Rng()
		}
	}
	return hashBase
}
