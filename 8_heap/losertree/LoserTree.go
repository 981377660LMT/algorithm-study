// https://oi-wiki.org/ds/loser-tree/
// https://grafana.com/blog/2024/04/23/the-loser-tree-data-structure-how-to-optimize-merges-and-make-your-programs-run-faster/
// https://leetcode.cn/problems/merge-k-sorted-lists/description/ 合并k个有序链表
// https://pkg.go.dev/github.com/grafana/phlare/pkg/util/loser#New 带有Close方法的LoserTree
// https://www.youtube.com/watch?v=rJfvv_c9mYU 用于日志存储项目grafana/loki中合并有序日志
// https://www.cise.ufl.edu/~sahni/cop3530/slides/lec252.pdf
//
// 败者树，一种为多路归并排序设计的高效数据结构.
// 使用败者树实现快速合并，optimize sorting and merging
// !不过感觉这个和堆实现的多路归并差别不大，理论上 loser tree 常数小一点.
//
// LoserTree(Tournament Tree) data structure, for fast k-way merge
// Loser tree, from https://en.wikipedia.org/wiki/K-way_merge_algorithm#Tournament_Tree
//
// api:
//  func NewLoserTree[E any](iters []Iterator[E], maxVal E, less func(E, E) bool) *LoserTree[E] {
//  func (t *LoserTree[E]) Winner() E
//  func (t *LoserTree[E]) Next() bool
//  func (t *LoserTree[E]) Push(iter Iterator[E])
//
// Loser tree, from https://en.wikipedia.org/wiki/K-way_merge_algorithm#Tournament_Tree

package main

type Iterator[E any] interface {
	Next() bool // Advances and returns true if there is a value at this new position.
	Value() E
}

type Node[E any] struct {
	index int         // This is the loser for all nodes except the 0th, where it is the winner.
	value E           // Value copied from the loser node, or winner for node 0.
	items Iterator[E] // Only populated for leaf nodes.
}

// A loser tree is a binary tree laid out such that nodes N and N+1 have parent N/2.
// We store M leaf nodes in positions M...2M-1, and M-1 internal nodes in positions 1..M-1.
// Node 0 is a special node, containing the winner of the contest.
type LoserTree[E any] struct {
	maxVal E
	less   func(E, E) bool
	nodes  []Node[E]
}

func NewLoserTree[E any](iters []Iterator[E], maxVal E, less func(E, E) bool) *LoserTree[E] {
	offset := len(iters)
	t := LoserTree[E]{
		maxVal: maxVal,
		less:   less,
		nodes:  make([]Node[E], offset*2),
	}
	for i, s := range iters {
		t.nodes[i+offset].items = s
		t.moveNext(i + offset) // Must call Next on each item so that At() has a value.
	}
	if offset > 0 {
		t.nodes[0].index = -1 // flag to be initialized on first call to Next().
	}
	return &t
}

func (t *LoserTree[E]) Winner() Iterator[E] {
	return t.nodes[t.nodes[0].index].items
}

func (t *LoserTree[E]) Next() bool {
	if len(t.nodes) == 0 {
		return false
	}
	if t.nodes[0].index == -1 { // If tree has not been initialized yet, do that.
		t.initialize()
		return t.nodes[t.nodes[0].index].index != -1
	}
	if t.nodes[t.nodes[0].index].index == -1 { // already exhausted
		return false
	}
	t.moveNext(t.nodes[0].index)
	t.replayGames(t.nodes[0].index)
	return t.nodes[t.nodes[0].index].index != -1
}

// Add a new sequence to the merge set
func (t *LoserTree[E]) Push(iter Iterator[E]) {
	// First, see if we can replace one that was previously finished.
	for newPos := len(t.nodes) / 2; newPos < len(t.nodes); newPos++ {
		if t.nodes[newPos].index == -1 {
			t.nodes[newPos].index = newPos
			t.nodes[newPos].items = iter
			t.moveNext(newPos)
			t.nodes[0].index = -1 // flag for re-initialize on next call to Next()
			return
		}
	}
	// We need to expand the tree. Pick the next biggest power of 2 to amortise resizing cost.
	size := 1
	for ; size <= len(t.nodes)/2; size *= 2 {
	}
	newPos := size + len(t.nodes)/2
	newNodes := make([]Node[E], size*2)
	// Copy data over and fix up the indexes.
	for i, n := range t.nodes[len(t.nodes)/2:] {
		newNodes[i+size] = n
		newNodes[i+size].index = i + size
	}
	t.nodes = newNodes
	t.nodes[newPos].index = newPos
	t.nodes[newPos].items = iter
	// Mark all the empty nodes we have added as finished.
	for i := newPos + 1; i < len(t.nodes); i++ {
		t.nodes[i].index = -1
		var zero E
		t.nodes[i].value = zero
	}
	t.moveNext(newPos)
	t.nodes[0].index = -1 // flag for re-initialize on next call to Next()
}

func (t *LoserTree[E]) initialize() {
	winners := make([]int, len(t.nodes))
	// Initialize leaf nodes as winners to start.
	for i := len(t.nodes) / 2; i < len(t.nodes); i++ {
		winners[i] = i
	}
	for i := len(t.nodes) - 2; i > 0; i -= 2 {
		// At each stage the winners play each other, and we record the loser in the node.
		loser, winner := t.playGame(winners[i], winners[i+1])
		p := parent(i)
		t.nodes[p].index = loser
		t.nodes[p].value = t.nodes[loser].value
		winners[p] = winner
	}
	t.nodes[0].index = winners[1]
	t.nodes[0].value = t.nodes[winners[1]].value
}

func (t *LoserTree[E]) moveNext(index int) bool {
	n := &t.nodes[index]
	if n.items.Next() {
		n.value = n.items.Value()
		return true
	}
	n.value = t.maxVal
	n.index = -1
	return false
}

// Starting at pos, re-consider all values up to the root.
func (t *LoserTree[E]) replayGames(pos int) {
	// At the start, pos is a leaf node, and is the winner at that level.
	n := parent(pos)
	for n != 0 {
		if t.less(t.nodes[n].value, t.nodes[pos].value) {
			loser := pos
			// Record pos as the loser here, and the old loser is the new winner.
			pos = t.nodes[n].index
			t.nodes[n].index = loser
			t.nodes[n].value = t.nodes[loser].value
		}
		n = parent(n)
	}
	// pos is now the winner; store it in node 0.
	t.nodes[0].index = pos
	t.nodes[0].value = t.nodes[pos].value
}

func (t *LoserTree[E]) playGame(a, b int) (loser, winner int) {
	if t.less(t.nodes[a].value, t.nodes[b].value) {
		return b, a
	}
	return a, b
}

func parent(i int) int { return i / 2 }
