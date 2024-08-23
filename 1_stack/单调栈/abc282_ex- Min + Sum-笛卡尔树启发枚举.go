// !O(nlogn)笛卡尔树启发式枚举，每次枚举较小的一半
//
// !一般的题目是以中点为mid分治，这里是以最小值为mid分治(就是笛卡尔树分治)
// 这样可以确定区间整个最小值为A[mid].

package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ex - Min + Sum (启发式枚举，枚举较小的一半)
// https://atcoder.jp/contests/abc282/tasks/abc282_h
// 给定两个长为n的数组，求区间个数，满足`A的子数组最小值+B的子数组之和<=k`
// n<=2e5，1<=A[i],B[i]<=1e9，1<=k<=1e18
// O(nlognlogn)
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	var k int
	fmt.Fscan(in, &n, &k)
	A, B := make([]int, n), make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &B[i])
	}

	leftMinA, rightMinA := GetRange32(A, false, false, true)
	presumB := NewPresumBisector(n, func(i int32) int { return B[i] })

	// 笛卡尔树的分治结果，mid为区间[start,end)的最小值.
	// 子数组必须包含mid，所以考虑枚举左端点/右端点.
	calc := func(start, mid, end int32) (count int) {
		len1, len2 := mid-start+1, end-mid
		upper := k - A[mid]
		if len1 < len2 {
			// 枚举左端点，看右端点到哪里
			for s := start; s <= mid; s++ {
				right := presumB.MaxRight(s, func(sum int, right int32) bool { return sum <= upper })
				right = min32(right, end)
				count += int(max32(0, right-mid))
			}
		} else {
			// 枚举右端点，看左端点到哪里
			for e := mid + 1; e <= end; e++ {
				left := presumB.MinLeft(e, func(sum int, left int32) bool { return sum <= upper })
				left = max32(left, start)
				count += int(max32(0, mid-left+1))
			}
		}
		return
	}

	res := 0
	for i := int32(0); i < n; i++ {
		left, right := leftMinA[i], rightMinA[i]+1
		res += calc(left, i, right)
	}
	fmt.Fprintln(out, res)
}

// 求每个元素作为最值的影响范围(区间)
func GetRange32(nums []int, isMax, isLeftStrict, isRightStrict bool) (leftMost, rightMost []int32) {
	compareLeft := func(stackValue, curValue int) bool {
		if isLeftStrict && isMax {
			return stackValue <= curValue
		} else if isLeftStrict && !isMax {
			return stackValue >= curValue
		} else if !isLeftStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	compareRight := func(stackValue, curValue int) bool {
		if isRightStrict && isMax {
			return stackValue <= curValue
		} else if isRightStrict && !isMax {
			return stackValue >= curValue
		} else if !isRightStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	n := int32(len(nums))
	leftMost, rightMost = make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		rightMost[i] = n - 1
	}

	stack := []int32{}
	for i := int32(0); i < n; i++ {
		for len(stack) > 0 && compareRight(nums[stack[len(stack)-1]], nums[i]) {
			rightMost[stack[len(stack)-1]] = i - 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && compareLeft(nums[stack[len(stack)-1]], nums[i]) {
			leftMost[stack[len(stack)-1]] = i + 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	return
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
