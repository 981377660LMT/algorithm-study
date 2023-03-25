// 可以区间排序的线段树/可排序线段树
// 单点修改+区间查询+区间排序

// API:
//  Set
//  SortInc/SortDec/SortRange
//  Query/QueryAll
//  GetEntries

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"runtime/debug"
	"strings"
)

// !因为都用的指针, 禁用会gc快很多
func init() {
	debug.SetGCPercent(-1)
}

func ABC237_G() {
	// https://atcoder.jp/contests/abc237/tasks/abc237_g
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q, x int
	fmt.Fscan(in, &n, &q, &x)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
		nums[i]--
	}
	seg := NewSortableSegmentTree(n, nums, nums)
	for i := 0; i < q; i++ {
		var op, l, r int
		fmt.Fscan(in, &op, &l, &r)
		l--
		if op == 1 {
			seg.SortInc(l, r)
		} else {
			seg.SortDec(l, r)
		}
	}

	keys, _ := seg.GetEntries()
	for i := 0; i < n; i++ {
		if keys[i]+1 == x {
			fmt.Fprintln(out, i+1)
			return
		}
	}
}

func main() {
	ABC237_G()
}

// 单点修改区间查询区间排序
type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }

type SortableSegmentTree struct {
	n, maxKey int
	unit      E
	rev       []bool
	unitNode  SNode
	ss        *fs
	seg       *SegmentTree
	root      []*SNode
}

type SNode struct {
	size    int
	x, revX E
	l, r    *SNode
}

func NewSortableSegmentTree(maxKey int, keys []int, values []E) *SortableSegmentTree {
	n := len(keys)
	res := &SortableSegmentTree{
		n:        n,
		maxKey:   maxKey,
		rev:      make([]bool, n),
		unitNode: SNode{size: 1, x: e(), revX: e()},
		ss:       nfs(n),
		unit:     e(),
		seg:      NewSegmentTree(values),
		root:     make([]*SNode, 0, n),
	}
	for i := 0; i < n; i++ {
		res.ss.Insert(i)
		newNode := res.unitNode // copy
		res.root = append(res.root, &newNode)
		res.setRec(res.root[i], 0, res.maxKey, keys[i], values[i])
	}
	return res
}

func (sg *SortableSegmentTree) Set(i, key int, value E) {
	sg.splitAt(i)
	sg.splitAt(i + 1)
	sg.rev[i] = false
	newNode := sg.unitNode
	sg.root[i] = &newNode
	sg.setRec(sg.root[i], 0, sg.maxKey, key, value)
	sg.seg.Set(i, value)
}

func (sg *SortableSegmentTree) Query(start, end int) E {
	sg.splitAt(start)
	sg.splitAt(end)
	return sg.seg.Query(start, end)
}

func (sg *SortableSegmentTree) QueryAll() E { return sg.seg.QueryAll() }

func (sg *SortableSegmentTree) SortRange(start, end int, reverse bool) {
	if reverse {
		sg.SortDec(start, end)
	} else {
		sg.SortInc(start, end)
	}
}

func (sg *SortableSegmentTree) SortInc(start, end int) {
	sg.splitAt(start)
	sg.splitAt(end)
	for {
		c := sg.root[start]
		i := sg.ss.Next(start + 1)
		if i == end {
			break
		}
		sg.root[start] = sg.merge(c, sg.root[i])
		sg.ss.Erase(i)
		sg.seg.Set(start, sg.unit)
	}
	sg.rev[start] = false
	sg.seg.Set(start, sg.root[start].x)
}

func (sg *SortableSegmentTree) SortDec(start, end int) {
	sg.SortInc(start, end)
	sg.rev[start] = true
	sg.seg.Set(start, sg.root[start].revX)
}

func (sg *SortableSegmentTree) GetEntries() (keys []int, values []E) {
	keys, values = make([]int, 0, sg.n), make([]E, 0, sg.n)
	var dfs func(np *SNode, l, r int, rev bool)
	dfs = func(np *SNode, l, r int, rev bool) {
		if np == nil {
			return
		}
		if l+1 == r {
			keys = append(keys, l)
			values = append(values, np.x)
			return
		}
		m := (l + r) / 2
		if !rev {
			dfs(np.l, l, m, rev)
			dfs(np.r, m, r, rev)
		} else {
			dfs(np.r, m, r, rev)
			dfs(np.l, l, m, rev)
		}
	}
	for i := 0; i < sg.n; i++ {
		if sg.ss.Has(i) {
			dfs(sg.root[i], 0, sg.maxKey, sg.rev[i])
		}
	}
	return
}

