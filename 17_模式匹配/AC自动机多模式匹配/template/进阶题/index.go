package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {

	// CF1202E()
	// P2292()
	// P2336()
	// CF163E()
	// CF1437G()
	// SP9941()
	// P2414()
	// CF1207G()
	// CF86C()
	CF547E()

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
	bit := NewBITRangeAddPointGet(n)

	add := func(index int) {
		if ok[index] {
			return
		}
		ok[index] = true
		pos := acm.WordPos[index]
		bit.AddRange(lid[pos], rid[pos], 1)
	}
	remove := func(index int) {
		if !ok[index] {
			return
		}
		ok[index] = false
		pos := acm.WordPos[index]
		bit.AddRange(lid[pos], rid[pos], -1)
	}
	getCount := func(pos int) int {
		return bit.Get(lid[pos])
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

// Death DBMS - 死亡笔记数据库管理系统
// https://www.luogu.com.cn/problem/CF1437G
// 给定m个字符串，每个字符串有一个权值.初始时，所有字符串的权值都为0.
// 给定q个操作，操作有两种类型：
// 1 i v 表示将第i个字符串的权值修改为v.
// 2 s 求所有是s的子串的字符串的权值最大值.
//
// 对所有 s 建出 AC 自动机。每次询问把 t 扔进去匹配，
// 每个节点在 fail 树上的祖先都是这个节点所代表字符串的子串，树剖取最大值即可。
// 注意重复的字符串,用一个可删除堆维护每个node处的最大值.
func CF1437G() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var m, q int
	fmt.Fscan(in, &m, &q)

	words := make([]string, m)
	acm := NewACAutoMatonArray(26, 97)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &words[i])
		acm.AddString(words[i])
	}
	acm.BuildSuffixLink(true)

	tree := NewTree(acm.Size())
	acm.Dp(func(from, to int) {
		tree.AddDirectedEdge(from, to, 1)
	})
	tree.Build(0)

	values := make([]int, m)
	heaps := make([]*ErasableHeap, acm.Size()) // 维护每个结点的最大值.
	for i := range heaps {
		heaps[i] = NewErasableHeap(func(a, b H) bool { return a > b }, []int{-INF})
	}
	for _, p := range acm.WordPos {
		heaps[p].Push(0)
	}

	seg := NewSegmentTree(acm.Size(), func(i int) E { return -INF })
	for i := 0; i < acm.Size(); i++ {
		seg.Set(tree.LID[i], heaps[i].Peek())
	}

	update := func(index, value int) {
		if values[index] == value {
			return
		}
		node := acm.WordPos[index]
		preValue := values[index]
		values[index] = value
		heaps[node].Erase(preValue)
		heaps[node].Push(value)
		lid := tree.LID[node]
		seg.Set(lid, heaps[node].Peek())
	}

	query := func(s string) int {
		res := -1
		pos := 0
		for _, v := range s {
			pos = acm.Move(pos, int(v))
			tree.EnumeratePathDecomposition(0, pos, true, func(start, end int) {
				cur := seg.Query(start, end)
				res = max(res, cur)
			})
		}
		return res
	}

	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 1 {
			var index, value int
			fmt.Fscan(in, &index, &value)
			index--
			update(index, value)
		} else {
			var s string
			fmt.Fscan(in, &s)
			fmt.Fprintln(out, query(s))
		}
	}
}

