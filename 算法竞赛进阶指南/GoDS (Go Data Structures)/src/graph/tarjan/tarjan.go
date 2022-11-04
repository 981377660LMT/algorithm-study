package tarjan

import (
	"container/heap"
	"io"
)

// SCC Tarjan
// 常数比 Kosaraju 略小（在 AtCoder 上的测试显示，5e5 的数据下比 Kosaraju 快了约 100ms）
// https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm
// https://oi-wiki.org/graph/scc/#tarjan
// https://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/TarjanSCC.java.html
// https://stackoverflow.com/questions/32750511/does-tarjans-scc-algorithm-give-a-topological-sort-of-the-scc
// 与最小割结合 https://www.luogu.com.cn/problem/P4126
func (*graph) sccTarjan(g [][]int, min func(int, int) int) (scc [][]int, sid []int) {
	dfn := make([]int, len(g)) // 值从 1 开始
	dfsClock := 0
	stk := []int{}
	inStk := make([]bool, len(g))
	var f func(int) int
	f = func(v int) int {
		dfsClock++
		dfn[v] = dfsClock
		lowV := dfsClock
		stk = append(stk, v)
		inStk[v] = true
		for _, w := range g[v] {
			if dfn[w] == 0 {
				lowW := f(w)
				lowV = min(lowV, lowW)
			} else if inStk[w] { // 找到 v 的反向边 v-w，用 dfn[w] 来更新 lowV
				lowV = min(lowV, dfn[w])
			}
		}
		if dfn[v] == lowV {
			comp := []int{}
			for {
				w := stk[len(stk)-1]
				stk = stk[:len(stk)-1]
				inStk[w] = false
				comp = append(comp, w)
				if w == v {
					break
				}
			}
			scc = append(scc, comp)
		}
		return lowV
	}
	for v, timestamp := range dfn {
		if timestamp == 0 {
			f(v)
		}
	}

	// 由于每个强连通分量都是在它的所有后继强连通分量被求出之后求得的
	// 上面得到的 scc 是拓扑序的逆序
	for i, n := 0, len(scc); i < n/2; i++ {
		scc[i], scc[n-1-i] = scc[n-1-i], scc[i]
	}

	sid = make([]int, len(g))
	for i, cp := range scc {
		for _, v := range cp {
			sid[v] = i
		}
	}

	return
}

// 割点（割顶） cut vertices / articulation points
// https://codeforces.com/blog/entry/68138
// https://oi-wiki.org/graph/cut/#_1
// low(v): 在不经过 v 父亲的前提下能到达的最小的时间戳
// 模板题 https://www.luogu.com.cn/problem/P3388
// LC928 https://leetcode-cn.com/problems/minimize-malware-spread-ii/
func (*graph) findCutVertices(n int, g [][]int, min func(int, int) int) (isCut []bool) {
	isCut = make([]bool, n)
	dfn := make([]int, n) // 值从 1 开始
	dfsClock := 0
	var f func(v, fa int) int
	f = func(v, fa int) int {
		dfsClock++
		dfn[v] = dfsClock
		lowV := dfsClock
		childCnt := 0
		for _, w := range g[v] {
			if dfn[w] == 0 {
				childCnt++
				lowW := f(w, v)
				if lowW >= dfn[v] { // 以 w 为根的子树中没有反向边能连回 v 的祖先（可以连到 v 上，这也算割顶）
					isCut[v] = true
				}
				lowV = min(lowV, lowW)
			} else if w != fa { // 找到 v 的反向边 v-w，用 dfn[w] 来更新 lowV
				lowV = min(lowV, dfn[w])
			}
		}
		if fa == -1 && childCnt == 1 { // 特判：只有一个儿子的树根，删除后并没有增加连通分量的个数，这种情况下不是割顶
			isCut[v] = false
		}
		return lowV
	}
	for v, timestamp := range dfn {
		if timestamp == 0 {
			f(v, -1)
		}
	}

	cuts := []int{}
	for v, is := range isCut {
		if is {
			cuts = append(cuts, v) // v+1
		}
	}

	return
}

