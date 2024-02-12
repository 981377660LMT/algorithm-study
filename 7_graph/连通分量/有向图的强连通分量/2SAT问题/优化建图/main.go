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

func main() {
	// CF1903FOnlineSCC()
	// CF1903FDivideInterval()
	AT_arc069_d()
	// P6378前后缀()
	// P6378OnlineSCC()
}

// 给你一个无向图，你需要求出一个"点覆盖"使得选出点编号之间的最小差距最大化。
// https://www.luogu.com.cn/problem/CF1903F
// 考虑 k 固定的时候怎么做。
// !点覆盖问题可以通过强制一个点不取另外一个点取转 2-SAT。
// 题目中的限制可以转化为两条：
// 对于一条边(u,v)
// !1.u不选=>v必选
// !2.u选了=>所有满足|u-v|<k的v都不能选(u!=v)
// kosaraju 在线求scc即可.
func CF1903FOnlineSCC() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(n int, edges [][2]int) int {
		vid := func(i int, b bool) int {
			if b {
				return 2*i + 1
			} else {
				return 2 * i
			}
		}
		// 所有点的差>=k时能否选出一个点覆盖.
		check := func(k int) bool {
			adjList := make([][]int, n)
			for _, e := range edges {
				u, v := e[0], e[1]
				adjList[u] = append(adjList[u], v)
				adjList[v] = append(adjList[v], u)
			}
			adjListPtr := make([][]int, 2)
			for i := 0; i < 2; i++ {
				adjListPtr[i] = make([]int, n)
				for j := 0; j < n; j++ {
					adjListPtr[i][j] = len(adjList[j]) - 1
				}
			}

			used := make([][]bool, 2)
			for i := 0; i < 2; i++ {
				used[i] = make([]bool, 2*n)
			}
			s0 := NewFastSetFrom(n, func(i int) bool { return true }) // !原图上的点
			s1 := NewFastSetFrom(n, func(i int) bool { return true }) // !反图上的点

			setUsed := func(v int, rev bool) {
				if rev {
					used[1][v] = true
				} else {
					used[0][v] = true
				}
			}

			findUnused := func(cur int, rev bool) int {
				b := cur & 1 // 0: false, 1: true
				cur >>= 1

				if !rev { // !原图上寻找未访问过的位置
					// 约束1：如果当前点不选(0), 那么所有的邻居必选(1)
					if b == 0 {
						nexts := adjList[cur]
						ptr := &adjListPtr[0][cur]
						for *ptr >= 0 {
							last := nexts[*ptr]
							*ptr--
							if id := vid(last, true); used[0][id] {
								continue
							} else {
								return id
							}
						}
					}

					if b == 1 {
						// 约束2:如果当点选了(1), 那么距离不超过k的点不能选(0)
						left := cur
						for left >= 0 {
							left = s0.Prev(left - 1)
							if left == -1 || left <= cur-k {
								break
							}
							if id := vid(left, false); used[0][id] {
								s0.Erase(left)
								continue
							} else {
								return id
							}
						}
						right := cur
						for right < n {
							right = s0.Next(right + 1)
							if right == n || right >= cur+k {
								break
							}
							if id := vid(right, false); used[0][id] {
								s0.Erase(right)
								continue
							} else {
								return id
							}
						}
					}
				} else { // !反图上寻找未访问过的位置
					if b == 1 {
						// 约束1：如果当前点选了1, 那么所有的邻居是0
						nexts := adjList[cur]
						ptr := &adjListPtr[1][cur]
						for *ptr >= 0 {
							last := nexts[*ptr]
							*ptr--
							if id := vid(last, false); used[1][id] {
								continue
							} else {
								return id
							}
						}
					}

					if b == 0 {
						// 约束2:如果当点是0, 那么距离不超过k的点为1
						left := cur
						for left >= 0 {
							left = s1.Prev(left - 1)
							if left == -1 || left <= cur-k {
								break
							}
							if id := vid(left, true); used[1][id] {
								s1.Erase(left)
								continue
							} else {
								return id
							}
						}
						right := cur
						for right < n {
							right = s1.Next(right + 1)
							if right == n || right >= cur+k {
								break
							}
							if id := vid(right, true); used[1][id] {
								s1.Erase(right)
								continue
							} else {
								return id
							}
						}
					}
				}

				return -1
			}

			_, belong := OnlineSCC(2*n, setUsed, findUnused)
			for i := 0; i < n; i++ {
				if belong[vid(i, false)] == belong[vid(i, true)] { // 有环
					return false
				}
			}
			return true
		}

		left, right := 1, n
		for left <= right {
			mid := (left + right) / 2
			if check(mid) {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return right
	}

	var T int
	fmt.Fscan(in, &T)
	for i := 0; i < T; i++ {
		var n, m int
		fmt.Fscan(in, &n, &m)
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &edges[j][0], &edges[j][1])
			edges[j][0]--
			edges[j][1]--
		}
		res := solve(n, edges)
		fmt.Fprintln(out, res)
	}
}

