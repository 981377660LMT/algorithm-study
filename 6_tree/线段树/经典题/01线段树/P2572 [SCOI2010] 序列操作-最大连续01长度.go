// P2572 [SCOI2010] 序列操作 (01线段树的基础上加上区间最大连续01长度)
// https://www.luogu.com.cn/problem/P2572
//
// 给定一个01序列，要求支持如下操作：
// 区间赋值
// 区间取反
// 查询区间内1的个数
// 查询区间内最多有多少个连续的1
//
// !让两个懒标记互斥，区间覆盖优先级高于区间翻转

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	bits := make([]int8, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &bits[i])
	}

	seg := NewRangeAssignRangeFlipRangeMaxContinuousOnes(n, func(i int32) int8 { return bits[i] })

	fillZero := func(start, end int32) {
		seg.FillZero(start, end)
	}

	fillOne := func(start, end int32) {
		seg.FillOne(start, end)
	}

	flip := func(start, end int32) {
		seg.Flip(start, end)
	}

	onesCount := func(start, end int32) int32 {
		return seg.OnesCount(start, end)
	}

	maxContinuousOnes := func(start, end int32) int32 {
		return seg.MaxContinuousOnes(start, end)
	}

	for i := int32(0); i < q; i++ {
		var op, l, r int32
		fmt.Fscan(in, &op, &l, &r)
		r++

		switch op {
		case 0:
			fillZero(l, r)
		case 1:
			fillOne(l, r)
		case 2:
			flip(l, r)
		case 3:
			fmt.Fprintln(out, onesCount(l, r))
		case 4:
			fmt.Fprintln(out, maxContinuousOnes(l, r))
		default:
			panic("unknown operation")
		}
	}
}

const INF32 int32 = 1 << 30

// RangeAssignRangeFlipRangeMaxContinuousOnes

type RangeAssignRangeFlipRangeMaxContinuousOnes struct {
	seg *lazySegTree32
}

func NewRangeAssignRangeFlipRangeMaxContinuousOnes(n int32, f func(int32) int8) *RangeAssignRangeFlipRangeMaxContinuousOnes {
	seg := newLazySegTree32(n, func(i int32) E { return fromElement(f(i)) })
	return &RangeAssignRangeFlipRangeMaxContinuousOnes{seg: seg}
}

func (r *RangeAssignRangeFlipRangeMaxContinuousOnes) FillZero(start, end int32) {
	r.seg.Update(start, end, Id{bit: 0, flip: 0})
}

func (r *RangeAssignRangeFlipRangeMaxContinuousOnes) FillOne(start, end int32) {
	r.seg.Update(start, end, Id{bit: 1, flip: 0})
}

func (r *RangeAssignRangeFlipRangeMaxContinuousOnes) Flip(start, end int32) {
	r.seg.Update(start, end, Id{bit: -1, flip: 1})
}

func (r *RangeAssignRangeFlipRangeMaxContinuousOnes) OnesCount(start, end int32) int32 {
	return r.seg.Query(start, end).count1
}

func (r *RangeAssignRangeFlipRangeMaxContinuousOnes) MaxContinuousOnes(start, end int32) int32 {
	return r.seg.Query(start, end).max1
}

func fromElement(v int8) E {
	res := E{}
	if v == 0 {
		res.pre0 = 1
		res.suf0 = 1
		res.max0 = 1
		res.count0 = 1
	} else {
		res.pre1 = 1
		res.suf1 = 1
		res.max1 = 1
		res.count1 = 1
	}
	return res
}

type E = struct {
	// 前缀后缀最大连续0/1个数，整个区间最大连续0/1个数, 0/1个数
	pre0, suf0, max0, count0 int32
	pre1, suf1, max1, count1 int32
}

type Id = struct {
	bit  int8 // 覆盖标记，0/1/-1
	flip int8 // 翻转标记, 0/1
}

func (*lazySegTree32) e() E {
	return E{}
}

func (*lazySegTree32) id() Id {
	return Id{bit: -1, flip: 0}
}

func (*lazySegTree32) op(a, b E) E {
	res := E{}

	res.pre0 = a.pre0
	if a.count1 == 0 {
		res.pre0 += b.pre0
	}
	res.suf0 = b.suf0
	if b.count1 == 0 {
		res.suf0 += a.suf0
	}
	res.max0 = max32(max32(a.max0, b.max0), a.suf0+b.pre0)
	res.count0 = a.count0 + b.count0

	res.pre1 = a.pre1
	if a.count0 == 0 {
		res.pre1 += b.pre1
	}
	res.suf1 = b.suf1
	if b.count0 == 0 {
		res.suf1 += a.suf1
	}
	res.max1 = max32(max32(a.max1, b.max1), a.suf1+b.pre1)
	res.count1 = a.count1 + b.count1

	return res
}

func (*lazySegTree32) mapping(f Id, g E) E {
	if f.bit == -1 {
		if f.flip == 0 {
			return g
		} else {
			g.pre0, g.pre1 = g.pre1, g.pre0
			g.suf0, g.suf1 = g.suf1, g.suf0
			g.max0, g.max1 = g.max1, g.max0
			g.count0, g.count1 = g.count1, g.count0
			return g
		}
	} else {
		size := g.count0 + g.count1
		if f.bit == 0 {
			g.pre0, g.pre1 = size, 0
			g.suf0, g.suf1 = size, 0
			g.max0, g.max1 = size, 0
			g.count0, g.count1 = size, 0
		} else {
			g.pre0, g.pre1 = 0, size
			g.suf0, g.suf1 = 0, size
			g.max0, g.max1 = 0, size
			g.count0, g.count1 = 0, size
		}
		return g
	}
}

func (*lazySegTree32) composition(f, g Id) Id {
	// !让两个懒标记互斥，不能同时存在
	if f.bit != -1 {
		g.bit = f.bit
		g.flip = 0
	} else {
		g.flip ^= f.flip
		if g.bit != -1 {
			g.bit ^= g.flip
			g.flip = 0
		}
	}
	return g
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

// !template
type lazySegTree32 struct {
	n    int32
	size int32
	log  int32
	data []E
	lazy []Id
}

func newLazySegTree32(n int32, f func(int32) E) *lazySegTree32 {
	tree := &lazySegTree32{}
	tree.n = n
	tree.log = int32(bits.Len32(uint32(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := int32(0); i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *lazySegTree32) Query(left, right int32) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}
func (tree *lazySegTree32) QueryAll() E {
	return tree.data[1]
}
func (tree *lazySegTree32) GetAll() []E {
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *lazySegTree32) Update(left, right int32, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := int32(1); i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

func (tree *lazySegTree32) Get(index int32) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *lazySegTree32) Set(index int32, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := int32(1); i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *lazySegTree32) pushUp(root int32) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *lazySegTree32) pushDown(root int32) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *lazySegTree32) propagate(root int32, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *lazySegTree32) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := int32(0); i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
}
