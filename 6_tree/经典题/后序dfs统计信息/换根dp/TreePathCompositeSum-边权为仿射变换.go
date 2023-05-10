// https://judge.yosupo.jp/problem/tree_path_composite_sum
// Tree Path Composite Sum
// 求每个点到其他点的距离之和，边权为仿射变换函数

package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	affine := make([]map[int][2]int, n)
	for i := range affine {
		affine[i] = make(map[int][2]int)
	}

	R := NewRerooting(n)
	for i := 0; i < n-1; i++ {
		var u, v, mul, add int
		fmt.Fscan(in, &u, &v, &mul, &add)
		R.AddEdge(u, v)
		affine[u][v] = [2]int{mul, add}
		affine[v][u] = [2]int{mul, add}
	}

	e := func(root int) E { return E{} }

	op := func(child1, child2 E) E {
		sum1, count1 := child1.sum, child1.count
		sum2, count2 := child2.sum, child2.count
		return E{(sum1 + sum2) % MOD, count1 + count2}
	}

	// dir: 0: cur -> parent, 1: parent -> cur
	composition := func(fromRes E, parent, cur int, dir uint8) E {
		from, to := parent, cur
		if dir == 0 {
			from, to = to, from
		}

		edge := affine[from][to]
		mul, add := edge[0], edge[1]
		sum, count := fromRes.sum, fromRes.count
		resSum := ((sum+values[from])*mul%MOD + (count+1)*add%MOD) % MOD
		resCount := count + 1
		return E{sum: resSum, count: resCount}
	}

	dp := R.ReRooting(e, op, composition)
	for i := 0; i < n; i++ {
		sum := (dp[i].sum + values[i]) % MOD
		fmt.Fprint(out, sum, " ")
	}
}

//
//
type E = struct{ sum, count int }
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
