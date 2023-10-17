package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://www.luogu.com.cn/problem/CF438D
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

	R := NewRangeModRangeSum(nums, 30)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			fmt.Fprintln(out, R.Query(start, end))
		} else if op == 2 {
			var start, end, mod int
			fmt.Fscan(in, &start, &end, &mod)
			start--
			R.Update(start, end, mod)
		} else {
			var pos, value int
			fmt.Fscan(in, &pos, &value)
			pos--
			R.Set(pos, value)
		}
	}

}

type RangeModRangeSum struct {
	nums       []int
	belong     []int
	blockStart []int
	blockEnd   []int
	blockSum   []int
	blockMax   []int
}

// blockSize = 30 (1e5时较快)
func NewRangeModRangeSum(nums []int, blockSize int) *RangeModRangeSum {
	nums = append(nums[:0:0], nums...)
	block := UseBlock(len(nums), blockSize)
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockSum := make([]int, blockCount)
	blockMax := make([]int, blockCount) // <mod 则跳过整块取模
	res := &RangeModRangeSum{
		nums:       nums,
		belong:     belong,
		blockStart: blockStart,
		blockEnd:   blockEnd,
		blockSum:   blockSum,
		blockMax:   blockMax,
	}
	for bid := 0; bid < blockCount; bid++ {
		res.rebuild(bid)
	}
	return res
}

// 0 start end : 查询区间[start,end)的和.
func (r *RangeModRangeSum) Query(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(r.nums) {
		end = len(r.nums)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := r.belong[start], r.belong[end-1]
	sum := 0
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			sum += r.nums[i]
		}
	} else {
		for i := start; i < r.blockEnd[bid1]; i++ {
			sum += r.nums[i]
		}
		for i := bid1 + 1; i < bid2; i++ {
			sum += r.blockSum[i]
		}
		for i := r.blockStart[bid2]; i < end; i++ {
			sum += r.nums[i]
		}
	}
	return sum
}

// 将区间[start,end)所有数模mod.
func (r *RangeModRangeSum) Update(start, end, mod int) {
	if start < 0 {
		start = 0
	}
	if end > len(r.nums) {
		end = len(r.nums)
	}
	if start >= end {
		return
	}
	bid1, bid2 := r.belong[start], r.belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			r.nums[i] %= mod
		}
		r.rebuild(bid1)
	} else {
		for i := start; i < r.blockEnd[bid1]; i++ {
			r.nums[i] %= mod
		}
		r.rebuild(bid1)
		for i := bid1 + 1; i < bid2; i++ {
			if r.blockMax[i] < mod {
				continue
			}
			for j := r.blockStart[i]; j < r.blockEnd[i]; j++ {
				r.nums[j] %= mod
			}
			r.rebuild(i)
		}
		for i := r.blockStart[bid2]; i < end; i++ {
			r.nums[i] %= mod
		}
		r.rebuild(bid2)
	}
}

// 单点修改 nums[pos] = value.
func (r *RangeModRangeSum) Set(pos, value int) {
	if pos < 0 || pos >= len(r.nums) {
		return
	}
	pre := r.nums[pos]
	if pre == value {
		return
	}
	r.nums[pos] = value
	r.rebuild(r.belong[pos])
}

func (r *RangeModRangeSum) rebuild(bid int) {
	r.blockSum[bid] = 0
	r.blockMax[bid] = 0
	for i := r.blockStart[bid]; i < r.blockEnd[bid]; i++ {
		r.blockSum[bid] += r.nums[i]
		r.blockMax[bid] = max(r.blockMax[bid], r.nums[i])
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
