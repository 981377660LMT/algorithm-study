// RangeToRangeGraph (区间图)
// !原图的连通分量/最短路在新图上仍然等价

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	// https://yukicoder.me/problems/no/1868
	// !给定一张有向图,每个点i可以向右达到i+1,i+2,...,targets[i]。求从0到n-1的最短路。
	// 解法1：每个点i连接targets[i],边权为1,所有i到i-1连边,边权为0。然后跑最短路。
	// 解法2：RangeToRangeGraph。每个点i连接i+1,i+2,...,targets[i]。然后跑最短路。
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	targets := make([]int, n-1) // !从i可以到 i+1, i+2, ..., targets[i]
	for i := range targets {
		fmt.Fscan(in, &targets[i])
		targets[i]-- // [0,n-1]内
	}

	R := NewRangeToRangeGraph(n)
	for i := 0; i < n-1; i++ {
		R.AddToRange(i, i+1, targets[i]+1, 1) // 左闭右开
	}
	adjList, newN := R.Build()

	dist, queue := make([]int, newN), NewDeque(newN)
	for i := range dist {
		dist[i] = INF
	}
	dist[0] = 0
	queue.Append(0)
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
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

	fmt.Fprintln(out, dist[n-1])
}

type RangeToRangeGraph struct {
	n     int
	nNode int
	edges [][3]int // [from, to, weight]
}

func NewRangeToRangeGraph(n int) *RangeToRangeGraph {
	g := &RangeToRangeGraph{
		n:     n,
		nNode: n * 3,
	}
	for i := 2; i < n+n; i++ {
		g.edges = append(g.edges, [3]int{g.toUpperIdx(i / 2), g.toUpperIdx(i), 0})
	}
	for i := 2; i < n+n; i++ {
		g.edges = append(g.edges, [3]int{g.toLowerIdx(i), g.toLowerIdx(i / 2), 0})
	}
	return g
}

// 添加有向边 from -> to, 权重为 weight.
func (g *RangeToRangeGraph) Add(from, to int, weight int) {
	g.edges = append(g.edges, [3]int{from, to, weight})
}

// 从区间 [fromStart, fromEnd) 中的每个点到 to 都添加一条有向边，权重为 weight.
func (g *RangeToRangeGraph) AddFromRange(fromStart, fromEnd, to int, weight int) {
	l, r := fromStart+g.n, fromEnd+g.n
	for l < r {
		if l&1 == 1 {
			g.Add(g.toLowerIdx(l), to, weight)
			l++
		}
		if r&1 == 1 {
			r--
			g.Add(g.toLowerIdx(r), to, weight)
		}
		l >>= 1
		r >>= 1
	}
}

// 从 from 到区间 [toStart, toEnd) 中的每个点都添加一条有向边，权重为 weight.
func (g *RangeToRangeGraph) AddToRange(from, toStart, toEnd int, weight int) {
	l, r := toStart+g.n, toEnd+g.n
	for l < r {
		if l&1 == 1 {
			g.Add(from, g.toUpperIdx(l), weight)
			l++
		}
		if r&1 == 1 {
			r--
			g.Add(from, g.toUpperIdx(r), weight)
		}
		l >>= 1
		r >>= 1
	}
}

// 从区间 [fromStart, fromEnd) 中的每个点到区间 [toStart, toEnd) 中的每个点都添加一条有向边，权重为 weight.
func (g *RangeToRangeGraph) AddRangeToRange(fromStart, fromEnd, toStart, toEnd int, weight int) {
	newNode := g.nNode
	g.nNode++
	g.AddFromRange(fromStart, fromEnd, newNode, weight)
	g.AddToRange(newNode, toStart, toEnd, 0)
}

// 返回`新图的有向邻接表和新图的节点数`.
func (g *RangeToRangeGraph) Build() (graph [][][2]int, vertex int) {
	graph = make([][][2]int, g.nNode)
	for _, e := range g.edges {
		u, v, w := e[0], e[1], e[2]
		graph[u] = append(graph[u], [2]int{v, w})
	}
	return graph, g.nNode
}

func (g *RangeToRangeGraph) toUpperIdx(i int) int {
	if i >= g.n {
		return i - g.n
	}
	return g.n + i
}

func (g *RangeToRangeGraph) toLowerIdx(i int) int {
	if i >= g.n {
		return i - g.n
	}
	return g.n + g.n + i
}

//
//
type D = int
type Deque struct{ l, r []D }

func NewDeque(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

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
