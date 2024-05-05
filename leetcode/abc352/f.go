package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/rand"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
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

// N 人の人がおり、人にはそれぞれ
// 1,2,…,N の番号が付けられています。

// N 人が競争を行い、順位が付きました。この順位に対して以下の情報が与えられています。

// それぞれの人に対して付けられた順位は相異なる
// 各
// 1≤i≤M について人
// A
// i
// ​
//   の順位を
// x、人
// B
// i
// ​
//   の順位を
// y とすると、
// x−y=C
// i
// ​
//   である
// ただし、この問題では与えられた情報に矛盾しないような順位付けが
// 1 つ以上存在するような入力のみが与えられます。

// N 個のクエリの答えを求めてください。
// i 番目のクエリの答えは以下により定まる整数です。

// 人
// i の順位が一意に定まるならば、その値を答えとする。そうでない場合、答えは
// −1 である。

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N, M := io.NextInt(), io.NextInt()
	A, B, C := make([]int, M), make([]int, M), make([]int, M)
	D := NewDifferSystem(int32(N) + 2)
	for i := 0; i < M; i++ {
		A[i], B[i], C[i] = io.NextInt(), io.NextInt(), io.NextInt()
		A[i]--
		B[i]--
		D.GreaterThanOrEqualTo(int32(A[i]), int32(B[i]), C[i])
	}

	time1 := time.Now()
	LOWER, UPPER := int32(N), int32(N+1)
	for i := 0; i < N; i++ {
		D.GreaterThanOrEqualTo(int32(i), LOWER, 0)
		D.LessThanOrEqualTo(int32(i), UPPER, N-1)
	}

	checkPerm := func(perm []int) bool {
		for i := 0; i < M; i++ {
			if perm[A[i]]-perm[B[i]] != C[i] {
				return false
			}
		}
		return true
	}

	res := make([]int, N)
	for i := 0; i < N; i++ {
		res[i] = -2
	}

	perm := make([]int, N)
	for i := 0; i < N; i++ {
		perm[i] = i
	}

	for i := 0; i < N+2; i++ {
		D.RunSince(UPPER)
		P := make([]int, N)
		for i := 0; i < N; i++ {
			P[i] = D.PossibleSolutionOf(int32(i)) + 1
		}
		if checkPerm(P) {
			for i := 0; i < N; i++ {
				if res[i] == -1 {
					continue
				}
				if res[i] == -2 {
					res[i] = perm[i] + 1
				} else if res[i] != perm[i]+1 {
					res[i] = -1
				}
			}
		}
	}

	for {
		rand.Shuffle(N, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
		if checkPerm(perm) {
			for i := 0; i < N; i++ {
				if res[i] == -1 {
					continue
				}
				if res[i] == -2 {
					res[i] = perm[i] + 1
				} else if res[i] != perm[i]+1 {
					res[i] = -1
				}
			}
		}
		if time.Since(time1) > 1800*time.Millisecond {
			break
		}
	}

	for i := 0; i < N; i++ {
		io.Println(res[i])
	}
}

const INF int = 2e18

// !默认最大化解.如果要最小化解, 需要令 y = -x, 然后求最大化解(即将原来的边反向).
type DifferSystem struct {
	nodes       []*DifferNode
	queue       *Deque
	n           int32
	allPos      bool
	hasSolution bool
}

func NewDifferSystem(n int32) *DifferSystem {
	res := &DifferSystem{
		nodes:  make([]*DifferNode, n),
		queue:  NewDeque(n),
		n:      n,
		allPos: true,
	}
	for i := int32(0); i < n; i++ {
		node := NewNode()
		node.id = i
		res.nodes[i] = node
	}
	return res
}

// d[i] - d[j] <= d
func (ds *DifferSystem) LessThanOrEqualTo(i, j int32, d int) {
	ds.nodes[j].adj = append(ds.nodes[j].adj, NewNeighbor(ds.nodes[i], d))
	ds.allPos = ds.allPos && d >= 0
}

// d[i] - d[j] >= d
func (ds *DifferSystem) GreaterThanOrEqualTo(i, j int32, d int) {
	ds.LessThanOrEqualTo(j, i, -d)
}

// d[i] - d[j] = d
func (ds *DifferSystem) EqualTo(i, j int32, d int) {
	ds.GreaterThanOrEqualTo(i, j, d)
	ds.LessThanOrEqualTo(i, j, d)
}

// d[i] - d[j] < d
func (ds *DifferSystem) LessThan(i, j int32, d int) {
	ds.LessThanOrEqualTo(i, j, d-1)
}

// d[i] - d[j] > d
func (ds *DifferSystem) GreaterThan(i, j int32, d int) {
	ds.GreaterThanOrEqualTo(i, j, d+1)
}

func (ds *DifferSystem) HasSolution() bool {
	ds._prepare(0)
	for i := int32(0); i < ds.n; i++ {
		ds.nodes[i].inQueue = true
		ds.queue.Append(ds.nodes[i])
	}
	ds.hasSolution = ds._spfa()
	return ds.hasSolution
}

// Find max(ai - aj), if INF is returned, it means no constraint between ai and aj.
func (ds *DifferSystem) FindMaxDifferenceBetween(i, j int32) int {
	ds.RunSince(j)
	return ds.nodes[i].dist
}

// Find min(ai - aj), if INF is returned, it means no constraint between ai and aj.
func (ds *DifferSystem) FindMinDifferenceBetween(i, j int32) int {
	r := ds.FindMaxDifferenceBetween(j, i) // !最小化解
	if r == INF {
		return INF
	}
	return -r
}

