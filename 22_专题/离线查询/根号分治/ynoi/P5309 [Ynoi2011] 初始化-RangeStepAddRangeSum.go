// P5309 [Ynoi2011] 初始化
// https://www.luogu.com.cn/problem/P5309

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

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
func RangeStepAddRangeSum(nums []int, operations [][4]int) []int {
	sqrt := int(math.Sqrt(float64(len(nums))) + 1)
	B := UseBlock(len(nums), sqrt)
	belong, blockStart, blockEnd, blockCount := B.belong, B.blockStart, B.blockEnd, B.blockCount
	blockSum := make([]int, blockCount)
	for i := range blockSum {
		for j := blockStart[i]; j < blockEnd[i]; j++ {
			blockSum[i] = (blockSum[i] + nums[j]) % MOD
		}
	}

	pre, suf := make([][]int, sqrt+1), make([][]int, sqrt+1)
	for i := range pre {
		pre[i] = make([]int, sqrt+1)
		suf[i] = make([]int, sqrt+1)
	}

	_sum := func(start, end int) int {
		bid1, bid2 := belong[start], belong[end-1]
		res := 0
		if bid1 == bid2 {
			for i := start; i < end; i++ {
				// res = (res + nums[i]) % MOD
				res += nums[i]
			}
		} else {
			for i := start; i < blockEnd[bid1]; i++ {
				// res = (res + nums[i]) % MOD
				res += nums[i]
			}
			for bid := bid1 + 1; bid < bid2; bid++ {
				// res = (res + blockSum[bid]) % MOD
				res += blockSum[bid]
			}
			for i := blockStart[bid2]; i < end; i++ {
				// res = (res + nums[i]) % MOD
				res += nums[i]
			}
		}
		return res % MOD
	}

	update := func(start, step, delta int) {
		if step >= sqrt {
			for i := start; i < len(nums); i += step {
				bid := belong[i]
				nums[i] += delta
				if nums[i] >= MOD {
					nums[i] %= MOD
				}
				blockSum[bid] += delta
				if blockSum[bid] >= MOD {
					blockSum[bid] %= MOD
				}
			}
		} else {
			curPre, curSuf := pre[step], suf[step]
			for i := start; i+1 < len(curPre); i++ {
				curPre[i+1] += delta
				if curPre[i+1] >= MOD {
					curPre[i+1] %= MOD
				}
			}
			for i := 0; i <= start; i++ {
				curSuf[i] += delta
				if curSuf[i] >= MOD {
					curSuf[i] %= MOD
				}
			}
		}
	}

	query := func(start, end int) int {
		res := _sum(start, end)
		// 加上每个step对应的和.
		for step := 1; step < sqrt; step++ {
			id1, id2 := start/step, (end-1)/step
			pos1, pos2 := start%step, (end-1)%step
			curPre, curSuf := pre[step], suf[step]
			if id1 == id2 {
				// res = (res + curPre[pos2+1] - curPre[pos1] + MOD) % MOD
				res += curPre[pos2+1] - curPre[pos1]

			} else {
				// res = (res + curSuf[pos1] + (id2-id1-1)*curPre[step] + curPre[pos2+1]) % MOD
				res += curSuf[pos1] + (id2-id1-1)*curPre[step] + curPre[pos2+1]
			}
		}

		res = res % MOD
		if res < 0 {
			res += MOD
		}
		return res
	}

	res := []int{}
	for _, op := range operations {
		kind := op[0]
		if kind == 0 {
			start, step, delta := op[1], op[2], op[3]
			update(start, step, delta)
		} else {
			start, end := op[1], op[2]
			res = append(res, query(start, end))
		}
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	operations := make([][4]int, q)
	for i := range operations {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var step, start, delta int
			fmt.Fscan(in, &step, &start, &delta)
			start--
			operations[i] = [4]int{0, start, step, delta}
		} else {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			operations[i] = [4]int{1, start, end, 0}
		}
	}

	res := RangeStepAddRangeSum(nums, operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
