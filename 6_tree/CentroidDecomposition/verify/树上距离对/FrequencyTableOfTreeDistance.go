// https://judge.yosupo.jp/problem/frequency_table_of_tree_distance
// 树上所有点对的距离表(距离指两点路径上的的边数,Frequency Table of Tree Distance)
// n<=2e5 O(nlognlogn)

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	g := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], Edge{to: v})
		g[v] = append(g[v], Edge{to: u})
	}

	pairDist := FreqTableOfTreeDistance(g)
	for _, v := range pairDist[1:] {
		fmt.Fprintln(out, v, " ")
	}
}

type Edge = struct{ to, weight int }

// pairDist[i] 表示距离为i的点对的数量
func FreqTableOfTreeDistance(g [][]Edge) (pairDist []int) {
	n := len(g)
	centTree, root := CentroidDecomposition(g)
	removed := make([]bool, n)
	pairDist = make([]int, 2*n) // 开大一倍,保证中间卷积后统计距离时不会越界

	var decomposition func(cur, pre int)
	decomposition = func(cur, pre int) {
		removed[cur] = true
		for _, next := range centTree[cur] {
			if !removed[next] {
				decomposition(next, cur)
			}
		}
		removed[cur] = false

		allDist := []int{1}

		for _, e := range g[cur] {
			next := e.to
			if next == pre || removed[next] {
				continue
			}

			dist := []int{0}
			queue := [][3]int{{next, cur, 1}} // collectDist by bfs
			for len(queue) > 0 {
				idx, par, dep := queue[0][0], queue[0][1], queue[0][2]
				queue = queue[1:]
				if len(allDist) <= dep {
					allDist = append(allDist, 0)
				}
				if len(dist) <= dep {
					dist = append(dist, 0)
				}
				allDist[dep]++
				dist[dep]++
				for _, e := range g[idx] {
					next := e.to
					if next == par || removed[next] {
						continue
					}
					queue = append(queue, [3]int{next, idx, dep + 1})
				}
			}

			res := Convolution(dist, dist)
			for i := range res {
				pairDist[i] -= res[i] // 同一个子树中都经过当前重心的距离对(不合法)
			}
		}
		res := Convolution(allDist, allDist)
		for i := range res {
			pairDist[i] += res[i] // 子树中所有经过当前重心的距离对
		}

	}

	decomposition(root, -1)
	pairDist = pairDist[:n]
	for i := range pairDist {
		pairDist[i] /= 2
	}
	return
}

// 计算 A(x) 和 B(x) 的卷积
//
//	c[i] = ∑a[k]*b[i-k], k=0..i
//	入参出参都是次项从低到高的系数
func Convolution(a, b []int) []int {
	n, m := len(a), len(b)
	limit := 1 << uint(bits.Len(uint(n+m-1)))
	f := newFFT(limit)
	cmplxA := make([]complex128, limit)
	for i, v := range a {
		cmplxA[i] = complex(float64(v), 0)
	}
	cmplxB := make([]complex128, limit)
	for i, v := range b {
		cmplxB[i] = complex(float64(v), 0)
	}
	f.dft(cmplxA)
	f.dft(cmplxB)
	for i := range cmplxA {
		cmplxA[i] *= cmplxB[i]
	}
	f.idft(cmplxA)
	conv := make([]int, n+m-1)
	for i := range conv {
		conv[i] = int(math.Round(real(cmplxA[i])))
	}
	return conv
}

// 计算多个多项式的卷积
//
//	入参出参都是次项从低到高的系数
func PolyConvolution(coefs [][]int) []int {
	n := len(coefs)
	if n == 1 {
		return coefs[0]
	}
	return Convolution(PolyConvolution(coefs[:n/2]), PolyConvolution(coefs[n/2:]))
}

// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
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

func (f *fft) transform(a, omega []complex128) {
	for i, j := 0, 0; i < f.n; i++ {
		if i > j {
			a[i], a[j] = a[j], a[i]
		}
		for l := f.n >> 1; ; l >>= 1 {
			j ^= l
			if j >= l {
				break
			}
		}
	}
	for l := 2; l <= f.n; l <<= 1 {
		m := l >> 1
		for st := 0; st < f.n; st += l {
			p := a[st:]
			for i := 0; i < m; i++ {
				t := omega[f.n/l*i] * p[m+i]
				p[m+i] = p[i] - t
				p[i] += t
			}
		}
	}
}

func (f *fft) dft(a []complex128) {
	f.transform(a, f.omega)
}

func (f *fft) idft(a []complex128) {
	f.transform(a, f.omegaInv)
	for i := range a {
		a[i] /= complex(float64(f.n), 0)
	}
}

// 树的重心分解, 返回点分树和点分树的根
//
//	!tree: `无向`树的邻接表.
//	centTree: 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
//	root: 点分树的根
func CentroidDecomposition(tree [][]Edge) (centTree [][]int, root int) {
	n := len(tree)
	subSize := make([]int, n)
	removed := make([]bool, n)
	centTree = make([][]int, n)
	var getSize func(cur, parent int) int
	var getCentroid func(cur, parent, mid int) int
	var build func(cur int) int

	getSize = func(cur, parent int) int {
		subSize[cur] = 1
		for _, e := range tree[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			subSize[cur] += getSize(next, cur)
		}
		return subSize[cur]
	}
	getCentroid = func(cur, parent, mid int) int {
		for _, e := range tree[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			if subSize[next] > mid {
				return getCentroid(next, cur, mid)
			}
		}
		return cur
	}
	build = func(cur int) int {
		centroid := getCentroid(cur, -1, getSize(cur, -1)/2)
		removed[centroid] = true
		for _, e := range tree[centroid] {
			next := e.to
			if !removed[next] {
				centTree[centroid] = append(centTree[centroid], build(next))
			}
		}
		removed[centroid] = false
		return centroid
	}

	root = build(0)
	return
}