// GRE - GRE Words
// https://www.luogu.com.cn/problem/SP9941
// 给定一个由字符串构成的序列，不同位置的字符串有自己权值。
// 现在让你选出一个子序列，使得在这个子序列中，
// !前面的串是后面的串的子串。
// 请你求满足条件的子序列的权值的最大值。
// 一个子序列权值是所有元素权值的和。
// 线段树优化dp
// dp[i] = max(dp[j]) + values[i], j 是 fail 树上 i 到根的路径上的点.
func SP9941() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	solve := func() int {
		var n int
		fmt.Fscan(in, &n)
		words := make([]string, n)
		values := make([]int, n)
		acm := NewACAutoMatonArray(26, 97)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &words[i])
			fmt.Fscan(in, &values[i])
			acm.AddString(words[i])
		}
		acm.BuildSuffixLink(true)

		tree := NewTree(acm.Size())
		acm.Dp(func(from, to int) {
			tree.AddDirectedEdge(from, to, 1)
		})
		tree.Build(0)

		dp := NewSegmentTree(acm.Size(), func(i int) E { return 0 })

		res := 0
		for i := 0; i < n; i++ {
			preMax := 0

			// query
			pos := 0
			for _, v := range words[i] {
				pos = acm.Move(pos, int(v))
				tree.EnumeratePathDecomposition(0, pos, true, func(start, end int) {
					cur := dp.Query(start, end)
					preMax = max(preMax, cur)
				})
			}

			// update
			res = max(res, preMax+values[i])
			dp.Update(tree.LID[pos], preMax+values[i])
		}

		return res
	}

	for i := 0; i < T; i++ {
		res := solve()
		fmt.Fprintf(out, "Case #%d: %d\n", i+1, res)
	}
}