// 给你一个无向图，你需要求出一个"点覆盖"使得选出点编号之间的最小差距最大化。
// https://www.luogu.com.cn/problem/CF1903F
// 考虑 k 固定的时候怎么做。
// 题目中的限制可以转化为两条：
// 对于一条边(u,v)
// !u不选=>v必选
// !u选了=>所有满足|u-v|<k的v都不能选(u!=v)
// 2sat问题，但是边数为O(nk)，每次连边是单点连边，或者点向区间连边，线段树优化建图即可
func CF1903FDivideInterval() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(n int32, edges [][2]int32) int32 {
		D := NewDivideInterval(n)

		// 所有点的差>=k时能否选出一个点覆盖.
		check := func(k int32) bool {
			S := NewTwoSat(D.Size())
			for _, e := range edges {
				u, v := e[0], e[1]
				S.AddOr(D.Id(u), D.Id(v))
			}

			D.PushDown(func(parent, child int32) { // 将线段树上一个点和其代表的所有点连边
				S.AddDirectedEdge(child, parent)
			})
			for i := int32(0); i < n; i++ {
				id := D.Id(i)
				a, b := i-k+1, i-1 // [i-k+1,i-1]
				D.EnumerateSegment(a, b+1, func(segmentId int32) {
					S.AddNand(id, segmentId)
				}, false)
				a, b = i+1, i+k-1 // [i+1,i+k-1]
				D.EnumerateSegment(a, b+1, func(segmentId int32) {
					S.AddNand(id, segmentId)
				}, false)
			}

			_, ok := S.Solve()
			return ok
		}

		left, right := int32(1), n
		for left <= right {
			mid := (left + right) / 2
			if check(mid) {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return right
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var n, m int32
		fmt.Fscan(in, &n, &m)
		edges := make([][2]int32, m)
		for j := int32(0); j < m; j++ {
			fmt.Fscan(in, &edges[j][0], &edges[j][1])
			edges[j][0]--
			edges[j][1]--
		}
		res := solve(n, edges)
		fmt.Fprintln(out, res)
	}
}

// [ARC069F] Flags
// https://www.luogu.com.cn/problem/AT_arc069_d
// Snuke 将n个标志放在一条线上。
// 第 i 个标志可以放置在坐标 xi​ 或坐标 yi​ 上。
// Snuke 认为当他们中的两个之间的最小距离 d 更大时，标志看起来更好。找出 d 的最大可能值。
// 2<=n<=1e4, 1<=xi,yi<=1e9
//
// 将原问题转化为2*n个顶点的问题，是否能选择每个顶点.xi对应顶点0-n-1, yi对应顶点n-2n-1
// !排序后，
// 对于一条边 (u,v)
// !1.u不选 => v 必选
// !2.u选了 => v不能选
// !3.u选了=>所有满足|u-v|<k的v都不能选(u!=v)
// 类似 CF1903FOnlineSCC
func AT_arc069_d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	xs, ys := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &ys[i])
	}

	vid := func(i int, b bool) int {
		if b {
			return 2*i + 1
		} else {
			return 2 * i
		}
	}
	neighbor := func(i int) int {
		if i < n {
			return i + n
		} else {
			return i - n
		}
	}
	// [[1 0] [1 2] [2 1] [3 0] [5 1] [9 2]]
	sortedPoints := make([][2]int, 0, 2*n) // (xi, id) or (yi, id)
	for i := 0; i < n; i++ {
		sortedPoints = append(sortedPoints, [2]int{xs[i], i})
		sortedPoints = append(sortedPoints, [2]int{ys[i], i + n})
	}
	sort.Slice(sortedPoints, func(i, j int) bool { return sortedPoints[i][0] < sortedPoints[j][0] })
	order := make([]int, 2*n)  // 顶点 => 排序后的位置
	rOrder := make([]int, 2*n) // 排序后的位置 => 顶点
	for i, p := range sortedPoints {
		id := p[1]
		order[id] = i
		rOrder[i] = id
	}

	// 最小距离>=k时，是否能满足条件.
	check := func(k int) bool {
		used := make([][]bool, 2)
		for i := 0; i < 2; i++ {
			used[i] = make([]bool, 2*2*n)
		}
		s0 := NewFastSetFrom(2*n, func(i int) bool { return true }) // !原图上未访问的点
		s1 := NewFastSetFrom(2*n, func(i int) bool { return true }) // !反图上未访问的点

		setUsed := func(cur int, rev bool) {
			if rev {
				used[1][cur] = true
			} else {
				used[0][cur] = true
			}
		}

		findUnused := func(cur int, rev bool) int {
			b := cur & 1 // 0: false, 1: true
			cur >>= 1
			curX := sortedPoints[order[cur]][0]

			if !rev {
				// !原图上寻找未访问过的位置。

				if b == 0 {
					next := neighbor(cur)
					if id := vid(next, true); !used[0][id] {
						return id
					}
				}

				if b == 1 {
					pos := order[cur]
					for pos >= 0 {
						pos = s0.Prev(pos - 1)
						if pos == -1 || curX-sortedPoints[pos][0] >= k {
							break
						}
						if id := vid(rOrder[pos], false); used[0][id] {
							s0.Erase(pos)
							continue
						} else {
							return id
						}
					}
					pos = order[cur]
					for pos < 2*n {
						pos = s0.Next(pos + 1)
						if pos == 2*n || sortedPoints[pos][0]-curX >= k {
							break
						}
						if id := vid(rOrder[pos], false); used[0][id] {
							s0.Erase(pos)
							continue
						} else {
							return id
						}
					}
				}

			} else {
				// !反图上寻找未访问过的位置

				if b == 1 {
					next := neighbor(cur)
					if id := vid(next, false); !used[1][id] {
						return id
					}
				}

				if b == 0 {

					pos := order[cur]
					for pos >= 0 {
						pos = s1.Prev(pos - 1)
						if pos == -1 || curX-sortedPoints[pos][0] >= k {
							break
						}
						if id := vid(rOrder[pos], true); used[1][id] {
							s1.Erase(pos)
							continue
						} else {
							return id
						}
					}
					pos = order[cur]
					for pos < 2*n {
						pos = s1.Next(pos + 1)
						if pos == 2*n || sortedPoints[pos][0]-curX >= k {
							break
						}
						if id := vid(rOrder[pos], true); used[1][id] {
							s1.Erase(pos)
							continue
						} else {
							return id
						}
					}
				}

			}

			return -1
		}

		_, belong := OnlineSCC(2*2*n, setUsed, findUnused)
		for i := 0; i < 2*n; i++ {
			if belong[vid(i, false)] == belong[vid(i, true)] { // 有环
				return false
			}
		}
		return true
	}

	left, right := 0, int(1e9+10)
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	fmt.Fprintln(out, right)
}

