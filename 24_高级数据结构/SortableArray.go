// 可排序数组(RangeSort)
// NewSortableArray(nums []int) *SortableArray
// Set(i, v int)
// Get(i int) int
// GetAll() []int
// SortInc(start, end int)
// SortDec(start, end int)
// !每次操作O(log^2(maxValue))
// https://www.luogu.com.cn/problem/P2824
// https://dpair.gitee.io/articles/odt/

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

func init() {
	debug.SetGCPercent(-1)
}

func demo() {
	nums := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sa := NewSortableArray(1000, nums)
	fmt.Println(sa)
	fmt.Println(sa.GetAll())
	sa.SortDec(0, 10)
	fmt.Println(sa.GetAll())
	sa.SortInc(0, 7)
	fmt.Println(sa.GetAll())
	sa.Set(3, 100)
	fmt.Println(sa.GetAll())
	fmt.Println(sa.Get(3))
}

func main() {
	// abc217e()
	abc237g()
	// demo()
	// Logu2824()
}

func abc217e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int32
	fmt.Fscan(in, &q)

	sa := NewSortableArray(1e9, make([]int32, q))
	left, right := int32(0), int32(0)
	append := func(x int32) {
		sa.Set(right, x)
		right++
	}
	popleft := func() (res int32) {
		res = sa.Get(left)
		left++
		return
	}
	sort := func() {
		sa.SortInc(left, right)
	}

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var x int32
			fmt.Fscan(in, &x)
			append(x)
		} else if op == 2 {
			fmt.Fprintln(out, popleft())
		} else {
			sort()
		}
	}
}

func abc237g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q, x int32
	fmt.Fscan(in, &n, &q, &x)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	sa := NewSortableArray(n, nums)
	for i := int32(0); i < q; i++ {
		var t, l, r int32
		fmt.Fscan(in, &t, &l, &r)
		l--
		if t == 1 {
			sa.SortInc(l, r)
		} else {
			sa.SortDec(l, r)
		}
	}

	res := sa.GetAll()
	for i := int32(0); i < n; i++ {
		if res[i] == x {
			fmt.Fprintln(out, i+1)
			return
		}
	}
}

func Logu2824() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	sa := NewSortableArray(n, nums)
	for i := int32(0); i < q; i++ {
		var op, l, r int32
		fmt.Fscan(in, &op, &l, &r)
		l--
		if op == 0 {
			sa.SortInc(l, r)
		} else {
			sa.SortDec(l, r)
		}
	}

	var k int32
	fmt.Fscan(in, &k)
	fmt.Fprintln(out, sa.Get(k-1))
}

type SortableArray struct {
	n        int32
	maxValue int32
	ss       *FastSet
	root     []*SNode // 动态开点线段树的结点
	rev      []bool
}

// 0<=nums[i]<=maxValue
func NewSortableArray(maxValue int32, nums []int32) *SortableArray {
	n := int32(len(nums))
	res := &SortableArray{n: n, maxValue: maxValue + 10, ss: NewFastSet(n)}
	res.init(nums)
	return res
}

// 0<=i<n
// 0<=v<=maxValue
func (sa *SortableArray) Set(i, v int32) {
	sa.splitAt(i)
	sa.splitAt(i + 1)
	sa.rev[i] = false
	sa.root[i].size = 0
	sa.root[i].l = nil
	sa.root[i].r = nil
	sa.setRec(sa.root[i], 0, sa.maxValue, v)
}

// 0<=i<n
func (sa *SortableArray) Get(i int32) int32 {
	p := sa.ss.Prev(i)
	k := i - p
	s := sa.root[p].size
	if sa.rev[p] {
		k = s - 1 - k
	}
	return sa.getDfs(sa.root[p], 0, sa.maxValue, k)
}

func (sa *SortableArray) GetAll() []int32 {
	res := make([]int32, 0, sa.n)
	for i := int32(0); i < sa.n; i++ {
		if sa.ss.Has(i) {
			sa.getAllDfs(sa.root[i], 0, sa.maxValue, sa.rev[i], &res)
		}
	}
	return res
}

func (sa *SortableArray) SortInc(start, end int32) {
	if start >= end {
		return
	}
	sa.splitAt(start)
	sa.splitAt(end)
	for {
		c := sa.root[start]
		i := sa.ss.Next(start + 1)
		if i == end {
			break
		}
		sa.root[start] = sa.merge(0, sa.maxValue, c, sa.root[i])
		sa.ss.Erase(i)
	}
	sa.rev[start] = false
}

func (sa *SortableArray) SortDec(start, end int32) {
	if start >= end {
		return
	}
	sa.SortInc(start, end)
	sa.rev[start] = true
}

func (sa *SortableArray) String() string {
	res := sa.GetAll()
	return fmt.Sprintf("SortableArray%v", res)
}

