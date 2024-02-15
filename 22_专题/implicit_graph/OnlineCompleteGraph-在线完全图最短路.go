// OnlineCompleteGraph-在线完全图bfs

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
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
func ComplementGraphUnionFind(n int, banEdges [][2]int) *UnionFindArray {
	ban := make([][]int, n)
	for _, e := range banEdges {
		u, v := e[0], e[1]
		ban[u] = append(ban[u], v)
		ban[v] = append(ban[v], u)
	}

	uf := NewUnionFindArray(n)
	fs := NewFastSet(n)
	for i := 0; i < n; i++ {
		fs.Insert(i)
	}

	stack := make([]int, n) // dfs
	ptr := 0
	for i := 0; i < n; i++ {
		if !fs.Has(i) {
			continue
		}

		stack[ptr] = i
		ptr++
		for ptr > 0 {
			ptr--
			leader := stack[ptr]

			var removed []int
			for _, to := range ban[leader] {
				if fs.Has(to) {
					removed = append(removed, to)
				}
			}

			for _, v := range removed {
				fs.Erase(v)
			}
			fs.Enumerate(0, n, func(to int) {
				fs.Erase(to)
				stack[ptr] = to
				ptr++
				uf.Union(leader, to)
			})
			for _, v := range removed {
				fs.Insert(v)
			}
		}
	}

	return uf
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int
	n    int
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{
		Part: n,
		n:    n,
		data: data,
	}
}

// 按秩合并.
func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	if cb != nil {
		cb(root1, root2)
	}
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetSize(key int) int {
	return -ufa.data[ufa.Find(key)]
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}

type FastSet struct {
	n, lg int
	seg   [][]int
	size  int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true

}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.size
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}
