// 维护区间和的 splay 树. monoid 满足交换律.
// api
// - NewSpalyTreeMonoidCommutative() *SplayTreeMonoidCommutative
// - NewRoot() *SplayNode
// - Build(n int32, f func(i int32) E) *SplayNode
// - Size(n *SplayNode) int32
// - Merge(l, r *SplayNode) *SplayNode
// - Split(root *SplayNode, k int32) (*SplayNode, *SplayNode)
// - EnumerateAll(root *SplayNode, f func(E))
// - Get(root **SplayNode, k int32) E
// - Set(root **SplayNode, k int32, x E)
// - Update(root **SplayNode, k int32, x E)
// - QueryRange(root **SplayNode, l, r int32) E
// - QueryAll(root *SplayNode) E
// - Reverse(root **SplayNode, l, r int32)
// - ReverseAll(root *SplayNode)
// - SplitMaxRightByValue(root *SplayNode, check func(E) bool) (*SplayNode, *SplayNode)
// - SplitMaxRightByValueAndCount(root *SplayNode, check func(E, int32) bool) (*SplayNode, *SplayNode)
// - SplitMaxRightBySum(root *SplayNode, check func(E) bool) (*SplayNode, *SplayNode)

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	yosupoRangeReverseRangeSum()
}

func demo() {
	S := NewSpalyTreeMonoidCommutative()
	nums := S.Build(10, func(i int32) E { return E(i) })
	S.Update(&nums, 0, 1)
	S.Update(&nums, 1, 2)
	S.Update(&nums, 2, 3)
	S.Set(&nums, 3, 4)

	getAll := func() []E {
		res := []E{}
		S.EnumerateAll(nums, func(x E) {
			res = append(res, x)
		})
		return res
	}

	fmt.Println(getAll())

	fmt.Println(S.QueryRange(&nums, 0, 4))
	fmt.Println(S.QueryAll(nums))
	fmt.Println(S.Get(&nums, 3))

}

// 区间翻转, 区间和(RangeFlipRangeSum)
// https://judge.yosupo.jp/problem/range_reverse_range_sum
func yosupoRangeReverseRangeSum() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, Q int
	fmt.Fscan(in, &N, &Q)
	A := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &A[i])
	}

	S := NewSpalyTreeMonoidCommutative()
	nums := S.Build(int32(N), func(i int32) E { return E(A[i]) })
	for i := 0; i < Q; i++ {
		var t, L, R int
		fmt.Fscan(in, &t, &L, &R)
		if t == 0 {
			S.Reverse(&nums, int32(L), int32(R))
		} else {
			fmt.Fprintln(out, S.QueryRange(&nums, int32(L), int32(R)))
		}
	}
}

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }

func NewSpalyTreeMonoidCommutative() *SpalyTreeMonoidCommutative {
	return &SpalyTreeMonoidCommutative{}
}

type SplayNode struct {
	rev     bool
	size    int32
	x, sum  E
	p, l, r *SplayNode
}

type SpalyTreeMonoidCommutative struct{}

func (st *SpalyTreeMonoidCommutative) NewRoot() *SplayNode {
	return nil
}

func (st *SpalyTreeMonoidCommutative) Build(n int32, f func(i int32) E) *SplayNode {
	var dfs func(l, r int32) *SplayNode
	dfs = func(l, r int32) *SplayNode {
		if l == r {
			return nil
		}
		if r == l+1 {
			return st.newNode(f(l))
		}
		m := (l + r) >> 1
		lRoot, rRoot := dfs(l, m), dfs(m+1, r)
		root := st.newNode(f(m))
		root.l, root.r = lRoot, rRoot
		if lRoot != nil {
			lRoot.p = root
		}
		if rRoot != nil {
			rRoot.p = root
		}
		st.nodePushup(root)
		return root
	}
	return dfs(0, n)
}

func (st *SpalyTreeMonoidCommutative) Size(n *SplayNode) int32 {
	if n == nil {
		return 0
	}
	return n.size
}

func (st *SpalyTreeMonoidCommutative) Merge(l, r *SplayNode) *SplayNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	st.splayKth(&r, 0)
	r.l = l
	l.p = r
	st.nodePushup(r)
	return r
}

func (st *SpalyTreeMonoidCommutative) Merge3(a, b, c *SplayNode) *SplayNode {
	return st.Merge(st.Merge(a, b), c)
}

func (st *SpalyTreeMonoidCommutative) Merge4(a, b, c, d *SplayNode) *SplayNode {
	return st.Merge(st.Merge(st.Merge(a, b), c), d)
}

func (st *SpalyTreeMonoidCommutative) Split(root *SplayNode, k int32) (*SplayNode, *SplayNode) {
	if k == 0 {
		return nil, root
	}
	if k == root.size {
		return root, nil
	}
	st.splayKth(&root, k-1)
	right := root.r
	root.r = nil
	root.p = nil
	st.nodePushup(root)
	return root, right
}

