// dynamicHashString/hashStringDynamic
// 动态哈希.
//
// api:
//  1. Insert(i, c): 在第i个位置插入字符c.
//  2. Pop(i): 删除第i个位置的字符.
//  3. Set(i, c): 将第i个位置的字符设置为c.
//  4. Get(start, end): 获取[start, end)的哈希值.
//  5. GetAll(): 获取所有字符的哈希值.
//  8. RotateRight(start, end, k): 将[start, end)的字符右旋k位.
//  9. RotateLeft(start, end, k): 将[start, end)的字符左旋k位.
//  10. RotateRightAll(k): 将所有字符右旋k位.
//  11. RotateLeftAll(k): 将所有字符左旋k位.

package main

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"
)

func main() {
	// demo()
	test()
	testTime()
}

func init() {
	debug.SetGCPercent(-1)
}

func demo() {
	s := "asezfvgbadpihoamgkcmco"
	base := NewDynamicHashStringBase(len(s), 37)
	hs := NewDynamicHashStringX(len(s), func(i int) uint { return uint(s[i]) }, base)
	fmt.Println(hs.Get(0, 1))
	fmt.Println(hs.Get(1, 2))
	fmt.Println(hs.Get(2, 3))
	hs.Set(0, 1)
	fmt.Println(hs.Get(0, 1))
	fmt.Println(hs.Get(0, 1))
	fmt.Println(hs.Get(1, 2))
	fmt.Println(hs.Get(2, 3))
}

