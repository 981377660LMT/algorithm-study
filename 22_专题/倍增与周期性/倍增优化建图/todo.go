
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
	getDis := func(v, w int) int { return dep[v] + dep[w] - dep[getLCA(v, w)]*2 }

	// EXTRA: 输入 v 和 to，to 可能是 v 的子孙，返回从 v 到 to 路径上的第二个节点（v 的一个儿子）
	// 如果 to 不是 v 的子孙，返回 -1
	// https://codeforces.com/problemset/problem/916/E
	// https://codeforces.com/problemset/problem/1702/G2
	down1 := func(v, to int) int {
		if dep[to] <= dep[v] {
			return -1
		}
		to = uptoDep(to, dep[v]+1)
		if pa[to][0] == v {
			return to
		}
		return -1
	}

	// EXTRA: 从 v 出发，向 to 方向走一步
	// 输入需要保证 v != to
	move1 := func(v, to int) int {
		if v == to {
			panic(-1)
		}
		if getLCA(v, to) == v { // to 在 v 下面
			return uptoDep(to, dep[v]+1)
		}
		// lca 在 v 上面
		return pa[v][0]
	}

	// EXTRA: 从 v 开始，向上跳 k 步
	// 不存在则返回 -1
	// O(1) 求法见长链剖分
	uptoKthPa := func(v, k int) int {
		for ; k > 0 && v != -1; k &= k - 1 {
			v = pa[v][bits.TrailingZeros(uint(k))]
		}
		return v
	}

	// EXTRA: 输入 v 和 w，返回 v 到 w 路径上的中点
	// 返回值是一个数组，因为可能有两个中点
	// 在有两个中点的情况下，保证返回值的第一个中点离 v 更近
	midPath := func(v, w int) []int {
		lca := getLCA(v, w)
		dv := dep[v] - dep[lca]
		dw := dep[w] - dep[lca]
		if dv == dw {
			return []int{lca}
		}
		if dv > dw {
			mid := uptoKthPa(v, (dv+dw)/2)
			if (dv+dw)%2 == 0 {
				return []int{mid}
			}
			return []int{mid, pa[mid][0]}
		} else {
			mid := uptoKthPa(w, (dv+dw)/2)
			if (dv+dw)%2 == 0 {
				return []int{mid}
			}
			return []int{pa[mid][0], mid} // pa[mid][0] 离 v 更近
		}
	}

	{
		// 加权树上二分
		var dep []int // 加权深度，dfs 预处理略
		// 从 v 开始向根移动至多 d 距离，返回最大移动次数，以及能移动到的离根最近的点
		// NOIP2012·提高 疫情控制 https://www.luogu.com.cn/problem/P1084
		// 变形 https://codeforces.com/problemset/problem/932/D
		// 点权写法 https://codeforces.com/problemset/problem/1059/E 2400
		uptoDep := func(v, d int) (int, int) {
			step := 0
			dv := dep[v]
			for i := mx - 1; i >= 0; i-- {
				if p := pa[v][i]; p != -1 && dv-dep[p] <= d {
					step |= 1 << i
					v = p
				}
			}
			return step, v
		}
		_ = uptoDep
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