// 基本的splay，不维护区间和.
//
// api:
//
//	NewSpalyTreeBasic() *SplayTreeBasic
//	NewRoot() *SplayNode
//	Build(n int32, f func(i int32) E) *SplayNode
//	Size(n *SplayNode) int32
//	Merge(l, r *SplayNode) *SplayNode
//	Split(root *SplayNode, k int32) (*SplayNode, *SplayNode)
//	EnumerateAll(root *SplayNode, f func(E))
//	Get(root **SplayNode, k int32) E
//	Set(root **SplayNode, k int32, x E)
//	Update(root **SplayNode, k int32, x E)
//	Reverse(root **SplayNode, l, r int32)
//	ReverseAll(root *SplayNode)
//  SplitMaxRightByValue(root *SplayNode, check func(E) bool) (*SplayNode, *SplayNode)
//  SplitMaxRightByValueAndCount(root *SplayNode, check func(E, int32) bool) (*SplayNode, *SplayNode)

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	demo()
	// arc153b()
	// abc350f()
}

func demo() {
	S := NewSplayTreeBasic()
	nums := S.Build(10, func(i int32) E { return E(i) })

	c, d := S.SplitMaxRightByValue(nums, func(e E) bool { return e < 5 })
	fmt.Println(S.GetAll(c))
	fmt.Println(S.GetAll(d))
}

// https://atcoder.jp/contests/arc153/tasks/arc153_b
func arc153b() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var H, W int32
	fmt.Fscan(in, &H, &W)
	G := make([]string, H)
	for i := int32(0); i < H; i++ {
		fmt.Fscan(in, &G[i])
	}

	S := NewSplayTreeBasic()

	A, B := make([]int32, H), make([]int32, W)
	for i := int32(0); i < H; i++ {
		A[i] = i
	}
	for i := int32(0); i < W; i++ {
		B[i] = i
	}

	root1 := S.Build(H, func(i int32) E { return E(A[i]) })
	root2 := S.Build(W, func(i int32) E { return E(B[i]) })

	var Q int
	fmt.Fscan(in, &Q)
	for i := 0; i < Q; i++ {
		var a, b int32
		fmt.Fscan(in, &a, &b)
		S.Reverse(&root1, 0, a)
		S.Reverse(&root1, a, H)
		S.Reverse(&root2, 0, b)
		S.Reverse(&root2, b, W)
	}

	A = S.GetAll(root1)
	B = S.GetAll(root2)

	res := make([][]string, H)
	for i := int32(0); i < H; i++ {
		res[i] = make([]string, W)
		for j := int32(0); j < W; j++ {
			res[i][j] = string(G[A[i]][B[j]])
		}
	}

	for i := int32(0); i < H; i++ {
		fmt.Fprintln(out, strings.Join(res[i], ""))
	}
}

// F - Transpose (反转括号内的大小写，区间反转)
// https://atcoder.jp/contests/abc350/tasks/abc350_f
func abc350f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	S := NewSplayTreeBasic()
	n := int32(len(s))
	ptr := int32(0)

	var dfs func(dep int32) *SplayNode
	dfs = func(dep int32) *SplayNode {
		res := S.NewRoot()
		for ptr < n {
			if s[ptr] == '(' {
				ptr++
				child := dfs(dep + 1)
				S.ReverseAll(child)
				res = S.Merge(res, child)
				if s[ptr] != ')' {
					panic("error")
				}
				ptr++
				continue
			}
			if s[ptr] == ')' {
				break
			}
			c := s[ptr]
			if dep&1 == 1 {
				c ^= 32
			}
			node := S.NewNode(E(c))
			res = S.Merge(res, node)
			ptr++
		}
		return res
	}

	root := dfs(0)
	sb := strings.Builder{}
	S.EnumerateAll(root, func(v E) { sb.WriteByte(byte(v)) })
	fmt.Fprintln(out, sb.String())
}

type E = int32

func NewSplayTreeBasic() *SplayTreeBasic {
	return &SplayTreeBasic{}
}

type SplayNode struct {
	rev     bool
	size    int32
	x       E
	p, l, r *SplayNode
}

type SplayTreeBasic struct{}

func (st *SplayTreeBasic) NewRoot() *SplayNode {
	return nil
}

