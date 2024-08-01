// 可撤销懒线段树.
// RollbackableSegmentTreeLazy/SegmentTreeLazyRollbackable
// !新增了 GetTime 和 Rollback 方法.

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	seg := NewLazySegTreeRollbackable(10, func(index int32) E { return 0 })
	fmt.Println(seg.GetAll())
	seg.Update(0, 5, 2)
	fmt.Println(seg.GetAll())
	fmt.Println(seg.Query(0, 5))
	time1, time2 := seg.GetTime()
	seg.Update(0, 5, 3)
	fmt.Println(seg.GetAll())
	seg.Rollback(time1, time2)
	fmt.Println(seg.GetAll())
}

// RangeAddRangeSum

type E = int
type Id = int

func (*LazySegTreeRollbackable) e() E   { return 0 }
func (*LazySegTreeRollbackable) id() Id { return 0 }
func (*LazySegTreeRollbackable) op(left, right E) E {
	return left + right
}
func (*LazySegTreeRollbackable) mapping(f Id, g E, size int) E {
	return f*size + g
}
func (*LazySegTreeRollbackable) composition(f, g Id) Id {
	return f + g
}
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func max32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

// !template
type LazySegTreeRollbackable struct {
	n    int32
	size int32
	log  int32
	data *rollbackArray[E]
	lazy *rollbackArray[Id]
}

func NewLazySegTreeRollbackable(n int32, f func(int32) E) *LazySegTreeRollbackable {
	tree := &LazySegTreeRollbackable{}
	tree.n = n
	tree.log = 1
	for 1<<tree.log < n {
		tree.log++
	}
	tree.size = 1 << tree.log
	data := make([]E, tree.size<<1)
	for i := range data {
		data[i] = tree.e()
	}
	for i := int32(0); i < n; i++ {
		data[tree.size+i] = f(i)
	}
	// pushUp
	for i := tree.size - 1; i >= 1; i-- {
		data[i] = tree.op(data[i<<1], data[i<<1|1])
	}
	tree.data = newRollbackArrayFrom(data)
	tree.lazy = newRollbackArray(tree.size, func(int32) Id { return tree.id() })
	return tree
}

func NewLazySegTreeRollbackableFrom(leaves []E) *LazySegTreeRollbackable {
	return NewLazySegTreeRollbackable(int32(len(leaves)), func(i int32) E { return leaves[i] })
}

func (tree *LazySegTreeRollbackable) GetTime() (dataTime, lazyTime int32) {
	return tree.data.GetTime(), tree.lazy.GetTime()
}

func (tree *LazySegTreeRollbackable) Rollback(dataTime, lazyTime int32) {
	tree.data.Rollback(dataTime)
	tree.lazy.Rollback(lazyTime)
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTreeRollbackable) Query(left, right int32) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data.Get(left))
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data.Get(right), smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}
func (tree *LazySegTreeRollbackable) QueryAll() E {
	return tree.data.Get(1)
}
func (tree *LazySegTreeRollbackable) GetAll() []E {
	for k := int32(1); k < tree.size; k++ {
		tree.pushDown(k)
	}
	res := tree.data.GetAllMut()[tree.size : tree.size+tree.n]
	return append(res[:0:0], res...)
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTreeRollbackable) Update(left, right int32, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := int32(1); i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTreeRollbackable) MinLeft(right int32, predicate func(data E) bool) int32 {
	if right == 0 {
		return 0
	}
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown((right - 1) >> i)
	}
	res := tree.e()
	for {
		right--
		for right > 1 && right&1 != 0 {
			right >>= 1
		}
		if !predicate(tree.op(tree.data.Get(right), res)) {
			for right < tree.size {
				tree.pushDown(right)
				right = right<<1 | 1
				if tmp := tree.op(tree.data.Get(right), res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - tree.size
		}
		res = tree.op(tree.data.Get(right), res)
		if (right & -right) == right {
			break
		}
	}
	return 0
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTreeRollbackable) MaxRight(left int32, predicate func(data E) bool) int32 {
	if left == tree.n {
		return tree.n
	}
	left += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}
	res := tree.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data.Get(left))) {
			for left < tree.size {
				tree.pushDown(left)
				left <<= 1
				if tmp := tree.op(res, tree.data.Get(left)); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - tree.size
		}
		res = tree.op(res, tree.data.Get(left))
		left++
		if (left & -left) == left {
			break
		}
	}
	return tree.n
}

func (tree *LazySegTreeRollbackable) Get(index int32) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data.Get(index)
}

func (tree *LazySegTreeRollbackable) Set(index int32, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data.Set(index, e)
	for i := int32(1); i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTreeRollbackable) Multiply(index int32, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data.Set(index, tree.op(tree.data.Get(index), e))
	for i := int32(1); i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTreeRollbackable) pushUp(root int32) {
	tree.data.Set(root, tree.op(tree.data.Get(root<<1), tree.data.Get(root<<1|1)))
}

func (tree *LazySegTreeRollbackable) pushDown(root int32) {
	if tmp := tree.lazy.Get(root); tmp != tree.id() {
		tree.propagate(root<<1, tmp)
		tree.propagate(root<<1|1, tmp)
		tree.lazy.Set(root, tree.id())
	}
}
func (tree *LazySegTreeRollbackable) propagate(root int32, f Id) {
	size := 1 << (tree.log - int32((bits.Len32(uint32(root)) - 1)) /**topbit**/)
	tree.data.Set(root, tree.mapping(f, tree.data.Get(root), size))
	if root < tree.size {
		tree.lazy.Set(root, tree.composition(f, tree.lazy.Get(root)))
	}
}

type historyItem[V any] struct {
	index int32
	value V
}

type rollbackArray[V any] struct {
	n       int32
	data    []V
	history []historyItem[V]
}

func newRollbackArray[V any](n int32, f func(index int32) V) *rollbackArray[V] {
	data := make([]V, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &rollbackArray[V]{
		n:    n,
		data: data,
	}
}

func newRollbackArrayFrom[V any](data []V) *rollbackArray[V] {
	return &rollbackArray[V]{n: int32(len(data)), data: data}
}

func (r *rollbackArray[V]) GetTime() int32 {
	return int32(len(r.history))
}

func (r *rollbackArray[V]) Rollback(time int32) {
	for int32(len(r.history)) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair.index] = pair.value
	}
}

func (r *rollbackArray[V]) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *rollbackArray[V]) Get(index int32) V {
	return r.data[index]
}

func (r *rollbackArray[V]) Set(index int32, value V) {
	r.history = append(r.history, historyItem[V]{index: index, value: r.data[index]})
	r.data[index] = value
}

func (r *rollbackArray[V]) GetAll() []V {
	return append(r.data[:0:0], r.data...)
}

func (r *rollbackArray[V]) GetAllMut() []V {
	return r.data
}

func (r *rollbackArray[V]) Len() int32 {
	return r.n
}
