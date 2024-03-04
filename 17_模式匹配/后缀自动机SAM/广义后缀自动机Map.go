// 广义后缀自动机
// https://zhuanlan.zhihu.com/p/34838533
// https://oi-wiki.org/string/general-sam/
// 每添加一个字符串之后把last设置为root就好

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

type SuffixAutomatonMapGeneral struct {
	Nodes []*Node
	n     int32 // 当前字符串长度
}

func NewSuffixAutomatonMapGeneral() *SuffixAutomatonMapGeneral {
	res := &SuffixAutomatonMapGeneral{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0, -1))
	return res
}

// !需要在插入新串之前将lastPos置为0.
// eg:
//
//	sam := NewSuffixAutomatonMapGeneral()
//	for _,word := range words {
//	  lastPos = 0
//	  for _,c := range word {
//	    lastPos = sam.Add(lastPos,c)
//	  }
//	}
func (sam *SuffixAutomatonMapGeneral) Add(lastPos int32, char int32) int32 {
	sam.n++

	// 判断当前转移结点是否存在.
	if tmp, ok := sam.Nodes[lastPos].Next[char]; ok {
		lastNode, nextNode := sam.Nodes[lastPos], sam.Nodes[tmp]
		if lastNode.MaxLen+1 == nextNode.MaxLen {
			return tmp
		} else {
			newQ := int32(len(sam.Nodes))
			sam.Nodes = append(sam.Nodes, sam.newNode(nextNode.Link, lastNode.MaxLen+1, -abs32(nextNode.End)))
			for k, v := range nextNode.Next {
				sam.Nodes[newQ].Next[k] = v
			}
			sam.Nodes[tmp].Link = newQ
			for lastPos != -1 && sam.Nodes[lastPos].Next[char] == tmp {
				sam.Nodes[lastPos].Next[char] = newQ
				lastPos = sam.Nodes[lastPos].Link
			}
			return newQ
		}
	}

	newNode := int32(len(sam.Nodes))
	// 新增一个实点以表示当前最长串
	sam.Nodes = append(sam.Nodes, sam.newNode(-1, sam.Nodes[lastPos].MaxLen+1, sam.Nodes[lastPos].End+1))
	p := lastPos
	for p != -1 {
		_, has := sam.Nodes[p].Next[char]
		if has {
			break
		}
		sam.Nodes[p].Next[char] = newNode
		p = sam.Nodes[p].Link
	}
	q := int32(0)
	if p != -1 {
		q = sam.Nodes[p].Next[char]
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
		for p != -1 && sam.Nodes[p].Next[char] == q {
			sam.Nodes[p].Next[char] = newQ
			p = sam.Nodes[p].Link
		}
	}
	return newNode
}

func (sam *SuffixAutomatonMapGeneral) AddString(s string) (lastPos int32) {
	lastPos = 0
	for _, c := range s {
		lastPos = sam.Add(lastPos, c)
	}
	return
}

func (sam *SuffixAutomatonMapGeneral) Size() int32 {
	return int32(len(sam.Nodes))
}

// 后缀链接树.也叫 parent tree.
func (sam *SuffixAutomatonMapGeneral) BuildTree() [][]int32 {
	n := int32(len(sam.Nodes))
	graph := make([][]int32, n)
	for v := int32(1); v < n; v++ {
		p := sam.Nodes[v].Link
		graph[p] = append(graph[p], v)
	}
	return graph
}

func (sam *SuffixAutomatonMapGeneral) BuildDAG() [][]int32 {
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
func (sam *SuffixAutomatonMapGeneral) GetDfsOrder() []int32 {
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
func (sam *SuffixAutomatonMapGeneral) GetEndPosSize(dfsOrder []int32) []int32 {
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
func (sam *SuffixAutomatonMapGeneral) RecoverSubstring(pos int32, len int32) (start, end int32) {
	end = abs32(sam.Nodes[pos].End) + 1
	start = end - len
	return
}

func (sam *SuffixAutomatonMapGeneral) DistinctSubstringAt(pos int32) int32 {
	if pos == 0 {
		return 0
	}
	return sam.Nodes[pos].MaxLen - sam.Nodes[sam.Nodes[pos].Link].MaxLen
}

// 本质不同的子串个数.
func (sam *SuffixAutomatonMapGeneral) DistinctSubstring() int {
	res := 0
	for i := 1; i < len(sam.Nodes); i++ {
		res += int(sam.DistinctSubstringAt(int32(i)))
	}
	return res
}

// 类似AC自动机转移，返回(转移后的位置, 转移后匹配的长度).
func (sam *SuffixAutomatonMapGeneral) Move(pos, len, char int32) (nextPos, nextLen int32) {
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

func (sa *SuffixAutomatonMapGeneral) newNode(link, maxLength, end int32) *Node {
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
	P6139()
}

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
	sam := NewSuffixAutomatonMapGeneral()
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		sam.AddString(s)
	}
	fmt.Fprintln(out, sam.DistinctSubstring())
	fmt.Fprintln(out, sam.Size())
}

// 多个字符串的最长公共子串
// 建立广义后缀自动机，然后拿size大小为n的节点的len更新答案

// P2336 [SCOI2012] 喵星球上的点名
// https://www.luogu.com.cn/problem/P2336
func P2336() {}
