// https://yukicoder.me/problems/no/114
// !一个带权的无向图上有k个关键点，求联通k个关键点最小的代价(边权之和)。
// n<=35,k<=n,m<=n*(n-1)/2

// k<=15时最小斯坦纳树
// k>15时枚举+并查集

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges = append(edges, []int{u, v, w})
	}
	criticals := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &criticals[i])
		criticals[i]--
	}

	fmt.Fprintln(out, solve(n, edges, criticals))
}

func solve(n int, edges [][]int, criticals []int) int {
	if len(criticals) <= 15 {
		return MinimumSteinerTree(n, edges, criticals)
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i][2] < edges[j][2]
	})
	notCriticals := []int{}
	set := map[int]struct{}{}
	for _, c := range criticals {
		set[c] = struct{}{}
	}
	for i := 0; i < n; i++ {
		if _, ok := set[i]; !ok {
			notCriticals = append(notCriticals, i)
		}
	}

	// 选择点集state时,连通criticals的最小边权之和
	cal := func(state int) int {
		uf := NewUnionFindArray(n)
		ok := 0
		cost := 0
		for _, e := range edges {
			u, v, w := e[0], e[1], e[2]
			if state&(1<<u) != 0 && state&(1<<v) != 0 && uf.Union(u, v) {
				cost += w
				ok++
			}
		}
		if ok == bits.OnesCount(uint(state))-1 {
			return cost
		}
		return INF
	}

	res := INF
	state := 0
	for _, c := range criticals {
		state |= 1 << c
	}
	for i := 0; i < (1 << len(notCriticals)); i++ {
		curState := state
		for j := 0; j < len(notCriticals); j++ {
			if i&(1<<j) != 0 {
				curState |= 1 << notCriticals[j]
			}
		}
		res = min(res, cal(curState))
	}

	return res
}

const INF int = 1e18

// 一个联通的无向带权图上有k个关键点 criticals，求联通k个关键点最小的代价(边权之和)。
func MinimumSteinerTree(n int, edges [][]int, criticals []int) int {
	k := len(criticals)
	visited := make([]bool, n)
	graph := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		graph[u] = append(graph[u], [2]int{v, w})
		graph[v] = append(graph[v], [2]int{u, w})
	}

	dp := make([][]int, 1<<k)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}

	for i := 0; i < k; i++ {
		dp[1<<i][criticals[i]] = 0
	}

	spfa := func(s int) {
		q := []int{}
		for i := 0; i < n; i++ {
			if dp[s][i] != INF {
				q = append(q, i)
				visited[i] = true
			}
		}

		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			visited[u] = false
			for _, e := range graph[u] {
				v, w := e[0], e[1]
				if dp[s][u]+w < dp[s][v] {
					dp[s][v] = dp[s][u] + w
					if !visited[v] {
						q = append(q, v)
						visited[v] = true
					}
				}
			}
		}
	}

	for s := 1; s < (1 << k); s++ {
		for t := s & (s - 1); t > 0; t = (t - 1) & s {
			if t < (s ^ t) {
				break
			}
			for i := 0; i < n; i++ {
				dp[s][i] = min(dp[s][i], dp[t][i]+dp[t^s][i])
			}
		}
		spfa(s)
	}

	res := INF
	for i := 0; i < n; i++ {
		res = min(res, dp[(1<<k)-1][i])
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// NewUnionFindWithCallback ...
func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type _UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	cb(root2, root1)
	return true
}

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *_UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
