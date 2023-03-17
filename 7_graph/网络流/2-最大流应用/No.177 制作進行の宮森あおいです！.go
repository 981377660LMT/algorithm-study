// No.177 制作進行の宮森あおいです！
// 我是制作进行的宫森葵

// 制作进行新人宮森葵要在两天之内完成w张画
// 现在她可以安排n位原画师和m位作画监督来完成这些画
// 第一天由原画师画,第二天由作画监督完成
// !第i位原画师一天可以画a[i]张画,第i位作画监督一天可以完成b[i]张画
// 但是,
// !第i位作画监督不能画qi位原画师画的画(q0,q1,...,qi-1)
// !问:宫森葵能否在两天之内完成w张画
// w<=1e4 n,m<=50

// 解:
// !问网络中最大流能否达到w

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var w int
	fmt.Fscan(in, &w)
	var n int
	fmt.Fscan(in, &n)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	var m int
	fmt.Fscan(in, &m)
	B := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &B[i])
	}

	ok := make([][]bool, m)
	for i := 0; i < m; i++ {
		ok[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			ok[i][j] = true
		}
	}

	for i := 0; i < m; i++ {
		var q int
		fmt.Fscan(in, &q)
		for j := 0; j < q; j++ {
			var x int
			fmt.Fscan(in, &x)
			ok[i][x-1] = false
		}
	}

	res := 制作進行の宮森あおいです(w, A, B, ok)
	if res {
		fmt.Fprintln(out, "SHIROBAKO")
	} else {
		fmt.Fprintln(out, "BANSAKUTSUKITA")
	}
}

func 制作進行の宮森あおいです(w int, A, B []int, ok [][]bool) bool {
	n, m := len(A), len(B)
	S, T := n+m, n+m+1
	flow := NewPushRelabel(n + m + 2)

	left := func(i int) int { return i }
	right := func(i int) int { return i + n }

	for i := 0; i < n; i++ {
		flow.AddEdge(S, left(i), A[i])
	}
	for i := 0; i < m; i++ {
		flow.AddEdge(right(i), T, B[i])
	}

	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			if ok[j][i] {
				flow.AddEdge(left(i), right(j), INF)
			}
		}
	}

	return flow.MaxFlow(S, T) >= w
}

const INF int = 1e18

type edge struct{ to, rid, cap int }

type PushRelabel struct {
	graph [][]edge
}

func NewPushRelabel(n int) *PushRelabel {
	return &PushRelabel{
		graph: make([][]edge, n),
	}
}

func (pr *PushRelabel) AddEdge(from, to, cap int) {
	pr.graph[from] = append(pr.graph[from], edge{to, len(pr.graph[to]), cap})
	pr.graph[to] = append(pr.graph[to], edge{from, len(pr.graph[from]) - 1, 0})
}

func (pr *PushRelabel) MaxFlow(s, t int) int {
	n := len(pr.graph)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}

	dist[t] = 0
	cd := make([]int, 2*n)
	queue := []int{t}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		cd[dist[v]]++
		for _, e := range pr.graph[v] {
			if w := e.to; dist[w] < 0 {
				dist[w] = dist[v] + 1
				queue = append(queue, w)
			}
		}
	}
	dist[s] = n

	exFlow := make([]int, n)
	pq := hp{d: dist}
	inQueue := make([]bool, n)
	push := func(v, f int, e *edge) {
		w := e.to
		e.cap -= f
		pr.graph[w][e.rid].cap += f
		exFlow[v] -= f
		exFlow[w] += f
		if w != s && w != t && !inQueue[w] {
			pq.push(w)
			inQueue[w] = true
		}
	}

	for i := range pr.graph[s] {
		if e := &pr.graph[s][i]; e.cap > 0 {
			push(s, e.cap, e)
		}
	}

	for pq.Len() > 0 {
		v := pq.pop()
		inQueue[v] = false
	o:
		for {
			for i := range pr.graph[v] {
				if e := &pr.graph[v][i]; e.cap > 0 && dist[e.to] < dist[v] {
					push(v, min(e.cap, exFlow[v]), e)
					if exFlow[v] == 0 {
						break o
					}
				}
			}
			dv := dist[v]
			if dv != -1 {
				if cd[dv]--; cd[dv] == 0 {
					for i, h := range dist {
						if i != s && i != t && dv < h && h <= n {
							dist[i] = n + 1
						}
					}
				}
			}
			minD := INF
			for _, e := range pr.graph[v] {
				if w := e.to; e.cap > 0 && dist[w] < minD {
					minD = dist[w]
				}
			}
			dist[v] = minD + 1
			cd[dist[v]]++
		}
	}

	return exFlow[t]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type hp struct{ vs, d []int }

func (h hp) Len() int              { return len(h.vs) }
func (h hp) Less(i, j int) bool    { return h.d[h.vs[i]] > h.d[h.vs[j]] }
func (h hp) Swap(i, j int)         { h.vs[i], h.vs[j] = h.vs[j], h.vs[i] }
func (h *hp) Push(v interface{})   { h.vs = append(h.vs, v.(int)) }
func (h *hp) Pop() (v interface{}) { a := h.vs; h.vs, v = a[:len(a)-1], a[len(a)-1]; return }
func (h *hp) push(v int)           { heap.Push(h, v) }
func (h *hp) pop() int             { return heap.Pop(h).(int) }
