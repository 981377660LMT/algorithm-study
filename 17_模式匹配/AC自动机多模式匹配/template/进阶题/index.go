package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {

	// CF1202E()
	// P2292()
	// P2336()
	CF163E()
}

// P2444 [POI2000] 病毒
// 给定一些01模式串,判断是否存在无限长01串不包含任何一个模式串.
// https://www.luogu.com.cn/problem/P2444
// sum(words[i].length) <= 3e4
// !“无限长”就是指删去不能经过的节点以后形成的有向图存在环
// 只要在AC自动机上，存在一个环，使得这个环上和从根到这个环的路径上，所有的点都不是危险节点，就有解
func P2444() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	acm := NewACAutoMatonArray(2, 48)
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(true)

	counter := acm.GetCounter()

	// 拓扑排序判环.
	adjList := make([][]int, acm.Size())
	deg := make([]int, acm.Size())
	for i := range adjList {
		if counter[i] == 0 {
			for j := 0; j < 2; j++ {
				next := acm.Move(i, j+48)
				if counter[next] == 0 {
					adjList[i] = append(adjList[i], next)
					deg[next]++
				}
			}
		}
	}

	queue := []int{}
	for i, v := range deg {
		if v == 0 {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, next := range adjList[v] {
			deg[next]--
			if deg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	for _, v := range deg {
		if v != 0 {
			fmt.Fprintln(out, "TAK") // 有环
			return
		}
	}
	fmt.Fprintln(out, "NIE")
}

// You Are Given Some Strings...
// https://www.luogu.com.cn/problem/CF1202E
// 给定一个目标串，和一些模式串.
// 求所有的模式串对(words[i], words[j])拼接后，在目标串中出现的次数之和.
// !等价于：对文本串的每个前缀求出有多少模式串是他的后缀（后面那个只要把串反转就好）
func CF1202E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var longer string
	fmt.Fscan(in, &longer)

	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	reverse := func(s string) string {
		res := make([]byte, len(s))
		for i := range s {
			res[len(s)-1-i] = s[i]
		}
		return string(res)
	}

	acm := NewACAutoMatonArray(26, 97)
	rAcm := NewACAutoMatonArray(26, 97)
	for _, word := range words {
		acm.AddString(word)
		rAcm.AddString(reverse(word))
	}
	acm.BuildSuffixLink(true)
	rAcm.BuildSuffixLink(true)

	counter := acm.GetCounter()
	rCounter := rAcm.GetCounter()

	pre := make([]int, len(longer))
	pos := 0
	for i := 0; i < len(longer); i++ {
		pos = acm.Move(pos, int(longer[i]))
		pre[i] = counter[pos]
	}

	suf := make([]int, len(longer))
	pos = 0
	for i := len(longer) - 1; i >= 0; i-- {
		pos = rAcm.Move(pos, int(longer[i]))
		suf[i] = rCounter[pos]
	}

	res := 0
	for i := 0; i < len(longer)-1; i++ {
		res += pre[i] * suf[i+1]
	}

	fmt.Fprintln(out, res)
}

// P2292 [HNOI2004] L 语言
// https://www.luogu.com.cn/problem/P2292
// 给定n个模式串，和q个文本串
// 对每个文本串，求出其最长的前缀，满足该前缀是由若干模式串（可以多次使用）首尾拼接而成的。
// !n<=20 len(word[i])<=20
// q<=50 len(text[i])<=1e6
// 状压+AC自动机
// dp[i]表示前i个字符是否能被理解.
func P2292() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}
	texts := make([]string, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &texts[i])
	}

	acm := NewACAutoMatonArray(26, 97)
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(true)

	lengthMask := make([]uint, acm.Size()) // !lengthMask[i] 表示第i个节点对应的模式串的长度的集合.
	for i, p := range acm.WordPos {
		lengthMask[p] |= 1 << len(words[i])
	}
	acm.Dp(func(from, to int) {
		lengthMask[to] |= lengthMask[from]
	})

	for _, text := range texts {
		pos, res := 0, -1
		dp := uint(1) // 保存前64个位置的dp值,自然溢出
		for i, v := range text {
			pos = acm.Move(pos, int(v))
			dp <<= 1
			if dp&lengthMask[pos] != 0 {
				dp |= 1
				res = i
			}
		}
		fmt.Fprintln(out, res+1)
	}
}

