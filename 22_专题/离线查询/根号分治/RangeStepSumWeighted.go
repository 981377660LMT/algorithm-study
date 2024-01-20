// https://codeforces.com/contest/1921/problem/F
// F. Sum of Progression
// 带权前缀和

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://codeforces.com/contest/1921/problem/F
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	solve := func() {
		var n, q int
		fmt.Fscan(in, &n, &q)
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &nums[i])
		}

		S := NewRangeStepSumWeighted(nums, 80)
		for i := 0; i < q; i++ {
			var start, step, length int
			fmt.Fscan(in, &start, &step, &length)
			start--
			fmt.Fprint(out, S.Query(start, start+(length-1)*step+1, step), " ")
		}
		fmt.Fprintln(out)
	}

	for i := 0; i < T; i++ {
		solve()
	}
}

type RangeStepSumWeighted struct {
	nums          []int
	stepThreshold int
	pre           [][]int
	sum           [][]int
}

// stepThreshold: 步长阈值,当步长小于该值时,使用dp数组预处理答案,否则直接遍历.
// 预处理时间空间复杂度均为`O(n*stepThreshold)`.
// 单次遍历时间复杂度为`O(n/stepThreshold)`.
// 取80较为合适.
func NewRangeStepSumWeighted(arr []int, stepThreshold int) *RangeStepSumWeighted {
	arr = append(arr[:0:0], arr...)
	n := len(arr)
	pre, sum := make([][]int, stepThreshold), make([][]int, stepThreshold)
	for i := 0; i < stepThreshold; i++ {
		pre[i] = make([]int, n+stepThreshold)
		sum[i] = make([]int, n+stepThreshold)
	}
	for step := 1; step < stepThreshold; step++ {
		tmpPre, tmpSum := pre[step], sum[step]
		for i := 0; i < n; i++ {
			tmpPre[i+step] = tmpPre[i] + arr[i]
			tmpSum[i+step] = tmpSum[i] + arr[i]*(i/step+1)
		}
	}
	return &RangeStepSumWeighted{nums: arr, stepThreshold: stepThreshold, pre: pre, sum: sum}
}

// 求 nums[start] + 2*nums[start+step] + 3*nums[start+2*step] + ... + 的和.(带权前缀和)
func (rss *RangeStepSumWeighted) Query(start, stop, step int) int {
	if start < 0 {
		start = 0
	}
	if stop > len(rss.nums) {
		stop = len(rss.nums)
	}
	if start >= stop {
		return 0
	}
	if step < rss.stepThreshold {
		count := (stop-start-1)/step + 1
		last := start + count*step
		return rss.sum[step][last] - rss.sum[step][start] - (rss.pre[step][last]-rss.pre[step][start])*(start/step)
	} else {
		sum := 0
		weight := 1
		for i := start; i < stop; i += step {
			sum += rss.nums[i] * weight
			weight++
		}
		return sum
	}
}

func (rss *RangeStepSumWeighted) String() string {
	return fmt.Sprintf("RangeStepSumWeighted{%v}", rss.nums)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