// P6378 [PA2010] Riddle
// https://www.luogu.com.cn/problem/P6378前后缀
// n 个点 m 条边的无向图被分成 k 个部分。每个部分包含一些点。
// !请选择一些关键点，使得每个部分`恰有`一个关键点，且每条边至少有一个端点是关键点。
// n,m<=1e6
//
// !注意到：`连向其余所有点`，可以看做连向这个点之前的前缀 和 这个点之后的后缀。
// 每一个部分进行一次前后缀优化即可.
// 0-n-1: 原图上的点
// n-2n-1: 前缀上的点
// 2n-3n-1: 后缀上的点
func P6378前后缀() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int32
	fmt.Fscan(in, &n, &m, &k)

	S := NewTwoSat(3 * n)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		S.AddOr(u, v)
	}

	for i := int32(0); i < k; i++ {
		var x int32
		fmt.Fscan(in, &x)
		g := make([]int32, x)
		for j := int32(0); j < x; j++ {
			fmt.Fscan(in, &g[j])
			g[j]--
		}

		for i := 0; i < len(g)-1; i++ {
			pre1, pre2 := g[i]+n, g[i+1]+n
			S.AddDirectedEdge(pre2, pre1)
			suf1, suf2 := g[i]+2*n, g[i+1]+2*n
			S.AddDirectedEdge(suf1, suf2)
		}
		for i, cur := range g {
			curPre, curSuf := cur+n, cur+2*n
			S.AddDirectedEdge(curPre, S.Rev(cur))
			S.AddDirectedEdge(curSuf, S.Rev(cur))
			// 每个组内恰好一个关键点.
			if i > 0 {
				S.AddDirectedEdge(cur, g[i-1]+n) // cur选了，则会连到前缀每个点的Rev(假)
			}
			if i < len(g)-1 {
				S.AddDirectedEdge(cur, g[i+1]+2*n) // cur选了，则会连到后缀每个点的Rev(假)
			}
		}
	}

	_, ok := S.Solve()
	if ok {
		fmt.Fprintln(out, "TAK")
	} else {
		fmt.Fprintln(out, "NIE")
	}
}

