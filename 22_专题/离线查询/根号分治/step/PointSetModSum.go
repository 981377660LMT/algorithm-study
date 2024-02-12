// PointSetModSum
// Set(index, value) 修改第index个数为value.
// Query(start, step) 询问所有下标模 step 的结果为 start 的位置的值之和.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	CF1207F()
	// luogu3396()
}

// CF1207F-Remainder Problem
// https://www.luogu.com.cn/problem/CF1207F
// 给你一个长为5e5的数组，初始时每个数都是0，现在有两种操作：
// 1 pos delta：将第pos个数加上delta
// 2 m v：询问所有下标模m的结果为v的位置的值之和
func CF1207F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const N = 500000
	const SQRT = 200
	nums := [N + 1]int{}
	dp := [SQRT][SQRT]int{}

	add := func(index, delta int) {
		nums[index] += delta
		for mod := 1; mod < SQRT; mod++ {
			dp[mod][index%mod] += delta
		}
	}

	query := func(start, mod int) int {
		if mod < SQRT {
			return dp[mod][start]
		}
		sum := 0
		for i := start; i <= N; i += mod {
			sum += nums[i]
		}
		return sum
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var pos, delta int
			fmt.Fscan(in, &pos, &delta)
			add(pos, delta)
		} else {
			var m, v int
			fmt.Fscan(in, &m, &v)
			fmt.Fprintln(out, query(v, m))
		}
	}
}

// P3396 哈希冲突
// https://www.luogu.com.cn/problem/P3396
// https://www.luogu.com.cn/blog/danieljiang/ha-xi-chong-tu-ti-xie-gen-hao-ke-ji
func luogu3396() {
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
	// dp[step][mod] 表示步长为step时，模为mod的所有数之和.
	dp [][]int
}

// stepThreshold: 步长阈值,当步长小于等于该值时,使用dp数组预处理答案,否则直接遍历.
// 预处理时间复杂度为`O(n*stepThreshold)`, 空间复杂度为O(`stepThreshold^2`)
// 单次遍历时间复杂度为`O(n/stepThreshold)`.
// 取50较为合适.
func NewPointSetModSum(nums []int, stepThreshold int) *PointSetModSum {
	if stepThreshold == -1 {
		stepThreshold = int(math.Sqrt(float64(len(nums))) + 1)
	}
	nums = append(nums[:0:0], nums...)
	n := len(nums)
	dp := make([][]int, stepThreshold)
	for step := 1; step <= stepThreshold; step++ {
		curSum := make([]int, step)
		for i := 0; i < n; i++ {
			curSum[i%step] += nums[i]
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
// !start<step.
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

func (pss *PointSetModSum) Get(index int) int {
	if index < 0 || index >= len(pss.nums) {
		return 0
	}
	return pss.nums[index]
}

func (pss *PointSetModSum) Add(index, delta int) {
	if index < 0 || index >= len(pss.nums) {
		return
	}
	if delta == 0 {
		return
	}
	pss.nums[index] += delta
	for step := 1; step <= pss.stepThreshold; step++ {
		pss.dp[step-1][index%step] += delta
	}
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