func (sg *SortableSegmentTree) String() string {
	keys, values := sg.GetEntries()
	sb := []string{}
	for i := 0; i < sg.n; i++ {
		sb = append(sb, fmt.Sprintf("%d: %d", keys[i], values[i]))
	}
	return strings.Join(sb, ", ")
}

func (sg *SortableSegmentTree) splitAt(x int) {
	if x == sg.n || sg.ss.Has(x) {
		return
	}
	a := sg.ss.Prev(x)
	b := sg.ss.Next(a + 1)
	sg.ss.Insert(x)
	if !sg.rev[a] {
		nl, nr := sg.split(sg.root[a], x-a)
		sg.root[a], sg.root[x] = nl, nr
		sg.rev[a], sg.rev[x] = false, false
		sg.seg.Set(a, sg.root[a].x)
		sg.seg.Set(x, sg.root[x].x)
	} else {
		nl, nr := sg.split(sg.root[a], b-x)
		sg.root[a], sg.root[x] = nr, nl
		sg.rev[a], sg.rev[x] = true, true
		sg.seg.Set(a, sg.root[a].revX)
		sg.seg.Set(x, sg.root[x].revX)
	}
}

func (sg *SortableSegmentTree) split(node *SNode, k int) (*SNode, *SNode) {
	if k == 0 {
		return nil, node
	}
	if k == node.size {
		return node, nil
	}
	s := 0
	if node.l != nil {
		s = node.l.size
	}
	newNode := sg.unitNode
	b := &newNode
	if k <= s {
		nl, nr := sg.split(node.l, k)
		b.l, b.r, node.l, node.r = nr, node.r, nl, nil
	} else {
		nl, nr := sg.split(node.r, k-s)
		node.r, b.l, b.r = nl, nil, nr
	}
	sg.update(node)
	sg.update(b)
	return node, b
}

func (sg *SortableSegmentTree) merge(a, b *SNode) *SNode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	a.l = sg.merge(a.l, b.l)
	a.r = sg.merge(a.r, b.r)
	sg.update(a)
	return a
}

func (sg *SortableSegmentTree) update(node *SNode) {
	if node.l == nil && node.r == nil {
		return
	}
	if node.l == nil {
		node.x = node.r.x
		node.revX = node.r.revX
		node.size = node.r.size
		return
	}
	if node.r == nil {
		node.x = node.l.x
		node.revX = node.l.revX
		node.size = node.l.size
		return
	}
	node.x = op(node.l.x, node.r.x)
	node.revX = op(node.r.revX, node.l.revX)
	node.size = node.l.size + node.r.size
}

func (sg *SortableSegmentTree) setRec(node *SNode, l, r, k int, x E) {
	if r == l+1 {
		node.x = x
		node.revX = x
		return
	}
	m := (l + r) / 2
	if k < m {
		if node.l == nil {
			newNode := sg.unitNode
			node.l = &newNode
		}
		sg.setRec(node.l, l, m, k, x)
	} else {
		if node.r == nil {
			newNode := sg.unitNode
			node.r = &newNode
		}
		sg.setRec(node.r, m, r, k, x)
	}
	sg.update(node)
}

type fs struct {
	n, lg int
	seg   [][]int
}

func nfs(n int) *fs {
	res := &fs{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)/64))
		n_ = (n_ + 63) / 64
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func (fs *fs) Has(i int) bool {
	return (fs.seg[0][i/64]>>(i%64))&1 != 0
}

func (fs *fs) Insert(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i/64] |= 1 << (i % 64)
		i /= 64
	}
}

func (fs *fs) Erase(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i/64] &= ^(1 << (i % 64))
		if fs.seg[h][i/64] != 0 {
			break
		}
		i /= 64
	}
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *fs) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		if i/64 == len(fs.seg[h]) {
			break
		}
		d := fs.seg[h][i/64] >> (i % 64)
		if d == 0 {
			i = i/64 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i *= 64
			i += fs.bsf(fs.seg[g][i/64])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *fs) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i/64] << (63 - i%64)
		if d == 0 {
			i = i/64 - 1
			continue
		}
		// find
		i += fs.bsr(d) - (64 - 1)
		for g := h - 1; g >= 0; g-- {
			i *= 64
			i += fs.bsr(fs.seg[g][i/64])
		}

		return i
	}

	return -1
}

func (*fs) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*fs) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}

const INF int = 1e18

type SegmentTree struct {
	n, size int
	seg     []E
	unit    E
}

func NewSegmentTree(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	res.unit = e()
	return res
}

func (st *SegmentTree) Set(index int, value E) {
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

func (st *SegmentTree) Update(index int, value E) {
	index += st.size
	st.seg[index] = op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int) E {
	leftRes, rightRes := st.unit, st.unit
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return op(leftRes, rightRes)
}

func (st *SegmentTree) QueryAll() E { return st.seg[1] }
