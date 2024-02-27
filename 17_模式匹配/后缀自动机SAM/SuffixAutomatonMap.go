// https://maspypy.github.io/library/string/suffix_automaton.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/sam.go

// Blumber 算法在线构建SAM

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://judge.yosupo.jp/problem/number_of_substrings
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	sa := NewSuffixAutomatonMap()
	for _, c := range s {
		sa.Add(c)
	}
	fmt.Fprintln(out, sa.CountSubstring())
}

// P4070 [SDOI2016] 生成魔咒
// https://www.luogu.com.cn/problem/P4070
// 在线求本质不同子串数.
// 按顺序在一个序列的末尾插入数字，每次求出插入后能得到的本质不同的子串个数。
//
// 插入每个字符后新增的子串个数为 len(cur) - len(link(cur))
func p4070() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
}

type Node struct {
	Next      map[int32]int32 // 孩子节点
	Link      int32           // 后缀链接
	MaxLength int32           // 当前节点对应的最长子串的长度
}

type SuffixAutomatonMap struct {
	Nodes  []*Node
	CurPos int32 // 当前插入的字符对应的节点
}

func NewSuffixAutomatonMap() *SuffixAutomatonMap {
	res := &SuffixAutomatonMap{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0))
	return res
}

func (sa *SuffixAutomatonMap) Add(ord int32) {
	newNode := int32(len(sa.Nodes))
	sa.Nodes = append(sa.Nodes, sa.newNode(-1, sa.Nodes[sa.CurPos].MaxLength+1))
	p := sa.CurPos
	for p != -1 {
		_, has := sa.Nodes[p].Next[ord]
		if has {
			break
		}
		sa.Nodes[p].Next[ord] = newNode
		p = sa.Nodes[p].Link
	}
	q := int32(0)
	if p != -1 {
		q = sa.Nodes[p].Next[ord]
	}
	if p == -1 || sa.Nodes[p].MaxLength+1 == sa.Nodes[q].MaxLength {
		sa.Nodes[newNode].Link = q
	} else {
		newQ := int32(len(sa.Nodes))
		sa.Nodes = append(sa.Nodes, sa.newNode(sa.Nodes[q].Link, sa.Nodes[p].MaxLength+1))
		for k, v := range sa.Nodes[q].Next {
			sa.Nodes[newQ].Next[k] = v
		}
		sa.Nodes[q].Link = newQ
		sa.Nodes[newNode].Link = newQ
		for p != -1 && sa.Nodes[p].Next[ord] == q {
			sa.Nodes[p].Next[ord] = newQ
			p = sa.Nodes[p].Link
		}
	}

	sa.CurPos = newNode
}

// 后缀链接树.也叫 parent tree.
func (sa *SuffixAutomatonMap) BuildTree() [][]int32 {
	n := int32(len(sa.Nodes))
	graph := make([][]int32, n)
	for v := int32(0); v < n; v++ {
		p := sa.Nodes[v].Link
		graph[p] = append(graph[p], v)
	}
	return graph
}

func (sa *SuffixAutomatonMap) BuildDAG() [][]int32 {
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
func (sa *SuffixAutomatonMap) CountSubstringAt(pos int32) int32 {
	if pos == 0 {
		return 0
	}
	return sa.Nodes[pos].MaxLength - sa.Nodes[sa.Nodes[pos].Link].MaxLength
}

// 本质不同的子串个数.
func (sa *SuffixAutomatonMap) CountSubstring() int {
	res := 0
	for i := 0; i < len(sa.Nodes); i++ {
		res += int(sa.CountSubstringAt(int32(i)))
	}
	return res
}

func (sa *SuffixAutomatonMap) newNode(link, maxLength int32) *Node {
	res := &Node{Next: make(map[int32]int32), Link: link, MaxLength: maxLength}
	return res
}
