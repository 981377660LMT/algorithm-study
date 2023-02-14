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
	diameter, path := td.CalDiameter()
	fmt.Fprintln(out, diameter, len(path))
	for _, v := range path {
		fmt.Fprint(out, v, " ")
	}
}

type TreeDiameter struct {
	n    int
	tree [][][]int
}

func NewTreeDiameter(n int) *TreeDiameter {
	return &TreeDiameter{n: n, tree: make([][][]int, n)}
}

// 添加一条无向边(u, v)，权值为w.
func (td *TreeDiameter) AddEdge(u, v, w int) {
	td.tree[u] = append(td.tree[u], []int{v, w})
	td.tree[v] = append(td.tree[v], []int{u, w})
}

// 求带权树的(直径长度, 直径路径).
func (td *TreeDiameter) CalDiameter() (int, []int) {
	u, _ := td.dfs(0)
	v, dist := td.dfs(u)
	diameter := dist[v]
	path := []int{v}
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