// 桥（割边）
// https://oi-wiki.org/graph/cut/#_4
// https://algs4.cs.princeton.edu/41graph/Bridge.java.html
// 模板题 LC1192 https://leetcode-cn.com/problems/critical-connections-in-a-network/
//       https://codeforces.com/problemset/problem/1000/E
// 题目推荐 https://cp-algorithms.com/graph/bridge-searching.html#toc-tgt-2
// 与 MST 结合 https://codeforces.com/problemset/problem/160/D
// 与最短路结合 https://codeforces.com/problemset/problem/567/E
// https://codeforces.com/problemset/problem/118/E
// todo 构造 https://codeforces.com/problemset/problem/550/D
func (*graph) findBridges(in io.Reader, n, m int) (isBridge []bool) {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	type neighbor struct{ to, eid int }
	type edge struct{ v, w int }

	g := make([][]neighbor, n)
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		var v, w int
		Fscan(in, &v, &w)
		v--
		w--
		g[v] = append(g[v], neighbor{w, i})
		g[w] = append(g[w], neighbor{v, i})
		edges[i] = edge{v, w}
	}
	isBridge = make([]bool, len(edges))
	dfn := make([]int, len(g)) // 值从 1 开始
	dfsClock := 0
	var f func(int, int) int
	f = func(v, fid int) int { // 使用 fid 而不是 fa，可以兼容重边的情况
		dfsClock++
		dfn[v] = dfsClock
		lowV := dfsClock
		for _, e := range g[v] {
			if w := e.to; dfn[w] == 0 {
				lowW := f(w, e.eid)
				if lowW > dfn[v] { // 以 w 为根的子树中没有反向边能连回 v 或 v 的祖先，所以 v-w 必定是桥
					isBridge[e.eid] = true
				}
				lowV = min(lowV, lowW)
			} else if e.eid != fid { // 找到 v 的反向边 v-w，用 dfn[w] 来更新 lowV
				lowV = min(lowV, dfn[w])
			}
		}
		return lowV
	}
	for v, timestamp := range dfn {
		if timestamp == 0 {
			f(v, -1)
		}
	}

	// EXTRA: 所有桥边的下标
	bridgeEIDs := []int{}
	for eid, b := range isBridge {
		if b {
			bridgeEIDs = append(bridgeEIDs, eid)
		}
	}

	return
}

// 无向图的双连通分量 Biconnected Components (BCC)          也叫重连通图
// v-BCC：任意割点都是至少两个不同 v-BCC 的公共点              广义圆方树
// https://oi-wiki.org/graph/bcc/
// https://www.csie.ntu.edu.tw/~hsinmu/courses/_media/dsa_13spring/horowitz_306_311_biconnected.pdf
// 好题 https://codeforces.com/problemset/problem/962/F
// https://leetcode-cn.com/problems/s5kipK/
// 结合树链剖分 https://codeforces.com/problemset/problem/487/E
/*
使用 https://csacademy.com/app/graph_editor/ 显示下面的样例
基础样例 - 一个割点两个简单环
1 2
2 3
3 4
4 1
3 5
5 6
6 7
7 3
基础样例 - 两个割点三个简单环（注意那条含有两个割点的边）
7 3
7 4
1 2
2 3
3 1
3 4
4 5
5 6
6 4
*/
func (G *graph) findVertexBCC(g [][]int, min func(int, int) int) (comps [][]int, bccIDs []int) {
	bccIDs = make([]int, len(g)) // ID 从 1 开始编号
	idCnt := 0
	isCut := make([]bool, len(g))

	dfn := make([]int, len(g)) // 值从 1 开始
	dfsClock := 0
	type edge struct{ v, w int } // eid
	stack := []edge{}
	var f func(v, fa int) int
	f = func(v, fa int) int {
		dfsClock++
		dfn[v] = dfsClock
		lowV := dfsClock
		childCnt := 0
		for _, w := range g[v] {
			e := edge{v, w} // ne.eid
			if dfn[w] == 0 {
				stack = append(stack, e)
				childCnt++
				lowW := f(w, v)
				if lowW >= dfn[v] {
					isCut[v] = true
					idCnt++
					comp := []int{}
					//eids := []int{}
					for {
						e, stack = stack[len(stack)-1], stack[:len(stack)-1]
						if bccIDs[e.v] != idCnt {
							bccIDs[e.v] = idCnt
							comp = append(comp, e.v)
						}
						if bccIDs[e.w] != idCnt {
							bccIDs[e.w] = idCnt
							comp = append(comp, e.w)
						}
						//eids = append(eids, e.eid)
						if e.v == v && e.w == w {
							break
						}
					}
					// 点数和边数相同，说明该 v-BCC 是一个简单环，且环上所有的边只属于一个简单环
					//if len(comp) == len(eids) {
					//	for _, eid := range eids {
					//		onSimpleCycle[eid] = true
					//	}
					//}
					comps = append(comps, comp)
				}
				lowV = min(lowV, lowW)
			} else if w != fa && dfn[w] < dfn[v] {
				stack = append(stack, e)
				lowV = min(lowV, dfn[w])
			}
		}
		if fa == -1 && childCnt == 1 {
			isCut[v] = false
		}
		return lowV
	}
	for v, timestamp := range dfn {
		if timestamp == 0 {
			if len(g[v]) == 0 { // 零度，即孤立点（isolated vertex）
				idCnt++
				bccIDs[v] = idCnt
				comps = append(comps, []int{v})
				continue
			}
			f(v, -1)
		}
	}

	// EXTRA: 缩点
	// BCC 和割点作为新图中的节点，并在每个割点与包含它的所有 BCC 之间连边
	cutIDs := make([]int, len(g))
	for i, is := range isCut {
		if is {
			idCnt++ // 接在 BCC 之后给割点编号
			cutIDs[i] = idCnt
		}
	}
	for v, cp := range comps {
		v++
		for _, w := range cp {
			if w = cutIDs[w]; w > 0 {
				// add(v,w); add(w,v) ...
			}
		}
	}

	return
}

