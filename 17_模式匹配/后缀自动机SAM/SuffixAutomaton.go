// https://www.luogu.com/article/w967d5rp
// https://www.bilibili.com/video/BV1S54y1G7P8
// https://www.cnblogs.com/Linshey/p/14219867.html
// https://maspypy.github.io/library/string/suffix_automaton.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/sam.go
// https://yutong.site/sam/ 可视化
// Blumber 算法在线构建SAM
//
// !子串是一个前缀的后缀.
//
// -s: aababa
//
// -fail tree
//
//	     !所有的叶子结点都是一个前缀，所有的前缀都是叶子结点(除开1号结点)，所有的叶子结点endPos大小都是1.
//		    link 指向比当前结点短的最长后缀endPos集.
//		    每个节点中的子串，都是以该节点为根的子树中的所有子串的"后缀"
//
//							      0,1,2,3,4,5 ""
//							     /              \
//							    /					  	   \
//							   /					  	    \
//							  /						  	  	 \
//							0,1,3,5 "a"(后)				 2,4 "b"
//							 /	 \	                / "ab"
//							/		  \						     /		  \
//						 /       \              /        \
//						/         \            /          \
//					1,"aa"   3,5 "ba"(后) 	2 "aab"      4 "bab"
//		                / "aba"                     "abab"
//						       /      \                    "aabab"
//						      /        \
//					       3 "aaba"   5 "baba"(后)
//		                         "ababa"
//		                        "aababa"
//
// note:
//  0. 后缀自动机 (Suffix Automaton, SAM) 是仅接受后缀且状态数最少的 DFA.
//  1. 每一个节点都表示一段子串，所有节点表示的子串们都是唯一的.
//     随着子串长度的减小，它有可能还会出现在其他的地方，于是它的endpos 就会多一些，就会分到其他的状态里。
//  2. len表示的是当前节点的最长长度，当前节点的子串长度范围是 [len-link.len+1, len]
//  3. endPos 集合的大小可以通过topo排序求出来，实际上用桶排实现
//     如果必须要求出 endPos 集合的话，可以用set实现树上自底向上启发式合并
//     如果需要每个点的endPos集合都需要求出来的话，可以用动态开点线段树维护endPos集合，然后使用线段树的合并进行更新
//  4. 子节点最短串的最长后缀=父结点最长串
//  5. 两个endPos集合要么包含要么不相交
//  6. 一个子串出现次数就是其对应 endPos 集合的size(注意不是长度范围).
//     !由于子串<=>前缀的后缀，
//     !可以先通过在 SAM 上找到该子串所处的节点，然后求以该节点为根的子树中，有多少个包含原串前缀的节点
//     !另一个含义——从SAM的根到这个结点的转移路径条数。
//  7. 可以把SAM理解为把某个串的所有子串建立AC自动机。
//     !8. 设 lcs(i,j) 为前缀i,j的最长公共后缀长度，其等于fail树上 LCA 的len 值。(反串可以将lcp转化为lcs)
//  9. 一个endpos等价类内的串的长度连续.
//     10.理解
//     - 从 SAM 的定义上理解：
//     SAM 可以看作一种加强版的 Trie，它可以高度压缩一个字符串的子串信息，
//     !一条从根出发到`终止结点`的路径对应了原串的一个后缀，而任意一个从根出发的路径对应了原串一个子串。
//     子串和从根出发的路径一一对应。在这种的理解下，每一个结点的含义并不是固定的，
//     它到底对应哪个子串取决于那条路径是怎么到达它的；而边有着确定的含义。
//     - 从 Parent Tree 的角度去理解连边的含义
//     两个不同等价类的Endpos集合要么无交集，要么相包含，因此可以建出一个由 Endpos集合的包含关系连结而成的树——Parent Tree
//     它的连边——后缀链接，若是向下看，是在一个等价类的前面加上一个字符，从而分成若干的其他等价类；
//     向上看，它是指向包含当前集合的最小的集合。
//     !而后缀自动机的连边是在一个等价类的后面加上一个字母，看看它会指向谁，显然对于同一个添上的字母，这个指向是唯一确定的。
//     - 从结点的含义去理解：
//     每一个结点都对应了一种子串，Parent Tree 的结点与 SAM 的结点一一对应
//     但是, 后缀自动机的边不同于 parent 树上的边
//
// !11. 转移边：parent树往下走代表往前加字符，SAM转移边往后走代表往后加字符
// !12. 子串是什么：
//
//	从SAM的DAG角度看，子串是后缀的一个前缀；
//	!从SAM的Parent Tree角度看，子串是前缀的一个后缀。
//
// !13. SAM 与AC自动机的相似性：
//
//	   AC自动机的失配链接和后缀自动机的后缀链接都有性质：
//	   指向的两个状态都满足"后者的代表串是前者的代表串的真后缀"。
//	   可以把 SAM 理解为把某个串的所有子串建立AC自动机.
//	14. 增量构造中，每次从后面加入一个字符, 有两件事要干：
//	   找出能转移到这个状态的状态，建立链接；确定这个状态的min，即找到它在parent树上的父亲。
//	15. 对于SAM任何一个节点u，从根到这个节点的路线有 `maxLen(u)-minLen(u)+1` 条，而这条路线则表示原字符串的一个子串，且各不相同.
//	16. 一般来讲,DAG上可能重复转移,是很难跑计数DP的。
//	   !但是我们知道后缀自动机的性质 : 任意两个节点的表示集合没有交。
//	   !所以我们只要统计路径数即可,不需要考虑重复问题。
//	   !17.可以通过parent树确定SAM的接受状态集合。找到MaxLen=n的结点，该结点到根的路径上的所有结点都是接受状态。
//
// applications:
//  1. 查找某个子串位于哪个节点 => 直接倍增往上跳到len[]合适的地方
//  2. 最长可重叠重复子串 => endPos集合大于等于2的那些节点的最大的范围
//  3. `在线`给出模式串的模式匹配问题(单模式串离线=>KMP，多模式串离线=>AC自动机，多模式串在线=>SAM)
//     一般有固定模式串的字符串处理问题和固定主串的字符串处理问题两大类问题。当固定模式串时，熟知的 AC 自动机算法便可以胜任这类问题。
//     如果主串固定，一般采用对主串构造后缀树、后缀自动机来解决这一类问题。
//  4. 两个字符s和t的最长公共子串 => 对s建立SAM，对t的每个前缀，在SAM中寻找这个前缀的最长后缀，类似AC自动机跳fail.
//  5. 最长不可重叠重复子串 => endPos 集合大于等于2，而且还需要考虑最靠右的那个位置和最靠左的那个位置之间的距离
//     if(sz[u] >= 2) res = max(res, min(maxLen[u], r[u] - l[u]));
//  6. 读入字符串时删除首部字符 => 记录已读入的字符串长度，若小于等于当前状态的 parent.MaxLen ，就转移到parent
//  7. 判断子串/后缀 => 建出文本串的SAM，将模式串分别输入SAM，若无法转移到则不是子串，否则是；若转移到接受状态则是后缀，否则不是。
//  8. 子串定位 => s[l:r]时s[:r]的一个后缀，len 不小于(r-l)的最近祖先.
//  9. 线段树合并维护每个点的endPos集合 => DFS 树结构时将父结点的线段树并上儿子的线段树，注意使用不销毁点写法的线段树合并。
//  10. 区间子“SAM”.
//  11. 多次询问区间本质不同子串 => 扫描线+LCT 维护区间子串信息
//  12. 维护endPos等价类最小End、最大End => 后缀链接树自底向上更新End下标信息
//  13. 判断子串s[a:b)是否为某个前缀s[:i)的子串 => 先定位子串到pos结点，然后判断该子串最早结束位置 endPosMinEnd[pos] <= i.
//  14. 求子串s[a:b)在s[c:d)中的出现次数 =>
//     先定位子串到pos结点，然后根据"子串是前缀的后缀"可知满足其他位置出现的s[a:b)都在pos的子树中。
//     合法的endPos满足：记长度 m=b-a，叶子结点的 prefixEnd 在 [c+m-1,d) 之间.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

