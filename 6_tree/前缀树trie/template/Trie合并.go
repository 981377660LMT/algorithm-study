// Trie合并(Trie启发式合并，复杂度O(logn))

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF778C()
}

// Peterson Polyglot
// https://www.luogu.com.cn/problem/CF778C
// 给一棵 trie 树，可以删掉某一层的所有节点和边。
// 被删除的节点的子节点会代替当前节点,形成新的一层节点,如果有相同的可以合并。
// 求删掉哪一层边后合并出的 trie 树最小？
//
// !考虑对于每一个深度为i的询问，答案为 "n - 所有depth=i的结点的子树合并后新增的结点数".
func CF778C() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)

	type trieNode struct{ children [SIGMA]int32 }
	newTrieNode := func() *trieNode { return &trieNode{} }
	nodePtr := int32(0) // 0 号为空节点
	nodes := make([]*trieNode, 2*(n+1))
	for i := range nodes {
		nodes[i] = newTrieNode()
	}
	subSum := make([]int32, n+1) // 每个深度对应的节点数

	alloc := func() int32 {
		nodePtr++
		return nodePtr
	}

	var mergeFn func(a, b int32) int32
	mergeFn = func(a, b int32) int32 {
		if a == 0 || b == 0 {
			return a + b
		}
		newNode := alloc()
		for i := int32(0); i < SIGMA; i++ {
			nodes[newNode].children[i] = mergeFn(nodes[a].children[i], nodes[b].children[i])
		}
		return newNode
	}

	var dfs func(cur int32, depth int32)
	dfs = func(cur int32, depth int32) {
		// 复用结点
		nodePtr = n + 1
		newNode := n + 1
		for i := int32(0); i < SIGMA; i++ {
			if next := nodes[cur].children[i]; next != 0 {
				newNode = mergeFn(newNode, next)
			}
		}
		subSum[depth] += nodePtr - (n + 1) // !这一层的孩子合并后，形成的新节点数(合并的结点数)
		for i := int32(0); i < SIGMA; i++ {
			if next := nodes[cur].children[i]; next != 0 {
				dfs(next, depth+1)
			}
		}
	}

	addEdge := func(u, v int32, c int32) {
		nodes[u].children[c-OFFSET] = v
	}

	for i := int32(0); i < n-1; i++ {
		var parent, child int32
		var char string
		fmt.Fscan(in, &parent, &child, &char)
		addEdge(parent, child, int32(char[0]))
	}
	dfs(1, 1)

	bestSum, bestDepth := int32(0), int32(0)
	for d := int32(1); d <= n; d++ {
		if subSum[d] > bestSum {
			bestSum = subSum[d]
			bestDepth = d
		}
	}

	fmt.Fprintln(out, n-bestSum)
	fmt.Fprintln(out, bestDepth)
}

// Trie合并.
// !由于每次合并向下递归时子树大小都会减少一半，所以合并最多O(logn)次.
const SIGMA int32 = 26
const OFFSET int32 = 'a'

// !将b和的c信息合并到a上, TODO.
func merge(a *TrieNode, b, c *TrieNode) {}

type TrieNode struct {
	children [SIGMA]*TrieNode // nil表示无子节点
}

func NewTrieNode() *TrieNode {
	return &TrieNode{}
}

// 用一个新的节点存合并的结果.
func Merge(a, b *TrieNode) *TrieNode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	newNode := NewTrieNode()
	merge(newNode, a, b)
	for i := int32(0); i < SIGMA; i++ {
		newNode.children[i] = Merge(a.children[i], b.children[i])
	}
	return newNode
}

// 把第二棵树直接合并到第一棵树上，比较省空间，缺点是会丢失合并前树的信息.
func MergeDestructively(a, b *TrieNode) *TrieNode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	merge(a, a, b)
	for i := int32(0); i < SIGMA; i++ {
		a.children[i] = MergeDestructively(a.children[i], b.children[i])
	}
	return a
}

func Insert(node *TrieNode, n int32, f func(i int32) int32) *TrieNode {
	for i := int32(0); i < n; i++ {
		char := f(i) - OFFSET
		if node.children[char] == nil {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}
	return node
}

func Search(node *TrieNode, n int32, f func(i int32) int32) *TrieNode {
	for i := int32(0); i < n; i++ {
		char := f(i) - OFFSET
		if node.children[char] == nil {
			return nil
		}
		node = node.children[char]
	}
	return node
}
