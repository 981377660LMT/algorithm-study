// RangeStepAddRangeSum
// 注意 !start<step.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	demo()
}

func demo() {
	S := NewRangeStepAddRangeSum(make([]int, 10), 30)
	fmt.Println(S.Query(0, 10))
	S.Update(3, 2, 1)
	fmt.Println(S.Query(0, 10)) // 1
}

const MOD int = 1e9 + 7

// 0 start step delta: 将 start, start+step, start+2*step, ... 加上delta.
// !0<=start<step
// 1 start end：查询[start, end)的和.
//
// !根号分治.对step的大小进行分治.
// 如果step>=根号n, 则直接暴力修改;
// 否则, 就以step为跳的周期，每个点统计累计修改总和.
// !为了优化，我们采用前缀后缀统计方法。(PointAddRangeSum O(1)查询O(n)修改)
// 1.通过前后缀和可以解决单点改、区间查的问题
// 2.维护原数组、分块数组和周期的前后缀和，修改时修改块或者周期的前后缀和，
// 查询时统计两侧不完整块、中间完整块和不同周期在查询区间内的修改总和即可。
func 初始化() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	S := NewRangeStepAddRangeSum(nums, 70)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var step, start, delta int
			fmt.Fscan(in, &step, &start, &delta)
			start--
			S.Update(start, step, delta)
		} else {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			fmt.Fprintln(out, S.Query(start, end))
		}
	}

}

type RangeStepAddRangeSum struct {
	nums                         []int
	stepThreshold                int
	blockCount                   int
	belong, blockStart, blockEnd []int
	blockSum                     []int
	pre, suf                     [][]int
}

func NewRangeStepAddRangeSum(nums []int, stepThreshold int) *RangeStepAddRangeSum {
	n := len(nums)
	newNums := make([]int, n)
	for i, v := range nums {
		newNums[i] = v % MOD
		if newNums[i] < 0 {
			newNums[i] += MOD
		}
	}
	blockCount := 1 + (n / stepThreshold)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * stepThreshold
		tmp := (i + 1) * stepThreshold
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / stepThreshold
	}
	blockSum := make([]int, blockCount)
	for i := range blockSum {
		for j := blockStart[i]; j < blockEnd[i]; j++ {
			blockSum[i] += newNums[j]
			if blockSum[i] >= MOD {
				blockSum[i] -= MOD
			}
		}
	}
	pre, suf := make([][]int, stepThreshold+1), make([][]int, stepThreshold+1)
	for i := range pre {
		pre[i] = make([]int, stepThreshold+1)
		suf[i] = make([]int, stepThreshold+1)
	}

	return &RangeStepAddRangeSum{
		nums:          newNums,
		stepThreshold: stepThreshold,
		blockCount:    blockCount,
		belong:        belong,
		blockStart:    blockStart,
		blockEnd:      blockEnd,
		blockSum:      blockSum,
		pre:           pre,
		suf:           suf,
	}
}

func (ss *RangeStepAddRangeSum) Query(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ss.nums) {
		end = len(ss.nums)
	}
	if start >= end {
		return 0
	}
	res := ss._sum(start, end)
	for step := 1; step < ss.stepThreshold; step++ {
		id1, id2 := start/step, (end-1)/step
		pos1, pos2 := start%step, (end-1)%step
		curPre, curSuf := ss.pre[step], ss.suf[step]
		if id1 == id2 {
			res += curPre[pos2+1] - curPre[pos1]
		} else {
			res += curSuf[pos1] + (id2-id1-1)*curPre[step] + curPre[pos2+1]
		}
	}
	res %= MOD
	if res < 0 {
		res += MOD
	}
	return res
}

// !start<step.
func (ss *RangeStepAddRangeSum) Update(start, step int, delta int) {
	delta %= MOD
	if delta < 0 {
		delta += MOD
	}
	if step >= ss.stepThreshold {
		for i := start; i < len(ss.nums); i += step {
			bid := ss.belong[i]
			ss.nums[i] += delta
			if ss.nums[i] >= MOD {
				ss.nums[i] -= MOD
			}
			ss.blockSum[bid] += delta
			if ss.blockSum[bid] >= MOD {
				ss.blockSum[bid] -= MOD
			}
		}
	} else {
		curPre, curSuf := ss.pre[step], ss.suf[step]
		for i := start; i+1 < len(curPre); i++ {
			curPre[i+1] += delta
			if curPre[i+1] >= MOD {
				curPre[i+1] -= MOD
			}
		}
		for i := 0; i <= start; i++ {
			curSuf[i] += delta
			if curSuf[i] >= MOD {
				curSuf[i] -= MOD
			}
		}
	}
}

func (ss *RangeStepAddRangeSum) _sum(start, end int) int {
	bid1, bid2 := ss.belong[start], ss.belong[end-1]
	res := 0
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			res += ss.nums[i]
		}
	} else {
		for i := start; i < ss.blockEnd[bid1]; i++ {
			res += ss.nums[i]
		}
		for bid := bid1 + 1; bid < bid2; bid++ {
			res += ss.blockSum[bid]
		}
		for i := ss.blockStart[bid2]; i < end; i++ {
			res += ss.nums[i]
		}
	}
	res %= MOD
	if res < 0 {
		res += MOD
	}
	return res
}
