package main

const INF int = 1e18

// Floyd 求多源最短路.
// 如果返回值dist[i][i]<0,则说明存在负环.
func Floyd(n int, edges [][3]int, directed bool) [][]int {
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = INF
		}
		dist[i][i] = 0
	}

	for _, road := range edges {
		u, v, w := road[0], road[1], road[2]
		dist[u][v] = min(w, dist[u][v]) // 有重边，取最小值
		if !directed {
			dist[v][u] = min(w, dist[v][u])
		}
	}

	// dis[k][i][j] 表示「经过若干个编号不超过 k 的节点」时，从 i 到 j 的最短路长度
	for k := 0; k < n; k++ { // 经过的中转点
		for i := 0; i < n; i++ {
			if dist[i][k] == INF { // 稀疏图优化
				continue
			}
			for j := 0; j < n; j++ {
				if dist[k][j] == INF { // 稀疏图优化
					continue
				}
				// 松弛：如果一条边可以被松弛了，说明这条边就没有必要留下了
				dist[i][j] = min(dist[i][j], dist[i][k]+dist[k][j])
			}
		}
	}
	return dist
}

type FloydStatic struct {
	built         bool
	n             int
	dist          [][]int
	pre           [][]int
	directedEdges [][3]int
}

func NewFloydStatic(n int) *FloydStatic {
	dist := make([][]int, n)
	pre := make([][]int, n)
	for i := 0; i < n; i++ {
		row1 := make([]int, n)
		row2 := make([]int, n)
		for j := 0; j < n; j++ {
			row1[j] = INF
			row2[j] = i
		}
		row1[i] = 0
		dist[i] = row1
		pre[i] = row2
	}
	return &FloydStatic{n: n, dist: dist, pre: pre}
}

func (fs *FloydStatic) AddEdge(u, v, w int) {
	fs.AddDirectedEdge(u, v, w)
	fs.AddDirectedEdge(v, u, w)
}

func (fs *FloydStatic) AddDirectedEdge(u, v, w int) {
	fs.built = false
	fs.directedEdges = append(fs.directedEdges, [3]int{u, v, w})
}

func (fs *FloydStatic) Build() {
	if fs.built {
		return
	}
	fs.built = true
	for i := 0; i < len(fs.directedEdges); i++ {
		e := &fs.directedEdges[i]
		u, v, w := e[0], e[1], e[2]
		fs.dist[u][v] = min(w, fs.dist[u][v])
	}

	for k := 0; k < fs.n; k++ {
		for i := 0; i < fs.n; i++ {
			if fs.dist[i][k] == INF {
				continue
			}
			for j := 0; j < fs.n; j++ {
				if fs.dist[k][j] == INF {
					continue
				}
				cand := fs.dist[i][k] + fs.dist[k][j]
				if fs.dist[i][j] > cand {
					fs.dist[i][j] = cand
					fs.pre[i][j] = fs.pre[k][j]
				}
			}
		}
	}
}

// 求出从`start`到`target`的最短路径长度,如果不存在这样的路径,返回-1.
func (fs *FloydStatic) Dist(start, target int) int {
	if !fs.built {
		fs.Build()
	}
	res := fs.dist[start][target]
	if res == INF {
		return -1
	}
	return res
}

// 求出从`start`到`target`的最短路径.如果不存在这样的路径,返回空数组.
func (fs *FloydStatic) GetPath(start, target int) []int {
	if !fs.built {
		fs.Build()
	}
	if fs.Dist(start, target) == -1 {
		return nil
	}
	cur := target
	path := []int{target}
	for cur != start {
		cur = fs.pre[start][cur]
		path = append(path, cur)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// O(n)判断是否存在负环.
func (fs *FloydStatic) HasNegativeCycle() bool {
	if !fs.built {
		fs.Build()
	}
	for i := 0; i < fs.n; i++ {
		if fs.dist[i][i] < 0 {
			return true
		}
	}
	return false
}

// 支持O(n^2)向图中添加边的Floyd.
type FloydDynamic struct {
	built         bool
	n             int
	dist          [][]int
	directedEdges [][3]int
}

func NewFloydDynamic(n int) *FloydDynamic {
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, n)
		for j := 0; j < n; j++ {
			row[j] = INF
		}
		row[i] = 0
		dist[i] = row
	}
	return &FloydDynamic{n: n, dist: dist}
}

func (fd *FloydDynamic) AddEdge(u, v, w int) {
	fd.AddDirectedEdge(u, v, w)
	fd.AddDirectedEdge(v, u, w)
}

func (fd *FloydDynamic) AddDirectedEdge(u, v, w int) {
	fd.built = false
	fd.directedEdges = append(fd.directedEdges, [3]int{u, v, w})
}

func (fd *FloydDynamic) Build() {
	if fd.built {
		return
	}
	fd.built = true
	for i := 0; i < len(fd.directedEdges); i++ {
		e := &fd.directedEdges[i]
		u, v, w := e[0], e[1], e[2]
		fd.dist[u][v] = min(w, fd.dist[u][v])
	}

	for k := 0; k < fd.n; k++ {
		for i := 0; i < fd.n; i++ {
			if fd.dist[i][k] == INF {
				continue
			}
			for j := 0; j < fd.n; j++ {
				if fd.dist[k][j] == INF {
					continue
				}

				cand := fd.dist[i][k] + fd.dist[k][j]
				if fd.dist[i][j] > cand {
					fd.dist[i][j] = cand
				}
			}
		}
	}
}

// 向边集中添加一条边.
func (fd *FloydDynamic) UpdateEdge(u, v, w int, directed bool) {
	fd.built = false
	if directed {
		fd._updateDirectedEdge(u, v, w)
	} else {
		fd._updateEdge(u, v, w)
	}
}

// 求出从`start`到`target`的最短路径长度,如果不存在这样的路径,返回-1.
func (fd *FloydDynamic) Dist(start, target int) int {
	if !fd.built {
		fd.Build()
	}
	res := fd.dist[start][target]
	if res == INF {
		return -1
	}
	return res
}

func (fd *FloydDynamic) HasNegativeCycle() bool {
	if !fd.built {
		fd.Build()
	}
	for i := 0; i < fd.n; i++ {
		if fd.dist[i][i] < 0 {
			return true
		}
	}
	return false
}

func (fd *FloydDynamic) _updateDirectedEdge(u, v, w int) {
	for i := 0; i < fd.n; i++ {
		for j := 0; j < fd.n; j++ {
			fd.dist[i][j] = min(fd.dist[i][j], fd.dist[i][u]+w+fd.dist[v][j])
		}
	}
}

func (fd *FloydDynamic) _updateEdge(u, v, w int) {
	for i := 0; i < fd.n; i++ {
		for j := 0; j < fd.n; j++ {
			fd.dist[i][j] = min(fd.dist[i][j], fd.dist[i][u]+w+fd.dist[v][j])
			fd.dist[i][j] = min(fd.dist[i][j], fd.dist[i][v]+w+fd.dist[u][j])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
