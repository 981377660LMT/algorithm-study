// MaxMatching
// MinVertexCover
// MaxIndependentSet
// MinEdgeCover
// DMDecomposition
// !MaxAntiChain
//
// !注意是无向图!!!

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// Movie()
	// MerryChristmas()
	// demo()
	// abc237ex()
	abc274g()
	// Yuki1479()
	// Yuki1744()
	// Yuki1745()

}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=1566
// 太郎君需要在夏假期间`每天`都看一部电影，每部电影上映时间为[a,b](1<=a<=b<=31).
// 如果是第一次看这部电影，他会得到100的幸福度，
// 如果是重复看，他会得到50的幸福度。他的目标是最大化他的总幸福度。
//
// n<=100
//
// 将每一天和每一部电影看作图的一个节点，如果某一天可以观看某部电影（即电影的上映日期包含这一天），
// 那么就在这两个节点之间连一条边。这样，我们就得到了一个二分图，
// 其中一部分节点代表日期，另一部分节点代表电影。
// 在这个二分图中，找到一个最大匹配，就相当于找到一个方式，
// 使得太郎君在每一天都看一部电影，并且尽可能地使得看的电影是第一次看，从而使得他得到的幸福度最大。
func Movie() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	OFFSET := 31
	intervals := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &intervals[i][0], &intervals[i][1])
		intervals[i][0]--
	}
	graph := make([][]int, OFFSET+n)
	isIn := make([]bool, OFFSET)
	for i := 0; i < n; i++ {
		start, end := intervals[i][0], intervals[i][1]
		for v := start; v < end; v++ {
			graph[v] = append(graph[v], OFFSET+i)
			graph[OFFSET+i] = append(graph[OFFSET+i], v)
			isIn[v] = true
		}
	}
	bm := NewBipartiteMatching(graph, nil)
	a := len(bm.MaxMatching())
	b := 0
	for _, v := range isIn {
		if v {
			b++
		}
	}
	fmt.Fprintln(out, (a+b)*50)
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2251
// 给定一个无向带权图(可能不连通)和q个查询.
// 每个查询形如(pos,time)，表示有一个在时间time向房屋pos送货的请求.
// 需要在图中放置人员送货，如何安排最少数量的人完成所有请求(每个人单位移动速度为1)。
// 同一地点和时间最多只有一个请求。
// n<=100 m<=1000 q<=1000
//
// !dag最小可相交路径覆盖=dag最长反链
func MerryChristmas() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(n, m, q int) {
		adjList := make([][][2]int, n)
		for i := 0; i < m; i++ {
			var u, v, w int
			fmt.Fscan(in, &u, &v, &w)
			adjList[u] = append(adjList[u], [2]int{v, w})
			adjList[v] = append(adjList[v], [2]int{u, w})
		}
		dist := make([][]int, n)
		for i := range dist {
			d, _ := Dijkstra(n, adjList, i)
			dist[i] = d
		}

		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			var pos, time int
			fmt.Fscan(in, &pos, &time)
			if pos == 0 && time == 0 {
				continue
			}
			queries[i] = [2]int{pos, time}
		}

		dag := make([][]int, q)
		// 连续解决任务
		for i := 0; i < q; i++ {
			for j := 0; j < q; j++ {
				if i == j {
					continue
				}
				pos1, time1 := queries[i][0], queries[i][1]
				pos2, time2 := queries[j][0], queries[j][1]
				if time1+dist[pos1][pos2] <= time2 {
					dag[i] = append(dag[i], j)
				}
			}
		}

		res := MaxAntiChain(q, dag)
		fmt.Fprintln(out, len(res))
	}

	for {
		var n, m, q int
		fmt.Fscan(in, &n, &m, &q)
		if n+m+q == 0 {
			break
		}
		solve(n, m, q)
	}
}

