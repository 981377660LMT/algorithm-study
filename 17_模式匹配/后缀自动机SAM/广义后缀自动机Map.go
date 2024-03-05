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
}

type SuffixAutomatonMapGeneral struct {
	Nodes []*Node
	n     int32 // 当前字符串长度
}

func NewSuffixAutomatonMapGeneral() *SuffixAutomatonMapGeneral {
	res := &SuffixAutomatonMapGeneral{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0))
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
			sam.Nodes = append(sam.Nodes, sam.newNode(nextNode.Link, lastNode.MaxLen+1))
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
	sam.Nodes = append(sam.Nodes, sam.newNode(-1, sam.Nodes[lastPos].MaxLen+1))
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
		sam.Nodes = append(sam.Nodes, sam.newNode(sam.Nodes[q].Link, sam.Nodes[p].MaxLen+1))
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

func (sam *SuffixAutomatonMapGeneral) AddString(n int32, f func(i int32) int32) (lastPos int32) {
	lastPos = 0
	for i := int32(0); i < n; i++ {
		lastPos = sam.Add(lastPos, f(i))
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

// 获取pattern在sam上的位置.如果不是子串，返回(0,false).
func (sam *SuffixAutomatonMapGeneral) GetPos(m int32, pattern func(i int32) int32) (pos int32, ok bool) {
	pos = 0
	for i := int32(0); i < m; i++ {
		if v, ok := sam.Nodes[pos].Next[pattern(i)]; ok {
			pos = v
		} else {
			return 0, false
		}
	}
	return pos, true
}

func (sa *SuffixAutomatonMapGeneral) newNode(link, maxLength int32) *Node {
	res := &Node{Next: make(map[int32]int32), Link: link, MaxLen: maxLength}
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
	// P2336()
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
		sam.AddString(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	}
	fmt.Fprintln(out, sam.DistinctSubstring())
	fmt.Fprintln(out, sam.Size())
}

// 多个字符串的最长公共子串
// 建立广义后缀自动机，然后拿size大小为n的节点的len更新答案

// P2336 [SCOI2012] 喵星球上的点名
// https://www.luogu.com.cn/problem/P2336
//
// 给定 n 个模板串，以及 q 个查询串.
// 查询：
// 1. 每一个查询串在多少个模板串中出现过(是多少个模板串的子串)
// 2. 每一个模板串包含多少个查询串
func P2336() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	firstName, secondName := make([][]int32, n), make([][]int32, n) // 每个模板串包括姓和名两个部分
	for i := int32(0); i < n; i++ {
		var k1 int32
		fmt.Fscan(in, &k1)
		cur1 := make([]int32, k1)
		for j := int32(0); j < k1; j++ {
			fmt.Fscan(in, &cur1[j])
		}
		firstName[i] = cur1
		var k2 int32
		fmt.Fscan(in, &k2)
		cur2 := make([]int32, k2)
		for j := int32(0); j < k2; j++ {
			fmt.Fscan(in, &cur2[j])
		}
		secondName[i] = cur2
	}
	patterns := make([][]int32, q)
	for i := int32(0); i < q; i++ {
		var k int32
		fmt.Fscan(in, &k)
		cur := make([]int32, k)
		for j := int32(0); j < k; j++ {
			fmt.Fscan(in, &cur[j])
		}
		patterns[i] = cur
	}

	sam := NewSuffixAutomatonMapGeneral()
	for i := int32(0); i < n; i++ {
		v1, v2 := firstName[i], secondName[i]
		sam.AddString(int32(len(v1)), func(i int32) int32 { return v1[i] })
		sam.AddString(int32(len(v2)), func(i int32) int32 { return v2[i] })
	}

	size := sam.Size()
	nodes := sam.Nodes
	res1 := make([]int32, q) // !每一个查询串在多少个模板串中出现过(是多少个模板串的子串)
	res2 := make([]int32, n) // !每一个模板串包含多少个查询串

	belongCount1 := make([]int32, size) // 每个状态属于多少个原串的子串(原串前缀的后缀)
	visitedTime := make([]int32, size)

	// 对文本串t[i]的每个前缀，在后缀链接树上向上跳标记每个endPos，表示该endPos包含了t[i]的子串.
	// 标记次数之和不超过O(Lsqrt(L)).
	for i := range visitedTime {
		visitedTime[i] = -1
	}
	markChain := func(sid int32, pos int32) {
		for pos >= 0 && visitedTime[pos] != sid {
			visitedTime[pos] = sid
			belongCount1[pos]++
			pos = nodes[pos].Link
		}
	}
	// 标记所有文本串的子串.
	for i := int32(0); i < n; i++ {
		pos1 := int32(0)
		w1 := firstName[i]
		for _, c := range w1 {
			pos1 = nodes[pos1].Next[c]
			markChain(int32(i), pos1)
		}
		pos2 := int32(0)
		w2 := secondName[i]
		for _, c := range w2 {
			pos2 = nodes[pos2].Next[c]
			markChain(int32(i), pos2)
		}
	}

	belongCount2 := make([]int32, size) // 每个endPos包含多少个查询串.
	for i, w := range patterns {
		pos, ok := sam.GetPos(int32(len(w)), func(i int32) int32 { return w[i] })
		fmt.Println(w, pos, ok, 111)
		if ok {
			res1[i] = belongCount1[pos]
			belongCount2[pos]++
		}
	}

	// 对每个模式串，向上跳统计包含了多少个查询串.
	for i := range visitedTime {
		visitedTime[i] = -1
	}
	queryChain := func(sid int32, pos int32) int32 {
		res := int32(0)
		for pos >= 0 && visitedTime[pos] != sid {
			visitedTime[pos] = sid
			res += belongCount2[pos]
			pos = nodes[pos].Link
		}
		return res
	}
	for i := int32(0); i < n; i++ {
		curRes := int32(0)

		v1 := firstName[i]
		pos1, _ := sam.GetPos(int32(len(v1)), func(i int32) int32 { return v1[i] })
		curRes += queryChain(int32(i), pos1)
		fmt.Println(curRes, i, v1, `1`)

		v2 := secondName[i]
		pos2, _ := sam.GetPos(int32(len(v2)), func(i int32) int32 { return v2[i] })
		curRes += queryChain(int32(i), pos2)
		fmt.Println(curRes, i, v2, `2`)

		res2[i] = curRes
	}

	for _, v := range res1 {
		fmt.Fprintln(out, v)
	}
	for _, v := range res2 {
		fmt.Fprint(out, v, " ")
	}
}
