// 前后缀优化建图3 - 区间向区间连边.
//
// 3 → 4 → 5 (前缀1)         9  → 10 → 11 (后缀2)
// ↑   ↑   ↑                 ↓    ↓     ↓
// 0   1   2          ---    0    1     2
// ↓   ↓   ↓                 ↑    ↑     ↑
// 6 ← 7 ← 8 (后缀1)         12 ← 13 ← 14 (前缀2)

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	solve1 := func(n int32, prefixToSuffix [][2]int32, suffixToPrefix [][2]int32, start int32) []int32 {
		graph := make([][]Neightbor, n)
		for _, e := range prefixToSuffix {
			prefixEnd, suffixStart := e[0], e[1]
			for from := int32(0); from < prefixEnd; from++ {
				for to := suffixStart; to < n; to++ {
					graph[from] = append(graph[from], Neightbor{to, 1})
				}
			}
		}
		for _, e := range suffixToPrefix {
			suffixStart, prefixEnd := e[0], e[1]
			for from := suffixStart; from < n; from++ {
				for to := int32(0); to < prefixEnd; to++ {
					graph[from] = append(graph[from], Neightbor{to, 1})
				}
			}
		}

		dist := Bfs0132(graph, start)
		return dist
	}

	solve2 := func(n int32, prefixToSuffix [][2]int32, suffixToPrefix [][2]int32, start int32) []int32 {
		rangeToRangeOpCount := int32(len(prefixToSuffix) + len(suffixToPrefix))
		R := NewRangeToRangeGraphOnPrefixSuffix3(n, rangeToRangeOpCount)
		newGraph := make([][]Neightbor, R.Size())
		R.Init(func(from, to int32) {
			newGraph[from] = append(newGraph[from], Neightbor{to, 0})
		})

		for _, e := range prefixToSuffix {
			R.AddPrefixToSuffix(e[0], e[1], func(from, to int32) {
				newGraph[from] = append(newGraph[from], Neightbor{to, 1})
			})
		}

		for _, e := range suffixToPrefix {
			R.AddSuffixToPrefix(e[0], e[1], func(from, to int32) {
				newGraph[from] = append(newGraph[from], Neightbor{to, 1})
			})
		}

		dist := Bfs0132(newGraph, start)
		dist = dist[:n]
		return dist
	}

	n := int32(rand.Intn(100) + 1)
	prefixToSuffix := make([][2]int32, 0)
	suffixToPrefix := make([][2]int32, 0)
	len1, len2 := rand.Intn(int(n))/4, rand.Intn(int(n))/4
	for i := 0; i < len1; i++ {
		prefixEnd, suffixStart := int32(rand.Intn(int(n))), int32(rand.Intn(int(n)))
		prefixToSuffix = append(prefixToSuffix, [2]int32{prefixEnd, suffixStart})
	}
	for i := 0; i < len2; i++ {
		suffixStart, prefixEnd := int32(rand.Intn(int(n))), int32(rand.Intn(int(n)))
		suffixToPrefix = append(suffixToPrefix, [2]int32{suffixStart, prefixEnd})
	}

	start := int32(rand.Intn(int(n)))
	dist1 := solve1(n, prefixToSuffix, suffixToPrefix, start)
	dist2 := solve2(n, prefixToSuffix, suffixToPrefix, start)
	for i := int32(0); i < n; i++ {
		if dist1[i] != dist2[i] {
			panic("not equal")
		}
	}
	fmt.Println("pass", len1, len2)
}

type RangeToRangeGraphOnPrefixSuffix3 struct {
	n int32
	// [0,n):原始点，[n,n*2):前缀1，[n*2,n*3):后缀1，[n*3,n*4):后缀2，[n*4,n*5):前缀2,[n*5,n*5+rangeToRangeOpCount):新建立的点.
	maxSize  int32
	allocPtr int32
}

// 新建一个区间图，n 为原图的节点数，rangeToRangeOpCount 为区间到区间的最大操作次数.
func NewRangeToRangeGraphOnPrefixSuffix3(n int32, rangeToRangeOpCount int32) *RangeToRangeGraphOnPrefixSuffix3 {
	return &RangeToRangeGraphOnPrefixSuffix3{n: n, maxSize: 5*n + rangeToRangeOpCount, allocPtr: 5 * n}
}

// 新图的结点数.前n个节点为原图的节点.
func (g *RangeToRangeGraphOnPrefixSuffix3) Size() int32 { return g.maxSize }

func (g *RangeToRangeGraphOnPrefixSuffix3) Init(f func(from, to int32)) {
	n1, n2, n3, n4 := g.n, g.n*2, g.n*3, g.n*4
	for i := int32(0); i < n1; i++ {
		f(i, i+n1)
		f(i, i+n2)
		f(i+n3, i)
		f(i+n4, i)
		if i > 0 {
			f(i+n1-1, i+n1)
			f(i+n3-1, i+n3)
			f(i+n2, i+n2-1)
			f(i+n4, i+n4-1)
		}
	}
}

// 添加有向边 from -> to.
func (g *RangeToRangeGraphOnPrefixSuffix3) Add(from, to int32, f func(from, to int32)) {
	f(from, to)
}

// [0,prefixEnd) -> to.
func (g *RangeToRangeGraphOnPrefixSuffix3) AddPrefixToPoint(prefixEnd int32, to int32, f func(from, to int32)) {
	if prefixEnd <= 0 {
		return
	}
	f(g.n+prefixEnd-1, to)
}

// [suffixStart,n) -> to.
func (g *RangeToRangeGraphOnPrefixSuffix3) AddSuffixToPoint(suffixStart int32, to int32, f func(from, to int32)) {
	if suffixStart >= g.n {
		return
	}
	f(g.n*2+suffixStart, to)
}

func (g *RangeToRangeGraphOnPrefixSuffix3) AddPointToPrefix(from int32, prefixEnd int32, f func(from, to int32)) {
	if prefixEnd <= 0 {
		return
	}
	f(from, g.n*4+prefixEnd-1)
}

func (g *RangeToRangeGraphOnPrefixSuffix3) AddPointToSuffix(from int32, suffixStart int32, f func(from, to int32)) {
	if suffixStart >= g.n {
		return
	}
	f(from, g.n*3+suffixStart)
}

func (g *RangeToRangeGraphOnPrefixSuffix3) AddPrefixToSuffix(prefixEnd, suffixStart int32, f func(from, to int32)) {
	newNode := g.allocPtr
	g.allocPtr++
	g.AddPrefixToPoint(prefixEnd, newNode, f)
	g.AddPointToSuffix(newNode, suffixStart, f)
}

func (g *RangeToRangeGraphOnPrefixSuffix3) AddSuffixToPrefix(suffixStart, prefixEnd int32, f func(from, to int32)) {
	newNode := g.allocPtr
	g.allocPtr++
	g.AddSuffixToPoint(suffixStart, newNode, f)
	g.AddPointToPrefix(newNode, prefixEnd, f)
}

const INF int32 = 1e9 + 10

type Neightbor struct {
	next int32
	dist int8
}

// 01bfs求最短路.
func Bfs0132(adjList [][]Neightbor, start int32) []int32 {
	n := int32(len(adjList))
	dist := make([]int32, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := NewDeque32(n)
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

func NewDeque32(cap int32) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

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
