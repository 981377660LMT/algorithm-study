// https://www.luogu.com.cn/blog/yszs/ac-zi-dong-ji-fou-guo-shi-jian-fail-shu-di-gong-ju-pi-liao
// - !fail[i]表示在trie树上的第i个点表示的前缀，它在trie树上的最长后缀是第fail[i]个点表示的前缀。
// - 子串 = 前缀的后缀
//   Trie树（AC自动机）的祖先节点 = 前缀
//   Fail树的祖先节点 = 后缀
//   字符串x在字符串y中出现的次数 = `Fail树中x的子树`与`Trie树中y到根的路径`的交点个数
//
// 1.dp类型题: 一般都是dfs(index,pos):长度为index的字符串，当前trie状态为pos.
//	枚举26种字符(字符集)转移.
//
// 2.自定义结点信息题：用WordPos初始化状态，用SuffixLink转移.
// 3.Trie图：用 Move 来获取每个点的临边.
// 4.失配树：用SuffixLink来获取每个点的父亲.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// SeparateString()
	// P5357()
	// P3966()
	// P3121()
	// P4052()
	P3041()
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

	minBorder := make([]int, acm.Size()) // 每个状态(前缀)的最短border
	for i := range minBorder {
		minBorder[i] = INF
	}
	for i, pos := range acm.WordPos {
		minBorder[pos] = min(minBorder[pos], len(forbidden[i]))
	}
	acm.Dp(func(from, to int) { minBorder[to] = min(minBorder[to], minBorder[from]) })

	res, left, pos := 0, 0, 0
	for right, char := range word {
		pos = acm.Move(pos, int(char))
		left = max(left, right-minBorder[pos]+2)
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

// P5357 【模板】AC 自动机（二次加强版）
// https://www.luogu.com.cn/problem/P5357
func P5357() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	acm := NewACAutoMatonArray(26, 'a')
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		acm.AddString(s)
	}
	acm.BuildSuffixLink(false)

	var s string
	fmt.Fscan(in, &s)

	indexes := acm.GetIndexes()
	res := make([]int, n)
	pos := 0
	for i := 0; i < len(s); i++ {
		pos = acm.Move(pos, int(s[i]))
		for _, j := range indexes[pos] {
			res[j]++
		}
	}
	for i := 0; i < n; i++ {
		fmt.Fprintln(out, res[i])
	}
}

// https://www.luogu.com.cn/problem/P3966
// 一篇论文是由许多单词组成但小张发现一个单词会在论文中出现很多次。
// 他想知道每个单词分别在论文中出现了多少次。
func P3966() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	acm := NewACAutoMatonArray(26, 97)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		words[i] = s
		acm.AddString(s)
	}
	acm.BuildSuffixLink(false)

	res := make([]int, n)
	indexes := acm.GetIndexes()
	pos := 0
	for _, w := range words {
		for i := range w {
			pos = acm.Move(pos, int(w[i]))
			for _, v := range indexes[pos] {
				res[v]++
			}
		}
		pos = 0 // !每个单词之间是独立的，所以要重置pos
	}

	for _, v := range res {
		fmt.Println(v)
	}
}

// https://www.luogu.com.cn/problem/P3121
// 在longer中不断删除toRemove中的字符(按顺序遍历toRemove查找到要删除的)，求剩下的字符串.
// 注意删除一个单词后可能会导致 s 中出现另一个toRemove中的单词.
// !保证toRemove中没有包含关系(子串关系)，即删除时保证只会删除一个单词.
func P3121() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var longer string
	fmt.Fscan(in, &longer)

	acm := NewACAutoMatonArray(26, 97)

	var n int
	fmt.Fscan(in, &n)
	toRemove := make([]string, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		toRemove[i] = s
		acm.AddString(s)
	}

	acm.BuildSuffixLink(true)

	posLen := make([]int, acm.Size()) // 每个状态对应的字符长度
	for i, p := range acm.WordPos {
		posLen[p] = len(toRemove[i])
	}

	pos := 0
	stack := make([]int, 0, len(longer))
	posRecord := make([]int, len(longer))
	for i := range longer {
		v := int(longer[i])
		pos = acm.Move(pos, v)
		posRecord[i] = pos
		stack = append(stack, i)
		if wordLen := posLen[pos]; wordLen > 0 { // 由于模式串不存在包含关系，所以不会出现多个模式串同时匹配的情况
			stack = stack[:len(stack)-wordLen]
			if len(stack) > 0 {
				pos = posRecord[stack[len(stack)-1]]
			} else {
				pos = 0
			}
		}
	}

	res := make([]byte, 0, len(stack))
	for _, v := range stack {
		res = append(res, longer[v])
	}
	fmt.Fprintln(out, string(res))
}

