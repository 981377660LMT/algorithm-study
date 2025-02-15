// RectangleUnionArea
// 线段树求矩形面积并
// https://leetcode.cn/problems/rectangle-area-ii/solutions/1827022/fen-xi-by-vclip-qdwl/

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	yosupoAreaOfUnionOfRectangles()
}

// https://leetcode.cn/problems/separate-squares-ii/
func separateSquares(squares [][]int) float64 {
	rectangles := make([]Rectangle, 0, len(squares))
	for _, s := range squares {
		x, y, l := s[0], s[1], s[2]
		xl, xr, yl, yr := x, x+l, y, y+l
		// 转置，转化为找一条垂直线，使得左边的矩形面积等于右边的矩形面积
		xl, yl = yl, xl
		xr, yr = yr, xr
		rectangles = append(rectangles, Rectangle{xl: xl, xr: xr, yl: yl, yr: yr})
	}

	R := NewRectangleUnionArea(rectangles)
	allSum := 0
	R.EnumerateX(func(xl, xr int, dy int) bool {
		allSum += dy * (xr - xl)
		return false
	})

	half := float64(allSum) / 2
	res := 0.0
	preSum := 0.0
	R.EnumerateX(func(xl, xr int, dy int) bool {
		dx := xr - xl
		curSum := float64(dx * dy)
		if dy > 0 && preSum+curSum >= half {
			res = float64(xl) + (half-preSum)/float64(dy)
			return true
		}
		preSum += curSum
		return false
	})
	return res
}

// https://judge.yosupo.jp/problem/area_of_union_of_rectangles
func yosupoAreaOfUnionOfRectangles() {
	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}
	_ = NextInt

	n := NextInt()
	rectangles := make([]Rectangle, 0, n)
	for i := 0; i < n; i++ {
		l, d, r, u := NextInt(), NextInt(), NextInt(), NextInt()
		rectangles = append(rectangles, Rectangle{xl: l, xr: r, yl: d, yr: u})
	}
	R := NewRectangleUnionArea(rectangles)

	res := 0
	R.EnumerateX(func(xl, xr int, dy int) bool {
		res += (xr - xl) * dy
		return false
	})
	fmt.Fprintln(out, res)
}

const INF = 1e18

type Rectangle struct{ xl, xr, yl, yr int }

type RectangleUnionArea struct {
	xs, ys         []int
	orderX, orderY []int32
	rankY          []int32
}

func NewRectangleUnionArea(rectangles []Rectangle) *RectangleUnionArea {
	res := &RectangleUnionArea{}
	if len(rectangles) == 0 {
		return res
	}
	res.build(rectangles)
	return res
}

// O(nlogn).
func (rua *RectangleUnionArea) EnumerateX(consumer func(xl, xr int, dy int) (shouldBreak bool)) {
	m := int32(len(rua.xs))
	if m == 0 {
		return
	}
	xs, ys, orderX, rankY := rua.xs, rua.ys, rua.orderX, rua.rankY
	seg := newLazySegTree32(m-1, func(i int32) E { return E{min: 0, minCount: ys[i+1] - ys[i]} })
	total := ys[m-1] - ys[0]
	for i := int32(0); i < m-1; i++ {
		var delta int
		if orderX[i]&1 == 0 {
			delta = 1
		} else {
			delta = -1
		}
		k := orderX[i] >> 1
		seg.Update(rankY[k<<1], rankY[k<<1|1], delta)
		yInfo := seg.QueryAll()
		dy := total
		if yInfo.min == 0 {
			dy -= yInfo.minCount
		}
		shouldBreak := consumer(xs[i], xs[i+1], dy)
		if shouldBreak {
			break
		}
	}
}

func (rua *RectangleUnionArea) build(rectangles []Rectangle) {
	m := int32(len(rectangles)) << 1
	xs, ys := make([]int, 0, m), make([]int, 0, m)
	for _, r := range rectangles {
		xs = append(xs, r.xl, r.xr)
		ys = append(ys, r.yl, r.yr)
	}
	orderX := argSort(m, func(i, j int32) bool { return xs[i] < xs[j] })
	orderY := argSort(m, func(i, j int32) bool { return ys[i] < ys[j] })
	rankY := make([]int32, m)
	for i := int32(0); i < m; i++ {
		rankY[orderY[i]] = i
	}
	xs = reArrage(xs, orderX)
	ys = reArrage(ys, orderY)
	rua.xs, rua.ys, rua.orderX, rua.orderY, rua.rankY = xs, ys, orderX, orderY, rankY
}

// RangeAddRangeMinMinCount

type E = struct {
	min      int // coverTimes
	minCount int // length
}
type Id = int

func (*lazySegTree32) e() E   { return E{min: INF} }
func (*lazySegTree32) id() Id { return 0 }
func (*lazySegTree32) op(left, right E) E {
	if left.min < right.min {
		return left
	}
	if right.min < left.min {
		return right
	}
	left.minCount += right.minCount
	return left
}
func (*lazySegTree32) mapping(f Id, g E, size int) E {
	g.min += f
	return g
}
func (*lazySegTree32) composition(f, g Id) Id {
	return f + g
}

type lazySegTree32 struct {
	n    int32
	size int32
	log  int32
	data []E
	lazy []Id
}

func newLazySegTree32(n int32, f func(int32) E) *lazySegTree32 {
	tree := &lazySegTree32{}
	tree.n = n
	tree.log = int32(bits.Len32(uint32(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := int32(0); i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

func (tree *lazySegTree32) QueryAll() E {
	return tree.data[1]
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *lazySegTree32) Update(left, right int32, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := int32(1); i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

func (tree *lazySegTree32) pushUp(root int32) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *lazySegTree32) pushDown(root int32) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *lazySegTree32) propagate(root int32, f Id) {
	size := 1 << (tree.log - int32((bits.Len32(uint32(root)) - 1)) /**topbit**/)
	tree.data[root] = tree.mapping(f, tree.data[root], size)
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func argSort(n int32, less func(i, j int32) bool) (order []int32) {
	order = make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return
}

func reArrage[T any](nums []T, order []int32) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}
