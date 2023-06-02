package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/staticrmq

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	rmq := NewLinearRMQ(n, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	for i := 0; i < q; i++ {
		var start, end int
		fmt.Fscan(in, &start, &end)
		fmt.Fprintln(out, nums[rmq.Query(start, end)])
	}
}

type LinearRMQ struct {
	small []int
	large [][]int
	less  func(i, j int) bool
}

// n: 序列长度.
// less: 入参为两个索引,返回值表示索引i处的值是否小于索引j处的值.
//  消除了泛型.
func NewLinearRMQ(n int, less func(i, j int) bool) *LinearRMQ {
	res := &LinearRMQ{less: less}
	stack := make([]int, 0, 32)
	small := make([]int, 0, n)
	var large [][]int
	large = append(large, make([]int, 0, n>>5))
	for i := 0; i < n; i++ {
		for len(stack) > 0 && !less(stack[len(stack)-1], i) {
			stack = stack[:len(stack)-1]
		}
		tmp := 0
		if len(stack) > 0 {
			tmp = small[stack[len(stack)-1]]
		}
		small = append(small, tmp|(1<<(i&31)))
		stack = append(stack, i)
		if (i+1)&31 == 0 {
			large[0] = append(large[0], stack[0])
			stack = stack[:0]
		}
	}

	for i := 1; (i << 1) <= n>>5; i <<= 1 {
		csz := n>>5 + 1 - (i << 1)
		v := make([]int, 0, csz)
		for k := 0; k < csz; k++ {
			back := large[len(large)-1]
			v = append(v, res._getMin(back[k], back[k+i]))
		}
		large = append(large, v)
	}

	res.small = small
	res.large = large
	return res
}

// 查询区间`[start, end)`中的最小值的索引.
func (rmq *LinearRMQ) Query(start, end int) (minIndex int) {
	if start >= end {
		panic(fmt.Sprintf("start(%d) should be less than end(%d)", start, end))
	}
	end--
	left := start>>5 + 1
	right := end >> 5
	if left < right {
		msb := bits.Len32(uint32(right-left)) - 1
		cache := rmq.large[msb]
		i := (left-1)<<5 + bits.TrailingZeros32(uint32(rmq.small[left<<5-1]&(^0<<start)&31))
		cand1 := rmq._getMin(i, cache[left])
		j := left<<5 + bits.TrailingZeros32(uint32(rmq.small[end]))
		cand2 := rmq._getMin(cache[right-(1<<msb)], j)
		return rmq._getMin(cand1, cand2)
	}
	if left == right {
		i := (left-1)<<5 + bits.TrailingZeros32(uint32(rmq.small[left<<5-1]&(^0<<start)&31))
		j := left<<5 + bits.TrailingZeros32(uint32(rmq.small[end]))
		return rmq._getMin(i, j)
	}
	return right<<5 + bits.TrailingZeros32(uint32(rmq.small[end]&(^0<<start)&31))
}

func (rmq *LinearRMQ) _getMin(i, j int) int {
	if rmq.less(i, j) {
		return i
	}
	return j
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
