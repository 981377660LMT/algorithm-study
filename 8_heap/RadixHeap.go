// Radix Heap（基数堆）是一种单调排列堆（添加的值必须比最后一个取出的值大）的一种。
// 由于比std::priority_queue更轻，因此可用于Dijkstra算法以实现常数倍的加速。

// https://nyaannyaan.github.io/library/data-structure/radix-heap.hpp
// RadixHeap<Key, Val>(): コンストラクタ。Keyは整数型のみを取る。
// push(Key, Val): ヒープにpushする。
// pop(): ヒープの要素のうち最小のものをpopして返す。
// size(): ヒープ内の要素数を返す。
// empty(): ヒープが空かどうかを返す。
// !key 必须是整数类型, 这里用int(实际应该用uint)

// !验证: 没有加速效果，反而变慢了

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_1_A
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, s int
	fmt.Fscan(in, &n, &m, &s)
	g := make([][][2]int, n)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		g[u] = append(g[u], [2]int{v, w})
	}
	dist := DijkRadixHeap(g, s)
	for _, d := range dist {
		if d == -1 {
			fmt.Fprintln(out, "INF")
		} else {
			fmt.Fprintln(out, d)
		}
	}
}

type Value = int

type RadixHeap struct {
	vs [64 + 1][]struct {
		key int
		val Value
	}
	ms   [64 + 1]int
	size int
	last int
}

func NewRadixHeap() *RadixHeap {
	res := &RadixHeap{}
	for i := range res.ms {
		res.ms[i] = 1 << 61
	}
	return res
}

func (rh *RadixHeap) Empty() bool {
	return rh.size == 0
}

func (rh *RadixHeap) Size() int {
	return rh.size
}

func (rh *RadixHeap) getBit(a int) int {
	return 64 - bits.LeadingZeros64(uint64(a))
}

func (rh *RadixHeap) Push(key int, val Value) {
	rh.size++
	b := rh.getBit(key ^ rh.last)
	rh.vs[b] = append(rh.vs[b], struct {
		key int
		val Value
	}{key, val})
	rh.ms[b] = min(rh.ms[b], key)
}

func (rh *RadixHeap) Pop() (key int, value Value) {
	if rh.ms[0] == 1<<61 {
		idx := 1
		for rh.ms[idx] == 1<<61 {
			idx++
		}
		rh.last = rh.ms[idx]
		for _, p := range rh.vs[idx] {
			b := rh.getBit(p.key ^ rh.last)
			rh.vs[b] = append(rh.vs[b], p)
			rh.ms[b] = min(p.key, rh.ms[b])
		}
		rh.vs[idx] = rh.vs[idx][:0]
		rh.ms[idx] = 1 << 61
	}
	rh.size--
	res := rh.vs[0][len(rh.vs[0])-1]
	rh.vs[0] = rh.vs[0][:len(rh.vs[0])-1]
	if len(rh.vs[0]) == 0 {
		rh.ms[0] = 1 << 61
	}
	return res.key, res.val
}

// RadixHeap优化的Dijkstra算法, 不可达时距离为-1.
func DijkRadixHeap(graph [][][2]int, start int) (dist []int) {
	n := len(graph)
	dist = make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	heap := NewRadixHeap()
	dist[start] = 0
	heap.Push(0, start)
	for !heap.Empty() {
		key, cur := heap.Pop()
		if dist[cur] < key {
			continue
		}
		for _, dst := range graph[cur] {
			next, weight := dst[0], dst[1]
			if dist[next] == -1 || dist[cur]+weight < dist[next] {
				dist[next] = dist[cur] + weight
				heap.Push(dist[next], next)
			}
		}
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
