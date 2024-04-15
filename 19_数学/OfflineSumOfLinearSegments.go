// Offline sum of linear segments （区間一次関数加算クエリの累積和を用いたオフライン処理）
// 区间一次函数叠加求和/区间等差数列求和

// AddLinear(start, end, k, b):
//
//	O(1) 为数组 [start, end) 的每个位置 i 加上 k*(i-start)+b.
//
// Run():
//
//	 O(n) 返回数组的累加和.
//
//	OfflineRangeLinearAddRangeSum

package main

import "fmt"

func main() {
	nums := NewOfflineSumOfLinearSegments(5)
	nums.AddLinear(0, 5, 1, 1)
	fmt.Println(nums.Run())
	nums.AddLinear(0, 2, 10, 10)
	fmt.Println(nums.Run())
}

type OfflineSumOfLinearSegments struct {
	n    int
	ret  []int
	upds [][]int // (start, end, b, k)
}

func NewOfflineSumOfLinearSegments(n int) *OfflineSumOfLinearSegments {
	return &OfflineSumOfLinearSegments{n: n, ret: make([]int, n)}
}

// Add `k*(x-start)+b` to A[x] for x in [start, end).
//  0 <= start <= end <= N.
func (o *OfflineSumOfLinearSegments) AddLinear(start, end, k, b int) {
	if start < end {
		o.upds = append(o.upds, []int{start, end, b, k})
	}
}

// Run offline queries.
//  Return A.
func (o *OfflineSumOfLinearSegments) Run() []int {
	if len(o.upds) == 0 {
		return o.ret
	}

	tmp := make([]int, o.n+1)
	l, r := 0, 0
	b, k := 0, 0
	for i := 0; i < len(o.upds); i++ {
		l, r, b, k = o.upds[i][0], o.upds[i][1], o.upds[i][2], o.upds[i][3]
		tmp[l+1] += k
		tmp[r] -= k
	}

	for i := 0; i < o.n; i++ {
		tmp[i+1] += tmp[i]
	}

	for i := 0; i < len(o.upds); i++ {
		l, r, b, k = o.upds[i][0], o.upds[i][1], o.upds[i][2], o.upds[i][3]
		tmp[l] += b
		tmp[r] -= b + (r-l-1)*k
	}

	o.upds = o.upds[:0]
	for i := 0; i < o.n; i++ {
		tmp[i+1] += tmp[i]
		o.ret[i] += tmp[i]
	}

	return o.ret
}
