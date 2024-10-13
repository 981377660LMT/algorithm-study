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
	// demo()
	// [[2,1]]

	// fmt.Println(countPaths(2, [][]int{{1, 2}}))
}

func demo() {
	{

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
		CentroidDecomposition1(
			n, tree,
			func(parent []int32, vertex []int32, l1, r1, l2, r2 int32) {
				fmt.Println(parent, vertex, l1, r1, l2, r2)
			},
		)
	}

	{
		{
			n := int32(1e5)
			edges := make([][]int32, n-1)
			for i := int32(0); i < n-1; i++ {
				edges[i] = []int32{0, i + 1} // star graph
			}
			tree := make([][]int32, n)
			for _, e := range edges {
				tree[e[0]] = append(tree[e[0]], e[1])
				tree[e[1]] = append(tree[e[1]], e[0])
			}
			count := 0
			vCount := 0
			vCounter := make([]int32, n)
			CentroidDecomposition1(
				n, tree,
				func(parent []int32, vertex []int32, l1, r1, l2, r2 int32) {
					count++
					vCount += len(vertex)
					for _, v := range vertex {
						vCounter[v]++
					}
				},
			)
			fmt.Println(count, vCount, vCounter[:100]) // 星图中，可以看到根结点被计算了O(n)次
		}
	}

}

// 2867. 统计树中的合法路径数目
var E *eratosthenesSieve

func init() { E = newEratosthenesSieve(1e5 + 10) }
func countPaths(n int, edges [][]int) int64 {
	tree := make([][]int32, n)
	for _, e := range edges {
		u, v := int32(e[0])-1, int32(e[1])-1
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	res := 0
	f := func(parent []int32, vertex []int32, l1, r1, l2, r2 int32) {
		m := int32(len(vertex))
		centroid := vertex[0]
		primeCounter := make([]int32, m)
		isCentroidPrime := E.IsPrime(int(centroid + 1))
		if isCentroidPrime {
			primeCounter[0] = 1
		}
		for i := int32(1); i < m; i++ {
			primeCounter[i] = primeCounter[parent[i]]
			if E.IsPrime(int(vertex[i] + 1)) {
				primeCounter[i]++
			}
		}

		left0, left1, right0, right1 := 0, 0, 0, 0
		for i := l1; i < r1; i++ {
			if primeCounter[i] == 0 {
				left0++
			} else if primeCounter[i] == 1 {
				left1++
			}
		}
		for i := l2; i < r2; i++ {
			if primeCounter[i] == 0 {
				right0++
			} else if primeCounter[i] == 1 {
				right1++
			}
		}

		if isCentroidPrime {
			res += left1 * right1
		} else {
			res += left1*right0 + left0*right1
		}
	}

	CentroidDecomposition1(int32(n), tree, f)

	// !特殊处理节点数<=2的情况
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		if E.IsPrime(u+1) != E.IsPrime(v+1) {
			res++
		}
	}

	return int64(res)
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
	var f func(parent []int32, vertex []int32, l1, r1, l2, r2 int32)
	f = func(parent []int32, vertex []int32, l1, r1, l2, r2 int32) {
		m := int32(len(vertex))
		dist := make([]int32, m)
		maxDist := int32(0)
		for i := int32(1); i < m; i++ {
			dist[i] = dist[parent[i]] + 1
			if dist[i] > maxDist {
				maxDist = dist[i]
			}
		}
		f, g := make([]int, maxDist+1), make([]int, maxDist+1)
		for i := l1; i < r1; i++ {
			f[dist[i]]++
		}
		for i := l2; i < r2; i++ {
			g[dist[i]]++
		}
		for len(f) > 0 && f[len(f)-1] == 0 {
			f = f[:len(f)-1]
		}
		for len(g) > 0 && g[len(g)-1] == 0 {
			g = g[:len(g)-1]
		}
		f = convolution(f, g)
		for i, v := range f {
			res[i] += v + v
		}
	}
	CentroidDecomposition1(n, tree, f)
	res[0] = int(n)
	res[1] = 2 * int(n-1)
	return res
}

// 1/3重心分解(1/3 Centroid Decomposition)
//
//		!f(parent, vertex, l1, r1, l2, r2) 处理经过重心的路径. 路径长度>=3，保证左、右子树节点数>=1.
//	  !注意需要特殊处理路径长度为2的情况.
//		  vertex[0] is the centroid of subtree.
//			[l1, r1): color 1
//			[l2, r2): color 2
//
// !example: https://maspypy.github.io/library/graph/tree_all_distances.hpp
// https://maspypy.com/%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%83%bb1-3%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%81%ae%e3%81%8a%e7%b5%b5%e6%8f%8f%e3%81%8d
// https://codeforces.com/blog/entry/104997
func CentroidDecomposition1(
	n int32, tree [][]int32,
	f func(parent []int32, vertex []int32, l1, r1, l2, r2 int32),
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
	centroidDecomposition1Dfs(parent, queue, f)
}

func centroidDecomposition1Dfs(
	parent []int32, vs []int32,
	f func(parent []int32, vertex []int32, l1, r1, l2, r2 int32),
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
	color := make([]int8, n)
	ord := make([]int32, n)
	for i := range color {
		color[i] = -1
		ord[i] = -1
	}
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
	f(par2, V2, 1, 1+n0, 1+n0, 1+n0+n1)
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
	if n <= 500 || m <= 500 {
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

// 埃氏筛
type eratosthenesSieve struct {
	minPrime []int
}

func newEratosthenesSieve(maxN int) *eratosthenesSieve {
	minPrime := make([]int, maxN+1)
	for i := range minPrime {
		minPrime[i] = i
	}
	upper := int(math.Sqrt(float64(maxN))) + 1
	for i := 2; i < upper; i++ {
		if minPrime[i] < i {
			continue
		}
		for j := i * i; j <= maxN; j += i {
			if minPrime[j] == j {
				minPrime[j] = i
			}
		}
	}
	return &eratosthenesSieve{minPrime}
}

func (es *eratosthenesSieve) IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	return es.minPrime[n] == n
}

func (es *eratosthenesSieve) GetPrimeFactors(n int) map[int]int {
	res := make(map[int]int)
	for n > 1 {
		m := es.minPrime[n]
		res[m]++
		n /= m
	}
	return res
}

func (es *eratosthenesSieve) GetPrimes() []int {
	res := []int{}
	for i, x := range es.minPrime {
		if i >= 2 && i == x {
			res = append(res, x)
		}
	}
	return res
}