// P2414 [NOI2011] 阿狸的打字机
// https://www.luogu.com.cn/problem/P2414
// 打字机上只有 28 个按键，分别印有 26 个小写英文字母和 B、P 两个字母。
// 经阿狸研究发现，这个打字机是这样工作的：
// - 输入小写字母，打字机的一个凹槽中会加入这个字母(这个字母加在凹槽的最后)。
// - 按一下印有 B 的按键，打字机凹槽中最后一个字母会消失。
// - 按一下印有 P 的按键，打字机会在纸上打印出凹槽中现有的所有字母并换行，但凹槽中的字母不会消失。
// 例如，阿狸输入 aPaPBbP，纸上被打印的字符如下：
//
// a
// aa
// ab
//
// 我们把纸上打印出来的字符串从 1 开始顺序编号，一直到 n。
// 打字机有一个非常有趣的功能，在打字机中暗藏一个带数字的小键盘，
// 在小键盘上输入两个数 (x,y)（其中 1≤x,y≤n），打字机会显示第 x 个打印的字符串在第 y 个打印的字符串中出现了多少次。
// !即：给你一颗 Trie，每次询问两个节点 x和y，求 x 代表的字符串在 y 代表的字符串中出现了多少次。
// !也即：给出若干个字符串，每次询问一个串在另一个串的出现次数。
// !等价于:fail树中x的子树(对应一些更长的后缀)与trie树中y到根节点的路径(对应一些更短的前缀)的公共结点数.
// 离线查询，将所有询问保存到y上，在 Trie树 上 dfs+回溯 即可.
// ACAM+树状数组
func P2414() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var command string
	fmt.Fscan(in, &command)
	acm := NewACAutoMatonArray(26, 97)
	wordToPos := []int{}
	pos := 0
	for _, v := range command {
		if v == 'B' {
			if p := acm.Parent[pos]; p != -1 {
				pos = p
			} else {
				pos = 0
			}
		} else if v == 'P' {
			wordToPos = append(wordToPos, pos)
		} else {
			pos = acm.AddChar(pos, int(v))
		}
	}
	acm.BuildSuffixLink(true)

	failTree := acm.BuildFailTree()
	trieTree := acm.BuildTrieTree()
	type query struct{ id, value int }
	queryGroup := make([][]query, len(trieTree))

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		node1, node2 := wordToPos[x], wordToPos[y]
		queryGroup[node2] = append(queryGroup[node2], query{id: i, value: node1})
	}

	lid, rid := make([]int, acm.Size()), make([]int, acm.Size())
	dfn := 0
	var dfsOrder func(cur, pre int)
	dfsOrder = func(cur, pre int) {
		lid[cur] = dfn
		dfn++
		for _, next := range failTree[cur] {
			if next != pre {
				dfsOrder(next, cur)
			}
		}
		rid[cur] = dfn
	}
	dfsOrder(0, -1)

	bit := NewBitArray(acm.Size())
	res := make([]int, q)
	var dfs func(cur int, pre int)
	dfs = func(cur int, pre int) {
		bit.Add(lid[cur], 1) // dfs序为fail树的dfs序
		for _, q := range queryGroup[cur] {
			qi, node := q.id, q.value
			res[qi] = bit.QueryRange(lid[node], rid[node])
		}
		for _, next := range trieTree[cur] {
			if next != pre {
				dfs(next, cur)
			}
		}
		bit.Add(lid[cur], -1)
	}
	dfs(0, -1)

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// Indie Album
// https://www.luogu.com.cn/problem/CF1207G
// 有q1次操作,操作有两种类型：
// 1 c : 新建一个字符c.
// 2 i c : 在第i次操作的串后面加上字符c.
// 接着是q2次询问,格式为：
// i t: 每次询问版本为i的串中，t串出现了多少次。
// q1,q2<=4e5, sum(len(text[i]))<=4e5
//
// !看见多字符串匹配，会想到AC自动机
// 相当于：给定一些(posOnTrie, posOnACM)对，查询posOnACM对应的串在posOnTrie对应的串中出现了多少次
// !离线查询，将查询挂在trieTree每个节点上，在trieTree上dfs，树状数组维护 failTree 的 dfs序.
// 类似阿狸的打字机,离线+树状数组.
func CF1207G() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	q1 := io.NextInt()
	trie1 := NewACAutoMatonArray(26, 97)
	posHistory := make([]int, q1)
	for i := 0; i < q1; i++ {
		kind := io.NextInt()
		if kind == 1 {
			c := io.Text()
			posHistory[i] = trie1.AddChar(0, int(c[0]))
		} else {
			version := io.NextInt() - 1
			c := io.Text()
			posHistory[i] = trie1.AddChar(posHistory[version], int(c[0]))
		}
	}
	// trie1.BuildSuffixLink(true)

	q2 := io.NextInt()
	acm2 := NewACAutoMatonArray(26, 97)
	queries := make([][2]int, q2) // !(posOnTrie, posOnACM)：查询posOnACM对应的串在posOnTrie对应的串中出现了多少次
	for i := 0; i < q2; i++ {
		version := io.NextInt() - 1
		text := io.Text()
		textPos := acm2.AddString(text)
		queries[i] = [2]int{posHistory[version], textPos}
	}
	acm2.BuildSuffixLink(true)

	failTree := acm2.BuildFailTree()
	lid, rid := make([]int, acm2.Size()), make([]int, acm2.Size())
	dfn := 0
	var dfsOrder func(cur, pre int)
	dfsOrder = func(cur, pre int) {
		lid[cur] = dfn
		dfn++
		for _, next := range failTree[cur] {
			if next != pre {
				dfsOrder(next, cur)
			}
		}
		rid[cur] = dfn
	}
	dfsOrder(0, -1) // failTree 的 dfs序
	bit := NewBitArray(acm2.Size())

	queryGroup := make([][]int, trie1.Size())
	for qid, query := range queries {
		triePos := query[0]
		queryGroup[triePos] = append(queryGroup[triePos], qid)
	}
	res := make([]int, q2)

	// 在 trie 上 dfs，计算 failTree 的某个节点的子树权值.
	var dfs func(triePos, acmPos int)
	dfs = func(triePos, acmPos int) {
		bit.Add(lid[acmPos], 1)

		for _, qid := range queryGroup[triePos] {
			qv := queries[qid][1]
			res[qid] = bit.QueryRange(lid[qv], rid[qv])
		}

		for c, next := range trie1.children[triePos] {
			if next != -1 {
				dfs(int(next), acm2.Move(acmPos, c+97))
			}
		}

		bit.Add(lid[acmPos], -1)
	}
	dfs(0, 0)

	for _, v := range res {
		io.Println(v)
	}
}

