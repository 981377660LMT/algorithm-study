// 二分图匹配
// https://ei1333.github.io/library/graph/flow/bipartite-matching.hpp

// O(V*E)
// BipartiteMatching(n):= 全体のグラフの頂点数を n で初期化する.
// add_edge(u, v):= 頂点 u, v 間に辺を張る.
// remove_edge(u, v):= 頂点 u, v 間の辺を削除する.
// bipartite_matching():= 二部グラフの最大マッチングを返す.
// erase_vertex(idx):= 頂点 idx を削除し, フローの変化量を返す(-1/0)
// add_vertex(idx):= 頂点 idx を追加し, フローの変化量を返す(0/1)
// getMatchingEdges():= マッチングに使った辺を出力する.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	Rearranging()
	// EdgeDemo()
}

func test() {
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

// https://atcoder.jp/contests/abc317/tasks/abc317_g
// 跑m次匈牙利，每跑一次就删去完美匹配的边
func Rearranging() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)
	grid := make([][]int, ROW)
	for i := 0; i < ROW; i++ {
		grid[i] = make([]int, COL)
		for j := 0; j < COL; j++ {
			var v int
			fmt.Fscan(in, &v)
			grid[i][j] = v - 1
		}
	}

	res := make([][]int, ROW)
	for i := 0; i < ROW; i++ {
		res[i] = make([]int, COL)
	}

	G := NewBipartiteMatching(ROW + ROW)
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			G.AddEdge(i, ROW+grid[i][j])
		}
	}

	for c := 0; c < COL; c++ {
		matching := G.MaxMatching()
		for _, e := range matching {
			r, v := e[0], e[1]-ROW
			res[r][c] = v
			G.RemoveEdge(r, e[1])
		}
	}

	fmt.Fprintln(out, "Yes")
	for _, row := range res {
		for _, v := range row {
			fmt.Fprint(out, v+1, " ")
		}
		fmt.Fprintln(out)
	}
}

func EdgeDemo() {
	G, _, rids := NewBipartiteMathingFromEdges(4, [][2]int{{0, 1}, {1, 2}, {2, 3}})
	M := G.MaxMatching()
	edges := [][2]int{}
	for _, e := range M {
		edges = append(edges, [2]int{rids[e[0]], rids[e[1]]})
	}
	fmt.Println(edges)
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

// left - right
// left: 左侧顶点. 0 <= left < L
// right: 右侧顶点. L <= right < n
func (bm *BipartiteMatching) AddEdge(left, right int) {
	bm.graph[left] = append(bm.graph[left], right)
	bm.graph[right] = append(bm.graph[right], left)
}

func (bm *BipartiteMatching) RemoveEdge(u, v int) {
	for i, to := range bm.graph[u] {
		if to == v {
			bm.graph[u] = append(bm.graph[u][:i], bm.graph[u][i+1:]...)
			break
		}
	}
	for i, to := range bm.graph[v] {
		if to == u {
			bm.graph[v] = append(bm.graph[v][:i], bm.graph[v][i+1:]...)
			break
		}
	}
}

func (bm *BipartiteMatching) MaxMatching() [][2]int {
	// reset match
	for i := range bm.match {
		bm.match[i] = -1
	}

	for i := 0; i < len(bm.graph); i++ {
		if !bm.alive[i] {
			continue
		}
		if bm.match[i] == -1 {
			bm.timestamp++
			bm._augment(i)
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

func (bm *BipartiteMatching) RemoveVertex(idx int) int {
	bm.alive[idx] = false
	if bm.match[idx] == -1 {
		return 0
	}
	bm.match[bm.match[idx]] = -1
	bm.timestamp++
	res := bm._augment(bm.match[idx])
	bm.match[idx] = -1
	if res {
		return 0
	}
	return -1
}

func (bm *BipartiteMatching) AddVertex(idx int) int {
	bm.alive[idx] = true
	bm.timestamp++
	res := bm._augment(idx)
	if res {
		return 1
	}
	return 0
}

// 返回匹配中使用的边.需要在MaxMatching()之后调用.
func (bm *BipartiteMatching) GetMatchingEdges() [][2]int {
	res := [][2]int{}
	for i := 0; i < len(bm.graph); i++ {
		if i < bm.match[i] {
			res = append(res, [2]int{i, bm.match[i]})
		}
	}
	return res
}

func (bm *BipartiteMatching) _augment(idx int) bool {
	bm.used[idx] = bm.timestamp
	for _, to := range bm.graph[idx] {
		toMatch := bm.match[to]
		if !bm.alive[to] {
			continue
		}
		if toMatch == -1 || (bm.used[toMatch] != bm.timestamp && bm._augment(toMatch)) {
			bm.match[idx] = to
			bm.match[to] = idx
			return true
		}
	}
	return false
}

// 从边创建二分图最大匹配.
//
//	Args:
//	 n: 顶点数.
//	 edges: 边集.
//	Returns:
//	 bm: 二分图最大匹配.
//	 ids: 原图中的点在二分图中的编号.左侧点编号为0~L-1, 右侧点编号为L~n-1.
//	 rids: 二分图中的点在原图中的编号.0~n-1 -> 0~n-1.
func NewBipartiteMathingFromEdges(n int, edges [][2]int) (bm *BipartiteMatching, ids, rids []int) {
	graph := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}
	colors, ok := IsBipartite(n, graph)
	if !ok {
		panic("not bipartite")
	}

	ids, rids = make([]int, n), make([]int, n)
	leftCount := 0
	for i := 0; i < n; i++ {
		if colors[i] == 0 {
			leftCount++
		}
	}

	left, right := 0, 0 // 规定左侧点颜色为0, 右侧点颜色为1
	for i := 0; i < n; i++ {
		if colors[i] == 0 {
			ids[i] = left
			rids[left] = i
			left++
		} else {
			ids[i] = right + leftCount
			rids[right+leftCount] = i
			right++
		}
	}

	bm = NewBipartiteMatching(n)
	for _, e := range edges {
		u, v := e[0], e[1]
		if colors[u] == 1 {
			u, v = v, u
		}
		bm.AddEdge(ids[u], ids[v])
	}
	return
}

// 判断是否是二分图.返回 (染色的01数组, 是否是二分图).
func IsBipartite(n int, graph [][]int) (colors []int, ok bool) {
	colors = make([]int, n)
	for i := 0; i < n; i++ {
		colors[i] = -1
	}
	for i := 0; i < n; i++ {
		if colors[i] == -1 {
			colors[i] = 0
			stack := []int{i}
			for len(stack) > 0 {
				cur := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, next := range graph[cur] {
					if colors[next] == -1 {
						colors[next] = 1 ^ colors[cur]
						stack = append(stack, next)
					} else if colors[next] == colors[cur] {
						return nil, false
					}
				}
			}
		}
	}
	return colors, true
}
