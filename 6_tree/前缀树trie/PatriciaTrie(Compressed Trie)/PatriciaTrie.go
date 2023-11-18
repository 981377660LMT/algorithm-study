// Patricia Trie/RadixTree/Compressed Trie
// https://sotanishy.github.io/cp-library-cpp/string/trie.cpp

package main

import "fmt"

type TrieNode struct {
	Parent      *TrieNode
	Children    map[int]*TrieNode
	PrefixCount int
	Accept      []int
	Value       string
}

func (tn *TrieNode) String() string {
	return fmt.Sprintf("%s(%d)", tn.Value, tn.Accept)
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{Children: make(map[int]*TrieNode)},
	}
}

func (t *Trie) Insert(s string, id int) {
	if len(s) == 0 {
		return
	}
	node := t.root
	node.PrefixCount++
	for _, v := range s {
		c := int(v)
		if _, ok := node.Children[c]; !ok {
			node.Children[c] = &TrieNode{Children: make(map[int]*TrieNode), Parent: node, Value: string(v)}
		}
		node = node.Children[c]
		node.PrefixCount++
	}
	node.Accept = append(node.Accept, id)
}

// 压缩 Trie 构建 Patricia Trie.
// 树的高度变为 sqrt(字符集大小).
// 时间复杂度 O(字符集大小).
func (t *Trie) Compress() {
	t._compress(t.root)
}

func (t *Trie) _compress(node *TrieNode) {
	// 压缩单个字符的节点
	for len(node.Accept) == 0 && len(node.Children) == 1 {
		for _, c := range node.Children {
			node.Children = c.Children
			node.Accept = c.Accept
			node.Value += c.Value
			for _, w := range node.Children {
				w.Parent = node
			}
			t._compress(node)
		}
	}
	for _, u := range node.Children {
		t._compress(u)
	}
}

func main() {
	trie := NewTrie()
	trie.Insert("hello", 1)
	trie.Insert("apple", 2)
	trie.Insert("banana", 3)
	trie.Compress()
	fmt.Println(trie.root) // Output: hello
	fmt.Println(trie.root.Children)
}
