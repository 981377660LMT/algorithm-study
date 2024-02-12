// https://atcoder.jp/contests/abc245/tasks/abc245_g
//
// 给 N 个点，M 条边，每个点的颜色（值域为 [1,K]）。并给定 L 个点作为特殊点。
// 询问 每个点到最近的与其颜色不同的特殊点的距离（无解输出 -1） 。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M, K, L int
	fmt.Fscan(in, &N, &M, &K, &L)
	colors := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &colors[i])
		colors[i]--
	}
	criticals := make([]int, L)
	for i := 0; i < L; i++ {
		fmt.Fscan(in, &criticals[i])
		criticals[i]--
	}
	edges := make([][3]int, M)
	for i := 0; i < M; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}

	res := Solve(N, edges, colors, criticals)
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

const INF int = 1e18

func Solve(
	n int, edges [][3]int, colors []int, criticals []int,
) []int {
	adjList := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adjList[u] = append(adjList[u], [2]int{v, w})
		adjList[v] = append(adjList[v], [2]int{u, w})
	}

	dist1, dist2 := make([]int, n), make([]int, n)
	fromColor1, fromColor2 := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		dist1[i], dist2[i] = INF, INF
		fromColor1[i], fromColor2[i] = -1, -1
	}
	pq := NewHeap(func(a, b H) bool { return a[0] < b[0] }, nil)
	for _, v := range criticals {
		dist1[v] = 0
		fromColor1[v] = colors[v]
		pq.Push(H{0, v, colors[v]})
	}

	for pq.Len() > 0 {
		item := pq.Pop()
		curDist, cur, curColor := item[0], item[1], item[2]
		if dist1[cur] != curDist && dist2[cur] != curDist {
			continue
		}
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			nextDist := curDist + weight
			if nextDist < dist1[next] {
				if fromColor1[next] != curColor {
					dist2[next] = dist1[next]
					fromColor2[next] = fromColor1[next]
				}
				dist1[next] = nextDist
				fromColor1[next] = curColor
				pq.Push(H{dist1[next], next, curColor})
			} else if nextDist < dist2[next] && fromColor1[next] != curColor {
				dist2[next] = nextDist
				fromColor2[next] = curColor
				pq.Push(H{dist2[next], next, curColor})
			}
		}
	}

	res := make([]int, n)
	for i := 0; i < n; i++ {
		if fromColor1[i] != colors[i] {
			res[i] = dist1[i]
		} else {
			res[i] = dist2[i]
		}
	}
	for i := 0; i < n; i++ {
		if res[i] == INF {
			res[i] = -1
		}
	}
	return res

}

type H = [3]int

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	heap.heapify()
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
