// https://maspypy.github.io/library/string/suffix_automaton.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/sam.go

// Blumber 算法在线构建SAM

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 线段树合并可以求出 SAM 的每个节点的 endPos 集合

// P3975 [TJOI2015] 弦论
// https://www.luogu.com.cn/problem/P3975
func P3975() {
	// in := bufio.NewReader(os.Stdin)
	// out := bufio.NewWriter(os.Stdout)
	// defer out.Flush()

	// var s string
	// fmt.Fscan(in, &s)
	// var b, k int
	// fmt.Fscan(in, &b, &k)

	// unique := b == 0

	// if !ok {
	// 	fmt.Fprintln(out, -1)
	// 	return
	// }
	// fmt.Fprintln(out, s[start:end])
}

// https://www.luogu.com.cn/problem/CF123D
// !枚举字符串 s 的每一个本质不同的子串 ss ，令 cnt(ss) 为子串 ss 在字符串 s 中出现的个数，求 ∑ cnt(ss)*(cnt(ss)+1)/2
func cf123d() {}

// https://judge.yosupo.jp/problem/number_of_substrings
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	sa := NewSuffixAutomaton()
	for _, c := range s {
		sa.Add(c)
	}
	fmt.Fprintln(out, sa.CountSubstring())
}

const SIGMA int32 = 26   // 字符集大小
const OFFSET int32 = 'a' // 字符集的起始字符

type Node struct {
	Next      [SIGMA]int32 // 孩子节点
	Link      int32        // 后缀链接
	MaxLength int32        // 当前节点对应的最长子串的长度
}

type SuffixAutomaton struct {
	Nodes  []*Node
	CurPos int32 // 当前插入的字符对应的节点
}

func NewSuffixAutomaton() *SuffixAutomaton {
	res := &SuffixAutomaton{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0))
	return res
}

func (sa *SuffixAutomaton) Add(ord int32) {
	c := ord - OFFSET
	newNode := int32(len(sa.Nodes))
	sa.Nodes = append(sa.Nodes, sa.newNode(-1, sa.Nodes[sa.CurPos].MaxLength+1))
	p := sa.CurPos
	for p != -1 && sa.Nodes[p].Next[c] == -1 {
		sa.Nodes[p].Next[c] = newNode
		p = sa.Nodes[p].Link
	}
	q := int32(0)
	if p != -1 {
		q = sa.Nodes[p].Next[c]
	}
	if p == -1 || sa.Nodes[p].MaxLength+1 == sa.Nodes[q].MaxLength {
		sa.Nodes[newNode].Link = q
	} else {
		newQ := int32(len(sa.Nodes))
		sa.Nodes = append(sa.Nodes, sa.newNode(sa.Nodes[q].Link, sa.Nodes[p].MaxLength+1))
		sa.Nodes[len(sa.Nodes)-1].Next = sa.Nodes[q].Next
		sa.Nodes[q].Link = newQ
		sa.Nodes[newNode].Link = newQ
		for p != -1 && sa.Nodes[p].Next[c] == q {
			sa.Nodes[p].Next[c] = newQ
			p = sa.Nodes[p].Link
		}
	}

	sa.CurPos = newNode
}

// 后缀链接树.也叫 parent tree.
func (sa *SuffixAutomaton) BuildTree() [][]int32 {
	n := int32(len(sa.Nodes))
	graph := make([][]int32, n)
	for v := int32(0); v < n; v++ {
		p := sa.Nodes[v].Link
		graph[p] = append(graph[p], v)
	}
	return graph
}

func (sa *SuffixAutomaton) BuildDAG() [][]int32 {
	n := int32(len(sa.Nodes))
	graph := make([][]int32, n)
	for v := int32(0); v < n; v++ {
		for _, to := range sa.Nodes[v].Next {
			if to != -1 {
				graph[v] = append(graph[v], to)
			}
		}
	}
	return graph
}

// pos 位置对应的子串个数.
// 用最长串的长度减去最短串的长度即可得到以当前节点为结尾的子串个数.
// 最长串的长度记录在节点的 MaxLength 中,最短串的长度可以通过link对应的节点的 MaxLength 加 1 得到.
func (sa *SuffixAutomaton) CountSubstringAt(pos int32) int32 {
	if pos == 0 {
		return 0
	}
	return sa.Nodes[pos].MaxLength - sa.Nodes[sa.Nodes[pos].Link].MaxLength
}

// 本质不同的子串个数.
func (sa *SuffixAutomaton) CountSubstring() int {
	res := 0
	for i := 0; i < len(sa.Nodes); i++ {
		res += int(sa.CountSubstringAt(int32(i)))
	}
	return res
}

func (sa *SuffixAutomaton) newNode(link, maxLength int32) *Node {
	res := &Node{Link: link, MaxLength: maxLength}
	for i := int32(0); i < SIGMA; i++ {
		res.Next[i] = -1
	}
	return res
}