// P6378 [PA2010] Riddle
// https://www.luogu.com.cn/problem/P6378
// n 个点 m 条边的无向图被分成 k 个部分。每个部分包含一些点。
// !请选择一些关键点，使得每个部分`恰有`一个关键点，且每条边至少有一个端点是关键点。
// n,m<=1e6
//
// kosaraju 隐式建图求解(implicit scc).
func P6378OnlineSCC() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int32
	fmt.Fscan(in, &n, &m, &k)

	adjList := make([][]int32, n)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
	}
	adjListPtr := make([][]int32, 2)
	for i := 0; i < 2; i++ {
		adjListPtr[i] = make([]int32, n)
		for j := int32(0); j < n; j++ {
			adjListPtr[i][j] = int32(len(adjList[j]) - 1)
		}
	}

	groups := make([][]int32, k)
	groupIds := make([]int32, n)
	for i := int32(0); i < k; i++ {
		var x int32
		fmt.Fscan(in, &x)
		groups[i] = make([]int32, x)
		for j := int32(0); j < x; j++ {
			fmt.Fscan(in, &groups[i][j])
			groups[i][j]--
			groupIds[groups[i][j]] = i
		}
	}
	groupsPtr := make([][]int32, 2)
	for i := 0; i < 2; i++ {
		groupsPtr[i] = make([]int32, k)
		for j := int32(0); j < k; j++ {
			groupsPtr[i][j] = int32(len(groups[j]) - 1)
		}
	}

	vid := func(i int32, b bool) int32 {
		if b {
			return 2*i + 1
		} else {
			return 2 * i
		}
	}

	used := make([][]bool, 2)
	for i := 0; i < 2; i++ {
		used[i] = make([]bool, 2*n)
	}
	setUsed := func(cur int, rev bool) {
		if !rev {
			used[0][cur] = true
		} else {
			used[1][cur] = true
		}
	}

	findUnused := func(cur int, rev bool) int {
		b := cur & 1 // 0: false, 1: true
		cur >>= 1

		if !rev {
			// !原图上寻找未访问过的位置

			// 如果当前点不选，那么邻居必选
			if b == 0 {
				nexts := adjList[cur]
				ptr := &adjListPtr[0][cur]
				for *ptr >= 0 {
					last := nexts[*ptr]
					*ptr--
					if id := vid(last, true); used[0][id] {
						continue
					} else {
						return int(id)
					}
				}
			}

			// 如果当前点选了，那么同一个组内的点不能选
			if b == 1 {
				g := groups[groupIds[cur]]
				ptr := &groupsPtr[0][groupIds[cur]]
				for *ptr >= 0 {
					last := g[*ptr]
					*ptr--
					if id := vid(last, false); used[0][id] {
						continue
					} else {
						return int(id)
					}
				}
			}
		} else {
			if b == 1 {
				nexts := adjList[cur]
				ptr := &adjListPtr[1][cur]
				for *ptr >= 0 {
					last := nexts[*ptr]
					*ptr--
					if id := vid(last, false); used[1][id] {
						continue
					} else {
						return int(id)
					}
				}
			}

			if b == 0 {
				g := groups[groupIds[cur]]
				ptr := &groupsPtr[1][groupIds[cur]]
				for *ptr >= 0 {
					last := g[*ptr]
					*ptr--
					if id := vid(last, true); used[1][id] {
						continue
					} else {
						return int(id)
					}
				}
			}
		}

		return -1
	}

	_, belong := OnlineSCC(int(2*n), setUsed, findUnused)
	ok := true
	for i := int32(0); i < n; i++ {
		if belong[vid(i, false)] == belong[vid(i, true)] {
			ok = false
			break
		}
	}
	if ok {
		fmt.Fprintln(out, "TAK")
	} else {
		fmt.Fprintln(out, "NIE")
	}
}

