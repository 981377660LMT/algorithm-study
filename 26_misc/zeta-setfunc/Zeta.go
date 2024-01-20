package main

import (
	"fmt"
	"math/bits"
)

func main() {
	Enumerate(3, func(s1, s2 int) {
		// do something
		if s1 > s2 {
			fmt.Println(s1, s2)
		}
	})

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	SubsetZeta(nums)
	fmt.Println(nums)
}

// SOS DP (Sum over Subsets)
func SubsetZeta(nums []int) {
	if len(nums) == 0 {
		return
	}
	log := bits.Len(uint(len(nums))) - 1
	if 1<<log != len(nums) {
		panic("len(nums) must be power of 2")
	}
	for n := 0; n < log; n++ {
		mask := 1 << n
		for s := 0; s < 1<<log; s++ {
			t := s ^ mask
			if s > t {
				nums[s] += nums[t] // add
			}
		}
	}
}

func SubsetMobius(nums []int) {
	if len(nums) == 0 {
		return
	}
	log := bits.Len(uint(len(nums))) - 1
	if 1<<log != len(nums) {
		panic("len(nums) must be power of 2")
	}
	for n := 0; n < log; n++ {
		mask := 1 << n
		for s := 0; s < 1<<log; s++ {
			t := s ^ mask
			if s > t {
				nums[s] -= nums[t] // inv
			}
		}
	}
}

func Enumerate(log int, f func(s1, s2 int)) {
	for n := 0; n < log; n++ {
		mask := 1 << n
		for s := 0; s < 1<<log; s++ {
			t := s ^ mask
			f(s, t)
		}
	}
}

func SuperSetZeta(nums []int) {
	if len(nums) == 0 {
		return
	}
	log := bits.Len(uint(len(nums))) - 1
	if 1<<log != len(nums) {
		panic("len(nums) must be power of 2")
	}
	for n := 0; n < log; n++ {
		mask := 1 << n
		for s := 0; s < 1<<log; s++ {
			t := s ^ mask
			if s < t {
				nums[s] += nums[t] // add
			}
		}
	}
}

func SupersetMobius(nums []int) {
	if len(nums) == 0 {
		return
	}
	log := bits.Len(uint(len(nums))) - 1
	if 1<<log != len(nums) {
		panic("len(nums) must be power of 2")
	}
	for n := 0; n < log; n++ {
		mask := 1 << n
		for s := 0; s < 1<<log; s++ {
			t := s ^ mask
			if s < t {
				nums[s] -= nums[t] // inv
			}
		}
	}
}
