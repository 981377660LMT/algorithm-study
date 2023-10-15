// RangeAddRangeSmaller

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// https://loj.ac/p/6278
// 0 start end delta
// 1 start end k
func RangeAddRangeSmaller(nums []int, operations [][4]int) []int {
	block := UseBlock(len(nums), int(math.Sqrt(float64(len(nums)))+1))

	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockLazy := make([]int, blockCount)
	blockSorted := make([][]int, blockCount)

	// 初始化/更新零散块后重构整个块
	rebuild := func(bid int) {
		blockSorted[bid] = make([]int, blockEnd[bid]-blockStart[bid])
		copy(blockSorted[bid], nums[blockStart[bid]:blockEnd[bid]])
		sort.Ints(blockSorted[bid])
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
			start, end, k := op[1], op[2], op[3]
			bid1, bid2 := belong[start], belong[end-1]
			smaller := 0
			if bid1 == bid2 {
				for i := start; i < end; i++ {
					if nums[i]+blockLazy[bid1] < k {
						smaller++
					}
				}
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					if nums[i]+blockLazy[bid1] < k {
						smaller++
					}
				}
				for i := bid1 + 1; i < bid2; i++ {
					smaller += bisectLeft(blockSorted[i], k-blockLazy[i])
				}
				for i := blockStart[bid2]; i < end; i++ {
					if nums[i]+blockLazy[bid2] < k {
						smaller++
					}
				}
			}

			res = append(res, smaller)
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
		var op, start, end, k int
		fmt.Fscan(in, &op, &start, &end, &k)
		start--
		if op == 0 {
			operations[i] = [4]int{op, start, end, k}
		} else {
			operations[i] = [4]int{op, start, end, k * k}
		}
	}

	res := RangeAddRangeSmaller(nums, operations)
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

func bisectLeft(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}
func bisectRight(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] <= target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}
