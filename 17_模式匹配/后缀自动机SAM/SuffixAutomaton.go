// https://maspypy.github.io/library/string/suffix_automaton.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/sam.go
// https://yutong.site/sam/ 可视化
// Blumber 算法在线构建SAM
//
// -s: aababa
//
// -fail tree
//
//		 每个节点后面表的是该类中最长的串。
//		 父亲是孩子的最长后缀( endPos 不同)，反串的parent树就满足 : 父亲是孩子的最长前缀 ( 不同 )。
//	   父亲的endPos是子节点的endPos求并(除开0).
//     每个节点中的子串，都是以该节点为根的子树中的所有子串的后缀
//     叶子节点一定包含原串的前缀（根据性质 1 反证可得），但原串的前缀不一定在叶子节点中（比如前缀 A 是前缀 AA 的后缀）
//     由于【子串】等价于【前缀的后缀】，因此求子串在原串中的出现次数，
//     可以先通过在 SAM 上找到该子串所处的节点，然后求以该节点为根的子树中，多有少个包含原串前缀的节点（性质 2）
//     这可以通过在 rev 上统计子树信息来预处理得到
//
//					      0,1,2,3,4,5 ""
//					     /              \
//					    /					  	   \
//					   /					  	    \
//					  /						  	  	 \
//					0,1,3,5 "a"				    2,4 "ab"
//					/		  \						    /		  \
//				 /       \             /       \
//				/         \           /         \
//			1,"aa"   3,5 "aba"	   2 "aab"   4 "aabab"
//				        /      \
//				       /        \
//			       3 "aaba"  5 "aababa"
//
// note:
//  0. 后缀自动机 (Suffix Automaton, SAM) 是仅接受后缀且状态数最少的 DFA
//  1. 每一个节点都表示一段子串，所有节点表示的子串们都是唯一的
//  2. len(cur)表示的是当前节点cur的最长长度，当前节点的子串长度范围是 [len(link)+1, len(cur)]
//  3. endPos 集合的大小可以通过topo排序求出来，实际上用桶排实现
//     如果必须要求出 endPos 集合的话，可以用set实现树上自底向上启发式合并
//     如果需要每个点的endPos集合都需要求出来的话，可以用动态开点线段树维护endPos集合，然后使用线段树的合并进行更新
//  4. 子节点最短串的最长后缀=父结点最长串
//  5. endPos集合要么包含要么不相交
//  6. 一个子串出现次数就是其对应 endPos 集合的元素个数。
//  7. 可以把SAM理解为把某个串的所有子串建立AC自动机。
//  8. 设 lcs(i,j) 为前缀i,j的最长公共后缀长度，其等于fail树上 LCA 的len 值。
//
// applications:
//  1. 查找某个子串位于哪个节点 => 直接倍增往上跳到len[]合适的地方
//  2. 最长可重叠重复子串 => endPos集合大于等于2的那些节点的最大的范围
//  3. `在线`给出模式串的模式匹配问题(单模式串离线=>KMP，多模式串离线=>AC自动机，多模式串在线=>SAM)
//     一般有固定模式串的字符串处理问题和固定主串的字符串处理问题两大类问题。当固定模式串时，熟知的 AC 自动机算法便可以胜任这类问题。
//     如果主串固定，一般采用对主串构造后缀树、后缀自动机来解决这一类问题。
//  4. 对于SAM任何一个节点u，从根到这个节点的路线有 `maxLen(u)-minLen(u)+1` 条，而这条路线则表示原字符串的一个子串，且各不相同.
//  5. 一个SAM 最多有 2n-1 个节点 和 3n-4 条转移边，最坏情况下为 abbb...bbb
//  6. 两个不同的 endPos 等价的字符串中，较短者总是较长者的真后缀。因此，等价类中没有等长的字符串。
//  7. 出现次数 => 在自动机上查找模式串 P 对应的节点，如果存在，则答案就是该节点的终点集合大小.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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
	Next [SIGMA]int32 // 转移边
	Link int32        // 后缀链接

	// 该节点（endpos 等价类）中最长的子串长度, link.MaxLen+1 为该节点（endpos 等价类）中最短的子串长度
	// 等价类大小为 MaxLen-link.MaxLen
	MaxLen int32
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

// 每次插入内部最多创建2个状态.
func (sa *SuffixAutomaton) Add(ord int32) {
	c := ord - OFFSET
	newNode := int32(len(sa.Nodes))
	sa.Nodes = append(sa.Nodes, sa.newNode(-1, sa.Nodes[sa.CurPos].MaxLen+1))
	p := sa.CurPos
	for p != -1 && sa.Nodes[p].Next[c] == -1 {
		sa.Nodes[p].Next[c] = newNode
		p = sa.Nodes[p].Link
	}
	q := int32(0)
	if p != -1 {
		q = sa.Nodes[p].Next[c]
	}
	if p == -1 || sa.Nodes[p].MaxLen+1 == sa.Nodes[q].MaxLen {
		sa.Nodes[newNode].Link = q
	} else {
		newQ := int32(len(sa.Nodes))
		sa.Nodes = append(sa.Nodes, sa.newNode(sa.Nodes[q].Link, sa.Nodes[p].MaxLen+1))
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
	for v := int32(1); v < n; v++ {
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

// TODO，类似kmp移动.
func (sa *SuffixAutomaton) Move(pos int32, input int32) int32 {}

// pos 位置对应的子串个数.
// 用最长串的长度减去最短串的长度即可得到以当前节点为结尾的子串个数.
// 最长串的长度记录在节点的 MaxLength 中,最短串的长度可以通过link对应的节点的 MaxLength 加 1 得到.
func (sa *SuffixAutomaton) CountSubstringAt(pos int32) int32 {
	if pos == 0 {
		return 0
	}
	return sa.Nodes[pos].MaxLen - sa.Nodes[sa.Nodes[pos].Link].MaxLen
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
	res := &Node{Link: link, MaxLen: maxLength}
	for i := int32(0); i < SIGMA; i++ {
		res.Next[i] = -1
	}
	return res
}
