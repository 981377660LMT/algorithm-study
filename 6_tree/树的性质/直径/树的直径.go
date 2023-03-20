// 求带权树的(直径长度, 直径路径)

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

	var n int
	fmt.Fscan(in, &n)
	td := NewTreeDiameter(n)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		td.AddEdge(u, v, w)
	}
	diameter, path := td.GetDiameter()
	fmt.Fprintln(out, diameter, len(path))
	for _, v := range path {
		fmt.Fprint(out, v, " ")
	}
}

type TreeDiameter struct {
	n    int
	tree [][][2]int
}

func NewTreeDiameter(n int) *TreeDiameter {
	return &TreeDiameter{n: n, tree: make([][][2]int, n)}
}

// 添加一条无向边(u, v)，权值为w.
func (td *TreeDiameter) AddEdge(u, v, w int) {
	td.tree[u] = append(td.tree[u], [2]int{v, w})
	td.tree[v] = append(td.tree[v], [2]int{u, w})
}

// 求带权树的(直径长度, 直径路径).
func (td *TreeDiameter) GetDiameter() (diameter int, path []int) {
	u, _ := td.dfs(0)
	v, dist := td.dfs(u)
	diameter = dist[v]
	path = []int{v}
	for u != v {
		for _, e := range td.tree[v] {
			if dist[e[0]]+e[1] == dist[v] {
				path = append(path, e[0])
				v = e[0]
				break
			}
		}
	}

	return diameter, path
}

// 求树的直径长度和个数.
// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
func (td *TreeDiameter) CountDiameter(start int) (diameter, dCount int) {
	weight := make([]map[int]int, td.n)
	for i := range weight {
		weight[i] = make(map[int]int)
	}
	for i, edges := range td.tree {
		for _, e := range edges {
			weight[i][e[0]] = e[1]
			weight[e[0]][i] = e[1]
		}
	}

	var dfs func(cur, pre int) (int, int)
	dfs = func(cur, pre int) (int, int) {
		maxDepth, count := 0, 1
		for _, e := range td.tree[cur] {
			next := e[0]
			if next != pre {
				nextD, nextC := dfs(next, cur)
				if cand := maxDepth + nextD; cand > diameter {
					diameter, dCount = cand, count*nextC
				} else if cand == diameter {
					dCount += count * nextC
				}
				if nextD > maxDepth {
					maxDepth, count = nextD, nextC
				} else if nextD == maxDepth {
					count += nextC
				}
			}
		}
		return maxDepth + weight[pre][cur], count
	}

	dfs(start, -1)
	return
}

// 求树的直径长度和在直径上的节点个数.
// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
func (td *TreeDiameter) CountVertexOnDiameter(start int) (diameter, vCount int) {
	weight := make([]map[int]int, td.n)
	for i := range weight {
		weight[i] = make(map[int]int)
	}
	for i, edges := range td.tree {
		for _, e := range edges {
			weight[i][e[0]] = e[1]
			weight[e[0]][i] = e[1]
		}
	}

	var dfs func(cur, pre int) (int, int)
	dfs = func(cur, pre int) (int, int) {
		maxDepth, count := 0, 0
		for _, e := range td.tree[cur] {
			next := e[0]
			if next != pre {
				nextD, nextC := dfs(next, cur)
				if cand := maxDepth + nextD; cand > diameter {
					diameter, vCount = cand, count+nextC+1 // 最长的链 + 当前链 + 当前节点
				} else if cand == diameter {
					vCount += nextC
				}
				if nextD > maxDepth {
					maxDepth, count = nextD, nextC
				} else if nextD == maxDepth {
					count += nextC
				}
			}
		}
		return maxDepth + weight[pre][cur], count + 1
	}
	dfs(start, -1)
	return
}

func (td *TreeDiameter) dfs(start int) (int, []int) {
	dist := make([]int, td.n)
	for i := range dist {
		dist[i] = -1
	}
	dist[start] = 0
	stack := []int{start}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, e := range td.tree[cur] {
			next, weight := e[0], e[1]
			if dist[next] != -1 {
				continue
			}
			dist[next] = dist[cur] + weight
			stack = append(stack, next)
		}
	}
	endPoint, maxDist := 0, -1
	for i, d := range dist {
		if d > maxDist {
			endPoint, maxDist = i, d
		}
	}
	return endPoint, dist
}
