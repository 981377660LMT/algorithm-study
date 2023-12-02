// 树上等高线汇集
// https://maspypy.com/%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%83%bb1-3%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%81%ae%e3%81%8a%e7%b5%b5%e6%8f%8f%e3%81%8d
package main

import "fmt"

func main() {
	g := [][]int{{1, 2}, {0, 3}, {0, 4}, {1}, {2}}
	centroidDecomposition(len(g), g, func(par, vs, color []int) {
		// par: 重心分解の親
		// vs: 重心分解の頂点集合
		// color: 重心分解の色
		fmt.Println(par, vs, color)
	})

}

func centroidDecomposition(n int, g [][]int, f func([]int, []int, []int)) {
	if n == 1 {
		return
	}
	V := make([]int, n)
	par := make([]int, n)
	for i := range par {
		par[i] = -1
	}
	l := 0
	r := 0
	V[r] = 0
	r++
	for l < r {
		v := V[l]
		l++
		for _, next := range g[v] {
			if next != par[v] {
				V[r] = next
				par[next] = v
				r++
			}
		}
	}
	if r != n {
		panic("r should be equal to n")
	}
	newIdx := make([]int, n)
	for i := 0; i < n; i++ {
		newIdx[V[i]] = i
	}
	tmp := make([]int, n)
	for i := 0; i < n; i++ {
		tmp[i] = -1
	}
	for i := 1; i < n; i++ {
		j := par[i]
		tmp[newIdx[i]] = newIdx[j]
	}
	par = tmp

	real := make([]int, n)
	for i := range real {
		real[i] = 1
	}
	centroidDecomposition2DFS(par, V, real, f)
}

func centroidDecomposition2DFS(par []int, vs []int, real []int, f func([]int, []int, []int)) {
	N := len(par)
	if N <= 1 {
		panic("N should be greater than or equal to 2")
	}
	if N == 2 {
		if real[0] != 0 && real[1] != 0 {
			color := []int{0, 1}
			f(par, vs, color)
		}
		return
	}
	c := -1
	sz := make([]int, N)
	for i := range sz {
		sz[i] = 1
	}
	for i := N - 1; i >= 0; i-- {
		if sz[i] >= (N+1)>>1 {
			c = i
			break
		}
		sz[par[i]] += sz[i]
	}
	color := make([]int, N)
	for i := range color {
		color[i] = -1
	}
	take := 0
	ord := make([]int, N)
	for i := range ord {
		ord[i] = -1
	}
	ord[c] = 0
	p := 1
	for v := 1; v < N; v++ {
		if par[v] == c && take+sz[v] <= (N-1)/2 {
			color[v] = 0
			ord[v] = p
			p++
			take += sz[v]
		}
	}
	for i := 1; i < N; i++ {
		if color[par[i]] == 0 {
			color[i] = 0
			ord[i] = p
			p++
		}
	}
	n0 := p - 1
	for a := par[c]; a != -1; a = par[a] {
		color[a] = 1
		ord[a] = p
		p++
	}
	for i := 0; i < N; i++ {
		if i != c && color[i] == -1 {
			color[i] = 1
			ord[i] = p
			p++
		}
	}
	if p != N {
		panic("p should be equal to N")
	}
	n1 := N - 1 - n0
	par0 := make([]int, n0+1)
	for i := range par0 {
		par0[i] = -1
	}
	par1 := make([]int, n1+1)
	for i := range par1 {
		par1[i] = -1
	}
	par2 := make([]int, N)
	for i := range par2 {
		par2[i] = -1
	}
	V0 := make([]int, n0+1)
	V1 := make([]int, n1+1)
	V2 := make([]int, N)
	rea0 := make([]int, n0+1)
	rea1 := make([]int, n1+1)
	rea2 := make([]int, N)
	for v := 0; v < N; v++ {
		i := ord[v]
		V2[i] = vs[v]
		rea2[i] = real[v]
		if color[v] != 1 {
			V0[i] = vs[v]
			rea0[i] = real[v]
		}
		if color[v] != 0 {
			V1[max(i-n0, 0)] = vs[v]
			rea1[max(i-n0, 0)] = real[v]
		}
	}
	for v := 1; v < N; v++ {
		a := ord[v]
		b := ord[par[v]]
		if a > b {
			a, b = b, a
		}
		par2[b] = a
		if color[v] != 1 && color[par[v]] != 1 {
			par0[b] = a
		}
		if color[v] != 0 && color[par[v]] != 0 {
			par1[max(b-n0, 0)] = max(a-n0, 0)
		}
	}
	if real[c] != 0 {
		color = make([]int, N)
		for i := 0; i < N; i++ {
			color[i] = -1
		}
		color[0] = 0
		for i := 1; i < N; i++ {
			if rea2[i] != 0 {
				color[i] = 1
			} else {
				color[i] = -1
			}
		}
		f(par2, V2, color)
		rea0[0] = 0
		rea1[0] = 0
		rea2[0] = 0
	}
	color = make([]int, N)
	for i := 0; i < N; i++ {
		color[i] = -1
	}
	for i := 1; i < N; i++ {
		if rea2[i] != 0 {
			if i <= n0 {
				color[i] = 0
			} else {
				color[i] = 1
			}
		}
	}
	f(par2, V2, color)
	centroidDecomposition2DFS(par0, V0, rea0, f)
	centroidDecomposition2DFS(par1, V1, rea1, f)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
