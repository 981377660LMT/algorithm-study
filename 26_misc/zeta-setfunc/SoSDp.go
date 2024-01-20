// 给定一个集合，对该集合的所有子集，计算该子集的所有子集之和
//（这个「和」不一定是加法，可以是其它的满足合并性质的统计量）

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// CF165E()
	// CF383E()
	// CF449D()
	// CF1234F()
	ARC100E()

	// 维护集合最大和次大的写法()
	// demo()
}

// https://www.luogu.com.cn/problem/CF165E
// 给定一个数组.对每个nums[i]，问是否存在nums[j]满足nums[i]&nums[j]=0.
// 存在则输出nums[j],否则输出-1
// !nums[i]<=4e6
//
// !其实就是问两个二进制集合有没有交。
func CF165E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	log := max(bits.Len(uint(maxs(nums...))), 1)
	dp := make([]int, 1<<log)
	for i := range dp {
		dp[i] = -1
	}
	for _, v := range nums {
		dp[v] = v
	}

	SosDp1(log, func(cur, sub int) {
		if dp[sub] != -1 {
			dp[cur] = dp[sub]
		}
	})

	mask := 1<<log - 1
	for _, v := range nums {
		fmt.Fprint(out, dp[mask^v], " ")
	}
}

// https://www.luogu.com.cn/problem/CF449D
// 求nums中按位与为0的非空集合个数，模1e9+7.
// (取反后等价于按位或为全集.)
// nums[i]<=1e6.
//
// !dp[i] 表示子集包含 i 的方案数.
func CF449D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	pow := func(x, n int) int {
		x %= MOD
		res := 1 % MOD
		for ; n > 0; n >>= 1 {
			if n&1 == 1 {
				res = res * x % MOD
			}
			x = x * x % MOD
		}
		return res
	}

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	log := max(bits.Len(uint(maxs(nums...))), 1)
	counter := make([]int, 1<<log)
	for _, v := range nums {
		counter[v]++
	}
	SosDp2(log, func(cur, super int) {
		counter[cur] = (counter[cur] + counter[super]) % MOD
	})

	dp := make([]int, 1<<log) // dp[i] 表示子集包含 i 的方案数.
	for i := range dp {
		dp[i] = pow(2, counter[i])
	}
	SosDp2(log, func(cur, super int) {
		dp[cur] = (dp[cur] - dp[super] + MOD) % MOD
	})

	fmt.Fprintln(out, dp[0])
}

// https://www.luogu.com.cn/problem/CF383E
// 给出 n 个长度为 3 的单词（每个字母都为 a 到 x 中的任意一个），
// 我们说该单词是正确的，当且仅当该单词中含有至少一个元音。元音可以是范围内的任意字母。
// 对于所有2^24种元音集合，求出所有情况下正确单词个数的平方的异或和。
//
// !平方的异或和不太能维护，考虑暴力枚举子集计算答案。
// 正难则反，求非法的数量.
func CF383E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := range words {
		fmt.Fscan(in, &words[i])
	}

	log := 24 // 'a'-'x' 24 位
	MASK := 1<<log - 1
	dp := make([]int, 1<<log) // !非法个数
	for _, word := range words {
		state := 0
		for _, c := range word {
			state |= 1 << (c - 'a')
		}
		dp[state^MASK]++
	}
	SosDp2(log, func(cur, super int) {
		dp[cur] += dp[super]
	})
	res := 0
	for i := 0; i < 1<<log; i++ {
		x := n - dp[i]
		res ^= x * x
	}
	fmt.Fprintln(out, res)
}

