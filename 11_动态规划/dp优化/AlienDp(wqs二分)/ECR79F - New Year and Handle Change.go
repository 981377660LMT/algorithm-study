// 给出一个01串(表现为大小写），可以最多选择 k 个 l 长度的子串，全部变为0或1。
// 求操作后min(0的个数，1的个数)的最小值。
// 1<=n,k,l<=1e6

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func NewYearandHandleChange(bits []byte, k, subLen int) int {
	n := len(bits)
	getDp := func(penalty int) [2]int {
		dp := make([][2]int, n+1) // dp[pos]:前pos个字符中1的最大次数,最大操作使用次数
		for i := range dp {
			dp[i][0] = -INF
		}

		dp[0] = [2]int{0, 0}
		for i := 0; i < n; i++ {
			next := min(i+subLen, n)
			cand := dp[i][0] + (next - i) - penalty // 用操作
			if cand > dp[next][0] {
				dp[next] = [2]int{cand, dp[i][1] + 1}
			} else if cand == dp[next][0] {
				dp[next][1] = max(dp[next][1], dp[i][1]+1) // !使用次数最大
			}

			cand = dp[i][0] + int(bits[i]) // 不用操作
			if cand > dp[i+1][0] {
				dp[i+1] = [2]int{cand, dp[i][1]}
			} else if cand == dp[i+1][0] {
				dp[i+1][1] = max(dp[i+1][1], dp[i][1]) // !使用次数最大
			}
		}

		return dp[n]
	}

	res1 := n - AliensDp(k, getDp)
	for i := 0; i < n; i++ {
		bits[i] ^= 1
	}
	res2 := n - AliensDp(k, getDp)
	return min(res1, res2)
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

// getDp：func(penalty int) [2]int: 每使用一次操作罚款 penalty 元, 返回 [子问题dp的`最大值`, `最大的`操作使用次数].
func AliensDp(k int, getDp func(penalty int) [2]int) int {
	left, right := 1, int(1e18)
	penalty := 0
	for left <= right {
		mid := (left + right) >> 1
		if cand := getDp(mid); cand[1] >= k {
			penalty = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	res := getDp(penalty)
	return res[0] + penalty*k
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, subLen int
	fmt.Fscan(in, &n, &k, &subLen)
	var s string
	fmt.Fscan(in, &s)
	bits := make([]byte, n)
	for i := 0; i < n; i++ {
		if 'a' <= s[i] && s[i] <= 'z' {
			bits[i] = 1
		} else {
			bits[i] = 0
		}
	}

	fmt.Fprintln(out, NewYearandHandleChange(bits, k, subLen))
}
