package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
	"strings"
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

// N 頂点の有向グラフ
// G
// S
// ​
//   があり、頂点には
// 1 から
// N までの番号が付けられています。
// G
// S
// ​
//   には
// N−1 本の辺があり、
// i 本目
// (1≤i≤N−1) の辺は頂点
// p
// i
// ​
//   (1≤p
// i
// ​
//  ≤i) から頂点
// i+1 に伸びています。

// N 頂点の有向グラフ
// G があり、頂点には
// 1 から
// N までの番号が付けられています。 最初、
// G は
// G
// S
// ​
//   と一致しています。
// G に関するクエリが
// Q 個与えられるので、与えられた順番に処理してください。クエリは次の
// 2 種類のいずれかです。

// 1 u v :
// G に頂点
// u から頂点
// v に伸びる辺を追加する。 このとき、以下の条件が満たされることが保証される。
// u
// 
// =v
// G
// S
// ​
//   上で頂点
// v からいくつかの辺を辿ることで頂点
// u に到達可能である
// 2 x :
// G 上で頂点
// x からいくつかの辺を辿ることで到達可能な頂点 (頂点
// x を含む) のうち、最も番号が小さい頂点の番号を出力する。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	uf := NewUnionFindArray(n)
	minChild := make([]int, n)
	for i := 0; i < n; i++ {
		minChild[i] = i
	}
	graph := make([][]int, n)
	for i := 0; i < n-1; i++ {
		pre := io.NextInt() - 1
		// pre -> i+1
		uf.Union(pre, i+1)
		graph[pre] = append(graph[pre], i+1)
	}

	// dfs找每个点里最小的点

	var dfs func(u int) int
	dfs = func(u int) int {
		minChild[u] = u
		for _, v := range graph[u] {
			minChild[u] = min(minChild[u], dfs(v))
		}
		return minChild[u]
	}

	q := io.NextInt()
	for i := 0; i < q; i++ {
		op := io.NextInt()
		if op == 1 {
			// from -> to 加边
			// 之前 to可以到达from
			from, to := io.NextInt()-1, io.NextInt()-1
			uf.Union(from, to)
		} else {
			// 查询x可以到达的最小的点
			// x := io.NextInt() - 1

		}
	}

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

type V = int

type leftistHeap struct {
	less         func(v1, v2 V) bool
	isPersistent bool
}

type heap struct {
	Value       V
	Id          int
	height      int // 维持平衡
	left, right *heap
}

// less: 小根堆返回 v1 < v2, 大根堆返回 v1 > v2
// isPersistent: 是否持久化, 如果是, 则每次合并操作都会返回一个新的堆, 否则会直接修改原堆.
func NewLeftistHeap(less func(v1, v2 V) bool, isPersisitent bool) *leftistHeap {
	return &leftistHeap{less: less, isPersistent: isPersisitent}
}

func (lh *leftistHeap) Alloc(value V, id int) *heap {
	res := &heap{Value: value, Id: id, height: 1}
	return res
}

func (lh *leftistHeap) Build(nums []V) []*heap {
	res := make([]*heap, len(nums))
	for i, num := range nums {
		res[i] = lh.Alloc(num, i)
	}
	return res
}

func (lh *leftistHeap) Push(heap *heap, value V, id int) *heap {
	return lh.Meld(heap, lh.Alloc(value, id))
}

func (lh *leftistHeap) Pop(heap *heap) *heap {
	return lh.Meld(heap.left, heap.right)
}

func (lh *leftistHeap) Top(heap *heap) V {
	return heap.Value
}

// 合并两个堆,返回合并后的堆.
func (lh *leftistHeap) Meld(heap1, heap2 *heap) *heap {
	if heap1 == nil {
		return heap2
	}
	if heap2 == nil {
		return heap1
	}
	if lh.less(heap2.Value, heap1.Value) {
		heap1, heap2 = heap2, heap1
	}
	heap1 = lh.clone(heap1)
	heap1.right = lh.Meld(heap1.right, heap2)
	if heap1.left == nil || heap1.left.height < heap1.right.height {
		heap1.left, heap1.right = heap1.right, heap1.left
	}
	heap1.height = 1
	if heap1.right != nil {
		heap1.height += heap1.right.height
	}
	return heap1
}

func (h *heap) String() string {
	var sb []string
	var dfs func(h *heap)
	dfs = func(h *heap) {
		if h == nil {
			return
		}
		sb = append(sb, fmt.Sprintf("%d", h.Value))
		dfs(h.left)
		dfs(h.right)
	}
	dfs(h)
	return strings.Join(sb, " ")
}

// 持久化,拷贝一份结点.
func (lh *leftistHeap) clone(h *heap) *heap {
	if h == nil || !lh.isPersistent {
		return h
	}
	res := &heap{height: h.height, Value: h.Value, Id: h.Id, left: h.left, right: h.right}
	return res
}

// NewUnionFindWithCallback ...
func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	cb(root2, root1)
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
