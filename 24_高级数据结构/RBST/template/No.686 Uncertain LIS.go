// No.686 Uncertain LIS (BoundedLIS)
// https://yukicoder.me/problems/no/686
// 带上下界限制的LIS
// !有一个长度为N的序列A，求A的一个严格递增子序列B，使得B的长度最大，且对于任意i(1≤i≤|B|)，都有L≤B[i]≤R。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int
	fmt.Fscan(in, &N)
	L, R := make([]int32, N), make([]int32, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &L[i], &R[i])
	}

	dp := NewRBSTMonoidlazy(false) // !dp[i] : LIS长度为i+1时，末尾元素的最小值
	root := dp.NewNode(INF32)

	for i := 0; i < N; i++ {
		if dp.QueryAll(root) != INF32 {
			root = dp.Merge(root, dp.NewNode(INF32))
		}
		var a, b, c, c1, c2 *node
		checkL := func(e int32) bool { return e < L[i] }
		checkR := func(e int32) bool { return e < R[i] }
		a, root = dp.SplitMaxRight(root, checkL)
		b, c = dp.SplitMaxRight(root, checkR)
		c1, c2 = dp.Split(c, 1)  // 删除[R[i], INF) 区间内的最小值
		b = dp.UpdateAll(b, 1)   // [L[i], R[i]) 区间+1
		c1 = dp.Set(c1, 0, L[i]) // 加入L[i]
		root = dp._merge4(a, c1, b, c2)
	}

	check := func(e int32) bool { return e < INF32 }
	n1, _ := dp.SplitMaxRight(root, check)
	if n1 == nil {
		fmt.Fprintln(out, 0)
	} else {
		fmt.Fprintln(out, n1.size)
	}
}

const INF32 int32 = 1e9 + 10

// RangeAddRangeMax
type E = int32
type Id = int32

func e() E                             { return 0 }
func id() Id                           { return 0 }
func op(a, b E) E                      { return maxi32(a, b) }
func mapping(f Id, g E, size uint32) E { return f + g }
func composition(f, g Id) Id           { return f + g }

type node struct {
	rev        uint8
	size       uint32
	value, sum E
	lazy       Id
	l, r       *node
}

type RBSTMonoidLazy struct {
	persistent bool
	x, y, z, w uint32
}

func NewRBSTMonoidlazy(persistent bool) *RBSTMonoidLazy {
	return &RBSTMonoidLazy{
		persistent: persistent,
		x:          123456789,
		y:          362436069,
		z:          521288629,
		w:          88675123,
	}
}

func (rbst *RBSTMonoidLazy) Build(n uint32, f func(uint32) E) *node {
	var dfs func(l, r uint32) *node
	dfs = func(l, r uint32) *node {
		if l == r {
			return nil
		}
		if r == l+1 {
			return rbst.NewNode(f(l))
		}
		mid := (l + r) >> 1
		lRoot := dfs(l, mid)
		rRoot := dfs(mid+1, r)
		root := rbst.NewNode(f(mid))
		root.l = lRoot
		root.r = rRoot
		rbst._update(root)
		return root
	}
	return dfs(0, n)
}

func (rbst *RBSTMonoidLazy) NewRoot() *node { return nil }

func (rbst *RBSTMonoidLazy) NewNode(v E) *node { return &node{value: v, sum: v, lazy: id(), size: 1} }

func (rbst *RBSTMonoidLazy) Merge(lRoot, rRoot *node) *node { return rbst._mergeRec(lRoot, rRoot) }
func (rbst *RBSTMonoidLazy) Split(root *node, k uint32) (*node, *node) {
	return rbst._splitRec(root, k)
}

func (rbst *RBSTMonoidLazy) Set(root *node, k uint32, v E) *node { return rbst._setRec(root, k, v) }
func (rbst *RBSTMonoidLazy) Get(root *node, k uint32) E          { return rbst._getRec(root, k, 0, id()) }
func (rbst *RBSTMonoidLazy) GetAll(root *node) []E {
	res := make([]E, 0, rbst.Size(root))
	var dfs func(root *node, rev uint8, lazy Id)
	dfs = func(root *node, rev uint8, lazy Id) {
		if root == nil {
			return
		}
		me := mapping(lazy, root.value, 1)
		lazy = composition(lazy, root.lazy)
		left, right := root.l, root.r
		if rev == 1 {
			left, right = right, left
		}
		nextRev := rev ^ root.rev
		dfs(left, nextRev, lazy)
		res = append(res, me)
		dfs(right, nextRev, lazy)
	}
	dfs(root, 0, id())
	return res
}

func (rbst *RBSTMonoidLazy) SplitMaxRight(root *node, check func(E) bool) (*node, *node) {
	x := e()
	return rbst._splitMaxRightRec(root, check, &x)
}

