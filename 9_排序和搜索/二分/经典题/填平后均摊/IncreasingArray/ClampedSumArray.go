// 要求数组有序.
// 如果数组无序，使用基于waveletMatrix实现的`RangeClampedSumOnline`.
//
// api:
// SumWithMin(start, end int32, min int) int
// SumWithMax(start, end int32, max int) int
// SumWithMinAndMax(start, end int32, min, max int) int
//
// CountRange(start, end int32, min, max int) int
// SumRange(start, end int32, min, max int) int
// CountAndSumRange(start, end int32, min, max int) (int, int)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// demo()

	abc373_e()
}

func demo() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	A := NewClampedSumArray(arr)
	fmt.Println(A.SumWithMinAndMax(0, 10, 2, 5))
}

// E - How to Win the Election (abc373 E) 投票问题
// https://atcoder.jp/contests/abc373/tasks/abc373_e
//
// n 个候选人， k 张票。
// 最终票数最多的 m 个候选人获胜，同票数的都会获胜，因此可能获胜的可能会超过 m 个 。
// 现已知部分投票情况。问每一个候选人，需要给他至少多少张票，才能使得他一定会获胜，
// 即无论剩余票数如何分配，他始终都是票数前 m 多的。
func abc373_e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M, K int
	fmt.Fscan(in, &N, &M, &K)
	A := make([]int, N)
	aSum := 0
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &A[i])
		aSum += A[i]
	}
	upper := K - aSum

	if M == N {
		for i := 0; i < N; i++ {
			fmt.Fprint(out, 0, " ")
		}
		return
	}

	sortedA := append(A[:0:0], A...)
	sort.Ints(sortedA)
	S := NewClampedSumArray(sortedA)
	for i := 0; i < N; i++ {
		check := func(mid int) bool { // A[i] 至少需要 mid 张票
			x := A[i] + mid
			sum := S.SumWithMax(N-M, N, x+1) // >=x+1
			if sortedA[N-M] <= A[i] {        // 需要排除 A[i] 本身
				sum -= A[i]
				sum += sortedA[N-M-1]
			}
			need := (x+1)*M - sum
			remain := upper - mid
			return need > remain
		}

		ok := false
		left, right := 0, upper
		for left <= right {
			mid := (left + right) / 2
			if check(mid) {
				right = mid - 1
				ok = true
			} else {
				left = mid + 1
			}
		}

		if ok {
			fmt.Fprint(out, left, " ")
		} else {
			fmt.Fprint(out, -1, " ")
		}
	}
}

type ClampedSumArray struct {
	Arr    []int
	Presum []int
}

func NewClampedSumArray(increasingArray []int) *ClampedSumArray {
	if !sort.IntsAreSorted(increasingArray) {
		panic("input array should be increasing")
	}
	presum := make([]int, len(increasingArray)+1)
	for i := 0; i < len(increasingArray); i++ {
		presum[i+1] = presum[i] + increasingArray[i]
	}
	return &ClampedSumArray{Arr: increasingArray, Presum: presum}
}

// [min, ?)
func (a *ClampedSumArray) SumWithMin(start, end int, v int) int {
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0
	}
	lessCount := sort.SearchInts(a.Arr[start:end], v)
	largerSum := a.Presum[end] - a.Presum[start+lessCount]
	return v*lessCount + largerSum
}

// (?, max]
func (a *ClampedSumArray) SumWithMax(start, end int, v int) int {
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0
	}
	lessCount := sort.SearchInts(a.Arr[start:end], v)
	lessSum := a.Presum[start+lessCount] - a.Presum[start]
	return lessSum + v*(end-start-lessCount)
}

// [min, max]
func (a *ClampedSumArray) SumWithMinAndMax(start, end int, min, max int) int {
	if min > max {
		return 0
	}
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0
	}
	posLow := sort.SearchInts(a.Arr[start:end], min)
	posUp := sort.SearchInts(a.Arr[start:end], max)
	return a.Presum[start+posUp] - a.Presum[start+posLow] + min*posLow + max*(end-start-posUp)
}

// !WaveletMatrixLike Api

// [start,end) x [y1,y2) 中的数的个数.
func (a *ClampedSumArray) CountRange(start, end int, y1, y2 int) int {
	count, _ := a.CountAndSumRange(start, end, y1, y2)
	return count
}

// [start,end) x [y1,y2) 中的数的和.
func (a *ClampedSumArray) SumRange(start, end int, y1, y2 int) int {
	_, sum := a.CountAndSumRange(start, end, y1, y2)
	return sum
}

// [start,end) x [y1,y2) 中的数的个数、和.
func (a *ClampedSumArray) CountAndSumRange(start, end int, y1, y2 int) (int, int) {
	if y1 >= y2 {
		return 0, 0
	}
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0, 0
	}
	nums := a.Arr[start:end]
	left := sort.SearchInts(nums, y1)
	right := sort.SearchInts(nums, y2)
	return right - left, a.Presum[start+right] - a.Presum[start+left]
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含，使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
