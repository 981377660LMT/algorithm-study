package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	loj_141()
}

func demo() {
	tree := NewDoubleEndPalindromeAutomaton(100, 100)
	tree.AppendLeft('a')
	tree.AppendLeft('b')
	tree.AppendLeft('a')
	tree.Append('b')

	// abab
	fmt.Println(tree.PalindromeSubstringCount())
	fmt.Println(tree.DistinctPalindromeSubstring())
	tree.Visit(func(node *TreeNode) {
		fmt.Println(node)
	})
}

// https://loj.ac/p/141
// 1：在字符串 s 的末尾添加一个字符串；
// 2：在字符串 s 的前端添加一个字符串的 反序；
// 3：查询字符串 s 的所有非空回文子串的数量。
func loj_141() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	var q int32
	fmt.Fscan(in, &q)

	n := int32(len(s))

	const S int32 = 4e5 + 10
	P := NewDoubleEndPalindromeAutomaton(S, S)
	for i := int32(0); i < n; i++ {
		P.Append(int32(s[i]))
	}

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var c string
			fmt.Fscan(in, &c)
			for _, v := range c {
				P.Append(v)
			}
		} else if op == 2 {
			var c string
			fmt.Fscan(in, &c)
			for _, v := range c {
				P.AppendLeft(v)
			}
		} else {
			fmt.Fprintln(out, P.PalindromeSubstringCount())
		}
	}
}

type TreeNode struct {
	Next  map[int32]*TreeNode
	Fail  *TreeNode
	Len   int32
	Depth int32
}

func (node *TreeNode) String() string {
	return fmt.Sprintf("TreeNode{Len: %d, Depth: %d, Next: %v}", node.Len, node.Depth, node.Next)
}

// 双端回文树.
type DoubleEndPalindromeAutomaton struct {
	odd, even                     *TreeNode
	data                          []int32
	frontSize, backSize           int32
	frontBuildLast, backBuildLast *TreeNode
	nodes                         []*TreeNode
	palindromeSubstringCount      int
}

func NewDoubleEndPalindromeAutomaton(frontAddition, backAddition int32) *DoubleEndPalindromeAutomaton {
	res := &DoubleEndPalindromeAutomaton{}
	cap_ := frontAddition + backAddition
	nodes := make([]*TreeNode, 0, cap_+2)
	data := make([]int32, cap_)
	zero := frontAddition
	frontSize, backSize := zero-1, zero
	odd := res._alloc()
	odd.Len = -1
	even := res._alloc()
	even.Fail = odd
	frontBuildLast, backBuildLast := odd, odd

	res.odd, res.even = odd, even
	res.data = data
	res.frontSize, res.backSize = frontSize, backSize
	res.frontBuildLast, res.backBuildLast = frontBuildLast, backBuildLast
	res.nodes = nodes
	return res
}

func (tree *DoubleEndPalindromeAutomaton) AppendLeft(char int32) {
	tree.data[tree.frontSize] = char
	tree.frontSize--
	trace := tree.frontBuildLast
	for tree.frontSize+2+trace.Len >= tree.backSize {
		trace = trace.Fail
	}
	for tree.data[tree.frontSize+trace.Len+2] != char {
		trace = trace.Fail
	}
	if v, ok := trace.Next[char]; ok {
		tree.frontBuildLast = v
	} else {
		now := tree._alloc()
		now.Len = trace.Len + 2
		trace.Next[char] = now
		if now.Len == 1 {
			now.Fail = tree.even
		} else {
			trace = trace.Fail
			for tree.data[tree.frontSize+trace.Len+2] != char {
				trace = trace.Fail
			}
			now.Fail = trace.Next[char]
		}
		now.Depth = now.Fail.Depth + 1
		tree.frontBuildLast = now
	}

	if tree.frontBuildLast.Len == tree.backSize-tree.frontSize-1 {
		tree.backBuildLast = tree.frontBuildLast
	}
	tree.palindromeSubstringCount += int(tree.frontBuildLast.Depth)
}

func (tree *DoubleEndPalindromeAutomaton) Append(char int32) {
	tree.data[tree.backSize] = char
	tree.backSize++
	trace := tree.backBuildLast
	for tree.backSize-2-trace.Len <= tree.frontSize {
		trace = trace.Fail
	}
	for tree.data[tree.backSize-trace.Len-2] != char {
		trace = trace.Fail
	}
	if v, ok := trace.Next[char]; ok {
		tree.backBuildLast = v
	} else {
		now := tree._alloc()
		now.Len = trace.Len + 2
		trace.Next[char] = now
		if now.Len == 1 {
			now.Fail = tree.even
		} else {
			trace = trace.Fail
			for tree.data[tree.backSize-trace.Len-2] != char {
				trace = trace.Fail
			}
			now.Fail = trace.Next[char]
		}
		now.Depth = now.Fail.Depth + 1
		tree.backBuildLast = now
	}

	if tree.backBuildLast.Len == tree.backSize-tree.frontSize-1 {
		tree.frontBuildLast = tree.backBuildLast
	}
	tree.palindromeSubstringCount += int(tree.backBuildLast.Depth)
}

// 回文子串总数.
func (tree *DoubleEndPalindromeAutomaton) PalindromeSubstringCount() int {
	return tree.palindromeSubstringCount
}

// 本质不同回文子串总数.
func (tree *DoubleEndPalindromeAutomaton) DistinctPalindromeSubstring() int32 {
	return int32(len(tree.nodes))
}

func (tree *DoubleEndPalindromeAutomaton) Visit(consumer func(*TreeNode)) {
	for i := len(tree.nodes) - 1; i >= 0; i-- {
		consumer(tree.nodes[i])
	}
}

func (tree *DoubleEndPalindromeAutomaton) _alloc() *TreeNode {
	res := &TreeNode{Next: make(map[int32]*TreeNode)}
	tree.nodes = append(tree.nodes, res)
	return res
}