// https://atcoder.jp/contests/abc237/tasks/abc237_h
// 给定一个字符串, 你需要从中选出若干回文子串, 并且使得选出的串不存在某一个是另一个的子串, 问最多能选出多少子串.
// n<=200
//
// 给定一些偏序包含关系,求最大独立集(互相无法到达).
// !遍历所有回文子串，如果j是i的子串，则连边 i->j，求dag最长反链即可.
func abc237ex() {
	zAlgo := func(s string) []int {
		n := len(s)
		if n == 0 {
			return nil
		}
		z := make([]int, n)
		j := 0
		for i := 1; i < n; i++ {
			var k int
			if j+z[j] <= i {
				k = 0
			} else {
				k = min(j+z[j]-i, z[i-j])
			}
			for i+k < n && s[k] == s[i+k] {
				k++
			}
			if j+z[j] < i+z[i] {
				j = i
			}
			z[i] = k
		}
		z[0] = n
		return z
	}

	// O(n+m)判断`shorter`是否是`longer`的子串.
	isSubstring := func(longer, shorter string) bool {
		if len(shorter) > len(longer) {
			return false
		}
		if len(shorter) == 0 {
			return true
		}
		n, m := len(longer), len(shorter)
		z := zAlgo(shorter + longer)
		for i := m; i < n+m; i++ {
			if z[i] >= m {
				return true
			}
		}
		return false
	}

	isPalindrome := func(s string) bool {
		n := len(s)
		for i := 0; i < n>>1; i++ {
			if s[i] != s[n-1-i] {
				return false
			}
		}
		return true
	}

	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	set := make(map[string]struct{})
	n := len(s)
	for start := 0; start < n; start++ {
		for end := start + 1; end <= n; end++ {
			cur := s[start:end]
			if isPalindrome(cur) {
				set[cur] = struct{}{}
			}
		}
	}
	allPalindromes := make([]string, 0, len(set))
	for k := range set {
		allPalindromes = append(allPalindromes, k)
	}

	dag := make([][]int, len(allPalindromes))
	for i := range dag {
		for j := range dag {
			if i == j {
				continue
			}
			if isSubstring(allPalindromes[i], allPalindromes[j]) {
				dag[i] = append(dag[i], j)
			}
		}
	}

	res := MaxAntiChain(len(dag), dag)
	fmt.Fprintln(out, len(res))
}

// G - Security Camera 3
// https://atcoder.jp/contests/abc274/tasks/abc274_g
// 网格图中有一些障碍,现在在这个图中放置一些监视器监视.
// 每个监视器可以监视一条直线(不被阻挡的情况下).
// 覆盖所有格子最少需要放置多少个监视器?
// ROW,COL<=300
// 障碍:"#"，无障碍:"."
func abc274g() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)
	grid := make([]string, ROW)
	for i := range grid {
		var row string
		fmt.Fscan(in, &row)
		grid[i] = row
	}

	OFFSET := ROW * COL
	graph := make([][]int, ROW*COL*2)
	for x := 0; x < ROW; x++ {
		for y := 0; y < COL; y++ {
			if grid[x][y] == '#' {
				continue
			}
			curX, curY := x, y
			for curX > 0 && grid[curX-1][y] == '.' {
				curX--
			}
			for curY > 0 && grid[x][curY-1] == '.' {
				curY--
			}
			// 要么取(curX, y)，要么取(x, curY)
			a := curX*COL + y
			b := x*COL + curY
			graph[a] = append(graph[a], b+OFFSET)
			graph[b+OFFSET] = append(graph[b+OFFSET], a)
		}
	}

	bm := NewBipartiteMatching(graph, nil)
	res := bm.MinVertexCover()
	fmt.Fprintln(out, len(res))
}

