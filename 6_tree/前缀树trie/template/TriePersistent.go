package main

import "fmt"

func main() {
	trie := NewTriePersistent(true)
	root := trie.NewRoot()
	root1 := trie.Insert(root, []byte("apple"))
	root2 := trie.Insert(root1, []byte("app"))
	root3 := trie.Remove(root2, []byte("app"))

	fmt.Println(trie.Find(root, []byte("app")))
	fmt.Println(trie.Find(root1, []byte("app")))
	fmt.Println(trie.Find(root2, []byte("app")))
	fmt.Println(trie.Find(root3, []byte("app")))
}

type TrieNode struct {
	PreCount  int // 多少单词以该结点为前缀
	WordCount int // 多少单词以该节点为结束
	Children  [26]*TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{}
}

type TriePersistent struct {
	persistent bool
}

func NewTriePersistent(persistent bool) *TriePersistent {
	return &TriePersistent{persistent: persistent}
}

func (trie *TriePersistent) NewRoot() *TrieNode {
	return nil
}

func (trie *TriePersistent) Insert(root *TrieNode, s []byte) *TrieNode {
	if len(s) == 0 {
		return root
	}
	if root == nil {
		root = NewTrieNode()
	}
	return trie._insert(root, s, 0)
}

// 需要保证s存在.
func (trie *TriePersistent) Remove(root *TrieNode, s []byte) *TrieNode {
	if len(s) == 0 {
		return root
	}
	if root == nil {
		panic("can not remove from nil root")
	}
	return trie._remove(root, s, 0)
}

func (trie *TriePersistent) Find(root *TrieNode, s []byte) (res *TrieNode, ok bool) {
	if root == nil || len(s) == 0 {
		return
	}
	for _, char := range s {
		ord_ := ord(char)
		if root.Children[ord_] == nil {
			return nil, false
		}
		root = root.Children[ord_]
	}
	return root, true
}

func (trie *TriePersistent) Enumerate(root *TrieNode, s []byte, f func(index int, node *TrieNode) bool) {
	if root == nil || len(s) == 0 {
		return
	}
	for i, char := range s {
		ord_ := ord(char)
		if root.Children[ord_] == nil {
			return
		}
		root = root.Children[ord_]
		if f(i, root) {
			return
		}
	}
}

func (trie *TriePersistent) Copy(node *TrieNode) *TrieNode {
	if node == nil || !trie.persistent {
		return node
	}
	return &TrieNode{
		PreCount:  node.PreCount,
		WordCount: node.WordCount,
		Children:  node.Children,
	}
}

func (trie *TriePersistent) _insert(root *TrieNode, s []byte, depth int) *TrieNode {
	root = trie.Copy(root)
	if depth == len(s) {
		root.WordCount++
		return root
	}
	ord_ := ord(s[depth])
	if root.Children[ord_] == nil {
		root.Children[ord_] = NewTrieNode()
	}
	root.Children[ord_] = trie._insert(root.Children[ord_], s, depth+1)
	root.Children[ord_].PreCount++
	return root
}

func (trie *TriePersistent) _remove(root *TrieNode, s []byte, depth int) *TrieNode {
	root = trie.Copy(root)
	if depth == len(s) {
		root.WordCount--
		return root
	}
	ord_ := ord(s[depth])
	root.Children[ord_] = trie._remove(root.Children[ord_], s, depth+1)
	root.Children[ord_].PreCount--
	return root
}

func ord(c byte) byte { return c - 'a' }
func chr(v byte) byte { return v + 'a' }

// 1804. 实现 Trie （前缀树） II
// https://leetcode.cn/problems/implement-trie-ii-prefix-tree/
type Trie struct {
	pt   *TriePersistent
	root *TrieNode
}

func Constructor() Trie {
	pt := NewTriePersistent(false)
	return Trie{
		pt:   pt,
		root: pt.NewRoot(),
	}
}

func (this *Trie) Insert(word string) {
	this.root = this.pt.Insert(this.root, []byte(word))
}

func (this *Trie) CountWordsEqualTo(word string) int {
	res := 0
	this.pt.Enumerate(this.root, []byte(word), func(index int, node *TrieNode) bool {
		if index == len(word)-1 {
			res = node.WordCount
		}
		return false
	})
	return res
}

func (this *Trie) CountWordsStartingWith(prefix string) int {
	res := 0
	this.pt.Enumerate(this.root, []byte(prefix), func(index int, node *TrieNode) bool {
		if index == len(prefix)-1 {
			res = node.PreCount
		}
		return false
	})
	return res
}

func (this *Trie) Erase(word string) {
	this.root = this.pt.Remove(this.root, []byte(word))
}

/**
 * Your Trie object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Insert(word);
 * param_2 := obj.CountWordsEqualTo(word);
 * param_3 := obj.CountWordsStartingWith(prefix);
 * obj.Erase(word);
 */