// https://www.luogu.com.cn/problem/CF1234F
// 给你一个字符串S，你可以翻转一次S的任意一个子串。
// 问翻转后S的子串中各个字符都不相同的最长子串长度。
// 字符集<=20
//
// !翻转操作本质上是在 S 中找两个不相交子串，且拼接后不含相同字符。
// !等价于:给定一个数组，求满足ai&aj==0的ai|aj二进制中1的个数的最大值.
// 我们可以枚举 ai，求出 ai的补集的子集中的二进制中 1 的个数的最大值，
// 其与 ai的二进制中 1 的个数之和的最大值就是答案。
func CF1234F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	set := make(map[int]struct{})
	n := len(s)
	for left := 0; left < n; left++ {
		state := 0
		for right := left; right < n; right++ {
			cur := 1 << (s[right] - 'a')
			if state&cur > 0 {
				break
			} else {
				state |= cur
			}
			set[state] = struct{}{}
		}
	}

	nums := make([]int, 0, len(set))
	for state := range set {
		nums = append(nums, state)
	}

	// 给定一个数组，求满足ai&aj==0的ai|aj二进制中1的个数的最大值.
	// nums[i]<=1e6
	solve := func(nums []int) int {
		log := max(bits.Len(uint(maxs(nums...))), 1)
		dp := make([]int, 1<<log)
		for _, v := range nums {
			dp[v] = bits.OnesCount(uint(v))
		}

		SosDp1(log, func(cur, sub int) {
			dp[cur] = max(dp[cur], dp[sub])
		})

		res := 0
		MASK := 1<<log - 1
		for s := range dp {
			res = max(res, dp[s]+dp[s^MASK])
		}
		return res
	}

	fmt.Fprintln(out, solve(nums))
}

// https://www.luogu.com.cn/problem/AT_arc100_c
// 给定一个长为2^n的数组，对每个1<=k<=2^n-1，
// 找出最大的nums[i]+nums[j] (i|j<=k)
// n<=18
// 维护集合最大和次大即可，最终输出的每一个答案都是一个前缀最大值。
func ARC100E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, 1<<n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	log := n
	type pair struct{ max1, max2 int }
	dp := make([]pair, 1<<log)
	for i, v := range nums {
		dp[i] = pair{max1: v, max2: -1}
	}
	SosDp1(log, func(cur, sub int) {
		curDp, subDp := dp[cur], dp[sub]
		if subDp.max2 > curDp.max1 {
			dp[cur] = subDp
		} else if subDp.max1 > curDp.max1 {
			dp[cur] = pair{subDp.max1, curDp.max1}
		} else if subDp.max1 > curDp.max2 {
			dp[cur].max2 = subDp.max1
		}
	})

	res := 0
	for i := 1; i < 1<<log; i++ {
		res = max(res, dp[i].max1+dp[i].max2)
		fmt.Fprintln(out, res)
	}
}

// https://atcoder.jp/contests/arc136/tasks/arc136_d
// 10进制的情形

// 与运算为0的二元组个数
// nums[i] <= 1e6
func 与运算为0的二元组(nums []int) int {
	log := max(bits.Len(uint(maxs(nums...))), 1)
	dp := make([]int, 1<<log)
	for _, v := range nums {
		dp[v]++
	}
	SosDp1(log, func(cur, sub int) {
		dp[cur] += dp[sub]
	})
	mask := 1<<log - 1
	res := 0
	for _, v := range nums {
		res += dp[mask^v]
	}
	return res / 2
}

// 从子集转移的 SOS DP
// SubSetZeta
func SosDp1(log int, f func(cur, sub int)) {
	for i := 0; i < log; i++ {
		mask := 1 << i
		for s := 0; s < 1<<log; s++ {
			s |= mask
			f(s, s^mask) // 将 s 的子集 s^1<<i 的统计量合并到 s 中
		}
	}
}

// 从超集转移的 SOS DP
// SuperSetZeta
func SosDp2(log int, f func(cur, super int)) {
	for i := 0; i < log; i++ {
		mask := 1 << i
		for s := 0; s < 1<<log; s++ {
			t := s ^ mask
			if s < t {
				f(s, t) // 将 s 的超集 s|1<<i 的统计量合并到 s 中
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
