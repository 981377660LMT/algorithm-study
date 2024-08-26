// https://codeforces.com/blog/entry/82094#comment-688448
// Kinetic Tournament (KTT，势能线段树, 区间一次函数最大值)
//
// (x[i],y[i]) からなる列. a>=0 であるときに y[i] :- y[i] + ax[i] + b という作用ができる
// x には単調性は要らない. x,sum(a):T1, y,sum(b):T2, T1*T1<=T2.
// https://codeforces.com/blog/entry/82094#comment-688448
// https://atcoder.jp/contests/jsc2024-final/tasks/jsc2024_final_d
// https://www.luogu.com.cn/article/u0ouhel6
//
//
// 给定两个数组 ks[i],bs[i].
// 给定一次函数序列 fi = ks[i]*a + bs[i] + b. 支持以下操作：
// 1. Set(i,k,b)：将下标为i所在的直线参数设置为k,b.
// 2. Query(l,r)：查询[l,r)中fi的最大值.(区间最优函数)
// 3. Update(l,r,a,b)：对于所有l≤i<r，将下标i所在的横坐标加a，纵坐标加b.

package main

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"
)

func main() {
	seg := NewSegmentTreeBeatsKineticMax(1, func(i int32) (int, int) { return 10, 2 })
	fmt.Println(seg.Query(0, 1)) // 2
	seg.Update(0, 1, 2, 0)
	fmt.Println(seg.Query(0, 1)) // 22
	seg.Update(0, 1, 0, 1)
	fmt.Println(seg.Query(0, 1)) // 23
}

const INF int = 1e18

type SegmentTreeBeatsKineticMax struct {
	tree *segmentTreeBeats
}

func NewSegmentTreeBeatsKineticMax(n int32, f func(int32) (k, b int)) *SegmentTreeBeatsKineticMax {
	res := &SegmentTreeBeatsKineticMax{}
	res.tree = newSegmentTreeBeats(n, func(i int32) E {
		x, y := f(i)
		return E{x: x, y: y, nextChange: INF}
	})
	return res
}

func NewSegmentTreeBeatsKineticMaxFrom(ks, bs []int) *SegmentTreeBeatsKineticMax {
	return NewSegmentTreeBeatsKineticMax(int32(len(ks)), func(i int32) (int, int) {
		return ks[i], bs[i]
	})
}

func (seg *SegmentTreeBeatsKineticMax) Get(i int32) int {
	e := seg.tree.Get(i)
	return e.y
}

func (seg *SegmentTreeBeatsKineticMax) Query(l, r int32) int {
	e := seg.tree.Query(l, r)
	return e.y
}

func (seg *SegmentTreeBeatsKineticMax) QueryAll() int {
	e := seg.tree.QueryAll()
	return e.y
}

func (seg *SegmentTreeBeatsKineticMax) Set(i int32, k, b int) {
	seg.tree.Set(i, E{x: k, y: b, nextChange: INF})
}

func (seg *SegmentTreeBeatsKineticMax) Update(l, r int32, a int, b int) {
	seg.tree.Update(l, r, Id{a: a, b: b})
}

type E struct {
	fail             bool
	x, y, nextChange int
}

type Id struct{ a, b int }

func (tree *segmentTreeBeats) e() E   { return E{y: -INF, nextChange: INF} }
func (tree *segmentTreeBeats) id() Id { return Id{} }
func (tree *segmentTreeBeats) op(l, r E) E {
	if l.y < r.y {
		l, r = r, l
	}
	nextChange := min(l.nextChange, r.nextChange)
	if l.x < r.x {
		nextChange = min(nextChange, (l.y-r.y)/(r.x-l.x)+1)
	}
	l.fail = false
	l.nextChange = nextChange
	return l
}

func (tree *segmentTreeBeats) mapping(lazy Id, data E, size int) E {
	if data.nextChange <= lazy.a {
		data.fail = true
		return data
	}
	data.y += lazy.a*data.x + lazy.b
	data.nextChange -= lazy.a
	return data
}

func (tree *segmentTreeBeats) composition(p, c Id) Id {
	p.a += c.a
	p.b += c.b
	return p
}

// atcoder::lazy_segtree に1行書き足すだけの抽象化 Segment Tree Beats
// https://rsm9.hatenablog.com/entry/2021/02/01/220408
type segmentTreeBeats struct {
	n, log, offset int32
	data           []E
	lazy           []Id
}

