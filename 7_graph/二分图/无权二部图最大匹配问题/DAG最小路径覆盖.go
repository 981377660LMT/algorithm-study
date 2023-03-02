// https://ei1333.github.io/library/graph/flow/bipartite-flow.hpp
// https://www.luogu.com.cn/problem/P2764
// https://zhuanlan.zhihu.com/p/125759333
// https://www.cnblogs.com/justPassBy/p/5369930.html

// DAG最小路径覆盖(DAG の最小パス被覆)
// DAG最小路径覆盖可以归结为二分图最大匹配问题

// 覆盖DAG所有点的`路径的集合`叫做DAG的路径覆盖,注意路径不相交
// 路径覆盖的路径数最少的集合叫做DAG的最小路径覆盖

// 做法:
// !原图每个点拆成出点和入点,如果有一条有向边u->v,那么就连一条边u.int->v.out
// 得到一个二分图.
// 跑一遍最大流，便能得到最大合并路径数，再用点数去减即得最小路径覆盖数。
// 从in点到out'点的每一条流，都代表着一次合并。而从源点只给每个点输送1单位流量，又保证了每个点只被经过一次。
// !最小路径覆盖=原图的结点数-新图的最大匹配数。
// n<=150 m<=6000

package main

import "fmt"

// func main() {
// 	// https://www.luogu.com.cn/problem/P2764
// 	in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	var n, m int
// 	fmt.Fscan(in, &n, &m)
// 	edges := make([][]int, m)
// 	for i := 0; i < m; i++ {
// 		edges[i] = make([]int, 2)
// 		fmt.Fscan(in, &edges[i][0], &edges[i][1])
// 		edges[i][0]--
// 		edges[i][1]--
// 	}

// 	count, paths := MinimumPathCovering(n, edges)
// 	for _, path := range paths {
// 		for _, v := range path {
// 			fmt.Fprint(out, v+1, " ")
// 		}
// 		fmt.Fprintln(out)
// 	}
// 	fmt.Fprintln(out, count)
// }

func main() {
	n := 6

	edges := [][]int{
		{0, 1},
		{0, 2},
		{1, 3},
		{2, 3},
		{3, 4},
		{3, 5}}

	fmt.Println(MinimumPathCovering(n, edges))  // 3 [[0 1 3 4] [2] [5]]
	fmt.Println(MinimumPathCovering2(n, edges)) // 2
}

// DAG最小不相交路径覆盖
func MinimumPathCovering(n int, edges [][]int) (count int, paths [][]int) {
	newEdges := make([][]int, 0, len(edges))
	bf := NewBipartiteFlow(n, n)
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		newEdges = append(newEdges, []int{u, v + n}) // 拆点 , A'in => B'out
		bf.AddEdge(u, v)
	}

	maxMathing := bf.MaxMatching()
	count = n - len(maxMathing)

	visited := make([]bool, n)
	ml := bf.MatchL
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		path := []int{i}
		for ml[path[len(path)-1]] > 0 {
			u := path[len(path)-1]
			v := ml[u]
			path = append(path, v)
			visited[v] = true
		}
		paths = append(paths, path)
	}

	return
}

// DAG最小可相交路径覆盖
// 先用floyd求出原图的传递闭包，即如果a到b有路径，那么就加边a->b。然后就转化成了最小不相交路径覆盖问题。
func MinimumPathCovering2(n int, edges [][]int) int {
	adjMatrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = make([]bool, n)
	}
	for _, edge := range edges {
		adjMatrix[edge[0]][edge[1]] = true
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if adjMatrix[i][k] && adjMatrix[k][j] {
					adjMatrix[i][j] = true
				}
			}
		}
	}
	newEdges := make([][]int, 0, len(edges))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if adjMatrix[i][j] {
				newEdges = append(newEdges, []int{i, j})
			}
		}
	}

	res, _ := MinimumPathCovering(n, newEdges)
	return res
}

type BipartiteFlow struct {
	N, M           int
	MatchL, MatchR []int
	timeStamp      int
	g, rg          [][]int
	dist           []int
	used           []int
	alive          []bool
	matched        bool
}

