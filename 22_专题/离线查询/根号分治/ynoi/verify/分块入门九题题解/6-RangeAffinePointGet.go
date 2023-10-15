// 给出一个长为 n 的数列，以及 n 个操作，操作涉及区间乘法，区间加法，单点询问。
// RangeAffinePointGet
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// 第一行输入一个数字 n。
// 第二行输入 n 个数字，第 i 个数字为 a_i，以空格隔开。
// 接下来输入 n 行询问，每行输入四个数字 \mathrm{opt}、l、r、c，以空格隔开。
// 若 \mathrm{opt} = 0，表示将位于 [l, r] 的之间的数字都加 c。
// 若 \mathrm{opt} = 1，表示将位于 [l, r] 的之间的数字都乘 c。
// 若 \mathrm{opt} = 2，表示询问 a_r 的值 \mathop{\mathrm{mod}} 10007（l 和 c 忽略）。
// !区间仿射变换，单点查询
// 注意修改非完整块的时候，要先标记下放再修改
const MOD int = 10007

func main() {
	// https://loj.ac/p/6283
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
			sqrt.Update(l, r, Id{mul: 1, add: c})
		} else if op == 1 {
			sqrt.Update(l, r, Id{mul: c, add: 0})
		} else {
			res := 0
			sqrt.Query(r-1, r, func(cur E) { res += cur })
			res %= MOD
			if res < 0 {
				res += MOD
			}
			fmt.Fprintln(out, res)
		}
	}
}

type E = int                     // 需要查询的值的类型
type Id = struct{ mul, add int } // 懒标记的类型

type Block struct {
	id, start, end int // 0 <= id < bs, 0 <= start <= end <= n
	nums           []E // block内的原序列

	sum  E
	lazy Id
}

func (b *Block) Init() {
	b.lazy = Id{mul: 1, add: 0}
	for i := range b.nums {
		b.nums[i] %= MOD
		b.sum = (b.sum + b.nums[i]) % MOD
	}
}

// 重构
func (b *Block) ReBuild() {
	b.sum = 0
	add, mul := b.lazy.add, b.lazy.mul
	for i := range b.nums {
		b.nums[i] = b.nums[i]*mul + add
		b.nums[i] %= MOD
		b.sum += b.nums[i]
		b.sum %= MOD
	}
	b.lazy = Id{mul: 1, add: 0}
}

func (b *Block) QueryPart(start, end int) E {
	res := 0
	add, mul := b.lazy.add, b.lazy.mul
	for i := start; i < end; i++ {
		res += b.nums[i]*mul + add
		res %= MOD
	}
	return res
}

func (b *Block) UpdatePart(start, end int, lazy Id) {
	b.ReBuild()
	mul, add := lazy.mul, lazy.add
	for i := start; i < end; i++ {
		b.nums[i] = b.nums[i]*mul + add
		b.nums[i] %= MOD
	}
}

func (b *Block) QueryAll() E {
	panic("not implemented")
}

func (b *Block) UpdateAll(lazy Id) {
	b.lazy.mul = b.lazy.mul * lazy.mul % MOD
	b.lazy.add = (b.lazy.add*lazy.mul + lazy.add) % MOD
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
	} else {
		s.bls[id1].UpdatePart(pos1, s.bs, lazy)
		for i := id1 + 1; i < id2; i++ {
			s.bls[i].UpdateAll(lazy)
		}
		s.bls[id2].UpdatePart(0, pos2, lazy)
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
