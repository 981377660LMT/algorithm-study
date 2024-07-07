// AC自动机上转移k次不含模式串的方案数.
// 矩阵快速幂
// 这里的图就是trie图.

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
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

const INF int = 1e18

// GEN - Text Generator
// https://www.luogu.com.cn/problem/SP1676
// P4052 [JSOI2007] 文本生成器 的加强版.
// 给定n(n<=10)个字符串(长度<=6)，要求构造一个长度为L(<=10^6)的字符串，至少包含n个字符串中任何一个，求方案数mod10007的值(所有字符都是大写字母)
// 输入为多组数据，每组数据第一行N,L，接下来N行每行一个字符串
//
// 用总方案数减去不合法方案数(不包含模式串).
// !等价于：有向图中转移k次后到达目的地的方案数.
// !其中adjMatrix表示转移矩阵，只有当u和v都不是接受状态，adjMatrix[u][v]才为true.
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	const MOD int = 10007
	matMul := func(m1, m2 [][]int, mod int) [][]int {
		res := make([][]int, len(m1))
		for i := range res {
			res[i] = make([]int, len(m2[0]))
		}
		for i := 0; i < len(m1); i++ {
			for k := 0; k < len(m2); k++ {
				for j := 0; j < len(m2[0]); j++ {
					res[i][j] = (res[i][j] + m1[i][k]*m2[k][j]) % mod
					if res[i][j] < 0 {
						res[i][j] += mod
					}
				}
			}
		}
		return res
	}

	matPow := func(mat [][]int, exp int, mod int) [][]int {
		n := len(mat)
		matCopy := make([][]int, n)
		for i := 0; i < n; i++ {
			matCopy[i] = make([]int, n)
			copy(matCopy[i], mat[i])
		}
		res := make([][]int, n)
		for i := 0; i < n; i++ {
			res[i] = make([]int, n)
			res[i][i] = 1
		}
		for exp > 0 {
			if exp&1 == 1 {
				res = matMul(res, matCopy, mod)
			}
			matCopy = matMul(matCopy, matCopy, mod)
			exp >>= 1
		}
		return res
	}

	qpow := func(a, b, mod int) int {
		a %= mod
		res := 1
		for b > 0 {
			if b&1 == 1 {
				res = res * a % mod
			}
			a = a * a % mod
			b >>= 1
		}
		return res
	}

	n, k := io.NextInt(), io.NextInt()
	acm := NewACAutoMatonArray(26, 'A')
	for i := 0; i < n; i++ {
		s := io.Text()
		acm.AddString(s)
	}
	acm.BuildSuffixLink(true)

	counter := acm.GetCounter()
	size := acm.Size()

	adjMatrix := make([][]int, size)
	for i := range adjMatrix {
		adjMatrix[i] = make([]int, size)
	}

	for cur := 0; cur < size; cur++ {
		if counter[cur] == 0 {
			for c := 0; c < 26; c++ {
				next := acm.children[cur][c]
				if counter[next] == 0 {
					adjMatrix[cur][next]++
				}
			}
		}
	}

	transition := matPow(adjMatrix, k, MOD)
	bad := 0
	for _, v := range transition[0] {
		bad = (bad + v) % MOD
	}

	res := (qpow(26, k, MOD) - bad) % MOD
	if res < 0 {
		res += MOD
	}

	io.Println(res)

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
func (trie *ACAutoMatonArray) GetIndexes() [][]int32 {
	res := make([][]int32, len(trie.children))
	for i, pos := range trie.WordPos {
		res[pos] = append(res[pos], int32(i))
	}
	for _, v := range trie.bfsOrder {
		if v != 0 {
			from, to := trie.suffixLink[v], v
			arr1, arr2 := res[from], res[to]
			arr3 := make([]int32, len(arr1)+len(arr2))
			i, j, k := 0, 0, 0
			for i < len(arr1) && j < len(arr2) {
				if arr1[i] < arr2[j] {
					arr3[k] = arr1[i]
					i++
				} else if arr1[i] > arr2[j] {
					arr3[k] = arr2[j]
					j++
				} else {
					arr3[k] = arr1[i]
					i++
					j++
				}
				k++
			}
			copy(arr3[k:], arr1[i:])
			k += len(arr1) - i
			copy(arr3[k:], arr2[j:])
			k += len(arr2) - j
			arr3 = arr3[:k:k]
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

// 返回str在trie树上的节点位置.如果不存在，返回0.
func (trie *ACAutoMatonArray) Search(str string) int {
	if len(str) == 0 {
		return 0
	}
	pos := 0
	for _, s := range str {
		if pos >= len(trie.children) || pos < 0 {
			return 0
		}
		ord := int32(s) - trie.offset
		if next := int(trie.children[pos][ord]); next == -1 {
			return 0
		} else {
			pos = next
		}
	}
	return pos
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