func (st *SplayTreeBasic) Build(n int32, f func(i int32) E) *SplayNode {
	var dfs func(l, r int32) *SplayNode
	dfs = func(l, r int32) *SplayNode {
		if l == r {
			return nil
		}
		if r == l+1 {
			return st.NewNode(f(l))
		}
		m := (l + r) >> 1
		lRoot, rRoot := dfs(l, m), dfs(m+1, r)
		root := st.NewNode(f(m))
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

func (st *SplayTreeBasic) Size(n *SplayNode) int32 {
	if n == nil {
		return 0
	}
	return n.size
}

func (st *SplayTreeBasic) Merge(l, r *SplayNode) *SplayNode {
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

func (st *SplayTreeBasic) Merge3(a, b, c *SplayNode) *SplayNode {
	return st.Merge(st.Merge(a, b), c)
}

func (st *SplayTreeBasic) Merge4(a, b, c, d *SplayNode) *SplayNode {
	return st.Merge(st.Merge(st.Merge(a, b), c), d)
}

func (st *SplayTreeBasic) Split(root *SplayNode, k int32) (*SplayNode, *SplayNode) {
	if k == 0 {
		return nil, root
	}
	if k == root.size {
		return root, nil
	}
	st.splayKth(&root, k-1)
	right := root.r
	root.r = nil
	right.p = nil
	st.nodePushup(root)
	return root, right
}

func (st *SplayTreeBasic) Split3(root *SplayNode, l, r int32) (*SplayNode, *SplayNode, *SplayNode) {
	var nm, nr *SplayNode
	root, nr = st.Split(root, r)
	root, nm = st.Split(root, l)
	return root, nm, nr
}

func (st *SplayTreeBasic) Split4(root *SplayNode, i, j, k int32) (*SplayNode, *SplayNode, *SplayNode, *SplayNode) {
	var d *SplayNode
	root, d = st.Split(root, k)
	a, b, c := st.Split3(root, i, j)
	return a, b, c, d
}

func (st *SplayTreeBasic) gotoBetween(root **SplayNode, l, r int32) {
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

func (st *SplayTreeBasic) EnumerateAll(root *SplayNode, f func(E)) {
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

func (st *SplayTreeBasic) GetAll(root *SplayNode) []E {
	if root == nil {
		return nil
	}
	res := make([]E, 0, root.size)
	st.EnumerateAll(root, func(v E) { res = append(res, v) })
	return res
}

func (st *SplayTreeBasic) Get(root **SplayNode, k int32) E {
	st.splayKth(root, k)
	return st.nodeGet(*root)
}

func (st *SplayTreeBasic) Set(root **SplayNode, k int32, x E) {
	st.splayKth(root, k)
	st.nodeSet(*root, x)
}

func (st *SplayTreeBasic) Update(root **SplayNode, k int32, x E) {
	st.splayKth(root, k)
	st.nodeSet(*root, x)
}

func (st *SplayTreeBasic) Reverse(root **SplayNode, l, r int32) {
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

func (st *SplayTreeBasic) ReverseAll(root *SplayNode) {
	if root != nil {
		st.nodeReverse(root)
	}
}

func (st *SplayTreeBasic) rotate(n *SplayNode) {
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

func (st *SplayTreeBasic) propFromRoot(c *SplayNode) {
	if c.p == nil {
		st.nodePushdown(c)
		return
	}
	st.propFromRoot(c.p)
	st.nodePushdown(c)
}

func (st *SplayTreeBasic) splay(me *SplayNode, propFromRootDone bool) {
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

func (st *SplayTreeBasic) splayKth(root **SplayNode, k int32) {
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
func (st *SplayTreeBasic) SplitMaxRightByValue(root *SplayNode, check func(E) bool) (*SplayNode, *SplayNode) {
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
func (st *SplayTreeBasic) SplitMaxRightByValueAndCount(root *SplayNode, check func(E, int32) bool) (*SplayNode, *SplayNode) {
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

func (st *SplayTreeBasic) findMaxRightByValue(root *SplayNode, check func(E) bool) *SplayNode {
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

func (st *SplayTreeBasic) findMaxRightByValueAndCount(root *SplayNode, check func(E, int32) bool) *SplayNode {
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

// 私有方法需要重写
func (st *SplayTreeBasic) NewNode(x E) *SplayNode {
	return &SplayNode{x: x, size: 1}
}

func (st *SplayTreeBasic) nodePushup(n *SplayNode) {
	n.size = 1
	if n.l != nil {
		n.size += n.l.size
	}
	if n.r != nil {
		n.size += n.r.size
	}
}

func (st *SplayTreeBasic) nodePushdown(n *SplayNode) {
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

func (st *SplayTreeBasic) nodeGet(n *SplayNode) E {
	return n.x
}

func (st *SplayTreeBasic) nodeSet(n *SplayNode, x E) {
	n.x = x
	st.nodePushup(n)
}

func (st *SplayTreeBasic) nodeReverse(n *SplayNode) {
	n.l, n.r = n.r, n.l
	n.rev = !n.rev
}
