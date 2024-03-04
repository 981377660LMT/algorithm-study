// 广义后缀自动机(General Suffix Automaton，GSA)
// https://www.yhzq-blog.cc/%e5%b9%bf%e4%b9%89%e5%90%8e%e7%bc%80%e8%87%aa%e5%8a%a8%e6%9c%ba%e6%80%bb%e7%bb%93/
// https://zhuanlan.zhihu.com/p/34838533
// https://oi-wiki.org/string/general-sam/
// https://www.luogu.com.cn/article/pm10t1pc
// https://www.luogu.com/article/w967d5rp
//
// note:
//
//  0. 一个能接受多个串所有子串的自动机。
//  1. 构建方式：
//     - 伪广义后缀自动机:
//     !如果给出的是多个字符串而不是一个trie，则可以使用.
//     对每个串，重复在同一个 SAM 上进行建立.
//     !每次建完一个串以后就把lastPos 指针移到root上面，接着建下一个串。
//     注意"插入字符串时需要看一下当前准备插入的位置是否已经有结点了".
//     如果有的话我们只需要在其基础上额外判断一下拆分 SAM 结点的情况；否则的话就和普通的 SAM 插入一模一样了
//     - Trie树上的广义后缀自动机：建立在 Trie 树上的 SAM 称为广义 SAM
//  !2. 自动机和广义后缀自动机中"用于构建"该自动机的所有串的所有前缀节点的树链的并的长度和是 O(L*sqrt(L)) 的。
//     !对文本串t[i]的每个前缀，在后缀链接树上向上跳标记每个endPos，表示该endPos包含了t[i]的子串.标记次数之和不超过O(Lsqrt(L)).
//     记count[u]为结点u的子树中的endPos的原串个数，则sum(count)的数量级为O(Lsqrt(L))，L为所有串长之和.
//     证明(利用根号分治)：
//     则若一个串长S>SQRT(L)，这样的串显然不超SQRT(L)个，而由于广义 SAM 上的节点数量级线性所以这里的总贡献数量级为O(LSQRT(L))。
//     而对于串长不超过SQRT(L)的串，贡献数量级为O(nSQRT(L))。
//     https://blog.csdn.net/ylsoi/article/details/94476894
//     https://ddosvoid.github.io/2020/10/18/%E6%B5%85%E8%B0%88%E6%A0%B9%E5%8F%B7%E7%AE%97%E6%B3%95/
//     喵星球上的点名 https://www.luogu.com.cn/problem/P2336
//     Sevenk Love Oimaster https://www.luogu.com.cn/problem/SP8093

package main

import (
	"bufio"
	"fmt"
	"os"
)

const SIGMA int32 = 26   // 字符集大小
const OFFSET int32 = 'a' // 字符集的起始字符

type Node struct {
	Next   [SIGMA]int32 // SAM 转移边
	Link   int32        // 后缀链接
	MaxLen int32        // 当前节点对应的最长子串的长度
}

type SuffixAutomaton struct {
	Nodes []*Node
	n     int32 // 当前字符串长度
}

func NewSuffixAutomatonGeneral() *SuffixAutomaton {
	res := &SuffixAutomaton{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0))
	return res
}

