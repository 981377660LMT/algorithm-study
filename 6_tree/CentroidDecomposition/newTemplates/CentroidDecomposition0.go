package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
)

func main() {
	yosupo()
}

func demo() {
	//   0
	//   |
	//   1
	//  / \
	//  2  3
	// / \
	// 4  5
	n := int32(6)
	edges := [][]int32{{0, 1}, {1, 2}, {1, 3}, {2, 4}, {2, 5}}
	tree := make([][]int32, n)
	for _, e := range edges {
		tree[e[0]] = append(tree[e[0]], e[1])
		tree[e[1]] = append(tree[e[1]], e[0])
	}
	CentroidDecomposition0(
		n, tree,
		func(parent []int32, vertex []int32, indptr []int32) {
			fmt.Println(parent, vertex, indptr)
		},
	)
}

// https://judge.yosupo.jp/problem/frequency_table_of_tree_distance
// 树上所有点对的距离表(距离指两点路径上的的边数,FrequencyTableofTreeDistance)
// n<=2e5 O(nlognlogn)
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	res := TreeAllDistances(n, tree, Convolution)
	res = res[1:]
	for i := range res {
		res[i] /= 2
	}
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

func TreeAllDistances(n int32, tree [][]int32, convolution func([]int, []int) []int) []int {
	res := make([]int, n)
	var f func(parent []int32, vertex []int32, indptr []int32)
	f = func(parent []int32, vertex []int32, indptr []int32) {
		m := int32(len(parent))
		dist := make([]int32, m)
		for i := int32(1); i < m; i++ {
			dist[i] = dist[parent[i]] + 1
		}
		calc := func(start, end int32, sign int) {
			max_ := int32(-1)
			for i := start; i < end; i++ {
				if dist[i] > max_ {
					max_ = dist[i]
				}
			}
			f := make([]int, max_+1)
			for i := start; i < end; i++ {
				f[dist[i]]++
			}
			f = convolution(f, f)
			for i := int32(0); i < min32(int32(len(f)), n); i++ {
				res[i] += sign * f[i]
			}
		}
		calc(0, m, 1)
		for k := 1; k < len(indptr)-1; k++ {
			calc(indptr[k], indptr[k+1], -1)
		}
	}
	CentroidDecomposition0(n, tree, f)
	return res
}

// 基于顶点的重心分解.
// f(parent, vertex, indptr).
// parent[i] is the parent of vertex[i].
// vertex[0] is the centroid of subtree.
// !vertex[indptr[i]:indptr[i+1]] (i>=1) is the subtree of vertex[i].
// !example: https://maspypy.github.io/library/test/library_checker/tree/frequency_table_of_tree_distance_0.test.cpp
// https://maspypy.com/%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%83%bb1-3%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%81%ae%e3%81%8a%e7%b5%b5%e6%8f%8f%e3%81%8d
func CentroidDecomposition0(
	n int32, tree [][]int32,
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
		for _, next := range tree[v] {
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

// ------------------------------ Convolution ------------------------------

// 计算 A(x) 和 B(x) 的卷积
//
//	c[i] = ∑a[k]*b[i-k], k=0..i
//	入参出参都是次项从低到高的系数
func Convolution(a, b []int) []int {
	n, m := len(a), len(b)
	if n == 0 || m == 0 {
		return nil
	}
	if n <= 1000 || m <= 1000 {
		return convolutionNaive(a, b)
	}
	limit := 1 << bits.Len(uint(n+m-1))
	A := make([]complex128, limit)
	for i, v := range a {
		A[i] = complex(float64(v), 0)
	}
	B := make([]complex128, limit)
	for i, v := range b {
		B[i] = complex(float64(v), 0)
	}
	t := newFFT(limit)
	t.dft(A)
	t.dft(B)
	for i := range A {
		A[i] *= B[i]
	}
	t.idft(A)
	conv := make([]int, n+m-1)
	for i := range conv {
		conv[i] = int(math.Round(real(A[i]))) // % mod
	}
	return conv
}

// 计算多个多项式的卷积
// 入参出参都是次项从低到高的系数
func MultiConvolution(coefs [][]int) []int {
	n := len(coefs)
	if n == 1 {
		return coefs[0]
	}
	return Convolution(MultiConvolution(coefs[:n/2]), MultiConvolution(coefs[n/2:]))
}

// https://github.com/EndlessCheng/codeforces-go/blob/5389a5dd32216aa3572260889a662cce28c1f1f5/copypasta/math_fft.go#L1
type fft struct {
	n               int
	omega, omegaInv []complex128
}

func newFFT(n int) *fft {
	omega := make([]complex128, n)
	omegaInv := make([]complex128, n)
	for i := range omega {
		sin, cos := math.Sincos(2 * math.Pi * float64(i) / float64(n))
		omega[i] = complex(cos, sin)
		omegaInv[i] = complex(cos, -sin)
	}
	return &fft{n, omega, omegaInv}
}

func (t *fft) transform(a, omega []complex128) {
	for i, j := 0, 0; i < t.n; i++ {
		if i > j {
			a[i], a[j] = a[j], a[i]
		}
		for l := t.n >> 1; ; l >>= 1 {
			j ^= l
			if j >= l {
				break
			}
		}
	}
	for l := 2; l <= t.n; l <<= 1 {
		m := l >> 1
		for st := 0; st < t.n; st += l {
			b := a[st:]
			for i := 0; i < m; i++ {
				d := omega[t.n/l*i] * b[m+i]
				b[m+i] = b[i] - d
				b[i] += d
			}
		}
	}
}

func (t *fft) dft(a []complex128) {
	t.transform(a, t.omega)
}

func (t *fft) idft(a []complex128) {
	t.transform(a, t.omegaInv)
	cn := complex(float64(t.n), 0)
	for i := range a {
		a[i] /= cn
	}
}

func convolutionNaive(a, b []int) []int {
	n, m := len(a), len(b)
	conv := make([]int, n+m-1)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			conv[i+j] += a[i] * b[j]
		}
	}
	return conv
}