// 指定左侧点数n，右侧点数m，初始化二分图最大流.
func NewBipartiteFlow(n, m int) *BipartiteFlow {
	g, rg := make([][]int, n), make([][]int, m)
	matchL, matchR := make([]int, n), make([]int, m)
	used, alive := make([]int, n), make([]bool, n)
	for i := 0; i < n; i++ {
		matchL[i] = -1
		alive[i] = true
	}
	for i := 0; i < m; i++ {
		matchR[i] = -1
	}

	return &BipartiteFlow{
		N:      n,
		M:      m,
		g:      g,
		rg:     rg,
		MatchL: matchL,
		MatchR: matchR,
		used:   used,
		alive:  alive,
	}
}

// 增加一条边u-v.u属于左侧点集，v属于右侧点集.
//  !0<=u<n,0<=v<m.
func (bf *BipartiteFlow) AddEdge(u, v int) {
	bf.g[u] = append(bf.g[u], v)
	bf.rg[v] = append(bf.rg[v], u)
}

// 求最大匹配.
//  返回(左侧点,右侧点)的匹配对.
//  !0<=左侧点<n,0<=右侧点<m.
func (bf *BipartiteFlow) MaxMatching() [][2]int {
	bf.matched = true
	for {
		bf.buildAugmentPath()
		bf.timeStamp++
		flow := 0
		for i := 0; i < bf.N; i++ {
			if bf.MatchL[i] == -1 {
				tmp := bf.findMinDistAugmentPath(i)
				if tmp {
					flow++
				}
			}
		}

		if flow == 0 {
			break
		}
	}

	res := [][2]int{}
	for i := 0; i < bf.N; i++ {
		if bf.MatchL[i] >= 0 {
			res = append(res, [2]int{i, bf.MatchL[i]})
		}
	}
	return res
}

// 构建残量图.
//  left: [0,n), right: [n,n+m), S: n+m, T: n+m+1
func (bf *BipartiteFlow) BuildRisidualGraph() [][]int {
	if !bf.matched {
		bf.MaxMatching()
	}

	S := bf.N + bf.M
	T := bf.N + bf.M + 1
	ris := make([][]int, bf.N+bf.M+2)
	for i := 0; i < bf.N; i++ {
		if bf.MatchL[i] == -1 {
			ris[S] = append(ris[S], i)
		} else {
			ris[i] = append(ris[i], S)
		}
	}

	for i := 0; i < bf.M; i++ {
		if bf.MatchR[i] == -1 {
			ris[i+bf.N] = append(ris[i+bf.N], T)
		} else {
			ris[T] = append(ris[T], i+bf.N)
		}
	}

	for i := 0; i < bf.N; i++ {
		for _, j := range bf.g[i] {
			if bf.MatchL[i] == j {
				ris[j+bf.N] = append(ris[j+bf.N], i)
			} else {
				ris[i] = append(ris[i], j+bf.N)
			}
		}
	}

	return ris
}

func (bf *BipartiteFlow) findResidualPath() []bool {
	res := bf.BuildRisidualGraph()
	que := []int{}
	visited := make([]bool, bf.N+bf.M+2)
	que = append(que, bf.N+bf.M)
	visited[bf.N+bf.M] = true
	for len(que) > 0 {
		idx := que[0]
		que = que[1:]
		for _, to := range res[idx] {
			if visited[to] {
				continue
			}
			visited[to] = true
			que = append(que, to)
		}
	}
	return visited
}

func (bf *BipartiteFlow) buildAugmentPath() {
	que := []int{}
	bf.dist = make([]int, len(bf.g))
	for i := 0; i < len(bf.g); i++ {
		bf.dist[i] = -1
	}
	for i := 0; i < bf.N; i++ {
		if bf.MatchL[i] == -1 {
			que = append(que, i)
			bf.dist[i] = 0
		}
	}
	for len(que) > 0 {
		a := que[0]
		que = que[1:]
		for _, b := range bf.g[a] {
			c := bf.MatchR[b]
			if c >= 0 && bf.dist[c] == -1 {
				bf.dist[c] = bf.dist[a] + 1
				que = append(que, c)
			}
		}
	}
}

func (bf *BipartiteFlow) findMinDistAugmentPath(a int) bool {
	bf.used[a] = bf.timeStamp
	for _, b := range bf.g[a] {
		c := bf.MatchR[b]
		if c < 0 || (bf.used[c] != bf.timeStamp && bf.dist[c] == bf.dist[a]+1 && bf.findMinDistAugmentPath(c)) {
			bf.MatchR[b] = a
			bf.MatchL[a] = b
			return true
		}
	}
	return false
}
