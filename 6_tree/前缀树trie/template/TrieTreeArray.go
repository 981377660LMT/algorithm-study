package main

import "fmt"

const SIGMA int32 = 26
const OFFSET int32 = 'a'

// !将b和的c信息合并到a上.
func merge(a *TrieNode, b, c *TrieNode) {}

type TrieNode struct {
	children [SIGMA]int32 // 0表示无子节点
	parent   int32        // -1表示无父节点
	endCount int32
}

func NewTrieNode(parent int32) *TrieNode {
	res := &TrieNode{parent: parent}
	return res
}

type TrieArray struct {
	nodes []*TrieNode
}

func NewTrieTree() *TrieArray {
	newNode := NewTrieNode(-1)
	return &TrieArray{nodes: []*TrieNode{newNode}}
}

func (t *TrieArray) Insert(n int32, f func(i int32) int32) (pos int32) {
	for i := int32(0); i < n; i++ {
		char := f(i) - OFFSET
		if tmp := t.nodes[pos].children[char]; tmp == 0 {
			next := t._alloc()
			t.nodes[pos].children[char] = next
			pos = next
		} else {
			pos = tmp
		}
	}
	t.nodes[pos].endCount++
	return pos
}

func (t *TrieArray) Search(n int32, f func(i int32) int32) (pos int32, ok bool) {
	nodes := t.nodes
	for i := int32(0); i < n; i++ {
		char := f(i) - OFFSET
		if next := nodes[pos].children[char]; next == 0 {
			return 0, false
		} else {
			pos = next
		}
	}
	return pos, true
}

func (t *TrieArray) BuildTree() [][]int32 {
	size := t.Size()
	nodes := t.nodes
	res := make([][]int32, size)
	for i := int32(1); i < size; i++ {
		p := nodes[i].parent
		res[p] = append(res[p], i)
	}
	return res
}

func (t *TrieArray) Size() int32 {
	return int32(len(t.nodes))
}

func (t *TrieArray) Empty(pos int32) bool {
	return len(t.nodes) == 1
}

func (t *TrieArray) _alloc() int32 {
	res := int32(len(t.nodes))
	t.nodes = append(t.nodes, NewTrieNode(-1))
	return res
}

// Trie合并.
// 用一个新的节点存合并的结果.
func (t *TrieArray) Merge(a, b int32) int32 {
	if a == 0 || b == 0 {
		return a + b
	}
	newNode := t._alloc()
	merge(t.nodes[newNode], t.nodes[a], t.nodes[b])
	for i := int32(0); i < SIGMA; i++ {
		t.nodes[newNode].children[i] = t.Merge(t.nodes[a].children[i], t.nodes[b].children[i])
	}
	return newNode
}

// 把第二棵树直接合并到第一棵树上，比较省空间，缺点是会丢失合并前树的信息.
func (t *TrieArray) MergeDestructively(a, b int32) int32 {
	if a == 0 || b == 0 {
		return a + b
	}
	merge(t.nodes[a], t.nodes[a], t.nodes[b])
	for i := int32(0); i < SIGMA; i++ {
		t.nodes[a].children[i] = t.MergeDestructively(t.nodes[a].children[i], t.nodes[b].children[i])
	}
	return a
}

func main() {
	trie := NewTrieTree()
	trie.Insert(3, func(i int32) int32 { return int32("abc"[i]) })
	fmt.Println(trie.Search(3, func(i int32) int32 { return int32("abc"[i]) }))
}