func (rbst *RBSTMonoidLazy) QueryRange(root *node, start, end uint32) E {
	if start < 0 {
		start = 0
	}
	if n := rbst.Size(root); end > n {
		end = n
	}
	if start >= end {
		return e()
	}
	return rbst._queryRec(root, start, end, 0)
}

func (rbst *RBSTMonoidLazy) QueryAll(root *node) E {
	if root == nil {
		return e()
	}
	return root.sum
}

func (rbst *RBSTMonoidLazy) Update(root *node, k uint32, x E) *node {
	return rbst._updateRec(root, k, x)
}

func (rbst *RBSTMonoidLazy) UpdateRange(root *node, start, end uint32, x Id) *node {
	if start < 0 {
		start = 0
	}
	if n := rbst.Size(root); end > n {
		end = n
	}
	if start >= end {
		return rbst._copyNode(root)
	}
	return rbst._updateRangeRec(root, start, end, x)
}

func (rbst *RBSTMonoidLazy) UpdateAll(root *node, x Id) *node {
	if root == nil {
		return nil
	}
	return rbst._updateRangeRec(root, 0, rbst.Size(root), x)
}

func (rbst *RBSTMonoidLazy) Reverse(root *node, start, end uint32) *node {
	if start < 0 {
		start = 0
	}
	if n := rbst.Size(root); end > n {
		end = n
	}
	if start >= end {
		return rbst._copyNode(root)
	}
	left, mid, right := rbst._split3(root, start, end)
	mid.rev ^= 1
	mid.l, mid.r = mid.r, mid.l
	return rbst._merge3(left, mid, right)
}

func (rbst *RBSTMonoidLazy) ReverseAll(root *node) *node {
	if root == nil {
		return nil
	}
	root = rbst._copyNode(root)
	root.rev ^= 1
	root.l, root.r = root.r, root.l
	return root
}

func (rbst *RBSTMonoidLazy) Size(root *node) uint32 {
	if root == nil {
		return 0
	}
	return root.size
}

func (rbst *RBSTMonoidLazy) CopyWithin(root *node, target, start, end uint32) *node {
	if !rbst.persistent {
		panic("CopyWithin only works on persistent RBST")
	}
	len := end - start
	p1Left, p1Right := rbst.Split(root, start)
	p2Left, p2Right := rbst.Split(p1Right, len)
	root = rbst.Merge(p1Left, rbst.Merge(p2Left, p2Right))
	p3Left, p3Right := rbst.Split(root, target)
	_, p4Right := rbst.Split(p3Right, len)
	root = rbst.Merge(p3Left, rbst.Merge(p2Left, p4Right))
	return root
}

func (rbst *RBSTMonoidLazy) Pop(root *node, k uint32) (*node, E) {
	n := rbst.Size(root)
	if k < 0 {
		k += n
	}
	x, y, z := rbst._split3(root, k, k+1)
	res := y.value
	newRoot := rbst.Merge(x, z)
	return newRoot, res
}

func (rbst *RBSTMonoidLazy) Erase(root *node, start, end uint32) (remain *node, erased *node) {
	x, y, z := rbst._split3(root, start, end)
	remain = rbst.Merge(x, z)
	erased = y
	return
}

func (rbst *RBSTMonoidLazy) Insert(root *node, k uint32, v E) *node {
	n := rbst.Size(root)
	if k < 0 {
		k += n
	}
	if k < 0 {
		k = 0
	}
	if k > n {
		k = n
	}
	x, y := rbst._splitRec(root, k)
	return rbst._merge3(x, rbst.NewNode(v), y)
}

func (rbst *RBSTMonoidLazy) RotateRight(root *node, start, end, k uint32) *node {
	if end-start <= 1 || k == 0 {
		return rbst._copyNode(root)
	}
	start++
	n := end - start + 1 - k%(end-start+1)
	x, y := rbst.Split(root, start-1)
	y, z := rbst.Split(y, n)
	z, p := rbst.Split(z, end-start+1-n)
	return rbst._merge4(x, z, y, p)
}

func (rbst *RBSTMonoidLazy) RotateLeft(root *node, start, end, k uint32) *node {
	if end-start <= 1 || k == 0 {
		return rbst._copyNode(root)
	}
	start++
	k %= (end - start + 1)
	if k == 0 {
		return rbst._copyNode(root)
	}
	x, y := rbst.Split(root, start-1)
	y, z := rbst.Split(y, k)
	z, p := rbst.Split(z, end-start+1-k)
	return rbst._merge4(x, z, y, p)
}

