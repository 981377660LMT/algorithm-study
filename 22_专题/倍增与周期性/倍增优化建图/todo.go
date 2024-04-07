
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/graph_tree.go#L2140

// 最近公共祖先 · 其一 · 基于树上倍增
// 【模板讲解】树上倍增算法（以及最近公共祖先）
// - 请看 https://leetcode.cn/problems/kth-ancestor-of-a-tree-node/solution/mo-ban-jiang-jie-shu-shang-bei-zeng-suan-v3rw/
// O(nlogn) 预处理，O(logn) 查询
// 适用于查询量和节点数等同的情形
// 适用于可以动态添加节点（挂叶子）的情形
// NOTE: 多个点的 LCA 等于 dfn_min 和 dfn_max 的 LCA
// https://oi-wiki.org/graph/lca/#_5
// 另见 mo.go 中的【树上莫队】
//
// 倍增 LC1483 https://leetcode.cn/problems/kth-ancestor-of-a-tree-node/
// 模板题 https://www.luogu.com.cn/problem/P3379
// https://codeforces.com/problemset/problem/33/D 2000
// https://codeforces.com/problemset/problem/1304/E 2000
// 到两点距离相同的点的数量 https://codeforces.com/problemset/problem/519/E 2100
// https://codeforces.com/problemset/problem/916/E 2400
// https://atcoder.jp/contests/arc060/tasks/arc060_c
// 路径点权乘积 https://ac.nowcoder.com/acm/contest/6913/C
//
// 维护元素和 LC2836 https://leetcode.cn/problems/maximize-value-of-function-in-a-ball-passing-game/
// 维护边权出现次数 LC2846 https://leetcode.cn/problems/minimum-edge-weight-equilibrium-queries-in-a-tree/
// 维护最大值（与 MST 结合）https://codeforces.com/problemset/problem/609/E
//    变体 https://codeforces.com/problemset/problem/733/F
// 维护最大值（与 MST 结合）LC1697 https://leetcode.cn/problems/checking-existence-of-edge-length-limited-paths/
// 维护最大值（与 MST 结合）LC1724（上面这题的在线版）https://leetcode.cn/problems/checking-existence-of-edge-length-limited-paths-ii/
// 维护最大值和严格次大值（严格次小 MST）：见 graph.go 中的 strictlySecondMST
// 维护前十大（点权）https://codeforces.com/problemset/problem/587/C
// 维护最大子段和 https://codeforces.com/contest/1843/problem/F2
// 维护从 x 往上有几个不同的 OR https://codeforces.com/contest/1878/problem/G
// 维护最大值 https://www.hackerearth.com/practice/algorithms/graphs/graph-representation/practice-problems/algorithm/optimal-connectivity-c6ae79ca/
// http://acm.hdu.edu.cn/showproblem.php?pid=7345
//
// 树上倍增-查询深度最小的未被标记的点 https://codeforces.com/problemset/problem/980/E
// 题目推荐 https://cp-algorithms.com/graph/lca.html#toc-tgt-2
// todo poj2763 poj1986 poj3728
func (*tree) lcaBinaryLifting(root int, g [][]int) {
	const mx = 17 // bits.Len(最大节点数)
	pa := make([][mx]int, len(g))
	dep := make([]int, len(g)) // 根节点的深度为 0
	var buildPa func(int, int)
	buildPa = func(v, p int) {
		pa[v][0] = p
		for _, w := range g[v] {
			if w != p {
				dep[w] = dep[v] + 1
				buildPa(w, v)
			}
		}
	}
	buildPa(root, -1)
	for i := 0; i+1 < mx; i++ {
		for v := range pa {
			if p := pa[v][i]; p != -1 {
				pa[v][i+1] = pa[p][i]
			} else {
				pa[v][i+1] = -1
			}
		}
	}
	// 从 v 开始，向上跳到指定深度 d
	// https://en.wikipedia.org/wiki/Level_ancestor_problem
	// https://codeforces.com/problemset/problem/1535/E
	uptoDep := func(v, d int) int {
		if d > dep[v] {
			panic(-1)
		}
		for k := uint(dep[v] - d); k > 0; k &= k - 1 {
			v = pa[v][bits.TrailingZeros(k)]
		}
		return v
	}
	getLCA := func(v, w int) int {
		if dep[v] > dep[w] {
			v, w = w, v
		}
		w = uptoDep(w, dep[v])
		if w == v {
			return v
		}
		for i := mx - 1; i >= 0; i-- {
			if pv, pw := pa[v][i], pa[w][i]; pv != pw {
				v, w = pv, pw
			}
		}
		return pa[v][0]
	}

	{
		// EXTRA: 倍增的时候维护其他属性，如边权最值等
		// 下面的代码来自 https://codeforces.com/problemset/problem/609/E
		// EXTRA: 额外维护最值边的下标，见 https://codeforces.com/contest/733/submission/120955685
		// 点权写法 https://codeforces.com/problemset/problem/1059/E 2400
		type nb struct{ to, wt int }
		var g [][]nb // read g ...

		const mx = 18
		type data int
		type pair struct {
			p     int
			maxWt data
		}
		pa := make([][mx]pair, len(g))
		dep := make([]int, len(g))
		var build func(v, p, d int)
		build = func(v, p, d int) {
			// 注意这里添加边权幺元
			pa[v][0].p = p
			dep[v] = d
			for _, e := range g[v] {
				if w := e.to; w != p {
					pa[w][0].maxWt = data(e.wt)
					build(w, v, d+1)
				}
			}
		}
		build(0, -1, 0)

		merge := func(a, b data) data {
			return data(max(int(a), int(b)))
		}

		for i := 0; i+1 < mx; i++ {
			for v := range pa {
				if p := pa[v][i]; p.p != -1 {
					pp := pa[p.p][i]
					pa[v][i+1] = pair{pp.p, merge(p.maxWt, pp.maxWt)}
				} else {
					pa[v][i+1].p = -1
				}
			}
		}

		// 求 LCA(v,w) 的同时，顺带求出 v-w 上的边权最值
		getLCA := func(v, w int) (lca int, maxWt data) {
			//pathLen := dep[v] + dep[w]
			if dep[v] > dep[w] {
				v, w = w, v
			}
			for k := dep[w] - dep[v]; k > 0; k &= k - 1 {
				p := pa[w][bits.TrailingZeros(uint(k))]
				maxWt = merge(maxWt, p.maxWt)
				w = p.p
			}
			if w != v {
				for i := mx - 1; i >= 0; i-- {
					if pv, pw := pa[v][i], pa[w][i]; pv.p != pw.p {
						maxWt = merge(maxWt, merge(pv.maxWt, pw.maxWt))
						v, w = pv.p, pw.p
					}
				}
				maxWt = merge(maxWt, merge(pa[v][0].maxWt, pa[w][0].maxWt))
				v = pa[v][0].p
			}
			// 如果是点权的话这里加上 maxWt = merge(maxWt, pa[v][0].maxWt)
			lca = v
			//pathLen -= dep[lca] * 2
			return
		}

		_ = getLCA
	}

	_ = []interface{}{getDis, uptoKthPa, down1, move1, midPath}
}