// !需要在插入新串之前将lastPos置为0.
// eg:
//
//	sam := NewSuffixAutomatonGeneral()
//	for _,word := range words {
//	  lastPos = 0
//	  for _,c := range word {
//	    lastPos = sam.Add(lastPos,c)
//	  }
//	}
func (sam *SuffixAutomaton) Add(lastPos int32, char int32) int32 {
	c := char - OFFSET
	sam.n++

	// 判断当前转移结点是否存在.
	if tmp := sam.Nodes[lastPos].Next[c]; tmp != -1 {
		lastNode, nextNode := sam.Nodes[lastPos], sam.Nodes[tmp]
		if lastNode.MaxLen+1 == nextNode.MaxLen {
			return tmp
		} else {
			newQ := int32(len(sam.Nodes))
			sam.Nodes = append(sam.Nodes, sam.newNode(nextNode.Link, lastNode.MaxLen+1))
			sam.Nodes[newQ].Next = nextNode.Next
			sam.Nodes[tmp].Link = newQ
			for lastPos != -1 && sam.Nodes[lastPos].Next[c] == tmp {
				sam.Nodes[lastPos].Next[c] = newQ
				lastPos = sam.Nodes[lastPos].Link
			}
			return newQ
		}
	}

	newNode := int32(len(sam.Nodes))
	// 新增一个实点以表示当前最长串
	sam.Nodes = append(sam.Nodes, sam.newNode(-1, sam.Nodes[lastPos].MaxLen+1))
	p := lastPos
	for p != -1 && sam.Nodes[p].Next[c] == -1 {
		sam.Nodes[p].Next[c] = newNode
		p = sam.Nodes[p].Link
	}
	q := int32(0)
	if p != -1 {
		q = sam.Nodes[p].Next[c]
	}
	if p == -1 || sam.Nodes[p].MaxLen+1 == sam.Nodes[q].MaxLen {
		sam.Nodes[newNode].Link = q
	} else {
		// 不够用，需要新增一个虚点
		newQ := int32(len(sam.Nodes))
		sam.Nodes = append(sam.Nodes, sam.newNode(sam.Nodes[q].Link, sam.Nodes[p].MaxLen+1))
		sam.Nodes[len(sam.Nodes)-1].Next = sam.Nodes[q].Next
		sam.Nodes[q].Link = newQ
		sam.Nodes[newNode].Link = newQ
		for p != -1 && sam.Nodes[p].Next[c] == q {
			sam.Nodes[p].Next[c] = newQ
			p = sam.Nodes[p].Link
		}
	}
	return newNode
}

func (sam *SuffixAutomaton) AddString(s string) (lastPos int32) {
	lastPos = 0
	for _, c := range s {
		lastPos = sam.Add(lastPos, c)
	}
	return
}

func (sam *SuffixAutomaton) Size() int32 {
	return int32(len(sam.Nodes))
}

// 后缀链接树.也叫 parent tree.
func (sam *SuffixAutomaton) BuildTree() [][]int32 {
	n := int32(len(sam.Nodes))
	graph := make([][]int32, n)
	for v := int32(1); v < n; v++ {
		p := sam.Nodes[v].Link
		graph[p] = append(graph[p], v)
	}
	return graph
}

func (sam *SuffixAutomaton) BuildDAG() [][]int32 {
	n := int32(len(sam.Nodes))
	graph := make([][]int32, n)
	for v := int32(0); v < n; v++ {
		for _, to := range sam.Nodes[v].Next {
			if to != -1 {
				graph[v] = append(graph[v], to)
			}
		}
	}
	return graph
}

// 将结点按照长度进行计数排序，返回后缀链接树的dfs顺序.
// 注意：后缀链接树上父亲的MaxLen值一定小于儿子，但不能认为编号小的节点MaxLen值也小.
// 常数比建图 + dfs 小.
func (sam *SuffixAutomaton) GetDfsOrder() []int32 {
	nodes, size, n := sam.Nodes, sam.Size(), sam.n
	counter := make([]int32, n+1)
	for i := int32(0); i < size; i++ {
		counter[nodes[i].MaxLen]++
	}
	for i := int32(1); i <= n; i++ {
		counter[i] += counter[i-1]
	}
	order := make([]int32, size)
	for i := size - 1; i >= 0; i-- {
		v := nodes[i].MaxLen
		counter[v]--
		order[counter[v]] = i
	}
	return order
}

func (sam *SuffixAutomaton) DistinctSubstringAt(pos int32) int32 {
	if pos == 0 {
		return 0
	}
	return sam.Nodes[pos].MaxLen - sam.Nodes[sam.Nodes[pos].Link].MaxLen
}

// 本质不同的子串个数.
func (sam *SuffixAutomaton) DistinctSubstring() int {
	res := 0
	for i := 1; i < len(sam.Nodes); i++ {
		res += int(sam.DistinctSubstringAt(int32(i)))
	}
	return res
}

