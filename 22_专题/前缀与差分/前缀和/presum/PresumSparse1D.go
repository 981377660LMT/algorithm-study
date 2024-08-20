package main

import (
	"fmt"
	"sort"
)

func main() {
	points := [][2]int{{1, 2}, {3, 4}, {1, 2}, {3, 4}, {500000, 500000}, {-500000, -500000}}
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumSparse1D(e, op, inv)
	S.Build(int32(len(points)), func(i int32) (int, int) { return points[i][0], points[i][1] })
	fmt.Println(S.Query(0, 10))
	fmt.Println(S.Query(3, 5))
	fmt.Println(S.Query(-500000, 500000))
	fmt.Println(S.Query(-500000, 500001098765876))
}

type PresumSparse1D[E any] struct {
	n      int32
	presum []E
	origin []int
	e      func() E
	op     func(a, b E) E
	inv    func(a E) E
}

func NewPresumSparse1D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *PresumSparse1D[E] {
	return &PresumSparse1D[E]{e: e, op: op, inv: inv}
}

// 一维数轴上的前缀和.
//
//	f: 返回第i个元素的横坐标值和权值.
func (p *PresumSparse1D[E]) Build(n int32, f func(i int32) (x int, e E)) {
	xs, es := make([]int, n), make([]E, n)
	for i := int32(0); i < n; i++ {
		xs[i], es[i] = f(i)
	}
	newXs, originX := discretize1D(xs)
	presum := make([]E, len(originX)+1)
	for i := range presum {
		presum[i] = p.e()
	}
	for i := int32(0); i < n; i++ {
		presum[newXs[i]+1] = p.op(presum[newXs[i]+1], es[i])
	}
	for i := 1; i < len(presum); i++ {
		presum[i] = p.op(presum[i], presum[i-1])
	}
	p.n = n
	p.presum = presum
	p.origin = originX
}

func (p *PresumSparse1D[E]) Query(start, end int) E {
	if start >= end {
		return p.e()
	}
	newStart, newEnd := p.compress(start), p.compress(end)
	return p.op(p.inv(p.presum[newStart]), p.presum[newEnd])
}

func (p *PresumSparse1D[E]) compress(x int) int32 {
	return int32(sort.SearchInts(p.origin, x))
}

func discretize1D(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}
