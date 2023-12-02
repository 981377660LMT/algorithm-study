// 树上等高线汇集
// https://maspypy.com/%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%83%bb1-3%e9%87%8d%e5%bf%83%e5%88%86%e8%a7%a3%e3%81%ae%e3%81%8a%e7%b5%b5%e6%8f%8f%e3%81%8d
// !注意不包含距离0,需额外处理.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// VertexGetRangeContourAdd()
	Yuki1038()
}

// https://judge.yosupo.jp/problem/vertex_add_range_contour_sum_on_tree
// 给定q个操作，操作有两种：
// 0 root x : 将root节点的值加上x (点权加)
// 1 root floor higher: 求出距离root节点距离在[floor,higher)之间的所有节点的值的和 (区间点权和)
// n<=1e5 q<=2e5
func VertexAddRangeContourSum() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	weights := make([]int, n)
	for i := range weights {
		fmt.Fscan(in, &weights[i])
	}
	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}

	C := NewContourQueryRange(n, tree)
	data := make([]int, C.Size())
	for i := 0; i < n; i++ {
		C.EnumeratePoint(i, func(pos int) {
			data[pos] += weights[i]
		})
	}

	bit := NewBitArrayFrom(len(data), func(i int) int { return data[i] })
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 0 {
			var root, x int
			fmt.Fscan(in, &root, &x)

			weights[root] += x
			C.EnumeratePoint(root, func(pos int) {
				bit.Add(pos, x)
			})

		} else {
			var root, floor, higher int
			fmt.Fscan(in, &root, &floor, &higher)
			res := 0

			if floor <= 0 && 0 < higher {
				res += weights[root]
			}
			C.EnumerateRange(root, floor, higher, func(start, end int) {
				res += bit.QueryRange(start, end)
			})

			fmt.Fprintln(out, res)
		}
	}
}

// https://judge.yosupo.jp/problem/vertex_get_range_contour_add_on_tree
// 给定q个操作，操作有两种：
// 0 root floor higher x: 距离root节点距离在[floor,higher)之间的所有节点的值加上x (区间点权加)
// 1 root : 求出root节点的值 (点权)
// n<=1e5 q<=2e5
func VertexGetRangeContourAdd() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	values := make([]int, n)
	for i := range values {
		fmt.Fscan(in, &values[i])
	}
	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}

	C := NewContourQueryRange(n, tree)
	bit := NewBitArray(C.Size() + 1)

	add := func(root, floor, higher, x int) {
		if floor <= 0 && 0 < higher {
			values[root] += x
		}
		C.EnumerateRange(root, floor, higher, func(start, end int) {
			bit.Add(start, x)
			bit.Add(end, -x)
		})
	}

	query := func(root int) int {
		res := values[root]
		C.EnumeratePoint(root, func(pos int) {
			res += bit.QueryPrefix(pos + 1)
		})
		return res
	}

	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 0 {
			var root, floor, higher, x int
			fmt.Fscan(in, &root, &floor, &higher, &x)
			add(root, floor, higher, x)
		} else {
			var root int
			fmt.Fscan(in, &root)
			fmt.Fprintln(out, query(root))
		}
	}
}

// https://yukicoder.me/problems/no/1038
// 给定一颗树，初始时点权为0，有q个操作.
// 每次操作输出顶点x的点权，并将距离x距离在[0,y]之间的所有节点的点权加上z.
func Yuki1038() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}
	weights := make([]int, n)

	C := NewContourQueryRange(n, tree)
	bit := NewBitArray(C.Size() + 1)
	for i := 0; i < q; i++ {
		var node, dist, delta int
		fmt.Fscan(in, &node, &dist, &delta)
		node--
		dist++

		res := weights[node]
		C.EnumeratePoint(node, func(pos int) {
			res += bit.QueryPrefix(pos + 1)
		})
		fmt.Fprintln(out, res)

		weights[node] += delta
		C.EnumerateRange(node, 0, dist, func(start, end int) {
			bit.Add(start, delta)
			bit.Add(end, -delta)
		})
	}

}

