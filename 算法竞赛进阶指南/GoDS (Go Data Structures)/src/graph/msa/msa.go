// !https: //github.dev/EndlessCheng/codeforces-go/blob/f032745f91b1817db522b2590d97970d93be8358/copypasta/graph.go#L2004
// 最小树形图：有向图最小生成树

// 最小树形图 (MSA, Minimum weight Spanning Arborescence)   有向图上的最小生成树 (DMST)
// O(nm) 朱刘算法（Edmonds 算法）
// https://en.wikipedia.org/wiki/Edmonds%27_algorithm
// https://oi-wiki.org/graph/dmst/
// todo 另外还有 Tarjan 的 O(m+nlogn) 算法
//  https://oi-wiki.org/graph/dmst/#tarjan-dmst
//
// 模板题 https://www.luogu.com.cn/problem/P4716
package msa

const INF int = 2e9

func msaEdmonds(n, root int, edges [][3]int) (res int64) {
	minW := make([]int, n)
	fa := make([]int, n)
	id := make([]int, n)
	rt := make([]int, n)

	for {
		for i := range minW {
			minW[i] = INF
		}
		for _, e := range edges {
			if v, w, wt := e[0], e[1], e[2]; wt < minW[w] {
				minW[w] = wt
				fa[w] = v
			}
		}
		for i, wt := range minW {
			if i != root && wt == INF {
				return -1
			}
		}
		cid := 0
		for i := range id {
			id[i] = -1
			rt[i] = -1
		}
		for i, wt := range minW {
			if i == root {
				continue
			}
			res += int64(wt)
			v := i
			for ; v != root && id[v] < 0 && rt[v] != i; v = fa[v] {
				rt[v] = i
			}
			if v != root && id[v] < 0 { // rt[v] == i，有环
				id[v] = cid
				for x := fa[v]; x != v; x = fa[x] {
					id[x] = cid
				}
				cid++
			}
		}
		if cid == 0 {
			return
		}
		for i, v := range id {
			if v < 0 {
				id[i] = cid
				cid++
			}
		}

		// 缩点
		tmp := edges
		edges = nil
		for _, e := range tmp {
			if v, w := id[e[0]], id[e[1]]; v != w {
				edges = append(edges, [3]int{v, w, e[2] - minW[e[1]]})
			}
		}
		root = id[root]
		minW = minW[:cid]
		id = id[:cid]
	}
}
