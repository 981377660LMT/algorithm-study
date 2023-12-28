// 二分图染色
// 二分图等价于不存在奇环.

package main

import "fmt"

func main() {
	fmt.Println(BipartiteVertexColoringRemoveOneEdge(3, 3, [][]int{{1, 2}, {2, 0}, {0, 1}}))
}

// 无向图二分图着色.
// 如果不是二分图，返回空数组.
func BipartiteVertexColoring(graph [][]int) (colors []int8) {
	n := len(graph)
	uf := NewUf(2 * n)
	for cur, nexts := range graph {
		for _, next := range nexts {
			if cur < next {
				uf.Union(cur+n, next)
				uf.Union(cur, next+n)
			}
		}
	}
	colors = make([]int8, 2*n)
	for i := range colors {
		colors[i] = -1
	}
	for v := 0; v < n; v++ {
		if root := uf.Find(v); root == v && colors[root] < 0 {
			colors[root] = 0
			colors[uf.Find(v+n)] = 1
		}
	}
	for v := 0; v < n; v++ {
		colors[v] = colors[uf.Find(v)]
	}
	colors = colors[:n]
	for v := 0; v < n; v++ {
		if uf.Find(v) == uf.Find(v+n) {
			return nil
		}
	}
	return
}

type Uf struct {
	data []int
}

func NewUf(n int) *Uf {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &Uf{data: data}
}

func (ufa *Uf) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *Uf) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

// 删除一条边，使得图变成二分图.
// 返回新的二分图着色.
// 如果无法删除，返回空数组.
// https://www.luogu.com.cn/problem/CF1680F
func BipartiteVertexColoringRemoveOneEdge(n, m int, graph [][]int) []int {
	vs := make([]int, 0, n)
	parent, color := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = -1
		color[i] = -1
	}
	visited := make(map[int]struct{}, m)
	dp0, dp1 := make([]int, n), make([]int, n)
	getEid := func(v1, v2 int) int {
		if v1 < v2 {
			v1, v2 = v2, v1
		}
		return v1*n + v2
	}

	odd := 0
	var dfs func(int)
	dfs = func(cur int) {
		vs = append(vs, cur)
		for _, next := range graph[cur] {
			eid := getEid(cur, next)
			if _, ok := visited[eid]; ok {
				continue
			}
			visited[eid] = struct{}{}
			if color[next] == -1 {
				parent[next] = cur
				color[next] = color[cur] ^ 1
				dfs(next)
			} else {
				if color[cur] != color[next] {
					dp0[cur]++
					dp0[next]--
				} else {
					dp1[cur]++
					dp1[next]--
					odd++
				}
			}
		}
	}
	for v := 0; v < n; v++ {
		if color[v] == -1 {
			color[v] = 0
			dfs(v)
		}
	}
	for i := n - 1; i >= 0; i-- {
		v := vs[i]
		p := parent[v]
		if p == -1 {
			continue
		}
		dp0[p] += dp0[v]
		dp1[p] += dp1[v]
	}
	if odd <= 1 {
		return color
	}
	for v := 0; v < n; v++ {
		if parent[v] == -1 || dp1[v] != odd || (dp0[v] != 0 && dp1[v] != 0) {
			continue
		}
		for _, w := range vs {
			if parent[w] == -1 {
				continue
			}
			var tmp int
			if w != v {
				tmp = 1
			}
			color[w] = color[parent[w]] ^ tmp
		}
		return color
	}

	return nil
}
