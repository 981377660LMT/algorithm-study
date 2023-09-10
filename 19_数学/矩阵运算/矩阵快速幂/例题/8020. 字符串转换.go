// 8020. 字符串轮转
// https://leetcode.cn/problems/string-transformation/description/
// 给你两个长度都为 n 的字符串 s 和 t 。你可以对字符串 s 执行以下操作：
//
// 将 s 长度为 l （0 < l < n）的 后缀字符串 删除，并将它添加在 s 的开头。
// 比方说，s = 'abcd' ，那么一次操作中，你可以删除后缀 'cd' ，并将它添加到 s 的开头，得到 s = 'cdab' 。
// 给你一个整数 k ，请你返回 恰好 k 次操作将 s 变为 t 的方案数。
//
// 由于答案可能很大，返回答案对 109 + 7 取余 后的结果。
//
// !记dp[i][0/1]为 `i` 次操作后 `等于/不等于t` 的方案数，count 为 `t` 在 `s` 的循环轮转中出现的次数
// !dp[i][0] = dp[i-1][0]*(count-1) + dp[i-1][1]*count
// !dp[i][1] = dp[i-1][0]*(n-count) + dp[i-1][1]*(n-count-1)
// 因此有状态转移方程
// dp[i] = T*dp[i-1],
//
//	 T = [
//		    [count-1, count],
//		    [n-count, n-count-1]
//		   ]

package main

import "fmt"

func main() {
	// s = "ababab", t = "ababab", k = 1
	fmt.Println(numberOfWays("ababab", "ababab", 1))
	// "abcd", t = "cdab", k = 2
	fmt.Println(numberOfWays("abcd", "cdab", 2))
}

const MOD int = 1e9 + 7

func numberOfWays(s string, t string, k int64) int {
	// !统计t在s的循环轮转中出现的次数,不包含(s==t)(即在(s+s)[1:]中出现的次数)
	countTInCyclicS := func(s, t string) int {
		allPos := IndexOfAll(s+s, t, 1)
		return len(allPos)
	}

	n := len(s)
	count := countTInCyclicS(s, t)

	var init M
	if s == t {
		init = M{{1}, {0}}
	} else {
		init = M{{0}, {1}}
	}
	T := M{
		{count - 1, count},
		{n - count, n - count - 1},
	}
	resT := MatPow(T, int(k), MOD)
	res := MatMul(resT, init, MOD)
	return res[0][0]
}

type M = [][]int

//
//
//
func NewIdentityMatrix(n int) M {
	res := make(M, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
		res[i][i] = 1
	}
	return res
}

func NewMatrix(row, col int) M {
	res := make(M, row)
	for i := range res {
		res[i] = make([]int, col)
	}
	return res
}

func MatMul(m1, m2 M, mod int) M {
	res := NewMatrix(len(m1), len(m2[0]))
	for i := 0; i < len(m1); i++ {
		for k := 0; k < len(m2); k++ {
			for j := 0; j < len(m2[0]); j++ {
				res[i][j] = (res[i][j] + m1[i][k]*m2[k][j]) % mod
				if res[i][j] < 0 {
					res[i][j] += mod
				}
			}
		}
	}
	return res
}

// matPow/matqpow
func MatPow(m1 M, exp, mod int) M {
	n := len(m1)
	e := NewIdentityMatrix(n)
	b := make(M, n)
	for i := 0; i < n; i++ {
		b[i] = make([]int, n)
		copy(b[i], m1[i])
	}
	for exp > 0 {
		if exp&1 == 1 {
			e = MatMul(e, b, mod)
		}
		b = MatMul(b, b, mod)
		exp >>= 1
	}
	return e
}

func GetNext(pattern string) []int {
	next := make([]int, len(pattern))
	j := 0
	for i := 1; i < len(pattern); i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = next[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		next[i] = j
	}
	return next
}

// `O(n+m)` 寻找 `shorter` 在 `longer` 中的所有匹配位置.
func IndexOfAll(longer string, shorter string, position int) []int {
	if len(shorter) == 0 {
		return []int{0}
	}
	if len(longer) < len(shorter) {
		return nil
	}
	res := []int{}
	next := GetNext(shorter)
	hitJ := 0
	for i := position; i < len(longer); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = next[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == len(shorter) {
			res = append(res, i-len(shorter)+1)
			hitJ = next[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}
