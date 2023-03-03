// # RectangleSum
// # https://judge.yosupo.jp/problem/point_add_rectangle_sum

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

//  0 x y w : 向平面上添加一个点(x,y)并且权值为w
//  1 left down right up : 查询矩形 left<=x<right down<=y<up 内的点权和
//  n,q<=1e5 0<=xi,yi<=1e9
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	points := make([][]int, 0, n)
	for i := 0; i < n; i++ {
		var x, y, w int
		fmt.Fscan(in, &x, &y, &w)
		points = append(points, []int{x, y, w})
	}
	queries := make([][]int, 0, q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var x, y, w int
			fmt.Fscan(in, &x, &y, &w)
			queries = append(queries, []int{op, x, y, w})
		} else {
			var left, down, right, up int
			fmt.Fscan(in, &left, &down, &right, &up)
			queries = append(queries, []int{op, left, down, right, up})
		}
	}

	for _, query := range queries {
		if query[0] == 0 {
			x, y := query[1], query[2]
			points = append(points, []int{x, y, 0}) // 0: 预先添加可能出现的点
		}
	}

	res := []int{}
	rectangleSum := NewPointAddRectangleSum(points)
	for _, query := range queries {
		if query[0] == 0 {
			x, y, w := query[1], query[2], query[3]
			rectangleSum.Add(x, y, w)
		} else {
			left, down, right, up := query[1], query[2], query[3], query[4]
			res = append(res, rectangleSum.Query(left, right, down, up))
		}
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type CompressedPointAddRectangleSum struct {
	xs, ys []int
	idxs   map[struct{ x, y int }]int
	mat    *pointAddRectangleSum
}

// 二维矩形区域计数
//  需要预先添加可能出现的所有点 (x,y,weight)
func NewPointAddRectangleSum(points [][]int) *CompressedPointAddRectangleSum {
	res := &CompressedPointAddRectangleSum{}

	points = append(points[:0:0], points...)
	sort.Slice(points, func(i, j int) bool {
		return points[i][0] < points[j][0]
	})

	xs, ys, ws := make([]int, len(points)), make([]int, len(points)), make([]int, len(points))
	for i, p := range points {
		xs[i], ys[i], ws[i] = p[0], p[1], p[2]
	}
	res.xs = xs

	set := make(map[int]struct{}, len(ys))
	for _, y := range ys {
		set[y] = struct{}{}
	}
	sortedSet := make([]int, 0, len(set))
	for y := range set {
		sortedSet = append(sortedSet, y)
	}
	sort.Ints(sortedSet)
	res.ys = sortedSet

	ids := make(map[struct{ x, y int }]int, len(points))
	for i, p := range points {
		x, y := p[0], p[1]
		ids[struct{ x, y int }{x, y}] = i
	}
	res.idxs = ids

	comp := make(map[int]int, len(sortedSet))
	for i, v := range sortedSet {
		comp[v] = i
	}

	newYs := make([]int, len(ys))
	for i, v := range ys {
		newYs[i] = comp[v]
	}

	maxLog := bits.Len(uint(len(sortedSet)))
	res.mat = newPointAddRectangleSum(newYs, ws, maxLog)
	return res
}

func (f *CompressedPointAddRectangleSum) Add(x, y, w int) {
	idx := f.idxs[struct{ x, y int }{x, y}]
	f.mat.PointAdd(idx, w)
}

// 求矩形x1<=x<x2,y1<=y<y2的权值和 注意是左闭右开
func (f *CompressedPointAddRectangleSum) Query(x1, x2, y1, y2 int) int {
	return f.rectangleSum(x1, x2, y2) - f.rectangleSum(x1, x2, y1)
}

func (f *CompressedPointAddRectangleSum) rectangleSum(l, r, upper int) int {
	l = sort.SearchInts(f.xs, l)
	r = sort.SearchInts(f.xs, r)
	upper = sort.SearchInts(f.ys, upper)
	return f.mat.RectangleSum(l, r, upper)
}

type bitVector struct {
	n     int
	block []int
	sum   []int // block ごとに立っている 1 の数の累積和
}

func newBitVector(n int) *bitVector {
	blockCount := (n + 63) >> 6
	return &bitVector{
		n:     n,
		block: make([]int, blockCount),
		sum:   make([]int, blockCount),
	}
}

func (f *bitVector) Set(i int) {
	f.block[i>>6] |= 1 << uint(i&63)
}

func (f *bitVector) Build() {
	for i := 0; i < len(f.block)-1; i++ {
		f.sum[i+1] = f.sum[i] + bits.OnesCount(uint(f.block[i]))
	}
}

// 统计 [0,end) 中 value 的个数
func (f *bitVector) Count(value, end int) int {
	if value == 1 {
		return f.count1(end)
	}
	return end - f.count1(end)
}

func (f *bitVector) count1(k int) int {
	mask := (1 << uint(k&63)) - 1
	return f.sum[k>>6] + bits.OnesCount(uint(f.block[k>>6]&mask))
}

type binaryIndexedTree struct {
	n   int
	bit []int
}

func newBinaryIndexedTree(n int) *binaryIndexedTree {
	return &binaryIndexedTree{
		n:   n,
		bit: make([]int, n+1),
	}
}

func (f *binaryIndexedTree) Build(a []int) {
	for i := 1; i < len(f.bit); i++ {
		f.bit[i] += a[i-1]
		if j := i + (i & -i); j < len(f.bit) {
			f.bit[j] += f.bit[i]
		}
	}
}

func (f *binaryIndexedTree) Add(i, x int) {
	i++
	for i <= f.n {
		f.bit[i] += x
		i += i & -i
	}
}

func (f *binaryIndexedTree) Sum(l, r int) int {
	return f.sum(r) - f.sum(l)
}

func (f *binaryIndexedTree) sum(i int) int {
	res := 0
	for i > 0 {
		res += f.bit[i]
		i -= i & -i
	}
	return res
}

type pointAddRectangleSum struct {
	n      int
	maxLog int
	arr    []int
	mat    []*bitVector
	zs     []int
	data   []*binaryIndexedTree
}

func newPointAddRectangleSum(array, ws []int, maxLog int) *pointAddRectangleSum {
	res := &pointAddRectangleSum{
		n:      len(array),
		maxLog: maxLog,
		arr:    array,
		mat:    make([]*bitVector, 0, maxLog),
		zs:     make([]int, 0, maxLog),
		data:   make([]*binaryIndexedTree, maxLog),
	}

	n := len(array)
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}

	for d := maxLog - 1; d >= 0; d-- {
		vec := newBitVector(n + 1)
		ls, rs := []int{}, []int{}
		for i, v := range order {
			if ((array[v] >> d) & 1) == 1 {
				rs = append(rs, v)
				vec.Set(i)
			} else {
				ls = append(ls, v)
			}
		}
		vec.Build()
		res.mat = append(res.mat, vec)
		res.zs = append(res.zs, len(ls))
		order = append(ls, rs...)
		res.data[len(res.data)-1-d] = newBinaryIndexedTree(n)
		nums := make([]int, n)
		for i, v := range order {
			nums[i] = ws[v]
		}
		res.data[len(res.data)-1-d].Build(nums)
	}

	return res
}

func (f *pointAddRectangleSum) PointAdd(k, val int) {
	y := f.arr[k]
	for d := 0; d < f.maxLog; d++ {
		if y>>(f.maxLog-d-1)&1 == 1 {
			k = f.mat[d].Count(1, k) + f.zs[d]
		} else {
			k = f.mat[d].Count(0, k)
		}

		f.data[d].Add(k, val)
	}
}

func (f *pointAddRectangleSum) RectangleSum(l, r, upper int) int {
	res := 0
	for d := 0; d < f.maxLog; d++ {
		if upper>>(f.maxLog-d-1)&1 == 1 {
			res += f.data[d].Sum(f.mat[d].Count(0, l), f.mat[d].Count(0, r))
			l = f.mat[d].Count(1, l) + f.zs[d]
			r = f.mat[d].Count(1, r) + f.zs[d]
		} else {
			l = f.mat[d].Count(0, l)
			r = f.mat[d].Count(0, r)
		}
	}
	return res
}