// https://yukicoder.me/problems/no/1479
// 给定一个矩阵, 你可以按照任意顺序执行以下操作任意多次：
// 选择矩阵的行。此时，将行中所有等于行元素最大值的元素全部变为0。
// 选择矩阵的列。此时，将列中所有等于列元素最大值的元素全部变为0。
// 你的目标是将矩阵的所有元素变为0。请求出达成目标所需的最小操作次数。
// ROW,COL<=500
// A[i][j]<=5e5
func Yuki1479() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	unique := func(a *[]int) {
		set := make(map[int]struct{})
		for _, v := range *a {
			set[v] = struct{}{}
		}
		allNums := make([]int, 0, len(set))
		for k := range set {
			allNums = append(allNums, k)
		}
		sort.Ints(allNums)
		*a = allNums
	}

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)

	mp := make(map[int][][2]int)
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			var v int
			fmt.Fscan(in, &v)
			mp[v] = append(mp[v], [2]int{i, j})
		}
	}

	res := 0
	for k, points := range mp {
		if k == 0 {
			continue
		}
		xs, ys := make([]int, len(points)), make([]int, len(points))
		for i, p := range points {
			xs[i], ys[i] = p[0], p[1]
		}
		unique(&xs)
		unique(&ys)
		g := make([][]int, len(xs)+len(ys))
		for _, p := range points {
			x, y := p[0], p[1]
			x = sort.SearchInts(xs, x)
			y = sort.SearchInts(ys, y)
			g[x] = append(g[x], len(xs)+y)
			g[len(xs)+y] = append(g[len(xs)+y], x)
		}
		bm := NewBipartiteMatching(g, nil)
		res += len(bm.MinVertexCover())
	}

	fmt.Fprintln(out, res)
}

// https://yukicoder.me/problems/no/1744
// 间谍分配任务。
// 每个任务最多只能分配给一个间谍，每个间谍最多只能分配一个任务(匹配)。
// !给出q个查询，每个查询形如(u,v)，问能否在不选择间谍u和v的情况下，达成最大匹配。
// n,m,q<=1e5
// DM分解.O(n+m+q*sqrt(n+m))
func Yuki1744() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var left, right, p int
	fmt.Fscan(in, &left, &right, &p)
	graph := make([][]int, left+right)
	edges := make([][2]int, p)
	for i := 0; i < p; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = [2]int{u, v + left}
		graph[u] = append(graph[u], v+left)
		graph[v+left] = append(graph[v+left], u)
	}

	bm := NewBipartiteMatching(graph, nil)
	compCount, belong := bm.DMDecomposition(edges)
	sameCounter := make([]int, compCount+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		if belong[u] == belong[v] {
			sameCounter[belong[u]]++
		}
	}
	for _, e := range edges {
		// 一定在最大匹配中，删除后最大匹配减少1.
		if belong[e[0]] == belong[e[1]] && sameCounter[belong[e[0]]] == 1 {
			fmt.Fprintln(out, "No")
		} else {
			fmt.Fprintln(out, "Yes")
		}
	}
}

