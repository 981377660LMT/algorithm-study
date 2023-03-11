// https://ei1333.github.io/library/graph/flow/bipartite-matching.hpp

// O(V*E)
// BipartiteMatching(n):= 全体のグラフの頂点数を n で初期化する.
// add_edge(u, v):= 頂点 u, v 間に辺を張る.
// bipartite_matching():= 二部グラフの最大マッチングを返す.
// add_vertex(idx):= 頂点 idx を追加し, フローの変化量を返す(0/1)
// erase_vertex(idx):= 頂点 idx を削除し, フローの変化量を返す(-1/0)
// output():= マッチングに使った辺を出力する.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_7_A
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var x, y, e int
	fmt.Fscan(in, &x, &y, &e)
	bm := NewBipartiteMatching(x + y)
	for i := 0; i < e; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		bm.AddEdge(u, v+x)
	}

	fmt.Fprintln(out, len(bm.MaxMatching()))
}

type BipartiteMatching struct {
	timestamp int
	graph     [][]int
	alive     []bool
	used      []int
	match     []int
}

func NewBipartiteMatching(n int) *BipartiteMatching {
	graph := make([][]int, n)
	alive := make([]bool, n)
	used := make([]int, n)
	match := make([]int, n)
	for i := 0; i < n; i++ {
		alive[i] = true
		match[i] = -1
	}
	return &BipartiteMatching{graph: graph, alive: alive, used: used, match: match}
}

// left <-> right
func (bm *BipartiteMatching) AddEdge(u, v int) {
	bm.graph[u] = append(bm.graph[u], v)
	bm.graph[v] = append(bm.graph[v], u)
}

func (bm *BipartiteMatching) MaxMatching() [][2]int {
	for i := 0; i < len(bm.graph); i++ {
		if !bm.alive[i] {
			continue
		}
		if bm.match[i] == -1 {
			bm.timestamp++
			bm.augment(i)
		}
	}

	res := [][2]int{}
	for i := 0; i < len(bm.graph); i++ {
		if i < bm.match[i] {
			res = append(res, [2]int{i, bm.match[i]})
		}
	}
	return res
}

func (bm *BipartiteMatching) AddVertex(idx int) int {
	bm.alive[idx] = true
	bm.timestamp++
	res := bm.augment(idx)
	if res {
		return 1
	}
	return 0
}

func (bm *BipartiteMatching) EraseVertex(idx int) int {
	bm.alive[idx] = false
	if bm.match[idx] == -1 {
		return 0
	}
	bm.match[bm.match[idx]] = -1
	bm.timestamp++
	res := bm.augment(bm.match[idx])
	bm.match[idx] = -1
	if res {
		return 0
	}
	return -1
}

func (bm *BipartiteMatching) augment(idx int) bool {
	bm.used[idx] = bm.timestamp
	for _, to := range bm.graph[idx] {
		toMatch := bm.match[to]
		if !bm.alive[to] {
			continue
		}
		if toMatch == -1 || (bm.used[toMatch] != bm.timestamp && bm.augment(toMatch)) {
			bm.match[idx] = to
			bm.match[to] = idx
			return true
		}
	}
	return false
}
