// https://www.luogu.com.cn/blog/yszs/ac-zi-dong-ji-fou-guo-shi-jian-fail-shu-di-gong-ju-pi-liao

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	SeparateString()
}

func demo() {
	acm := NewACAutoMatonArray(26, 97)
	acm.AddString("abc")
	acm.AddString("bc")
	acm.BuildSuffixLink(false)
	fmt.Println(acm.GetCounter())
}

const INF int = 1e18

// 2781. 最长合法子字符串的长度
// https://leetcode.cn/problems/length-of-the-longest-valid-substring/
// 给你一个字符串 word 和一个字符串数组 forbidden 。
// 如果一个字符串不包含 forbidden 中的任何字符串，我们称这个字符串是 合法 的。
// 请你返回字符串 word 的一个 最长合法子字符串 的长度。
// 子字符串 指的是一个字符串中一段连续的字符，它可以为空。
//
// 1 <= word.length <= 1e5
// word 只包含小写英文字母。
// 1 <= forbidden.length <= 1e5
// !1 <= forbidden[i].length <= 1e5
// !sum(len(forbidden)) <= 1e7
// forbidden[i] 只包含小写英文字母。
//
// 思路:
// 类似字符流, 需要处理出每个位置为结束字符的包含至少一个模式串的`最短后缀`.
// !那么此时左端点就对应这个位置+1
func longestValidSubstring(word string, forbidden []string) int {
	acm := NewACAutoMatonArray(26, 97)
	for _, w := range forbidden {
		acm.AddString(w)
	}
	acm.BuildSuffixLink(false)

	minLen := make([]int, acm.Size()) // 每个状态匹配到的模式串的最小长度
	for i := range minLen {
		minLen[i] = INF
	}
	for i, pos := range acm.Words {
		minLen[pos] = min(minLen[pos], len(forbidden[i]))
	}
	acm.Dp(func(from, to int) { minLen[to] = min(minLen[to], minLen[from]) })

	res, left, pos := 0, 0, 0
	for right, char := range word {
		pos = acm.Move(pos, int(char))
		left = max(left, right-minLen[pos]+2)
		res = max(res, right-left+1)
	}
	return res
}

// 一个文本串，可以将一个字符改成*，
// 保证都是小写字母组成。问最少改多少次，可以完全不包含若干模式串。
// !AC自动机(多模式串匹配)
// https://atcoder.jp/contests/abc268/submissions/34752897
func SeparateString() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	acm := NewACAutoMatonArray(26, 97)
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(false)
	accept := acm.GetCounter()
	pos := 0
	res := 0
	for i := 0; i < len(s); i++ {
		pos = acm.Move(pos, int(s[i]))
		if accept[pos] > 0 {
			res++
			pos = 0
		}
	}

	fmt.Fprintln(out, res)
}

// 1032. 字符流
// https://leetcode.cn/problems/stream-of-characters/description/
type StreamChecker struct {
	*ACAutoMatonArray
	pos   int // state in AhoCorasick
	count []int
}

// words = ["abc", "xyz"] 且字符流中逐个依次加入 4 个字符 'a'、'x'、'y' 和 'z' ，
// 你所设计的算法应当可以检测到 "axyz" 的后缀 "xyz" 与 words 中的字符串 "xyz" 匹配。
func Constructor(words []string) StreamChecker {
	acm := NewACAutoMatonArray(26, 'a')
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(false)
	return StreamChecker{ACAutoMatonArray: acm, count: acm.GetCounter()}
}

// 从字符流中接收一个新字符，如果字符流中的任一非空后缀能匹配 words 中的某一字符串，返回 true ；否则，返回 false。
func (this *StreamChecker) Query(letter byte) bool {
	this.pos = this.Move(this.pos, int(letter))
	return this.count[this.pos] > 0
}

// https://leetcode.cn/problems/multi-search-lcci/
// 给定一个较长字符串big和一个包含较短字符串的数组smalls，
// 设计一个方法，根据smalls中的每一个较短字符串，对big进行搜索。
// !输出smalls中的字符串在big里出现的所有位置positions，
// 其中positions[i]为smalls[i]出现的所有位置。
func multiSearch(big string, smalls []string) [][]int {
	acm := NewACAutoMatonArray(26, 'a')
	for _, s := range smalls {
		acm.AddString(s)
	}
	acm.BuildSuffixLink(false)

	pos := 0
	res := make([][]int, len(smalls))
	indexes := acm.GetIndexes()
	for i := 0; i < len(big); i++ {
		pos = acm.Move(pos, int(big[i]))
		for _, j := range indexes[pos] {
			res[j] = append(res[j], i-len(smalls[j])+1) // !i-len(smalls[j])+1: 模式串在big中的起始位置
		}
	}
	return res
}

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个字符串.
type ACAutoMatonArray struct {
	Words              []int   // words[i] 表示加入的第i个模式串对应的节点编号.
	Parent             []int   // parent[v] 表示节点v的父节点.
	sigma              int     // 字符集大小.
	offset             int     // 字符集的偏移量.
	children           [][]int // children[v][c] 表示节点v通过字符c转移到的节点.
	suffixLink         []int   // 又叫fail.指向当前节点最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
	bfsOrder           []int   // 结点的拓扑序,0表示虚拟节点.
	needUpdateChildren bool    // 是否需要更新children数组.
}

