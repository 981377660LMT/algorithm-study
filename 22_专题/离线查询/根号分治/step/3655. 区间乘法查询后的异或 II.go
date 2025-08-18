// 3655. 区间乘法查询后的异或 II
// https://leetcode.cn/problems/xor-after-range-multiplication-queries-ii/description/
// 给你一个长度为 n 的整数数组 nums 和一个大小为 q 的二维整数数组 queries，其中 queries[i] = [li, ri, ki, vi]。
//
// 对于每个查询，需要按以下步骤依次执行操作：
//
// 设定 idx = li。
// 当 idx <= ri 时：
// 更新：nums[idx] = (nums[idx] * vi) % (1e9 + 7)。
// 将 idx += ki。
// 在处理完所有查询后，返回数组 nums 中所有元素的 按位异或 结果。
//
// - step 大 -> 暴力
// - step 小 -> 这样的step不多，对每种step，差分求每个位置处的乘积

package main

import (
	"math"
)

const mod int = 1e9 + 7

func Pow(base, exp, mod int) int {
	base %= mod
	res := 1 % mod
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}

func xorAfterQueries(nums []int, queries [][]int) int {
	n := len(nums)
	B := int(math.Sqrt(float64(n))) + 1

	resMul := make([]int, n)
	for i := range resMul {
		resMul[i] = 1
	}

	var big [][]int
	small := make([][][]int, B+1)
	for _, q := range queries {
		if q[2] <= B {
			k := q[2]
			small[k] = append(small[k], q)
		} else {
			big = append(big, q)
		}
	}

	for _, q := range big {
		l, r, k, v := q[0], q[1], q[2], q[3]
		for idx := l; idx <= r; idx += k {
			resMul[idx] = resMul[idx] * v % mod
		}
	}

	for k := 1; k <= B; k++ {
		if len(small[k]) == 0 {
			continue
		}

		diff := make([]int, n+k+1)
		for i := range diff {
			diff[i] = 1
		}
		for _, q := range small[k] {
			l, r, v := q[0], q[1], q[3]
			inv := Pow(v, mod-2, mod)
			start, end := l, l+((r-l)/k)*k+k
			diff[start] = diff[start] * v % mod
			diff[end] = diff[end] * inv % mod
		}

		origin := make([]int, k)
		for i := range origin {
			origin[i] = 1
		}
		for i := 0; i < n; i++ {
			m := i % k
			origin[m] = origin[m] * diff[i] % mod
			resMul[i] = resMul[i] * origin[m] % mod
		}
	}

	res := 0
	for i, v := range nums {
		v = v * resMul[i] % mod
		res ^= v
	}
	return res
}
