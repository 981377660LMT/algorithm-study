// !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L2739

package tarjan

import (
	"fmt"
	"io"
)

// !SCC Tarjan (Tarjan 求有向图的强联通分量，缩点成拓扑图)
// 常数比 Kosaraju 略小（在 AtCoder 上的测试显示，5e5 的数据下比 Kosaraju 快了约 100ms）
// https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm
// https://oi-wiki.org/graph/scc/#tarjan
// https://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/TarjanSCC.java.html
// https://stackoverflow.com/questions/32750511/does-tarjans-scc-algorithm-give-a-topological-sort-of-the-scc
// 与最小割结合 https://www.luogu.com.cn/problem/P4126
func sccTarjan(graph [][]int, min func(int, int) int) (scc [][]int, sid []int) {
	dfn := make([]int, len(graph)) // 值从 1 开始
	dfsClock := 0
	stack := []int{}
	inStack := make([]bool, len(graph))
	var dfs func(int) int
	dfs = func(cur int) int {
		dfsClock++
		dfn[cur] = dfsClock
		curLow := dfsClock
		stack = append(stack, cur)
		inStack[cur] = true
		for _, next := range graph[cur] {
			if dfn[next] == 0 {
				nextLow := dfs(next)
				curLow = min(curLow, nextLow)
			} else if inStack[next] { // 找到 cur 的反向边 cur-next，用 dfn[next] 来更新 curLow
				curLow = min(curLow, dfn[next])
			}
		}
		if dfn[cur] == curLow {
			group := []int{}
			for {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inStack[top] = false
				group = append(group, top)
				if top == cur {
					break
				}
			}
			scc = append(scc, group)
		}
		return curLow
	}

	for i, timestamp := range dfn {
		if timestamp == 0 {
			dfs(i)
		}
	}

	// 由于每个强连通分量都是在它的所有后继强连通分量被求出之后求得的
	// 上面得到的 scc 是拓扑序的逆序
	for i, n := 0, len(scc); i < n/2; i++ {
		scc[i], scc[n-1-i] = scc[n-1-i], scc[i]
	}

	sid = make([]int, len(graph))
	for i, group := range scc {
		for _, v := range group {
			sid[v] = i
		}
	}

	{
		// EXTRA: 缩点: 将边 v-w 转换成 sid[v]-sid[w]
		// 缩点后得到了一张 DAG，点的编号范围为 [0,len(scc)-1]
		// !注意这样可能会产生重边，不能有重边时可以用 map 或对每个点排序去重
		// 模板题 点权 https://www.luogu.com.cn/problem/P3387
		// 		 边权 https://codeforces.com/contest/894/problem/E
		// 检测路径是否可达/唯一/无穷 https://codeforces.com/problemset/problem/1547/G
		m := len(scc)
		graph2 := make([][]int, m)
		deg := make([]int, m)
		for i, nexts := range graph {
			sid1 := sid[i]
			for next := range nexts {
				sid2 := sid[next]
				if sid1 != sid2 {
					graph2[sid1] = append(graph2[sid1], sid2)
					deg[sid2]++
				} else {
					// 这里可以记录自环（指 len(scc) == 1 但是有自环）、汇合同一个 SCC 的权值等 ...
				}
			}
		}
	}

	return
}

