// 复杂度有问题，不要使用

// Package blockchain is a faithful – yet idiomatic – Go translation of the
// Java “BlockChain” framework you provided.  It keeps the same O(N/B) asymptotic
// guarantees, but embraces Go conventions:
//
//	• The Java Pair<T,U> has been replaced with ordinary multiple return values.
//	• Methods that can safely be no‑ops (Reverse/AfterPartialUpdate/…)
//	  are optional – implementers may leave them empty.
//	• The interface hierarchy is flattened: a single generic Block interface
//	  is used everywhere, so callers never need to know the concrete block type.
//
// The code targets Go 1.22+ (for type parameters).
package blockchain

import "fmt"

// ---------- helper utilities -------------------------------------------------

// enter reports whether segment [l,r] is entirely covered by [L,R].
func enter(L, R, l, r int) bool { return L <= l && R >= r }

// leave reports whether segment [l,r] is completely outside [L,R].
func leave(L, R, l, r int) bool { return L > r || R < l }

// ---------- the Block interface ---------------------------------------------

// Block is the building‑block interface every concrete block must satisfy.
// E – aggregator type used by queries (usually *YourSumStruct)
// Id – “update” object (lazy‑tag, delta, …)
// V – element stored inside the block (int, struct{}, …)
type Block[E, Id, V any] interface {
	// Split the block after n elements; return (left,right).
	Split(n int) (left, right Block[E, Id, V])

	// Merge merges b into the receiver and returns the merged block.
	Merge(b Block[E, Id, V]) Block[E, Id, V]

	// Optional point operations (no‑ops by default).
	Insert(index int, e V) // 0 means before first element
	Delete(index int)
	Get(index int) V
	Reverse()

	// Range ⭢ aggregator / updater hooks.
	FullyQuery(sum E)
	PartialQuery(index int, sum E)
	FullyUpdate(upd Id)
	PartialUpdate(index int, upd Id)

	// Optional “striped” hooks (no‑ops by default).
	AfterPartialUpdate()
	BeforePartialQuery()
}

// ---------- intrusive list node ---------------------------------------------

type node[S, U, E any] struct {
	prev, next *node[S, U, E]
	size       int
	data       Block[S, U, E]
}

func link[S, U, E any](a, b *node[S, U, E]) {
	a.next = b
	b.prev = a
}

// ---------- BlockChain -------------------------------------------------------

// BlockChain keeps elements in blocks of ≤B, guaranteeing O(N/B) operations.
type BlockChain[S, U, E any] struct {
	head, tail *node[S, U, E] // sentinels
	blockSize  int            // B
	size       int            // total #elements
}

// NewEmpty returns an *empty* chain with one 0‑sized block obtained from supplier.
func NewEmpty[S, U, E any](B int, supplier func() Block[S, U, E]) *BlockChain[S, U, E] {
	bc := &BlockChain[S, U, E]{blockSize: B}
	bc.head, bc.tail = &node[S, U, E]{}, &node[S, U, E]{}
	link(bc.head, bc.tail)

	n := &node[S, U, E]{size: 0, data: supplier()}
	link(bc.head, n)
	link(n, bc.tail)
	return bc
}

// NewFilled builds a chain of n elements, partitioned in blocks of size ≤B.
// supplier(l,r) must return a block holding elements [l,r] (0‑based, inclusive).
func NewFilled[S, U, E any](
	n, B int,
	supplier func(l, r int) Block[S, U, E],
) *BlockChain[S, U, E] {
	bc := &BlockChain[S, U, E]{blockSize: B, size: n}
	bc.head, bc.tail = &node[S, U, E]{}, &node[S, U, E]{}
	link(bc.head, bc.tail)

	for i := 0; i < n; i += B {
		l, r := i, min(i+B-1, n-1)
		block := supplier(l, r)
		node := &node[S, U, E]{size: r - l + 1, data: block}
		link(bc.tail.prev, node)
		link(node, bc.tail)
	}
	return bc
}

// --------------------- public high‑level operations --------------------------

// Size returns the total number of stored elements.
func (bc *BlockChain[S, U, E]) Size() int { return bc.size }

// Get element at index (0‑based).
func (bc *BlockChain[S, U, E]) Get(index int) E {
	for n := bc.head.next; n != bc.tail; n = n.next {
		if index >= n.size {
			index -= n.size
			continue
		}
		n.data.BeforePartialQuery()
		return n.data.Get(index)
	}

	panic("index out of range")
}

