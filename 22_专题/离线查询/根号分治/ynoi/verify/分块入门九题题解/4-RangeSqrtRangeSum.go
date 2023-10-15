// RangeSqrtRangeSum
// 区间开方，区间和

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// https://loj.ac/d?problemId=6281
// https://www.luogu.com.cn/problem/P4145
// P4145 上帝造题的七分钟 2 / 花神游历各国
// 0 start end : 区间开方 (int(math.Sqrt(float64(nums[i])))
// 1 start end : 区间求和
// !需要维护块内是否全1, 区间更新时`把已经不能再开方的块进行跳过`
func RangeSqrtRangeSum(nums []int, operations [][3]int) []int {
	block := UseBlock(len(nums), int(math.Sqrt(float64(len(nums)))+1))

	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockSum := make([]int, blockCount)
	blockAllOne := make([]bool, blockCount)

	// 初始化/更新零散块后重构整个块
	rebuild := func(bid int) {
		blockSum[bid] = 0
		for i := blockStart[bid]; i < blockEnd[bid]; i++ {
			blockSum[bid] += nums[i]
		}
	}
	for bid := 0; bid < blockCount; bid++ {
		rebuild(bid)
	}

	res := []int{}
	for _, op := range operations {
		kind := op[0]
		if kind == 0 {
			start, end := op[1], op[2]
			bid1, bid2 := belong[start], belong[end-1]
			if bid1 == bid2 {
				for i := start; i < end; i++ {
					blockSum[bid1] -= nums[i]
					nums[i] = int(math.Sqrt(float64(nums[i])))
					blockSum[bid1] += nums[i]
				}
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					blockSum[bid1] -= nums[i]
					nums[i] = int(math.Sqrt(float64(nums[i])))
					blockSum[bid1] += nums[i]
				}
				for i := bid1 + 1; i < bid2; i++ {
					if blockAllOne[i] {
						continue
					}
					// !像更新散块一样更新
					allOne := true
					for j := blockStart[i]; j < blockEnd[i]; j++ {
						blockSum[i] -= nums[j]
						sqrt := int(math.Sqrt(float64(nums[j])))
						nums[j] = sqrt
						blockSum[i] += sqrt
						if sqrt > 1 {
							allOne = false
						}
					}
					blockAllOne[i] = allOne
				}
				for i := blockStart[bid2]; i < end; i++ {
					blockSum[bid2] -= nums[i]
					nums[i] = int(math.Sqrt(float64(nums[i])))
					blockSum[bid2] += nums[i]
				}
			}
		} else {
			start, end := op[1], op[2]
			bid1, bid2 := belong[start], belong[end-1]
			sum := 0
			if bid1 == bid2 {
				for i := start; i < end; i++ {
					sum += nums[i]
				}
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					sum += nums[i]
				}
				for i := bid1 + 1; i < bid2; i++ {
					sum += blockSum[i]
				}
				for i := blockStart[bid2]; i < end; i++ {
					sum += nums[i]
				}
			}

			res = append(res, sum)
		}
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)
	operations := make([][3]int, q)
	for i := range operations {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var start, end int
			fmt.Fscan(in, &start, &end)
			if start > end {
				start, end = end, start
			}
			start--
			operations[i] = [3]int{0, start, end}
		} else {
			var start, end int
			fmt.Fscan(in, &start, &end)
			if start > end {
				start, end = end, start
			}
			start--
			operations[i] = [3]int{1, start, end}
		}
	}

	res := RangeSqrtRangeSum(nums, operations)
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
