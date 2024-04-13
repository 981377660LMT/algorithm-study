// https://www.acwing.com/blog/content/28060/
package main

import (
	"fmt"
	"math"
)

func main() {
	test()
}

type IBlock[S, U, V any] interface {
	// 分裂整块，block1包含前k个元素，block2包含后面的元素.分裂后，可以回收原块内存.
	Split(k int32) (block1 IBlock[S, U, V], block2 IBlock[S, U, V])
	// 合并两个块，合并后，可以回收原块内存.
	Merge(other IBlock[S, U, V]) IBlock[S, U, V]
	// 在index下标之前插入元素.
	InserBefore(index int32, e V)
	Delete(index int32)
	Get(index int32) V
	Reverse()

	FullyQuery(sum S)
	PartialQuery(index int32, sum S)
	FullyUpdate(lazy U)
	PartialUpdate(index int32, lazy U)
	BeforePartialQuery()
	AfterPartialUpdate()
}

// 块状链表.
type BlockChain[S, U, V any] struct {
	Head, Tail *linkedNode[IBlock[S, U, V]]
	b, size    int32
}

// blockSize: 2 * (int32(math.Sqrt(float64(n))) + 1)
func NewBlockChain[S, U, V any](n int32, blockSupplier func(start, end int32) IBlock[S, U, V], blockSize int32) *BlockChain[S, U, V] {
	if blockSize == -1 {
		blockSize = 2 * (int32(math.Sqrt(float64(n))) + 1)
	}
	res := &BlockChain[S, U, V]{
		b:    blockSize,
		size: n,
		Head: &linkedNode[IBlock[S, U, V]]{},
		Tail: &linkedNode[IBlock[S, U, V]]{},
	}
	linkNode(res.Head, res.Tail)
	for start := int32(0); start < n; start += blockSize {
		end := start + blockSize
		if end > n {
			end = n
		}
		block := blockSupplier(start, end)
		node := &linkedNode[IBlock[S, U, V]]{}
		node.data = block
		node.size = end - start
		linkNode(res.Tail.prev, node)
		linkNode(node, res.Tail)
	}
	return res
}

func (bc *BlockChain[S, U, V]) _new2(b int32, supplier func() IBlock[S, U, V]) *BlockChain[S, U, V] {
	res := &BlockChain[S, U, V]{
		b:    b,
		Head: &linkedNode[IBlock[S, U, V]]{},
		Tail: &linkedNode[IBlock[S, U, V]]{},
	}
	block := supplier()
	node := &linkedNode[IBlock[S, U, V]]{}
	node.data = block
	res.b = b
	linkNode(res.Tail.prev, node)
	linkNode(node, res.Tail)
	return res
}

func (bc *BlockChain[S, U, V]) _new3(b int32, begin, end *linkedNode[IBlock[S, U, V]]) *BlockChain[S, U, V] {
	res := &BlockChain[S, U, V]{
		b:    b,
		Head: &linkedNode[IBlock[S, U, V]]{},
		Tail: &linkedNode[IBlock[S, U, V]]{},
	}
	linkNode(res.Head, begin)
	linkNode(end, res.Tail)
	res._maintain()
	return res
}

func (bc *BlockChain[S, U, V]) Get(index int32) V {
	if index < 0 {
		index += bc.size
	}
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		if node.size <= index {
			index -= node.size
			continue
		}
		node.data.BeforePartialQuery()
		return node.data.Get(index)
	}
	panic("Index out of bounds")
}

func (bc *BlockChain[S, U, V]) PrefixSize(block IBlock[S, U, V], include bool) int32 {
	res := int32(0)
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		if node.data == block {
			if include {
				res += node.size
			}
			break
		}
		res += node.size
	}
	return res
}

func (bc *BlockChain[S, U, V]) Split(k int32, blockSupplier func() IBlock[S, U, V]) (first, second *BlockChain[S, U, V]) {
	k++
	if k == 1 {
		return bc._new2(bc.b, blockSupplier), bc
	}
	if k > bc.size {
		return bc, bc._new2(bc.b, blockSupplier)
	}
	head := bc._splitKth(k)
	end := bc.Tail.prev
	linkNode(head.prev, bc.Tail)
	b := bc._new3(bc.b, head, end)
	bc._maintain()
	return bc, b
}

