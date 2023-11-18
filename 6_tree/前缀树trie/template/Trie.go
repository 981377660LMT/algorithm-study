/* eslint-disable @typescript-eslint/no-non-null-assertion */
// 前缀树是一个树状的数据结构，
// 用于高效地存储和检索一系列字符串的前缀。
// 前缀树有许多应用，如自动补全和拼写检查

package main

import "fmt"

func main() {
	trie := NewTrie()
	trie.Insert("hello")
	trie.Insert("apple")
	trie.Insert("banana")
	trie.Enumerate("app", func(index int, node *TrieNode) bool {
		fmt.Println(index, node)
		return false
	})
}

type TrieNode struct {
	PreCount  int // 多少单词以该结点为前缀
	WordCount int // 多少单词以该节点为结束
	Parent    *TrieNode
	Children  map[byte]*TrieNode
}

func NewTrieNode(parent *TrieNode) *TrieNode {
	return &TrieNode{
		Parent:   parent,
		Children: make(map[byte]*TrieNode),
	}
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(nil),
	}
}

func (t *Trie) Insert(s string) {
	if len(s) == 0 {
		return
	}
	root := t.root
	for i := 0; i < len(s); i++ {
		char := s[i]
		if _, ok := root.Children[char]; !ok {
			root.Children[char] = NewTrieNode(root)
		}
		root.Children[char].PreCount++
		root = root.Children[char]
	}
	root.WordCount++
}

func (t *Trie) Remove(s string) {
	if len(s) == 0 {
		return
	}
	root := t.root
	for i := 0; i < len(s); i++ {
		char := s[i]
		root.Children[char].PreCount--
		root = root.Children[char]
	}
	root.WordCount--
}

func (t *Trie) Find(s string) (node *TrieNode, ok bool) {
	if len(s) == 0 {
		return
	}
	root := t.root
	for i := 0; i < len(s); i++ {
		char := s[i]
		if _, check := root.Children[char]; !check {
			return
		}
		root = root.Children[char]
	}
	return root, true
}

func (t *Trie) Enumerate(s string, f func(index int, node *TrieNode) bool) {
	if len(s) == 0 {
		return
	}
	root := t.root
	for i := 0; i < len(s); i++ {
		char := s[i]
		if _, ok := root.Children[char]; !ok {
			return
		}
		root = root.Children[char]
		if f(i, root) {
			return
		}
	}
}
