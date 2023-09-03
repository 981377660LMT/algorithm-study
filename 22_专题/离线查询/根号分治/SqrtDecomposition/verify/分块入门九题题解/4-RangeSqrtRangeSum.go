package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// 给出一个长为 n 的数列 a_1 \ldots a_n，以及 n 个操作，操作涉及区间开方，区间求和。
// 若 \mathrm{opt} = 0，表示将位于 [l, r] 的之间的数字都开方后向下取整。
// 若 \mathrm{opt} = 1，表示询问位于 [l, r] 的所有数字的和。
// RangeSqrtRangeSum
// 区间开方，区间求和
//
// !需要维护块内是否全1, 区间更新时把已经不能再开方的块进行跳过
func main() {
	// https://loj.ac/d?problemId=6281
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]E, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	sqrt := NewSqrtDecomposition(nums, 1+int(math.Sqrt(float64(n))))
	for i := 0; i < n; i++ {
		var op, l, r, c int
		fmt.Fscan(in, &op, &l, &r, &c)
		l--

		if op == 0 {
			sqrt.Update(l, r, c)
		} else {
			res := 0
			sqrt.Query(l, r, func(cur E) { res += cur })
			fmt.Fprintln(out, res)
		}
	}
}

type E = int  // 需要查询的值的类型
type Id = int // 懒标记的类型

// 例子:区间求和
type Block struct {
	// dont modify
	id, start, end int // 0 <= id < bs, 0 <= start <= end <= n
	nums           []E // block内的原序列

	// !to do
	sum    E
	allOne bool
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
func (b *Block) QueryPart(start, end int) E {
	res := 0
	for i := start; i < end; i++ {
		res += b.nums[i]
	}
	return res
}
func (b *Block) UpdatePart(start, end int, lazy Id) {
	for i := start; i < end; i++ {
		b.sum -= b.nums[i]
		b.nums[i] = int(math.Sqrt(float64(b.nums[i])))
		b.sum += b.nums[i]
	}
}
func (b *Block) QueryAll() E {
	return b.sum
}
func (b *Block) UpdateAll(lazy Id) {
	if b.allOne {
		return
	}
	b.allOne = true
	b.sum = 0
	for i := range b.nums {
		sqrt_ := int(math.Sqrt(float64(b.nums[i])))
		b.nums[i] = sqrt_
		b.sum += sqrt_
		if sqrt_ > 1 {
			b.allOne = false
		}
	}
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
