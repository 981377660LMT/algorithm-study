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
//	    link 指向比当前结点短的最长后缀endPos集.
//	    每个节点中的子串，都是以该节点为根的子树中的所有子串的"后缀"
//
//						      0,1,2,3,4,5 ""
//						     /              \
//						    /					  	   \
//						   /					  	    \
//						  /						  	  	 \
//						0,1,3,5 "a"(后)				 2,4 "b"
//						 /	 \	                / "ab"
//						/		  \						     /		  \
//					 /       \              /        \
//					/         \            /          \
//				1,"aa"   3,5 "ba"(后) 	2 "aab"      4 "bab"
//	                / "aba"                     "abab"
//					       /      \                    "aabab"
//					      /        \
//				       3 "aaba"   5 "baba"(后)
//	                         "ababa"
//	                        "aababa"
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
//     !8. 设 lcs(i,j) 为前缀i,j的最长公共后缀长度，其等于fail树上 LCA 的len 值。
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
//     !11. 转移边：parent树往下走代表往前加字符，SAM转移边往后走代表往后加字符
//     !12. 子串是什么：
//          从SAM的DAG角度看，子串是后缀的一个前缀；
//          !从SAM的Parent Tree角度看，子串是前缀的一个后缀。
//     !13. SAM 与AC自动机的相似性：
//     AC自动机的失配链接和后缀自动机的后缀链接都有性质：
//     指向的两个状态都满足"后者的代表串是前者的代表串的真后缀"。
//     可以把 SAM 理解为把某个串的所有子串建立AC自动机.
//  14. 增量构造中，每次从后面加入一个字符, 有两件事要干：
//     找出能转移到这个状态的状态，建立链接；确定这个状态的min，即找到它在parent树上的父亲。
//  15. 对于SAM任何一个节点u，从根到这个节点的路线有 `maxLen(u)-minLen(u)+1` 条，而这条路线则表示原字符串的一个子串，且各不相同.
//  16. 一般来讲,DAG上可能重复转移,是很难跑计数DP的。
//     !但是我们知道后缀自动机的性质 : 任意两个节点的表示集合没有交。
//     !所以我们只要统计路径数即可,不需要考虑重复问题。
//  !17.可以通过parent树确定SAM的接受状态集合。找到MaxLen=n的结点，该结点到根的路径上的所有结点都是接受状态。
//  !18.子串是
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

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const SIGMA int32 = 26   // 字符集大小
const OFFSET int32 = 'a' // 字符集的起始字符

type Node struct {
	Next   [SIGMA]int32 // SAM 转移边
	Link   int32        // 后缀链接
	MaxLen int32        // 当前节点对应的最长子串的长度
	End    int32        // 最长的字符在原串的下标, 实点下标为非负数, 虚点下标为负数
}

type SuffixAutomaton struct {
	Nodes   []*Node
	LastPos int32 // 当前插入的字符对应的节点(实点，原串的一个前缀)
	n       int32 // 当前字符串长度
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

// TODO: 线段树合并维护 EndPos 集合
func (sam *SuffixAutomaton) GetEndPos() {}

// TODO: 快速定位子串
// 倍增往上跳到len[]合适的地方
func (sam *SuffixAutomaton) GetNodeBySubstring(start, end int32) {}

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

	// cf235c()
	// cf802I()

	// number_of_substrings()
	// longest_common_substring()

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

// Martian Strings
// https://www.luogu.com.cn/problem/CF149E
// 可以找到两个不相交的区间，满足这两个区间对应的子串拼起来和 wi相同
func cf149e() {}

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