const INF int32 = 1e9 + 10

const SIGMA int32 = 2    // 字符集大小
const OFFSET int32 = 'a' // 字符集的起始字符

type Node struct {
	Next   [SIGMA]int32 // SAM 转移边
	Link   int32        // 后缀链接
	MaxLen int32        // 当前节点对应的最长子串的长度
	End    int32        // 最长的字符在原串的下标, 实点下标为非负数, 虚点下标为负数
}

func (n *Node) String() string {
	return fmt.Sprintf("{Next:%v, Link:%v, MaxLen:%v, End:%v}", n.Next, n.Link, n.MaxLen, n.End)
}

type SuffixAutomaton struct {
	Nodes    []*Node
	LastPos  int32 // 当前插入的字符对应的节点(实点，原串的一个前缀)
	n        int32 // 当前字符串长度
	doubling *DoublingSimple
}

func NewSuffixAutomaton() *SuffixAutomaton {
	res := &SuffixAutomaton{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0, -1))
	return res
}

// 每次插入会增加一个实点，可能增加一个虚点.
// 返回当前前缀对应的节点编号(lastPos).
func (sam *SuffixAutomaton) Add(char int32) int32 {
	c := char - OFFSET
	newNode := int32(len(sam.Nodes))
	// 新增一个实点以表示当前最长串
	sam.Nodes = append(sam.Nodes, sam.newNode(-1, sam.Nodes[sam.LastPos].MaxLen+1, sam.Nodes[sam.LastPos].End+1))
	p := sam.LastPos
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
		sam.Nodes = append(sam.Nodes, sam.newNode(sam.Nodes[q].Link, sam.Nodes[p].MaxLen+1, -abs32(sam.Nodes[q].End)))
		sam.Nodes[len(sam.Nodes)-1].Next = sam.Nodes[q].Next
		sam.Nodes[q].Link = newQ
		sam.Nodes[newNode].Link = newQ
		for p != -1 && sam.Nodes[p].Next[c] == q {
			sam.Nodes[p].Next[c] = newQ
			p = sam.Nodes[p].Link
		}
	}
	sam.n++
	sam.LastPos = newNode
	return sam.LastPos
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

