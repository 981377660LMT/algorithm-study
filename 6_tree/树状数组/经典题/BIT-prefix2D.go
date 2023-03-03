// 二维树状数组(树套树)
// https://atcoder.jp/contests/abc266/submissions/34262996 (树状数组套树状数组)
// https://www.luogu.com.cn/problem/P3810 (三维偏序)

package main

import (
	"fmt"
	"sort"
)

func main() {
	xs := []int{99988124, 1223, 12345678, 123}
	ys := []int{100012, 2345, 1212, 1111}
	values := []int{100012, 2345, 1212, 1111}

	// 离散化
	sortedX, mp1 := sortedSet(xs)
	sortedY, mp2 := sortedSet(ys)

	bit2d := NewBIT2D(len(sortedX), len(sortedY))
	for i := 0; i < len(xs); i++ {
		bit2d.Update(mp1[xs[i]], mp2[ys[i]], values[i])
	}
	fmt.Println(bit2d.Query(100000005, 12345678))
}

func sortedSet(nums []int) (sorted []int, rank map[int]int) {
	set := make(map[int]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	sorted = make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	rank = make(map[int]int, len(sorted))
	for i, v := range sorted {
		rank[v] = i
	}
	return
}

// !幺元为0, 二元操作为op
var op = max

type BIT2D struct {
	h, w int
	data []*bit1D
}

// 二维树状数组,初始化height个一维树状数组,每个一维树状数组管理的纵坐标长度为width.
//  需要提前将所有点的坐标离散化.
func NewBIT2D(height, width int) *BIT2D {
	bits := make([]*bit1D, height)
	for i := range bits {
		bits[i] = newBIT(width)
	}
	return &BIT2D{height, width, bits}
}

// 单点更新(x,y)处的元素,x和y都是离散化后的坐标.
// 0 <= x < h, 0 <= y < w
func (f *BIT2D) Update(x, y, value int) {
	for x++; x <= f.h; x += x & -x {
		f.data[x-1].Update(y, value)
	}
}

// 查询前缀区间 [0,rightX) * [0,rightY) 的值, rightX和rightY都是离散化后的坐标.
// 0 <= rightX <= h, 0 <= rightY <= w
func (f *BIT2D) Query(rightX, rightY int) int {
	res := 0
	if rightX > f.h {
		rightX = f.h
	}
	if rightY > f.w {
		rightY = f.w
	}
	for ; rightX > 0; rightX -= rightX & -rightX {
		res = op(res, f.data[rightX-1].Query(rightY))
	}
	return res
}

type bit1D struct {
	n    int
	data map[int]int
}

func newBIT(n int) *bit1D {
	return &bit1D{n, make(map[int]int)}
}

// 单点更新index处的元素.
// 0 <= index < n
func (f *bit1D) Update(index int, value int) {
	for index++; index <= f.n; index += index & -index {
		f.data[index-1] = op(f.data[index-1], value)
	}
}

// 查询前缀区间 [0,right) 的值.
// 0 <= right <= n
func (f *bit1D) Query(right int) int {
	res := 0
	if right > f.n {
		right = f.n
	}
	for ; right > 0; right -= right & -right {
		res = op(res, f.data[right-1])
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