// kosaraju 在线求有向图的强连通分量.常用于2-sat优化建图问题.
//
//	setUsed(v, rev)：将 v 设置为已使用, rev 表示是否是反图
//	findUnused(v, rev)：返回未使用过的点中与 v 相邻的点, rev 表示是否是反图.不存在时返回 -1.
//
// 返回强连通分量的个数和每个点所属的分量编号.
// !注意按照0到count-1的遍历结果为拓扑序逆序.
//
// 步骤:
// https://www.cnblogs.com/RioTian/p/14026585.html
// 1.通过后序遍历的方式遍历整个有向图，并且维护每个点的出栈顺序
// 2.根据出栈顺序从大到小遍历反图
// 3.对点u来说，在遍历反图时所有能够到达的v都和u在一个强连通分量当中
func OnlineSCC(n int, setUsed func(cur int, rev bool), findUnused func(cur int, rev bool) int) (count int, belong []int) {
	belong = make([]int, n)

	stack := make([]int, n)
	ptr := n

	visited := make([]bool, n)
	var dfs1 func(v int)
	dfs1 = func(v int) { // 原图
		visited[v] = true
		setUsed(v, false)
		for {
			to := findUnused(v, false)
			if to == -1 {
				break
			}
			dfs1(to)
		}
		ptr--
		stack[ptr] = v
	}
	for v := 0; v < n; v++ {
		if !visited[v] {
			dfs1(v)
		}
	}

	visited = make([]bool, n)
	var dfs2 func(v int)
	dfs2 = func(v int) { // 反图
		visited[v] = true
		belong[v] = count
		setUsed(v, true)
		for {
			to := findUnused(v, true)
			if to == -1 {
				break
			}
			dfs2(to)
		}
	}
	for _, v := range stack {
		if !visited[v] {
			dfs2(v)
			count++
		}
	}

	return
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

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
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

type DivideInterval struct {
	Offset int32 // 线段树中一共offset+n个节点,offset+i对应原来的第i个节点.
	n      int32
}

// 线段树分割区间.
// 将长度为n的序列搬到长度为offset+n的线段树上, 以实现快速的区间操作.
func NewDivideInterval(n int32) *DivideInterval {
	offset := int32(1)
	for offset < n {
		offset <<= 1
	}
	return &DivideInterval{Offset: offset, n: n}
}

// 获取原下标为i的元素在树中的(叶子)编号.
func (d *DivideInterval) Id(rawIndex int32) int32 {
	return rawIndex + d.Offset
}

// O(logn) 顺序遍历`[start,end)`区间对应的线段树节点.
// sorted表示是否按照节点编号从小到大的顺序遍历.
func (d *DivideInterval) EnumerateSegment(start, end int32, f func(segmentId int32), sorted bool) {
	if start < 0 {
		start = 0
	}
	if end > d.n {
		end = d.n
	}
	if start >= end {
		return
	}

	if sorted {
		for _, i := range d.getSegmentIds(start, end) {
			f(i)
		}
	} else {
		for start, end = start+d.Offset, end+d.Offset; start < end; start, end = start>>1, end>>1 {
			if start&1 == 1 {
				f(start)
				start++
			}
			if end&1 == 1 {
				end--
				f(end)
			}
		}
	}
}

func (d *DivideInterval) EnumeratePoint(index int32, f func(segmentId int32)) {
	if index < 0 || index >= d.n {
		return
	}
	index += d.Offset
	for index > 0 {
		f(index)
		index >>= 1
	}
}

// O(n) 从根向叶子方向push.
func (d *DivideInterval) PushDown(f func(parent, child int32)) {
	for p := int32(1); p < d.Offset; p++ {
		f(p, p<<1)
		f(p, p<<1|1)
	}
}

// O(n) 从叶子向根方向update.
func (d *DivideInterval) PushUp(f func(parent, child1, child2 int32)) {
	for p := d.Offset - 1; p > 0; p-- {
		f(p, p<<1, p<<1|1)
	}
}

// 线段树的节点个数.
func (d *DivideInterval) Size() int32 {
	return d.Offset + d.n
}

func (d *DivideInterval) IsLeaf(segmentId int32) bool {
	return segmentId >= d.Offset
}

func (d *DivideInterval) Depth(u int32) int32 {
	if u == 0 {
		return 0
	}
	return int32(bits.LeadingZeros32(uint32(u))) - 1
}

// 线段树(完全二叉树)中两个节点的最近公共祖先(两个二进制数字的最长公共前缀).
func (d *DivideInterval) Lca(u, v int32) int32 {
	if u == v {
		return u
	}
	if u > v {
		u, v = v, u
	}
	depth1 := d.Depth(u)
	depth2 := d.Depth(v)
	diff := u ^ (v >> (depth2 - depth1))
	if diff == 0 {
		return u
	}
	len := bits.Len32(uint32(diff))
	return u >> len
}

func (d *DivideInterval) getSegmentIds(start, end int32) []int32 {
	if !(0 <= start && start <= end && end <= d.n) {
		return nil
	}
	var leftRes, rightRes []int32
	for start, end = start+d.Offset, end+d.Offset; start < end; start, end = start>>1, end>>1 {
		if start&1 == 1 {
			leftRes = append(leftRes, start)
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = append(rightRes, end)
		}
	}
	for i := len(rightRes) - 1; i >= 0; i-- {
		leftRes = append(leftRes, rightRes[i])
	}
	return leftRes
}

type TwoSat struct {
	size  int32
	graph [][]int32
}

func NewTwoSat(n int32) *TwoSat {
	return &TwoSat{size: n, graph: make([][]int32, n*2)}
}

// u -> v <=> !v -> !u
func (ts *TwoSat) AddIf(u, v int32) {
	ts.AddDirectedEdge(u, v)
	ts.AddDirectedEdge(ts.Rev(v), ts.Rev(u))
}

// u,v 中至少有一个为真.
// u or v <=> !u -> v
func (ts *TwoSat) AddOr(u, v int32) {
	ts.AddIf(ts.Rev(u), v)
}

// u,v 中恰好有一个为真, 一个为假.
// u xor v <=> u -> !v, v -> !u, !u -> v, !v -> u
func (ts *TwoSat) AddXor(u, v int32) {
	ts.AddOr(u, v)
	ts.AddNand(u, v)
}

// u,v 不同时为真.
// u nand v <=> u -> !v
func (ts *TwoSat) AddNand(u, v int32) {
	ts.AddIf(u, ts.Rev(v))
}

// 手动添加边(不推荐).常用于优化建图时.
func (ts *TwoSat) AddDirectedEdge(u, v int32) {
	ts.graph[u] = append(ts.graph[u], v)
}

// u <=> !u -> u
func (ts *TwoSat) SetTrue(u int32) {
	ts.AddDirectedEdge(ts.Rev(u), u)
}

// !u <=> u -> !u
func (ts *TwoSat) SetFalse(u int32) {
	ts.AddDirectedEdge(u, ts.Rev(u))
}

func (ts *TwoSat) Rev(u int32) int32 {
	if u >= ts.size {
		return u - ts.size
	}
	return u + ts.size
}

func (ts *TwoSat) Solve() (res []bool, ok bool) {
	_, belong := StronglyConnectedComponentInt32(ts.graph)
	res = make([]bool, ts.size)
	for i := int32(0); i < ts.size; i++ {
		if belong[i] == belong[ts.Rev(i)] {
			return
		}
		res[i] = belong[i] > belong[ts.Rev(i)]
	}
	ok = true
	return
}

// 有向图强连通分量分解.
// 返回强连通分量的个数和每个点所属的分量编号.
// !注意按照0到count-1的遍历结果为拓扑序逆序.
func StronglyConnectedComponentInt32(graph [][]int32) (count int32, belong []int32) {
	n := int32(len(graph))
	belong = make([]int32, n)
	low := make([]int32, n)
	order := make([]int32, n)
	for i := range order {
		order[i] = -1
	}
	now := int32(0)
	path := []int32{}

	var dfs func(int32)
	dfs = func(v int32) {
		low[v] = now
		order[v] = now
		now++
		path = append(path, v)
		for _, to := range graph[v] {
			if order[to] == -1 {
				dfs(to)
				low[v] = min32(low[v], low[to])
			} else {
				low[v] = min32(low[v], order[to])
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

	for i := int32(0); i < n; i++ {
		if order[i] == -1 {
			dfs(i)
		}
	}
	for i := int32(0); i < n; i++ {
		belong[i] = count - 1 - belong[i]
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