// 返回每个节点的endPos集合大小.
// !注意：0号结点(空串)大小为n，有时需要置为0.
func (sam *SuffixAutomaton) GetEndPosSize(dfsOrder []int32) []int32 {
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

// 每个endPos中(最长串)最小的end.
// 返回每个节点对应子串中在原串中的最小结束位置i(0<=i<n).不存在则为INF.
func (sam *SuffixAutomaton) GetEndPosMinEnd(dfsOrder []int32) []int32 {
	size := sam.Size()
	res := make([]int32, size)
	for i := range res {
		res[i] = INF
	}
	for i := size - 1; i >= 1; i-- {
		cur := dfsOrder[i]
		res[cur] = abs32(sam.Nodes[cur].End)
		pre := sam.Nodes[cur].Link
		res[pre] = min32(res[pre], res[cur])
	}
	return res
}

// 每个endPos中(最长串)最大的end.
// 返回每个节点对应子串中在原串中的最大结束位置i(0<=i<n).不存在则为-1.
func (sam *SuffixAutomaton) GetEndPosMaxEnd(dfsOrder []int32) []int32 {
	size := sam.Size()
	res := make([]int32, size)
	for i := range res {
		res[i] = -1
	}
	for i := size - 1; i >= 1; i-- {
		cur := dfsOrder[i]
		res[cur] = abs32(sam.Nodes[cur].End)
		pre := sam.Nodes[cur].Link
		res[pre] = max32(res[pre], res[cur])
	}
	return res
}

// 线段树合并维护 endPos 集合.
// 将父结点的线段树并上儿子的线段树，使用不销毁点写法的线段树合并。
func (sam *SuffixAutomaton) GetEndPos(dfsOrder []int32) (seg *SegmentTreeOnRange, nodes []*SegNode) {
	size := sam.Size()
	seg = NewSegmentTreeOnRange(0, sam.n-1)
	nodes = make([]*SegNode, size)
	for i := int32(0); i < size; i++ {
		nodes[i] = seg.Alloc()
		if end := sam.Nodes[i].End; end >= 0 {
			seg.Set(nodes[i], end, 1)
		}
	}
	for i := size - 1; i >= 1; i-- {
		cur := dfsOrder[i]
		pre := sam.Nodes[cur].Link
		nodes[pre] = seg.Merge(nodes[pre], nodes[cur])
	}
	return
}

// 快速定位子串, 可以与其它字符串算法配合使用.
// 倍增往上跳到 MaxLen>=end-start 的最后一个节点.
// start: 子串起始位置, end: 子串结束位置, endPosOfEnd: 子串结束位置在fail树上的位置.
func (sam *SuffixAutomaton) LocateSubstring(start, end int32, endPosOfEnd int32) (pos int32) {
	target := end - start
	_, pos = sam.Doubling().MaxStep(endPosOfEnd, func(p int32) bool { return sam.Nodes[p].MaxLen >= target })
	return
}

// 给定结点编号和子串长度，返回该子串的起始和结束位置.
func (sam *SuffixAutomaton) RecoverSubstring(pos int32, len int32) (start, end int32) {
	end = abs32(sam.Nodes[pos].End) + 1
	start = end - len
	return
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

// 类似AC自动机转移，输入一个字符，返回(转移后的位置, 转移后匹配的"最长后缀"长度).
// pos: 当前状态, len: 当前匹配的长度, char: 输入字符.
func (sam *SuffixAutomaton) Move(pos, len, char int32) (nextPos, nextLen int32) {
	char -= OFFSET
	if tmp := sam.Nodes[pos].Next[char]; tmp != -1 {
		nextPos = tmp
		nextLen = len + 1
	} else {
		for pos != -1 && sam.Nodes[pos].Next[char] == -1 {
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

// 删除当前模式串首部字符，返回(转移后的位置, 转移后匹配的"最长后缀"长度).
// pos: 当前状态, len: 当前匹配的长度, patternLen: 当前模式串长度.
func (sam *SuffixAutomaton) MoveLeft(pos, len, patternLen int32) (nextPos, nextLen int32) {
	if len < patternLen { // 没有完全匹配，可以不删字符，匹配到的首字母是模式串某个后缀的首字母
		return pos, len
	}
	if len == 0 {
		return 0, 0
	}
	len--
	node := sam.Nodes[pos]
	if len == sam.Nodes[node.Link].MaxLen {
		pos = node.Link
	}
	return pos, len
}

// 给定模式串pattern，返回模式串的每个非空前缀s[:i+1]与SAM文本串的最长公共后缀长度.
func (sam *SuffixAutomaton) LongestCommonSuffix(m int32, pattern func(i int32) int32) []int32 {
	res := make([]int32, m)
	pos, len := int32(0), int32(0)
	for i := int32(0); i < m; i++ {
		pos, len = sam.Move(pos, len, pattern(i))
		res[i] = len
	}
	return res
}

// 后缀连接树上倍增.
func (sam *SuffixAutomaton) Doubling() *DoublingSimple {
	if sam.doubling == nil {
		size := sam.Size()
		doubling := NewDoubling(size, int(size))
		for i := int32(1); i < size; i++ {
			doubling.Add(i, sam.Nodes[i].Link)
		}
		doubling.Build()
		sam.doubling = doubling
	}
	return sam.doubling
}

func (sam *SuffixAutomaton) Print(s string) {
	dfsOrder := sam.GetDfsOrder()
	tree := sam.BuildTree()
	fmt.Println("---")
	fmt.Println("Fail Tree")
	for i := range tree {
		if len(tree[i]) > 0 {
			fmt.Println(i, tree[i])
		}
	}
	fmt.Println("---")
	fmt.Println("EndPos")
	fmt.Println("pos,longest,minLen,maxLen,prefixEnd")
	for i := int32(1); i < int32(len(dfsOrder)); i++ {
		pos := dfsOrder[i]
		link := sam.Nodes[pos].Link
		minLen, maxLen := sam.Nodes[link].MaxLen+1, sam.Nodes[pos].MaxLen
		start, end := sam.RecoverSubstring(pos, maxLen)
		sub := s[start:end]
		prefixEnd := sam.Nodes[pos].End
		fmt.Println(fmt.Sprintf("%v,%v,%v,%v,%v", pos, sub, minLen, maxLen, prefixEnd))
	}
	fmt.Println("---")
}

func (sam *SuffixAutomaton) newNode(link, maxLen, end int32) *Node {
	res := &Node{Link: link, MaxLen: maxLen, End: end}
	for i := int32(0); i < SIGMA; i++ {
		res.Next[i] = -1
	}
	return res
}

type S = int32

// SparseTable 稀疏表: st[j][i] 表示区间 [i, i+2^j) 的贡献值.
type SparseTable struct {
	st     [][]S
	lookup []int
	e      func() S
	op     func(S, S) S
	n      int
}

func NewSparseTable(n int, f func(int) S, e func() S, op func(S, S) S) *SparseTable {
	res := &SparseTable{}

	b := bits.Len(uint(n))
	st := make([][]S, b)
	for i := range st {
		st[i] = make([]S, n)
	}
	for i := 0; i < n; i++ {
		st[0][i] = f(i)
	}
	for i := 1; i < b; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	lookup := make([]int, n+1)
	for i := 2; i < len(lookup); i++ {
		lookup[i] = lookup[i>>1] + 1
	}
	res.st = st
	res.lookup = lookup
	res.e = e
	res.op = op
	res.n = n
	return res
}

func NewSparseTableFrom(leaves []S, e func() S, op func(S, S) S) *SparseTable {
	return NewSparseTable(len(leaves), func(i int) S { return leaves[i] }, e, op)
}

// 查询区间 [start, end) 的贡献值.
func (st *SparseTable) Query(start, end int) S {
	if start >= end {
		return st.e()
	}
	b := st.lookup[end-start]
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (st *SparseTable) MaxRight(left int, check func(e S) bool) int {
	if left == st.n {
		return st.n
	}
	ok, ng := left, st.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(st.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (st *SparseTable) MinLeft(right int, check func(e S) bool) int {
	if right == 0 {
		return 0
	}
	ok, ng := right, -1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(st.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

type E = int32

type SegNode struct {
	count                 E
	leftChild, rightChild *SegNode
}

func (n *SegNode) String() string {
	return fmt.Sprintf("%v", n.count)
}

type SegmentTreeOnRange struct {
	min, max int32
}

// 指定闭区间[min,max]建立权值线段树.
func NewSegmentTreeOnRange(min, max int32) *SegmentTreeOnRange {
	return &SegmentTreeOnRange{min: min, max: max}
}

// NewRoot().
func (sm *SegmentTreeOnRange) Alloc() *SegNode {
	return &SegNode{}
}

// 权值线段树求第 k 小.
// 调用前需保证 1 <= k <= node.value.
func (sm *SegmentTreeOnRange) Kth(node *SegNode, k int32) (value int32, ok bool) {
	if k < 1 || k > sm._eval(node) {
		return 0, false
	}
	return sm._kth(k, node, sm.min, sm.max), true
}

func (sm *SegmentTreeOnRange) Get(node *SegNode, index int32) E {
	return sm._get(node, index, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) Set(node *SegNode, index int32, value E) {
	sm._set(node, index, value, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) Query(node *SegNode, left, right int32) E {
	return sm._query(node, left, right, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) QueryAll(node *SegNode) E {
	return sm._eval(node)
}

func (sm *SegmentTreeOnRange) Update(node *SegNode, index int32, count E) {
	sm._update(node, index, count, sm.min, sm.max)
}

// 用一个新的节点存合并的结果，会生成重合节点数量的新节点.
func (sm *SegmentTreeOnRange) Merge(a, b *SegNode) *SegNode {
	return sm._merge(a, b, sm.min, sm.max)
}

// 把第二棵树直接合并到第一棵树上，比较省空间，缺点是会丢失合并前树的信息.
func (sm *SegmentTreeOnRange) MergeDestructively(a, b *SegNode) *SegNode {
	return sm._mergeDestructively(a, b, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) Enumerate(node *SegNode, f func(i int32, count E)) {
	sm._enumerate(node, sm.min, sm.max, f)
}

func (sm *SegmentTreeOnRange) _kth(k int32, node *SegNode, left, right int32) int32 {
	if left == right {
		return left
	}
	mid := (left + right) >> 1
	if leftCount := sm._eval(node.leftChild); leftCount >= k {
		return sm._kth(k, node.leftChild, left, mid)
	} else {
		return sm._kth(k-leftCount, node.rightChild, mid+1, right)
	}
}

func (sm *SegmentTreeOnRange) _get(node *SegNode, index int32, left, right int32) E {
	if node == nil {
		return 0
	}
	if left == right {
		return node.count
	}
	mid := (left + right) >> 1
	if index <= mid {
		return sm._get(node.leftChild, index, left, mid)
	} else {
		return sm._get(node.rightChild, index, mid+1, right)
	}
}

func (sm *SegmentTreeOnRange) _query(node *SegNode, L, R int32, left, right int32) E {
	if node == nil {
		return 0
	}
	if L <= left && right <= R {
		return node.count
	}
	mid := (left + right) >> 1
	if R <= mid {
		return sm._query(node.leftChild, L, R, left, mid)
	}
	if L > mid {
		return sm._query(node.rightChild, L, R, mid+1, right)
	}
	return sm._query(node.leftChild, L, R, left, mid) + sm._query(node.rightChild, L, R, mid+1, right)
}

func (sm *SegmentTreeOnRange) _set(node *SegNode, index int32, count E, left, right int32) {
	if left == right {
		node.count = count
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._set(node.leftChild, index, count, left, mid)
	} else {
		if node.rightChild == nil {
			node.rightChild = sm.Alloc()
		}
		sm._set(node.rightChild, index, count, mid+1, right)
	}
	node.count = sm._eval(node.leftChild) + sm._eval(node.rightChild)
}

func (sm *SegmentTreeOnRange) _update(node *SegNode, index int32, count E, left, right int32) {
	if left == right {
		node.count += count
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._update(node.leftChild, index, count, left, mid)
	} else {
		if node.rightChild == nil {
			node.rightChild = sm.Alloc()
		}
		sm._update(node.rightChild, index, count, mid+1, right)
	}
	node.count = sm._eval(node.leftChild) + sm._eval(node.rightChild)
}

func (sm *SegmentTreeOnRange) _enumerate(node *SegNode, left, right int32, f func(i int32, count E)) {
	if node == nil {
		return
	}
	if left == right {
		f(left, node.count)
		return
	}
	mid := (left + right) >> 1
	sm._enumerate(node.leftChild, left, mid, f)
	sm._enumerate(node.rightChild, mid+1, right, f)
}

func (sm *SegmentTreeOnRange) _merge(a, b *SegNode, left, right int32) *SegNode {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	newNode := sm.Alloc()
	if left == right {
		newNode.count = a.count + b.count
		return newNode
	}
	mid := (left + right) >> 1
	newNode.leftChild = sm._merge(a.leftChild, b.leftChild, left, mid)
	newNode.rightChild = sm._merge(a.rightChild, b.rightChild, mid+1, right)
	newNode.count = sm._eval(newNode.leftChild) + sm._eval(newNode.rightChild)
	return newNode
}

func (sm *SegmentTreeOnRange) _mergeDestructively(a, b *SegNode, left, right int32) *SegNode {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	if left == right {
		a.count += b.count
		return a
	}
	mid := (left + right) >> 1
	a.leftChild = sm._mergeDestructively(a.leftChild, b.leftChild, left, mid)
	a.rightChild = sm._mergeDestructively(a.rightChild, b.rightChild, mid+1, right)
	a.count = sm._eval(a.leftChild) + sm._eval(a.rightChild)
	return a
}

func (sm *SegmentTreeOnRange) _eval(node *SegNode) E {
	if node == nil {
		return 0
	}
	return node.count
}

type DoublingSimple struct {
	n   int32
	log int32
	to  []int32
}

func NewDoubling(n int32, maxStep int) *DoublingSimple {
	res := &DoublingSimple{n: n, log: int32(bits.Len(uint(maxStep)))}
	size := n * res.log
	res.to = make([]int32, size)
	for i := int32(0); i < size; i++ {
		res.to[i] = -1
	}
	return res
}

func (d *DoublingSimple) Add(from, to int32) {
	d.to[from] = to
}

func (d *DoublingSimple) Build() {
	n := d.n
	for k := int32(0); k < d.log-1; k++ {
		for v := int32(0); v < n; v++ {
			w := d.to[k*n+v]
			next := (k+1)*n + v
			if w == -1 {
				d.to[next] = -1
				continue
			}
			d.to[next] = d.to[k*n+w]
		}
	}
}

// 求从 `from` 状态开始转移 `step` 次的最终状态的编号。
// 不存在时返回 -1。
func (d *DoublingSimple) Jump(from int32, step int) (to int32) {
	to = from
	for k := int32(0); k < d.log; k++ {
		if to == -1 {
			return
		}
		if step&(1<<k) != 0 {
			to = d.to[k*d.n+to]
		}
	}
	return
}

// 求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step` 以及最终状态的编号。
func (d *DoublingSimple) MaxStep(from int32, check func(next int32) bool) (step int, to int32) {
	for k := d.log - 1; k >= 0; k-- {
		tmp := d.to[k*d.n+from]
		if tmp == -1 {
			continue
		}
		if check(tmp) {
			step |= 1 << k
			from = tmp
		}
	}
	to = from
	return
}

// 注意f(i)>=0.
func NewWaveletMatrix32(n int32, f func(i int32) int32) *WaveletMatrix32 {
	dataCopy := make([]int32, n)
	max_ := int32(0)
	for i := int32(0); i < n; i++ {
		v := f(i)
		if v > max_ {
			max_ = v
		}
		dataCopy[i] = v
	}
	maxLog := int32(bits.Len32(uint32(max_)) + 1)
	mat := make([]*BitVector32, maxLog)
	zs := make([]int32, maxLog)
	buff1 := make([]int32, maxLog)
	buff2 := make([]int32, maxLog)

	ls, rs := make([]int32, n), make([]int32, n)
	for dep := int32(0); dep < maxLog; dep++ {
		mat[dep] = NewBitVector32(n + 1)
		p, q := int32(0), int32(0)
		for i := int32(0); i < n; i++ {
			k := (dataCopy[i] >> (maxLog - dep - 1)) & 1
			if k == 1 {
				rs[q] = dataCopy[i]
				mat[dep].Set(i)
				q++
			} else {
				ls[p] = dataCopy[i]
				p++
			}
		}

		zs[dep] = p
		mat[dep].Build()
		ls = dataCopy
		for i := int32(0); i < q; i++ {
			dataCopy[p+i] = rs[i]
		}
	}

	return &WaveletMatrix32{
		n:      n,
		maxLog: maxLog,
		mat:    mat,
		zs:     zs,
		buff1:  buff1,
		buff2:  buff2,
	}
}

type WaveletMatrix32 struct {
	n            int32
	maxLog       int32
	mat          []*BitVector32
	zs           []int32
	buff1, buff2 []int32
}

// [start, end) 内的 value 的個数.
func (w *WaveletMatrix32) Count(start, end, value int32) int32 {
	return w.count(value, end) - w.count(value, start)
}

// [start, end) 内 [lower, upper) 的个数.
func (w *WaveletMatrix32) CountRange(start, end, lower, upper int32) int32 {
	return w.freqDfs(0, start, end, 0, lower, upper)
}

// 第k(0-indexed)个value的位置.
func (w *WaveletMatrix32) Index(value, k int32) int32 {
	w.count(value, w.n)
	for dep := w.maxLog - 1; dep >= 0; dep-- {
		bit := (value >> uint(w.maxLog-dep-1)) & 1
		k = w.mat[dep].IndexWithStart(bit, k, w.buff1[dep])
		if k < 0 || k >= w.buff2[dep] {
			return -1
		}
		k -= w.buff1[dep]
	}
	return k
}

func (w *WaveletMatrix32) IndexWithStart(value, k, start int32) int32 {
	return w.Index(value, k+w.count(value, start))
}

// [start, end) 内第k(0-indexed)小的数(kth).
func (w *WaveletMatrix32) Kth(start, end, k int32) int32 {
	return w.KthMax(start, end, end-start-k-1)
}

// [start, end) 内第k(0-indexed)大的数.
func (w *WaveletMatrix32) KthMax(start, end, k int32) int32 {
	if k < 0 || k >= end-start {
		return -1
	}
	res := int32(0)
	for dep := int32(0); dep < w.maxLog; dep++ {
		p, q := w.mat[dep].Count(1, start), w.mat[dep].Count(1, end)
		if k < q-p {
			start = w.zs[dep] + p
			end = w.zs[dep] + q
			res |= 1 << uint(w.maxLog-dep-1)
		} else {
			k -= q - p
			start -= p
			end -= q
		}
	}
	return res
}

// [start, end) 内第k(0-indexed)小的数(kth).
func (w *WaveletMatrix32) KthMin(start, end, k int32) int32 {
	return w.KthMax(start, end, end-start-k-1)
}

// [start, end) 中比 value 严格小的数, 不存在返回 -INF.
func (w *WaveletMatrix32) Lower(start, end, value int32) int32 {
	k := w.lt(start, end, value)
	if k != 0 {
		return w.KthMin(start, end, k-1)
	}
	return -INF
}

// [start, end) 中比 value 严格大的数, 不存在返回 INF.
func (w *WaveletMatrix32) Higher(start, end, value int32) int32 {
	k := w.le(start, end, value)
	if k == end-start {
		return INF
	}
	return w.KthMin(start, end, k)
}

// [start, end) 中不超过 value 的最大值, 不存在返回 -INF.
func (w *WaveletMatrix32) Floor(start, end, value int32) int32 {
	count := w.Count(start, end, value)
	if count > 0 {
		return value
	}
	return w.Lower(start, end, value)
}

// [start, end) 中不小于 value 的最小值, 不存在返回 INF.
func (w *WaveletMatrix32) Ceiling(start, end, value int32) int32 {
	count := w.Count(start, end, value)
	if count > 0 {
		return value
	}
	return w.Higher(start, end, value)
}

func (w *WaveletMatrix32) access(k int32) int32 {
	res := int32(0)
	for dep := int32(0); dep < w.maxLog; dep++ {
		bit := w.mat[dep].Get(k)
		res = (res << 1) | bit
		k = w.mat[dep].Count(bit, k) + w.zs[dep]*dep
	}
	return res
}

func (w *WaveletMatrix32) count(value, end int32) int32 {
	left, right := int32(0), end
	for dep := int32(0); dep < w.maxLog; dep++ {
		w.buff1[dep] = left
		w.buff2[dep] = right
		bit := (value >> (w.maxLog - dep - 1)) & 1
		left = w.mat[dep].Count(bit, left) + w.zs[dep]*bit
		right = w.mat[dep].Count(bit, right) + w.zs[dep]*bit
	}
	return right - left
}

func (w *WaveletMatrix32) freqDfs(d, left, right, val, a, b int32) int32 {
	if left == right {
		return 0
	}
	if d == w.maxLog {
		if a <= val && val < b {
			return right - left
		}
		return 0
	}

	nv := (1 << (w.maxLog - d - 1)) | val
	nnv := ((1 << (w.maxLog - d - 1)) - 1) | nv
	if nnv < a || b <= val {
		return 0
	}
	if a <= val && nnv < b {
		return right - left
	}
	lc, rc := w.mat[d].Count(1, left), w.mat[d].Count(1, right)
	return w.freqDfs(d+1, left-lc, right-rc, val, a, b) + w.freqDfs(d+1, lc+w.zs[d], rc+w.zs[d], nv, a, b)
}

func (w *WaveletMatrix32) ll(left, right, v int32) (int32, int32) {
	res := int32(0)
	for dep := int32(0); dep < w.maxLog; dep++ {
		w.buff1[dep] = left
		w.buff2[dep] = right
		bit := v >> uint(w.maxLog-dep-1) & 1
		if bit == 1 {
			res += right - left + w.mat[dep].Count(1, left) - w.mat[dep].Count(1, right)
		}
		left = w.mat[dep].Count(bit, left) + w.zs[dep]*bit
		right = w.mat[dep].Count(bit, right) + w.zs[dep]*bit
	}
	return res, right - left
}

func (w *WaveletMatrix32) lt(left, right, v int32) int32 {
	a, _ := w.ll(left, right, v)
	return a
}

func (w *WaveletMatrix32) le(left, right, v int32) int32 {
	a, b := w.ll(left, right, v)
	return a + b
}

type BitVector32 struct {
	n     int32
	block []int
	sum   []int
}

func NewBitVector32(n int32) *BitVector32 {
	blockCount := (n + 63) >> 6
	return &BitVector32{
		n:     n,
		block: make([]int, blockCount+1),
		sum:   make([]int, blockCount+1),
	}
}

func (f *BitVector32) Set(i int32) {
	f.block[i>>6] |= 1 << uint(i&63)
}

func (f *BitVector32) Build() {
	for i := 0; i < len(f.block)-1; i++ {
		f.sum[i+1] = f.sum[i] + bits.OnesCount(uint(f.block[i]))
	}
}

func (f *BitVector32) Get(i int32) int32 {
	return int32((f.block[i>>6] >> (i & 63))) & 1
}

func (f *BitVector32) Count(value, end int32) int32 {
	mask := (1 << uint(end&63)) - 1
	res := int32(f.sum[end>>6] + bits.OnesCount(uint(f.block[end>>6]&mask)))
	if value == 1 {
		return res
	}
	return end - res
}

func (f *BitVector32) Index(value, k int32) int32 {
	if k < 0 || f.Count(value, f.n) <= k {
		return -1
	}

	left, right := int32(0), f.n
	for right-left > 1 {
		mid := (left + right) >> 1
		if f.Count(value, mid) >= k+1 {
			right = mid
		} else {
			left = mid
		}
	}
	return right - 1
}

func (f *BitVector32) IndexWithStart(value, k, start int32) int32 {
	return f.Index(value, k+f.Count(value, start))
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
	// P3975()
	// P6640()

	// cf149e()
	// cf235c()
	cf653f()
	// cf802I()
	// CF1037H()
	// cf1073g()
	// cf1207g()

	// number_of_substrings()
	// longest_common_substring()

	// nowCoder37092C()
	// nowCoder37092D()
	// nowCoder37092E()
	// nowCoder37092J()

	// testMerge()
}

// P3975 [TJOI2015] 弦论(字典序第k小子串)
// https://www.luogu.com.cn/problem/P3975
//
// 1. 求出后缀链接树上每个endPos集合的size.
// 2. 将endPosSize通过SAM转移边统计到结点上, 即这个节点下面总共有多少个子串。
// 3. !在SAM上按照字典序往下dfs匹配,如果子树内子串个数小于k则跳过，否则答案在这个结点中。
// 4. 根据转移到的SAM上的位置和子串长度即可还原子串.
func P3975() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	// 本质相同的子串在不同位置出现算相同, endPosSize 除开根节点都为1，根节点为0.
	// 本质不同的子串在不同位置出现算不同, 保持现有的 endPosSize, 根节点为0.
	solve := func(s string, k int, unique bool) (start, end int32, ok bool) {
		sam := NewSuffixAutomaton()
		for _, c := range s {
			sam.Add(c)
		}
		size := sam.Size()
		dfsOrder := sam.GetDfsOrder()
		endPosSize := sam.GetEndPosSize(dfsOrder)
		if unique {
			for i := int32(1); i < size; i++ {
				endPosSize[i] = 1
			}
		}
		endPosSize[0] = 0

		samSubSize := make([]int, size) // 每个sam结点往后包含的子串个数.
		for i := size - 1; i >= 0; i-- {
			cur := dfsOrder[i]
			samSubSize[cur] = int(endPosSize[cur])
			nexts := &sam.Nodes[cur].Next
			for j := int32(0); j < SIGMA; j++ {
				if nexts[j] != -1 {
					samSubSize[cur] += samSubSize[nexts[j]]
				}
			}
		}

		remain := k
		if remain > samSubSize[0] {
			return
		}

		pos := int32(0)
		length := int32(0) // SAM上转移的长度/子串长度
		for remain > int(endPosSize[pos]) {
			remain -= int(endPosSize[pos])
			length++
			nexts := &sam.Nodes[pos].Next
			for i := int32(0); i < SIGMA; i++ {
				if nexts[i] != -1 {
					if tmp := samSubSize[nexts[i]]; remain > tmp {
						remain -= tmp
					} else {
						pos = nexts[i] // 答案在这个节点下方
						break
					}
				}
			}
		}

		start, end = sam.RecoverSubstring(pos, length)
		ok = true
		return
	}

	var s string
	fmt.Fscan(in, &s)
	var b, k int
	fmt.Fscan(in, &b, &k)
	unique := b == 0
	start, end, ok := solve(s, k, unique)
	if !ok {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, s[start:end])
	}
}

// P6640 [BJOI2020] 封印 (sam+RMQ)
// https://www.luogu.com.cn/problem/P6640
//
// 给定两个字符串s和t，q次查询s[start:end)和t的最长公共子串长度.
//
// !1. 对t建SAM，求出s的每个前缀与t的最长公共后缀长度lcs[i].
// !2. 对每个询问，答案为 `max(min(lcs[i], i-start+1) for i in range(start, end))`，不好处理.
// !3. 考虑二分答案长度mid，则只需要判定`[start+mid-1,end)`区间内的lcs[i]最大值是否不小于mid即可。
func P6640() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	fmt.Fscan(in, &s, &t)
	sam := NewSuffixAutomaton()
	for _, c := range t {
		sam.Add(c)
	}

	// s的每个前缀s[:i+1]与t的最长公共后缀长度.
	lcs := sam.LongestCommonSuffix(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	rmq := NewSparseTableFrom(lcs, func() int32 { return 0 }, max32)
	query := func(start, end int32) (res int32) {
		// 暴力：
		// for i := start; i < end; i++ {
		// 	res = max32(res, min32(lcs[i], i-start+1))
		// }

		check := func(mid int32) bool {
			return rmq.Query(int(start+mid-1), int(end)) >= mid
		}
		left, right := int32(1), end-start
		for left <= right {
			mid := (left + right) / 2
			if check(mid) {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return int32(right)
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var start, end int32
		fmt.Fscan(in, &start, &end)
		start--
		fmt.Fprintln(out, query(start, end))
	}
}

// Martian Strings (火星字符串)
// https://www.luogu.com.cn/problem/CF149E
// 给定一个主串s和q个模式串.
// 问：对每个模式串pi，问主串中是否存在两个不相交的非空字串，拼起来和模式串相同。
//
// !需要知道，模式串的每个前缀在主串中的最小结束位置和后缀在主串中的最大起始位置.
// 考虑对s的正串和反串分别建后缀自动机.
// 正串sam的结点维护子串结束的最小位置，反串sam的结点维护子串起始的最大位置.
// 对每个模式串，将其前缀在正串上匹配记录最小值，后缀在反串上匹配记录最大值.
// 枚举分割点，如果两个区间的最小值和最大值都不为-1且最小值小于等于最大值，则找到了答案.
func cf149e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int32 = 1e9 + 10

	var s string
	fmt.Fscan(in, &s)
	var q int32
	fmt.Fscan(in, &q)

	n := int32(len(s))
	sam1 := NewSuffixAutomaton() // 正串sam
	for _, c := range s {
		sam1.Add(c)
	}
	minEnd := make([]int32, sam1.Size()) // 正串sam的结点维护子串结束的最小索引(包含)
	for i := range minEnd {
		minEnd[i] = INF
	}
	dfsOrder1 := sam1.GetDfsOrder()
	for i := sam1.Size() - 1; i >= 1; i-- {
		cur := dfsOrder1[i]
		minEnd[cur] = abs32(sam1.Nodes[cur].End)
		pre := sam1.Nodes[cur].Link
		minEnd[pre] = min32(minEnd[pre], minEnd[cur])
	}

	sam2 := NewSuffixAutomaton() // 反串sam
	for i := len(s) - 1; i >= 0; i-- {
		sam2.Add(int32(s[i]))
	}
	maxStart := make([]int32, sam2.Size()) // 反串sam的结点维护子串起始的最大索引(包含)
	for i := range maxStart {
		maxStart[i] = -1
	}
	dfsOrder2 := sam2.GetDfsOrder()
	for i := sam2.Size() - 1; i >= 1; i-- {
		cur := dfsOrder2[i]
		maxStart[cur] = n - 1 - abs32(sam2.Nodes[cur].End)
		pre := sam2.Nodes[cur].Link
		maxStart[pre] = max32(maxStart[pre], maxStart[cur])
	}

	query := func(pattern string) bool {
		m := int32(len(pattern))
		preMinEnd := make([]int32, m+1)
		sufMaxStart := make([]int32, m+1)
		for i := range preMinEnd {
			preMinEnd[i] = INF
			sufMaxStart[i] = -1
		}

		pos1 := int32(0)
		for i := int32(0); i < m; i++ {
			if pos1 = sam1.Nodes[pos1].Next[int32(pattern[i])-OFFSET]; pos1 != -1 {
				preMinEnd[i+1] = minEnd[pos1]
			} else {
				break
			}
		}
		pos2 := int32(0)
		for i := m - 1; i >= 0; i-- {
			if pos2 = sam2.Nodes[pos2].Next[int32(pattern[i])-OFFSET]; pos2 != -1 {
				sufMaxStart[i+1] = maxStart[pos2]
			} else {
				break
			}
		}

		for i := int32(1); i < m; i++ {
			if preMinEnd[i] != INF && sufMaxStart[i+1] != -1 && preMinEnd[i] < sufMaxStart[i+1] {
				return true
			}
		}
		return false
	}

	res := int32(0)
	for i := int32(0); i < q; i++ {
		var pattern string
		fmt.Fscan(in, &pattern)
		if query(pattern) {
			res++
		}
	}
	fmt.Fprintln(out, res)
}

// Cyclical Quest
// https://www.luogu.com.cn/problem/CF235C
// https://www.cnblogs.com/h-lka/p/15169021.html
// !给定一个主串S和n个询问串，求每个询问串的所有循环同构"去重"后在主串中出现的次数总和。
//
// !循环就是把询问串第一个字符拿去放在最后面(目前匹配到的串的第一个字符删掉,然后再加上一个)。
// !在前面删除字符类似 parent tree 上跳的操作，后面加字符类似SAM上的转移操作。
// 如果同一个询问串多次匹配到同一个节点,贡献只算一次,具体可以打标记实现。
func cf235c() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var text string
	fmt.Fscan(in, &text)
	var q int
	fmt.Fscan(in, &q)
	words := make([]string, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &words[i])
	}

	sam := NewSuffixAutomaton()
	for _, c := range text {
		sam.Add(c)
	}

	endPosSize := sam.GetEndPosSize(sam.GetDfsOrder())
	size := sam.Size()
	visitedTime := make([]int32, size)
	for i := range visitedTime {
		visitedTime[i] = -1
	}
	for i, word := range words {
		m := int32(len(word))
		pos, len_ := int32(0), int32(0)
		res := 0
		for j := int32(0); j < m*2; j++ {
			var v int32
			if j < m {
				v = int32(word[j])
			} else {
				v = int32(word[j-m])
			}

			pos, len_ = sam.Move(pos, len_, v)
			if j >= m { // 移除模式串首部字符
				pos, len_ = sam.MoveLeft(pos, len_, m+1)
			}
			if j >= m-1 {
				if len_ == m && visitedTime[pos] < int32(i) {
					visitedTime[pos] = int32(i)
					res += int(endPosSize[pos])
				}
			}
		}

		fmt.Fprintln(out, res)
	}
}

// Paper task (本质不同的合法括号子串计数)
// https://www.luogu.com.cn/problem/CF653F
// 给定一个长度为n的括号串，问有多少种不同的合法括号子串
// n<=5e5.
//
// 参考 https://atcoder.jp/contests/abc223/tasks/abc223_f
// !令(为1，)为-1，括号序列合法当且仅当
// - 和为0.
// - "后缀和"都不大于0 => 后缀和的最小值不大于0.
//
// 问题转化为：每个endPos上，end 固定时，求合法的start.
// 先用最小值<=0条件二分出最左端的start起点，再统计使得[start:end)和为0的start的个数。
func cf653f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF32 int32 = 1e9 + 10
	var n int32
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	sam := NewSuffixAutomaton()
	for _, c := range s {
		if c == '(' {
			sam.Add('a')
		} else {
			sam.Add('b')
		}
	}
	preSum := make([]int32, n+1)
	for i, c := range s {
		if c == '(' {
			preSum[i+1] = preSum[i] + 1
		} else {
			preSum[i+1] = preSum[i] - 1
		}
	}

	mp := make(map[int32][]int32) // group by preSum
	for i := int32(0); i < n+1; i++ {
		mp[preSum[i]] = append(mp[preSum[i]], i)
	}

	minSt := NewSparseTable(len(preSum), func(i int) S { return preSum[i] }, func() S { return INF32 }, func(a, b S) S { return min32(a, b) })

	// 求sortedArray中[min,max]范围内的元素个数.
	rangeCount := func(sortedArray []int32, min, max int32) int32 {
		if min > max {
			return 0
		}
		a := sort.Search(len(sortedArray), func(i int) bool { return sortedArray[i] >= min })
		b := sort.Search(len(sortedArray), func(i int) bool { return sortedArray[i] > max })
		return int32(b - a)
	}

	_ = rangeCount
	cal := func(pos int32) int32 {
		link := sam.Nodes[pos].Link
		minLen, maxLen := sam.Nodes[link].MaxLen+1, sam.Nodes[pos].MaxLen
		end := abs32(sam.Nodes[pos].End) + 1 // 每个endPos随意选取一个结束位置
		minStart := end - maxLen
		_ = maxLen
		maxStart := end - minLen

		// end固定时，求合法的start范围.
		minLeft := int32(minSt.MinLeft(int(end+1), func(e S) bool { return e <= 0 }))
		res := int32(0)
		for i := max32(minLeft, minStart); i <= maxStart; i++ {
			curSum := 0
			for j := i; j < end; j++ {
				if s[j] == '(' {
					curSum++
				} else {
					curSum--
				}
			}
			if curSum == 0 {
				res++
			}
		}
		fmt.Println(pos, minLeft, end, res, 999)
		return res

		// sum := preSum[end]
		// return rangeCount(mp[sum], max32(minLeft, minStart), maxStart)
	}

	res := 0
	for i := int32(1); i < sam.Size(); i++ {
		res += int(cal(i))
	}
	fmt.Fprintln(out, res)
}

// Fake News (hard)
// https://www.luogu.com.cn/problem/CF802I
// 求本质不同子串出现次数的平方和
func cf802I() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func() {
		var s string
		fmt.Fscan(in, &s)
		sam := NewSuffixAutomaton()
		for _, c := range s {
			sam.Add(c)
		}
		dfsOrder := sam.GetDfsOrder()
		endPosSize := sam.GetEndPosSize(dfsOrder)
		res := 0
		for i := int32(1); i < sam.Size(); i++ {
			size, length := int(endPosSize[i]), int(sam.DistinctSubstringAt(i))
			res += size * size * length
		}
		fmt.Fprintln(out, res)
	}

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		solve()
	}
}

// Security (线段树合并)
// https://www.luogu.com.cn/problem/CF1037H
// https://www.cnblogs.com/Troverld/p/14605723.html
//
// 在后缀自动机的 DAWG 上贪心。使用线段树合并判断当前字符串是否作为[l,r]的子串出现过 。
// https://codeforces.com/contest/1037/submission/147520554
func CF1037H() {}

// Yet Another LCP Problem
// https://www.luogu.com.cn/problem/CF1073G
// https://www.cnblogs.com/Troverld/p/14605733.html
// 给定一个长为n的字符串s和q次询问.
// 每次询问给出两个数组A和B，求两两后缀最长公共前缀之和 ∑lcp(s[A[i]:], s[B[j]:]).
func cf1073g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
}

// Indie Album
// https://www.luogu.com.cn/problem/CF1207G
// https://www.cnblogs.com/Troverld/p/14605729.html
func cf1207g() {}

// P3804 【模板】后缀自动机（SAM）
// https://www.luogu.com.cn/problem/P3804
// 给定一个长度为 n 的只包含小写字母的字符串 s。
// !对于所有 s 的出现次数不为 1 的子串，设其 value值为`该子串出现的次数 × 该子串的长度`。
// 请计算，value 的最大值是多少。
// n <= 1e6
//
// !一个子串必然是一个后缀的前缀，所以这个子串的出现次数，就是子树中实点的个数.
func P3804() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	if len(s) <= 1 {
		fmt.Fprintln(out, 0)
		return
	}

	sam := NewSuffixAutomaton()
	for _, c := range s {
		sam.Add(c)
	}

	dfsOrder := sam.GetDfsOrder()
	endPosSize := sam.GetEndPosSize(dfsOrder)
	res := 0
	for i := int32(1); i < sam.Size(); i++ {
		if endPosSize[i] > 1 {
			res = max(res, int(sam.Nodes[i].MaxLen)*int(endPosSize[i]))
		}
	}
	fmt.Fprintln(out, res)
}

// https://judge.yosupo.jp/problem/number_of_substrings
func number_of_substrings() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	sa := NewSuffixAutomaton()
	for _, c := range s {
		sa.Add(c)
	}
	fmt.Fprintln(out, sa.DistinctSubstring())
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

	sam := NewSuffixAutomaton()
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

// 葫芦的考验之定位子串
// https://ac.nowcoder.com/acm/contest/37092/C
// q次查询子串s[start:end)的出现次数.
func nowCoder37092C() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var q int32
	fmt.Fscan(in, &q)

	sam := NewSuffixAutomaton()
	prefixEnd := make([]int32, len(s))
	for i, c := range s {
		pos := sam.Add(c)
		prefixEnd[i] = pos
	}
	dfsOrder := sam.GetDfsOrder()
	endPosSize := sam.GetEndPosSize(dfsOrder)

	query := func(start, end int32) int32 {
		pos := sam.LocateSubstring(start, end, prefixEnd[end-1])
		return endPosSize[pos]
	}

	for i := int32(0); i < q; i++ {
		var start, end int32
		fmt.Fscan(in, &start, &end)
		start--
		fmt.Fprintln(out, query(start, end))
	}
}

// Typewriter (后缀自动机优化dp)
// https://ac.nowcoder.com/acm/contest/37092/D
// https://blog.csdn.net/swustzhaoxingda/article/details/98024818
// 初始时有一个空字符串，有两种操作：
// 1. 在字符串末尾添加一个字符，代价为cost1；
// 2. 复制一个已经出现的子串到末尾，代价为cost2。
// 问最小代价使得字符串变为目标字符串s。
//
// 考虑配るdp，dp[i]表示变为s[:i]的最小代价.
// 使用操作1，有 dp[i+1] = dp[i] + cost1。
// 使用操作2，有 dp[i+len] = min(dp[i+len], dp[i]+cost2)，
// 其中len是一个最大的长度使得s[i:i+len]是s[:i]的一个子串。
// 二分len，然后判断s[i:i+len]是否是前缀s[:i]的一个子串即可。
// 时间复杂度O(nlogn^2)
func nowCoder37092D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int = 1e18

	var s string
	fmt.Fscan(in, &s)
	var cost1, cost2 int
	fmt.Fscan(in, &cost1, &cost2)

	n := int32(len(s))
	sam := NewSuffixAutomaton()
	prefixEnd := make([]int32, n)
	for i, c := range s {
		pos := sam.Add(c)
		prefixEnd[i] = pos
	}
	endPosMinEnd := sam.GetEndPosMinEnd(sam.GetDfsOrder())

	// 判断s[start:end]是否是前缀s[:prefix]的一个子串.
	isSubstringOfPrefix := func(start, end int32, prefix int32) bool {
		if end-start > prefix {
			return false
		}
		pos := sam.LocateSubstring(start, end, prefixEnd[end-1])
		return endPosMinEnd[pos]+1 <= prefix
	}

	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0

	for i := int32(0); i < n; i++ {
		dp[i+1] = min(dp[i+1], dp[i]+cost1) // 操作1

		ok := false
		left, right := int32(1), min32(n-i, i)
		for left <= right {
			mid := (left + right) / 2
			if isSubstringOfPrefix(i, i+mid, i) {
				left = mid + 1
				ok = true
			} else {
				right = mid - 1
			}
		}
		if ok {
			dp[i+right] = min(dp[i+right], dp[i]+cost2) // 操作2
		}
	}

	fmt.Fprintln(out, dp[n])
}

// 葫芦的考验之定位子串2.0
// https://ac.nowcoder.com/acm/contest/37092/E
// https://www.luogu.com.cn/problem/P4094 二分后转化为本题
// q次查询子串s[start1:end1)在子串s[start2:end2)中出现的次数.
//
// 先定位子串s[start1:end1)的endPos结点，这个结点子树都包含了这个子串。
// 转化成dfs序后waveletMatrix查询即可.
func nowCoder37092E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var q int32
	fmt.Fscan(in, &q)

	sam := NewSuffixAutomaton()
	prefixEnd := make([]int32, len(s))
	for i, c := range s {
		pos := sam.Add(c)
		prefixEnd[i] = pos
	}

	failTree := sam.BuildTree()
	size := sam.Size()
	lid, rid := make([]int32, size), make([]int32, size)
	dfn := int32(0)
	var dfs func(cur int32)
	dfs = func(cur int32) {
		lid[cur] = dfn
		dfn++
		for _, next := range failTree[cur] {
			dfs(next)
		}
		rid[cur] = dfn
	}
	dfs(0)

	data := make([]int32, size)
	for i, pos := range prefixEnd {
		data[lid[pos]] = int32(i + 1)
	}
	wm := NewWaveletMatrix32(size, func(i int32) int32 { return data[i] })

	// CountSubstringInRange
	query := func(start1, end1 int32, start2, end2 int32) int32 {
		if (end1 - start1) > (end2 - start2) {
			return 0
		}
		pos := sam.LocateSubstring(start1, end1, prefixEnd[end1-1])
		len_ := end1 - start1
		return wm.CountRange(lid[pos], rid[pos], 1+(start2+len_-1), 1+end2)
	}

	for i := int32(0); i < q; i++ {
		var start1, end1, start2, end2 int32
		fmt.Fscan(in, &start1, &end1, &start2, &end2)
		start1--
		start2--
		fmt.Fprintln(out, query(start1, end1, start2, end2))
	}
}

// AStringGame (dag+grundy数)
// https://ac.nowcoder.com/acm/contest/37092/J
//
// 用sam跑出DAG图，再求出sg值。
// SG值异或和等于0时先手必输，否则先手必胜。
func nowCoder37092J() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(text string, patterns []string) bool {
		sam := NewSuffixAutomaton()
		for _, c := range text {
			sam.Add(c)
		}

		memo := make([]int32, sam.Size())
		for i := range memo {
			memo[i] = -1
		}

		var grundy func(int32) int32
		grundy = func(state int32) int32 {
			if memo[state] != -1 {
				return memo[state]
			}

			nextStates := make(map[int32]struct{})
			for _, next := range sam.Nodes[state].Next {
				if next != -1 {
					nextStates[grundy(next)] = struct{}{}
				}
			}

			mex := int32(0)
			for {
				if _, ok := nextStates[mex]; !ok {
					break
				}
				mex++
			}
			memo[state] = mex
			return mex
		}

		sg := int32(0) // 初始状态的sg值为转移状态的sg值的异或和
		for _, pattern := range patterns {
			pos := int32(0)
			for _, c := range pattern {
				pos = sam.Nodes[pos].Next[c-OFFSET]
				if pos == -1 {
					break
				}
			}

			if pos != -1 {
				sg ^= grundy(pos)
			}
		}

		return sg != 0
	}

	for {
		var text string
		fmt.Fscan(in, &text)
		// 判断是否还有输入
		if len(text) == 0 {
			break
		}
		var n int32
		fmt.Fscan(in, &n)
		patterns := make([]string, n)
		for i := int32(0); i < n; i++ {
			fmt.Fscan(in, &patterns[i])
		}

		firstWin := solve(text, patterns)
		if firstWin {
			fmt.Fprintln(out, "Alice")
		} else {
			fmt.Fprintln(out, "Bob")
		}
	}
}

func testMerge() {
	s := "aababa"
	sam := NewSuffixAutomaton()
	for _, c := range s {
		sam.Add(c)
	}
	dfsOrder := sam.GetDfsOrder()
	seg, endPos := sam.GetEndPos(dfsOrder)
	for i, node := range endPos {
		var pos []int32
		seg.Enumerate(node, func(j int32, _ int32) {
			pos = append(pos, j)
		})
		fmt.Println(i, pos)
	}
	fmt.Println(sam.GetEndPos(dfsOrder))
}