func (rbst *RBSTMonoidLazy) RotateRightAll(root *node, k uint32) *node {
	n := rbst.Size(root)
	if k >= n {
		k %= n
	}
	if k == 0 {
		return rbst._copyNode(root)
	}
	a, b := rbst.Split(root, n-k)
	return rbst.Merge(b, a)
}

func (rbst *RBSTMonoidLazy) RotateLeftAll(root *node, k uint32) *node {
	n := rbst.Size(root)
	if k >= n {
		k %= n
	}
	if k == 0 {
		return rbst._copyNode(root)
	}
	a, b := rbst.Split(root, k)
	return rbst.Merge(b, a)
}

func (rbst *RBSTMonoidLazy) _merge3(a, b, c *node) *node {
	return rbst.Merge(rbst.Merge(a, b), c)
}

func (rbst *RBSTMonoidLazy) _merge4(a, b, c, d *node) *node {
	return rbst.Merge(rbst.Merge(rbst.Merge(a, b), c), d)
}

func (rbst *RBSTMonoidLazy) _split3(root *node, l, r uint32) (*node, *node, *node) {
	root, nr := rbst.Split(root, r)
	root, nm := rbst.Split(root, l)
	return root, nm, nr
}

func (rbst *RBSTMonoidLazy) _split4(root *node, i, j, k uint32) (*node, *node, *node, *node) {
	root, d := rbst.Split(root, k)
	a, b, c := rbst._split3(root, i, j)
	return a, b, c, d
}

func (rbst *RBSTMonoidLazy) _copyNode(n *node) *node {
	if n == nil || !rbst.persistent {
		return n
	}
	return &node{
		l: n.l, r: n.r,
		value: n.value, sum: n.sum,
		lazy: n.lazy,
		size: n.size, rev: n.rev,
	}
}

func (rbst *RBSTMonoidLazy) _nextRand() uint32 {
	t := rbst.x ^ (rbst.x << 11)
	rbst.x, rbst.y, rbst.z = rbst.y, rbst.z, rbst.w
	rbst.w = (rbst.w ^ (rbst.w >> 19)) ^ (t ^ (t >> 8))
	return rbst.w
}

func (rbst *RBSTMonoidLazy) _propagate(c *node) {
	blLazy := c.lazy != id()
	blRev := c.rev == 1
	if blLazy || blRev {
		c.l, c.r = rbst._copyNode(c.l), rbst._copyNode(c.r)
	}
	if blRev {
		if left := c.l; left != nil {
			left.rev ^= 1
			left.l, left.r = left.r, left.l
		}
		if right := c.r; right != nil {
			right.rev ^= 1
			right.l, right.r = right.r, right.l
		}
		c.rev = 0
	}
	if blLazy {
		if left := c.l; left != nil {
			left.value = mapping(c.lazy, left.value, 1)
			left.sum = mapping(c.lazy, left.sum, left.size)
			left.lazy = composition(c.lazy, left.lazy)
		}
		if right := c.r; right != nil {
			right.value = mapping(c.lazy, right.value, 1)
			right.sum = mapping(c.lazy, right.sum, right.size)
			right.lazy = composition(c.lazy, right.lazy)
		}
		c.lazy = id()
	}
}

func (rbst *RBSTMonoidLazy) _update(c *node) {
	c.size = 1
	c.sum = c.value
	if left := c.l; left != nil {
		c.size += left.size
		c.sum = op(left.sum, c.sum)
	}
	if right := c.r; right != nil {
		c.size += right.size
		c.sum = op(c.sum, right.sum)
	}
}

func (rbst *RBSTMonoidLazy) _mergeRec(lRoot, rRoot *node) *node {
	if lRoot == nil || rRoot == nil {
		if lRoot == nil {
			return rRoot
		}
		return lRoot
	}
	sl, sr := lRoot.size, rRoot.size
	if rbst._nextRand()%(sl+sr) < sl {
		rbst._propagate(lRoot)
		lRoot = rbst._copyNode(lRoot)
		lRoot.r = rbst._mergeRec(lRoot.r, rRoot)
		rbst._update(lRoot)
		return lRoot
	}
	rbst._propagate(rRoot)
	rRoot = rbst._copyNode(rRoot)
	rRoot.l = rbst._mergeRec(lRoot, rRoot.l)
	rbst._update(rRoot)
	return rRoot
}

func (rbst *RBSTMonoidLazy) _splitRec(root *node, k uint32) (*node, *node) {
	if root == nil {
		return nil, nil
	}
	rbst._propagate(root)
	sl := uint32(0)
	if root.l != nil {
		sl = root.l.size
	}
	if k <= sl {
		nl, nr := rbst._splitRec(root.l, k)
		root = rbst._copyNode(root)
		root.l = nr
		rbst._update(root)
		return nl, root
	}
	nl, nr := rbst._splitRec(root.r, k-(1+sl))
	root = rbst._copyNode(root)
	root.r = nl
	rbst._update(root)
	return root, nr
}