func newSegmentTreeBeats(n int32, f func(int32) E) *segmentTreeBeats {
	res := &segmentTreeBeats{}
	res.n = n
	res.log = 1
	for 1<<res.log < n {
		res.log++
	}
	res.offset = 1 << res.log
	res.data = make([]E, res.offset<<1)
	for i := range res.data {
		res.data[i] = res.e()
	}
	res.lazy = make([]Id, res.offset)
	for i := range res.lazy {
		res.lazy[i] = res.id()
	}
	for i := int32(0); i < n; i++ {
		res.data[res.offset+i] = f(i)
	}
	for i := res.offset - 1; i >= 1; i-- {
		res.pushUp(i)
	}
	return res
}

func (tree *segmentTreeBeats) Get(index int32) E {
	index += tree.offset
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

func (tree *segmentTreeBeats) Set(index int32, e E) {
	index += tree.offset
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := int32(1); i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *segmentTreeBeats) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > tree.n {
		end = tree.n
	}
	if start >= end {
		return tree.e()
	}
	start += tree.offset
	end += tree.offset
	for i := tree.log; i >= 1; i-- {
		if ((start >> i) << i) != start {
			tree.pushDown(start >> i)
		}
		if ((end >> i) << i) != end {
			tree.pushDown((end - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for start < end {
		if start&1 != 0 {
			sml = tree.op(sml, tree.data[start])
			start++
		}
		if end&1 != 0 {
			end--
			smr = tree.op(tree.data[end], smr)
		}
		start >>= 1
		end >>= 1
	}
	return tree.op(sml, smr)
}

func (tree *segmentTreeBeats) QueryAll() E {
	return tree.data[1]
}

func (tree *segmentTreeBeats) Update(start, end int32, f Id) {
	if start < 0 {
		start = 0
	}
	if end > tree.n {
		end = tree.n
	}
	if start >= end {
		return
	}
	start += tree.offset
	end += tree.offset
	for i := tree.log; i >= 1; i-- {
		if ((start >> i) << i) != start {
			tree.pushDown(start >> i)
		}
		if ((end >> i) << i) != end {
			tree.pushDown((end - 1) >> i)
		}
	}
	l2, r2 := start, end
	for start < end {
		if start&1 != 0 {
			tree.propagate(start, f)
			start++
		}
		if end&1 != 0 {
			end--
			tree.propagate(end, f)
		}
		start >>= 1
		end >>= 1
	}
	start = l2
	end = r2
	for i := int32(1); i <= tree.log; i++ {
		if ((start >> i) << i) != start {
			tree.pushUp(start >> i)
		}
		if ((end >> i) << i) != end {
			tree.pushUp((end - 1) >> i)
		}
	}
}

func (tree *segmentTreeBeats) MinLeft(right int32, predicate func(data E) bool) int32 {
	if right == 0 {
		return 0
	}
	right += tree.offset
	for i := tree.log; i >= 1; i-- {
		tree.pushDown((right - 1) >> i)
	}
	res := tree.e()
	for {
		right--
		for right > 1 && right&1 != 0 {
			right >>= 1
		}
		if !predicate(tree.op(tree.data[right], res)) {
			for right < tree.offset {
				tree.pushDown(right)
				right = right<<1 | 1
				if predicate(tree.op(tree.data[right], res)) {
					res = tree.op(tree.data[right], res)
					right--
				}
			}
			return right + 1 - tree.offset
		}
		res = tree.op(tree.data[right], res)
		if (right & -right) == right {
			break
		}
	}
	return 0
}

func (tree *segmentTreeBeats) MaxRight(left int32, predicate func(data E) bool) int32 {
	if left == tree.n {
		return tree.n
	}
	left += tree.offset
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}
	res := tree.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.offset {
				tree.pushDown(left)
				left <<= 1
				if predicate(tree.op(res, tree.data[left])) {
					res = tree.op(res, tree.data[left])
					left++
				}
			}
			return left - tree.offset
		}
		res = tree.op(res, tree.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return tree.n
}

func (tree *segmentTreeBeats) GetAll() []E {
	for i := int32(1); i < tree.offset; i++ {
		tree.pushDown(i)
	}
	return tree.data[tree.offset : tree.offset+tree.n]
}

func (tree *segmentTreeBeats) pushUp(root int32) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}

func (tree *segmentTreeBeats) pushDown(root int32) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}

func (tree *segmentTreeBeats) propagate(root int32, lazy Id) {
	size := 1 << (tree.log - int32((bits.Len32(uint32(root)) - 1)))
	tree.data[root] = tree.mapping(lazy, tree.data[root], size)
	if root < tree.offset {
		tree.lazy[root] = tree.composition(lazy, tree.lazy[root])
		// !区别于普通线段树的地方.
		if tree.data[root].fail {
			tree.pushDown(root)
			tree.pushUp(root)
		}
	}
}

func (tree *segmentTreeBeats) String() string {
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

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}