// P4052 [JSOI2007] 文本生成器
// https://www.luogu.com.cn/problem/P4052
// !给定一些模式串，求长度为m的所有文本串的个数，且该文本串至少包括一个模式串，答案对10007取模
// n<=60,targetLen<=100,len(word[i])<=100.
// 合法的情况总数并不好求，考虑求出不合法的情况（即不存在一个子串等于模式串），最后用总数26^m减去即可。
func P4052() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e4 + 7

	acm := NewACAutoMatonArray(26, 'A')

	var n, targetLen int
	fmt.Fscan(in, &n, &targetLen)
	words := make([]string, n)
	lengthSum := 0
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		words[i] = s
		acm.AddString(s)
		lengthSum += len(s)
	}
	acm.BuildSuffixLink(true)

	size := acm.Size()
	counter := acm.GetCounter()

	memo := make([][]int, targetLen)
	for i := range memo {
		row := make([]int, size)
		for j := range row {
			row[j] = -1
			memo[i] = row
		}
	}
	var dfs func(index int, pos int) int
	dfs = func(index, pos int) int {
		if index == targetLen {
			return 1
		}
		if tmp := memo[index][pos]; tmp != -1 {
			return tmp
		}
		res := 0
		for v := 65; v < 65+26; v++ {
			nextPos := acm.Move(pos, v)
			if counter[nextPos] > 0 {
				continue
			}
			res += dfs(index+1, nextPos)
			res %= MOD
		}

		memo[index][pos] = res
		return res
	}
	bad := dfs(0, 0)

	qpow := func(a, b int) int {
		res := 1
		for b > 0 {
			if b&1 == 1 {
				res = res * a % MOD
			}
			a = a * a % MOD
			b >>= 1
		}
		return res
	}

	res := (qpow(26, targetLen) - bad) % MOD
	if res < 0 {
		res += MOD
	}

	fmt.Println(res)
}

// P3041 [USACO12JAN] Video Game G
// https://www.luogu.com.cn/problem/P3041
// 给出n个字典单词，问长度为k的字符串最多可以包含多少个字典单词。
// n<=20,k<=1000,len(word[i])<=15
func P3041() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	words := make([]string, n)
	acm := NewACAutoMatonArray(26, 'A')
	for i := range words {
		var s string
		fmt.Fscan(in, &s)
		words[i] = s
		acm.AddString(s)
	}
	acm.BuildSuffixLink(true)
	size := acm.Size()
	counter := acm.GetCounter()
	memo := make([][]int, k)
	for i := range memo {
		row := make([]int, size)
		for j := range row {
			row[j] = -1
			memo[i] = row
		}
	}
	var dfs func(index int, pos int) int
	dfs = func(index, pos int) int {
		if index == k {
			return 0
		}
		if tmp := memo[index][pos]; tmp != -1 {
			return tmp
		}
		res := 0
		for v := 65; v < 65+26; v++ {
			nextPos := acm.Move(pos, v)
			res = max(res, dfs(index+1, nextPos)+counter[nextPos])
		}
		memo[index][pos] = res
		return res
	}

	res := dfs(0, 0)
	fmt.Println(res)
}

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个前缀.
type ACAutoMatonArray struct {
	WordPos            []int     // WordPos[i] 表示加入的第i个模式串对应的节点编号.
	Parent             []int     // parent[v] 表示节点v的父节点.
	sigma              int32     // 字符集大小.
	offset             int32     // 字符集的偏移量.
	children           [][]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	suffixLink         []int32   // 又叫fail.指向当前trie节点(对应一个前缀)的最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
	bfsOrder           []int32   // 结点的拓扑序,0表示虚拟节点.
	needUpdateChildren bool      // 是否需要更新children数组.
}

