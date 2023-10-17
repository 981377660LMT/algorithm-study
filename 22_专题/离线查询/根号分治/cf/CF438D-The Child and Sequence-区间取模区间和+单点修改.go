// https://www.luogu.com.cn/problem/CF438D
// 区间取模区间和，单点修改

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 0 start end : 查询区间[start,end)的和.
// 1 start end mod: 将区间[start,end)所有数模mod.
// 2 pos value : 单点修改 nums[pos] = value.
// n,q<=1e5, nums[i]<=1e9, mod<=1e9
//
// 取模的结论： x mod p < x/2 (p<x)，所以取模也是最多logx次就不变了.
// !即：如果一个数在取模后改变了，那么它必定缩小至少一半.
// !需要维护块内最大值, 而且如果区间最大值小于模数，那取模就没有意义了，直接跳过.
func RangeModRangeSum(nums []int, operations [][4]int) []int {
	block := UseBlock(len(nums), 40) // 块大小为40，跑得很快

	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockSum := make([]int, blockCount)
	blockMax := make([]int, blockCount) // <mod 则跳过整块取模

	// 初始化/更新零散块后重构整个块
	rebuild := func(bid int) {
		blockSum[bid] = 0
		blockMax[bid] = 0
		for i := blockStart[bid]; i < blockEnd[bid]; i++ {
			blockSum[bid] += nums[i]
			blockMax[bid] = max(blockMax[bid], nums[i])
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
		} else if kind == 1 {
			start, end, mod := op[1], op[2], op[3]
			bid1, bid2 := belong[start], belong[end-1]
			if bid1 == bid2 {
				for i := start; i < end; i++ {
					nums[i] %= mod
				}
				rebuild(bid1)
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					nums[i] %= mod
				}
				rebuild(bid1)
				for i := bid1 + 1; i < bid2; i++ {
					if blockMax[i] < mod {
						continue
					}
					for j := blockStart[i]; j < blockEnd[i]; j++ {
						nums[j] %= mod
					}
					rebuild(i)
				}
				for i := blockStart[bid2]; i < end; i++ {
					nums[i] %= mod
				}
				rebuild(bid2)
			}
		} else {
			pos, target := op[1], op[2]
			pre := nums[pos]
			if pre == target {
				continue
			}
			nums[pos] = target
			rebuild(belong[pos])
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
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			operations[i] = [4]int{0, start, end}
		} else if op == 2 {
			var start, end, mod int
			fmt.Fscan(in, &start, &end, &mod)
			start--
			operations[i] = [4]int{1, start, end, mod}
		} else {
			var pos, value int
			fmt.Fscan(in, &pos, &value)
			pos--
			operations[i] = [4]int{2, pos, value}
		}
	}

	res := RangeModRangeSum(nums, operations)
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