// After invoking this method, the value of i is max(ai - aj).
func (ds *DifferSystem) RunSince(j int32) bool {
	ds._prepare(INF)
	ds.queue.Clear()
	ds.queue.Append(ds.nodes[j])
	ds.nodes[j].dist = 0
	ds.nodes[j].inQueue = true
	ds.hasSolution = ds._spfa()
	return ds.hasSolution
}

func (ds *DifferSystem) PossibleSolutionOf(i int32) int {
	return ds.nodes[i].dist
}

func (ds *DifferSystem) Clear(n int32) {
	ds.n = n
	ds.allPos = true
	for i := int32(0); i < n; i++ {
		ds.nodes[i].adj = ds.nodes[i].adj[:0]
	}
}

func (ds *DifferSystem) _spfa() bool {
	if ds.allPos {
		ds._dijkstra()
		return true
	}
	queue := ds.queue
	for queue.Size() > 0 {
		head := queue.PopLeft()
		head.inQueue = false
		if head.times >= ds.n {
			return false
		}
		for _, edge := range head.adj {
			node := edge.next
			if node.dist <= edge.weight+head.dist {
				continue
			}
			node.dist = edge.weight + head.dist
			if node.inQueue {
				continue
			}
			node.times++
			node.inQueue = true
			queue.AppendLeft(node)
		}
	}
	return true
}

func (ds *DifferSystem) _prepare(initDist int) {
	ds.queue.Clear()
	for i := int32(0); i < ds.n; i++ {
		node := ds.nodes[i]
		node.dist = initDist
		node.times = 0
		node.inQueue = false
	}
}

func (ds *DifferSystem) _dijkstra() {
	pq := NewErasableHeapGeneric[*DifferNode](
		func(a, b *DifferNode) bool {
			if a.dist == b.dist {
				return a.id < b.id
			}
			return a.dist < b.dist
		},
		ds.queue.GetAll()...,
	)
	for pq.Len() > 0 {
		head := pq.Pop()
		for _, e := range head.adj {
			if e.next.dist <= head.dist+e.weight {
				continue
			}
			pq.Erase(e.next)
			e.next.dist = head.dist + e.weight
			pq.Push(e.next)
		}
	}
}

func (ds *DifferSystem) String() string {
	sb := strings.Builder{}
	for i := int32(0); i < ds.n; i++ {
		for _, edge := range ds.nodes[i].adj {
			sb.WriteString(edge.String())
			sb.WriteByte('\n')
		}
	}
	sb.WriteString("-------------\n")
	if !ds.hasSolution {
		sb.WriteString("impossible")
	} else {
		for i := int32(0); i < ds.n; i++ {
			sb.WriteString(fmt.Sprintf("a%d=%d\n", i, ds.nodes[i].dist))
		}
	}
	return sb.String()
}

type Neighbor struct {
	next   *DifferNode
	weight int
}

func NewNeighbor(next *DifferNode, weight int) *Neighbor {
	return &Neighbor{next: next, weight: weight}
}

func (n *Neighbor) String() string {
	return fmt.Sprintf("next=%d, weight=%d", n.next.id, n.weight)
}

type DifferNode struct {
	adj     []*Neighbor
	dist    int
	inQueue bool
	times   int32
	id      int32
}

func NewNode() *DifferNode {
	return &DifferNode{}
}

type ErasableHeapGeneric[H comparable] struct {
	data   *HeapGeneric[H]
	erased *HeapGeneric[H]
	size   int
}

func NewErasableHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *ErasableHeapGeneric[H] {
	return &ErasableHeapGeneric[H]{NewHeapGeneric(less, nums...), NewHeapGeneric(less), len(nums)}
}

// 从堆中删除一个元素,要保证堆中存在该元素.
func (h *ErasableHeapGeneric[H]) Erase(value H) {
	h.erased.Push(value)
	h.normalize()
	h.size--
}

func (h *ErasableHeapGeneric[H]) Push(value H) {
	h.data.Push(value)
	h.normalize()
	h.size++
}

func (h *ErasableHeapGeneric[H]) Pop() (value H) {
	value = h.data.Pop()
	h.normalize()
	h.size--
	return
}

func (h *ErasableHeapGeneric[H]) Peek() (value H) {
	value = h.data.Top()
	return
}

func (h *ErasableHeapGeneric[H]) Len() int {
	return h.size
}

func (h *ErasableHeapGeneric[H]) Clear() {
	h.data.Clear()
	h.erased.Clear()
	h.size = 0
}

func (h *ErasableHeapGeneric[H]) normalize() {
	for h.data.Len() > 0 && h.erased.Len() > 0 && h.data.Top() == h.erased.Top() {
		h.data.Pop()
		h.erased.Pop()
	}
}

type HeapGeneric[H comparable] struct {
	data []H
	less func(a, b H) bool
}

func NewHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *HeapGeneric[H] {
	nums = append(nums[:0:0], nums...)
	heap := &HeapGeneric[H]{less: less, data: nums}
	if len(nums) > 1 {
		heap.heapify()
	}
	return heap
}

func (h *HeapGeneric[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *HeapGeneric[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *HeapGeneric[H]) Top() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *HeapGeneric[H]) Len() int { return len(h.data) }

func (h *HeapGeneric[H]) Clear() {
	h.data = h.data[:0]
}

func (h *HeapGeneric[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *HeapGeneric[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *HeapGeneric[H]) pushDown(root int) {
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

type D = *DifferNode
type Deque struct{ l, r []D }

func NewDeque(cap int32) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q *Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q *Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q *Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q *Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q *Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}

func (q *Deque) Clear() {
	q.l = q.l[:0]
	q.r = q.r[:0]
}

func (q *Deque) GetAll() []D {
	return append(q.l, q.r...)
}
