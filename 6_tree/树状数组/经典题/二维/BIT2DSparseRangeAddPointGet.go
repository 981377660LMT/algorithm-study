// Add : 区间加
// Get : 单点查询

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// https://judge.yosupo.jp/problem/rectangle_add_point_get
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	initData := make([][5]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &initData[i][0], &initData[i][1], &initData[i][2], &initData[i][3], &initData[i][4])
	}

	queries := make([][5]int, q)
	var xs, ys []int
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 0 {
			var left, down, right, up, value int
			fmt.Fscan(in, &left, &down, &right, &up, &value)
			queries[i] = [5]int{left, down, right, up, value}
		} else {
			var x, y int
			fmt.Fscan(in, &x, &y)
			xs = append(xs, x)
			ys = append(ys, y)
			queries[i] = [5]int{x, y, -1, -1, -1}
		}
	}

	tree := NewBIT2DSparseRangeAddPointGet(xs, ys, true)
	for i := 0; i < n; i++ {
		a, b, c, d, value := initData[i][0], initData[i][1], initData[i][2], initData[i][3], initData[i][4]
		tree.Update(a, c, b, d, value)
	}
	for i := 0; i < q; i++ {
		a, b, c, d, value := queries[i][0], queries[i][1], queries[i][2], queries[i][3], queries[i][4]
		if value == -1 {
			fmt.Fprintln(out, tree.Get(a, b))
		} else {
			tree.Update(a, c, b, d, value)
		}
	}

}

type Able = int

// 需要是阿贝尔群(满足交换律)
func e() Able           { return 0 }
func op(a, b Able) Able { return a + b }
func inv(a Able) Able   { return -a }

type BIT2DSparseRangeAddPointGet struct {
	n          int
	keyX       []int
	keyY       []int
	minX       int
	indptr     []int
	data       []Able
	discretize bool
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewBIT2DSparseRangeAddPointGet(xs, ys []int, discretize bool) *BIT2DSparseRangeAddPointGet {
	res := &BIT2DSparseRangeAddPointGet{discretize: discretize}
	res._build(xs, ys)
	return res
}

func (bit *BIT2DSparseRangeAddPointGet) Get(x, y int) Able {
	res := e()
	i := bit._xtoi(x)
	for i < bit.n {
		res = op(res, bit._getI(i, y))
		i += ((i + 1) & -(i + 1))
	}
	return res
}

// 区间[startX, endX) * [startY, endY) 加上 value.
func (bit *BIT2DSparseRangeAddPointGet) Update(startX, endX, startY, endY int, value Able) {
	L := bit._xtoi(startX) - 1
	R := bit._xtoi(endX) - 1
	neg := inv(value)
	for L < R {
		bit._add(R, startY, endY, value)
		R -= ((R + 1) & -(R + 1))
	}
	for R < L {
		bit._add(L, startY, endY, neg)
		L -= ((L + 1) & -(L + 1))
	}
}

func (bit *BIT2DSparseRangeAddPointGet) _getI(i, y int) Able {
	res := e()
	lid, n := bit.indptr[i], bit.indptr[i+1]-bit.indptr[i]
	j := bisectLeft(bit.keyY, y, lid, lid+n-1) - lid
	for j < n {
		res = op(res, bit.data[lid+j])
		j += ((j + 1) & -(j + 1))
	}
	return res
}

func (bit *BIT2DSparseRangeAddPointGet) _add(i int, ly, ry int, val Able) {
	neg := inv(val)
	lid, n := bit.indptr[i], bit.indptr[i+1]-bit.indptr[i]
	left := bisectLeft(bit.keyY, ly, lid, lid+n-1) - lid - 1
	right := bisectLeft(bit.keyY, ry, lid, lid+n-1) - lid - 1
	for left < right {
		bit.data[lid+right] = op(val, bit.data[lid+right])
		right -= ((right + 1) & -(right + 1))
	}
	for right < left {
		bit.data[lid+left] = op(neg, bit.data[lid+left])
		left -= ((left + 1) & -(left + 1))
	}
}

func (bit *BIT2DSparseRangeAddPointGet) _build(X, Y []int) {
	if len(X) != len(Y) {
		panic("Lengths of X, Y must be equal.")
	}

	if bit.discretize {
		bit.keyX = unique(X)
		bit.n = len(bit.keyX)
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
			bit.minX = min_
			bit.n = max_ - min_ + 1
		}
		bit.keyX = make([]int, bit.n)
		for i := 0; i < bit.n; i++ {
			bit.keyX[i] = bit.minX + i
		}
	}

	N := bit.n
	keyYRaw := make([][]int, N)
	indices := argSort(Y)
	for _, i := range indices {
		ix := bit._xtoi(X[i])
		y := Y[i]
		for ix < N {
			kY := keyYRaw[ix]
			if len(kY) == 0 || kY[len(kY)-1] < y {
				keyYRaw[ix] = append(keyYRaw[ix], y)
			}
			ix += ((ix + 1) & -(ix + 1))
		}
	}

	bit.indptr = make([]int, N+1)
	for i := 0; i < N; i++ {
		bit.indptr[i+1] = bit.indptr[i] + len(keyYRaw[i])
	}
	bit.keyY = make([]int, bit.indptr[N])
	bit.data = make([]Able, bit.indptr[N])
	for i := range bit.data {
		bit.data[i] = e()
	}

	for i := 0; i < N; i++ {
		for j := 0; j < bit.indptr[i+1]-bit.indptr[i]; j++ {
			bit.keyY[bit.indptr[i]+j] = keyYRaw[i][j]
		}
	}
}

func (bit *BIT2DSparseRangeAddPointGet) _xtoi(x int) int {
	if bit.discretize {
		return bisectLeft(bit.keyX, x, 0, len(bit.keyX)-1)
	}
	tmp := x - bit.minX
	if tmp < 0 {
		tmp = 0
	} else if tmp > bit.n {
		tmp = bit.n
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
