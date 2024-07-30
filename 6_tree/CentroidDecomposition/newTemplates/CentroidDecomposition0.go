package main

func main() {

}

// 基于顶点的重心分解
// f(parent, vertex, color):
// https://maspypy.com/%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%83%bb1-3%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%81%ae%e3%81%8a%e7%b5%b5%e6%8f%8f%e3%81%8d
func CentroidDecomposition0(
	n int32, g [][]int32,
	f func(parent []int32, vertex []int32, indptr []int32),
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
	centroidDecomposition0Dfs(parent, queue, f)
}

func centroidDecomposition0Dfs(
	parent []int32, vs []int32,
	f func(parent []int32, vertex []int32, indptr []int32),
) {
	n := int32(len(parent))
	if n < 1 {
		panic("n must be at least 1")
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
	vertex := []int32{c}
	nc := int32(1)
	for v := int32(1); v < n; v++ {
		if parent[v] == c {
			vertex = append(vertex, v)
			color[v] = nc
			nc++
		}
	}
	if c > 0 {
		for a := parent[c]; a != -1; a = parent[a] {
			color[a] = nc
			vertex = append(vertex, a)
		}
		nc++
	}
	for i := int32(0); i < n; i++ {
		if i != c && color[i] == 0 {
			color[i] = color[parent[i]]
			vertex = append(vertex, i)
		}
	}
	indptr := make([]int32, nc+1)
	for i := int32(0); i < n; i++ {
		indptr[1+color[i]]++
	}
	for i := int32(0); i < nc; i++ {
		indptr[i+1] += indptr[i]
	}
	counter := append(indptr[:0:0], indptr...)
	ord := make([]int32, n)
	for _, v := range vertex {
		ord[counter[color[v]]] = v
		counter[color[v]]++
	}
	newIdx := make([]int32, n)
	for i := int32(0); i < n; i++ {
		newIdx[ord[i]] = i
	}
	name := make([]int32, n)
	for i := int32(0); i < n; i++ {
		name[newIdx[i]] = vs[i]
	}
	{
		tmp := make([]int32, n)
		for i := range tmp {
			tmp[i] = -1
		}
		for i := int32(1); i < n; i++ {
			a, b := newIdx[i], newIdx[parent[i]]
			if a > b {
				a, b = b, a
			}
			tmp[b] = a
		}
		parent, tmp = tmp, parent
	}
	f(parent, name, indptr)
	for k := int32(1); k < nc; k++ {
		left, right := indptr[k], indptr[k+1]
		par1 := make([]int32, right-left)
		for i := range par1 {
			par1[i] = -1
		}
		name1 := append(par1[:0:0], par1...)
		name1[0] = name[0]
		for i := left; i < right; i++ {
			name1[i-left] = name[i]
		}
		for i := left; i < right; i++ {
			par1[i-left] = max32(parent[i]-left, -1)
		}
		centroidDecomposition0Dfs(par1, name1, f)
	}
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
