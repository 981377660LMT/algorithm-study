package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	P2371()
}

// P2371 [国家集训队] 墨墨的等式
// https://www.luogu.com.cn/problem/P2371
// No.1782 ManyCoins
// https://yukicoder.me/submissions/1008663
// 给定n个系数coeffs和上下界lower,upper
// !对于 lower<=k<=upper 求有多少个k能够满足
// !a0*x0+a1*x1+...+an*xn=k
// n<=12 0<=ai<=5e5 1<=lower<=upper<=2^63-1
// !时间复杂度：O(n*ai)
func P2371() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, lower, upper int
	fmt.Fscan(in, &n, &lower, &upper)
	coeffs := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &coeffs[i])
	}

	retain(&coeffs, func(i int) bool { return coeffs[i] != 0 })
	if len(coeffs) == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	base, dp := ModShortestPath(coeffs)
	res := 0
	for i := range base {
		if upper >= dp[i] {
			res += (upper-dp[i])/base + 1
		}
		if lower > dp[i] {
			res -= (lower-1-dp[i])/base + 1
		}
	}
	fmt.Fprintln(out, res)
}

const INF int = 1e18

// 确定线性组合∑ai*xi的可能取到的值(ai非负).
// base:最小的非零ai
// dp[i]:dp[i]记录的是最小的x,满足x=i(mod base)且x能被系数coeffs线性表出(xi非负)
func ModShortestPath(coeffs []int) (int, []int) {
	coeffs = append(coeffs[:0:0], coeffs...)
	retain(&coeffs, func(i int) bool { return coeffs[i] > 0 })
	if len(coeffs) == 0 {
		return 0, nil
	}

	base := mins(coeffs...)
	dp := make([]int, base)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for _, v := range coeffs {
		cycle := gcd(base, v)
		for j := 0; j < cycle; j++ {
			for cur, count := j, 0; count < 2; {
				next := (cur + v) % base
				dp[next] = min(dp[next], dp[cur]+v)
				cur = next
				if cur == j {
					count++
				}
			}
		}
	}
	return base, dp
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
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

func retain(arr *[]int, f func(index int) bool) {
	ptr := 0
	n := len(*arr)
	for i := 0; i < n; i++ {
		if f(i) {
			(*arr)[ptr] = (*arr)[i]
			ptr++
		}
	}
	*arr = (*arr)[:ptr:ptr]
}
