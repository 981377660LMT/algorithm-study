// https://maspypy.github.io/library/string/suffix_automaton.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/sam.go
// https://www.bilibili.com/video/BV1S54y1G7P8
// https://www.cnblogs.com/Linshey/p/14219867.html
// Blumber 算法在线构建SAM
//
// !子串是一个前缀的后缀.
//
// -s: aababa
//
// -fail tree
//
//     link 指向比当前结点短的最长后缀endPos集.
//     每个节点中的子串，都是以该节点为根的子树中的所有子串的"后缀"
//
//					      0,1,2,3,4,5 ""
//					     /              \
//					    /					  	   \
//					   /					  	    \
//					  /						  	  	 \
//					0,1,3,5 "a"				    2,4 "ab"
//					 /	 \	               /   "b"
//					/		  \						    /		  \
//				 /       \             /       \
//				/         \           /         \
//			1,"aa"   3,5 "aba"	   2 "aab"   4 "aabab"
//                 /  "ba"                  "abab"
//				        /      \                   "bab"
//				       /        \
//			       3 "aaba"  5 "aababa"
//                          "ababa"
//                           "baba"
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
//  5. endPos集合要么包含要么不相交
//  6. 一个子串出现次数就是其对应 endPos 集合的size(注意不是长度范围).
//    !由于子串<=>前缀的后缀，
//    !可以先通过在 SAM 上找到该子串所处的节点，然后求以该节点为根的子树中，有多少个包含原串前缀的节点
//    !另一个含义——从SAM的根到这个结点的转移路径条数。
//  7. 可以把SAM理解为把某个串的所有子串建立AC自动机。
//  !8. 设 lcs(i,j) 为前缀i,j的最长公共后缀长度，其等于fail树上 LCA 的len 值。
//  9. 一个endpos等价类内的串的长度连续.
//  10.理解
//     - 从 SAM 的定义上理解：
//       SAM 可以看作一种加强版的 Trie，它可以高度压缩一个字符串的子串信息，
//       !一条从根出发到`终止结点`的路径对应了原串的一个后缀，而任意一个从根出发的路径对应了原串一个子串。
//       子串和从根出发的路径一一对应。在这种的理解下，每一个结点的含义并不是固定的，
//       它到底对应哪个子串取决于那条路径是怎么到达它的；而边有着确定的含义。
//     - 从 Parent Tree 的角度去理解连边的含义
//       两个不同等价类的Endpos集合要么无交集，要么相包含，因此可以建出一个由 Endpos集合的包含关系连结而成的树——Parent Tree
//       它的连边——后缀链接，若是向下看，是在一个等价类的前面加上一个字符，从而分成若干的其他等价类；
//       向上看，它是指向包含当前集合的最小的集合。
//       !而后缀自动机的连边是在一个等价类的后面加上一个字母，看看它会指向谁，显然对于同一个添上的字母，这个指向是唯一确定的。
//     - 从结点的含义去理解：
//       每一个结点都对应了一种子串，Parent Tree 的结点与 SAM 的结点一一对应
//  !11.转移边：parent树往下走代表往前加字符，SAM转移边往后走代表往后加字符
//  !12. 从SAM的DAG角度看，子串是后缀的一个前缀；
//	     !从SAM的Parent Tree角度看，子串是前缀的一个后缀。
//  !13. AC自动机的失配链接和后缀自动机的后缀链接都有性质：指向的两个状态都满足"后者的代表串是前者的代表串的真后缀"。
//
// applications:
//  1. 查找某个子串位于哪个节点 => 直接倍增往上跳到len[]合适的地方
//  2. 最长可重叠重复子串 => endPos集合大于等于2的那些节点的最大的范围
//  3. `在线`给出模式串的模式匹配问题(单模式串离线=>KMP，多模式串离线=>AC自动机，多模式串在线=>SAM)
//     一般有固定模式串的字符串处理问题和固定主串的字符串处理问题两大类问题。当固定模式串时，熟知的 AC 自动机算法便可以胜任这类问题。
//     如果主串固定，一般采用对主串构造后缀树、后缀自动机来解决这一类问题。
//  4. 对于SAM任何一个节点u，从根到这个节点的路线有 `maxLen(u)-minLen(u)+1` 条，而这条路线则表示原字符串的一个子串，且各不相同.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	P3804()
}

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

	endPosSize := sam.GetEndPosSize()
	res := 0
	for i := int32(1); i < sam.Size(); i++ {
		if endPosSize[i] > 1 {
			res = max(res, int(sam.Nodes[i].MaxLen)*int(endPosSize[i]))
		}
	}
	fmt.Fprintln(out, res)
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
	Next   [SIGMA]int32 // 孩子节点
	Link   int32        // 后缀链接
	MaxLen int32        // 当前节点对应的最长子串的长度
	End    int32        // 字符在原串的下标, 实点下标为非负数, 虚点下标为负数
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
func (sam *SuffixAutomaton) GetEndPosSize() []int32 {
	size := sam.Size()
	endPosSize := make([]int32, size)
	dfsOrder := sam.GetDfsOrder()
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
func (sam *SuffixAutomaton) GetNodeBySubstring(start, end int32) {}

// pos 位置对应的子串个数.
// 用最长串的长度减去最短串的长度即可得到以当前节点为结尾的子串个数.
// 最长串的长度记录在节点的 MaxLength 中,最短串的长度可以通过link对应的节点的 MaxLength 加 1 得到.
func (sam *SuffixAutomaton) CountSubstringAt(pos int32) int32 {
	if pos == 0 {
		return 0
	}
	return sam.Nodes[pos].MaxLen - sam.Nodes[sam.Nodes[pos].Link].MaxLen
}

// 本质不同的子串个数.
func (sam *SuffixAutomaton) CountSubstring() int {
	res := 0
	for i := 0; i < len(sam.Nodes); i++ {
		res += int(sam.CountSubstringAt(int32(i)))
	}
	return res
}

func (sam *SuffixAutomaton) newNode(link, maxLength, end int32) *Node {
	res := &Node{Link: link, MaxLen: maxLength, End: end}
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