func NewACAutoMatonArray(sigma, offset int) *ACAutoMatonArray {
	res := &ACAutoMatonArray{sigma: int32(sigma), offset: int32(offset)}
	res.newNode()
	return res
}

// 添加一个字符串，返回最后一个字符对应的节点编号.
func (trie *ACAutoMatonArray) AddString(str string) int {
	if len(str) == 0 {
		return 0
	}
	pos := 0
	for _, s := range str {
		ord := int32(s) - trie.offset
		if trie.children[pos][ord] == -1 {
			trie.children[pos][ord] = trie.newNode()
			trie.Parent[len(trie.Parent)-1] = pos
		}
		pos = int(trie.children[pos][ord])
	}
	trie.WordPos = append(trie.WordPos, pos)
	return pos
}

// 功能与 AddString 相同.
func (trie *ACAutoMatonArray) AddFrom(n int, getOrd func(i int) int) int {
	if n == 0 {
		return 0
	}
	pos := 0
	for i := 0; i < n; i++ {
		s := getOrd(i)
		ord := int32(s) - trie.offset
		if trie.children[pos][ord] == -1 {
			trie.children[pos][ord] = trie.newNode()
			trie.Parent[len(trie.Parent)-1] = pos
		}
		pos = int(trie.children[pos][ord])
	}
	trie.WordPos = append(trie.WordPos, pos)
	return pos
}

// 在pos位置添加一个字符，返回新的节点编号.
func (trie *ACAutoMatonArray) AddChar(pos int, ord int) int {
	ord -= int(trie.offset)
	if trie.children[pos][ord] != -1 {
		return int(trie.children[pos][ord])
	}
	trie.children[pos][ord] = trie.newNode()
	trie.Parent[len(trie.Parent)-1] = pos
	return int(trie.children[pos][ord])
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *ACAutoMatonArray) Move(pos int, ord int) int {
	ord -= int(trie.offset)
	if trie.needUpdateChildren {
		return int(trie.children[pos][ord])
	}
	for {
		nexts := trie.children[pos]
		if nexts[ord] != -1 {
			return int(nexts[ord])
		}
		if pos == 0 {
			return 0
		}
		pos = int(trie.suffixLink[pos])
	}
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *ACAutoMatonArray) Size() int {
	return len(trie.children)
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *ACAutoMatonArray) BuildSuffixLink(needUpdateChildren bool) {
	trie.needUpdateChildren = needUpdateChildren
	trie.suffixLink = make([]int32, len(trie.children))
	for i := range trie.suffixLink {
		trie.suffixLink[i] = -1
	}
	trie.bfsOrder = make([]int32, len(trie.children))
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

// 获取每个状态包含的模式串的个数.
func (trie *ACAutoMatonArray) GetCounter() []int {
	counter := make([]int, len(trie.children))
	for _, pos := range trie.WordPos {
		counter[pos]++
	}
	for _, v := range trie.bfsOrder {
		if v != 0 {
			counter[v] += counter[trie.suffixLink[v]]
		}
	}
	return counter
}

// 获取每个状态包含的模式串的索引.
func (trie *ACAutoMatonArray) GetIndexes() [][]int {
	res := make([][]int, len(trie.children))
	for i, pos := range trie.WordPos {
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

// 按照拓扑序进行转移(EnumerateFail).
func (trie *ACAutoMatonArray) Dp(f func(from, to int)) {
	for _, v := range trie.bfsOrder {
		if v != 0 {
			f(int(trie.suffixLink[v]), int(v))
		}
	}
}

func (trie *ACAutoMatonArray) BuildFailTree() [][]int {
	res := make([][]int, trie.Size())
	trie.Dp(func(pre, cur int) {
		res[pre] = append(res[pre], cur)
	})
	return res
}

func (trie *ACAutoMatonArray) BuildTrieTree() [][]int {
	res := make([][]int, trie.Size())
	for i := 1; i < trie.Size(); i++ {
		res[trie.Parent[i]] = append(res[trie.Parent[i]], i)
	}
	return res
}

func (trie *ACAutoMatonArray) newNode() int32 {
	trie.Parent = append(trie.Parent, -1)
	nexts := make([]int32, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.children = append(trie.children, nexts)
	return int32(len(trie.children) - 1)
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