func (rbst *RBSTMonoidLazy) _setRec(root *node, k uint32, v E) *node {
	if root == nil {
		return root
	}
	rbst._propagate(root)
	sl := uint32(0)
	if root.l != nil {
		sl = root.l.size
	}
	if k < sl {
		root = rbst._copyNode(root)
		root.l = rbst._setRec(root.l, k, v)
		rbst._update(root)
		return root
	}
	if k == sl {
		root = rbst._copyNode(root)
		root.value = v
		rbst._update(root)
		return root
	}
	root = rbst._copyNode(root)
	root.r = rbst._setRec(root.r, k-(1+sl), v)
	rbst._update(root)
	return root
}

func (rbst *RBSTMonoidLazy) _getRec(root *node, k uint32, rev uint8, lazy Id) E {
	left, right := root.l, root.r
	if rev == 1 {
		left, right = right, left
	}
	sl := uint32(0)
	if left != nil {
		sl = left.size
	}
	if k == sl {
		return mapping(lazy, root.value, 1)
	}
	lazy = composition(root.lazy, lazy)
	rev ^= root.rev
	if k < sl {
		return rbst._getRec(left, k, rev, lazy)
	}
	return rbst._getRec(right, k-(1+sl), rev, lazy)
}

func (rbst *RBSTMonoidLazy) _splitMaxRightRec(root *node, check func(E) bool, x *E) (*node, *node) {
	if root == nil {
		return nil, nil
	}
	rbst._propagate(root)
	root = rbst._copyNode(root)
	y := op(*x, root.sum)
	if check(y) {
		*x = y
		return root, nil
	}
	left, right := root.l, root.r
	if left != nil {
		y = op(*x, left.sum)
		if !check(y) {
			n1, n2 := rbst._splitMaxRightRec(left, check, x)
			root.l = n2
			rbst._update(root)
			return n1, root
		}
		*x = y
	}
	y = op(*x, root.value)
	if !check(y) {
		root.l = nil
		rbst._update(root)
		return left, root
	}
	*x = y
	n1, n2 := rbst._splitMaxRightRec(right, check, x)
	root.r = n1
	rbst._update(root)
	return root, n2
}

func (rbst *RBSTMonoidLazy) _queryRec(root *node, l, r uint32, rev uint8) E {
	if l == 0 && r == root.size {
		return root.sum
	}
	left, right := root.l, root.r
	if rev == 1 {
		left, right = right, left
	}
	leftSize := rbst.Size(left)
	nextRev := rev ^ root.rev
	res := e()
	if l < leftSize {
		y := rbst._queryRec(left, l, min32(r, leftSize), nextRev)
		res = op(res, mapping(root.lazy, y, min32(r, leftSize)-l))
	}
	if l <= leftSize && leftSize < r {
		res = op(res, root.value)
	}
	k := 1 + leftSize
	if k < r {
		y := rbst._queryRec(right, max32(k, l)-k, r-k, nextRev)
		res = op(res, mapping(root.lazy, y, r-max32(k, l)))
	}
	return res
}

func (rbst *RBSTMonoidLazy) _updateRec(root *node, k uint32, x E) *node {
	if root == nil {
		return root
	}
	rbst._propagate(root)
	sl := uint32(0)
	if root.l != nil {
		sl = root.l.size
	}
	if k < sl {
		root = rbst._copyNode(root)
		root.l = rbst._updateRec(root.l, k, x)
		rbst._update(root)
		return root
	}
	if k == sl {
		root = rbst._copyNode(root)
		root.value = op(root.value, x)
		rbst._update(root)
		return root
	}
	root = rbst._copyNode(root)
	root.r = rbst._updateRec(root.r, k-(1+sl), x)
	rbst._update(root)
	return root
}

func (rbst *RBSTMonoidLazy) _updateRangeRec(root *node, l, r uint32, a Id) *node {
	rbst._propagate(root)
	root = rbst._copyNode(root)
	if l == 0 && r == root.size {
		root.value = mapping(a, root.value, 1)
		root.sum = mapping(a, root.sum, root.size)
		root.lazy = a
		return root
	}
	leftSize := rbst.Size(root.l)
	if l < leftSize {
		root.l = rbst._updateRangeRec(root.l, l, min32(r, leftSize), a)
	}
	if l <= leftSize && leftSize < r {
		root.value = mapping(a, root.value, 1)
	}
	k := 1 + leftSize
	if k < r {
		root.r = rbst._updateRangeRec(root.r, max32(k, l)-k, r-k, a)
	}
	rbst._update(root)
	return root
}

func min32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func maxi32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