func (st *SpalyTreeMonoidCommutative) Split3(root *SplayNode, l, r int32) (*SplayNode, *SplayNode, *SplayNode) {
	var nm, nr *SplayNode
	root, nr = st.Split(root, r)
	root, nm = st.Split(root, l)
	return root, nm, nr
}

func (st *SpalyTreeMonoidCommutative) Split4(root *SplayNode, i, j, k int32) (*SplayNode, *SplayNode, *SplayNode, *SplayNode) {
	var d *SplayNode
	root, d = st.Split(root, k)
	a, b, c := st.Split3(root, i, j)
	return a, b, c, d
}

func (st *SpalyTreeMonoidCommutative) EnumerateAll(root *SplayNode, f func(E)) {
	var dfs func(*SplayNode)
	dfs = func(root *SplayNode) {
		if root == nil {
			return
		}
		st.nodePushdown(root)
		dfs(root.l)
		f(st.nodeGet(root))
		dfs(root.r)
	}
	dfs(root)
}
func (st *SpalyTreeMonoidCommutative) GetAll(root *SplayNode) []E {
	if root == nil {
		return nil
	}
	res := make([]E, 0, root.size)
	st.EnumerateAll(root, func(v E) { res = append(res, v) })
	return res
}
func (st *SpalyTreeMonoidCommutative) Get(root **SplayNode, k int32) E {
	st.splayKth(root, k)
	return st.nodeGet(*root)
}

func (st *SpalyTreeMonoidCommutative) Set(root **SplayNode, k int32, x E) {
	st.splayKth(root, k)
	st.nodeSet(*root, x)
}

func (st *SpalyTreeMonoidCommutative) Update(root **SplayNode, k int32, x E) {
	st.splayKth(root, k)
	st.nodeUpdate(*root, x)
}

func (st *SpalyTreeMonoidCommutative) Reverse(root **SplayNode, l, r int32) {
	if *root == nil {
		return
	}
	if l < 0 {
		l = 0
	}
	if s := (*root).size; r > s {
		r = s
	}
	if l >= r {
		return
	}
	st.gotoBetween(root, l, r)
	st.nodeReverse(*root)
	st.splay(*root, true)
}

func (st *SpalyTreeMonoidCommutative) ReverseAll(root *SplayNode) {
	if root != nil {
		st.nodeReverse(root)
	}
}

func (st *SpalyTreeMonoidCommutative) QueryRange(root **SplayNode, l, r int32) E {
	if *root == nil {
		return e()
	}
	if l < 0 {
		l = 0
	}
	if s := (*root).size; r > s {
		r = s
	}
	if l >= r {
		return e()
	}
	st.gotoBetween(root, l, r)
	res := (*root).sum
	st.splay(*root, true)
	return res
}

func (st *SpalyTreeMonoidCommutative) QueryAll(root *SplayNode) E {
	if root == nil {
		return e()
	}
	return root.sum
}

func (st *SpalyTreeMonoidCommutative) gotoBetween(root **SplayNode, l, r int32) {
	if l == 0 && r == (*root).size {
		return
	}
	if l == 0 {
		st.splayKth(root, r)
		*root = (*root).l
		return
	}
	if r == (*root).size {
		st.splayKth(root, l-1)
		*root = (*root).r
		return
	}
	st.splayKth(root, r)
	rp := *root
	(*root) = rp.l
	(*root).p = nil
	st.splayKth(root, l-1)
	(*root).p = rp
	rp.l = *root
	st.nodePushup(rp)
	*root = (*root).r
}

func (st *SpalyTreeMonoidCommutative) rotate(n *SplayNode) {
	var pp, p, c *SplayNode
	p = n.p
	pp = p.p
	if p.l == n {
		c = n.r
		n.r = p
		p.l = c
	} else {
		c = n.l
		n.l = p
		p.r = c
	}
	if pp != nil && pp.l == p {
		pp.l = n
	}
	if pp != nil && pp.r == p {
		pp.r = n
	}
	n.p = pp
	p.p = n
	if c != nil {
		c.p = p
	}
}

func (st *SpalyTreeMonoidCommutative) propFromRoot(c *SplayNode) {
	if c.p == nil {
		st.nodePushdown(c)
		return
	}
	st.propFromRoot(c.p)
	st.nodePushdown(c)
}

func (st *SpalyTreeMonoidCommutative) splay(me *SplayNode, propFromRootDone bool) {
	if !propFromRootDone {
		st.propFromRoot(me)
	}
	st.nodePushdown(me)
	for me.p != nil {
		p := me.p
		pp := p.p
		if pp == nil {
			st.rotate(me)
			st.nodePushup(p)
			break
		}
		same := (p.l == me && pp.l == p) || (p.r == me && pp.r == p)
		if same {
			st.rotate(p)
			st.rotate(me)
		} else {
			st.rotate(me)
		}
		st.nodePushup(pp)
		st.nodePushup(p)
	}
	st.nodePushup(me)
}

