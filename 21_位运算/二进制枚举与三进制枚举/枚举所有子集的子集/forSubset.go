// forSubset枚举某个状态的所有子集(枚举子集的子集)

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	state := 0b1101
	for g1 := state; g1 >= 0; {
		if g1 == state || g1 == 0 { // 排除空集和全集
			g1--
			continue
		}
		g2 := state ^ g1
		fmt.Println(g1, g2)
		if g1 == 0 {
			g1 = -1
		} else {
			g1 = (g1 - 1) & state
		}
	}

	EnumerateSubsetOfState1(0b1101, func(y int) bool {
		fmt.Println(y)
		return false
	})

	EnumerateSubsetOfState2(0b1101, func(y int) bool {
		fmt.Println(y)
		return false
	})

	EnumerateSupersetOfState(4, 0b1101, func(y int) bool {
		fmt.Println(y)
		return false
	})

	EnumerateSubsetOfSizeK(4, 2, func(x uint) bool {
		fmt.Println(x, "see")
		return false
	})
}

// 升序枚举state所有子集的子集.
//
//	0b1101 -> 0,1,4,5,8,9,12,13.
func EnumerateSubsetOfState1(state int, f func(subset int) (shouldBreak bool)) {
	for y := 0; ; y = (y - state) & state {
		if f(y) {
			break
		}
		if y == state {
			break
		}
	}
}

// 降序枚举state所有子集的子集.
//
//	0b1101 -> 13,12,9,8,5,4,1,0.
func EnumerateSubsetOfState2(state int, f func(subset int) (shouldBreak bool)) {
	for y := state; ; y = (y - 1) & state {
		if f(y) {
			break
		}
		if y == 0 {
			break
		}
	}
}

// 升序枚举state的所有超集.
//
//	0b1101 -> 13,15.
func EnumerateSupersetOfState(n int, state int, f func(superset int) (shouldBreak bool)) {
	upper := 1 << n
	for y := state; y < upper; y = (y + 1) | state {
		if f(y) {
			break
		}
	}
}

// 遍历n个元素的集合中大小为k的子集(combinations)
//
//	一共有C(n,k)个子集.
//	C(4,2) -> 3,5,6,9,10,12.
func EnumerateSubsetOfSizeK(n int, k int, f func(subset uint) (shouldBreak bool)) {
	if k <= 0 || k > n {
		return
	}
	upper := uint(1 << n)
	for x := uint((1 << k) - 1); x < upper; {
		if f(x) {
			break
		}
		t := x | (x - 1)
		// nextCombination (gosper hack)
		x = (t + 1) | (((^t & -^t) - 1) >> (bits.TrailingZeros(x) + 1))
	}
}
