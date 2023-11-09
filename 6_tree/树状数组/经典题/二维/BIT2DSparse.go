// !谨慎使用 (2e5, 1700ms; 1e5, 800ms)
// Add : 单点加
// Query : 区间和
// QueryPrefix : 前缀和

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// https://judge.yosupo.jp/problem/point_add_rectangle_sum

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	xs, ys, ws := make([]int, n), make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &ys[i], &ws[i])
	}
	query := make([][4]int, q)
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 0 {
			var x, y, w int
			fmt.Fscan(in, &x, &y, &w)
			xs = append(xs, x)
			ys = append(ys, y)
			ws = append(ws, 0)
			query[i] = [4]int{-1, x, y, w}
		} else {
			var a, b, c, d int
			fmt.Fscan(in, &a, &b, &c, &d)
			query[i] = [4]int{a, c, b, d}
		}
	}

	tree := NewBIT2DSparseWithWeights(xs, ys, ws, true)
	for i := 0; i < q; i++ {
		a, b, c, d := query[i][0], query[i][1], query[i][2], query[i][3]
		if a == -1 {
			tree.Update(b, c, d)
		} else {
			fmt.Fprintln(out, tree.QueryRange(a, b, c, d))
		}
	}
}

func main2() {
	// https://judge.yosupo.jp/problem/rectangle_sum
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	xs, ys, ws := make([]int, n), make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &ys[i], &ws[i])
	}
	tree := NewBIT2DSparseWithWeights(xs, ys, ws, true)
	for i := 0; i < q; i++ {
		var l, d, r, u int
		fmt.Fscan(in, &l, &d, &r, &u)
		fmt.Fprintln(out, tree.QueryRange(l, r, d, u))
	}
}

func demo() {
	xs, ys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	tree := NewBIT2DSparse(xs, ys, false)
	tree.Update(1, 1, 91)
	fmt.Println(tree.QueryRange(1, 2, 1, 2))
	fmt.Println(tree.QueryPrefix(1, 2))
}

type Able = int

// 需要是阿贝尔群(满足交换律)
func e() Able           { return 0 }
func op(a, b Able) Able { return a + b }
func inv(a Able) Able   { return -a }

