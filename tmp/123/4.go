package main

import (
	"math/bits"
	"sort"
)

type StaticRectangleSum struct {
	points [][3]int
	xs     []int
	ys     []int
	wm     *waveletMatrix
}

func NewStaticRectangleCount() *StaticRectangleSum {
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

	set := make(map[int]struct{})
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

type Point = struct{ x, y, w int }
type Query = struct{ lx, rx, ly, ry int }

// Point: (x,y,w)
// Query: [lx,rx) * [ly,ry)
func RectangleSum(points []Point, queries []Query) []int {
	n := len(points)
	q := len(queries)
	res := make([]int, q)
	if n == 0 || q == 0 {
		return res
	}

	sort.Slice(points, func(i, j int) bool { return points[i].y < points[j].y })
	ys := make([]int, 0, n)
	for i := range points {
		if len(ys) == 0 || ys[len(ys)-1] != points[i].y {
			ys = append(ys, points[i].y)
		}
		points[i].y = len(ys) - 1
	}

	type Q struct {
		x    int
		d, u int
		t    bool
		idx  int
	}

	qs := make([]Q, 0, q+q)
	for i := 0; i < q; i++ {
		query := queries[i]
		d := sort.SearchInts(ys, query.ly)
		u := sort.SearchInts(ys, query.ry)
		qs = append(qs, Q{x: query.lx, d: d, u: u, t: false, idx: i})
		qs = append(qs, Q{x: query.rx, d: d, u: u, t: true, idx: i})
	}

	sort.Slice(points, func(i, j int) bool { return points[i].x < points[j].x })
	sort.Slice(qs, func(i, j int) bool { return qs[i].x < qs[j].x })

	j := 0
	bit := newBinaryIndexedTree(len(ys))
	for _, query := range qs {
		for j < n && points[j].x < query.x {
			bit.Apply(points[j].y, points[j].w)
			j++
		}
		if query.t {
			res[query.idx] += bit.ProdRange(query.d, query.u)
		} else {
			res[query.idx] -= bit.ProdRange(query.d, query.u)
		}
	}

	return res
}

type binaryIndexedTree struct {
	n    int
	log  int
	data []int
}

func newBinaryIndexedTree(n int) *binaryIndexedTree {
	return &binaryIndexedTree{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

func newBinaryIndexedTreeFrom(arr []int) *binaryIndexedTree {
	res := newBinaryIndexedTree(len(arr))
	res.build(arr)
	return res
}

func (b *binaryIndexedTree) Apply(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

func (b *binaryIndexedTree) Prod(r int) int {
	res := 0
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

func (b *binaryIndexedTree) ProdRange(l, r int) int {
	return b.Prod(r) - b.Prod(l)
}

func (b *binaryIndexedTree) build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}

func numberOfPairs(points [][]int) int {
	ps := make([]Point, 0)
	for _, p := range points {
		ps = append(ps, Point{p[0], p[1], 1})
	}
	qs := make([]Query, 0)
	for i, p1 := range points { // liupengsay
		x1, y1 := p1[0], p1[1]
		for j, p2 := range points { // 小羊肖恩
			x2, y2 := p2[0], p2[1]
			if i != j && x1 <= x2 && y1 >= y2 {
				qs = append(qs, Query{x1, x2 + 1, y2, y1 + 1})
			}
		}
	}
	r := RectangleSum(ps, qs)
	res := 0
	for _, v := range r {
		if v == 2 {
			res++
		}
	}
	return res
}
