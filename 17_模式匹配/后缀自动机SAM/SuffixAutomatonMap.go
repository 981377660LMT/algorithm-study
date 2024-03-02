// https://maspypy.github.io/library/string/suffix_automaton.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/sam.go

// Blumber 算法在线构建SAM

package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	Next   map[int32]int32 // SAM 转移边
	Link   int32           // 后缀链接
	MaxLen int32           // 当前节点对应的最长子串的长度
	End    int32           // 最长的字符在原串的下标, 实点下标为非负数, 虚点下标为负数
}

type SuffixAutomatonMap struct {
	Nodes   []*Node
	LastPos int32 // 当前插入的字符对应的节点(实点，原串的一个前缀)
	n       int32 // 当前字符串长度
}

func NewSuffixAutomatonMap() *SuffixAutomatonMap {
	res := &SuffixAutomatonMap{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0, -1))
	return res
}

// 每次插入会增加一个实点，可能增加一个虚点.
func (sam *SuffixAutomatonMap) Add(ord int32) int32 {
	newNode := int32(len(sam.Nodes))
	// 新增一个实点以表示当前最长串
	sam.Nodes = append(sam.Nodes, sam.newNode(-1, sam.Nodes[sam.LastPos].MaxLen+1, sam.Nodes[sam.LastPos].End+1))
	p := sam.LastPos
	for p != -1 {
		_, has := sam.Nodes[p].Next[ord]
		if has {
			break
		}
		sam.Nodes[p].Next[ord] = newNode
		p = sam.Nodes[p].Link
	}
	q := int32(0)
	if p != -1 {
		q = sam.Nodes[p].Next[ord]
	}
	if p == -1 || sam.Nodes[p].MaxLen+1 == sam.Nodes[q].MaxLen {
		sam.Nodes[newNode].Link = q
	} else {
		// 不够用，需要新增一个虚点
		newQ := int32(len(sam.Nodes))
		sam.Nodes = append(sam.Nodes, sam.newNode(sam.Nodes[q].Link, sam.Nodes[p].MaxLen+1, -abs32(sam.Nodes[q].End)))
		for k, v := range sam.Nodes[q].Next {
			sam.Nodes[newQ].Next[k] = v
		}
		sam.Nodes[q].Link = newQ
		sam.Nodes[newNode].Link = newQ
		for p != -1 && sam.Nodes[p].Next[ord] == q {
			sam.Nodes[p].Next[ord] = newQ
			p = sam.Nodes[p].Link
		}
	}

	sam.n++
	sam.LastPos = newNode
	return newNode
}

func (sam *SuffixAutomatonMap) Size() int32 {
	return int32(len(sam.Nodes))
}

// 后缀链接树.也叫 parent tree.
func (sam *SuffixAutomatonMap) BuildTree() [][]int32 {
	n := int32(len(sam.Nodes))
	graph := make([][]int32, n)
	for v := int32(1); v < n; v++ {
		p := sam.Nodes[v].Link
		graph[p] = append(graph[p], v)
	}
	return graph
}

