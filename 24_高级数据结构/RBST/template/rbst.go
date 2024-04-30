// 构建api:
//  1. Build(n, f) -> root
//  2. NewRoot() -> root
//  3. NewNode(v) -> node
//
// 分裂/拼接api:
//  1. Merge(left, right) -> root
//  2. Split(root, k) -> [0,k) and [k,n)
//
// 查询/更新api:
//
//  1. Set(root, k, v) -> root
//  2. Get(root, k) -> v
//  3. GetAll(root) -> []v
//  4. SplitMaxRight(root,check) -> left,right
//
// 操作api:
//  1. Reverse(root, start, end) -> root
//  2. ReverseAll(root) -> root
//  3. Size(root) -> size
//  4. CopyWithin(root, target, start, end) -> root (持久化为true时)
//  5. Pop(root, k) -> root, v
//  6. Erase(root, start, end) -> remain, erased
//  7. Insert(root, k, v) -> root
//  8. RotateRight(root, start, end, k) -> root
//  9. RotateLeft(root, start, end, k) -> root
//  10. RotateRightAll(root, k) -> root
//  11. RotateLeftAll(root, k) -> root

// !因为支持可持久化，所有修改操作都必须返回新的root.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	test()
	testTime()
}

type E = int
type node struct {
	l, r *node
	v    E
	size uint32
	rev  uint8
}

type RBST struct {
	persistent bool
	x, y, z, w uint32
}

func NewRBST(persistent bool) *RBST {
	return &RBST{
		persistent: persistent,
		x:          123456789,
		y:          362436069,
		z:          521288629,
		w:          88675123,
	}
}