// https://leetcode.cn/problems/sum-of-scores-of-built-strings/description/
func sumScores(s string) int64 {
	n := len(s)
	base := NewDynamicHashStringBase(n, 0)
	hasher := NewDynamicHashStringX(n, func(i int) uint { return uint(s[i]) }, base)
	countPre := func(curLen, start int) int {
		left, right := 1, curLen
		for left <= right {
			mid := (left + right) >> 1
			hash1 := hasher.Get(start, start+mid)
			hash2 := hasher.Get(0, mid)
			if hash1 == hash2 {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}

		return right
	}

	res := 0
	for i := 1; i < n+1; i++ {
		if s[0] != s[n-i] {
			continue
		}
		count := countPre(i, n-i)
		res += count
	}

	return int64(res)

}

const (
	hashStringMod    uint = (1 << 61) - 1
	hashStringMask30 uint = (1 << 30) - 1
	hashStringMask31 uint = (1 << 31) - 1
	hashStringMASK61 uint = hashStringMod
)

type DynamicHashStringBase struct {
	n    int
	powb []uint
}

// base: 0 表示随机生成
func NewDynamicHashStringBase(n int, base uint) *DynamicHashStringBase {
	res := &DynamicHashStringBase{}
	if base == 0 {
		base = uint(37 + rand.Intn(1e9))
	}
	powb := make([]uint, n+1)
	powb[0] = 1
	for i := 1; i <= n; i++ {
		powb[i] = res.Mul(powb[i-1], base)
	}
	res.n = n
	res.powb = powb
	return res
}

// h1 <- h2. len(h2) == k.
func (hsb *DynamicHashStringBase) Concat(h1, h2, h2Len uint) uint {
	return hsb.Mod(hsb.Mul(h1, hsb.powb[h2Len]) + h2)
}

// a*b % (2^61-1)
func (hsb *DynamicHashStringBase) Mul(a, b uint) uint {
	au := a >> 31
	ad := a & hashStringMask31
	bu := b >> 31
	bd := b & hashStringMask31
	mid := ad*bu + au*bd
	midu := mid >> 30
	midd := mid & hashStringMask30
	return hsb.Mod(au*bu<<1 + midu + (midd << 31) + ad*bd)
}

// x % (2^61-1)
func (hsb *DynamicHashStringBase) Mod(x uint) uint {
	xu := x >> 61
	xd := x & hashStringMASK61
	res := xu + xd
	if res >= hashStringMod {
		res -= hashStringMod
	}
	return res
}

type DynamicHashStringX struct {
	n    int
	base *DynamicHashStringBase
	rbst *RBSTMonoidNoncommutative
	root *node
}

func NewDynamicHashStringX(n int, f func(i int) uint, base *DynamicHashStringBase) *DynamicHashStringX {
	res := &DynamicHashStringX{n: n, base: base}
	res.rbst = NewRBSTMonoidNoncommutative(
		E{},
		func(a, b E) E {
			return E{hash: base.Concat(a.hash, b.hash, uint(b.len)), len: a.len + b.len}
		},
		false, // 不使用持久化
	)
	res.root = res.rbst.Build(uint32(n), func(i uint32) E { return fromElement(f(int(i))) })
	return res
}

func (hs *DynamicHashStringX) Insert(index int, c uint) {
	hs.root = hs.rbst.Insert(hs.root, uint32(index), fromElement(c))
	hs.n++
}

func (hs *DynamicHashStringX) Pop(index int) uint {
	hs.n--
	a, b := hs.rbst.Pop(hs.root, uint32(index))
	hs.root = a
	return b.hash
}

func (hs *DynamicHashStringX) Get(start, end int) uint {
	if start < 0 {
		start = 0
	}
	if end > hs.n {
		end = hs.n
	}
	if start >= end {
		return 0
	}
	if start == 0 && end == hs.n {
		return hs.GetAll()
	}
	return hs.rbst.QueryRange(hs.root, uint32(start), uint32(end)).hash
}

func (hs *DynamicHashStringX) GetAll() uint { return hs.rbst.QueryAll(hs.root).hash }

func (hs *DynamicHashStringX) Set(index int, c uint) {
	if index < 0 || index >= hs.n {
		return
	}
	hs.root = hs.rbst.Set(hs.root, uint32(index), fromElement(c))
}

func (hs *DynamicHashStringX) Size() int { return hs.n }

func (hs *DynamicHashStringX) RotateRight(start, end, k int) {
	if start < 0 {
		start = 0
	}
	if end > hs.n {
		end = hs.n
	}
	if start >= end {
		return
	}
	if end-start <= 1 || k == 0 {
		return
	}
	hs.root = hs.rbst.RotateRight(hs.root, uint32(start), uint32(end), uint32(k))
}

func (hs *DynamicHashStringX) RotateLeft(start, end, k int) {
	if start < 0 {
		start = 0
	}
	if end > hs.n {
		end = hs.n
	}
	if start >= end {
		return
	}
	if end-start <= 1 || k == 0 {
		return
	}
	hs.root = hs.rbst.RotateLeft(hs.root, uint32(start), uint32(end), uint32(k))
}

func (hs *DynamicHashStringX) RotateRightAll(k int) {
	hs.root = hs.rbst.RotateRightAll(hs.root, uint32(k))
}

func (hs *DynamicHashStringX) RotateLeftAll(k int) {
	hs.root = hs.rbst.RotateLeftAll(hs.root, uint32(k))
}

type E = struct {
	hash uint
	len  int32
}

func fromElement(v uint) E { return E{hash: v, len: 1} }

type node struct {
	l, r   *node
	v, sum E
	size   uint32
}

type RBSTMonoidNoncommutative struct {
	e          E
	op         func(a, b E) E
	persistent bool
	x, y, z, w uint32
}

func NewRBSTMonoidNoncommutative(e E, op func(a, b E) E, persistent bool) *RBSTMonoidNoncommutative {
	return &RBSTMonoidNoncommutative{
		e:          e,
		op:         op,
		persistent: persistent,
		x:          123456789,
		y:          362436069,
		z:          521288629,
		w:          88675123,
	}
}

func (rbst *RBSTMonoidNoncommutative) Build(n uint32, f func(uint32) E) *node {
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

func (rbst *RBSTMonoidNoncommutative) NewRoot() *node { return nil }

func (rbst *RBSTMonoidNoncommutative) NewNode(v E) *node {
	return &node{v: v, sum: v, size: 1}
}

func (rbst *RBSTMonoidNoncommutative) Merge(lRoot, rRoot *node) *node {
	return rbst._mergeRec(lRoot, rRoot)
}
func (rbst *RBSTMonoidNoncommutative) Split(root *node, k uint32) (*node, *node) {
	return rbst._splitRec(root, k)
}

func (rbst *RBSTMonoidNoncommutative) Set(root *node, k uint32, v E) *node {
	return rbst._setRec(root, k, v)
}
func (rbst *RBSTMonoidNoncommutative) Get(root *node, k uint32) E { return rbst._getRec(root, k) }
func (rbst *RBSTMonoidNoncommutative) GetAll(root *node) []E {
	res := make([]E, 0, rbst.Size(root))
	var dfs func(root *node)
	dfs = func(root *node) {
		if root == nil {
			return
		}
		dfs(root.l)
		res = append(res, root.v)
		dfs(root.r)
	}
	dfs(root)
	return res
}

func (rbst *RBSTMonoidNoncommutative) SplitMaxRight(root *node, check func(E) bool) (*node, *node) {
	x := rbst.e
	return rbst._splitMaxRightRec(root, check, &x)
}

func (rbst *RBSTMonoidNoncommutative) QueryRange(root *node, start, end uint32) E {
	if start < 0 {
		start = 0
	}
	if n := rbst.Size(root); end > n {
		end = n
	}
	if start >= end {
		return rbst.e
	}
	return rbst._queryRec(root, start, end)
}

func (rbst *RBSTMonoidNoncommutative) QueryAll(root *node) E {
	if root == nil {
		return rbst.e
	}
	return root.sum
}

func (rbst *RBSTMonoidNoncommutative) Update(root *node, k uint32, x E) *node {
	return rbst._updateRec(root, k, x)
}

func (rbst *RBSTMonoidNoncommutative) Size(root *node) uint32 {
	if root == nil {
		return 0
	}
	return root.size
}

func (rbst *RBSTMonoidNoncommutative) CopyWithin(root *node, target, start, end uint32) *node {
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

func (rbst *RBSTMonoidNoncommutative) Pop(root *node, k uint32) (*node, E) {
	n := rbst.Size(root)
	if k < 0 {
		k += n
	}
	x, y, z := rbst._split3(root, k, k+1)
	res := y.v
	newRoot := rbst.Merge(x, z)
	return newRoot, res
}

func (rbst *RBSTMonoidNoncommutative) Erase(root *node, start, end uint32) (remain *node, erased *node) {
	x, y, z := rbst._split3(root, start, end)
	remain = rbst.Merge(x, z)
	erased = y
	return
}

func (rbst *RBSTMonoidNoncommutative) Insert(root *node, k uint32, v E) *node {
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

func (rbst *RBSTMonoidNoncommutative) RotateRight(root *node, start, end, k uint32) *node {
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

func (rbst *RBSTMonoidNoncommutative) RotateLeft(root *node, start, end, k uint32) *node {
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

func (rbst *RBSTMonoidNoncommutative) RotateRightAll(root *node, k uint32) *node {
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

func (rbst *RBSTMonoidNoncommutative) RotateLeftAll(root *node, k uint32) *node {
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

func (rbst *RBSTMonoidNoncommutative) _merge3(a, b, c *node) *node {
	return rbst.Merge(rbst.Merge(a, b), c)
}

func (rbst *RBSTMonoidNoncommutative) _merge4(a, b, c, d *node) *node {
	return rbst.Merge(rbst.Merge(rbst.Merge(a, b), c), d)
}

func (rbst *RBSTMonoidNoncommutative) _split3(root *node, l, r uint32) (*node, *node, *node) {
	root, nr := rbst.Split(root, r)
	root, nm := rbst.Split(root, l)
	return root, nm, nr
}

func (rbst *RBSTMonoidNoncommutative) _split4(root *node, i, j, k uint32) (*node, *node, *node, *node) {
	root, d := rbst.Split(root, k)
	a, b, c := rbst._split3(root, i, j)
	return a, b, c, d
}

func (rbst *RBSTMonoidNoncommutative) _copyNode(n *node) *node {
	if n == nil || !rbst.persistent {
		return n
	}
	return &node{l: n.l, r: n.r, v: n.v, sum: n.sum, size: n.size}
}

func (rbst *RBSTMonoidNoncommutative) _nextRand() uint32 {
	t := rbst.x ^ (rbst.x << 11)
	rbst.x, rbst.y, rbst.z = rbst.y, rbst.z, rbst.w
	rbst.w = (rbst.w ^ (rbst.w >> 19)) ^ (t ^ (t >> 8))
	return rbst.w
}

func (rbst *RBSTMonoidNoncommutative) _update(c *node) {
	c.size = 1
	c.sum = c.v
	if left := c.l; left != nil {
		c.size += left.size
		c.sum = rbst.op(left.sum, c.sum)
	}
	if right := c.r; right != nil {
		c.size += right.size
		c.sum = rbst.op(c.sum, right.sum)
	}
}

func (rbst *RBSTMonoidNoncommutative) _mergeRec(lRoot, rRoot *node) *node {
	if lRoot == nil {
		return rRoot
	}
	if rRoot == nil {
		return lRoot
	}
	sl, sr := lRoot.size, rRoot.size
	if rbst._nextRand()%(sl+sr) < sl {
		lRoot = rbst._copyNode(lRoot)
		lRoot.r = rbst._mergeRec(lRoot.r, rRoot)
		rbst._update(lRoot)
		return lRoot
	}
	rRoot = rbst._copyNode(rRoot)
	rRoot.l = rbst._mergeRec(lRoot, rRoot.l)
	rbst._update(rRoot)
	return rRoot
}

func (rbst *RBSTMonoidNoncommutative) _splitRec(root *node, k uint32) (*node, *node) {
	if root == nil {
		return nil, nil
	}
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

func (rbst *RBSTMonoidNoncommutative) _setRec(root *node, k uint32, v E) *node {
	if root == nil {
		return root
	}
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
		root.v = v
		rbst._update(root)
		return root
	}
	root = rbst._copyNode(root)
	root.r = rbst._setRec(root.r, k-(1+sl), v)
	rbst._update(root)
	return root
}

func (rbst *RBSTMonoidNoncommutative) _getRec(root *node, k uint32) E {
	left, right := root.l, root.r
	sl := uint32(0)
	if left != nil {
		sl = left.size
	}
	if k == sl {
		return root.v
	}
	if k < sl {
		return rbst._getRec(left, k)
	}
	return rbst._getRec(right, k-(1+sl))
}

func (rbst *RBSTMonoidNoncommutative) _splitMaxRightRec(root *node, check func(E) bool, x *E) (*node, *node) {
	if root == nil {
		return nil, nil
	}
	root = rbst._copyNode(root)
	y := rbst.op(*x, root.sum)
	if check(y) {
		*x = y
		return root, nil
	}
	left, right := root.l, root.r
	if left != nil {
		y = rbst.op(*x, left.sum)
		if !check(y) {
			n1, n2 := rbst._splitMaxRightRec(left, check, x)
			root.l = n2
			rbst._update(root)
			return n1, root
		}
		*x = y
	}
	y = rbst.op(*x, root.v)
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

func (rbst *RBSTMonoidNoncommutative) _queryRec(root *node, start, end uint32) E {
	if start == 0 && end == root.size {
		return root.sum
	}
	left, right := root.l, root.r
	sl := uint32(0)
	if left != nil {
		sl = left.size
	}
	res := rbst.e
	if start < sl {
		y := rbst._queryRec(left, start, min32(end, sl))
		res = rbst.op(res, y)
	}
	if start <= sl && sl < end {
		res = rbst.op(res, root.v)
	}
	k := 1 + sl
	if k < end {
		y := rbst._queryRec(right, max32(k, start)-k, end-k)
		res = rbst.op(res, y)
	}
	return res
}

func (rbst *RBSTMonoidNoncommutative) _updateRec(root *node, k uint32, x E) *node {
	if root == nil {
		return root
	}
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
		root.v = rbst.op(root.v, x)
		rbst._update(root)
		return root
	}
	root = rbst._copyNode(root)
	root.r = rbst._updateRec(root.r, k-(1+sl), x)
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

func test() {
	swapRange := func(arr []int, start, end int) {
		if start >= end {
			return
		}
		end--
		for start < end {
			arr[start], arr[end] = arr[end], arr[start]
			start++
			end--
		}
	}
	_ = swapRange

	rotateLeft := func(arr []int, start, end, step int) {
		n := end - start
		if n <= 1 || step == 0 {
			return
		}
		if step >= n {
			step %= n
		}
		swapRange(arr, start, start+step)
		swapRange(arr, start+step, end)
		swapRange(arr, start, end)
	}
	_ = rotateLeft

	rotateRight := func(arr []int, start, end, step int) {
		n := end - start
		if n <= 1 || step == 0 {
			return
		}
		if step >= n {
			step %= n
		}
		swapRange(arr, start, end-step)
		swapRange(arr, end-step, end)
		swapRange(arr, start, end)
	}
	_ = rotateRight

	for i := 0; i < 500; i++ {
		nums := make([]int, 500)
		for j := 0; j < 500; j++ {
			nums[j] = rand.Intn(100)
		}
		base := NewDynamicHashStringBase(len(nums), 0)
		hasher := NewDynamicHashStringX(len(nums), func(i int) uint { return uint(nums[i]) }, base)
		for j := 0; j < 500; j++ {
			// Set
			index := rand.Intn(len(nums))
			value := rand.Intn(100)
			hasher.Set(index, uint(value))
			nums[index] = value

			// Get
			start, end := rand.Intn(len(nums)), rand.Intn(len(nums))
			if start > end {
				start, end = end, start
			}
			var sum uint
			for k := start; k < end; k++ {
				sum = base.Concat(sum, uint(nums[k]), 1)
			}
			if sum != hasher.Get(start, end) {
				fmt.Println("Get error")
				panic("Get error")
			}

			// GetAll
			{
				sum = 0
				for k := 0; k < len(nums); k++ {
					sum = base.Concat(sum, uint(nums[k]), 1)
				}
				if sum != hasher.GetAll() {
					fmt.Println("GetAll error")
					panic("GetAll error")
				}
			}

			// Pop
			index = rand.Intn(len(nums))
			if hasher.Pop(index) != uint(nums[index]) {
				fmt.Println("Pop error")
				panic("Pop error")
			}
			nums = append(nums[:index], nums[index+1:]...)

			// Insert
			index = rand.Intn(len(nums))
			value = rand.Intn(100)
			hasher.Insert(index, uint(value))
			nums = append(nums[:index], append([]int{value}, nums[index:]...)...)
		}

		// RotateRight
		{
			start, end := rand.Intn(len(nums)), rand.Intn(len(nums))
			if start > end {
				start, end = end, start
			}
			k := rand.Intn(len(nums))
			hasher.RotateRight(start, end, k)
			rotateRight(nums, start, end, k)
		}

		// RotateLeft
		{
			start, end := rand.Intn(len(nums)), rand.Intn(len(nums))
			if start > end {
				start, end = end, start
			}
			k := rand.Intn(len(nums))
			hasher.RotateLeft(start, end, k)
		}

		// RotateRightAll
		{
			k := rand.Intn(len(nums))
			hasher.RotateRightAll(k)
			rotateRight(nums, 0, len(nums), k)
		}

		// RotateLeftAll
		{
			k := rand.Intn(len(nums))
			hasher.RotateLeftAll(k)
			rotateLeft(nums, 0, len(nums), k)
		}
	}

	fmt.Println("pass")
}

func testTime() {
	n := 200000
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = rand.Intn(100)
	}
	time1 := time.Now()
	base := NewDynamicHashStringBase(len(nums), 0)
	hasher := NewDynamicHashStringX(len(nums), func(i int) uint { return uint(nums[i]) }, base)

	for i := 0; i < n; i++ {
		hasher.Set(i, uint(nums[i]))
		hasher.Get(0, i+1)
		hasher.GetAll()
		hasher.Pop(i)
		hasher.Insert(i, uint(nums[i]))
		hasher.RotateRight(0, i+1, i)
		hasher.RotateLeft(0, i+1, i)
		hasher.RotateRightAll(i)
		hasher.RotateLeftAll(i)
	}
	fmt.Println(time.Since(time1)) // 1.331472166s
}
