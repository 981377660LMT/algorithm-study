package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// Mass Change Queries
// https://www.luogu.com.cn/problem/CF911G
// 给出一个数列,有q个操作,每种操作是把区间[start,end)中等于x的数改成y.输出q步操作完的数列.
// !1<=x<=y<=100.
// !线段树里面维护每个节点对应的区间里面，每个数将会变成哪个数
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int8, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
		nums[i]--
	}
	seg := NewSegmentTreeDual(n)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var start, end int
		var x, y int8
		fmt.Fscan(in, &start, &end, &x, &y)
		x--
		y--
		start--

		unit := seg.id()
		unit[x] = y
		seg.Update(start, end, unit)
	}

	res := seg.GetAll()
	for i := 0; i < n; i++ {
		fmt.Fprint(out, res[i][nums[i]]+1, " ")
	}
}

const MAX int8 = 100

type Id = [MAX]int8

const COMMUTATIVE = false

func (*SegmentTreeDual) id() Id {
	res := Id{}
	for i := int8(0); i < MAX; i++ {
		res[i] = i
	}
	return res
}
func (*SegmentTreeDual) composition(f, g Id) Id {
	for i := int8(0); i < MAX; i++ {
		g[i] = f[g[i]]
	}
	return g
}

type SegmentTreeDual struct {
	n            int
	size, height int
	lazy         []Id
	unit         Id
}

func NewSegmentTreeDual(n int) *SegmentTreeDual {
	res := &SegmentTreeDual{}
	size := 1
	height := 0
	for size < n {
		size <<= 1
		height++
	}
	lazy := make([]Id, 2*size)
	unit := res.id()
	for i := 0; i < 2*size; i++ {
		lazy[i] = unit
	}
	res.n = n
	res.size = size
	res.height = height
	res.lazy = lazy
	res.unit = unit
	return res
}
func (seg *SegmentTreeDual) Get(index int) Id {
	index += seg.size
	for i := seg.height; i > 0; i-- {
		seg.propagate(index >> i)
	}
	return seg.lazy[index]
}
func (seg *SegmentTreeDual) GetAll() []Id {
	for i := 0; i < seg.size; i++ {
		seg.propagate(i)
	}
	res := make([]Id, seg.n)
	copy(res, seg.lazy[seg.size:seg.size+seg.n])
	return res
}
func (seg *SegmentTreeDual) Update(left, right int, value Id) {
	if left < 0 {
		left = 0
	}
	if right > seg.n {
		right = seg.n
	}
	if left >= right {
		return
	}
	left += seg.size
	right += seg.size
	if !COMMUTATIVE {
		for i := seg.height; i > 0; i-- {
			if (left>>i)<<i != left {
				seg.propagate(left >> i)
			}
			if (right>>i)<<i != right {
				seg.propagate((right - 1) >> i)
			}
		}
	}
	for left < right {
		if left&1 > 0 {
			seg.lazy[left] = seg.composition(value, seg.lazy[left])
			left++
		}
		if right&1 > 0 {
			right--
			seg.lazy[right] = seg.composition(value, seg.lazy[right])
		}
		left >>= 1
		right >>= 1
	}
}
func (seg *SegmentTreeDual) propagate(k int) {
	if seg.lazy[k] != seg.unit {
		seg.lazy[k<<1] = seg.composition(seg.lazy[k], seg.lazy[k<<1])
		seg.lazy[k<<1|1] = seg.composition(seg.lazy[k], seg.lazy[k<<1|1])
		seg.lazy[k] = seg.unit
	}
}
func (st *SegmentTreeDual) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := 0; i < st.n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(fmt.Sprint(st.Get(i)))
	}
	buf.WriteByte(']')
	return buf.String()
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