func (sam *SuffixAutomatonMap) BuildDAG() [][]int32 {
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
func (sam *SuffixAutomatonMap) GetDfsOrder() []int32 {
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

// 返回每个节点的endPos集合大小.
// !注意：0号结点(空串)大小为n，有时需要置为0.
func (sam *SuffixAutomatonMap) GetEndPosSize(dfsOrder []int32) []int32 {
	size := sam.Size()
	endPosSize := make([]int32, size)
	for i := size - 1; i >= 1; i-- {
		cur := dfsOrder[i]
		if sam.Nodes[cur].End >= 0 { // 实点
			endPosSize[cur]++
		}
		pre := sam.Nodes[cur].Link
		endPosSize[pre] += endPosSize[cur]
	}
	return endPosSize
}

// 给定结点编号和子串长度，返回该子串的起始和结束位置.
func (sam *SuffixAutomatonMap) RecoverSubstring(pos int32, len int32) (start, end int32) {
	end = abs32(sam.Nodes[pos].End) + 1
	start = end - len
	return
}

func (sam *SuffixAutomatonMap) DistinctSubstringAt(pos int32) int32 {
	if pos == 0 {
		return 0
	}
	return sam.Nodes[pos].MaxLen - sam.Nodes[sam.Nodes[pos].Link].MaxLen
}

// 本质不同的子串个数.
func (sam *SuffixAutomatonMap) DistinctSubstring() int {
	res := 0
	for i := 1; i < len(sam.Nodes); i++ {
		res += int(sam.DistinctSubstringAt(int32(i)))
	}
	return res
}

// 类似AC自动机转移，返回(转移后的位置, 转移后匹配的长度).
func (sam *SuffixAutomatonMap) Move(pos, len, char int32) (nextPos, nextLen int32) {
	if v, ok := sam.Nodes[pos].Next[char]; ok {
		nextPos = v
		nextLen = len + 1
	} else {
		for pos != -1 {
			if _, ok := sam.Nodes[pos].Next[char]; ok {
				break
			}
			pos = sam.Nodes[pos].Link
		}
		if pos == -1 {
			nextPos = 0
			nextLen = 0
		} else {
			nextPos = sam.Nodes[pos].Next[char]
			nextLen = sam.Nodes[pos].MaxLen + 1
		}
	}
	return
}

func (sa *SuffixAutomatonMap) newNode(link, maxLength, end int32) *Node {
	res := &Node{Next: make(map[int32]int32), Link: link, MaxLen: maxLength, End: end}
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
	// p4070()
	longest_common_substring()
}

// P4070 [SDOI2016] 生成魔咒
// https://www.luogu.com.cn/problem/P4070
// 在线求本质不同子串数.
// 按顺序在一个序列的末尾插入数字，每次求出插入后能得到的本质不同的子串个数。
//
// !插入每个字符后新增的子串个数为 len(cur) - len(link(cur))
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

	sam := NewSuffixAutomatonMap()
	curSum := 0
	for _, num := range nums {
		pos := sam.Add(int32(num))
		node := sam.Nodes[pos]
		curSum += int(node.MaxLen - sam.Nodes[node.Link].MaxLen)
		fmt.Fprintln(out, curSum)
	}
}

// https://judge.yosupo.jp/problem/longest_common_substring
// https://oi-wiki.org/string/sam/#%E4%B8%A4%E4%B8%AA%E5%AD%97%E7%AC%A6%E4%B8%B2%E7%9A%84%E6%9C%80%E9%95%BF%E5%85%AC%E5%85%B1%E5%AD%90%E4%B8%B2
func longest_common_substring() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	fmt.Fscan(in, &s, &t)
	swap := false
	if len(s) > len(t) { // 对短串建立SAM
		s, t = t, s
		swap = true
	}

	sam := NewSuffixAutomatonMap()
	for _, c := range s {
		sam.Add(c)
	}

	pos, len_ := int32(0), int32(0)
	bestPos, bestLen := int32(0), int32(0)
	bestTEnd := int32(0)
	for i, c := range t {
		pos, len_ = sam.Move(pos, len_, c) // !前缀t[:i+1]匹配到的最长后缀，长为len_，对应SAM上的结点pos
		if len_ > bestLen {
			bestPos, bestLen = pos, len_
			bestTEnd = int32(i + 1)
		}
	}

	if bestLen == 0 {
		fmt.Fprintln(out, 0, 0, 0, 0)
		return
	}
	sStart, sEnd := sam.RecoverSubstring(bestPos, bestLen)
	tStart, tEnd := bestTEnd-bestLen, bestTEnd
	if swap {
		sStart, sEnd, tStart, tEnd = tStart, tEnd, sStart, sEnd
	}
	fmt.Fprintln(out, sStart, sEnd, tStart, tEnd)
}

// 多串最长公共子串 (len(words)>=2)
// https://leetcode.cn/problems/longest-common-subpath/description/
// 1. 把第一个串当做匹配串,其他的串建立 SAM
// 2. 把第一个串(的每个前缀)在每个 SAM 上跑匹配，并记录每个前缀能匹配到的最长后缀.
func MultiLCS(words [][]int32) (res int) {

}