// !注意不包含距离0.
type ContourQueryRange struct {
	_n          int
	_v          []int
	_comp       []int
	_dep        []int
	_infoIdx    []int
	_infoIndptr []int
	_compRange  []int
}

func NewContourQueryRange(n int, graph [][]int) *ContourQueryRange {
	p := 0
	compRange := []int{0}
	V := []int{}
	comp := []int{}
	dep := []int{}
	f := func(par []int, vs []int, color []int) {
		n := len(par)
		dist := make([]int, n)
		A := make([]int, 0)
		B := make([]int, 0)
		for v := 1; v < n; v++ {
			dist[v] = dist[par[v]] + 1
		}
		for v := 0; v < n; v++ {
			if color[v] == 0 {
				A = append(A, v)
			}
			if color[v] == 1 {
				B = append(B, v)
			}
		}
		if len(A) == 0 || len(B) == 0 {
			return
		}
		mxA := 0
		mxB := 0
		for _, v := range A {
			V = append(V, vs[v])
			comp = append(comp, p)
			dep = append(dep, dist[v])
			mxA = max(mxA, dist[v])
		}
		compRange = append(compRange, compRange[len(compRange)-1]+mxA+1)
		p++
		for _, v := range B {
			V = append(V, vs[v])
			comp = append(comp, p)
			dep = append(dep, dist[v])
			mxB = max(mxB, dist[v])
		}
		compRange = append(compRange, compRange[len(compRange)-1]+mxB+1)
		p++
	}
	centroidDecomposition(n, graph, f)
	infoIndptr := make([]int, n+1)
	for _, v := range V {
		infoIndptr[v+1]++
	}
	for v := 0; v < n; v++ {
		infoIndptr[v+1] += infoIndptr[v]
	}
	counter := append([]int{}, infoIndptr...)
	infoIdx := make([]int, infoIndptr[len(infoIndptr)-1])
	for i := 0; i < len(V); i++ {
		infoIdx[counter[V[i]]] = i
		counter[V[i]]++
	}
	return &ContourQueryRange{
		_n:          n,
		_v:          V,
		_comp:       comp,
		_dep:        dep,
		_infoIdx:    infoIdx,
		_infoIndptr: infoIndptr,
		_compRange:  compRange,
	}
}

func (cqr *ContourQueryRange) Size() int {
	return cqr._compRange[len(cqr._compRange)-1]
}

func (cqr *ContourQueryRange) EnumerateRange(node int, start int, end int, f func(int, int)) {
	for k := cqr._infoIndptr[node]; k < cqr._infoIndptr[node+1]; k++ {
		idx := cqr._infoIdx[k]
		p := cqr._comp[idx] ^ 1
		lo := start - cqr._dep[idx]
		hi := end - cqr._dep[idx]
		L := cqr._compRange[p]
		R := cqr._compRange[p+1]
		n := R - L
		lo = max(lo, 0)
		hi = min(hi, n)
		if lo < hi {
			f(L+lo, L+hi)
		}
	}
}

func (cqr *ContourQueryRange) EnumeratePoint(v int, f func(int)) {
	for k := cqr._infoIndptr[v]; k < cqr._infoIndptr[v+1]; k++ {
		idx := cqr._infoIdx[k]
		p := cqr._comp[idx]
		f(cqr._compRange[p] + cqr._dep[idx])
	}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int, f func(i int) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITArray) QueryAll() int {
	return b.total
}

func (b *BITArray) MaxRight(check func(index, preSum int) bool) int {
	i := 0
	s := 0
	k := 1
	for 2*k <= b.n {
		k *= 2
	}
	for k > 0 {
		if i+k-1 < b.n {
			t := s + b.data[i+k-1]
			if check(i+k, t) {
				i += k
				s = t
			}
		}
		k >>= 1
	}
	return i
}

// 0/1树状数组查找第 k(0-based) 个1的位置.
// UpperBound.
func (b *BITArray) Kth(k int) int {
	return b.MaxRight(func(index, preSum int) bool { return preSum <= k })
}

func (b *BITArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}
