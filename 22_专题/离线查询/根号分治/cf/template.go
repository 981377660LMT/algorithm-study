// RangeReplaceRangeSum

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// 0 start end delta
// 1 start end k
func RangeReplaceRangeSum(nums []int, operations [][3]int) []int {
	block := UseBlock(len(nums), int(math.Sqrt(float64(len(nums)))+1))

	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount

	// 初始化/更新零散块后重构整个块
	rebuild := func(bid int) {
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
			} else {
			}
		} else {
			start, end := op[1], op[2]
			bid1, bid2 := belong[start], belong[end-1]
			if bid1 == bid2 {
			} else {

			}
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

	operations := make([][3]int, q)
	for i := range operations {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			operations[i] = [3]int{0, start, end}
		} else {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			operations[i] = [3]int{1, start, end}
		}
	}

	res := RangeReplaceRangeSum(nums, operations)
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
