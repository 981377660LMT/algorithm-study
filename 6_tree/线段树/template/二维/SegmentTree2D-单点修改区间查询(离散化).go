// 二维线段树区间修改单点查询
// Query(lx, rx, ly, ry) 查询区间[lx, rx) * [ly, ry)的值.
// Update(x, y, val) 将点(x, y)的值加上(op) val.
// !每次修改/查询 O(lognlogn)

// 1e5:1200ms; 2e5:2200ms

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main2() {
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

	tree := NewSegmentTree2DWithWeights(xs, ys, ws, true)
	for i := 0; i < q; i++ {
		a, b, c, d := query[i][0], query[i][1], query[i][2], query[i][3]
		if a == -1 {
			tree.Update(b, c, d)
		} else {
			fmt.Fprintln(out, tree.Query(a, b, c, d))
		}
	}
}

func main() {
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
	tree := NewSegmentTree2DWithWeights(xs, ys, ws, true)
	for i := 0; i < q; i++ {
		var l, d, r, u int
		fmt.Fscan(in, &l, &d, &r, &u)
		fmt.Fprintln(out, tree.Query(l, r, d, u))
	}
}

// 需要满足交换律.
type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }

type SegmentTree2D struct {
	n          int
	keyX       []int
	keyY       []int
	minX       int
	indptr     []int
	data       []E
	discretize bool
	unit       E
}

// discretize:
//  为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//  为 false 时不对x维度二分离散化,而是直接用x的值作为下标(所有x给一个偏移量minX),
//  x 维度数组长度为最大值减最小值.
func NewSegmentTree2D(xs, ys []int, discretize bool) *SegmentTree2D {
	res := &SegmentTree2D{discretize: discretize, unit: e()}
	ws := make([]E, len(xs))
	for i := range ws {
		ws[i] = res.unit
	}
	res._build(xs, ys, ws)
	return res
}

// discretize:
//  为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//  为 false 时不对x维度二分离散化,而是直接用x的值作为下标(所有x给一个偏移量minX),
//  x 维度数组长度为最大值减最小值.
func NewSegmentTree2DWithWeights(xs, ys []int, ws []E, discretize bool) *SegmentTree2D {
	res := &SegmentTree2D{discretize: discretize, unit: e()}
	res._build(xs, ys, ws)
	return res
}

// 点 (x,y) 的值加上 val.
func (t *SegmentTree2D) Update(x, y int, val E) {
	i := t._xtoi(x)
	i += t.n
	for i > 0 {
		t._add(i, y, val)
		i >>= 1
	}
}

// [lx,rx) * [ly,ry)
func (t *SegmentTree2D) Query(lx, rx, ly, ry int) E {
	L := t._xtoi(lx) + t.n
	R := t._xtoi(rx) + t.n
	val := t.unit
	for L < R {
		if L&1 == 1 {
			val = op(val, t._prodI(L, ly, ry))
			L++
		}
		if R&1 == 1 {
			R--
			val = op(t._prodI(R, ly, ry), val)
		}
		L >>= 1
		R >>= 1
	}
	return val
}

func (t *SegmentTree2D) _add(i int, y int, val E) {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	j := bisectLeft(t.keyY, y, lid, lid+n-1) - lid
	offset := 2 * lid
	j += n
	for j > 0 {
		t.data[offset+j] = op(t.data[offset+j], val)
		j >>= 1
	}
}

func (t *SegmentTree2D) _prodI(i int, ly, ry int) E {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	left := bisectLeft(t.keyY, ly, lid, lid+n-1) - lid + n
	right := bisectLeft(t.keyY, ry, lid, lid+n-1) - lid + n
	offset := 2 * lid
	val := t.unit
	for left < right {
		if left&1 == 1 {
			val = op(val, t.data[offset+left])
			left++
		}
		if right&1 == 1 {
			right--
			val = op(t.data[offset+right], val)
		}
		left >>= 1
		right >>= 1
	}
	return val
}

func (seg *SegmentTree2D) _build(X, Y []int, wt []E) {
	if len(X) != len(Y) || len(X) != len(wt) {
		panic("Lengths of X, Y, and wt must be equal.")
	}

	if seg.discretize {
		seg.keyX = unique(X)
		seg.n = len(seg.keyX)
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
			seg.minX = min_
			seg.n = max_ - min_ + 1
		}
		seg.keyX = make([]int, seg.n)
		for i := 0; i < seg.n; i++ {
			seg.keyX[i] = seg.minX + i
		}
	}

	N := seg.n
	keyYRaw := make([][]int, N+N)
	datRaw := make([][]E, N+N)
	indices := argSort(Y)

	for _, i := range indices {
		ix, y := seg._xtoi(X[i]), Y[i]
		ix += N
		for ix > 0 {
			KY := keyYRaw[ix]
			if len(KY) == 0 || KY[len(KY)-1] < y {
				keyYRaw[ix] = append(keyYRaw[ix], y)
				datRaw[ix] = append(datRaw[ix], wt[i])
			} else {
				datRaw[ix][len(datRaw[ix])-1] = op(datRaw[ix][len(datRaw[ix])-1], wt[i])
			}
			ix >>= 1
		}
	}

	seg.indptr = make([]int, N+N+1)
	for i := 0; i < N+N; i++ {
		seg.indptr[i+1] = seg.indptr[i] + len(keyYRaw[i])
	}
	fullN := seg.indptr[N+N]
	seg.keyY = make([]int, fullN)
	seg.data = make([]E, 2*fullN)

	for i := 0; i < N+N; i++ {
		off := 2 * seg.indptr[i]
		n := seg.indptr[i+1] - seg.indptr[i]
		for j := 0; j < n; j++ {
			seg.keyY[seg.indptr[i]+j] = keyYRaw[i][j]
			seg.data[off+n+j] = datRaw[i][j]
		}
		for j := n - 1; j >= 1; j-- {
			seg.data[off+j] = op(seg.data[off+2*j], seg.data[off+2*j+1])
		}
	}
}

func (seg *SegmentTree2D) _xtoi(x int) int {
	if seg.discretize {
		return bisectLeft(seg.keyX, x, 0, len(seg.keyX)-1)
	}
	tmp := x - seg.minX
	if tmp < 0 {
		tmp = 0
	} else if tmp > seg.n {
		tmp = seg.n
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
