// 前后缀优化建图2 - 区间向点连边.
//
// 3 → 4 → 5 (前缀)
// ↑   ↑   ↑
// 0   1   2
// ↓   ↓   ↓
// 6 ← 7 ← 8 (后缀)

package main

import "fmt"

func main() {
	n := int32(3)
	R := NewRangeToRangeGraphOnPrefixSuffix2(n)
	newGraph := make([][]Neighbor, R.Size())
	R.Init(func(from, to int32) { // from -> to
		newGraph[from] = append(newGraph[from], Neighbor{to, 0})
	})

	R.AddPrefixToPoint(2, 2, func(from, to int32) { // from -> [0,prefixEnd)
		newGraph[from] = append(newGraph[from], Neighbor{to, 1})
	})
	R.AddSuffixToPoint(0, 1, func(from, to int32) { // from -> [suffixStart,n)
		newGraph[from] = append(newGraph[from], Neighbor{to, 1})
	})
	fmt.Println(newGraph)

	dist := Bfs0132(newGraph, 2)
	fmt.Println(dist[:n])
	dist = Bfs0132(newGraph, 0)
	fmt.Println(dist[:n])

}

type RangeToRangeGraphOnPrefixSuffix2 struct {
	n       int32
	maxSize int32 // [0,n):原始点，[n,n*2):前缀点，[n*2,n*3):后缀点.

}

// 新建一个区间图，n 为原图的节点数.
func NewRangeToRangeGraphOnPrefixSuffix2(n int32) *RangeToRangeGraphOnPrefixSuffix2 {
	return &RangeToRangeGraphOnPrefixSuffix2{n: n, maxSize: 3 * n}
}

// 新图的结点数.前n个节点为原图的节点.
func (g *RangeToRangeGraphOnPrefixSuffix2) Size() int32 { return g.maxSize }

func (g *RangeToRangeGraphOnPrefixSuffix2) Init(f func(from, to int32)) {
	n1, n2 := g.n, g.n*2
	for i := int32(0); i < n1; i++ {
		f(i, i+n1)
		f(i, i+n2)
		if i > 0 {
			f(i+n1-1, i+n1)
			f(i+n2, i+n2-1)
		}
	}
}

// 添加有向边 from -> to.
func (g *RangeToRangeGraphOnPrefixSuffix2) Add(from, to int32, f func(from, to int32)) {
	f(from, to)
}

// [0,prefixEnd) -> to.
func (g *RangeToRangeGraphOnPrefixSuffix2) AddPrefixToPoint(prefixEnd int32, to int32, f func(from, to int32)) {
	if prefixEnd <= 0 {
		return
	}
	f(g.n+prefixEnd-1, to)
}

// [suffixStart,n) -> to.
func (g *RangeToRangeGraphOnPrefixSuffix2) AddSuffixToPoint(suffixStart int32, to int32, f func(from, to int32)) {
	if suffixStart >= g.n {
		return
	}
	f(g.n*2+suffixStart, to)
}

const INF int32 = 1e9 + 10

type Neighbor struct {
	next int32
	dist int8
}

// 01bfs求最短路.
func Bfs0132(adjList [][]Neighbor, start int32) []int32 {
	n := int32(len(adjList))
	dist := make([]int32, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := NewDeque32(n)
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
