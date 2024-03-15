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

// 0 x y w : 向平面上添加一个点(x,y)并且权值为w
// 1 left down right up : 查询矩形 left<=x<right down<=y<up 内的点权和
// n,q<=1e5 0<=xi,yi<=1e9
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	points := make([][]int32, 0, n)
	for i := int32(0); i < n; i++ {
		var x, y, w int32
		fmt.Fscan(in, &x, &y, &w)
		points = append(points, []int32{x, y, w})
	}
	queries := make([][]int32, 0, q)
	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var x, y, w int32
			fmt.Fscan(in, &x, &y, &w)
			queries = append(queries, []int32{op, x, y, w})
		} else {
			var left, down, right, up int32
			fmt.Fscan(in, &left, &down, &right, &up)
			queries = append(queries, []int32{op, left, down, right, up})
		}
	}

	for _, query := range queries {
		if query[0] == 0 {
			x, y := query[1], query[2]
			points = append(points, []int32{x, y, 0}) // 0: 预先添加可能出现的点
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

type pair = struct{ x, y int32 }
type CompressedPointAddRectangleSum struct {
	xs, ys []int32
	idxs   map[pair]int32
	mat    *pointAddRectangleSum
}

// 二维矩形区域计数
//
//	需要预先添加可能出现的所有点 (x,y,weight)
func NewPointAddRectangleSum(points [][]int32) *CompressedPointAddRectangleSum {
	res := &CompressedPointAddRectangleSum{}

	points = append(points[:0:0], points...)
	sort.Slice(points, func(i, j int) bool {
		return points[i][0] < points[j][0]
	})

	xs, ys, ws := make([]int32, len(points)), make([]int32, len(points)), make([]int32, len(points))
	for i, p := range points {
		xs[i], ys[i], ws[i] = p[0], p[1], p[2]
	}
	res.xs = xs

	set := make(map[int32]struct{}, len(ys))
	for _, y := range ys {
		set[y] = struct{}{}
	}
	sortedSet := make([]int32, 0, len(set))
	for y := range set {
		sortedSet = append(sortedSet, y)
	}
	sort.Slice(sortedSet, func(i, j int) bool { return sortedSet[i] < sortedSet[j] })
	res.ys = sortedSet

	type pair = struct{ x, y int32 }
	ids := make(map[pair]int32, len(points))
	for i, p := range points {
		x, y := p[0], p[1]
		ids[pair{x, y}] = int32(i)
	}
	res.idxs = ids

	comp := make(map[int32]int32, len(sortedSet))
	for i, v := range sortedSet {
		comp[v] = int32(i)
	}

	newYs := make([]int32, len(ys))
	for i, v := range ys {
		newYs[i] = comp[v]
	}

	maxLog := int32(bits.Len(uint(len(sortedSet))))
	res.mat = newPointAddRectangleSum(newYs, ws, maxLog)
	return res
}

func (f *CompressedPointAddRectangleSum) Add(x, y, w int32) {
	idx := f.idxs[pair{x, y}]
	f.mat.PointAdd(idx, int(w))
}

// 求矩形x1<=x<x2,y1<=y<y2的权值和 注意是左闭右开
func (f *CompressedPointAddRectangleSum) Query(x1, x2, y1, y2 int32) int {
	return f.rectangleSum(x1, x2, y2) - f.rectangleSum(x1, x2, y1)
}

func (f *CompressedPointAddRectangleSum) rectangleSum(l, r, upper int32) int {
	l = int32(sort.Search(len(f.xs), func(i int) bool { return f.xs[i] >= l }))
	r = int32(sort.Search(len(f.xs), func(i int) bool { return f.xs[i] >= r }))
	upper = int32(sort.Search(len(f.ys), func(i int) bool { return f.ys[i] >= upper }))
	return f.mat.RectangleSum(l, r, upper)
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

type binaryIndexedTree struct {
	n   int32
	bit []int
}

func newBinaryIndexedTree(n int32) *binaryIndexedTree {
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

func (f *binaryIndexedTree) Add(i int32, x int) {
	i++
	for i <= f.n {
		f.bit[i] += x
		i += i & -i
	}
}

func (f *binaryIndexedTree) Sum(l, r int32) int {
	if l == 0 {
		return f.sum(r)
	}
	pos, neg := 0, 0
	for r > l {
		pos += f.bit[r]
		r &= r - 1
	}
	for l > r {
		neg += f.bit[l]
		l &= l - 1
	}
	return pos - neg
}

func (f *binaryIndexedTree) sum(i int32) int {
	res := 0
	for i > 0 {
		res += f.bit[i]
		i -= i & -i
	}
	return res
}

type pointAddRectangleSum struct {
	n      int32
	maxLog int32
	arr    []int32
	mat    []*bitVector
	zs     []int32
	data   []*binaryIndexedTree
}

func newPointAddRectangleSum(arr, ws []int32, maxLog int32) *pointAddRectangleSum {
	res := &pointAddRectangleSum{
		n:      int32(len(arr)),
		maxLog: maxLog,
		arr:    arr,
		mat:    make([]*bitVector, 0, maxLog),
		zs:     make([]int32, 0, maxLog),
		data:   make([]*binaryIndexedTree, maxLog),
	}

	n := int32(len(arr))
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}

	for d := maxLog - 1; d >= 0; d-- {
		vec := newBitVector(n + 1)
		ls, rs := []int32{}, []int32{}
		for i, v := range order {
			if ((arr[v] >> d) & 1) == 1 {
				rs = append(rs, v)
				vec.Set(int32(i))
			} else {
				ls = append(ls, v)
			}
		}
		vec.Build()
		res.mat = append(res.mat, vec)
		res.zs = append(res.zs, int32(len(ls)))
		order = append(ls, rs...)
		res.data[int32(len(res.data))-1-d] = newBinaryIndexedTree(n)
		nums := make([]int, n)
		for i, v := range order {
			nums[i] = int(ws[v])
		}
		res.data[int32(len(res.data))-1-d].Build(nums)
	}

	return res
}

func (f *pointAddRectangleSum) PointAdd(k int32, val int) {
	y := f.arr[k]
	for d := int32(0); d < f.maxLog; d++ {
		if y>>(f.maxLog-d-1)&1 == 1 {
			k = f.mat[d].Count(1, k) + f.zs[d]
		} else {
			k = f.mat[d].Count(0, k)
		}
		f.data[d].Add(k, val)
	}
}

func (f *pointAddRectangleSum) RectangleSum(l, r, upper int32) int {
	res := 0
	for d := int32(0); d < f.maxLog; d++ {
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
