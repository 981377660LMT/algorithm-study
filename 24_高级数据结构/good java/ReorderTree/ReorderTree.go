// 维护一个多重集合，支持以下操作：
// 1. 插入元素.
// 2. 弹出最大的 k 个元素，如果其中有元素<=0，则操作失败，否则将它们减一后重新插入.

package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

func init() {
	debug.SetGCPercent(-1)
}

// 2790. 长度递增组的最大数目
// https://leetcode.cn/problems/maximum-number-of-groups-with-increasing-length
// 给你一个下标从 0 开始、长度为 n 的数组 usageLimits 。
// 你的任务是使用从 0 到 n - 1 的数字创建若干组，并确保每个数字 i 在 所有组 中使用的次数总共不超过 usageLimits[i] 次。
// 此外，还必须满足以下条件：
// 每个组必须由 不同 的数字组成，也就是说，单个组内不能存在重复的数字。
// 每个组（除了第一个）的长度必须 严格大于 前一个组。
// 在满足所有条件的情况下，以整数形式返回可以创建的最大组数。
func maxIncreasingGroups(usageLimits []int) int {
	n := int32(len(usageLimits))
	tree := NewReorderTree()
	for _, limit := range usageLimits {
		tree.Add(int32(limit))
	}
	for size := int32(1); size <= n; size++ {
		ok := tree.ReduceTopK(size)
		if !ok {
			return int(size - 1) // rollback
		}
	}
	return int(n)
}

func demo() {
	container := NewReorderTree()
	container.Add(1)
	container.Add(2)
	container.Add(3)
	container.Add(4)

	container.Enumerate(func(x int32) { fmt.Println(x) })
	fmt.Println(container.ReduceTopK(4))
	fmt.Println(container.ReduceTopK(4))

	container.Enumerate(func(x int32) { fmt.Println(x) })
}

type ReorderTree struct {
	root *Treap
}

func NewReorderTree() *ReorderTree {
	return &ReorderTree{root: NIL}
}

func (t *ReorderTree) Add(x int32) {
	s0, s1 := SplitByKey(t.root, x)
	node := NewTreap(x)
	node.PushUp()
	t.root = Merge(s0, node)
	t.root = Merge(t.root, s1)
}

func (t *ReorderTree) ReduceTopK(k int32) bool {
	if t.root.size < k {
		return false
	}
	s0, s1 := SplitByRank(t.root, t.root.size-k)
	minKey := GetKeyByRank(s1, 1)
	if minKey <= 0 {
		t.root = Merge(s0, s1)
		return false
	}
	s2, s3 := SplitByKey(s0, minKey-1)
	s1.Propagate(-1)
	s4, s5 := SplitByKey(s1, minKey-1)
	s0 = Merge(s2, s4)
	s1 = Merge(s3, s5)
	t.root = Merge(s0, s1)
	return true
}

func (t *ReorderTree) Enumerate(f func(int32)) {
	Enumerate(t.root, f)
}

func (t *ReorderTree) GetAll() []int32 {
	res := make([]int32, 0, t.root.size)
	t.Enumerate(func(x int32) { res = append(res, x) })
	return res
}

func (t *ReorderTree) Size() int32 {
	return t.root.size
}

func createNIL() *Treap {
	res := &Treap{}
	res.left, res.right = res, res
	return res
}

var random = NewRandom()
var NIL = createNIL()

type Treap struct {
	left, right *Treap
	size        int32
	key         int32
	lazy        int32
}

func NewTreap(value int32) *Treap {
	return &Treap{left: NIL, right: NIL, size: 1, key: value}
}

func (p *Treap) Clone() *Treap {
	if p == NIL {
		return p
	}
	return &Treap{left: p.left, right: p.right, size: p.size, key: p.key, lazy: p.lazy}
}

func (p *Treap) PushUp() {
	if p == NIL {
		return
	}
	p.size = p.left.size + p.right.size + 1
}

func (p *Treap) PushDown() {
	if p == NIL {
		return
	}
	if p.lazy != 0 {
		p.left.Propagate(p.lazy)
		p.right.Propagate(p.lazy)
		p.lazy = 0
	}
}

func (p *Treap) Propagate(x int32) {
	p.lazy += x
	p.key += x
}

func Build(f func(i int32) *Treap, l, r int32) *Treap {
	if l > r {
		return NIL
	}
	mid := (l + r) >> 1
	root := f(mid)
	root.left = Build(f, l, mid-1)
	root.right = Build(f, mid+1, r)
	root.PushUp()
	return root
}

func SplitByRank(root *Treap, rank int32) (*Treap, *Treap) {
	if root == NIL {
		return NIL, NIL
	}
	root.PushDown()
	var res0, res1 *Treap
	if root.left.size >= rank {
		res0, res1 = SplitByRank(root.left, rank)
		root.left = res1
		res1 = root
	} else {
		res0, res1 = SplitByRank(root.right, rank-(root.size-root.right.size))
		root.right = res0
		res0 = root
	}
	root.PushUp()
	return res0, res1
}

func Merge(a, b *Treap) *Treap {
	if a == NIL {
		return b
	}
	if b == NIL {
		return a
	}
	if int(random.Rng()*(uint64(a.size)+uint64(b.size))>>32) < int(a.size) {
		a.PushDown()
		a.right = Merge(a.right, b)
		a.PushUp()
		return a
	} else {
		b.PushDown()
		b.left = Merge(a, b.left)
		b.PushUp()
		return b
	}
}

func GetValueByRank(root *Treap, k int32) int32 {
	for root.size > 1 {
		if root.left.size >= k {
			root = root.left
		} else {
			k -= root.left.size
			if k == 1 {
				break
			}
			k--
			root = root.right
		}
	}
	return root.key
}

// <= key, > key
func SplitByKey(root *Treap, key int32) (*Treap, *Treap) {
	if root == NIL {
		return NIL, NIL
	}
	root.PushDown()
	var res0, res1 *Treap
	if root.key > key {
		res0, res1 = SplitByKey(root.left, key)
		root.left = res1
		res1 = root
	} else {
		res0, res1 = SplitByKey(root.right, key)
		root.right = res0
		res0 = root
	}
	root.PushUp()
	return res0, res1
}

func GetKeyByRank(treap *Treap, k int32) int32 {
	for treap.size > 1 {
		treap.PushDown()
		if treap.left.size >= k {
			treap = treap.left
		} else {
			k -= treap.left.size
			if k == 1 {
				break
			}
			k--
			treap = treap.right
		}
	}
	return treap.key
}

func Clone(root *Treap) *Treap {
	if root == NIL {
		return NIL
	}
	clone := root.Clone()
	clone.left = Clone(root.left)
	clone.right = Clone(root.right)
	return clone
}

func Enumerate(root *Treap, f func(int32)) {
	if root == NIL {
		return
	}
	root.PushDown()
	Enumerate(root.left, f)
	f(root.key)
	Enumerate(root.right, f)
}

type Random struct {
	seed     uint64
	hashBase uint64
}

func NewRandom() *Random                 { return &Random{seed: uint64(time.Now().UnixNano()/2 + 1)} }
func NewRandomWithSeed(seed int) *Random { return &Random{seed: uint64(seed)} }

func (r *Random) Rng() uint64 {
	r.seed ^= r.seed << 7
	r.seed ^= r.seed >> 9
	return r.seed & 0xFFFFFFFF
}
