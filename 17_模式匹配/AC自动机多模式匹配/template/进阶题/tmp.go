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

const INF int = 1e18

func main() {

	// P2336()

	// SP1676()
	SP9941()

	// CF86C()

	// CF547E()

	// CF1207G()
}

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

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个字符串.
type ACAutoMatonArray struct {
	WordPos            []int     // WordPos[i] 表示加入的第i个模式串对应的节点编号.
	Parent             []int     // parent[v] 表示节点v的父节点.
	Link               []int32   // 又叫fail.指向当前trie节点(对应一个前缀)的最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
	Children           [][]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	BfsOrder           []int32   // 结点的拓扑序,0表示虚拟节点.
	sigma              int32     // 字符集大小.
	offset             int32     // 字符集的偏移量.
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
		if trie.Children[pos][ord] == -1 {
			trie.Children[pos][ord] = trie.newNode()
			trie.Parent[len(trie.Parent)-1] = pos
		}
		pos = int(trie.Children[pos][ord])
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
		if trie.Children[pos][ord] == -1 {
			trie.Children[pos][ord] = trie.newNode()
			trie.Parent[len(trie.Parent)-1] = pos
		}
		pos = int(trie.Children[pos][ord])
	}
	trie.WordPos = append(trie.WordPos, pos)
	return pos
}

// 在pos位置添加一个字符，返回新的节点编号.
func (trie *ACAutoMatonArray) AddChar(pos int, ord int) int {
	ord -= int(trie.offset)
	if trie.Children[pos][ord] != -1 {
		return int(trie.Children[pos][ord])
	}
	trie.Children[pos][ord] = trie.newNode()
	trie.Parent[len(trie.Parent)-1] = pos
	return int(trie.Children[pos][ord])
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *ACAutoMatonArray) Move(pos int, ord int) int {
	ord -= int(trie.offset)
	if trie.needUpdateChildren {
		return int(trie.Children[pos][ord])
	}
	for {
		nexts := trie.Children[pos]
		if nexts[ord] != -1 {
			return int(nexts[ord])
		}
		if pos == 0 {
			return 0
		}
		pos = int(trie.Link[pos])
	}
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *ACAutoMatonArray) Size() int {
	return len(trie.Children)
}

func (trie *ACAutoMatonArray) Empty() bool {
	return len(trie.Children) == 1
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *ACAutoMatonArray) BuildSuffixLink(needUpdateChildren bool) {
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

func (trie *ACAutoMatonArray) Clear() {
	trie.WordPos = trie.WordPos[:0]
	trie.Parent = trie.Parent[:0]
	trie.Children = trie.Children[:0]
	trie.Link = trie.Link[:0]
	trie.BfsOrder = trie.BfsOrder[:0]
	trie.newNode()
}

// 获取每个状态包含的模式串的个数.
func (trie *ACAutoMatonArray) GetCounter() []int {
	counter := make([]int, len(trie.Children))
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

// 获取每个状态包含的模式串的索引.
func (trie *ACAutoMatonArray) GetIndexes() [][]int {
	res := make([][]int, len(trie.Children))
	for i, pos := range trie.WordPos {
		res[pos] = append(res[pos], i)
	}
	for _, v := range trie.BfsOrder {
		if v != 0 {
			from, to := trie.Link[v], v
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
	for _, v := range trie.BfsOrder {
		if v != 0 {
			f(int(trie.Link[v]), int(v))
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
		if pos >= len(trie.Children) || pos < 0 {
			return 0
		}
		ord := int32(s) - trie.offset
		if next := int(trie.Children[pos][ord]); next == -1 {
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
	trie.Children = append(trie.Children, nexts)
	return int32(len(trie.Children) - 1)
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

type Tree struct {
	Tree                 [][][2]int // (next, weight)
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int // 欧拉序[in,out)
	IdToNode             []int
	top, heavySon        []int
	timer                int
}

func NewTree(n int) *Tree {
	tree := make([][][2]int, n)
	lid := make([]int, n)
	rid := make([]int, n)
	IdToNode := make([]int, n)
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
		IdToNode:      IdToNode,
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
			return tree.IdToNode[tree.LID[root]-k]
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
				res = append(res, tree.IdToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.IdToNode[i])
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
	w := tree.IdToNode[k]
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
	tree.IdToNode[tree.timer] = cur
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
