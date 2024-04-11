// 前后缀优化建图1 - 点向区间连边.
//
// PointToPrefixSuffix
// 3 → 4 → 5 (后缀)
// ↓   ↓   ↓
// 0   1   2
// ↑   ↑   ↑
// 6 ← 7 ← 8 (前缀)

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF32 int32 = 1e9 + 10

func main() {
	yuki1868()
}

func yuki1868() {
	// https://yukicoder.me/problems/no/1868
	// !给定一张有向图,每个点i可以向右达到i+1,i+2,...,targets[i]。求从0到n-1的最短路。(前后缀优化建图)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	targets := make([]int32, n-1) // !从i可以到 i+1, i+2, ..., targets[i]
	for i := range targets {
		fmt.Fscan(in, &targets[i])
		targets[i]-- // [0,n-1]内
	}

	R := NewRangeToRangeGraphOnPrefixSuffix1(n)
	adjList := make([][]Neighbor, R.Size())
	R.Init(func(from, to int32) {
		adjList[from] = append(adjList[from], Neighbor{to, 0})
	})
	for i := int32(0); i < n-1; i++ {
		R.AddPointToPrefix(i, targets[i]+1, func(from, to int32) {
			adjList[from] = append(adjList[from], Neighbor{to, 1})
		})
	}

	dist := Bfs0132(adjList, 0)
	fmt.Fprintln(out, dist[n-1])
}

func demo() {
	n := int32(3)
	R := NewRangeToRangeGraphOnPrefixSuffix1(n)
	newGraph := make([][]Neighbor, R.Size())
	R.Init(func(from, to int32) { // from -> to

		newGraph[from] = append(newGraph[from], Neighbor{to, 0})
	})

	R.AddPointToPrefix(2, 2, func(from, to int32) { // from -> [0,prefixEnd)
		newGraph[from] = append(newGraph[from], Neighbor{to, 1})
	})
	R.AddPointToSuffix(0, 0, func(from, to int32) { // from -> [suffixStart,n)
		newGraph[from] = append(newGraph[from], Neighbor{to, 1})
	})
	fmt.Println(newGraph)

	dist := Bfs0132(newGraph, 2)
	fmt.Println(dist[:n])
	dist = Bfs0132(newGraph, 0)
	fmt.Println(dist[:n])

}

type RangeToRangeGraphOnPrefixSuffix1 struct {
	n       int32
	maxSize int32 // [0,n):原始点，[n,n*2):后缀点，[n*2,n*3):前缀点.
}

// 新建一个区间图，n 为原图的节点数.
func NewRangeToRangeGraphOnPrefixSuffix1(n int32) *RangeToRangeGraphOnPrefixSuffix1 {
	return &RangeToRangeGraphOnPrefixSuffix1{n: n, maxSize: 3 * n}
}

// 新图的结点数.前n个节点为原图的节点.
func (g *RangeToRangeGraphOnPrefixSuffix1) Size() int32 { return g.maxSize }

func (g *RangeToRangeGraphOnPrefixSuffix1) Init(f func(from, to int32)) {
	n1, n2 := g.n, g.n*2
	for i := int32(0); i < n1; i++ {
		f(i+n1, i)
		f(i+n2, i)
		if i > 0 {
			f(i+n1-1, i+n1)
			f(i+n2, i+n2-1)
		}
	}
}

// 添加有向边 from -> to.
func (g *RangeToRangeGraphOnPrefixSuffix1) Add(from, to int32, f func(from, to int32)) {
	f(from, to)
}

// from -> [0,prefixEnd)
func (g *RangeToRangeGraphOnPrefixSuffix1) AddPointToPrefix(from int32, prefixEnd int32, f func(from, to int32)) {
	if prefixEnd <= 0 {
		return
	}
	f(from, g.n*2+prefixEnd-1)
}

// from -> [suffixStart,n)
func (g *RangeToRangeGraphOnPrefixSuffix1) AddPointToSuffix(from int32, suffixStart int32, f func(from, to int32)) {
	if suffixStart >= g.n {
		return
	}
	f(from, g.n+suffixStart)
}

type Neighbor struct {
	next int32
	dist int8
}

// 01bfs求最短路.
func Bfs0132(adjList [][]Neighbor, start int32) []int32 {
	n := int32(len(adjList))
	dist := make([]int32, n)
	for i := range dist {
		dist[i] = INF32
	}
	dist[start] = 0
	queue := NewDeque(n)
	queue.Append(start)
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		for _, edge := range adjList[cur] {
			next, weight := edge.next, edge.dist
			cand := dist[cur] + int32(weight)
			if cand < dist[next] {
				dist[next] = cand
				if weight == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}
	return dist
}

type D = int32
type Deque struct{ l, r []D }

func NewDeque(cap int32) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
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

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
