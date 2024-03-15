// 静态二维矩形区间计数(二维偏序)
// n<=2e5 xi,yi,wi<=1e9
// 如果离线可以使用cdq分治

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

// https://judge.yosupo.jp/problem/rectangle_sum
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	SC := NewRectangleSum()
	for i := 0; i < n; i++ {
		var x, y, w int32
		fmt.Fscan(in, &x, &y, &w)
		SC.AddPoint(x, y, w)
	}
	SC.Build()

	for i := 0; i < q; i++ {
		// left<=x<right, down<=y<up
		var left, down, right, up int32
		fmt.Fscan(in, &left, &down, &right, &up)
		fmt.Fprintln(out, SC.Query(left, right, down, up))
	}
}

// RectangleSumStatic
type StaticRectangleSum struct {
	points [][3]int32
	xs     []int32
	ys     []int32
	wm     *waveletMatrix
}

func NewRectangleSum() *StaticRectangleSum {
	return &StaticRectangleSum{
		points: [][3]int32{},
	}
}

func (s *StaticRectangleSum) AddPoint(x, y, w int32) {
	s.points = append(s.points, [3]int32{x, y, w})
}

func (s *StaticRectangleSum) Build() {
	sort.Slice(s.points, func(i, j int) bool {
		return s.points[i][0] < s.points[j][0]
	})
	n := len(s.points)
	xs, ys, ws := make([]int32, n), make([]int32, n), make([]int32, n)
	for i, p := range s.points {
		xs[i], ys[i], ws[i] = p[0], p[1], p[2]
	}
	s.xs = xs

	set := make(map[int32]struct{}, len(ys))
	for _, y := range ys {
		set[y] = struct{}{}
	}
	sortedSet := make([]int32, 0, len(set))
	for y := range set {
		sortedSet = append(sortedSet, y)
	}
	sort.Slice(sortedSet, func(i, j int) bool { return sortedSet[i] < sortedSet[j] })
	s.ys = sortedSet

	comp := make(map[int32]int32, len(sortedSet))
	for i, y := range sortedSet {
		comp[y] = int32(i)
	}

	newYs := make([]int32, len(ys))
	for i, y := range ys {
		newYs[i] = comp[y]
	}

	maxLog := bits.Len(uint(len(sortedSet)))
	s.wm = newWaveletMatrix(newYs, ws, int32(maxLog))
}

// 求矩形x1<=x<x2,y1<=y<y2的权值和 注意是左闭右开
func (s *StaticRectangleSum) Query(x1, x2, y1, y2 int32) int {
	return s.rectSum(x1, x2, y2) - s.rectSum(x1, x2, y1)
}

func (s *StaticRectangleSum) rectSum(left, right, upper int32) int {
	left = int32(sort.Search(len(s.xs), func(i int) bool { return s.xs[i] >= left }))
	right = int32(sort.Search(len(s.xs), func(i int) bool { return s.xs[i] >= right }))
	upper = int32(sort.Search(len(s.ys), func(i int) bool { return s.ys[i] >= upper }))
	return s.wm.RectSum(left, right, upper)
}

func newWaveletMatrix(ys, ws []int32, maxLog int32) *waveletMatrix {
	n := int32(len(ys))
	mat := make([]*bitVector, 0, maxLog)
	zs := make([]int32, 0, maxLog)
	data := make([][]int, maxLog)
	for i := range data {
		data[i] = make([]int, n+1)
	}

	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}

	for d := maxLog - 1; d >= 0; d-- {
		vec := newBitVector(n + 1)
		ls, rs := make([]int32, 0, n), make([]int32, 0, n)
		for i, val := range order {
			if (ys[val]>>uint(d))&1 == 1 {
				rs = append(rs, val)
				vec.Set(int32(i))
			} else {
				ls = append(ls, val)
			}
		}
		vec.Build()
		mat = append(mat, vec)
		zs = append(zs, int32(len(ls)))
		order = append(ls, rs...)
		for i, val := range order {
			data[maxLog-d-1][i+1] = data[maxLog-d-1][i] + int(ws[val])
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
	n      int32
	maxLog int32
	mat    []*bitVector
	zs     []int32
	data   [][]int
}

func (w *waveletMatrix) RectSum(left, right, upper int32) int {
	res := 0
	for d := int32(0); d < w.maxLog; d++ {
		if (upper>>(w.maxLog-d-1))&1 == 1 {
			res += int(w.data[d][w.mat[d].Count(0, right)])
			res -= int(w.data[d][w.mat[d].Count(0, left)])
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
	n     int32
	block []int
	sum   []int32
}

func newBitVector(n int32) *bitVector {
	blockCount := (n + 63) >> 6
	return &bitVector{
		n:     n,
		block: make([]int, blockCount+1),
		sum:   make([]int32, blockCount+1),
	}
}

func (f *bitVector) Set(i int32) {
	f.block[i>>6] |= 1 << uint(i&63)
}

func (f *bitVector) Build() {
	for i := 0; i < len(f.block)-1; i++ {
		f.sum[i+1] = f.sum[i] + int32(bits.OnesCount(uint(f.block[i])))
	}
}

// 统计 [0,end) 中 value 的个数
func (f *bitVector) Count(value, end int32) int32 {
	if value == 1 {
		return f.count1(end)
	}
	return end - f.count1(end)
}

func (f *bitVector) count1(k int32) int32 {
	mask := (1 << uint(k&63)) - 1
	return f.sum[k>>6] + int32(bits.OnesCount(uint(f.block[k>>6]&mask)))
}
