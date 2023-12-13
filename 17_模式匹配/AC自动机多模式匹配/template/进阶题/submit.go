package main

import (
	"bufio"
	"fmt"
	stdio "io"
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
			bit.Add(lid[pos], 1)
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