// e-BCC：删除无向图中所有的割边后，剩下的每一个 CC 都是 e-BCC
// 缩点后形成一颗 bridge tree
// 模板题 https://codeforces.com/problemset/problem/1000/E
// 较为综合的一道题 http://codeforces.com/problemset/problem/732/F
func (G *graph) findEdgeBCC(in io.Reader, n, m int) (comps [][]int, bccIDs []int) {
	type neighbor struct{ to, eid int }
	type edge struct{ v, w int }
	g := make([][]neighbor, n)
	edges := make([]edge, m)

	// *copy* 包含读图
	isBridge := G.findBridges(in, n, m)

	// 求原图中每个点的 bccID
	bccIDs = make([]int, len(g))
	idCnt := 0
	var comp []int
	var f2 func(int)
	f2 = func(v int) {
		bccIDs[v] = idCnt
		comp = append(comp, v)
		for _, e := range g[v] {
			if w := e.to; !isBridge[e.eid] && bccIDs[w] == 0 {
				f2(w)
			}
		}
	}
	for i, id := range bccIDs {
		if id == 0 {
			idCnt++
			comp = []int{}
			f2(i)
			comps = append(comps, comp)
		}
	}

	// EXTRA: 缩点，复杂度 O(M)
	// 遍历 edges，若两端点的 bccIDs 不同（割边）则建边
	g2 := make([][]int, idCnt)
	for _, e := range edges {
		if v, w := bccIDs[e.v]-1, bccIDs[e.w]-1; v != w {
			g2[v] = append(g2[v], w)
			g2[w] = append(g2[w], v)
		}
	}

	// 也可以遍历 isBridge，割边两端点 bccIDs 一定不同
	for eid, b := range isBridge {
		if b {
			e := edges[eid]
			v, w := bccIDs[e.v]-1, bccIDs[e.w]-1
			g2[v] = append(g2[v], w)
			g2[w] = append(g2[w], v)
		}
	}

	return
}

type vdPair struct {
	v   int
	dis int64
}
type vdHeap []vdPair

func (h vdHeap) Len() int              { return len(h) }
func (h vdHeap) Less(i, j int) bool    { return h[i].dis < h[j].dis }
func (h vdHeap) Swap(i, j int)         { h[i], h[j] = h[j], h[i] }
func (h *vdHeap) Push(v interface{})   { *h = append(*h, v.(vdPair)) }
func (h *vdHeap) Pop() (v interface{}) { a := *h; *h, v = a[:len(a)-1], a[len(a)-1]; return }
func (h *vdHeap) push(v vdPair)        { heap.Push(h, v) }
func (h *vdHeap) pop() vdPair          { return heap.Pop(h).(vdPair) }
