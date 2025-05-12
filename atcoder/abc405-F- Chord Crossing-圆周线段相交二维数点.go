// F - Chord Crossing(圆周线段相交二维数点)
// https://atcoder.jp/contests/abc405/tasks/abc405_f
//
// 在一个圆周上等间隔排列了 **2N** 个点，这些点从某个点开始按顺时针方向依次编号为 `1, 2, ..., 2N`。
//
// 有 **M** 条线段 `1, 2, ..., M`，其中第 `i` 条线段连接了点 `A[i]` 和点 `B[i]`。
// 以下条件成立：
// - `A[i]` 和 `B[i]` 是 **不同的偶数**。
// - 任意两条线段都 **不共享端点**。
//
// 接下来有 **Q** 个查询。对于第 `j` 个查询，给出两个 **不同的奇数** `C[j]` 和 `D[j]`，请回答以下问题：
//
// > 在给定的 **M** 条线段中，有多少条线段与连接点 `C[j]` 和 `D[j]` 的线段有公共点？
//
// 1. 二维数点 https://atcoder.jp/contests/abc405/editorial/13017
// 2. 转换成区间 -> 线段树 https://atcoder.jp/contests/abc405/editorial/13011
// 3. 转化成树上距离查询 不好理解

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M int
	fmt.Fscan(in, &N, &M)
	A, B := make([]int, M), make([]int, M)
	for i := 0; i < M; i++ {
		fmt.Fscan(in, &A[i], &B[i])
		A[i]--
		B[i]--
	}

	R := NewRectangleSum()
	for i := 0; i < M; i++ {
		R.AddPoint(A[i], B[i], 1)
	}
	R.Build()

	query := func(c, d int) int {
		// 两条线段(a,b)和(c,d)相交的条件：
		// 1. a < c < b < d 或
		// 2. c < a < d < b

		// 计算满足a < c < b < d的线段数量
		count1 := R.Query(0, c, c, d)

		// 计算满足c < a < d < b的线段数量
		count2 := R.Query(c, d, d, 2*N)

		return count1 + count2
	}

	var Q int
	fmt.Fscan(in, &Q)
	for i := 0; i < Q; i++ {
		var c, d int
		fmt.Fscan(in, &c, &d)
		c--
		d--
		fmt.Fprintln(out, query(c, d))
	}
}

// RectangleSumStatic
type StaticRectangleSum struct {
	points [][3]int
	xs     []int
	ys     []int
	wm     *waveletMatrix
}

func NewRectangleSum() *StaticRectangleSum {
	return &StaticRectangleSum{
		points: [][3]int{},
	}
}

func (s *StaticRectangleSum) AddPoint(x, y, w int) {
	s.points = append(s.points, [3]int{x, y, w})
}

func (s *StaticRectangleSum) Build() {
	sort.Slice(s.points, func(i, j int) bool {
		return s.points[i][0] < s.points[j][0]
	})
	n := len(s.points)
	xs, ys, ws := make([]int, n), make([]int, n), make([]int, n)
	for i, p := range s.points {
		xs[i], ys[i], ws[i] = p[0], p[1], p[2]
	}
	s.xs = xs

	set := make(map[int]struct{}, len(ys))
	for _, y := range ys {
		set[y] = struct{}{}
	}
	sortedSet := make([]int, 0, len(set))
	for y := range set {
		sortedSet = append(sortedSet, y)
	}
	sort.Ints(sortedSet)
	s.ys = sortedSet

	comp := make(map[int]int, len(sortedSet))
	for i, y := range sortedSet {
		comp[y] = i
	}

	newYs := make([]int, len(ys))
	for i, y := range ys {
		newYs[i] = comp[y]
	}

	maxLog := bits.Len(uint(len(sortedSet)))
	s.wm = newWaveletMatrix(newYs, ws, maxLog)
}

// 求矩形x1<=x<x2,y1<=y<y2的权值和 注意是左闭右开
func (s *StaticRectangleSum) Query(x1, x2, y1, y2 int) int {
	return s.rectSum(x1, x2, y2) - s.rectSum(x1, x2, y1)
}

func (s *StaticRectangleSum) rectSum(left, right, upper int) int {
	left = sort.SearchInts(s.xs, left)
	right = sort.SearchInts(s.xs, right)
	upper = sort.SearchInts(s.ys, upper)
	return s.wm.RectSum(left, right, upper)
}

func newWaveletMatrix(ys, ws []int, maxLog int) *waveletMatrix {
	n := len(ys)
	mat := make([]*bitVector, 0, maxLog)
	zs := make([]int, 0, maxLog)
	data := make([][]int, maxLog)
	for i := range data {
		data[i] = make([]int, n+1)
	}

	order := make([]int, n)
	for i := range order {
		order[i] = i
	}

	for d := maxLog - 1; d >= 0; d-- {
		vec := newBitVector(n + 1)
		ls, rs := make([]int, 0, n), make([]int, 0, n)
		for i, val := range order {
			if (ys[val]>>uint(d))&1 == 1 {
				rs = append(rs, val)
				vec.Set(i)
			} else {
				ls = append(ls, val)
			}
		}
		vec.Build()
		mat = append(mat, vec)
		zs = append(zs, len(ls))
		order = append(ls, rs...)
		for i, val := range order {
			data[maxLog-d-1][i+1] = data[maxLog-d-1][i] + ws[val]
		}
	}

	return &waveletMatrix{
		n:      n,
		maxLog: maxLog,
		mat:    mat,
		zs:     zs,
		data:   data,
	}
}

type waveletMatrix struct {
	n      int
	maxLog int
	mat    []*bitVector
	zs     []int
	data   [][]int
}

func (w *waveletMatrix) RectSum(left, right, upper int) int {
	res := 0
	for d := 0; d < w.maxLog; d++ {
		if (upper>>(w.maxLog-d-1))&1 == 1 {
			res += w.data[d][w.mat[d].Count(0, right)]
			res -= w.data[d][w.mat[d].Count(0, left)]
			left = w.mat[d].Count(1, left) + w.zs[d]
			right = w.mat[d].Count(1, right) + w.zs[d]
		} else {
			left = w.mat[d].Count(0, left)
			right = w.mat[d].Count(0, right)
		}
	}
	return res
}

type bitVector struct {
	n     int
	block []int
	sum   []int
}

func newBitVector(n int) *bitVector {
	blockCount := (n + 63) >> 6
	return &bitVector{
		n:     n,
		block: make([]int, blockCount+1),
		sum:   make([]int, blockCount+1),
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
