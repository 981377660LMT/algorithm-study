// https://maspypy.github.io/library/string/suffix_automaton.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/sam.go
// 完全没有理解

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	sa := NewSuffixAutomaton()
	for _, c := range s {
		sa.Add(string(c))
	}
	fmt.Fprintln(out, sa.CountSubstring())
}

const SIZE int = 26     // 字符集大小
const MARGIN byte = 'a' // 字符集的起始字符

type SuffixAutomaton struct {
	size   int // 字符集大小
	margin byte
	nodes  []*Node
	last   int // 文字列全体を入れたときの行き先
}

type Node struct {
	next [SIZE]int // automaton の遷移先
	link int       // suffix link
	size int       // node が受理する最長文字列の長さ
}

func NewSuffixAutomaton() *SuffixAutomaton {
	res := &SuffixAutomaton{size: SIZE, margin: MARGIN}
	res.nodes = append(res.nodes, res.newNode(-1, 0))
	return res
}

func (sa *SuffixAutomaton) Add(char string) {
	c := char[0] - sa.margin
	newNode := len(sa.nodes)
	sa.nodes = append(sa.nodes, sa.newNode(-1, sa.nodes[sa.last].size+1))
	p := sa.last
	for p != -1 && sa.nodes[p].next[c] == -1 {
		sa.nodes[p].next[c] = newNode
		p = sa.nodes[p].link
	}
	q := 0
	if p != -1 {
		q = sa.nodes[p].next[c]
	}

	if p == -1 || sa.nodes[p].size+1 == sa.nodes[q].size {
		sa.nodes[newNode].link = q
	} else {
		newQ := len(sa.nodes)
		sa.nodes = append(sa.nodes, sa.newNode(sa.nodes[q].link, sa.nodes[p].size+1))
		sa.nodes[len(sa.nodes)-1].next = sa.nodes[q].next
		sa.nodes[q].link = newQ
		sa.nodes[newNode].link = newQ
		for p != -1 && sa.nodes[p].next[c] == q {
			sa.nodes[p].next[c] = newQ
			p = sa.nodes[p].link
		}
	}

	sa.last = newNode
}

func (sa *SuffixAutomaton) CalcDAG() [][]int {
	n := len(sa.nodes)
	graph := make([][]int, n)
	for v := 0; v < n; v++ {
		for _, to := range sa.nodes[v].next {
			if to != -1 {
				graph[v] = append(graph[v], to)
			}
		}
	}
	return graph
}

func (sa *SuffixAutomaton) CalcTree() [][]int {
	n := len(sa.nodes)
	graph := make([][]int, n)
	for v := 1; v < n; v++ {
		p := sa.nodes[v].link
		graph[p] = append(graph[p], v)
	}
	return graph
}

// あるノードについて、最短と最長の文字列長が分かればよい。
// 最長は size が持っている
// 最短は、suffix link 先の最長に 1 を加えたものである。
func (sa *SuffixAutomaton) CountSubstringAt(p int) int {
	if p == 0 {
		return 0
	}
	return sa.nodes[p].size - sa.nodes[sa.nodes[p].link].size
}

// 本质不同的子串个数.
func (sa *SuffixAutomaton) CountSubstring() int {
	res := 0
	for i := 0; i < len(sa.nodes); i++ {
		res += sa.CountSubstringAt(i)
	}
	return res
}

func (sa *SuffixAutomaton) newNode(link, size int) *Node {
	res := &Node{link: link, size: size}
	for i := 0; i < sa.size; i++ {
		res.next[i] = -1
	}
	return res
}