// 无向图的割点（割顶） cut vertices / articulation points
// https://codeforces.com/blog/entry/68138
// https://oi-wiki.org/graph/cut/#_1
// low(v): 在不经过 v 父亲的前提下能到达的最小的时间戳
// 模板题 https://www.luogu.com.cn/problem/P3388
// LC928 https://leetcode-cn.com/problems/minimize-malware-spread-ii/
func findCutVertices(n int, graph [][]int, min func(int, int) int) (isCut []bool) {
	isCut = make([]bool, n)
	dfn := make([]int, n) // 值从 1 开始
	dfsClock := 0
	var dfs func(cur, pre int) int
	dfs = func(cur, pre int) int {
		dfsClock++
		dfn[cur] = dfsClock
		curLow := dfsClock
		childCount := 0
		for _, next := range graph[cur] {
			if dfn[next] == 0 {
				childCount++
				nextLow := dfs(next, cur)
				if nextLow >= dfn[cur] { // 以 next 为根的子树中没有反向边能连回 cur 的祖先（可以连到 cur 上，这也算割顶）
					isCut[cur] = true
				}
				curLow = min(curLow, nextLow)
			} else if next != pre { // 找到 cur 的反向边 cur-next，用 dfn[next] 来更新 curLow
				curLow = min(curLow, dfn[next])
			}
		}
		if pre == -1 && childCount == 1 { // 特判：只有一个儿子的树根，删除后并没有增加连通分量的个数，这种情况下不是割顶
			isCut[cur] = false
		}
		return curLow
	}

	for i, timestamp := range dfn {
		if timestamp == 0 {
			dfs(i, -1)
		}
	}

	cuts := []int{}
	for i, v := range isCut {
		if v {
			cuts = append(cuts, i) // v+1
		}
	}

	return
}