func (bc *BlockChain[S, U, V]) MergeDestructively(other *BlockChain[S, U, V]) *BlockChain[S, U, V] {
	linkNode(bc.Tail.prev, other.Head.next)
	bc.Tail = other.Tail
	bc.size += other.size
	return bc
}

func (bc *BlockChain[S, U, V]) InsertBefore(index int32, e V) {
	if index < 0 {
		index += bc.size
	}
	if index < 0 {
		index = 0
	}
	if index > bc.size {
		index = bc.size
	}
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		if node.size < index {
			index -= node.size
			continue
		}
		node.data.InserBefore(index, e)
		node.size++
		break
	}
	bc._maintain()
}

func (bc *BlockChain[S, U, V]) Delete(index int32) {
	if index < 0 {
		index += bc.size
	}
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		if node.size <= index {
			index -= node.size
			continue
		}
		node.data.Delete(index)
		node.size--
		break
	}
	bc._maintain()
}

func (bc *BlockChain[S, U, V]) Update(start, end int32, update U) {
	if start < 0 {
		start += bc.size
	}
	if end > bc.size {
		end = bc.size
	}
	if start >= end {
		return
	}
	end--
	offset := int32(0)
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		left := offset
		right := offset + node.size - 1
		offset += node.size
		if bc._enter(start, end, left, right) {
			node.data.FullyUpdate(update)
		} else if bc._leave(start, end, left, right) {
			continue
		} else {
			for i := max32(left, start); i <= min32(right, end); i++ {
				node.data.PartialUpdate(i-left, update)
			}
			node.data.AfterPartialUpdate()
		}
	}
}

func (bc *BlockChain[S, U, V]) Query(start, end int32, sum S) {
	if start < 0 {
		start += bc.size
	}
	if end > bc.size {
		end = bc.size
	}
	if start >= end {
		return
	}
	end--
	offset := int32(0)
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		left := offset
		right := offset + node.size - 1
		offset += node.size
		if bc._enter(start, end, left, right) {
			node.data.FullyQuery(sum)
		} else if bc._leave(start, end, left, right) {
			continue
		} else {
			node.data.BeforePartialQuery()
			for i := max32(left, start); i <= min32(right, end); i++ {
				node.data.PartialQuery(i-left, sum)
			}
		}
	}
}

func (bc *BlockChain[S, U, V]) RotateLeft(k int32) {
	if k < 0 {
		k += bc.size
	}
	if k >= bc.size {
		k %= bc.size
	}
	if k == 0 {
		return
	}
	k++
	node := bc._splitKth(k)
	h1 := bc.Head.next
	e1 := node.prev
	h2 := node
	e2 := bc.Tail.prev
	linkNode(bc.Head, h2)
	linkNode(e2, h1)
	linkNode(e1, bc.Tail)
	bc._maintain()
}

func (bc *BlockChain[S, U, V]) Reverse(start, end int32) {
	if start >= end {
		return
	}
	end--
	left := bc._splitKth(start + 1)
	right := bc._splitKth(end + 2).prev
	begin := left.prev
	endNode := right.next
	right.next = nil
	bc._reverse(left, nil)
	linkNode(begin, right)
	linkNode(left, endNode)
	bc._maintain()
}

func (bc *BlockChain[S, U, V]) GetAll() []V {
	res := make([]V, 0, bc.size)
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		for i := int32(0); i < node.size; i++ {
			res = append(res, node.data.Get(i))
		}
	}
	return res
}

func (bc *BlockChain[S, U, V]) EnumerateBlock(f func(block IBlock[S, U, V])) {
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		f(node.data)
	}
}

func (bc *BlockChain[S, U, V]) Len() int32 {
	return bc.size
}

func (bc *BlockChain[S, U, V]) _reverse(root *linkedNode[IBlock[S, U, V]], p *linkedNode[IBlock[S, U, V]]) {
	if root == nil {
		return
	}
	bc._reverse(root.next, root)
	root.data.Reverse()
	root.prev = root.next
	root.next = p
}

