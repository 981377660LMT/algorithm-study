// E. LIS of Sequence
// https://www.luogu.com.cn/problem/CF486E
// 给你一个长度为n的序列a1,a2,...,an，你需要把这n个元素分成三类：1，2，3：
// 1:所有的最长上升子序列都不包含这个元素
// 2:有但非所有的最长上升子序列包含这个元素
// 3:所有的最长上升子序列都包含这个元素
// !注意lis是严格上升的

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	res := LISOfSequnce(nums, true)
	for _, v := range res {
		fmt.Fprint(out, v, "")
	}
}

type State uint8

const (
	NotInLIS State = 1 + iota
	InSomeLIS
	InAllLIS
)

func LISOfSequnce(arr []int, strict bool) []State {
	arr = append(arr[:0:0], arr...)
	lis, dp1 := LISDp(arr, strict)
	for i, j := int32(0), int32(len(arr)-1); i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	for i := range arr {
		arr[i] = -arr[i]
	}
	_, dp2 := LISDp(arr, strict)
	for i, j := int32(0), int32(len(arr)-1); i < j; i, j = i+1, j-1 {
		dp2[i], dp2[j] = dp2[j], dp2[i]
	}

	n := int32(len(arr))
	counter := make([]int32, n)
	for i := int32(0); i < n; i++ {
		if dp1[i]+dp2[i] == lis {
			counter[dp1[i]]++
		}
	}
	res := make([]State, n)
	for i := int32(0); i < n; i++ {
		if dp1[i]+dp2[i] < lis {
			res[i] = NotInLIS
		} else if counter[dp1[i]] == 1 {
			res[i] = InAllLIS
		} else {
			res[i] = InSomeLIS
		}
	}
	return res
}

const INF int = 2e18

// 返回每个位置为结尾的LIS长度(不包括自身).
func LISDp(nums []int, strict bool) (int32, []int32) {
	n := int32(len(nums))
	dp := make([]int, n)
	for i := range dp {
		dp[i] = INF
	}
	lis := int32(0)
	lisRank := make([]int32, n)
	var f func([]int, int) int
	if strict {
		f = sort.SearchInts
	} else {
		f = func(a []int, x int) int {
			return sort.SearchInts(a, x+1)
		}
	}
	for i := int32(0); i < n; i++ {
		pos := int32(f(dp, nums[i]))
		dp[pos] = nums[i]
		lisRank[i] = pos
		if lis < pos {
			lis = pos
		}
	}
	return lis, lisRank
}
