// 可能的子集和,子集和输出方案
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// !能否分成两个和相等的子集(天平称重问题)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
		sum += nums[i]
	}

	_, ok := subsetSum(nums, sum/2)
	can := ok || sum == 0
	if can && sum%2 == 0 {
		fmt.Fprintln(out, "possible")
	} else {
		fmt.Fprintln(out, "impossible")
	}
}

// 能否用nums中的若干个数凑出和为target
//  O(n*max(nums)))
//  target小于等于0时返回无解
func subsetSum(nums []int, target int) (res []int, ok bool) {
	if target <= 0 {
		return
	}

	n := len(nums)
	max_ := 0
	for _, v := range nums {
		max_ = max(max_, v)
	}
	right, curSum := 0, 0
	for right < n && curSum+nums[right] <= target {
		curSum += nums[right]
		right++
	}
	if right == n && curSum != target {
		return
	}

	offset := target - max_ + 1
	dp := make([]int, 2*max_)
	for i := range dp {
		dp[i] = -1
	}
	pre := make([][]int, n)
	for i := range pre {
		pre[i] = make([]int, 2*max_)
		for j := range pre[i] {
			pre[i][j] = -1
		}
	}

	dp[curSum-offset] = right
	for i := right; i < n; i++ {
		ndp := make([]int, len(dp))
		copy(ndp, dp)
		p := pre[i]
		a := nums[i]
		for j := 0; j < max_; j++ {
			if ndp[j+a] < dp[j] {
				ndp[j+a] = dp[j]
				p[j+a] = -2
			}
		}
		for j := 2*max_ - 1; j >= max_; j-- {
			for k := ndp[j] - 1; k >= max(dp[j], 0); k-- {
				if ndp[j-nums[k]] < k {
					ndp[j-nums[k]] = k
					p[j-nums[k]] = k
				}
			}
		}
		dp = ndp
	}

	if dp[max_-1] == -1 {
		return
	}

	used := make([]bool, n)
	i, j := n-1, max_-1
	for i >= right {
		p := pre[i][j]
		if p == -2 {
			used[i] = !used[i]
			j -= nums[i]
			i--
		} else if p == -1 {
			i--
		} else {
			used[p] = !used[p]
			j += nums[p]
		}
	}

	for i >= 0 {
		used[i] = !used[i]
		i--
	}

	for i := 0; i < n; i++ {
		if used[i] {
			res = append(res, i)
		}
	}

	ok = true
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
