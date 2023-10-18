package main

import (
	"fmt"
	"time"
)

const MOD int = 1e9 + 7

// 1714. 数组中特殊等间距元素的和
// https://leetcode.cn/problems/sum-of-special-evenly-spaced-elements-in-array/
func solve(nums []int, queries [][]int) []int {
	n := len(nums)
	R := NewRangeStepSum(nums, 40)
	res := make([]int, len(queries))
	for qi, query := range queries {
		start, stop, step := query[0], n, query[1]
		res[qi] = R.Query(start, stop, step) % MOD
	}
	return res
}

type RangeStepSum struct {
	nums          []int
	stepThreshold int
	// dp[step][start] 表示步长为step,起点为start的所有元素的和.
	// `dp[step][start] = dp[step][start+step] + nums[start]`.
	dp [][]int
}

// stepThreshold: 步长阈值,当步长小于等于该值时,使用dp数组预处理答案,否则直接遍历.
// 预处理时间空间复杂度均为`O(n*stepThreshold)`.
// 单次遍历时间复杂度为`O(n/stepThreshold)`.
// 取40较为合适.
func NewRangeStepSum(arr []int, stepThreshold int) *RangeStepSum {
	n := len(arr)
	dp := make([][]int, stepThreshold)
	for step := 1; step <= stepThreshold; step++ {
		curSum := make([]int, n+1)
		for start := n - 1; start >= 0; start-- {
			curSum[start] = curSum[min(n, start+step)] + arr[start]
		}
		dp[step-1] = curSum
	}
	return &RangeStepSum{nums: arr, stepThreshold: stepThreshold, dp: dp}
}

func (rss *RangeStepSum) Query(start, stop, step int) int {
	if start < 0 {
		start = 0
	}
	if stop > len(rss.nums) {
		stop = len(rss.nums)
	}
	if start >= stop {
		return 0
	}
	if step <= rss.stepThreshold {
		curDp := rss.dp[step-1]
		// 找到 >=stop 的第一个形为start+k*step的数
		div := (stop - start + step - 1) / step
		nextStart := min(start+div*step, len(rss.nums))
		return curDp[start] - curDp[nextStart]
	}
	sum := 0
	for i := start; i < stop; i += step {
		sum += rss.nums[i]
	}
	return sum
}

func (rss *RangeStepSum) String() string {
	return fmt.Sprintf("RangeStepSum{%v}", rss.nums)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {

	checkWithBruteForce := func() {
		n := int(1e3)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = i
		}
		rss := NewRangeStepSum(arr, 100)
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				for step := 1; step <= 100; step++ {
					expect := 0
					for k := i; k < j; k += step {
						expect += arr[k]
					}
					actual := rss.Query(i, j, step)
					if expect != actual {
						panic(fmt.Sprintf("expect %d, actual %d", expect, actual))
					}
				}
			}
		}
		fmt.Println("checkWithBruteForce passed!")
	}

	testTime := func() {
		n := int(2e5)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = i
		}
		time1 := time.Now()
		rss := NewRangeStepSum(arr, 80)
		for i := 0; i < n; i++ {
			rss.Query(0, n, i+1)
		}
		time2 := time.Now()
		fmt.Println(time2.Sub(time1))
	}

	checkWithBruteForce()
	testTime()

}
