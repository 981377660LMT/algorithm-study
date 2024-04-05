package main

type TrieNode[K comparable] struct {
	children map[K]int32
	parent   int32
	endCount int32
}

func NewTrieNode[K comparable](parent int32) *TrieNode[K] {
	return &TrieNode[K]{children: make(map[K]int32), parent: parent}
}

type TrieTreeMap[K comparable] struct {
	nodes []*TrieNode[K]
}

func NewTrieTree[K comparable]() *TrieTreeMap[K] {
	newNode := NewTrieNode[K](-1)
	return &TrieTreeMap[K]{nodes: []*TrieNode[K]{newNode}}
}

func (t *TrieTreeMap[K]) Insert(n int32, f func(i int32) K) (pos int32) {
	for i := int32(0); i < n; i++ {
		char := f(i)
		if next, ok := t.nodes[pos].children[char]; ok {
			pos = next
		} else {
			next = int32(len(t.nodes))
			t.nodes[pos].children[char] = next
			t.nodes = append(t.nodes, NewTrieNode[K](pos))
			pos = next
		}
	}
	t.nodes[pos].endCount++
	return pos
}

func (t *TrieTreeMap[K]) Search(n int32, f func(i int32) K) (pos int32, ok bool) {
	nodes := t.nodes
	for i := int32(0); i < n; i++ {
		char := f(i)
		if next, ok := nodes[pos].children[char]; ok {
			pos = next
		} else {
			return 0, false
		}
	}
	return pos, true
}

func (t *TrieTreeMap[K]) BuildTree() [][]int32 {
	size := t.Size()
	nodes := t.nodes
	res := make([][]int32, size)
	for i := int32(1); i < size; i++ {
		p := nodes[i].parent
		res[p] = append(res[p], i)
	}
	return res
}

func (t *TrieTreeMap[K]) Size() int32 {
	return int32(len(t.nodes))
}

func (t *TrieTreeMap[K]) Empty(pos int32) bool {
	return len(t.nodes) == 1
}
