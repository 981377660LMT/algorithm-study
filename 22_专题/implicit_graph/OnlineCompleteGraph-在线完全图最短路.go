// OnlineCompleteGraph-在线完全图bfs

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 1e18
const MOD int = 998244353

func main() {
	// demo()
	// CountingShortestPaths()
	// SafetyJourney()
	// demo()
	CF920E()
}

func demo() {
	fmt.Println(ComplementGraphBfs(5, 0, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}))
	fmt.Println(ComplementGraphDistance(4, [][2]int{{0, 1}, {1, 2}, {2, 3}}))
	uf := ComplementGraphUnionFind(4, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}})
	fmt.Println(uf)
}

// G - Counting Shortest Paths (完全图dp、完全图最短路)
// https://atcoder.jp/contests/abc319/tasks/abc319_g
// 在一个n个顶点的无向无权完全图中删除m条边,求从START到TARGET的最路径数模998244353.
//
// bfs最短路计数分解成两个问题:
//  1. 在线bfs求出START到其他点的最短路;
//  2. START到TARGET的路径上所有的边 u->v 满足 dist[u]+1 == dist[v].(最短路的充要条件)
//     !可以按照距离从小到大dp.用总数减去不合法的路径数即可.
//
// 参考:
// E - Safety Journey
// https://atcoder.jp/contests/abc212/tasks/abc212_e
// 转移边数很多的dp问题 => 正难则反.
func CountingShortestPaths() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	banEdges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &banEdges[i][0], &banEdges[i][1])
		banEdges[i][0]--
		banEdges[i][1]--
	}
	START, TARGET := 0, n-1

	dist, _ := ComplementGraphBfs(n, START, banEdges)
	if dist[TARGET] == INF {
		fmt.Fprintln(out, -1)
		return
	}

	ban := make([][]int, n) // 每个点的禁止转移的邻居
	for _, e := range banEdges {
		u, v := e[0], e[1]
		ban[u] = append(ban[u], v)
		ban[v] = append(ban[v], u)
	}

	groupByDist := make([][]int, n) // 距离START为0,1,2,...,n-1的点
	for i := 0; i < n; i++ {
		if d := dist[i]; d != INF {
			groupByDist[d] = append(groupByDist[d], i)
		}
	}

	dp := make([]int, n) // 到达i点的路径数
	dp[START] = 1
	for d := 0; d < n-1; d++ {
		preCount := 0
		for _, pre := range groupByDist[d] {
			preCount += dp[pre]
			preCount %= MOD
		}
		for _, cur := range groupByDist[d+1] {
			dp[cur] = preCount
			for _, pre := range ban[cur] {
				if dist[pre]+1 == dist[cur] {
					dp[cur] -= dp[pre]
					dp[cur] %= MOD
					if dp[cur] < 0 {
						dp[cur] += MOD
					}
				}
			}
		}
	}

	fmt.Fprintln(out, dp[TARGET])
}

// E - Safety Journey (完全图dp)
// https://atcoder.jp/contests/abc212/tasks/abc212_e
// 一张n个点的完全图,删去m条边,一共走k步,求从START到TARGET的方案数.
// n,k,m<=5000
func SafetyJourney() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	banEdegs := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &banEdegs[i][0], &banEdegs[i][1])
		banEdegs[i][0]--
		banEdegs[i][1]--
	}

	modAdd := func(a, b int) int {
		res := (a + b) % MOD
		if res < 0 {
			res += MOD
		}
		return res
	}

	START, TARGET := 0, 0

	ban := make([][]int, n) // 每个点的禁止转移的邻居
	for _, e := range banEdegs {
		u, v := e[0], e[1]
		ban[u] = append(ban[u], v)
		ban[v] = append(ban[v], u)
	}

	dp := make([]int, n)
	dp[START] = 1
	for step := 0; step < k; step++ {
		preCount := 0
		for _, pre := range dp {
			preCount = modAdd(preCount, pre)
		}
		ndp := make([]int, n)
		for i := 0; i < n; i++ {
			ndp[i] = preCount
			ndp[i] = modAdd(ndp[i], -dp[i]) // !不能从同一点转移过来
			for _, pre := range ban[i] {    // !不能从不通的路转移过来
				ndp[i] = modAdd(ndp[i], -dp[pre])
			}
		}

		dp = ndp
	}

	fmt.Fprintln(out, dp[TARGET])
}

// Connected Components?
// https://www.luogu.com.cn/problem/CF920E
// 给定一张n个点的完全图,删去m条边(剩余 (n*(n-1)/2 - m) 条边),求剩下的连通块数,以及每个连通块的大小.
func CF920E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	banEdges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &banEdges[i][0], &banEdges[i][1])
		banEdges[i][0]--
		banEdges[i][1]--
	}

	uf := ComplementGraphUnionFind(n, banEdges)
	groups := uf.GetGroups()
	sizes := make([]int, 0, len(groups))
	for _, v := range groups {
		sizes = append(sizes, len(v))
	}
	sort.Ints(sizes)
	fmt.Fprintln(out, len(sizes))
	for _, v := range sizes {
		fmt.Fprint(out, v, " ")
	}
}

