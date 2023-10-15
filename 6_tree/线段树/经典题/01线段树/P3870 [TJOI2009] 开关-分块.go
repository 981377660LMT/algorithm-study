// P3870 [TJOI2009] 开关
// https://www.luogu.com.cn/problem/P3870
// 01 线段树/ bitset / 分块 三种方法

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// 0 start end: 将区间 [start, end) 内的所有数取反.
// 1 start end：询问区间 [start, end) 内 1 的个数.
func RangeFlipRangeOnesCount(nums []int, operations [][3]int) []int {
	block := UseBlock(len(nums), int(math.Sqrt(float64(len(nums)))+1))
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockFlip := make([]int, blockCount)
	blockOnes := make([]int, blockCount)

	rebuild := func(bid int) {
		blockOnes[bid] = 0
		for i := blockStart[bid]; i < blockEnd[bid]; i++ {
			blockOnes[bid] += nums[i]
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
					// !修改时需要维护 blockOnes.
					blockOnes[bid1] -= nums[i] ^ blockFlip[bid1]
					nums[i] ^= 1
					blockOnes[bid1] += nums[i] ^ blockFlip[bid1]
				}
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					blockOnes[bid1] -= nums[i] ^ blockFlip[bid1]
					nums[i] ^= 1
					blockOnes[bid1] += nums[i] ^ blockFlip[bid1]
				}
				for bid := bid1 + 1; bid < bid2; bid++ {
					blockFlip[bid] ^= 1
					blockOnes[bid] = (blockEnd[bid] - blockStart[bid]) - blockOnes[bid]
				}
				for i := blockStart[bid2]; i < end; i++ {
					blockOnes[bid2] -= nums[i] ^ blockFlip[bid2]
					nums[i] ^= 1
					blockOnes[bid2] += nums[i] ^ blockFlip[bid2]
				}
			}
		} else {
			start, end := op[1], op[2]
			bid1, bid2 := belong[start], belong[end-1]
			ones := 0
			if bid1 == bid2 {
				for i := start; i < end; i++ {
					ones += nums[i] ^ blockFlip[bid1]
				}
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					ones += nums[i] ^ blockFlip[bid1]
				}
				for bid := bid1 + 1; bid < bid2; bid++ {
					ones += blockOnes[bid]
				}
				for i := blockStart[bid2]; i < end; i++ {
					ones += nums[i] ^ blockFlip[bid2]
				}
			}

			res = append(res, ones)
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

	operations := make([][3]int, q)
	for i := range operations {
		var op, start, end int
		fmt.Fscan(in, &op, &start, &end)
		start--
		operations[i] = [3]int{op, start, end}
	}

	nums := make([]int, n)
	res := RangeFlipRangeOnesCount(nums, operations)
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