type BIT2DSparse struct {
	n          int
	keyX       []int
	keyY       []int
	minX       int
	indptr     []int
	data       []Able
	discretize bool
	unit       Able
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewBIT2DSparse(xs, ys []int, discretize bool) *BIT2DSparse {
	res := &BIT2DSparse{discretize: discretize, unit: e()}
	ws := make([]Able, len(xs))
	for i := range ws {
		ws[i] = res.unit
	}
	res._build(xs, ys, ws)
	return res
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewBIT2DSparseWithWeights(xs, ys []int, ws []Able, discretize bool) *BIT2DSparse {
	res := &BIT2DSparse{discretize: discretize, unit: e()}
	res._build(xs, ys, ws)
	return res
}

// 点 (x,y) 的值加上 val.
func (fwt *BIT2DSparse) Update(x, y int, val Able) {
	i := fwt._xtoi(x)
	for ; i < fwt.n; i += ((i + 1) & -(i + 1)) {
		fwt._add(i, y, val)
	}
}

// [lx,rx) * [ly,ry)
func (t *BIT2DSparse) QueryRange(lx, rx, ly, ry int) Able {
	pos, neg := t.unit, t.unit
	l, r := t._xtoi(lx)-1, t._xtoi(rx)-1
	for l < r {
		pos = op(pos, t._prodI(r, ly, ry))
		r -= ((r + 1) & -(r + 1))
	}
	for r < l {
		neg = op(neg, t._prodI(l, ly, ry))
		l -= ((l + 1) & -(l + 1))
	}
	return op(pos, inv(neg))
}

// [0,rx) * [0,ry)
func (t *BIT2DSparse) QueryPrefix(rx, ry int) Able {
	pos := t.unit
	r := t._xtoi(rx) - 1
	for r >= 0 {
		pos = op(pos, t._prefixProdI(r, ry))
		r -= ((r + 1) & -(r + 1))
	}
	return pos
}

func (t *BIT2DSparse) _add(i int, y int, val Able) {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	j := bisectLeft(t.keyY, y, lid, lid+n-1) - lid
	for j < n {
		t.data[lid+j] = op(t.data[lid+j], val)
		j += ((j + 1) & -(j + 1))
	}
}

func (t *BIT2DSparse) _prodI(i int, ly, ry int) Able {
	pos, neg := t.unit, t.unit
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	left := bisectLeft(t.keyY, ly, lid, lid+n-1) - lid - 1
	right := bisectLeft(t.keyY, ry, lid, lid+n-1) - lid - 1
	for left < right {
		pos = op(pos, t.data[lid+right])
		right -= ((right + 1) & -(right + 1))
	}
	for right < left {
		neg = op(neg, t.data[lid+left])
		left -= ((left + 1) & -(left + 1))
	}
	return op(pos, inv(neg))
}

func (t *BIT2DSparse) _prefixProdI(i int, ry int) Able {
	pos := t.unit
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	R := bisectLeft(t.keyY, ry, lid, lid+n-1) - lid - 1
	for R >= 0 {
		pos = op(pos, t.data[lid+R])
		R -= ((R + 1) & -(R + 1))
	}
	return pos
}

func (ft *BIT2DSparse) _build(X, Y []int, wt []Able) {
	if len(X) != len(Y) || len(X) != len(wt) {
		panic("Lengths of X, Y, and wt must be equal.")
	}

	if ft.discretize {
		ft.keyX = unique(X)
		ft.n = len(ft.keyX)
	} else {
		if len(X) > 0 {
			min_, max_ := 0, 0
			for _, x := range X {
				if x < min_ {
					min_ = x
				}
				if x > max_ {
					max_ = x
				}
			}
			ft.minX = min_
			ft.n = max_ - min_ + 1
		}
		ft.keyX = make([]int, ft.n)
		for i := 0; i < ft.n; i++ {
			ft.keyX[i] = ft.minX + i
		}
	}

	N := ft.n
	keyYRaw := make([][]int, N)
	datRaw := make([][]Able, N)
	indices := argSort(Y)

	for _, i := range indices {
		ix := ft._xtoi(X[i])
		y := Y[i]
		for ix < N {
			kY := keyYRaw[ix]
			if len(kY) == 0 || kY[len(kY)-1] < y {
				keyYRaw[ix] = append(keyYRaw[ix], y)
				datRaw[ix] = append(datRaw[ix], wt[i])
			} else {
				datRaw[ix][len(datRaw[ix])-1] = op(datRaw[ix][len(datRaw[ix])-1], wt[i])
			}
			ix += ((ix + 1) & -(ix + 1))
		}
	}

	ft.indptr = make([]int, N+1)
	for i := 0; i < N; i++ {
		ft.indptr[i+1] = ft.indptr[i] + len(keyYRaw[i])
	}
	ft.keyY = make([]int, ft.indptr[N])
	ft.data = make([]Able, ft.indptr[N])

	for i := 0; i < N; i++ {
		for j := 0; j < ft.indptr[i+1]-ft.indptr[i]; j++ {
			ft.keyY[ft.indptr[i]+j] = keyYRaw[i][j]
			ft.data[ft.indptr[i]+j] = datRaw[i][j]
		}
	}

	for i := 0; i < N; i++ {
		n := ft.indptr[i+1] - ft.indptr[i]
		for j := 0; j < n-1; j++ {
			k := j + ((j + 1) & -(j + 1))
			if k < n {
				ft.data[ft.indptr[i]+k] = op(ft.data[ft.indptr[i]+k], ft.data[ft.indptr[i]+j])
			}
		}
	}
}

func (ft *BIT2DSparse) _xtoi(x int) int {
	if ft.discretize {
		return bisectLeft(ft.keyX, x, 0, len(ft.keyX)-1)
	}
	tmp := x - ft.minX
	if tmp < 0 {
		tmp = 0
	} else if tmp > ft.n {
		tmp = ft.n
	}
	return tmp
}

func bisectLeft(nums []int, x int, left, right int) int {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func unique(nums []int) (sorted []int) {
	set := make(map[int]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	sorted = make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	return
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

func argSort(nums []int) []int {
	order := make([]int, len(nums))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}
