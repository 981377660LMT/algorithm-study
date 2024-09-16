// SegmentTreeAddChminChmaxWithCost

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// demo()
	abc365_f()
}

func demo() {
	S := NewSegmentTreeAddChminChmaxWithCost(10, func(i int32) E {
		return NewE()
	})

	S.Set(0, NewEWith(0, 0, 0, 8, 5))
	S.Set(1, NewEWith(0, 0, 0, 4, 5))

	e := S.Query(0, 2)
	fmt.Println(e.Eval(4))
	fmt.Println(e.EvalCost(4))
}

// F - Takahashi on Grid (网格图最短路)
// https://atcoder.jp/contests/abc365/tasks/abc365_f
// 给定一个无穷大的的网格，每个单元格为空或墙壁.
// 给定n个[L,R]，表示第i列的[L,R]之间的单元格是空地，其他是墙壁
//
// 现有q次询问，每次询问给定起点和终点
// 从起点出发，规定每次只能走到相邻的空单元格上，问最少走多少步能走到终点。
func abc365_f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int32
	fmt.Fscan(in, &N)

	L, R := make([]int, N), make([]int, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &L[i], &R[i])
	}

	// 每个元素为列的信息.
	seg := NewSegmentTreeAddChminChmaxWithCost(N, func(i int32) E {
		// x -> clamp(x, L[i], R[i])
		return NewEWith(0, 0, 0, L[i], R[i])
	})

	var Q int
	fmt.Fscan(in, &Q)
	for i := 0; i < Q; i++ {
		var x1, y1, x2, y2 int
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		x1, x2 = x1-1, x2-1
		if x1 > x2 {
			x1, x2 = x2, x1
			y1, y2 = y2, y1
		}

		f := seg.Query(int32(x1), int32(x2))

		// ma、mi: 起点为y1时, 因chmax、chmin导致的需要的移动量.
		add, ma, mi := f.EvalCost(y1)
		res := ma + mi
		y1 += add + ma - mi // 最终y1的位置

		res += abs(y1 - y2)
		res += abs(x1 - x2)
		fmt.Fprintln(out, res)
	}
}

const INF int = 1e18
const INF32 int32 = 1 << 30

// https://atcoder.jp/contests/abc365/editorial/10582
// MonoidAddChminChmaxWithCost
// x变化时, add导致的变化量为a, chmax导致的变化量为b+max(x1-x,0), chmin导致的变化量为c+max(x-x2,0)
type E struct{ a, b, c, x1, x2 int }

func NewE() E {
	return E{
		x1: -INF, x2: INF,
	}
}

func NewEWith(a, b, c, x1, x2 int) E {
	return E{a: a, b: b, c: c, x1: x1, x2: x2}
}

// 起点为y时, add、chmax、chmin导致的移动量.
func (e E) EvalCost(y int) (add int, chmax int, chmin int) {
	return e.a, e.b + max(0, e.x1-y), e.c + max(0, y-e.x2)
}

// 起点为y，返回操作后的位置.
func (e E) Eval(y int) int {
	add, ma, mi := e.a, e.b+max(0, e.x1-y), e.c+max(0, y-e.x2)
	return y + add + ma - mi
}

func (*SegmentTreeAddChminChmaxWithCost) e() E { return NewE() }
func (*SegmentTreeAddChminChmaxWithCost) op(a, b E) E {
	res := NewE()
	x1, x2 := a.x1, a.x2
	y1 := b.x1 - a.a - a.b + a.c
	y2 := b.x2 - a.a - a.b + a.c
	res.a = a.a + b.a
	if y1 < x1 {
		res.b = a.b + b.b
		res.x1 = x1
	} else if y1 < x2 {
		res.b = a.b + b.b
		res.x1 = y1
	} else {
		res.b = a.b + b.b + y1 - x2
		res.x1 = x2
	}
	if y2 < x1 {
		res.c = a.c + b.c + x1 - y2
		res.x2 = x1
	} else if y2 < x2 {
		res.c = a.c + b.c
		res.x2 = y2
	} else {
		res.c = a.c + b.c
		res.x2 = x2
	}
	return res
}

// SegmentTreeClamp.
type SegmentTreeAddChminChmaxWithCost struct {
	n, size int32
	seg     []E
}

func NewSegmentTreeAddChminChmaxWithCost(n int32, f func(int32) E) *SegmentTreeAddChminChmaxWithCost {
	res := &SegmentTreeAddChminChmaxWithCost{}
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func NewSegmentTreeAddChminChmaxWithCostFrom(leaves []E) *SegmentTreeAddChminChmaxWithCost {
	res := &SegmentTreeAddChminChmaxWithCost{}
	n := int32(len(leaves))
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTreeAddChminChmaxWithCost) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeAddChminChmaxWithCost) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeAddChminChmaxWithCost) Update(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTreeAddChminChmaxWithCost) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTreeAddChminChmaxWithCost) QueryAll() E { return st.seg[1] }
func (st *SegmentTreeAddChminChmaxWithCost) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreeAddChminChmaxWithCost) MaxRight(left int32, predicate func(E) bool) int32 {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if tmp := st.op(res, st.seg[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreeAddChminChmaxWithCost) MinLeft(right int32, predicate func(E) bool) int32 {
	if right == 0 {
		return 0
	}
	right += st.size
	res := st.e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(st.op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := st.op(st.seg[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