// Insert e *after* “index” elements – i.e. 0 inserts at front.
func (bc *BlockChain[S, U, E]) Insert(index int, e E) {
	for n := bc.head.next; n != bc.tail; n = n.next {
		if index > n.size {
			index -= n.size
			continue
		}
		n.data.Insert(index, e)
		n.size++
		break
	}
	bc.maintain()
}

// Delete element at index.
func (bc *BlockChain[S, U, E]) Delete(index int) {
	for n := bc.head.next; n != bc.tail; n = n.next {
		if index >= n.size {
			index -= n.size
			continue
		}
		n.data.Delete(index)
		n.size--
		break
	}
	bc.maintain()
}

// Split splits the chain so that the k‑th element (1‑based) becomes the first
// element of the *right* part. Returns (left,right).
func (bc *BlockChain[S, U, E]) Split(
	k int,
	supplier func() Block[S, U, E],
) (*BlockChain[S, U, E], *BlockChain[S, U, E]) {
	if k <= 1 {
		return NewEmpty(bc.blockSize, supplier), bc
	}
	if k > bc.size {
		return bc, NewEmpty(bc.blockSize, supplier)
	}

	headNode := bc.splitKth(k)
	endNode := bc.tail.prev

	// cut out [headNode … endNode]
	headNode.prev.next = bc.tail
	bc.tail.prev = headNode.prev

	// build right chain
	right := &BlockChain[S, U, E]{blockSize: bc.blockSize}
	right.head, right.tail = &node[S, U, E]{}, &node[S, U, E]{}
	link(right.head, headNode)
	link(endNode, right.tail)

	// fix sizes
	bc.maintain()
	right.maintain()
	return bc, right
}

// Update applies upd to range [L,R] (inclusive).
func (bc *BlockChain[S, U, E]) Update(L, R int, upd U) {
	offset := 0
	for n := bc.head.next; n != bc.tail; n = n.next {
		l, r := offset, offset+n.size-1
		offset += n.size
		switch {
		case enter(L, R, l, r):
			n.data.FullyUpdate(upd)
		case leave(L, R, l, r):
			// nothing
		default:
			for i := max(l, L); i <= min(r, R); i++ {
				n.data.PartialUpdate(i-l, upd)
			}
			n.data.AfterPartialUpdate()
		}
	}
}

// Query folds range [L,R] into “sum”.
func (bc *BlockChain[S, U, E]) Query(L, R int, sum S) {
	offset := 0
	for n := bc.head.next; n != bc.tail; n = n.next {
		l, r := offset, offset+n.size-1
		offset += n.size
		switch {
		case enter(L, R, l, r):
			n.data.FullyQuery(sum)
		case leave(L, R, l, r):
			// nothing
		default:
			n.data.BeforePartialQuery()
			for i := max(l, L); i <= min(r, R); i++ {
				n.data.PartialQuery(i-l, sum)
			}
		}
	}
}

// Rotate makes the k‑th element (1‑based) the new first element.
func (bc *BlockChain[S, U, E]) Rotate(k int) {
	node := bc.splitKth(k)
	if node == bc.head.next { // already front
		return
	}
	h1, e1 := bc.head.next, node.prev
	h2, e2 := node, bc.tail.prev

	link(bc.head, h2)
	link(e2, h1)
	link(e1, bc.tail)
	bc.maintain()
}

// Reverse reverses [l,r] (0‑based, inclusive).
func (bc *BlockChain[S, U, E]) Reverse(l, r int) {
	if l >= r {
		return
	}
	left := bc.splitKth(l + 1)
	right := bc.splitKth(r + 2).prev

	begin, end := left.prev, right.next
	right.next = nil

	var rec func(cur, nxt *node[S, U, E])
	rec = func(cur, nxt *node[S, U, E]) {
		if cur == nil {
			return
		}
		rec(cur.next, cur)
		cur.data.Reverse()
		cur.prev, cur.next = cur.next, nxt
	}
	rec(left, nil)

	link(begin, right)
	link(left, end)
	bc.maintain()
}

// String renders the chain (debug‑friendly, not O(1)).
func (bc *BlockChain[S, U, E]) String() string {
	var buf []byte
	buf = append(buf, '[')
	for n := bc.head.next; n != bc.tail; n = n.next {
		buf = append(buf, '<')
		buf = append(buf, []byte(fmt.Sprint(n.data))...)
		buf = append(buf, '>', ',')
	}
	if len(buf) > 1 && buf[len(buf)-1] == ',' {
		buf = buf[:len(buf)-1]
	}
	buf = append(buf, ']')
	return string(buf)
}

