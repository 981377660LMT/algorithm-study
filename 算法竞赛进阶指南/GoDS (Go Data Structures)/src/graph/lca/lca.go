package lca

func LCA() {

}
func (*tree) lcaBinarySearch(n, root int, g [][]int, max func(int, int) int) {
	const mx = 17 // bits.Len(最大节点数)
	pa := make([][mx]int, n)
	dep := make([]int, n)
	var build func(v, p, d int)
	build = func(v, p, d int) {
		pa[v][0] = p
		dep[v] = d
		for _, w := range g[v] {
			if w != p {
				build(w, v, d+1)
			}
		}
	}
	build(root, -1, 0)
	// 倍增
	for i := 0; i+1 < mx; i++ {
		for v := range pa {
			if p := pa[v][i]; p != -1 {
				pa[v][i+1] = pa[p][i]
			} else {
				pa[v][i+1] = -1
			}
		}
	}
	// 从 v 开始向上跳 k 步，不存在返回 -1
	// O(1) 求法见长链剖分
	uptoKthPa := func(v, k int) int {
		for i := 0; i < mx && v != -1; i++ {
			if k>>i&1 > 0 {
				v = pa[v][i]
			}
		}
		return v
	}
	// 从 v 开始向上跳到指定深度 d，d<=dep[v]
	// https://en.wikipedia.org/wiki/Level_ancestor_problem
	// https://codeforces.com/problemset/problem/1535/E
	uptoDep := func(v, d int) int {
		for i := 0; i < mx; i++ {
			if (dep[v]-d)>>i&1 > 0 {
				v = pa[v][i]
				//if v == -1 { panic(-9) }
			}
		}
		return v
	}
	_lca := func(v, w int) int {
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
	disVW := func(v, w int) int { return dep[v] + dep[w] - dep[_lca(v, w)]<<1 }

	// EXTRA: 输入 u 和 v，u 是 v 的祖先，返回 u 到 v 路径上的第二个节点
	down := func(u, v int) int {
		// assert dep[u] < dep[v]
		v = uptoDep(v, dep[u]+1)
		if pa[v][0] == u {
			return v
		}
		return -1
	}

	{
		// 加权树上二分
		var dep []int64 // 加权深度，dfs 预处理略
		// 从 v 开始向根移动至多 d 距离，返回最大移动次数，以及能移动到的离根最近的点
		// 变形 https://codeforces.com/problemset/problem/932/D
		uptoDep := func(v int, d int64) (int, int) {
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
		type nb struct{ to, wt int }
		g := make([][]nb, n)
		// read g ...
		const mx = 18
		type pair struct{ p, maxWt int }
		pa := make([][mx]pair, n)
		dep := make([]int, n)
		var build func(v, p, d int)
		build = func(v, p, d int) {
			pa[v][0].p = p
			dep[v] = d
			for _, e := range g[v] {
				if w := e.to; w != p {
					pa[w][0].maxWt = e.wt
					build(w, v, d+1)
				}
			}
		}
		build(0, -1, 0)

		for i := 0; i+1 < mx; i++ {
			for v := range pa {
				if p := pa[v][i]; p.p != -1 {
					pp := pa[p.p][i]
					pa[v][i+1] = pair{pp.p, max(p.maxWt, pp.maxWt)}
				} else {
					pa[v][i+1].p = -1
				}
			}
		}

		// 求 LCA(v,w) 的同时，顺带求出 v-w 上的边权最值
		_lca := func(v, w int) (lca, maxWt int) {
			if dep[v] > dep[w] {
				v, w = w, v
			}
			for i := 0; i < mx; i++ {
				if (dep[w]-dep[v])>>i&1 > 0 {
					p := pa[w][i]
					maxWt = max(maxWt, p.maxWt)
					w = p.p
				}
			}
			if w != v {
				for i := mx - 1; i >= 0; i-- {
					if pv, pw := pa[v][i], pa[w][i]; pv.p != pw.p {
						maxWt = max(maxWt, max(pv.maxWt, pw.maxWt))
						v, w = pv.p, pw.p
					}
				}
				maxWt = max(maxWt, max(pa[v][0].maxWt, pa[w][0].maxWt))
				v = pa[v][0].p
			}
			// 如果是点权的话这里加上 maxWt = max(maxWt, pa[v][0].maxWt)
			lca = v
			return
		}

		_ = _lca
	}

	_ = []interface{}{disVW, uptoKthPa, down}
}

// 最近公共祖先 · 其二 · 基于 RMQ
// O(nlogn) 预处理，O(1) 查询
// 由于预处理 ST 表是基于一个长度为 2n 的序列，所以常数上是比倍增算法要大的。内存占用也比倍增要大一倍左右（这点可忽略）
// 优点是查询的复杂度低，适用于查询量大的情形
// https://oi-wiki.org/graph/lca/#rmq
func (*tree) lcaRMQ(n, root int, g [][]int) {
	vs := make([]int, 0, 2*n-1)  // 欧拉序列
	pos := make([]int, n)        // pos[v] 表示 v 在 vs 中第一次出现的位置编号
	dep := make([]int, 0, 2*n-1) // 深度序列，和欧拉序列一一对应
	disRoot := make([]int, n)    // disRoot[v] 表示 v 到 root 的距离
	var build func(v, p, d int)  // 若有边权需额外传参 dis
	build = func(v, p, d int) {
		pos[v] = len(vs)
		vs = append(vs, v)
		dep = append(dep, d)
		disRoot[v] = d
		for _, w := range g[v] {
			if w != p {
				build(w, v, d+1) // d+e.wt
				vs = append(vs, v)
				dep = append(dep, d)
			}
		}
	}
	build(root, -1, 0)
	type pair struct{ v, i int }
	const mx = 17 // bits.Len(最大节点数)
	var st [][mx]pair
	stInit := func(a []int) {
		n := len(a)
		st = make([][mx]pair, n)
		for i, v := range a {
			st[i][0] = pair{v, i}
		}
		for j := 1; 1<<j <= n; j++ {
			for i := 0; i+1<<j <= n; i++ {
				if a, b := st[i][j-1], st[i+1<<(j-1)][j-1]; a.v < b.v {
					st[i][j] = a
				} else {
					st[i][j] = b
				}
			}
		}
	}
	stInit(dep)
	stQuery := func(l, r int) int { // [l,r) 注意 l r 是从 0 开始算的
		k := bits.Len(uint(r-l)) - 1
		a, b := st[l][k], st[r-1<<k][k]
		if a.v < b.v {
			return a.i
		}
		return b.i
	}
	// 注意下标的换算，打印 LCA 的话要 +1
	_lca := func(v, w int) int {
		pv, pw := pos[v], pos[w]
		if pv > pw {
			pv, pw = pw, pv
		}
		return vs[stQuery(pv, pw+1)]
	}
	_d := func(v, w int) int { return disRoot[v] + disRoot[w] - disRoot[_lca(v, w)]<<1 }

	_ = _d
}

// 最近公共祖先 · 其三 · Tarjan 离线算法
// 时间和空间复杂度均为 O(n+q)
// 虽然用了并查集但是由于数据的特殊性，操作的均摊结果是 O(1) 的，见 https://core.ac.uk/download/pdf/82125836.pdf
// https://oi-wiki.org/graph/lca/#tarjan
// https://cp-algorithms.com/graph/lca_tarjan.html
// 扩展：Tarjan RMQ https://codeforces.com/blog/entry/48994
func (*tree) lcaTarjan(in io.Reader, n, q, root int) []int {
	g := make([][]int, n)
	for i := 1; i < n; i++ {
		v, w := 0, 0
		Fscan(in, &v, &w)
		v--
		w--
		g[v] = append(g[v], w)
		g[w] = append(g[w], v)
	}

	lca := make([]int, q)
	dis := make([]int, q) // dis(q.v,q.w)
	type query struct{ w, i int }
	qs := make([][]query, n)
	for i := 0; i < q; i++ {
		v, w := 0, 0
		Fscan(in, &v, &w)
		v--
		w--
		if v != w {
			qs[v] = append(qs[v], query{w, i})
			qs[w] = append(qs[w], query{v, i})
		} else {
			// do v==w...
			lca[i] = v
			dis[i] = 0
		}
	}

	pa := make([]int, n)
	for i := range pa {
		pa[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if pa[x] != x {
			pa[x] = find(pa[x])
		}
		return pa[x]
	}

	dep := make([]int, n)
	color := make([]int8, n)
	var _f func(v, d int)
	_f = func(v, d int) {
		dep[v] = d
		color[v] = 1
		for _, w := range g[v] {
			if color[w] == 0 {
				_f(w, d+1)
				pa[w] = v
			}
		}
		for _, q := range qs[v] {
			if w := q.w; color[w] == 2 {
				// do(v, w, lcaVW)...
				lcaVW := find(w)
				lca[q.i] = lcaVW
				dis[q.i] = dep[v] + dep[w] - dep[lcaVW]<<1
			}
		}
		color[v] = 2
	}
	_f(root, 0)
	return lca
}
