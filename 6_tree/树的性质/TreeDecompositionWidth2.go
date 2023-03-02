// https://ei1333.hateblo.jp/entry/2020/02/12/150319
// https://ei1333.github.io/library/graph/others/tree-decomposition-width-2.hpp
// TreeDecompositionWidth2
// !無向グラフが与えられたとき, 木幅が 2以下か判定し, 2以下の場合は木幅 2 以下の木分解を構成する.
// 非連結だとバグります！(適当にダミー辺を加えて連結にしてね)

package main

func main() {

}

type DecompNode struct{ bag, child []int }
type TreeDecompositionWidth2 struct{ g [][]int }

func NewTreeDecompositionWidth2(n int) *TreeDecompositionWidth2 {
	return &TreeDecompositionWidth2{make([][]int, n)}
}

func (t *TreeDecompositionWidth2) AddEdge(a, b int) {
	t.g[a] = append(t.g[a], b)
	t.g[b] = append(t.g[b], a)
}

func (t *TreeDecompositionWidth2) Build() []DecompNode {
	n := len(t.g)
	used, deg := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		used[i] = -1
	}
	que := []int{}
	for i := 0; i < n; i++ {
		deg[i] = len(t.g[i])
		if deg[i] <= 2 {
			que = append(que, i)
		}
	}
	exists := make([]map[int]struct{}, n)
	for i := 0; i < n; i++ {
		exists[i] = map[int]struct{}{}
		for _, j := range t.g[i] {
			if i < j {
				exists[i][j] = struct{}{}
			}
		}
	}

	res := []DecompNode{}
	res = append(res, DecompNode{})
	for len(que) > 0 {
		idx := que[0]
		que = que[1:]
		if deg[idx] > 2 || used[idx] != -1 {
			continue
		}
		var r DecompNode
		r.bag = append(r.bag, idx)
		u, v := -1, -1
		for _, to := range t.g[idx] {
			if used[to] == -1 {
				if u == -1 {
					u = to
				} else {
					v = to
				}
				r.bag = append(r.bag, to)
			} else if used[to] >= 0 {
				r.child = append(r.child, to)
				used[to] = -2
			}
		}

		if u == -1 {
		} else if v == -1 {
			deg[u]--
		} else {
			if u > v {
				u, v = v, u
			}
			if _, ok := exists[u][v]; !ok {
				t.g[u] = append(t.g[u], v)
				t.g[v] = append(t.g[v], u)
				exists[u][v] = struct{}{}
			} else {
				deg[u]--
				deg[v]--
			}
		}

		for _, to := range t.g[idx] {
			if deg[to] <= 2 {
				que = append(que, to)
			}
		}

		used[idx] = len(res)
		deg[idx] = 0
		res = append(res, r)
	}

	for i := 0; i < n; i++ {
		if deg[i] > 0 {
			return nil
		}
	}
	res[0] = res[len(res)-1]
	res = res[:len(res)-1]
	return res
}