// Genetic engineering
// https://www.luogu.com.cn/problem/CF86C
// 给定w个由ATCG组成的字符串words，求构造一个长度为n的ATCG字符串s，使得:
// 对字符串s任意一个位置i，存在left和right 满足 s[left:right] 至少与一个word匹配.(left<=i<=right)
// (即：字符串完全被这些DNA片段覆盖)
// 求方案总数对1e9+9取模.
// n<=1000,w<=10,len(words[i]<=10).
// dp[index][pos][need]: 前index个字符，当前在ac自动机位置为pos,未匹配need个字符
// !需要判断：匹配好的长度加上最长的DNA的后缀长度是否大于等于整个当前字符串的长度即可（等价于前一部分和后一部分正好拼接或覆盖）
// 如果下一步结点x'是某DNA的后缀且长度比k大, 则将need置为0
func CF86C() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	const MOD int = 1e9 + 9
	getOrd := func(v byte) int {
		if v == 'A' {
			return 0
		} else if v == 'C' {
			return 1
		} else if v == 'G' {
			return 2
		} else {
			return 3
		}
	}
	n, w := io.NextInt(), io.NextInt()
	words := make([]string, w)
	acm := NewACAutoMatonArray(4, 0)
	for i := 0; i < w; i++ {
		word := io.Text()
		words[i] = word
		acm.AddFrom(len(word), func(i int) int { return getOrd(word[i]) })
	}
	acm.BuildSuffixLink(true)

	maxLen := make([]int, acm.Size()) // 每个节点匹配到的最大长度
	for i, p := range acm.WordPos {
		maxLen[p] = max(maxLen[p], len(words[i]))
	}
	acm.Dp(func(from, to int) {
		maxLen[to] = max(maxLen[to], maxLen[from])
	})

	m := acm.Size()
	upper := 0
	for _, v := range words {
		upper = max(upper, len(v))
	}

	{
		// 记忆化搜索求解
		m := acm.Size()
		memo := make([][][]int, n)
		for i := range memo {
			inner1 := make([][]int, m)
			for j := range inner1 {
				inner2 := make([]int, upper+1)
				for k := range inner2 {
					inner2[k] = -1
				}
				inner1[j] = inner2
			}
			memo[i] = inner1
		}

		var dfs func(index, pos, need int) int
		dfs = func(index, pos, need int) int {
			if index == n {
				if need == 0 {
					return 1
				}
				return 0
			}
			if tmp := memo[index][pos][need]; tmp != -1 {
				return tmp
			}
			res := 0
			for c := 0; c < 4; c++ {
				nextPos := acm.Move(pos, c)
				coverLen := maxLen[nextPos]
				nextNeed := need + 1
				if coverLen >= nextNeed {
					nextNeed = 0
				}
				if nextNeed <= upper {
					res += dfs(index+1, nextPos, nextNeed)
					res %= MOD
				}
			}
			memo[index][pos][need] = res
			return res
		}

		io.Println(dfs(0, 0, 0))

	}

	{
		dp := func() {
			makeDp := func(row, col int) [][]int {
				res := make([][]int, row)
				for i := range res {
					res[i] = make([]int, col)
				}
				return res
			}

			dp := makeDp(m, upper+1)
			dp[0][0] = 1
			for i := 0; i < n; i++ {
				ndp := makeDp(m, upper+1)
				for pos := 0; pos < m; pos++ {
					for k := 0; k < upper+1; k++ {
						for c := 0; c < 4; c++ {
							nextPos := acm.Move(pos, c)
							nextK := k + 1
							if maxLen[nextPos] >= nextK {
								nextK = 0
							}
							if nextK <= upper {
								ndp[nextPos][nextK] += dp[pos][k]
								ndp[nextPos][nextK] %= MOD
							}
						}
					}
				}
				dp = ndp
			}

			res := 0
			for i := range dp {
				res = (res + dp[i][0]) % MOD
			}

			io.Println(res)
		}
		_ = dp
	}

}

