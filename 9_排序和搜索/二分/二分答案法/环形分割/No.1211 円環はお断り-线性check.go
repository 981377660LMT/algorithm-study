// No.1211 円環はお断り(圆环，环上最大化最小值)
// https://yukicoder.me/problems/no/1211
// !给定一个环形数组,分成k个非空连续子数组,使得这k个子数组的和的最小值最大,求出这个最大的最小值.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	presum := make([]int, n+1)
	for i := 0; i < n; i++ {
		presum[i+1] = presum[i] + nums[i]
	}

	cost := func(start, end int) int {
		return presum[end] - presum[start]
	}

	check := func(mid int) bool {
		{
			// 先求解链上的问题(剪枝)
			count := 0
			left := 0
			for right := 0; right < n; right++ {
				if cost(left, right+1) >= mid {
					count++
					left = right + 1
				}
			}
			if count >= k {
				return true
			}
			if count <= k-2 {
				return false
			}
		}

		next := make([]int, n)
		right := 0
		for left := 0; left < n; left++ {
			for right < n && cost(left, right) < mid {
				right++
			}
			if cost(left, right) >= mid {
				next[left] = right
			} else {
				next[left] = -1
			}
		}

		type dpItem struct{ count, next int }
		dp := make([]dpItem, n+1)
		dp[n] = dpItem{next: n}
		for i := n - 1; i >= 0; i-- {
			if next[i] == -1 {
				dp[i] = dpItem{next: i}
			} else {
				dp[i] = dp[next[i]]
				dp[i].count++
			}
		}

		for i := 0; i < n; i++ {
			count := dp[i].count
			if count <= k-2 {
				break
			}
			end := dp[i].next
			if cost(0, i)+cost(end, n) >= mid {
				count++
			}
			if count >= k {
				return true
			}
		}

		return false
	}

	left, right := 1, presum[n]/k+1
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// !此时最小值为right.
	fmt.Fprintln(out, right)
}