// 完全图最短路.
//
//	给定一个无向无权的完全图，求出完全图上从start到其他点的最短路.不可达的点距离为INF.
//	banEdges是禁止通行的边.
//	O(V+len(banEdges)).
func ComplementGraphBfs(n int, start int, banEdges [][2]int) (dist []int, pre []int) {
	banGraph := make([][]int, n)
	for _, e := range banEdges {
		u, v := e[0], e[1]
		banGraph[u] = append(banGraph[u], v)
		banGraph[v] = append(banGraph[v], u)
	}

	dist = make([]int, n)
	pre = make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		pre[i] = -1
	}
	dist[start] = 0
	queue := make([]int, 0, n)
	queue = append(queue, start)

	notNeightBor := make([]bool, n)
	unVisited := make([]int, 0, n-1)
	for i := 0; i < n; i++ {
		if i != start {
			unVisited = append(unVisited, i)
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, u := range banGraph[cur] {
			notNeightBor[u] = true
		}

		nextUnvisited := []int{}
		for _, next := range unVisited {
			if notNeightBor[next] {
				nextUnvisited = append(nextUnvisited, next) // findUnvisited
			} else {
				// setVisited
				dist[next] = dist[cur] + 1
				pre[next] = cur
				queue = append(queue, next)
			}
		}
		unVisited = nextUnvisited

		for _, u := range banGraph[cur] {
			notNeightBor[u] = false
		}
	}

	return
}

// 完全图距离>=2的点对.
// 返回: [u,v,dist(u,v)].
func ComplementGraphDistance(n int, banEdges [][2]int) (res [][3]int) {
	deg := make([]int, n)
	ban := make([][]int, n)
	for _, e := range banEdges {
		u, v := e[0], e[1]
		ban[u] = append(ban[u], v)
		ban[v] = append(ban[v], u)
		deg[u]++
		deg[v]++
	}

	minArg := 0 // 度数最小的点
	for i := 0; i < n; i++ {
		if deg[i] < deg[minArg] {
			minArg = i
		}
	}

	removed := make([]bool, n)
	for _, v := range ban[minArg] {
		removed[v] = true
	}

	for _, e := range banEdges {
		u, v := e[0], e[1]
		if removed[u] || removed[v] {
			continue
		}
		res = append(res, [3]int{u, v, 2}) // u -> minArg -> v
	}

	for _, v := range ban[minArg] {
		dist, _ := ComplementGraphBfs(n, v, banEdges)
		for i := 0; i < n; i++ {
			if dist[i] <= 1 {
				continue
			}
			if removed[i] && v >= i {
				continue
			}
			res = append(res, [3]int{v, i, dist[i]})
		}
	}

	return
}

// 完全图并查集.
// 维护一个 set，保存当前未访问过的点。每一次dfs从未访问过的点出发，遍历到一个节点后删除对应元素。
func ComplementGraphUnionFind(n int, banEdges [][2]int) *UnionFindArraySimple {
	ban := make([][]int, n)
	for _, e := range banEdges {
		u, v := e[0], e[1]
		ban[u] = append(ban[u], v)
		ban[v] = append(ban[v], u)
	}
	queue, toVisit := make([]int, n), make([]int, n)
	bad := make([]bool, n)
	ptr := 0
	for i := 0; i < n; i++ {
		toVisit[ptr] = i
		ptr++
	}
	uf := NewUnionFindArraySimple(n)
	for ptr > 0 {
		head, tail := 0, 0
		queue[tail] = toVisit[ptr-1]
		ptr--
		tail++
		for head < tail {
			cur := queue[head]
			head++
			for _, next := range ban[cur] {
				bad[next] = true
			}
			for i := ptr - 1; i >= 0; i-- {
				next := toVisit[i]
				if bad[next] {
					continue
				}
				queue[tail] = next
				tail++
				toVisit[i], toVisit[ptr-1] = toVisit[ptr-1], toVisit[i]
				ptr--
				uf.Union(cur, next)
			}
			for _, next := range ban[cur] {
				bad[next] = false
			}
		}
	}
	return uf
}

type UnionFindArraySimple struct {
	Part int
	n    int
	data []int32
}

func NewUnionFindArraySimple(n int) *UnionFindArraySimple {
	data := make([]int32, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple) Union(key1 int, key2 int) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

func (u *UnionFindArraySimple) Find(key int) int {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = int32(u.Find(int(u.data[key])))
	return int(u.data[key])
}

func (u *UnionFindArraySimple) GetSize(key int) int {
	return int(-u.data[u.Find(key)])
}

func (u *UnionFindArraySimple) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for k := range u.data {
		root := u.Find(k)
		groups[root] = append(groups[root], k)
	}
	return groups
}
