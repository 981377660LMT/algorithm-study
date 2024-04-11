// https://github.com/981377660LMT/codeforces-go/blob/master/copypasta/graph_tree.go

// LCA 应用：树上差分
// 操作为更新 v-w 路径上的点权或边权（初始为 0）
// 点权时 diff[lca] -= val
// 边权时 diff[lca] -= 2 * val（定义 diff 为点到父亲的差分值）
// https://www.luogu.com.cn/blog/RPdreamer/ci-fen-and-shu-shang-ci-fen
// todo https://loj.ac/d/1698
// 模板题（边权）https://codeforces.com/problemset/problem/191/C
func (*tree) differenceInTree(in io.Reader, n, root int, g [][]int) []int {
	_lca := func(v, w int) (_ int) { return }

	diff := make([]int, n)
	update := func(v, w int, val int) {
		lca := _lca(v, w)
		diff[v] += val
		diff[w] += val
		diff[lca] -= val // 点权
		//diff[lca] -= 2 * val // 边权
	}
	var q int
	Fscan(in, &q)
	for i := 0; i < q; i++ {
		var v, w, val int
		Fscan(in, &v, &w, &val)
		v--
		w--
		update(v, w, val)
	}

	// 自底向上求出每个点的点权/边权
	ans := make([]int, n)
	var f func(v, fa int) int
	f = func(v, fa int) int {
		sum := diff[v]
		for _, w := range g[v] {
			if w != fa {
				// 边权的话在这里记录 ans
				//s := f(w, v)
				//ans[e.eid] = s
				//sum += s
				sum += f(w, v)
			}
		}
		// 点权的话在这里记录 ans
		ans[v] = sum
		return sum
	}
	f(root, -1)

	return ans
}

// type E = int32

// func e() E          { return 0 }
// func op(e1, e2 E) E { return e1 + e2 }
// func inv(e E) E     { return -e }

// type RangeAddPointGetTreeOfflineEdge struct {
// 	tree [][]int32
// 	root int32
// 	lca  *LCAHLD
// 	diff []E
// }

// // 树上差分离线版.区间加,单点查询.
// func NewRangeAddPointGetTreeOfflineEdge(tree [][]int32, root int32) *RangeAddPointGetTreeOfflineEdge {
// 	n := len(tree)
// 	diff := make([]E, n)
// 	for i := 0; i < n; i++ {
// 		diff[i] = e()
// 	}
// 	lca := NewLCA(tree, root)
// 	return &RangeAddPointGetTreeOfflineEdge{
// 		tree: tree,
// 		root: root,
// 		lca:  lca,
// 		diff: diff,
// 	}
// }

// // 路径上所有点加上delta.
// // sum[u]++, sum[v]++, sum[lca]-=2
// func (r *RangeAddPointGetTreeOfflineEdge) AddPath(u, v int32, delta E) {
// 	r.diff[u] = op(r.diff[u], delta)
// 	r.diff[v] = op(r.diff[v], delta)
// 	lca := r.lca.LCA(u, v)
// 	r.diff[lca] = op(r.diff[lca], op(inv(delta), inv(delta)))
// }

// // O(n)构建，求出每个结点与父亲结点的边权.
// func (r *RangeAddPointGetTreeOfflineEdge) GetAll() []E {
// 	res := make([]E, len(r.diff))
// 	parent := r.lca.Parent

// 	var dfs func(int32) E
// 	dfs = func(cur int32) E {
// 		sum := r.diff[cur]
// 		for _, next := range r.tree[cur] {
// 			if next != parent[cur] {
// 				nextSum := dfs(next)
// 				res[next] = nextSum
// 				sum = op(sum, nextSum)
// 			}
// 		}
// 		return sum
// 	}
// 	dfs(r.root)

// 	return res
// }
