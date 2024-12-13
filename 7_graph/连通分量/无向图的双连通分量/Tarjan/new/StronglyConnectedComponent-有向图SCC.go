// StronglyConnectedComponent-有向图SCC

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yuki1813()
	// yosupo()
	// yuki1293()
}

func yuki1813() {
	// https://yukicoder.me/problems/no/1813
	// 不等关系:有向边; 全部相等:强连通(环)
	// 给定一个DAG 求将DAG变为一个环(强连通分量)的最少需要添加的边数
	// !答案为 `max(入度为0的点的个数, 出度为0的点的个数)`

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	graph := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		graph[u] = append(graph[u], v)
	}
	count, belong := StronglyConnectedComponent(graph)
	if count == 1 { // 缩成一个点了,说明是强连通的
		fmt.Fprintln(out, 0)
		return
	}

	dag := SCCDag(graph, count, belong)
	indeg, outDeg := make([]int, count), make([]int, count)
	for i := 0; i < count; i++ {
		for _, next := range dag[i] {
			indeg[next]++
			outDeg[i]++
		}
	}

	in0, out0 := 0, 0
	for i := 0; i < count; i++ {
		if indeg[i] == 0 {
			in0++
		}
		if outDeg[i] == 0 {
			out0++
		}
	}

	fmt.Fprintln(out, max(in0, out0))
}

func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	graph := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		graph[u] = append(graph[u], v)
	}

	count, belong := StronglyConnectedComponent(graph)
	group := make([][]int, count)
	for i := 0; i < n; i++ {
		group[belong[i]] = append(group[belong[i]], i)
	}
	fmt.Fprintln(out, count)
	for _, p := range group {
		fmt.Fprint(out, len(p))
		for _, v := range p {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}

func yuki1293() {
	// https://yukicoder.me/problems/no/1293
	// No.1293 2種類の道路-SCC
	// 无向图中有两种路径,各有road1,road2条
	// 求有多少个二元组(a,b),满足从a到b经过 '若干条第一种路径+若干条第二种路径'

	// !每个点i拆成点2*i和点2*i+1,2*i->2*i+1
	// !第一种路径: 2*i<->2*j
	// !第二种路径: 2*i+1<->2*j+1
	// 然后对每个顶点求出有多少个可以到达自己

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, road1, road2 int
	fmt.Fscan(in, &n, &road1, &road2)
	graph := make([][]int, 2*n)
	for i := 0; i < road1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		graph[2*a] = append(graph[2*a], 2*b)
		graph[2*b] = append(graph[2*b], 2*a)
	}
	for i := 0; i < road2; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		graph[2*a+1] = append(graph[2*a+1], 2*b+1)
		graph[2*b+1] = append(graph[2*b+1], 2*a+1)
	}
	for i := 0; i < n; i++ {
		graph[2*i] = append(graph[2*i], 2*i+1)
	}

	count, belong := StronglyConnectedComponent(graph)
	dag := SCCDag(graph, count, belong)
	dp := make([]int, count)
	for i := 0; i < n; i++ {
		dp[belong[2*i]]++
	}
	for i := 0; i < count; i++ {
		for _, to := range dag[i] {
			dp[to] += dp[i]
		}
	}

	res := 0
	for i := 0; i < n; i++ {
		res += dp[belong[2*i+1]] - 1 // !减去自己到自己的路径1
	}
	fmt.Fprintln(out, res)
}

// 3383. 施法所需最低符文数量
// https://leetcode.cn/problems/minimum-runes-to-add-to-cast-spell/description/
// !在不可达子DAG中，有多少个入度为0的SCC，就需要为其添加多少条新的有向符文
func minRunesToAdd(n int, crystals []int, flowFrom []int, flowTo []int) int {
	graph := make([][]int, n)
	for i := range flowFrom {
		from_, to_ := flowFrom[i], flowTo[i]
		graph[from_] = append(graph[from_], to_)
	}

	sccCount, belong := StronglyConnectedComponent(graph)
	dag := SCCDag(graph, sccCount, belong)

	reachable := make([]bool, sccCount)
	{
		var queue []int
		for _, c := range crystals {
			queue = append(queue, belong[c])
			reachable[belong[c]] = true
		}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range dag[u] {
				if !reachable[v] {
					reachable[v] = true
					queue = append(queue, v)
				}
			}
		}
	}

	indeg := make([]int, sccCount)
	{
		for i := 0; i < sccCount; i++ {
			for _, v := range dag[i] {
				indeg[v]++
			}
		}
	}

	res := 0
	for i := 0; i < sccCount; i++ {
		if !reachable[i] && indeg[i] == 0 {
			res++
		}
	}
	return res
}

// 有向图强连通分量分解.
func StronglyConnectedComponent(graph [][]int) (count int, belong []int) {
	n := len(graph)
	belong = make([]int, n)
	low := make([]int, n)
	order := make([]int, n)
	for i := range order {
		order[i] = -1
	}
	now := 0
	path := []int{}

	var dfs func(int)
	dfs = func(v int) {
		low[v] = now
		order[v] = now
		now++
		path = append(path, v)
		for _, to := range graph[v] {
			if order[to] == -1 {
				dfs(to)
				low[v] = min(low[v], low[to])
			} else {
				low[v] = min(low[v], order[to])
			}
		}
		if low[v] == order[v] {
			for {
				u := path[len(path)-1]
				path = path[:len(path)-1]
				order[u] = n
				belong[u] = count
				if u == v {
					break
				}
			}
			count++
		}
	}

	for i := 0; i < n; i++ {
		if order[i] == -1 {
			dfs(i)
		}
	}
	for i := 0; i < n; i++ {
		belong[i] = count - 1 - belong[i]
	}
	return
}

// 有向图的强连通分量缩点.
func SCCDag(graph [][]int, count int, belong []int) (dag [][]int) {
	dag = make([][]int, count)
	adjSet := make([]map[int]struct{}, count)
	for i := 0; i < count; i++ {
		adjSet[i] = make(map[int]struct{})
	}
	for cur, nexts := range graph {
		for _, next := range nexts {
			if bid1, bid2 := belong[cur], belong[next]; bid1 != bid2 {
				adjSet[bid1][bid2] = struct{}{}
			}
		}
	}
	for i := 0; i < count; i++ {
		for next := range adjSet[i] {
			dag[i] = append(dag[i], next)
		}
	}
	return
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