// https://www.luogu.com.cn/problem/P2336
// P2336 [SCOI2012] 喵星球上的点名
// 喵星球上的老师会选择 m 个串来点名，每次读出一个串的时候，如果这个串是一个喵星人的姓或名的子串，那么这个喵星人就必须答到。
// !把每文本串的姓和名中间加一个字符变成一个串，于是这个题就变成了：
// !给一些文本串和一些模式串，求每个模式串在多少个文本串中出现过，以及每个文本串里有多少个模式串。
// bitset加速
func P2336() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	texts := make([][]int, n)
	for i := 0; i < n; i++ {
		var len1 int
		fmt.Fscan(in, &len1)
		nums1 := make([]int, len1)
		for j := 0; j < len1; j++ {
			fmt.Fscan(in, &nums1[j])
		}
		var len2 int
		fmt.Fscan(in, &len2)
		nums2 := make([]int, len2)
		for j := 0; j < len2; j++ {
			fmt.Fscan(in, &nums2[j])
		}
		texts[i] = append(texts[i], nums1...)
		texts[i] = append(texts[i], -1e9) // -1表示分隔符
		texts[i] = append(texts[i], nums2...)
	}

	patterns := make([][]int, m)
	for i := 0; i < m; i++ {
		var len int
		fmt.Fscan(in, &len)
		patterns[i] = make([]int, len)
		for j := 0; j < len; j++ {
			fmt.Fscan(in, &patterns[i][j])
		}
	}

	acm := NewACAutoMatonMap()
	for _, p := range patterns {
		acm.AddString(p)
	}
	acm.BuildSuffixLink()

	indexes := acm.GetIndexes()
	freq1 := make([]int, m)
	freq2 := make([]int, n)

	bs := NewBS(m)
	for i, text := range texts {
		pos := 0
		for _, v := range text {
			pos = acm.Move(pos, v)
			for _, p := range indexes[pos] {
				bs.Set(p)
			}
		}
		bs.ForEach(func(p int) {
			freq1[p]++
			freq2[i]++
		})
		bs.Clear()
	}

	for _, v := range freq1 {
		fmt.Fprintln(out, v)
	}
	for _, v := range freq2 {
		fmt.Fprint(out, v, " ")
	}
}

