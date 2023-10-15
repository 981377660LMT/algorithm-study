// RangeAddRangeModSum

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// https://loj.ac/p/6280
// 0 start end delta
// 1 start end mod
func RangeAddRangeModSum(nums []int, operations [][4]int) []int {
	block := UseBlock(len(nums), int(math.Sqrt(float64(len(nums)))+1))

	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockLazy := make([]int, blockCount)
	blockSum := make([]int, blockCount)

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
			start, end, delta := op[1], op[2], op[3]
			bid1, bid2 := belong[start], belong[end-1]
			if bid1 == bid2 {
				for i := start; i < end; i++ {
					nums[i] += delta
				}
				rebuild(bid1)
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					nums[i] += delta
				}
				rebuild(bid1)
				for i := bid1 + 1; i < bid2; i++ {
					blockLazy[i] += delta
				}
				for i := blockStart[bid2]; i < end; i++ {
					nums[i] += delta
				}
				rebuild(bid2)
			}
		} else {
			start, end, mod := op[1], op[2], op[3]
			bid1, bid2 := belong[start], belong[end-1]
			modSum := 0
			if bid1 == bid2 {
				for i := start; i < end; i++ {
					modSum += nums[i] + blockLazy[bid1]
					modSum %= mod
				}
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					modSum += nums[i] + blockLazy[bid1]
					modSum %= mod
				}
				for i := bid1 + 1; i < bid2; i++ {
					modSum += blockSum[i] + blockLazy[i]*(blockEnd[i]-blockStart[i])
					modSum %= mod
				}
				for i := blockStart[bid2]; i < end; i++ {
					modSum += nums[i] + blockLazy[bid2]
					modSum %= mod
				}
			}

			if modSum < 0 {
				modSum += mod
			}
			res = append(res, modSum)
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

	operations := make([][4]int, n)
	for i := range operations {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var start, end, delta int
			fmt.Fscan(in, &start, &end, &delta)
			start--
			operations[i] = [4]int{0, start, end, delta}
		} else {
			var start, end, mod int
			fmt.Fscan(in, &start, &end, &mod)
			start--
			operations[i] = [4]int{1, start, end, mod + 1}
		}
	}

	res := RangeAddRangeModSum(nums, operations)
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