// https://yukicoder.me/problems/no/1745
// 间谍分配任务。
// 每个任务最多只能分配给一个间谍，每个间谍最多只能分配一个任务(匹配)。
// !给出q个查询，每个查询形如(u,v)，问选择间谍u和v的情况下能否达成最大匹配。
// n,m,q<=1e5
// DM分解.O(n+m+q*sqrt(n+m))
func Yuki1745() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var left, right, p int
	fmt.Fscan(in, &left, &right, &p)
	graph := make([][]int, left+right)
	edges := make([][2]int, p)
	for i := 0; i < p; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = [2]int{u, v + left}
		graph[u] = append(graph[u], v+left)
		graph[v+left] = append(graph[v+left], u)
	}

	bm := NewBipartiteMatching(graph, nil)
	_, belong := bm.DMDecomposition(edges)

	for _, e := range edges {
		// 可能在最大匹配中
		if belong[e[0]] == belong[e[1]] {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

const INF int = 1e18

// `无向`二分图最大匹配.
type BipartiteMatching struct {
	n       int32
	graph   [][]int
	color   []int8
	dist    []int32
	match   []int32
	visited []bool
}

// `无向`二分图最大匹配.
// colors: 如果为nil，内部会自动计算.
func NewBipartiteMatching(graph [][]int, colors []int8) *BipartiteMatching {
	n := len(graph)
	bm := &BipartiteMatching{
		n:       int32(n),
		graph:   graph,
		dist:    make([]int32, n),
		match:   make([]int32, n),
		visited: make([]bool, n),
	}
	if colors != nil {
		bm.color = colors
	} else {
		bm.color = BipartiteVertexColoring(n, graph)
	}
	if n > 0 && len(bm.color) == 0 {
		panic("not bipartite graph")
	}
	for i := range bm.dist {
		bm.dist[i] = -1
	}
	for i := range bm.match {
		bm.match[i] = -1
	}
	for {
		bm.bfs()
		for i := range bm.visited {
			bm.visited[i] = false
		}
		flow := 0
		for v := int32(0); v < bm.n; v++ {
			if bm.color[v] == 0 && bm.match[v] == -1 && bm.dfs(v) {
				flow++
			}
		}
		if flow == 0 {
			break
		}
	}
	return bm
}

// 最大匹配.
func (bm *BipartiteMatching) MaxMatching() (res [][2]int) {
	for v := int32(0); v < bm.n; v++ {
		if v < bm.match[v] {
			res = append(res, [2]int{int(v), int(bm.match[v])})
		}
	}
	return
}

// 最小点覆盖.
// 最小点覆盖是指在二部图中，可以覆盖所有边的最小顶点集合。
func (bm *BipartiteMatching) MinVertexCover() (res []int) {
	for v := int32(0); v < bm.n; v++ {
		if (bm.color[v] != 0) != (bm.dist[v] == -1) {
			res = append(res, int(v))
		}
	}
	return
}

// 最大独立集.
// 没有相邻顶点的最大顶点集合。
func (bm *BipartiteMatching) MaxIndependentSet() (res []int) {
	for v := int32(0); v < bm.n; v++ {
		if (bm.color[v] != 0) == (bm.dist[v] == -1) {
			res = append(res, int(v))
		}
	}
	return
}

// 最小边覆盖.返回排序后的边的编号.
// 可以覆盖所有顶点的最小边集合。
func (bm *BipartiteMatching) MinEdgeCover(edges [][2]int) (res []int) {
	done := make([]bool, bm.n)
	for ei, e := range edges {
		u, v := e[0], e[1]
		if done[u] || done[v] {
			continue
		}
		if bm.match[u] == int32(v) {
			res = append(res, ei)
			done[u] = true
			done[v] = true
		}
	}
	for ei, e := range edges {
		u, v := e[0], e[1]
		if !done[u] {
			res = append(res, ei)
			done[u] = true
		}
		if !done[v] {
			res = append(res, ei)
			done[v] = true
		}
	}
	sort.Ints(res)
	return
}

func (bm *BipartiteMatching) Debug() {
	fmt.Println("match", bm.match)
	fmt.Println("MinVertexCoverr", bm.MinVertexCover())
	fmt.Println("MaxIndependentSet", bm.MaxIndependentSet())
}

// Dulmage–Mendelsohn decomposition （DM分解）
// https://ei1333.github.io/library/graph/flow/bipartite-flow.hpp
// https://hitonanode.github.io/cplib-cpp/graph/dulmage_mendelsohn_decomposition.hpp.html
// 在残量网络的基础上添加虚拟源点和汇点，进行强连通分量分解.
// 性质:
//  1. 边是否存在与最大匹配中：
//     如果一条边的两个顶点所在的强连通分量相同，那么可以作为最大匹配的一条边(可能在最大匹配中)。
//     如果一条边的两个顶点所在的强连通分量相同并且没有其他边具有相同的id，那么它一定在任何最大匹配中。
//     如果一条边的两个顶点所在的强连通分量不同，那么不可能作为最大匹配的一条边。
//     如果一条边在某个最大匹配中，那么它一定在任何最大匹配中。
//  2. 点是否存在于最大匹配中：
//     如果一个顶点的颜色为 0，并且它的 id 值在 1 到 compCount 的闭区间范围内，那么这个顶点一定会被用于最大匹配。
//     如果一个顶点的颜色为 1，并且它的 id 值在 0 到 compCount-1 的闭区间范围内，那么这个顶点一定会被用于最大匹配。
//  3. 如果一条边从 color=0 的顶点到 color=1 的顶点，那么这条边的左顶点的 id 值应该小于或等于右顶点的 id 值。
func (bm *BipartiteMatching) DMDecomposition(edges [][2]int) (compCount int, belong []int) {
	belong = make([]int, bm.n)
	for i := range belong {
		belong[i] = -1
	}
	queue := []int{}
	add := func(v, x int) {
		if belong[v] == -1 {
			belong[v] = x
			queue = append(queue, v)
		}
	}
	for v := 0; v < int(bm.n); v++ {
		if bm.match[v] == -1 && bm.color[v] == 0 {
			add(v, 0)
		}
	}
	for v := 0; v < int(bm.n); v++ {
		if bm.match[v] == -1 && bm.color[v] == 1 {
			add(v, INF)
		}
	}
	for len(queue) > 0 {
		v := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if bm.match[v] != -1 {
			add(int(bm.match[v]), belong[v])
		}
		if bm.color[v] == 0 && belong[v] == 0 {
			for _, to := range bm.graph[v] {
				add(to, belong[v])
			}
		}
		if bm.color[v] == 1 && belong[v] == INF {
			for _, to := range bm.graph[v] {
				add(to, belong[v])
			}
		}
	}

	// 残留图的强连通分量分解.
	vs := []int{}
	for v := 0; v < int(bm.n); v++ {
		if belong[v] == -1 {
			vs = append(vs, v)
		}
	}
	m := len(vs)
	dg := make([][]int, m)
	for i := range dg {
		v := vs[i]
		if bm.match[v] != -1 {
			j := sort.SearchInts(vs, int(bm.match[v]))
			dg[i] = append(dg[i], j)
		}
		if bm.color[v] == 0 {
			for _, to := range bm.graph[v] {
				if belong[to] != -1 || to == int(bm.match[v]) {
					continue
				}
				j := sort.SearchInts(vs, to)
				dg[i] = append(dg[i], j)
			}
		}
	}

	compCount, comp := StronglyConnectedComponent(dg)
	compCount++

	for i := 0; i < m; i++ {
		belong[vs[i]] = 1 + comp[i]
	}
	for v := 0; v < int(bm.n); v++ {
		if belong[v] == INF {
			belong[v] = compCount
		}
	}
	return
}

func (bm *BipartiteMatching) bfs() {
	for i := range bm.dist {
		bm.dist[i] = -1
	}
	queue := []int32{}
	for v := int32(0); v < bm.n; v++ {
		if bm.color[v] == 0 && bm.match[v] == -1 {
			queue = append(queue, v)
			bm.dist[v] = 0
		}
	}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range bm.graph[v] {
			bm.dist[to] = 0
			w := bm.match[to]
			if w != -1 && bm.dist[w] == -1 {
				bm.dist[w] = bm.dist[v] + 1
				queue = append(queue, w)
			}
		}
	}
}

func (bm *BipartiteMatching) dfs(v int32) bool {
	bm.visited[v] = true
	for _, to := range bm.graph[v] {
		w := bm.match[to]
		if w == -1 || (!bm.visited[w] && bm.dist[w] == bm.dist[v]+1 && bm.dfs(w)) {
			bm.match[to] = v
			bm.match[v] = int32(to)
			return true
		}
	}
	return false
}

// 无向图二分图着色.
// 如果不是二分图，返回空数组.
func BipartiteVertexColoring(n int, graph [][]int) (colors []int8) {
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

// 有向图强连通分量.
func StronglyConnectedComponent(graph [][]int) (compCount int, belong []int) {
	n32 := int32(len(graph))
	compId := int32(0)
	comp := make([]int32, n32)
	low := make([]int32, n32)
	ord := make([]int32, n32)
	for i := range ord {
		ord[i] = -1
	}
	path := []int32{}
	now := int32(0)

	var dfs func(int32)
	dfs = func(v int32) {
		low[v] = now
		ord[v] = now
		now++
		path = append(path, v)
		for _, to := range graph[v] {
			if ord[to] == -1 {
				dfs(int32(to))
				if low[v] > low[to] {
					low[v] = low[to]
				}
			} else if low[v] > ord[to] {
				low[v] = ord[to]
			}
		}
		if low[v] == ord[v] {
			for {
				u := path[len(path)-1]
				path = path[:len(path)-1]
				ord[u] = n32
				comp[u] = compId
				if u == v {
					break
				}
			}
			compId++
		}
	}

	for v := int32(0); v < n32; v++ {
		if ord[v] == -1 {
			dfs(v)
		}
	}

	compCount = int(compId)
	belong = make([]int, n32)
	for v := int32(0); v < n32; v++ {
		belong[v] = compCount - 1 - int(comp[v])
	}
	return
}

func SccDag(graph [][]int, compCount int, belong []int) (dag [][]int) {
	unique := func(nums []int32) []int32 {
		set := make(map[int32]struct{})
		for _, v := range nums {
			set[v] = struct{}{}
		}
		res := make([]int32, 0, len(set))
		for k := range set {
			res = append(res, k)
		}
		return res
	}

	edges := make([][]int32, compCount)
	for cur, nexts := range graph {
		curComp := belong[cur]
		for _, next := range nexts {
			nextComp := belong[next]
			if curComp != nextComp {
				edges[curComp] = append(edges[curComp], int32(nextComp))
			}
		}
	}

	dag = make([][]int, compCount)
	for cur := 0; cur < compCount; cur++ {
		edges[cur] = unique(edges[cur])
		for _, next := range edges[cur] {
			dag[cur] = append(dag[cur], int(next))
		}
	}

	return
}

// dag最长反链(最大独立集).
func MaxAntiChain(n int, dag [][]int) []int {
	newGraph := make([][]int, n+n)
	for i := 0; i < n; i++ {
		for _, to := range dag[i] {
			newGraph[i] = append(newGraph[i], to+n)
		}
	}
	bm := NewBipartiteMatching(newGraph, nil)
	cover := bm.MinVertexCover()
	ok := make([]bool, n)
	for i := range ok {
		ok[i] = true
	}
	for _, v := range cover {
		ok[v%n] = false
	}
	antichain := []int{}
	for v := 0; v < n; v++ {
		if ok[v] {
			antichain = append(antichain, v)
		}
	}
	return antichain
}

type Edge = [2]int

func Dijkstra(n int, adjList [][]Edge, start int) (dist, preV []int) {
	type pqItem struct{ node, dist int }
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	preV = make([]int, n)
	for i := range preV {
		preV[i] = -1
	}

	pq := nhp(func(a, b H) int {
		return a.(pqItem).dist - b.(pqItem).dist
	}, nil)
	pq.Push(pqItem{start, 0})

	for pq.Len() > 0 {
		curNode := pq.Pop().(pqItem)
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge[0], edge[1]
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				preV[next] = cur
				pq.Push(pqItem{next, cand})
			}
		}
	}

	return
}

type H = interface{}

// Should return a number:
//
//	negative , if a < b
//	zero     , if a == b
//	positive , if a > b
type Comparator func(a, b H) int

func nhp(comparator Comparator, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{comparator: comparator, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data       []H
	comparator Comparator
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		return
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Peek() (value H) {
	if h.Len() == 0 {
		return
	}
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.comparator(h.data[root], h.data[parent]) < 0; parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.comparator(h.data[left], h.data[minIndex]) < 0 {
			minIndex = left
		}

		if right < n && h.comparator(h.data[right], h.data[minIndex]) < 0 {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
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