// e-Government
// https://www.luogu.com.cn/problem/CF163E
// https://codeforces.com/contest/163/submission/215316998
// 给定包含k个字符串的容器S。
// 开始时，所有字符串都被启用。
// 有q个操作，操作有三种类型：
// 以‘？’开头的操作为询问操作，询问当前容器S中的每一个字符串匹配询问字符串的次数之和；
// 以‘+’开头的操作为set操作，表示将编号为i的字符串启用；
// 以‘-’开头的操作为reset操作，表示将编号为i的字符串禁用。
// 注意当编号为i的字符串已经在容器中时，允许存在添加编号为i的字符串，删除亦然。
// !不删除字符时，答案等于fail树中某个结点到根节点路径上点权之和;
// !删除字符后，相当于将子树中所有结点权值减去1.
// 树状数组维护.
// 树状数组+AC自动机
func CF163E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q, k int
	fmt.Fscan(in, &q, &k)
	words := make([]string, k)
	acm := NewACAutoMatonArray(26, 97)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &words[i])
		acm.AddString(words[i])
	}
	acm.BuildSuffixLink(true)

	n := acm.Size()
	failTree := acm.BuildFailTree()

	// dfs序
	lid, rid := make([]int, n), make([]int, n)
	dfn := 0
	var dfs func(cur, pre int)
	dfs = func(cur, pre int) {
		lid[cur] = dfn
		dfn++
		for _, next := range failTree[cur] {
			if next != pre {
				dfs(next, cur)
			}
		}
		rid[cur] = dfn
	}
	dfs(0, -1)

	ok := make([]bool, n)
	bit := NewBITRangeAddRangeSumArray(n)
	// bit := NewBitArray(n)

	// TODO: bit问题
	add := func(index int) {
		if ok[index] {
			return
		}
		ok[index] = true
		pos := acm.WordPos[index]
		bit.AddRange(lid[pos], rid[pos], 1)
		// bit.Add(lid[pos], 1)
		// bit.Add(rid[pos], -1)
	}
	remove := func(index int) {
		if !ok[index] {
			return
		}
		ok[index] = false
		pos := acm.WordPos[index]
		bit.AddRange(lid[pos], rid[pos], -1)
		// bit.Add(lid[pos], -1)
		// bit.Add(rid[pos], 1)

	}
	// 到根节点的路径上的点权和.
	getCount := func(pos int) int {
		return bit.QueryRange(0, lid[pos]+1)
		// return bit.QueryPrefix(lid[pos] + 1)
	}
	for i := 0; i < k; i++ {
		add(i)
	}

	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op[0] == '+' {
			id, _ := strconv.Atoi(op[1:])
			id--
			add(id)
		} else if op[0] == '-' {
			id, _ := strconv.Atoi(op[1:])
			id--
			remove(id)
		} else {
			text := op[1:]
			res, pos := 0, 0
			for _, v := range text {
				pos = acm.Move(pos, int(v))
				res += getCount(pos)
			}
			fmt.Println(res)
		}
	}
}

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个字符串.
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
			trie.Parent = append(trie.Parent, pos)
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
	trie.Parent = append(trie.Parent, pos)
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