func (rbst *RBST) Build(n uint32, f func(uint32) E) *node {
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

func (rbst *RBST) NewRoot() *node { return nil }

func (rbst *RBST) NewNode(v E) *node { return &node{v: v, size: 1} }

func (rbst *RBST) Merge(lRoot, rRoot *node) *node            { return rbst._mergeRec(lRoot, rRoot) }
func (rbst *RBST) Split(root *node, k uint32) (*node, *node) { return rbst._splitRec(root, k) }

func (rbst *RBST) Set(root *node, k uint32, v E) *node { return rbst._setRec(root, k, v) }
func (rbst *RBST) Get(root *node, k uint32) E          { return rbst._getRec(root, k, 0) }
func (rbst *RBST) GetAll(root *node) []E {
	res := make([]E, 0, rbst.Size(root))
	var dfs func(root *node, rev uint8)
	dfs = func(root *node, rev uint8) {
		if root == nil {
			return
		}
		if rev == 1 {
			dfs(root.r, rev^root.rev)
			res = append(res, root.v)
			dfs(root.l, rev^root.rev)
		} else {
			dfs(root.l, rev^root.rev)
			res = append(res, root.v)
			dfs(root.r, rev^root.rev)
		}
	}
	dfs(root, 0)
	return res
}

func (rbst *RBST) SplitMaxRight(root *node, check func(E) bool) (*node, *node) {
	return rbst._splitMaxRightRec(root, check)
}

func (rbst *RBST) Reverse(root *node, start, end uint32) *node {
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

func (rbst *RBST) ReverseAll(root *node) *node {
	if root == nil {
		return nil
	}
	root = rbst._copyNode(root)
	root.rev ^= 1
	root.l, root.r = root.r, root.l
	return root
}

func (rbst *RBST) Size(root *node) uint32 {
	if root == nil {
		return 0
	}
	return root.size
}

func (rbst *RBST) CopyWithin(root *node, target, start, end uint32) *node {
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

func (rbst *RBST) Pop(root *node, k uint32) (*node, E) {
	n := rbst.Size(root)
	if k < 0 {
		k += n
	}
	x, y, z := rbst._split3(root, k, k+1)
	res := y.v
	newRoot := rbst.Merge(x, z)
	return newRoot, res
}

func (rbst *RBST) Erase(root *node, start, end uint32) (remain *node, erased *node) {
	x, y, z := rbst._split3(root, start, end)
	remain = rbst.Merge(x, z)
	erased = y
	return
}

func (rbst *RBST) Insert(root *node, k uint32, v E) *node {
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

func (rbst *RBST) RotateRight(root *node, start, end, k uint32) *node {
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

func (rbst *RBST) RotateLeft(root *node, start, end, k uint32) *node {
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

func (rbst *RBST) RotateRightAll(root *node, k uint32) *node {
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

func (rbst *RBST) RotateLeftAll(root *node, k uint32) *node {
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

func (rbst *RBST) _merge3(a, b, c *node) *node {
	return rbst.Merge(rbst.Merge(a, b), c)
}

func (rbst *RBST) _merge4(a, b, c, d *node) *node {
	return rbst.Merge(rbst.Merge(rbst.Merge(a, b), c), d)
}

func (rbst *RBST) _split3(root *node, l, r uint32) (*node, *node, *node) {
	root, nr := rbst.Split(root, r)
	root, nm := rbst.Split(root, l)
	return root, nm, nr
}

func (rbst *RBST) _split4(root *node, i, j, k uint32) (*node, *node, *node, *node) {
	root, d := rbst.Split(root, k)
	a, b, c := rbst._split3(root, i, j)
	return a, b, c, d
}

func (rbst *RBST) _copyNode(n *node) *node {
	if n == nil || !rbst.persistent {
		return n
	}
	return &node{l: n.l, r: n.r, v: n.v, size: n.size, rev: n.rev}
}

func (rbst *RBST) _nextRand() uint32 {
	t := rbst.x ^ (rbst.x << 11)
	rbst.x, rbst.y, rbst.z = rbst.y, rbst.z, rbst.w
	rbst.w = (rbst.w ^ (rbst.w >> 19)) ^ (t ^ (t >> 8))
	return rbst.w
}

func (rbst *RBST) _propagate(c *node) {
	if c.rev == 1 {
		if c.l != nil {
			c.l = rbst._copyNode(c.l)
			c.l.rev ^= 1
			c.l.l, c.l.r = c.l.r, c.l.l
		}
		if c.r != nil {
			c.r = rbst._copyNode(c.r)
			c.r.rev ^= 1
			c.r.l, c.r.r = c.r.r, c.r.l
		}
		c.rev = 0
	}
}

func (rbst *RBST) _update(c *node) {
	c.size = 1
	if c.l != nil {
		c.size += c.l.size
	}
	if c.r != nil {
		c.size += c.r.size
	}
}

func (rbst *RBST) _mergeRec(lRoot, rRoot *node) *node {
	if lRoot == nil {
		return rRoot
	}
	if rRoot == nil {
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

func (rbst *RBST) _splitRec(root *node, k uint32) (*node, *node) {
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

func (rbst *RBST) _setRec(root *node, k uint32, v E) *node {
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
		root.v = v
		rbst._update(root)
		return root
	}
	root = rbst._copyNode(root)
	root.r = rbst._setRec(root.r, k-(1+sl), v)
	rbst._update(root)
	return root
}

func (rbst *RBST) _getRec(root *node, k uint32, rev uint8) E {
	if rev == 1 {
		left, right := root.r, root.l
		sl := uint32(0)
		if left != nil {
			sl = left.size
		}
		if k == sl {
			return root.v
		}
		rev ^= root.rev
		if k < sl {
			return rbst._getRec(left, k, rev)
		}
		return rbst._getRec(right, k-(1+sl), rev)
	} else {
		left, right := root.l, root.r
		sl := uint32(0)
		if left != nil {
			sl = left.size
		}
		if k == sl {
			return root.v
		}
		rev ^= root.rev
		if k < sl {
			return rbst._getRec(left, k, rev)
		}
		return rbst._getRec(right, k-(1+sl), rev)
	}
}

func (rbst *RBST) _splitMaxRightRec(root *node, check func(E) bool) (*node, *node) {
	if root == nil {
		return nil, nil
	}
	rbst._propagate(root)
	root = rbst._copyNode(root)
	if check(root.v) {
		n1, n2 := rbst._splitMaxRightRec(root.r, check)
		root.r = n1
		rbst._update(root)
		return root, n2
	}
	n1, n2 := rbst._splitMaxRightRec(root.l, check)
	root.l = n2
	rbst._update(root)
	return n1, root
}

func test() {
	for i := 0; i < 10; i++ {
		n := uint32(rand.Intn(1000)) + 500
		nums := make([]int, n)
		for i := 0; i < int(n); i++ {
			nums[i] = rand.Intn(100)
		}
		rbst := NewRBST(true)
		root := rbst.Build(n, func(i uint32) E { return E(nums[i]) })

		for j := 0; j < 500; j++ {
			// Get
			{

				k := rand.Intn(int(len(nums)))
				if rbst.Get(root, uint32(k)) != E(nums[k]) {
					fmt.Println("Get Error")
					panic("Get Error")
				}
			}

			// Set
			{

				k := rand.Intn(int(len(nums)))
				v := rand.Intn(100)
				root = rbst.Set(root, uint32(k), E(v))
				nums[k] = v
			}

			// GetAll
			{
				res1 := rbst.GetAll(root)
				res2 := make([]int, len(nums))
				copy(res2, nums)
				if len(res1) != len(res2) {
					fmt.Println("GetAll Error")
					panic("GetAll Error")
				}
				for i := 0; i < len(res1); i++ {
					if res1[i] != E(res2[i]) {
						fmt.Println("GetAll Error")
						panic("GetAll Error")
					}
				}
			}

			// Reverse
			{
				l, r := rand.Intn(len(nums)), rand.Intn(len(nums))
				if l > r {
					l, r = r, l
				}
				root = rbst.Reverse(root, uint32(l), uint32(r))
				for i, j := l, r-1; i < j; i, j = i+1, j-1 {
					nums[i], nums[j] = nums[j], nums[i]
				}
			}

			// ReverseAll
			{
				root = rbst.ReverseAll(root)
				for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
					nums[i], nums[j] = nums[j], nums[i]
				}
			}

			// CopyWithin
			{
				target, start, end := rand.Intn(len(nums)), rand.Intn(len(nums)), rand.Intn(len(nums))
				if start > end {
					start, end = end, start
				}
				if target+end-start <= len(nums) {
					root = rbst.CopyWithin(root, uint32(target), uint32(start), uint32(end))
					copy(nums[target:target+end-start], nums[start:end])
				}
			}

			// Pop
			{
				k := rand.Intn(len(nums))
				newRoot, v := rbst.Pop(root, uint32(k))
				root = newRoot
				if v != E(nums[k]) {
					fmt.Println("Pop Error", v, nums[k])
					panic("Pop Error")
				}
				nums = append(nums[:k], nums[k+1:]...)
			}

			// Insert
			{
				k := rand.Intn(len(nums))
				v := rand.Intn(100)
				root = rbst.Insert(root, uint32(k), E(v))
				nums = append(nums[:k], append([]int{v}, nums[k:]...)...)
			}

			// RotateRight
			{
				start, end, k := rand.Intn(len(nums)), rand.Intn(len(nums)), rand.Intn(len(nums))
				if start > end {
					start, end = end, start
				}
				root = rbst.RotateRight(root, uint32(start), uint32(end), uint32(k))
				rotateRight(nums, start, end, k)
			}

			// RotateLeft
			{
				start, end, k := rand.Intn(len(nums)), rand.Intn(len(nums)), rand.Intn(len(nums))
				if start > end {
					start, end = end, start
				}
				root = rbst.RotateLeft(root, uint32(start), uint32(end), uint32(k))
				rotateLeft(nums, start, end, k)
			}

			// RotateRightAll
			{
				k := rand.Intn(len(nums))
				root = rbst.RotateRightAll(root, uint32(k))
				rotateRight(nums, 0, len(nums), k)
			}

			// RotateLeftAll
			{
				k := rand.Intn(len(nums))
				root = rbst.RotateLeftAll(root, uint32(k))
				rotateLeft(nums, 0, len(nums), k)
			}
		}
	}
	fmt.Println("Pass")
}

func testTime() {
	n := uint32(2e5)
	nums := make([]int, n)
	for i := 0; i < int(n); i++ {
		nums[i] = rand.Intn(100)
	}
	rbst := NewRBST(false)
	root := rbst.Build(n, func(i uint32) E { return E(nums[i]) })

	time1 := time.Now()
	for j := uint32(0); j < n; j++ {
		root = rbst.Set(root, j, int(j))
		rbst.Get(root, j)
		root = rbst.Reverse(root, 0, j)
		root = rbst.ReverseAll(root)
		root, _ = rbst.Pop(root, j)
		root = rbst.Insert(root, j, int(j))
		root = rbst.RotateRight(root, 0, j, j)
		root = rbst.RotateLeft(root, 0, j, j)
		root = rbst.RotateRightAll(root, j)
		root = rbst.RotateLeftAll(root, j)
	}
	fmt.Println("Time1:", time.Since(time1)) // Time1: 637.91425ms
}

func swapRange(arr []int, start, end int) {
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

func rotateLeft(arr []int, start, end, step int) {
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

func rotateRight(arr []int, start, end, step int) {
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