// Mike and Friends
// https://www.luogu.com.cn/problem/CF547E
// 给定n个字符串words和q个查询，每个查询为：
// !(left, right, index) 查询 words[index]在 [left,right] 中出现了多少次(0<=left<=right<n).
// 将区间查询转换为两个前缀的差.
// 类似阿狸的打字机，模式串 fail 树向下，文本串 trie 树向上
// 沿着trie文本串从根到结束位置点权+1，查询时为fail树某节点子树和.
func CF547E() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, q := io.NextInt(), io.NextInt()
	words := make([]string, n)
	acm := NewACAutoMatonArray(26, 97)
	for i := 0; i < n; i++ {
		words[i] = io.Text()
		acm.AddString(words[i])
	}
	acm.BuildSuffixLink(true)

	queries := make([][3]int, q)
	leftQueryGroup := make([][]int, n)
	rightQueryGroup := make([][]int, n)
	for i := 0; i < q; i++ {
		left, right, index := io.NextInt(), io.NextInt(), io.NextInt()
		left--
		right--
		index--
		queries[i] = [3]int{left, right, index}

		rightQueryGroup[right] = append(rightQueryGroup[right], i)
		if left > 0 {
			leftQueryGroup[left-1] = append(leftQueryGroup[left-1], i)
		}
	}

	failTree := acm.BuildFailTree()
	lid, rid := make([]int, acm.Size()), make([]int, acm.Size())
	dfn := 0
	var dfsOrder func(cur, pre int)
	dfsOrder = func(cur, pre int) {
		lid[cur] = dfn
		dfn++
		for _, next := range failTree[cur] {
			if next != pre {
				dfsOrder(next, cur)
			}
		}
		rid[cur] = dfn
	}
	dfsOrder(0, -1)
	bit := NewBitArray(acm.Size())

	res := make([]int, q)
	for i := 0; i < n; i++ {
		pos := 0
		for _, v := range words[i] {
			pos = acm.Move(pos, int(v))
			bit.Add(lid[pos], 1) // !沿着trie走，从根节点到endPos链上的点加1
		}
		for _, qid := range rightQueryGroup[i] {
			index := queries[qid][2]
			node := acm.WordPos[index]
			res[qid] += bit.QueryRange(lid[node], rid[node])
		}
		for _, qid := range leftQueryGroup[i] {
			index := queries[qid][2]
			node := acm.WordPos[index]
			res[qid] -= bit.QueryRange(lid[node], rid[node])
		}
	}

	for _, v := range res {
		io.Println(v)
	}
}

// Legen...
// https://www.luogu.com.cn/problem/CF696D
func CF696D() {
}

// Duff is Mad
// https://www.luogu.com.cn/problem/CF587F
func CF587F() {
}

// Digits of Number Pi
// https://www.luogu.com.cn/problem/CF585F
func CF585F() {
}

// Birthday
// https://www.luogu.com.cn/problem/CF590E
func CF590E() {
}