func (trie *ACAutoMatonArray) newNode() int32 {
	trie.Parent = append(trie.Parent, -1)
	nexts := make([]int32, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.children = append(trie.children, nexts)
	return int32(len(trie.children) - 1)
}

type T = int

type ACAutoMatonMap struct {
	WordPos    []int         // WordPos[i] 表示加入的第i个模式串对应的节点编号.
	children   []map[T]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	suffixLink []int32       // 又叫fail.指向当前节点最长真后缀对应结点.
	bfsOrder   []int32       // 结点的拓扑序,0表示虚拟节点.
}

func NewACAutoMatonMap() *ACAutoMatonMap {
	return &ACAutoMatonMap{
		WordPos:  []int{},
		children: []map[T]int32{{}},
	}
}

func (ac *ACAutoMatonMap) AddString(s []T) int {
	if len(s) == 0 {
		return 0
	}
	pos := 0
	for i := 0; i < len(s); i++ {
		ord := s[i]
		nexts := ac.children[pos]
		if next, ok := nexts[ord]; ok {
			pos = int(next)
		} else {
			nextState := len(ac.children)
			nexts[ord] = int32(nextState)
			pos = nextState
			ac.children = append(ac.children, map[T]int32{})
		}
	}
	ac.WordPos = append(ac.WordPos, pos)
	return pos
}

func (ac *ACAutoMatonMap) AddChar(pos int, ord T) int {
	nexts := ac.children[pos]
	if next, ok := nexts[ord]; ok {
		return int(next)
	}
	nextState := len(ac.children)
	nexts[ord] = int32(nextState)
	ac.children = append(ac.children, map[T]int32{})
	return nextState
}

func (ac *ACAutoMatonMap) Move(pos int, ord T) int {
	pos32 := int32(pos)
	for {
		nexts := ac.children[pos32]
		if next, ok := nexts[ord]; ok {
			return int(next)
		}
		if pos32 == 0 {
			return 0
		}
		pos32 = (ac.suffixLink[pos32])
	}
}

func (ac *ACAutoMatonMap) BuildSuffixLink() {
	ac.suffixLink = make([]int32, len(ac.children))
	for i := range ac.suffixLink {
		ac.suffixLink[i] = -1
	}
	ac.bfsOrder = make([]int32, len(ac.children))
	head, tail := 0, 1
	for head < tail {
		v := ac.bfsOrder[head]
		head++
		for char, next := range ac.children[v] {
			ac.bfsOrder[tail] = next
			tail++
			f := ac.suffixLink[v]
			for f != -1 {
				if _, ok := ac.children[f][char]; ok {
					break
				}
				f = ac.suffixLink[f]
			}
			if f == -1 {
				ac.suffixLink[next] = 0
			} else {
				ac.suffixLink[next] = ac.children[f][char]
			}
		}
	}
}

func (ac *ACAutoMatonMap) GetCounter() []int {
	counter := make([]int, len(ac.children))
	for _, pos := range ac.WordPos {
		counter[pos]++
	}
	for _, v := range ac.bfsOrder {
		if v != 0 {
			counter[v] += counter[ac.suffixLink[v]]
		}
	}
	return counter
}

func (ac *ACAutoMatonMap) GetIndexes() [][]int {
	res := make([][]int, len(ac.children))
	for i, pos := range ac.WordPos {
		res[pos] = append(res[pos], i)
	}
	for _, v := range ac.bfsOrder {
		if v != 0 {
			from, to := ac.suffixLink[v], v
			arr1, arr2 := res[from], res[to]
			arr3 := make([]int, 0, len(arr1)+len(arr2))
			i, j := 0, 0
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

func (ac *ACAutoMatonMap) Dp(f func(from, to int)) {
	for _, v := range ac.bfsOrder {
		if v != 0 {
			f(int(ac.suffixLink[v]), int(v))
		}
	}
}

func (ac *ACAutoMatonMap) BuildFailTree() [][]int {
	res := make([][]int, ac.Size())
	ac.Dp(func(pre, cur int) {
		res[pre] = append(res[pre], cur)
	})
	return res
}

func (ac *ACAutoMatonMap) Size() int {
	return len(ac.children)
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

type BS []uint

func NewBS(n int) BS { return make(BS, n>>6+1) } // (n+64-1)>>6

func (b BS) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b BS) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b BS) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b BS) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0

func (b BS) Copy() BS {
	res := make(BS, len(b))
	copy(res, b)
	return res
}

func (bs BS) Clear() {
	for i := range bs {
		bs[i] = 0
	}
}

func (b BS) ForEach(f func(p int)) {
	for i, v := range b {
		for ; v != 0; v &= v - 1 {
			j := i<<6 | bits.TrailingZeros(v)
			f(j)
		}
	}
}

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int, f func(i int) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITArray) QueryAll() int {
	return b.total
}

func (b *BITArray) MaxRight(check func(index, preSum int) bool) int {
	i := 0
	s := 0
	k := 1
	for 2*k <= b.n {
		k *= 2
	}
	for k > 0 {
		if i+k-1 < b.n {
			t := s + b.data[i+k-1]
			if check(i+k, t) {
				i += k
				s = t
			}
		}
		k >>= 1
	}
	return i
}

// 0/1 树状数组查找第 k(0-based) 个1的位置.
// UpperBound.
func (b *BITArray) Kth(k int) int {
	return b.MaxRight(func(index, preSum int) bool { return preSum <= k })
}

func (b *BITArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}

// !Range Add Range Sum, 0-based.
type BITRangeAddRangeSumArray struct {
	n    int
	bit0 *BITArray
	bit1 *BITArray
}

func NewBITRangeAddRangeSumArray(n int) *BITRangeAddRangeSumArray {
	return &BITRangeAddRangeSumArray{
		n:    n,
		bit0: NewBitArray(n),
		bit1: NewBitArray(n),
	}
}

func NewBITRangeAddRangeSumFrom(n int, f func(index int) int) *BITRangeAddRangeSumArray {
	return &BITRangeAddRangeSumArray{
		n:    n,
		bit0: NewBitArrayFrom(n, f),
		bit1: NewBitArray(n),
	}
}

func (b *BITRangeAddRangeSumArray) Add(index int, delta int) {
	b.bit0.Add(index, delta)
}

func (b *BITRangeAddRangeSumArray) AddRange(start, end int, delta int) {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return
	}
	b.bit0.Add(start, -delta*start)
	b.bit0.Add(end, delta*end)
	b.bit1.Add(start, delta)
	b.bit1.Add(end, -delta)
}

func (b *BITRangeAddRangeSumArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	rightRes := b.bit1.QueryPrefix(end)*end + b.bit0.QueryPrefix(end)
	leftRes := b.bit1.QueryPrefix(start)*start + b.bit0.QueryPrefix(start)
	return rightRes - leftRes
}

func (b *BITRangeAddRangeSumArray) String() string {
	res := []string{}
	for i := 0; i < b.n; i++ {
		res = append(res, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BITRangeAddRangeSumArray: [%v]", strings.Join(res, ", "))
}

type Tree struct {
	Tree                 [][][2]int // (next, weight)
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int // 欧拉序[in,out)
	idToNode             []int
	top, heavySon        []int
	timer                int
}

func NewTree(n int) *Tree {
	tree := make([][][2]int, n)
	lid := make([]int, n)
	rid := make([]int, n)
	idToNode := make([]int, n)
	top := make([]int, n)   // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n) // 深度
	depthWeighted := make([]int, n)
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	return &Tree{
		Tree:          tree,
		Depth:         depth,
		DepthWeighted: depthWeighted,
		Parent:        parent,
		LID:           lid,
		RID:           rid,
		idToNode:      idToNode,
		top:           top,
		heavySon:      heavySon,
	}
}

// 添加无向边 u-v, 边权为w.
func (tree *Tree) AddEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Tree[v] = append(tree.Tree[v], [2]int{u, w})
}

// 添加有向边 u->v, 边权为w.
func (tree *Tree) AddDirectedEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
}

// root:0-based
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *Tree) Build(root int) {
	if root != -1 {
		tree.build(root, -1, 0, 0)
		tree.markTop(root, root)
	} else {
		for i := 0; i < len(tree.Tree); i++ {
			if tree.Parent[i] == -1 {
				tree.build(i, -1, 0, 0)
				tree.markTop(i, i)
			}
		}
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree) Id(root int) (int, int) {
	return tree.LID[root], tree.RID[root]
}

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *Tree) Eid(u, v int) int {
	if tree.LID[u] > tree.LID[v] {
		return tree.LID[u]
	}
	return tree.LID[v]
}

func (tree *Tree) LCA(u, v int) int {
	for {
		if tree.LID[u] > tree.LID[v] {
			u, v = v, u
		}
		if tree.top[u] == tree.top[v] {
			return u
		}
		v = tree.Parent[tree.top[v]]
	}
}

func (tree *Tree) RootedLCA(u, v int, root int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *Tree) RootedParent(u int, root int) int {
	return tree.Jump(u, root, 1)
}

