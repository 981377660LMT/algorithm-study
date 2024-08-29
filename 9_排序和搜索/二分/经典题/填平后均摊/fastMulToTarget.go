package main

import (
	"fmt"
	"math/bits"
)

func main() {
	F := FastMulToTarget(2, 100)
	fmt.Println(F(1, 100)) // 7 128
}

// 给定 multiplier >= 2
// O(1) 计算把任意正整数 x 通过不断乘 multiplier，直到 >= y，需要乘多少次
// 返回最小的 e，满足 x * multiplier^e >= y
// 额外返回 powM = multiplier^e
// LC3266 https://leetcode.cn/problems/final-array-state-after-k-multiplication-operations-ii/
// https://github.com/EndlessCheng/codeforces-go/commit/8648991a126963d33f0880a620a1a38209e4bc6d#diff-4fe0d06c08ff6757f221629de960dd041a8a803b8e1870495b422856f3a0e3b8R1196
func FastMulToTarget(multiplier int, maxY int) func(x, y int) (e, powM int) {
	if multiplier < 2 {
		panic("multiplier must be at least 2")
	}
	maxY++
	type ep struct{ e, powM int }
	ePowM := make([]ep, 0, bits.Len(uint(maxY)))
	for pow2, powM, e := 1, 1, 0; pow2 <= maxY; pow2 <<= 1 {
		if powM < pow2 {
			powM *= multiplier
			e++
		}
		ePowM = append(ePowM, ep{e, powM})
	}

	fastMul := func(x, y int) (e, powM int) {
		if x >= y {
			return 0, 1
		}
		p := ePowM[bits.Len(uint(y))-bits.Len(uint(x))]
		e, powM = p.e, p.powM
		if powM/multiplier*x >= y {
			powM /= multiplier
			e--
		} else if x*powM < y {
			powM *= multiplier
			e++
		}
		return
	}

	return fastMul
}
