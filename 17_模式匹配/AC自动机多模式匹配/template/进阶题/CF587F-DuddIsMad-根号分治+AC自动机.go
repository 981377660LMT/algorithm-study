package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

// CF587F-DuddIsMad-根号分治+AC自动机
// https://www.luogu.com.cn/problem/CF587F
// !给定n个字符串，q次询问[start:end)中的字符串在s[index]中出现次数之和.
// 字符串总长度<=1e5.
//
// 按照查询的文本串的长度进行根号分治.
// !文本串为短串时，将所有模式串中的短串按照扫描线的顺序加入.扫描线+前缀和+树状数组处理.
// 对每个(短)模式串，加入时fail树子树和+1, 文本串查询时累加其前缀结点之和.类似mike and friends中的技巧.
// !文本串为长串时，这样的串不超过sqrt(n)个.
// 转化为求出每个模式串在文本串中出现的次数(P3966 单词).
// 文本串每个前缀+1，dfs求出子树和即可.
// ps: 也可以kmp处理出每个模式串在此长串中出现次数.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const THRESHOLD int32 = 300

	var n, q int32
	fmt.Fscan(in, &n, &q)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	acam := NewACAutoMatonArray32()
	for _, w := range words {
		acam.AddString(w)
	}
	acam.BuildSuffixLink(true)
	wordPos := acam.WordPos

	queries := make([][3]int32, q)              // left, right, index
	leftQueryGroupOfShort := make([][]int32, n) // !用于处理短文本串的匹配
	rightQueryGroupOfShort := make([][]int32, n)
	queryGroupOfLong := make([][]int32, n) // !用于处理长文本串的匹配

	for qi := int32(0); qi < q; qi++ {
		var left, right, index int32
		fmt.Fscan(in, &left, &right, &index)
		left--
		right--
		index--
		queries[qi] = [3]int32{left, right, index}
		word := words[index]
		if int32(len(word)) < THRESHOLD {
			rightQueryGroupOfShort[right] = append(rightQueryGroupOfShort[right], qi)
			if left > 0 {
				leftQueryGroupOfShort[left-1] = append(leftQueryGroupOfShort[left-1], qi)
			}
		} else {
			queryGroupOfLong[index] = append(queryGroupOfLong[index], qi)
		}
	}

	size := acam.Size()
	failTree := acam.BuildFailTree()
	lid, rid := make([]int32, size), make([]int32, size) // failTree 的 dfs序
	dfn := int32(0)
	var dfsOrder func(cur int32)
	dfsOrder = func(cur int32) {
		lid[cur] = dfn
		dfn++
		for _, next := range failTree[cur] {
			dfsOrder(next)
		}
		rid[cur] = dfn
	}
	dfsOrder(0)

	// 文本串为短串时，将所有模式串按照扫描线的顺序加入
	res := make([]int, q)
	bit := NewBITRangeBlockRangeAddPointGet32(size)
	for i := int32(0); i < n; i++ {
		curPos := wordPos[i]
		bit.AddRange(lid[curPos], rid[curPos], 1) // 模式串的子树+1，查询时和为文本串每个前缀结点之和

		queryChain := func(wid int32) int { // 前缀结点权值之和
			word := words[wid]
			pos := int32(0)
			curSum := 0
			for _, c := range word {
				pos = acam.Move(pos, c)
				curSum += bit.Get(lid[pos])
			}
			return curSum
		}

		for _, qid := range rightQueryGroupOfShort[i] {
			wid := queries[qid][2]
			res[qid] += queryChain(wid)
		}
		for _, qid := range leftQueryGroupOfShort[i] {
			wid := queries[qid][2]
			res[qid] -= queryChain(wid)
		}
	}

	// 文本串为长串时，这样的串不超过sqrt(n)个
	// 转化为求出每个模式串在文本串中出现的次数
	for i := int32(0); i < n; i++ {
		if len(queryGroupOfLong[i]) == 0 {
			continue
		}
		word := words[i]
		values := make([]int32, size)
		pos := int32(0)
		for _, c := range word {
			pos = acam.Move(pos, c)
			values[pos]++
		}
		subValues := make([]int32, size)
		var dfs func(cur int32)
		dfs = func(cur int32) {
			subValues[cur] = values[cur]
			for _, next := range failTree[cur] {
				dfs(next)
				subValues[cur] += subValues[next]
			}
		}
		dfs(0)
		preSum := make([]int, n+1)
		for j := int32(0); j < n; j++ {
			pos := wordPos[j]
			preSum[j+1] = preSum[j] + int(subValues[pos])
		}
		for _, qid := range queryGroupOfLong[i] {
			left, right := queries[qid][0], queries[qid][1]
			res[qid] = preSum[right+1] - preSum[left]
		}
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

const SIGMA int32 = 26  // 字符集大小
const OFFSET int32 = 97 // 字符集偏移量

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个前缀.
type ACAutoMatonArray32 struct {
	WordPos            []int32        // WordPos[i] 表示加入的第i个模式串对应的节点编号(单词结点).
	Parent             []int32        // parent[v] 表示节点v的父节点.
	Link               []int32        // 又叫fail.指向当前trie节点(对应一个前缀)的最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
	Children           [][SIGMA]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	BfsOrder           []int32        // 结点的拓扑序,0表示虚拟节点.
	needUpdateChildren bool           // 是否需要更新children数组.
}

func NewACAutoMatonArray32() *ACAutoMatonArray32 {
	res := &ACAutoMatonArray32{}
	res.Clear()
	return res
}

// 添加一个字符串，返回最后一个字符对应的节点编号.
func (trie *ACAutoMatonArray32) AddString(str string) int32 {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for _, s := range str {
		ord := s - OFFSET
		if trie.Children[pos][ord] == -1 {
			trie.Children[pos][ord] = trie.newNode()
			trie.Parent[len(trie.Parent)-1] = pos
		}
		pos = trie.Children[pos][ord]
	}
	trie.WordPos = append(trie.WordPos, pos)
	return pos
}

// 在pos位置添加一个字符，返回新的节点编号.
func (trie *ACAutoMatonArray32) AddChar(pos, ord int32) int32 {
	ord -= OFFSET
	if trie.Children[pos][ord] != -1 {
		return trie.Children[pos][ord]
	}
	trie.Children[pos][ord] = trie.newNode()
	trie.Parent[len(trie.Parent)-1] = pos
	return trie.Children[pos][ord]
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *ACAutoMatonArray32) Move(pos, ord int32) int32 {
	ord -= OFFSET
	if trie.needUpdateChildren {
		return trie.Children[pos][ord]
	}
	for {
		nexts := trie.Children[pos]
		if nexts[ord] != -1 {
			return nexts[ord]
		}
		if pos == 0 {
			return 0
		}
		pos = trie.Link[pos]
	}
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *ACAutoMatonArray32) Size() int32 {
	return int32(len(trie.Children))
}

func (trie *ACAutoMatonArray32) Empty() bool {
	return len(trie.Children) == 1
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *ACAutoMatonArray32) BuildSuffixLink(needUpdateChildren bool) {
	trie.needUpdateChildren = needUpdateChildren
	trie.Link = make([]int32, len(trie.Children))
	for i := range trie.Link {
		trie.Link[i] = -1
	}
	trie.BfsOrder = make([]int32, len(trie.Children))
	head, tail := 0, 0
	trie.BfsOrder[tail] = 0
	tail++
	for head < tail {
		v := trie.BfsOrder[head]
		head++
		for i, next := range trie.Children[v] {
			if next == -1 {
				continue
			}
			trie.BfsOrder[tail] = next
			tail++
			f := trie.Link[v]
			for f != -1 && trie.Children[f][i] == -1 {
				f = trie.Link[f]
			}
			trie.Link[next] = f
			if f == -1 {
				trie.Link[next] = 0
			} else {
				trie.Link[next] = trie.Children[f][i]
			}
		}
	}
	if !needUpdateChildren {
		return
	}
	for _, v := range trie.BfsOrder {
		for i, next := range trie.Children[v] {
			if next == -1 {
				f := trie.Link[v]
				if f == -1 {
					trie.Children[v][i] = 0
				} else {
					trie.Children[v][i] = trie.Children[f][i]
				}
			}
		}
	}
}

func (trie *ACAutoMatonArray32) Clear() {
	trie.WordPos = trie.WordPos[:0]
	trie.Parent = trie.Parent[:0]
	trie.Children = trie.Children[:0]
	trie.Link = trie.Link[:0]
	trie.BfsOrder = trie.BfsOrder[:0]
	trie.newNode()
}

// 获取每个状态包含的模式串的个数.
func (trie *ACAutoMatonArray32) GetCounter() []int32 {
	counter := make([]int32, len(trie.Children))
	for _, pos := range trie.WordPos {
		counter[pos]++
	}
	for _, v := range trie.BfsOrder {
		if v != 0 {
			counter[v] += counter[trie.Link[v]]
		}
	}
	return counter
}

// 获取每个状态包含的模式串的索引.(模式串长度和较小时使用)
func (trie *ACAutoMatonArray32) GetIndexes() [][]int32 {
	res := make([][]int32, len(trie.Children))
	for i, pos := range trie.WordPos {
		res[pos] = append(res[pos], int32(i))
	}
	for _, v := range trie.BfsOrder {
		if v != 0 {
			from, to := trie.Link[v], v
			arr1, arr2 := res[from], res[to]
			arr3 := make([]int32, 0, len(arr1)+len(arr2))
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

// 按照拓扑序进行转移(EnumerateFail).
func (trie *ACAutoMatonArray32) Dp(f func(from, to int32)) {
	for _, v := range trie.BfsOrder {
		if v != 0 {
			f(trie.Link[v], v)
		}
	}
}

func (trie *ACAutoMatonArray32) BuildFailTree() [][]int32 {
	res := make([][]int32, trie.Size())
	trie.Dp(func(pre, cur int32) {
		res[pre] = append(res[pre], cur)
	})
	return res
}

func (trie *ACAutoMatonArray32) BuildTrieTree() [][]int32 {
	res := make([][]int32, trie.Size())
	for i := int32(1); i < trie.Size(); i++ {
		res[trie.Parent[i]] = append(res[trie.Parent[i]], i)
	}
	return res
}

func (trie *ACAutoMatonArray32) newNode() int32 {
	trie.Parent = append(trie.Parent, -1)
	nexts := [SIGMA]int32{}
	for i := range nexts {
		nexts[i] = -1
	}
	trie.Children = append(trie.Children, nexts)
	return int32(len(trie.Children) - 1)
}

type Str = string

func GetNext32(pattern Str) []int32 {
	n := int32(len(pattern))
	next := make([]int32, n)
	j := int32(0)
	for i := int32(1); i < n; i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = int32(next[j-1])
		}
		if pattern[i] == pattern[j] {
			j++
		}
		next[i] = j
	}
	return next
}

func CountIndexOfAll32(longer Str, shorter Str, position int32, nexts []int32) int32 {
	if len(shorter) == 0 {
		return 0
	}
	if len(longer) < len(shorter) {
		return 0
	}
	res := int32(0)
	if nexts == nil {
		nexts = GetNext32(shorter)
	}
	n := int32(len(longer))
	m := int32(len(shorter))
	hitJ := int32(0)
	for i := position; i < n; i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == m {
			res++
			hitJ = nexts[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

// 基于分块实现的`树状数组`.
// `O(1)`单点查询，`O(sqrt(n))`区间加.
// 一般配合莫队算法使用.
type BITRangeBlockRangeAddPointGet32 struct {
	_n          int32
	_belong     []int32
	_blockStart []int32
	_blockEnd   []int32
	_nums       []int
	_blockLazy  []int
}

func NewBITRangeBlockRangeAddPointGet32(n int32) *BITRangeBlockRangeAddPointGet32 {
	blockSize := int32(math.Sqrt(float64(n)) + 1)
	blockCount := 1 + (n / blockSize)
	belong := make([]int32, n)
	for i := range belong {
		belong[i] = int32(i) / blockSize
	}
	blockStart := make([]int32, blockCount)
	blockEnd := make([]int32, blockCount)
	for i := range blockStart {
		blockStart[i] = int32(i) * blockSize
		tmp := (int32(i) + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	nums := make([]int, n)
	blockSum := make([]int, blockCount)
	res := &BITRangeBlockRangeAddPointGet32{
		_n:          n,
		_belong:     belong,
		_blockStart: blockStart,
		_blockEnd:   blockEnd,
		_nums:       nums,
		_blockLazy:  blockSum,
	}
	return res
}

func NewBITRangeBlockRangeAddPointGet32From(n int32, f func(i int32) int) *BITRangeBlockRangeAddPointGet32 {
	res := NewBITRangeBlockRangeAddPointGet32(n)
	res.Build(n, f)
	return res
}

func (b *BITRangeBlockRangeAddPointGet32) Get(index int32) int {
	if index < 0 || index >= b._n {
		panic("index out of range")
	}
	return b._nums[index] + b._blockLazy[b._belong[index]]
}

func (b *BITRangeBlockRangeAddPointGet32) AddRange(start, end int32, delta int) {
	if start < 0 {
		start = 0
	}
	if end > b._n {
		end = b._n
	}
	if start >= end {
		return
	}
	bid1 := b._belong[start]
	bid2 := b._belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			b._nums[i] += delta
		}
		return
	}
	for i := start; i < b._blockEnd[bid1]; i++ {
		b._nums[i] += delta
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		b._blockLazy[bid] += delta
	}
	for i := b._blockStart[bid2]; i < end; i++ {
		b._nums[i] += delta
	}
}

func (b *BITRangeBlockRangeAddPointGet32) Build(n int32, f func(i int32) int) {
	if n != b._n {
		panic("array length mismatch n")
	}
	for i := range b._nums {
		b._nums[i] = f(int32(i))
	}
	for i := range b._blockLazy {
		b._blockLazy[i] = 0
	}
}

func (b *BITRangeBlockRangeAddPointGet32) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := range b._nums {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", b.Get(int32(i))))
	}
	sb.WriteString("}")
	return sb.String()
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

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}
