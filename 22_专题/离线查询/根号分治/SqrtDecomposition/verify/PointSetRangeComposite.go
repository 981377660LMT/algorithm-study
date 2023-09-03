// SqrtDecomposition
// https://nyaannyaan.github.io/library/data-structure/square-root-decomposition.hpp

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const MOD int = 998244353

func main() {
	// https://judge.yosupo.jp/problem/point_set_range_composite
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	nums := make([]E, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i].mul, &nums[i].add)
	}

	sqrt := NewSqrtDecomposition(nums, 1+int(math.Sqrt(float64(n))))
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 0 {
			var i, mul, add int
			fmt.Fscan(in, &i, &mul, &add)
			sqrt.Update(i, i+1, E{mul, add})
		} else {
			var start, end, x int
			fmt.Fscan(in, &start, &end, &x)
			affine := e()
			sqrt.Query(start, end, func(blockRes E) {
				affine = op(affine, blockRes)
			})
			fmt.Fprintln(out, (affine.mul*x+affine.add)%MOD)
		}
	}
}

type E = struct{ mul, add int }
type Id = E

// 仿射变换
func e() E        { return E{1, 0} }
func op(a, b E) E { return E{a.mul * b.mul % MOD, (a.add*b.mul + b.add) % MOD} }

type Block struct {
	// dont modify
	id, start, end int // block id => nums[start:end]
	nums           []E // block内的原序列

	// !to do
	sum E
	// lazyAdd Id
}

// 初始化块内数据(只会调用一次)
func (b *Block) Init() {
	b.Build()
}

// 重构
func (b *Block) Build() {
	b.sum = e()
	for i := range b.nums {
		b.sum = op(b.sum, b.nums[i])
	}
}
func (b *Block) UpdateAll(lazy Id) {}
func (b *Block) UpdatePart(start, end int, lazy Id) {
	for i := start; i < end; i++ {
		b.nums[i] = lazy
	}
}
func (b *Block) QueryAll() E { return b.sum }
func (b *Block) QueryPart(start, end int) E {
	res := e()
	for i := start; i < end; i++ {
		res = op(res, b.nums[i])
	}
	return res
}

//
//
//
// dont modify the template below
//
//
//

type SqrtDecomposition struct {
	n      int
	bs     int
	bls    []Block
	belong []int
}

// 指定维护的序列和分块大小初始化.
//
//	blockSize:分块大小,一般取根号n(300)
func NewSqrtDecomposition(nums []E, blockSize int) *SqrtDecomposition {
	nums = append(nums[:0:0], nums...)
	res := &SqrtDecomposition{
		n:      len(nums),
		bs:     blockSize,
		bls:    make([]Block, len(nums)/blockSize+1),
		belong: make([]int, len(nums)),
	}
	for i := range res.belong {
		res.belong[i] = i / blockSize
	}
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
//
//	0<=start<=end<=n
func (s *SqrtDecomposition) Update(start, end int, lazy Id) {
	if start >= end {
		return
	}
	id1, id2 := s.belong[start], s.belong[end-1]
	pos1, pos2 := start-s.bs*id1, end-s.bs*id2
	if id1 == id2 {
		s.bls[id1].UpdatePart(pos1, pos2, lazy)
		s.bls[id1].Build()
	} else {
		s.bls[id1].UpdatePart(pos1, s.bs, lazy)
		s.bls[id1].Build()
		for i := id1 + 1; i < id2; i++ {
			s.bls[i].UpdateAll(lazy)
		}
		s.bls[id2].UpdatePart(0, pos2, lazy)
		s.bls[id2].Build()
	}
}

// 查询左闭右开区间[start,end)的值.
//
//	0<=start<=end<=n
func (s *SqrtDecomposition) Query(start, end int, cb func(blockRes E)) {
	if start >= end {
		return
	}
	id1, id2 := s.belong[start], s.belong[end-1]
	pos1, pos2 := start-s.bs*id1, end-s.bs*id2
	if id1 == id2 {
		cb(s.bls[id1].QueryPart(pos1, pos2))
		return
	}
	cb(s.bls[id1].QueryPart(pos1, s.bs))
	for i := id1 + 1; i < id2; i++ {
		cb(s.bls[i].QueryAll())
	}
	cb(s.bls[id2].QueryPart(0, pos2))
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