func (st *SpalyTreeMonoidCommutative) splayKth(root **SplayNode, k int32) {
	for {
		st.nodePushdown(*root)
		sl := st.Size((*root).l)
		if k == sl {
			break
		}
		if k < sl {
			*root = (*root).l
		} else {
			k -= sl + 1
			*root = (*root).r
		}
	}
	st.splay(*root, true)
}

// 分离出的左侧节点值满足check函数.
func (st *SpalyTreeMonoidCommutative) SplitMaxRightByValue(root *SplayNode, check func(E) bool) (*SplayNode, *SplayNode) {
	if root == nil {
		return nil, nil
	}
	c := st.findMaxRightByValue(root, check)
	if c == nil {
		st.splay(root, true)
		return nil, root
	}
	st.splay(c, true)
	right := c.r
	if right == nil {
		return c, nil
	}
	right.p = nil
	c.r = nil
	st.nodePushup(c)
	return c, right
}

// 分离出的左侧节点之和与
func (st *SpalyTreeMonoidCommutative) SplitMaxRightByValueAndCount(root *SplayNode, check func(E, int32) bool) (*SplayNode, *SplayNode) {
	if root == nil {
		return nil, nil
	}
	c := st.findMaxRightByValueAndCount(root, check)
	if c == nil {
		st.splay(root, true)
		return nil, root
	}
	st.splay(c, true)
	right := c.r
	if right == nil {
		return c, nil
	}
	right.p = nil
	c.r = nil
	st.nodePushup(c)
	return c, right
}

// 分离出的左侧节点之和满足check函数.
func (st *SpalyTreeMonoidCommutative) SplitMaxRightBySum(root *SplayNode, check func(E) bool) (*SplayNode, *SplayNode) {
	if root == nil {
		return nil, nil
	}
	c := st.findMaxRightBySum(root, check)
	if c == nil {
		st.splay(root, true)
		return nil, root
	}
	st.splay(c, true)
	right := c.r
	if right == nil {
		return c, nil
	}
	right.p = nil
	c.r = nil
	st.nodePushup(c)
	return c, right
}

func (st *SpalyTreeMonoidCommutative) findMaxRightByValue(root *SplayNode, check func(E) bool) *SplayNode {
	var lastOk, last *SplayNode
	for root != nil {
		last = root
		st.nodePushdown(root)
		if check(root.x) {
			lastOk = root
			root = root.r
		} else {
			root = root.l
		}
	}
	st.splay(last, true)
	return lastOk
}

func (st *SpalyTreeMonoidCommutative) findMaxRightByValueAndCount(root *SplayNode, check func(E, int32) bool) *SplayNode {
	var lastOk, last *SplayNode
	var n int32
	for root != nil {
		last = root
		st.nodePushdown(root)
		ns := st.Size(root.l)
		if check(root.x, n+ns+1) {
			lastOk = root
			n += ns + 1
			root = root.r
		} else {
			root = root.l
		}
	}
	st.splay(last, true)
	return lastOk
}

func (st *SpalyTreeMonoidCommutative) findMaxRightBySum(root *SplayNode, check func(E) bool) *SplayNode {
	prod := e()
	var lastOk, last *SplayNode
	for root != nil {
		last = root
		st.nodePushdown(root)
		lprod := prod
		if root.l != nil {
			lprod = op(lprod, root.l.sum)
		}
		lprod = op(lprod, root.x)
		if check(lprod) {
			prod = lprod
			lastOk = root
			root = root.r
		} else {
			root = root.l
		}
	}
	st.splay(last, true)
	return lastOk
}

// 私有方法需要重写
func (st *SpalyTreeMonoidCommutative) newNode(x E) *SplayNode {
	return &SplayNode{x: x, sum: x, size: 1}
}

func (st *SpalyTreeMonoidCommutative) nodePushup(n *SplayNode) {
	n.size = 1
	n.sum = n.x
	if n.l != nil {
		n.size += n.l.size
		n.sum = op(n.l.sum, n.sum)
	}
	if n.r != nil {
		n.size += n.r.size
		n.sum = op(n.sum, n.r.sum)
	}
}

func (st *SpalyTreeMonoidCommutative) nodePushdown(n *SplayNode) {
	if n.rev {
		if left := n.l; left != nil {
			left.rev = !left.rev
			left.l, left.r = left.r, left.l
		}
		if right := n.r; right != nil {
			right.rev = !right.rev
			right.l, right.r = right.r, right.l
		}
		n.rev = false
	}
}

func (st *SpalyTreeMonoidCommutative) nodeGet(n *SplayNode) E {
	return n.x
}

func (st *SpalyTreeMonoidCommutative) nodeSet(n *SplayNode, x E) {
	n.x = x
	st.nodePushup(n)
}

func (st *SpalyTreeMonoidCommutative) nodeUpdate(n *SplayNode, x E) {
	n.x = op(n.x, x)
	st.nodePushup(n)
}

func (st *SpalyTreeMonoidCommutative) nodeReverse(n *SplayNode) {
	n.l, n.r = n.r, n.l
	n.rev = !n.rev
}
