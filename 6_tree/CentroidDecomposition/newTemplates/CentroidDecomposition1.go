package main

func main() {

}

// 1/3重心分解(1/3 Centroid Decomposition)
//
//	f(parent, vertex, n1, n2):
//	 [1,1+n1]: color 1
//	 [1+n1,1+n1+n2]: color 2
//
// https://maspypy.com/%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%83%bb1-3%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%81%ae%e3%81%8a%e7%b5%b5%e6%8f%8f%e3%81%8d
// https://codeforces.com/blog/entry/104997
func CentroidDecomposition1(
	n int32, g [][]int32,
	f func(parent []int32, vertex []int32, n1, n2 int32),
) {
	if n == 1 {
		return
	}
	queue := make([]int32, n)
	parent := make([]int32, n)
	for i := range parent {
		parent[i] = -1
	}
	l := int32(0)
	r := int32(0)
	queue[r] = int32(0)
	r++
	for l < r {
		v := queue[l]
		l++
		for _, next := range g[v] {
			if next != parent[v] {
				queue[r] = next
				parent[next] = v
				r++
			}
		}
	}
	if r != n {
		panic("r should be equal to n")
	}
	newIdx := make([]int32, n)
	for i := int32(0); i < n; i++ {
		newIdx[queue[i]] = i
	}
	tmp := make([]int32, n)
	for i := int32(0); i < n; i++ {
		tmp[i] = -1
	}
	for i := int32(1); i < n; i++ {
		j := parent[i]
		tmp[newIdx[i]] = newIdx[j]
	}
	parent = tmp
	centroidDecomposition1Dfs(parent, queue, f)
}

func centroidDecomposition1Dfs(
	parent []int32, vs []int32,
	f func(parent []int32, vertex []int32, n1, n2 int32),
) {
	n := int32(len(parent))
	if n < 2 {
		panic("n must be at least 2")
	}
	if n == 2 {
		return
	}
	c := int32(-1)
	size := make([]int32, n)
	for i := range size {
		size[i] = 1
	}
	for i := n - 1; i >= 0; i-- {
		if size[i] >= (n+1)>>1 {
			c = i
			break
		}
		size[parent[i]] += size[i]
	}
	color := make([]int32, n)
	for i := range color {
		color[i] = -1
	}
	ord := append(color[:0:0], color...)
	take := int32(0)
	ord[c] = 0
	p := int32(1)
	for v := int32(1); v < n; v++ {
		if parent[v] == c && take+size[v] <= (n-1)>>1 {
			color[v] = 0
			ord[v] = p
			p++
			take += size[v]
		}
	}
	for i := int32(1); i < n; i++ {
		if color[parent[i]] == 0 {
			color[i] = 0
			ord[i] = p
			p++
		}
	}
	n0 := p - 1
	for a := parent[c]; a != -1; a = parent[a] {
		color[a] = 1
		ord[a] = p
		p++
	}
	for i := int32(0); i < n; i++ {
		if i != c && color[i] == -1 {
			color[i] = 1
			ord[i] = p
			p++
		}
	}
	n1 := n - 1 - n0
	par0 := make([]int32, n0+1)
	for i := range par0 {
		par0[i] = -1
	}
	par1 := make([]int32, n1+1)
	for i := range par1 {
		par1[i] = -1
	}
	par2 := make([]int32, n)
	for i := range par2 {
		par2[i] = -1
	}
	V0 := make([]int32, n0+1)
	V1 := make([]int32, n1+1)
	V2 := make([]int32, n)
	for v := int32(0); v < n; v++ {
		i := ord[v]
		V2[i] = vs[v]
		if color[v] != 1 {
			V0[i] = vs[v]
		}
		if color[v] != 0 {
			V1[max32(i-n0, 0)] = vs[v]
		}
	}
	for v := int32(1); v < n; v++ {
		a, b := ord[v], ord[parent[v]]
		if a > b {
			a, b = b, a
		}
		par2[b] = a
		if color[v] != 1 && color[parent[v]] != 1 {
			par0[b] = a
		}
		if color[v] != 0 && color[parent[v]] != 0 {
			par1[max32(b-n0, 0)] = max32(a-n0, 0)
		}
	}
	f(par2, V2, n0, n1)
	centroidDecomposition1Dfs(par0, V0, f)
	centroidDecomposition1Dfs(par1, V1, f)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