// 无向图的桥（割边）
// https://oi-wiki.org/graph/cut/#_4
// https://algs4.cs.princeton.edu/41graph/Bridge.java.html
// 模板题 LC1192 https://leetcode-cn.com/problems/critical-connections-in-a-network/
//       https://codeforces.com/problemset/problem/1000/E
// 题目推荐 https://cp-algorithms.com/graph/bridge-searching.html#toc-tgt-2
// 与 MST 结合 https://codeforces.com/problemset/problem/160/D
// 与最短路结合 https://codeforces.com/problemset/problem/567/E
// https://codeforces.com/problemset/problem/118/E
// todo 构造 https://codeforces.com/problemset/problem/550/D
func findBridges(in io.Reader, n, m int) (isBridge []bool) {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	type neighbor struct{ to, eid int }
	type edge struct{ v, w int }

	graph := make([][]neighbor, n)
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		var v, w int
		fmt.Fscan(in, &v, &w)
		v--
		w--
		graph[v] = append(graph[v], neighbor{w, i})
		graph[w] = append(graph[w], neighbor{v, i})
		edges[i] = edge{v, w}
	}

	isBridge = make([]bool, len(edges))
	dfn := make([]int, len(graph)) // 值从 1 开始
	dfsClock := 0
	var dfs func(int, int) int
	dfs = func(cur, fid int) int { // 使用 fid 而不是 fa，可以兼容重边的情况
		dfsClock++
		dfn[cur] = dfsClock
		curLow := dfsClock
		for _, edge := range graph[cur] {
			if next := edge.to; dfn[next] == 0 {
				nextLow := dfs(next, edge.eid)
				if nextLow > dfn[cur] { // 以 next 为根的子树中没有反向边能连回 cur 或 cur 的祖先，所以 cur-next 必定是桥
					isBridge[edge.eid] = true
				}
				curLow = min(curLow, nextLow)
			} else if edge.eid != fid { // 找到 cur 的反向边 cur-next，用 dfn[next] 来更新 curLow
				curLow = min(curLow, dfn[next])
			}
		}
		return curLow
	}

	for v, timestamp := range dfn {
		if timestamp == 0 {
			dfs(v, -1)
		}
	}

	// EXTRA: 所有桥边的下标
	bridgeEIDs := []int{}
	for eid, v := range isBridge {
		if v {
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
func findVertexBCC(graph [][]int, min func(int, int) int) (groups [][]int, vbccId []int) {
	vbccId = make([]int, len(graph)) // ID 从 1 开始编号
	idCount := 0
	isCut := make([]bool, len(graph))

	dfn := make([]int, len(graph)) // 值从 1 开始
	dfsClock := 0
	type edge struct{ v, w int } // eid
	stack := []edge{}
	var dfs func(cur, pre int) int
	dfs = func(cur, pre int) int {
		dfsClock++
		dfn[cur] = dfsClock
		curLow := dfsClock
		childCount := 0
		for _, next := range graph[cur] {
			e := edge{cur, next} // ne.eid
			if dfn[next] == 0 {
				stack = append(stack, e)
				childCount++
				nextLow := dfs(next, cur)
				if nextLow >= dfn[cur] {
					isCut[cur] = true
					idCount++
					group := []int{}
					//eids := []int{}
					for {
						e, stack = stack[len(stack)-1], stack[:len(stack)-1]
						if vbccId[e.v] != idCount {
							vbccId[e.v] = idCount
							group = append(group, e.v)
						}
						if vbccId[e.w] != idCount {
							vbccId[e.w] = idCount
							group = append(group, e.w)
						}
						//eids = append(eids, e.eid)
						if e.v == cur && e.w == next {
							break
						}
					}
					// 点数和边数相同，说明该 v-BCC 是一个简单环，且环上所有的边只属于一个简单环
					//if len(comp) == len(eids) {
					//	for _, eid := range eids {
					//		onSimpleCycle[eid] = true
					//	}
					//}
					groups = append(groups, group)
				}
				curLow = min(curLow, nextLow)
			} else if next != pre && dfn[next] < dfn[cur] {
				stack = append(stack, e)
				curLow = min(curLow, dfn[next])
			}
		}
		if pre == -1 && childCount == 1 {
			isCut[cur] = false
		}
		return curLow
	}

	for i, timestamp := range dfn {
		if timestamp == 0 {
			if len(graph[i]) == 0 { // 零度，即孤立点（isolated vertex）
				idCount++
				vbccId[i] = idCount
				groups = append(groups, []int{i})
				continue
			}
			dfs(i, -1)
		}
	}

	{
		// EXTRA: 缩点成树
		// !BCC 和割点作为新图中的节点，并在每个割点与包含它的所有 BCC 之间连边
		// !bcc1 - 割点1 - bcc2 - 割点2 - ...
		cutId := make([]int, len(graph))
		for i, v := range isCut {
			if v {
				idCount++ // !接在 BCC 之后给割点编号
				cutId[i] = idCount
			}
		}

		for vbcc, group := range groups {
			vbcc++
			for _, v := range group {
				if v = cutId[v]; v > 0 {
					// add(v,w); add(w,v) ...
				}
			}
		}
	}

	return
}

// e-BCC：删除无向图中所有的割边后，剩下的每一个 CC 都是 e-BCC
// 缩点后形成一颗 bridge tree
// 模板题 https://codeforces.com/problemset/problem/1000/E
// 较为综合的一道题 http://codeforces.com/problemset/problem/732/F
func findEdgeBCC(in io.Reader, n, m int) (groups [][]int, ebccId []int) {
	type neighbor struct{ to, eid int }
	type edge struct{ v, w int }
	graph := make([][]neighbor, n)
	edges := make([]edge, m)

	// *copy* 包含读图
	isBridge := findBridges(in, n, m)

	// 求原图中每个点的 bccID
	ebccId = make([]int, len(graph))
	idCount := 0
	var group []int
	var dfs2 func(int)
	dfs2 = func(cur int) {
		ebccId[cur] = idCount
		group = append(group, cur)
		for _, e := range graph[cur] {
			if next := e.to; !isBridge[e.eid] && ebccId[next] == 0 {
				dfs2(next)
			}
		}
	}

	for i, id := range ebccId {
		if id == 0 {
			idCount++
			group = []int{}
			dfs2(i)
			groups = append(groups, group)
		}
	}

	{
		// EXTRA: 缩点成树，复杂度 O(M)
		// !ebcc1 - ebcc2 - ebcc3 - ...
		graph2 := make([][]int, idCount)

		// 遍历 edges，若两端点的 bccIDs 不同（割边）则建边
		for _, e := range edges {
			if v, w := ebccId[e.v]-1, ebccId[e.w]-1; v != w {
				graph2[v] = append(graph2[v], w)
				graph2[w] = append(graph2[w], v)
			}
		}

		// 也可以遍历 isBridge，割边两端点 bccIDs 一定不同
		for eid, b := range isBridge {
			if b {
				e := edges[eid]
				v, w := ebccId[e.v]-1, ebccId[e.w]-1
				graph2[v] = append(graph2[v], w)
				graph2[w] = append(graph2[w], v)
			}
		}
	}

	return
}
