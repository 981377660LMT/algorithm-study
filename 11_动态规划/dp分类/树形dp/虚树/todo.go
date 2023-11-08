package main

import "sort"

func main() {

}

// LCA+DFN：虚树 Virtual Tree / Auxiliary Tree
// https://oi-wiki.org/graph/virtual-tree/ 栈相比两次排序，效率更高
// https://www.luogu.com.cn/blog/SSerxhs/qian-tan-xu-shu
// https://www.luogu.com.cn/problem/P5891 https://class.luogu.com.cn/classroom/lgr66
// https://codeforces.com/problemset/problem/613/D
// https://www.luogu.com.cn/problem/P4103
// https://www.luogu.com.cn/problem/P7409
// https://www.luogu.com.cn/problem/P3233
// https://www.luogu.com.cn/problem/P2495
func (*tree) virtualTree(g [][]int) {
	dfn := make([]int, len(g))
	ts := 0
	_ = ts
	// 向上查找<lcaBinarySearch>
	// buildPa 开头添加：
	// dfn[v] = ts; ts++
	var _lca func(int, int) int

	vt := make([][]int, len(g))
	const root = 0
	st := []int{root} // 用根节点作为栈底哨兵
	makeVT := func(nodes []int) [][]int {
		sort.Slice(nodes, func(i, j int) bool { return dfn[nodes[i]] < dfn[nodes[j]] })
		vt[root] = vt[root][:0]
		st = st[:1]
		for _, v := range nodes {
			if v == root {
				continue
			}
			// ... 某些题目需要判断 v 和 pa[v][0] 是否都在 nodes 中
			vt[v] = vt[v][:0]
			lca := _lca(st[len(st)-1], v)
			if lca != st[len(st)-1] {
				// 回溯
				for dfn[st[len(st)-2]] > dfn[lca] {
					top := st[len(st)-1]
					st = st[:len(st)-1]
					p := st[len(st)-1]
					vt[p] = append(vt[p], top)
				}
				if lca != st[len(st)-2] { // lca 不在栈中（首次遇到）
					vt[lca] = vt[lca][:0]
					vt[lca] = append(vt[lca], st[len(st)-1])
					st[len(st)-1] = lca // 加到栈中
				} else { // lca 已经在栈中
					vt[lca] = append(vt[lca], st[len(st)-1])
					st = st[:len(st)-1]
				}
			}
			st = append(st, v)
		}
		// 最后的回溯
		for i := 1; i < len(st); i++ {
			vt[st[i-1]] = append(vt[st[i-1]], st[i])
		}
		return vt
	}

	_ = makeVT
}