func NewACAutoMatonArray(sigma, offset int) *ACAutoMatonArray {
	res := &ACAutoMatonArray{sigma: sigma, offset: offset}
	res.newNode()
	return res
}

func (trie *ACAutoMatonArray) AddString(str string) int {
	if len(str) == 0 {
		return 0
	}
	pos := 0
	for _, s := range str {
		ord := int(s) - trie.offset
		if trie.children[pos][ord] == -1 {
			trie.children[pos][ord] = trie.newNode()
			trie.Parent = append(trie.Parent, pos)
		}
		pos = trie.children[pos][ord]
	}
	trie.Words = append(trie.Words, pos)
	return pos
}

func (trie *ACAutoMatonArray) AddChar(pos int, ord int) int {
	ord -= trie.offset
	if trie.children[pos][ord] != -1 {
		return trie.children[pos][ord]
	}
	trie.children[pos][ord] = trie.newNode()
	trie.Parent = append(trie.Parent, pos)
	return trie.children[pos][ord]
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *ACAutoMatonArray) Move(pos int, ord int) int {
	ord -= trie.offset
	if trie.needUpdateChildren {
		return trie.children[pos][ord]
	}
	for {
		nexts := trie.children[pos]
		if nexts[ord] != -1 {
			return nexts[ord]
		}
		if pos == 0 {
			return 0
		}
		pos = trie.suffixLink[pos]
	}
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *ACAutoMatonArray) Size() int {
	return len(trie.children)
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组.
// !设置为false更快.
func (trie *ACAutoMatonArray) BuildSuffixLink(needUpdateChildren bool) {
	trie.needUpdateChildren = needUpdateChildren
	trie.suffixLink = make([]int, len(trie.children))
	for i := range trie.suffixLink {
		trie.suffixLink[i] = -1
	}
	trie.bfsOrder = make([]int, len(trie.children))
	head, tail := 0, 0
	trie.bfsOrder[tail] = 0
	tail++
	for head < tail {
		v := trie.bfsOrder[head]
		head++
		for i, next := range trie.children[v] {
			if next == -1 {
				continue
			}
			trie.bfsOrder[tail] = next
			tail++
			f := trie.suffixLink[v]
			for f != -1 && trie.children[f][i] == -1 {
				f = trie.suffixLink[f]
			}
			trie.suffixLink[next] = f
			if f == -1 {
				trie.suffixLink[next] = 0
			} else {
				trie.suffixLink[next] = trie.children[f][i]
			}
		}
	}
	if !needUpdateChildren {
		return
	}
	for _, v := range trie.bfsOrder {
		for i, next := range trie.children[v] {
			if next == -1 {
				f := trie.suffixLink[v]
				if f == -1 {
					trie.children[v][i] = 0
				} else {
					trie.children[v][i] = trie.children[f][i]
				}
			}
		}
	}
}

// 获取每个状态匹配到的模式串的个数.
func (trie *ACAutoMatonArray) GetCounter() []int {
	counter := make([]int, len(trie.children))
	for _, pos := range trie.Words {
		counter[pos]++
	}
	for _, v := range trie.bfsOrder {
		if v != 0 {
			counter[v] += counter[trie.suffixLink[v]]
		}
	}
	return counter
}

// 获取每个状态匹配到的模式串的索引.
func (trie *ACAutoMatonArray) GetIndexes() [][]int {
	res := make([][]int, len(trie.children))
	for i, pos := range trie.Words {
		res[pos] = append(res[pos], i)
	}
	for _, v := range trie.bfsOrder {
		if v != 0 {
			from, to := trie.suffixLink[v], v
			arr1, arr2 := res[from], res[to]
			arr3 := make([]int, 0, len(arr1)+len(arr2))
			i, j := 0, 0
			for i < len(arr1) && j < len(arr2) {
				for i < len(arr1) && j < len(arr2) {
					if arr1[i] < arr2[j] {
						arr3 = append(arr3, arr1[i])
						i++
					} else if arr1[i] > arr2[j] {
						arr3 = append(arr3, arr2[j])
						j++
					} else {
						arr3 = append(arr3, arr1[i])
						i++
						j++
					}
				}
			}
			for i < len(arr1) {
				arr3 = append(arr3, arr1[i])
				i++
			}
			for j < len(arr2) {
				arr3 = append(arr3, arr2[j])
				j++
			}
			res[to] = arr3
		}
	}
	return res
}

// 按照拓扑序进行转移.
func (trie *ACAutoMatonArray) Dp(f func(from, to int)) {
	for _, v := range trie.bfsOrder {
		if v != 0 {
			f(trie.suffixLink[v], v)
		}
	}
}

func (trie *ACAutoMatonArray) newNode() int {
	trie.Parent = append(trie.Parent, -1)
	nexts := make([]int, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.children = append(trie.children, nexts)
	return len(trie.children) - 1
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