// Exam
// https://www.luogu.com.cn/problem/CF1483F
func CF1483F() {
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
	res.Clear()
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

func (trie *ACAutoMatonArray) Empty() bool {
	return len(trie.children) == 1
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

func (trie *ACAutoMatonArray) Clear() {
	trie.WordPos = trie.WordPos[:0]
	trie.Parent = trie.Parent[:0]
	trie.children = trie.children[:0]
	trie.suffixLink = trie.suffixLink[:0]
	trie.bfsOrder = trie.bfsOrder[:0]
	trie.newNode()
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

// 功能与 AddString 相同.
func (ac *ACAutoMatonMap) AddFrom(n int, getOrd func(i int) int) int {
	if n == 0 {
		return 0
	}
	pos := 0
	for i := 0; i < n; i++ {
		ord := getOrd(i)
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

func (ac *ACAutoMatonMap) BuildTrieTree() [][]int {
	res := make([][]int, ac.Size())
	var dfs func(int)
	dfs = func(cur int) {
		for _, next := range ac.children[cur] {
			res[cur] = append(res[cur], int(next))
			dfs(int(next))
		}
	}
	dfs(0)
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

type BITRangeAddPointGetArray struct {
	bit *BITArray
}

func NewBITRangeAddPointGet(n int) *BITRangeAddPointGetArray {
	return &BITRangeAddPointGetArray{bit: NewBitArray(n)}
}

func NewBITRangeAddPointGetFrom(n int, f func(i int) int) *BITRangeAddPointGetArray {
	return &BITRangeAddPointGetArray{bit: NewBitArrayFrom(n, f)}
}

func (b *BITRangeAddPointGetArray) AddRange(start, end int, delta int) {
	if start < 0 {
		start = 0
	}
	if end > b.bit.n {
		end = b.bit.n
	}
	if start >= end {
		return
	}
	b.bit.Add(start, delta)
	b.bit.Add(end, -delta)
}

func (b *BITRangeAddPointGetArray) Get(index int) int {
	return b.bit.QueryPrefix(index + 1)
}

func (b *BITRangeAddPointGetArray) String() string {
	res := []string{}
	for i := 0; i < b.bit.n; i++ {
		res = append(res, fmt.Sprintf("%d", b.Get(i)))
	}
	return fmt.Sprintf("BITRangeAddPointGetArray: [%v]", strings.Join(res, ", "))
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

type S = int

func (*BITPrefixArray) e() S        { return 0 }
func (*BITPrefixArray) op(a, b S) S { return max(a, b) }

type BITPrefixArray struct {
	n    int
	data []S
}

func NewBITPrefixArray(n int) *BITPrefixArray {
	res := &BITPrefixArray{}
	data := make([]S, n)
	for i := range data {
		data[i] = res.e()
	}
	res.n = n
	res.data = data
	return res
}

func NewBITPrefixFrom(n int, f func(index int) S) *BITPrefixArray {
	res := &BITPrefixArray{}
	total := res.e()
	data := make([]S, n)
	for i := range data {
		data[i] = f(i)
		total = res.op(total, data[i])
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = res.op(data[j-1], data[i-1])
		}
	}
	res.n = n
	res.data = data
	return res
}

// 单点更新index处的元素.
// 0 <= index < n
func (f *BITPrefixArray) Update(index int, value S) {
	for index++; index <= f.n; index += index & -index {
		f.data[index-1] = f.op(f.data[index-1], value)
	}
}

// 查询前缀区间 [0,right) 的值.
// 0 <= end <= n
func (f *BITPrefixArray) QueryPrefix(end int) S {
	if end > f.n {
		end = f.n
	}
	res := f.e()
	for ; end > 0; end &= end - 1 {
		res = f.op(res, f.data[end-1])
	}
	return res
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

const INF int = 1e18

// PointSetRangeMax

type E = int

func (*SegmentTree) e() E        { return -INF }
func (*SegmentTree) op(a, b E) E { return max(a, b) }

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(n int, f func(i int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}

func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTree) QueryAll() E { return st.seg[1] }

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if tmp := st.op(res, st.seg[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	right += st.size
	res := st.e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(st.op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := st.op(st.seg[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
}

type H = int

type ErasableHeap struct {
	base   *Heap
	erased *Heap
}

func NewErasableHeap(less func(a, b H) bool, nums []H) *ErasableHeap {
	return &ErasableHeap{NewHeap(less, nums), NewHeap(less, nil)}
}

// 从堆中删除一个元素,要保证堆中存在该元素.
func (h *ErasableHeap) Erase(value H) {
	h.erased.Push(value)
	h.normalize()
}

func (h *ErasableHeap) Push(value H) {
	h.base.Push(value)
	h.normalize()
}

func (h *ErasableHeap) Pop() (value H) {
	value = h.base.Pop()
	h.normalize()
	return
}

func (h *ErasableHeap) Peek() (value H) {
	value = h.base.Top()
	return
}

func (h *ErasableHeap) Len() int {
	return h.base.Len()
}

func (h *ErasableHeap) normalize() {
	for h.base.Len() > 0 && h.erased.Len() > 0 && h.base.Top() == h.erased.Top() {
		h.base.Pop()
		h.erased.Pop()
	}
}

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	if len(nums) > 1 {
		heap.heapify()
	}
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Top() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}

		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
