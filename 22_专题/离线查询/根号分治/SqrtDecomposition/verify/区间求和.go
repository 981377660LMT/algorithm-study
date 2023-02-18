// SqrtDecomposition
// https://nyaannyaan.github.io/library/data-structure/square-root-decomposition.hpp

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	// https://www.luogu.com.cn/problem/P3372
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]E, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	sqrt := NewSqrtDecomposition(nums, int(math.Sqrt(float64(n))))
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var l, r, add int
			fmt.Fscan(in, &l, &r, &add)
			l--
			sqrt.Update(l, r, add)
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			fmt.Fprintln(out, sqrt.Query(l, r))
		}
	}
}

type E = int
type Id = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }

// 例子:区间求和
type Block struct {
	// dont modify
	id, start, end int // block id => nums[start:end]
	nums           []E // block内的原序列

	// !to do
	sum     E
	lazyAdd Id
}

// 初始化块内数据(只会调用一次)
func (b *Block) Init() {
	b.Build()
}

// 重构
func (b *Block) Build() {
	b.sum = 0
	for i := range b.nums {
		b.sum += b.nums[i]
	}
}
func (b *Block) QueryAll() E { return b.sum + b.lazyAdd*(b.end-b.start) }
func (b *Block) QueryPart(start, end int) E {
	res := 0
	for i := start; i < end; i++ {
		res += b.nums[i] + b.lazyAdd
	}
	return res
}
func (b *Block) UpdateAll(lazy Id) { b.lazyAdd += lazy }
func (b *Block) UpdatePart(start, end int, lazy Id) {
	for i := start; i < end; i++ {
		b.nums[i] += lazy
	}
	b.Build() // !注意重构
}

//
//
//
// dont modify the template below
//
//
//

type SqrtDecomposition struct {
	n   int
	bs  int
	bls []Block
}

// 指定维护的序列和分块大小初始化.
//  blockSize:分块大小,一般取根号n(300)
func NewSqrtDecomposition(nums []E, blockSize int) *SqrtDecomposition {
	nums = append(nums[:0:0], nums...)
	res := &SqrtDecomposition{n: len(nums), bs: blockSize, bls: make([]Block, len(nums)/blockSize+1)}
	for i := range res.bls {
		res.bls[i].id = i
		res.bls[i].start = i * blockSize
		res.bls[i].end = min((i+1)*blockSize, len(nums))
		res.bls[i].nums = nums[res.bls[i].start:res.bls[i].end]
		res.bls[i].Init()
	}
	return res
}

// 更新左闭右开区间[start,end)的值.
//  0<=start<=end<=n
func (s *SqrtDecomposition) Update(start, end int, lazy Id) {
	if start/s.bs == end/s.bs {
		s.bls[start/s.bs].UpdatePart(start%s.bs, end%s.bs, lazy)
	} else {
		s.bls[start/s.bs].UpdatePart(start%s.bs, s.bs, lazy)
		for i := start/s.bs + 1; i < end/s.bs; i++ {
			s.bls[i].UpdateAll(lazy)
		}
		s.bls[end/s.bs].UpdatePart(0, end%s.bs, lazy)
	}
}

// 查询左闭右开区间[start,end)的值.
//  0<=start<=end<=n
func (s *SqrtDecomposition) Query(start, end int) E {
	if start/s.bs == end/s.bs {
		return s.bls[start/s.bs].QueryPart(start%s.bs, end%s.bs)
	}
	res := e()
	res = op(res, s.bls[start/s.bs].QueryPart(start%s.bs, s.bs))
	for i := start/s.bs + 1; i < end/s.bs; i++ {
		res = op(res, s.bls[i].QueryAll())
	}
	res = op(res, s.bls[end/s.bs].QueryPart(0, end%s.bs))
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
