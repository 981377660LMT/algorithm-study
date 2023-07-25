// 可排序数组(RangeSort)
// NewSortableArray(nums []int) *SortableArray
// Set(i, v int)
// Get(i int) int
// GetAll() []int
// SortInc(start, end int)
// SortDec(start, end int)
// !每次操作O(log^2(maxValue))

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
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
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
	abc217e()
	// abc237g()
	// demo()

}

func abc217e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	sa := NewSortableArray(1e9, make([]int, q))
	left, right := 0, 0
	append := func(x int) {
		sa.Set(right, x)
		right++
	}
	popleft := func() (res int) {
		res = sa.Get(left)
		left++
		return
	}
	sort := func() {
		sa.SortInc(left, right)
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var x int
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

	var n, q, x int
	fmt.Fscan(in, &n, &q, &x)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	sa := NewSortableArray(n, nums)
	for i := 0; i < q; i++ {
		var t, l, r int
		fmt.Fscan(in, &t, &l, &r)
		l--
		if t == 1 {
			sa.SortInc(l, r)
		} else {
			sa.SortDec(l, r)
		}
	}

	res := sa.GetAll()
	for i := 0; i < n; i++ {
		if res[i] == x {
			fmt.Fprintln(out, i+1)
			return
		}
	}
}

type SortableArray struct {
	n        int
	maxValue int
	ss       *FastSet
	root     []*SNode // 动态开点线段树的结点
	rev      []bool
}

// 0<=nums[i]<=maxValue
func NewSortableArray(maxValue int, nums []int) *SortableArray {
	n := len(nums)
	res := &SortableArray{n: n, maxValue: maxValue + 10, ss: NewFastSet(n)}
	res.init(nums)
	return res
}

// 0<=i<n
// 0<=v<=maxValue
func (sa *SortableArray) Set(i, v int) {
	sa.splitAt(i)
	sa.splitAt(i + 1)
	sa.rev[i] = false
	sa.root[i].size = 0
	sa.root[i].l = nil
	sa.root[i].r = nil
	sa.setRec(sa.root[i], 0, sa.maxValue, v)
}

// 0<=i<n
func (sa *SortableArray) Get(i int) int {
	p := sa.ss.Prev(i)
	k := i - p
	s := sa.root[p].size
	if sa.rev[p] {
		k = s - 1 - k
	}
	return sa.getDfs(sa.root[p], 0, sa.maxValue, k)
}

func (sa *SortableArray) GetAll() []int {
	res := make([]int, 0, sa.n)
	for i := 0; i < sa.n; i++ {
		if sa.ss.Has(i) {
			sa.getAllDfs(sa.root[i], 0, sa.maxValue, sa.rev[i], &res)
		}
	}
	return res
}

func (sa *SortableArray) SortInc(start, end int) {
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

func (sa *SortableArray) SortDec(start, end int) {
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

func (sa *SortableArray) getDfs(node *SNode, l, r, k int) int {
	if r == l+1 {
		return l
	}
	m := (l + r) >> 1
	s := 0
	if node.l != nil {
		s = node.l.size
	}
	if k < s {
		return sa.getDfs(node.l, l, m, k)
	}
	return sa.getDfs(node.r, m, r, k-s)
}

func (sa *SortableArray) getAllDfs(node *SNode, l, r int, rev bool, key *[]int) {
	if node == nil || node.size == 0 {
		return
	}
	if r == l+1 {
		for i := 0; i < node.size; i++ {
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

func (sa *SortableArray) init(nums []int) {
	sa.rev = make([]bool, sa.n)
	sa.root = make([]*SNode, 0, sa.n)
	for i := 0; i < sa.n; i++ {
		sa.ss.Insert(i)
		sa.root = append(sa.root, &SNode{})
		sa.setRec(sa.root[i], 0, sa.maxValue, nums[i])
	}
}

func (sa *SortableArray) splitAt(x int) {
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

func (sa *SortableArray) split(node *SNode, l, r, k int) (*SNode, *SNode) {
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
	s := 0
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

func (sa *SortableArray) merge(l, r int, a, b *SNode) *SNode {
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

func (sa *SortableArray) setRec(node *SNode, l, r, k int) {
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
	size int
	l, r *SNode
}

type FastSet struct {
	n, lg int
	seg   [][]int
}

func NewFastSet(n int) *FastSet {
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
	res.lg = len(seg)
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
}

func (fs *FastSet) Erase(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] &= ^(1 << (i & 63))
		if fs.seg[h][i>>6] != 0 {
			break
		}
		i >>= 6
	}
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		if i>>6 == len(fs.seg[h]) {
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
func (fs *FastSet) Prev(i int) int {
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
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}
