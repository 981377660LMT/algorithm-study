// 换根dp

package main

// https://leetcode.cn/problems/difference-between-maximum-and-minimum-price-sum/
// !求每个点作为根节点时，到叶子节点的最大点权和(不包括自身)
func maxOutput(n int, edges [][]int, price []int) int64 {
	e := func(root int) E {
		return 0
	}
	op := func(child1, child2 E) E {
		return max(child1, child2)
	}
	composition := func(fromRes E, parent, cur int, direction uint8) E {
		if direction == 0 {
			return fromRes + price[cur] // child -> parent
		}
		return fromRes + price[parent] // parent -> child
	}

	getSubSize := func() []int {
		subSize := make([]int, n)
		var dfs func(cur, parent int) int
		dfs = func(cur, parent int) int {
			res := 1
			for _, next := range edges[cur] {
				if next != parent {
					res += dfs(next, cur)
				}
			}
			subSize[cur] = res
			return res
		}
		dfs(0, -1)
		return subSize
	}

	R := NewRerooting(n)
	for _, edge := range edges {
		R.AddEdge(edge[0], edge[1])
	}

	_ = getSubSize
	dp := R.ReRooting(e, op, composition)
	return int64(maxs(dp...))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

//
//
type E = int
type Rerooting struct {
	Tree [][]int
	n    int
}

func NewRerooting(n int) *Rerooting {
	return &Rerooting{Tree: make([][]int, n), n: n}
}

// 添加一条无向边.
func (r *Rerooting) AddEdge(u, v int) {
	r.Tree[u] = append(r.Tree[u], v)
	r.Tree[v] = append(r.Tree[v], u)
}

func (r *Rerooting) ReRooting(e func(root int) E, op func(child1, child2 E) E, composition func(fromRes E, parent, cur int, direction uint8) E) []E {
	parents := make([]int, r.n)
	for i := range parents {
		parents[i] = -1
	}
	order := []int{0} // root = 0
	stack := []int{0}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, next := range r.Tree[cur] {
			if next != parents[cur] {
				parents[next] = cur
				order = append(order, next)
				stack = append(stack, next)
			}
		}
	}

	dp1, dp2 := make([]E, r.n), make([]E, r.n)
	for i := range dp1 {
		dp1[i] = e(i)
		dp2[i] = e(i)
	}
	for i := r.n - 1; i >= 0; i-- {
		cur := order[i]
		res := e(cur)
		for _, next := range r.Tree[cur] {
			if next != parents[cur] {
				dp2[next] = res
				res = op(res, composition(dp1[next], cur, next, 0))
			}
		}

		res = e(cur)
		for j := len(r.Tree[cur]) - 1; j >= 0; j-- {
			next := r.Tree[cur][j]
			if next != parents[cur] {
				dp2[next] = op(res, dp2[next])
				res = op(res, composition(dp1[next], cur, next, 0))
			}
		}

		dp1[cur] = res
	}

	for i := 1; i < r.n; i++ {
		newRoot := order[i]
		parent := parents[newRoot]
		dp2[newRoot] = composition(op(dp2[newRoot], dp2[parent]), parent, newRoot, 1)
		dp1[newRoot] = op(dp1[newRoot], dp2[newRoot])
	}

	return dp1
}