// --------------------- internal helpers -------------------------------------

func (bc *BlockChain[S, U, E]) splitNode(n *node[S, U, E], leftSize int) {
	left, right := n.data.Split(leftSize)
	newNode := &node[S, U, E]{size: n.size - leftSize, data: right}

	link(newNode, n.next)
	link(n, newNode)

	n.data, n.size = left, leftSize
}

func (bc *BlockChain[S, U, E]) splitKth(k int) *node[S, U, E] {
	for n := bc.head.next; n != bc.tail; n = n.next {
		if k > n.size {
			k -= n.size
			continue
		}
		if k != 1 {
			bc.splitNode(n, k-1)
			n = n.next
		}
		return n
	}
	return bc.tail
}

func (bc *BlockChain[S, U, E]) mergeNodes(a, b *node[S, U, E]) {
	link(a, b.next)
	a.data = a.data.Merge(b.data)
	a.size += b.size
}

func (bc *BlockChain[S, U, E]) maintain() {
	bc.size = 0
	for n := bc.head.next; n != bc.tail; n = n.next {
		bc.size += n.size
		switch {
		case n.size >= 2*bc.blockSize:
			bc.splitNode(n, bc.blockSize)
		case n.prev != bc.head && n.size+n.prev.size <= bc.blockSize:
			bc.mergeNodes(n.prev, n)
			n = n.prev // step back – current node collapsed into prev
		}
	}
}

// ---------- misc helpers -----------------------------------------------------

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ---------------- block‑local update ----------------------------------------

type upd struct{ val int } // 单点赋值用的「更新」对象

// ---------------- xorBlock ---------------------------------------------------

// xorBlock 满足 Block[*int, upd, int] 接口。
// • arr 保存元素
// • x 保存整个块的异或值，便于 O(1) fullyQuery
type xorBlock struct {
	arr []int
	x   int
}

// 保证 xorBlock 实现接口
var _ Block[*int, upd, int] = (*xorBlock)(nil)

func buildBlock(nums []int) *xorBlock {
	x := 0
	for _, v := range nums {
		x ^= v
	}
	return &xorBlock{arr: append([]int(nil), nums...), x: x}
}

// ---- 必需接口方法 ----

// Split: 前 n 个给左块，其余给右块
func (b *xorBlock) Split(n int) (Block[*int, upd, int], Block[*int, upd, int]) {
	left := buildBlock(b.arr[:n])
	right := buildBlock(b.arr[n:])
	return left, right
}

// Merge: 直接拼接切片
func (b *xorBlock) Merge(other Block[*int, upd, int]) Block[*int, upd, int] {
	o := other.(*xorBlock)
	b.arr = append(b.arr, o.arr...)
	b.x ^= o.x
	return b
}

// --- 可选辅助 ---

func (b *xorBlock) Get(idx int) int { return b.arr[idx] }

func (b *xorBlock) Reverse() {
	for l, r := 0, len(b.arr)-1; l < r; l, r = l+1, r-1 {
		b.arr[l], b.arr[r] = b.arr[r], b.arr[l]
	}
	// b.x 不变
}

// ---- 查询 / 更新 ----

func (b *xorBlock) FullyQuery(sum *int) { *sum ^= b.x }

func (b *xorBlock) PartialQuery(idx int, sum *int) { *sum ^= b.arr[idx] }

func (b *xorBlock) FullyUpdate(_ upd) {
	// 本题没有整块赋值需求，留空
}

func (b *xorBlock) PartialUpdate(idx int, u upd) {
	b.x ^= b.arr[idx] // 去掉旧值
	b.arr[idx] = u.val
	b.x ^= u.val // 加入新值
}

// -------------------------------------------------------------

// getResults 按题目要求返回所有 xor 查询结果
func getResults(nums []int, queries [][]int) []int {
	const B = 350 // √1e5 级别经验值

	// 建链：supplier(l,r) 返回包含区间 [l,r] 的块
	bc := NewFilled(len(nums), B, func(l, r int) Block[*int, upd, int] {
		return buildBlock(nums[l : r+1])
	})

	var res []int
	for _, q := range queries {
		switch q[0] {
		case 1: // 单点赋值
			idx, val := q[1], q[2]
			bc.Update(idx, idx, upd{val})

		case 2: // 区间异或查询
			l, r := q[1], q[2]
			ans := 0
			bc.Query(l, r, &ans)
			res = append(res, ans)

		case 3: // 区间反转
			l, r := q[1], q[2]
			bc.Reverse(l, r)
		}
	}
	return res
}
