// https://maspypy.github.io/library/string/suffix_automaton.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/sam.go
// Blumber 算法在线构建SAM
//
// note:
//  1. 每一个节点都表示一段子串，所有节点表示的子串们都是唯一的
//  2. len表示的是当前节点的最长长度，当前节点的子串长度范围是 [len-link.len+1, len]
//  3. endPos 集合的大小可以通过topo排序求出来，实际上用桶排实现
//     如果必须要求出 endPos 集合的话，可以用set实现树上自底向上启发式合并
//     如果需要每个点的endPos集合都需要求出来的话，可以用动态开点线段树维护endPos集合，然后使用线段树的合并进行更新
//
// applications:
//  1. 查找某个子串位于哪个节点 => 直接倍增往上跳到len[]合适的地方
//  2. 最长可重叠重复子串 => endPos集合大于等于2的那些节点的最大的范围

package main

import (
	"bufio"
	"fmt"
	"os"
)

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

// Martian Strings
// https://www.luogu.com.cn/problem/CF149E
// 可以找到两个不相交的区间，满足这两个区间对应的子串拼起来和 wi相同
func cf149e() {}

// Fake News (hard)
// https://www.luogu.com.cn/problem/CF802I
// 给出 s，求所有 s 的本质不同子串 ss 在 s 中的出现次数平方和，重复的子串只算一次。
func cf802I() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func() {
		var s string
		fmt.Fscan(in, &s)
	}

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		solve()
	}
}

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
