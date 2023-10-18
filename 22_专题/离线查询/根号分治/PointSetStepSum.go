// P3396 哈希冲突
// https://www.luogu.com.cn/problem/P3396
// https://www.luogu.com.cn/blog/danieljiang/ha-xi-chong-tu-ti-xie-gen-hao-ke-ji

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

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	S := NewPointSetModSum(nums, 50)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "A" {
			var step, start int
			fmt.Fscan(in, &step, &start)
			start--
			if start < 0 {
				start += step
			}
			fmt.Fprintln(out, S.Query(start, step))
		} else {
			var pos, target int
			fmt.Fscan(in, &pos, &target)
			pos--
			S.Set(pos, target)
		}
	}
}

type PointSetModSum struct {
	nums          []int
	stepThreshold int
	// dp[step][start] 表示步长为step,起点为start的所有元素的和.
	// `dp[step][start] = dp[step][start+step] + nums[start]`.
	dp [][]int
}

// stepThreshold: 步长阈值,当步长小于等于该值时,使用dp数组预处理答案,否则直接遍历.
// 预处理时间空间复杂度均为`O(n*stepThreshold)`.
// 单次遍历时间复杂度为`O(n/stepThreshold)`.
// 取50较为合适.
func NewPointSetModSum(nums []int, stepThreshold int) *PointSetModSum {
	n := len(nums)
	dp := make([][]int, stepThreshold)
	for step := 1; step <= stepThreshold; step++ {
		curSum := make([]int, n+1)
		for start := n - 1; start >= 0; start-- {
			curSum[start] = curSum[min(n, start+step)] + nums[start]
		}
		dp[step-1] = curSum
	}
	return &PointSetModSum{nums: nums, stepThreshold: stepThreshold, dp: dp}
}

func (pss *PointSetModSum) Set(index, value int) {
	if index < 0 || index >= len(pss.nums) {
		return
	}
	pre := pss.nums[index]
	if pre == value {
		return
	}
	pss.nums[index] = value
	delta := value - pre
	for step := 1; step <= pss.stepThreshold; step++ {
		pss.dp[step-1][index%step] += delta
	}
}

// 查询 sum(nums[start::step]).
func (pss *PointSetModSum) Query(start, step int) int {
	if start < 0 {
		start = 0
	}
	if step <= pss.stepThreshold {
		return pss.dp[step-1][start]
	}
	sum := 0
	for i := start; i < len(pss.nums); i += step {
		sum += pss.nums[i]
	}
	return sum
}

func (pss *PointSetModSum) String() string {
	return fmt.Sprintf("PointSetModSum{nums: %v}", pss.nums)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
