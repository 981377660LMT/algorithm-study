// 光速幂
// 已知base和模数mod，求base^n 模mod
// !O(sqrt(maxN))预处理,O(1)查询

package main

import (
	"math"
)

func main() {

	fp := NewFastPow(2, 1e5)
	println(fp.Pow(1e18))
}

// https://leetcode.cn/problems/count-collisions-of-monkeys-on-a-polygon/

var fp2 = NewFastPow(2, 1e9)

func monkeyMove(n int) int {
	res := fp2.Pow(n-1) - 2
	if res < 0 {
		res += MOD
	}
	return res
}

const MOD int = 1e9 + 7

type E = int

func (*FastPow) e() E        { return 1 }
func (*FastPow) op(a, b E) E { return a * b % MOD }

// 光速幂.
type FastPow struct {
	max     int
	divData []E
	modData []E
}

// O(sqrt(maxN))预处理,O(1)查询.
//
//	base: 幂运算的基.
//	maxN: 最大的幂.
func NewFastPow(base E, maxN int) *FastPow {
	max := int(math.Ceil(math.Sqrt(float64(maxN))))
	res := &FastPow{max: max, divData: make([]E, max+1), modData: make([]E, max+1)}
	cur := res.e()
	for i := 0; i <= max; i++ {
		res.modData[i] = cur
		cur = res.op(cur, base)
	}
	cur = res.e()
	last := res.modData[max]
	for i := 0; i <= max; i++ {
		res.divData[i] = cur
		cur = res.op(cur, last)
	}
	return res
}

// n<=maxN.
func (fp *FastPow) Pow(n int) E {
	return fp.op(fp.divData[n/fp.max], fp.modData[n%fp.max])
}