func (sa *SortableArray) getDfs(node *SNode, l, r, k int32) int32 {
	if r == l+1 {
		return l
	}
	m := (l + r) >> 1
	s := int32(0)
	if node.l != nil {
		s = node.l.size
	}
	if k < s {
		return sa.getDfs(node.l, l, m, k)
	}
	return sa.getDfs(node.r, m, r, k-s)
}

func (sa *SortableArray) getAllDfs(node *SNode, l, r int32, rev bool, key *[]int32) {
	if node == nil || node.size == 0 {
		return
	}
	if r == l+1 {
		for i := int32(0); i < node.size; i++ {
			*key = append(*key, l)
		}
		return
	}
	m := (l + r) >> 1
	if !rev {
		sa.getAllDfs(node.l, l, m, rev, key)
		sa.getAllDfs(node.r, m, r, rev, key)
	} else {
		sa.getAllDfs(node.r, m, r, rev, key)
		sa.getAllDfs(node.l, l, m, rev, key)
	}
}

func (sa *SortableArray) init(nums []int32) {
	sa.rev = make([]bool, sa.n)
	sa.root = make([]*SNode, 0, sa.n)
	for i := int32(0); i < sa.n; i++ {
		sa.ss.Insert(i)
		sa.root = append(sa.root, &SNode{})
		sa.setRec(sa.root[i], 0, sa.maxValue, nums[i])
	}
}

func (sa *SortableArray) splitAt(x int32) {
	if x == sa.n || sa.ss.Has(x) {
		return
	}
	a := sa.ss.Prev(x)
	b := sa.ss.Next(a + 1)
	sa.ss.Insert(x)
	root := sa.root
	if !sa.rev[a] {
		nl, nr := sa.split(root[a], 0, sa.maxValue, x-a)
		root[a], root[x] = nl, nr
		sa.rev[a], sa.rev[x] = false, false
	} else {
		nl, nr := sa.split(root[a], 0, sa.maxValue, b-x)
		root[a], root[x] = nr, nl
		sa.rev[a], sa.rev[x] = true, true
	}
}

func (sa *SortableArray) split(node *SNode, l, r, k int32) (*SNode, *SNode) {
	if k == 0 {
		return nil, node
	}
	if k == node.size {
		return node, nil
	}
	if r == l+1 {
		s := node.size
		node.size = k
		return node, &SNode{size: s - k}
	}
	s := int32(0)
	if node.l != nil {
		s = node.l.size
	}
	b := &SNode{}
	m := (l + r) >> 1
	if k <= s {
		nl, nr := sa.split(node.l, l, m, k)
		b.l, b.r = nr, node.r
		node.l, node.r = nl, nil
	} else {
		nl, nr := sa.split(node.r, m, r, k-s)
		node.r = nl
		b.l, b.r = nil, nr
	}
	sa.update(node)
	sa.update(b)
	return node, b
}

func (sa *SortableArray) merge(l, r int32, a, b *SNode) *SNode {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	if r == l+1 {
		a.size += b.size
		return a
	}
	m := (l + r) >> 1
	a.l = sa.merge(l, m, a.l, b.l)
	a.r = sa.merge(m, r, a.r, b.r)
	sa.update(a)
	return a
}

func (sa *SortableArray) update(node *SNode) {
	if node.l == nil && node.r == nil {
		return
	}
	if node.l == nil {
		node.size = node.r.size
		return
	}
	if node.r == nil {
		node.size = node.l.size
		return
	}
	node.size = node.l.size + node.r.size
}

func (sa *SortableArray) setRec(node *SNode, l, r, k int32) {
	if r == l+1 {
		node.size = 1
		return
	}
	m := (l + r) >> 1
	if k < m {
		if node.l == nil {
			node.l = &SNode{}
		}
		sa.setRec(node.l, l, m, k)
	} else {
		if node.r == nil {
			node.r = &SNode{}
		}
		sa.setRec(node.r, m, r, k)
	}
	sa.update(node)
}

type SNode struct {
	size int32
	l, r *SNode
}

type FastSet struct {
	n, lg int32
	seg   [][]int
}

func NewFastSet(n int32) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = int32(len(seg))
	return res
}

func (fs *FastSet) Has(i int32) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int32) {
	for h := int32(0); h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
}

func (fs *FastSet) Erase(i int32) {
	for h := int32(0); h < fs.lg; h++ {
		fs.seg[h][i>>6] &= ^(1 << (i & 63))
		if fs.seg[h][i>>6] != 0 {
			break
		}
		i >>= 6
	}
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int32) int32 {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := int32(0); h < fs.lg; h++ {
		if i>>6 == int32(len(fs.seg[h])) {
			break
		}
		d := fs.seg[h][i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int32) int32 {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := int32(0); h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int32, f func(i int32)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := int32(0); i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(int(i)))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (*FastSet) bsr(x int) int32 {
	return int32(bits.Len64(uint64(x)) - 1)
}

func (*FastSet) bsf(x int) int32 {
	return int32(bits.TrailingZeros64(uint64(x)))
}