func (tree *Tree) Dist(u, v int, weighted bool) int {
	if weighted {
		return tree.DepthWeighted[u] + tree.DepthWeighted[v] - 2*tree.DepthWeighted[tree.LCA(u, v)]
	}
	return tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)]
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
//	kthAncestor(root,0) == root
func (tree *Tree) KthAncestor(root, k int) int {
	if k > tree.Depth[root] {
		return -1
	}
	for {
		u := tree.top[root]
		if tree.LID[root]-k >= tree.LID[u] {
			return tree.idToNode[tree.LID[root]-k]
		}
		k -= tree.LID[root] - tree.LID[u] + 1
		root = tree.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (tree *Tree) Jump(from, to, step int) int {
	if step == 1 {
		if from == to {
			return -1
		}
		if tree.IsInSubtree(to, from) {
			return tree.KthAncestor(to, tree.Depth[to]-tree.Depth[from]-1)
		}
		return tree.Parent[from]
	}
	c := tree.LCA(from, to)
	dac := tree.Depth[from] - tree.Depth[c]
	dbc := tree.Depth[to] - tree.Depth[c]
	if step > dac+dbc {
		return -1
	}
	if step <= dac {
		return tree.KthAncestor(from, step)
	}
	return tree.KthAncestor(to, dac+dbc-step)
}

func (tree *Tree) CollectChild(root int) []int {
	res := []int{}
	for _, e := range tree.Tree[root] {
		next := e[0]
		if next != tree.Parent[root] {
			res = append(res, next)
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree) GetPathDecomposition(u, v int, vertex bool) [][2]int {
	up, down := [][2]int{}, [][2]int{}
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			down = append(down, [2]int{tree.LID[tree.top[v]], tree.LID[v]})
			v = tree.Parent[tree.top[v]]
		} else {
			up = append(up, [2]int{tree.LID[u], tree.LID[tree.top[u]]})
			u = tree.Parent[tree.top[u]]
		}
	}
	edgeInt := 1
	if vertex {
		edgeInt = 0
	}
	if tree.LID[u] < tree.LID[v] {
		down = append(down, [2]int{tree.LID[u] + edgeInt, tree.LID[v]})
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		up = append(up, [2]int{tree.LID[u], tree.LID[v] + edgeInt})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree) EnumeratePathDecomposition(u, v int, vertex bool, f func(start, end int)) {
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			a, b := tree.LID[tree.top[v]], tree.LID[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = tree.Parent[tree.top[v]]
		} else {
			a, b := tree.LID[u], tree.LID[tree.top[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = tree.Parent[tree.top[u]]
		}
	}

	edgeInt := 1
	if vertex {
		edgeInt = 0
	}

	if tree.LID[u] < tree.LID[v] {
		a, b := tree.LID[u]+edgeInt, tree.LID[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		a, b := tree.LID[u], tree.LID[v]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

func (tree *Tree) GetPath(u, v int) []int {
	res := []int{}
	composition := tree.GetPathDecomposition(u, v, true)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, tree.idToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.idToNode[i])
			}
		}
	}
	return res
}

// 以root为根时,结点v的子树大小.
func (tree *Tree) SubSize(v, root int) int {
	if root == -1 {
		return tree.RID[v] - tree.LID[v]
	}
	if v == root {
		return len(tree.Tree)
	}
	x := tree.Jump(v, root, 1)
	if tree.IsInSubtree(v, x) {
		return tree.RID[v] - tree.LID[v]
	}
	return len(tree.Tree) - tree.RID[x] + tree.LID[x]
}

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *Tree) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链底端节点.
func (tree *Tree) GetHeavyPath(start int) []int {
	heavyPath := []int{start}
	cur := start
	for tree.heavySon[cur] != -1 {
		cur = tree.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

// 结点v的重儿子.如果没有重儿子,返回-1.
func (tree *Tree) GetHeavyChild(v int) int {
	k := tree.LID[v] + 1
	if k == len(tree.Tree) {
		return -1
	}
	w := tree.idToNode[k]
	if tree.Parent[w] == v {
		return w
	}
	return -1
}

func (tree *Tree) ELID(u int) int {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *Tree) ERID(u int) int {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *Tree) build(cur, pre, dep, dist int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, e := range tree.Tree[cur] {
		next, weight := e[0], e[1]
		if next != pre {
			nextSize := tree.build(next, cur, dep+1, dist+weight)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	tree.Depth[cur] = dep
	tree.DepthWeighted[cur] = dist
	tree.heavySon[cur] = heavySon
	tree.Parent[cur] = pre
	return subSize
}

func (tree *Tree) markTop(cur, top int) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.idToNode[tree.timer] = cur
	tree.timer++
	heavySon := tree.heavySon[cur]
	if heavySon != -1 {
		tree.markTop(heavySon, top)
		for _, e := range tree.Tree[cur] {
			next := e[0]
			if next != heavySon && next != tree.Parent[cur] {
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
}