// 获取pattern在sam上的位置.
func (sam *SuffixAutomaton) GetPos(pattern string) (pos int32, ok bool) {
	pos = 0
	for _, c := range pattern {
		pos = sam.Nodes[pos].Next[c-OFFSET]
		if pos == -1 {
			return -1, false
		}
	}
	return pos, true
}

func (sam *SuffixAutomaton) newNode(link, maxLen int32) *Node {
	res := &Node{Link: link, MaxLen: maxLen}
	for i := int32(0); i < SIGMA; i++ {
		res.Next[i] = -1
	}
	return res
}

func abs32(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func main() {
	// P6139()
	SP8093()

}

// P3181 [HAOI2016] 找相同字符 (分别维护不同串的 size)
// https://www.luogu.com.cn/problem/P3181
// 求两个字符串的相同子串数量。
func P3181() {}

// P6139 【模板】广义后缀自动机（广义 SAM）
// https://www.luogu.com.cn/problem/P6139
// 求多个字符串的本质不同子串个数.
// 同时需要输出广义后缀自动机的点数.
func P6139() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	sam := NewSuffixAutomatonGeneral()
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		sam.AddString(s)
	}
	fmt.Fprintln(out, sam.DistinctSubstring())
	fmt.Fprintln(out, sam.Size())
}

// bzoj 3926 [Zjoi2015]诸神眷顾的幻想乡 (树上本质不同路径数)
// 给出一颗叶子结点不超过 20 个的无根树，每个节点上都有一个不超过 10 的数字，求树上本质不同的路径个数（两条路径相同定义为：其路径上所有节点上的数字依次相连组成的字符串相同）。
func bzoj3926() {}

// JZPGYZ - Sevenk Love Oimaster
// https://www.luogu.com.cn/problem/SP8093
// 给定 n 个模板串，以及 q 个查询串
// 依次查询每一个查询串是多少个模板串的子串
//
// 对原串建一个广义SAM，然后把每一个原串放到SAM上跑一跑，记录一下每一个状态属于多少个原串，用belongCount表示。
// 这样的话查询串直接在SAM上跑，如果失配输出0，否则直接输出记录在上面的belongCount就好了。
func SP8093() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	texts := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &texts[i])
	}
	patterns := make([]string, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &patterns[i])
	}

	sam := NewSuffixAutomatonGeneral()
	for _, v := range texts {
		sam.AddString(v)
	}

	size := sam.Size()
	nodes := sam.Nodes
	belongCount := make([]int32, size) // 每个状态属于多少个原串
	visitedTime := make([]int32, size)
	for i := range visitedTime {
		visitedTime[i] = -1
	}

	// 对文本串t[i]的每个前缀，在后缀链接树上向上跳标记每个endPos，表示该endPos包含了t[i]的子串.
	// 标记次数之和不超过O(Lsqrt(L)).
	markChain := func(sid int32, pos int32) {
		for pos >= 0 && visitedTime[pos] != sid {
			visitedTime[pos] = sid
			belongCount[pos]++
			pos = nodes[pos].Link
		}
	}

	// 查询模式串是多少个文本串的子串.
	query := func(pattern string) int32 {
		pos := int32(0)
		for _, c := range pattern {
			pos = nodes[pos].Next[c-OFFSET]
			if pos == -1 {
				return 0
			}
		}
		return belongCount[pos]
	}

	// 标记所有文本串的子串.
	for i, w := range texts {
		pos := int32(0)
		for _, c := range w {
			pos = nodes[pos].Next[c-OFFSET]
			markChain(int32(i), pos)
		}
	}

	for _, w := range patterns {
		fmt.Fprintln(out, query(w))
	}
}

// Good Substrings
// https://www.luogu.com.cn/problem/CF316G2
func CF316G2() {}

// Match & Catch
// https://www.luogu.com.cn/problem/CF427D
// 给定两个字符串，求最短的满足各只出现一次的连续公共字串
func CF427D() {}

// Forensic Examination [CF666E] (线段树合并维护 endPosSize)

// G. Death DBMS (死亡笔记数据库管理系统)
// https://codeforces.com/problemset/problem/1437/G
func CF1437G() {}
