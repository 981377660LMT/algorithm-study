package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

// from https://atcoder.jp/users/ccppjsrb
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

// Indie Album (阿狸的打字机，一个串在另一个串中出现了多少次)
// https://www.luogu.com.cn/problem/CF1207G
// 有q1次操作,操作有两种类型：
// 1 c : 新建一个字符c.
// 2 i c : 在第i次操作的串后面加上字符c.
// 接着是q2次询问,格式为：
// i t: 每次询问版本为i的串中，t串出现了多少次。
// q1,q2<=4e5, sum(len(text[i]))<=4e5
//
// !看见多字符串匹配，会想到AC自动机
// 相当于：给定一些(triePos, failTreePos)对，查询failTreePos对应的串在triePos对应的串中出现了多少次
//
// 1.现在trie上找到所有模式串的位置；
// !2.离线查询，将查询挂在trieTree每个节点上.
// 3.在trie上dfs，计算 failTree 节点的子树权值.
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	q1 := io.NextInt()
	acm := NewACAutoMatonArray32()
	wordPos := make([]int32, q1)
	for i := 0; i < q1; i++ {
		kind := io.NextInt()
		if kind == 1 {
			c := io.Text()
			wordPos[i] = acm.AddChar(0, int32(c[0]))
		} else {
			version := io.NextInt() - 1
			c := io.Text()
			wordPos[i] = acm.AddChar(wordPos[version], int32(c[0]))
		}
	}

	q2 := io.NextInt()
	queries := make([][2]int32, q2) // !(posOnTrie, posOnACM)：查询posOnACM对应的串在posOnTrie对应的串中出现了多少次
	for i := 0; i < q2; i++ {
		version := io.NextInt() - 1
		text := io.Text()
		textPos := acm.AddString(text)
		queries[i] = [2]int32{wordPos[version], textPos}
	}
	acm.BuildSuffixLink(true)

	size := acm.Size()
	trieTree := acm.BuildTrieTree()
	failTree := acm.BuildFailTree()
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

	bit := NewBitArray(size)
	queryGroup := make([][]int, size)
	for qid, query := range queries {
		triePos := query[0]
		queryGroup[triePos] = append(queryGroup[triePos], qid)
	}
	res := make([]int, q2)

	// 在 trie 上 dfs，计算 failTree 的某个节点的子树权值.
	var dfs func(cur int32)
	dfs = func(cur int32) {
		bit.Add(lid[cur], 1)
		for _, qid := range queryGroup[cur] {
			qv := queries[qid][1]
			res[qid] = bit.QueryRange(lid[qv], rid[qv])
		}
		for _, next := range trieTree[cur] {
			dfs(next)
		}
		bit.Add(lid[cur], -1)
	}
	dfs(0)

	for _, v := range res {
		io.Println(v)
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

// 返回str在trie树上的节点位置.如果不存在，返回0.
func (trie *ACAutoMatonArray32) Search(str string) int32 {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for _, char := range str {
		if pos >= int32(len(trie.Children)) || pos < 0 {
			return 0
		}
		ord := char - OFFSET
		if next := trie.Children[pos][ord]; next == -1 {
			return 0
		} else {
			pos = next
		}
	}
	return pos
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

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int32
	total int
	data  []int
}

func NewBitArray(n int32) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int32, f func(i int32) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := int32(1); i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int32, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int32) int {
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
func (b *BITArray) QueryRange(start, end int32) int {
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
