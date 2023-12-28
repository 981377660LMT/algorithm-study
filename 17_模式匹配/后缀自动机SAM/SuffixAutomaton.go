// https://maspypy.github.io/library/string/suffix_automaton.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/sam.go
// 在线构建SAM

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
	for i := 0; i < len(s); i++ {
		sa.Add(s[i])
	}
	fmt.Fprintln(out, sa.CountSubstring())
}

const SIGMA int32 = 26   // 字符集大小
const OFFSET int32 = 'a' // 字符集的起始字符

type Node struct {
	next      [SIGMA]int32 // automaton の遷移先
	link      int32        // suffix link
	maxLength int32        // node が受理する最長文字列の長さ
}

type SuffixAutomaton struct {
	nodes []*Node
	last  int32 // 文字列全体を入れたときの行き先
}

func NewSuffixAutomaton() *SuffixAutomaton {
	res := &SuffixAutomaton{}
	res.nodes = append(res.nodes, res.newNode(-1, 0))
	return res
}

func (sa *SuffixAutomaton) Add(char byte) {
	c := int32(char) - OFFSET
	newNode := int32(len(sa.nodes))
	sa.nodes = append(sa.nodes, sa.newNode(-1, sa.nodes[sa.last].maxLength+1))
	p := sa.last
	for p != -1 && sa.nodes[p].next[c] == -1 {
		sa.nodes[p].next[c] = newNode
		p = sa.nodes[p].link
	}
	q := int32(0)
	if p != -1 {
		q = sa.nodes[p].next[c]
	}
	if p == -1 || sa.nodes[p].maxLength+1 == sa.nodes[q].maxLength {
		sa.nodes[newNode].link = q
	} else {
		newQ := int32(len(sa.nodes))
		sa.nodes = append(sa.nodes, sa.newNode(sa.nodes[q].link, sa.nodes[p].maxLength+1))
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

// 后缀链接树.也叫 parent tree.
func (sa *SuffixAutomaton) BuildTree() [][]int {
	n := len(sa.nodes)
	graph := make([][]int, n)
	for v := 1; v < n; v++ {
		p := sa.nodes[v].link
		graph[p] = append(graph[p], v)
	}
	return graph
}

func (sa *SuffixAutomaton) BuildDAG() [][]int {
	n := len(sa.nodes)
	graph := make([][]int, n)
	for v := 0; v < n; v++ {
		for _, to := range sa.nodes[v].next {
			if to != -1 {
				graph[v] = append(graph[v], int(to))
			}
		}
	}
	return graph
}

// あるノードについて、最短と最長の文字列長が分かればよい。
// 最長は size が持っている
// 最短は、suffix link 先の最長に 1 を加えたものである。
func (sa *SuffixAutomaton) CountSubstringAt(v int) int {
	if v == 0 {
		return 0
	}
	return int(sa.nodes[v].maxLength - sa.nodes[sa.nodes[v].link].maxLength)
}

// 本质不同的子串个数.
func (sa *SuffixAutomaton) CountSubstring() int {
	res := 0
	for i := 0; i < len(sa.nodes); i++ {
		res += sa.CountSubstringAt(i)
	}
	return res
}

func (sa *SuffixAutomaton) newNode(link, size int32) *Node {
	res := &Node{link: link, maxLength: size}
	for i := int32(0); i < SIGMA; i++ {
		res.next[i] = -1
	}
	return res
}
