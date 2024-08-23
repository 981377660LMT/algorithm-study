// PresumWithBisect/PresumWithMinLeft/PresumWithMaxLeft/PresumBisector

package main

import "fmt"

func main() {
	// [0,1,2,3,4,5,6,7,8,9]
	P := NewPresumBisector(10, func(i int32) int { return int(i) })
	fmt.Println(P.Query(0, 10))                                                     // 45
	fmt.Println(P.MaxRight(0, func(sum int, right int32) bool { return sum < 10 })) // 10
	fmt.Println(P.MinLeft(10, func(sum int, left int32) bool { return sum < 10 }))  // 0
}

// 带有二分的前缀和，要求元素为非负数.
type PresumBisector struct {
	n      int32
	presum []int
}

func NewPresumBisector(n int32, f func(i int32) int) *PresumBisector {
	presum := make([]int, n+1)
	for i := int32(1); i <= n; i++ {
		presum[i] = presum[i-1] + f(i-1)
	}
	return &PresumBisector{n: n, presum: presum}
}

func (p *PresumBisector) Query(start, end int32) int {
	if start <= 0 {
		start = 0
	}
	if end > p.n {
		end = p.n
	}
	if start >= end {
		return 0
	}
	return p.presum[end] - p.presum[start]
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (p *PresumBisector) MaxRight(left int32, check func(sum int, right int32) bool) int32 {
	if left >= p.n {
		return p.n
	}
	ok, ng := left, p.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(p.Query(left, mid), mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (p *PresumBisector) MinLeft(right int32, check func(sum int, left int32) bool) int32 {
	if right <= 0 {
		return 0
	}
	ok, ng := right, int32(-1)
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(p.Query(mid, right), mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}
