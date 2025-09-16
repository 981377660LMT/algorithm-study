// Zeta transform and Mobius transform on set functions
// SOS DP (Sum over Subsets)
// SuperSet Zeta/Mobius Transform

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// Enumerate(3, func(s1, s2 int) {
	// 	// do something
	// 	if s1 > s2 {
	// 		fmt.Println(s1, s2)
	// 	}
	// })

	// nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	// SubsetZeta(nums)
	// fmt.Println(nums)

	abc423_f()
}

const INF int = 4e18

// F - Loud Cicada
// https://atcoder.jp/contests/abc423/tasks/abc423_f
// 在 AtCoder 岛上，栖息着 N 种蝉。第 i 种蝉 (1≤i≤N) 仅在年份是 A_i 的倍数时才会大量出现。
// 请计算在 1 年到 Y 年这 Y 年间，有多少个年份恰好有 M 种蝉大量出现。
// 1 ≤ M ≤ N ≤ 20
// 1 ≤ Y ≤ 10¹⁸
// 1 ≤ Aᵢ ≤ 10¹⁸ (1≤i≤N)
// 所有输入均为整数。
func abc423_f() {
	gcd := func(a, b int) int {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}
	lcm := func(a, b int) int {
		if a == INF || b == INF {
			return INF
		}
		a /= gcd(a, b)
		if a >= (INF+b-1)/b {
			return INF
		}
		return a * b
	}

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, y int
	fmt.Fscan(in, &n, &m, &y)
	subsetLcm := make([]int, 1<<n)
	subsetLcm[0] = 1
	for i := 0; i < n; i++ {
		var a int
		fmt.Fscan(in, &a)
		for s := 0; s < (1 << i); s++ {
			subsetLcm[s|(1<<i)] = lcm(subsetLcm[s], a)
		}
	}

	dp := make([]int, 1<<n) // 每个子集的“至少”发生次数
	for i := 0; i < (1 << n); i++ {
		dp[i] = y / subsetLcm[i]
	}
	SupersetMobius(dp) // 容斥
	res := 0
	for i := 0; i < (1 << n); i++ {
		if bits.OnesCount(uint(i)) == m {
			res += dp[i]
		}
	}
	fmt.Fprintln(out, res)
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
