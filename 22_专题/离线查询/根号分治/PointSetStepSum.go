package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc335f()
	// demo()
	// luogu3396()

}

func demo() {
	S := NewPointSetStepSum([]int{1, 1, 1, 1, 1}, 50)
	fmt.Println(S)
	// fmt.Println(S.Query(0, 2)
	S.Set(0, 2)
	fmt.Println(S)

}

// F - Hop Sugoroku
// https://atcoder.jp/contests/abc335/tasks/abc335_f
// 加速:
//
// dp = [1]*n
//
// for i in range(n-1,-1,-1):
//
//	for j in range(i+A[i],n,A[i]):
//		dp[i] += dp[j]
//
// print(dp[0])
func abc335f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	dp := make([]int, n)
	for i := 0; i < n; i++ {
		dp[i] = 1
	}
	S := NewPointSetStepSum(dp, 50)
	for i := n - 1; i >= 0; i-- {
		sum := S.Query(i+nums[i], nums[i]) % MOD
		S.Add(i, sum)
		fmt.Println(S, sum, i+nums[i], nums[i], S.Query(1, 1))
	}
	fmt.Fprintln(out, S.Get(0)%MOD)
	// fmt.Println(S)
	// dp := make([]int, n)
	// for i := 0; i < n; i++ {
	// 	dp[i] = 1
	// }
	// for i := n - 1; i >= 0; i-- {
	// 	for j := i + nums[i]; j < n; j += nums[i] {
	// 		dp[i] += dp[j]
	// 		dp[i] %= MOD
	// 	}
	// }
	// fmt.Fprintln(out, dp[0])

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

	S := NewPointSetStepSum(nums, 50)
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

type PointSetStepSum struct {
	nums          []int
	stepThreshold int
	// dp[step][mod] 表示步长为step时，模为mod的所有数之和.
	dp [][]int
}

// stepThreshold: 步长阈值,当步长小于等于该值时,使用dp数组预处理答案,否则直接遍历.
// 预处理时间空间复杂度均为`O(n*stepThreshold)`.
// 单次遍历时间复杂度为`O(n/stepThreshold)`.
// 取50较为合适.
func NewPointSetStepSum(nums []int, stepThreshold int) *PointSetStepSum {
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
	return &PointSetStepSum{nums: nums, stepThreshold: stepThreshold, dp: dp}
}

func (pss *PointSetStepSum) Set(index, value int) {
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

func (pss *PointSetStepSum) Get(index int) int {
	if index < 0 || index >= len(pss.nums) {
		return 0
	}
	return pss.nums[index]
}

func (pss *PointSetStepSum) Add(index, delta int) {
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

// 查询 sum(nums[start::step]).
// !step<=start.
func (pss *PointSetStepSum) Query(start, step int) int {
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

func (pss *PointSetStepSum) String() string {
	return fmt.Sprintf("PointSetModSum{nums: %v}", pss.nums)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