func (bc *BlockChain[S, U, V]) _split(node *linkedNode[IBlock[S, U, V]], k int32) {
	post := &linkedNode[IBlock[S, U, V]]{}
	block1, block2 := node.data.Split(k)
	linkNode(post, node.next)
	linkNode(node, post)
	post.data = block2
	post.size = node.size - k
	node.data = block1
	node.size = k
}

func (bc *BlockChain[S, U, V]) _mergeNode(a, b *linkedNode[IBlock[S, U, V]]) {
	linkNode(a, b.next)
	a.data = a.data.Merge(b.data)
	a.size += b.size
}

func (bc *BlockChain[S, U, V]) _maintain() {
	bc.size = 0
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		bc.size += node.size
		if node.size >= 2*bc.b {
			bc._split(node, bc.b)
		} else if node.prev != bc.Head && node.size+node.prev.size <= bc.b {
			bc._mergeNode(node.prev, node)
		}
	}
}

func (bc *BlockChain[S, U, V]) _splitKth(k int32) *linkedNode[IBlock[S, U, V]] {
	for node := bc.Head.next; node != bc.Tail; node = node.next {
		if node.size < k {
			k -= node.size
			continue
		}
		if k != 1 {
			bc._split(node, k-1)
			node = node.next
		}
		return node
	}
	return bc.Tail
}

func (bc *BlockChain[S, U, V]) _enter(L, R, l, r int32) bool {
	return L <= l && r <= R
}

func (bc *BlockChain[S, U, V]) _leave(L, R, l, r int32) bool {
	return l > R || r < L
}

type linkedNode[V any] struct {
	prev, next *linkedNode[V]
	size       int32
	data       V
}

func (node *linkedNode[V]) String() string {
	return fmt.Sprint(node.data)
}

func linkNode[V any](a, b *linkedNode[V]) {
	b.prev = a
	a.next = b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// application
type INumberArrayBlock = IBlock[*int, int, int]
type NumberArrayBlock struct {
	data []int
	sum  int
}

func NewNumberArrayBlock(data []int) *NumberArrayBlock {
	res := &NumberArrayBlock{data: data}
	for _, v := range data {
		res.sum += v
	}
	return res
}

func (nab *NumberArrayBlock) Split(k int32) (block1 INumberArrayBlock, block2 INumberArrayBlock) {
	a, b := nab.data[:k], nab.data[k:]
	return NewNumberArrayBlock(a), NewNumberArrayBlock(b)
}

func (nab *NumberArrayBlock) Merge(other INumberArrayBlock) INumberArrayBlock {
	nab.data = append(nab.data, other.(*NumberArrayBlock).data...)
	for _, v := range other.(*NumberArrayBlock).data {
		nab.sum += v
	}
	return nab
}

func (nab *NumberArrayBlock) InserBefore(index int32, e int) {
	nab.data = append(nab.data[:index], append([]int{e}, nab.data[index:]...)...)
	nab.sum += e
}

func (nab *NumberArrayBlock) Delete(index int32) {
	nab.sum -= nab.data[index]
	nab.data = append(nab.data[:index], nab.data[index+1:]...)
}

func (nab *NumberArrayBlock) Get(index int32) int {
	return nab.data[index]
}

func (nab *NumberArrayBlock) Reverse() {
	for i, j := 0, len(nab.data)-1; i < j; i, j = i+1, j-1 {
		nab.data[i], nab.data[j] = nab.data[j], nab.data[i]
	}
}

func (nab *NumberArrayBlock) FullyQuery(sum *int) {
	*sum += nab.sum
}

func (nab *NumberArrayBlock) PartialQuery(index int32, sum *int) {
	*sum += nab.data[index]
}

func (nab *NumberArrayBlock) FullyUpdate(lazy int) {}

func (nab *NumberArrayBlock) PartialUpdate(index int32, lazy int) {}

func (nab *NumberArrayBlock) BeforePartialQuery() {}

func (nab *NumberArrayBlock) AfterPartialUpdate() {}

func test() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	blockChain := NewBlockChain(
		10,
		func(start, end int32) IBlock[*int, int, int] {
			return NewNumberArrayBlock(nums[start:end])
		},
		-1,
	)

	sum := new(int)
	blockChain.Query(0, 9, sum)
	fmt.Println(*sum)
}

// 406. 根据身高重建队列
// https://leetcode.cn/problems/queue-reconstruction-by-height/
