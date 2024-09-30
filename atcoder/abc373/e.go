package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"sort"
	"strconv"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N, M, K := io.NextInt(), io.NextInt(), io.NextInt()
	A := make([]int, N)
	aSum := 0
	for i := 0; i < N; i++ {
		A[i] = io.NextInt()
		aSum += A[i]
	}
	upper := K - aSum

	if M == N {
		for i := 0; i < N; i++ {
			io.Print(0, " ")
		}
		return
	}

	sortedA := make([]int, N)
	copy(sortedA, A)
	sort.Ints(sortedA)
	presum := make([]int, N+1)
	presum[0] = sortedA[0]
	for i := 1; i < N; i++ {
		presum[i] = presum[i-1] + sortedA[i]
	}

	sl := NewIncreasingArray(sortedA)
	for i := 0; i < N; i++ {
		check := func(mid int) (ok bool) {
			x := A[i] + mid
			remain := upper - mid
			// >=x+1
			sum := sl.SumWithUpClampRange(x+1, N-M, N)
			if sortedA[N-M] <= A[i] {
				sum -= A[i]
				sum += sortedA[N-M-1]
			}
			diff := (x+1)*M - sum
			ok = diff > remain
			return
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
			io.Print(left, " ")
		} else {
			io.Print(-1, " ")
		}
	}
}

type IncreasingArray struct {
	Arr    []int
	Presum []int
}

func NewIncreasingArray(increasingArray []int) *IncreasingArray {
	if !sort.IntsAreSorted(increasingArray) {
		panic("input array should be increasing")
	}
	presum := make([]int, len(increasingArray)+1)
	for i := 0; i < len(increasingArray); i++ {
		presum[i+1] = presum[i] + increasingArray[i]
	}
	return &IncreasingArray{Arr: increasingArray, Presum: presum}
}

// 每次选取最矮的矩形中编号最小的，并把它的高度+1
// 返回第K次选取的矩形的高度和编号.
// k>=1.
func (a *IncreasingArray) Increase(k int) (value, pos int) {
	if k <= 0 {
		panic("k should be >=1")
	}
	// !二分无法与哪个数齐平
	right := MaxRight(0, func(r int) bool { return a.Arr[r-1]*r-a.Presum[r] < k }, len(a.Arr))
	filled := a.Arr[right-1]*right - a.Presum[right]
	remain := k - filled
	div, mod := remain/right, remain%right
	value = a.Arr[right-1] + div
	if mod > 0 {
		value++
	}
	pos = mod - 1
	if pos < 0 {
		pos += right
	}
	return
}

func (a *IncreasingArray) IncreaseForMin(k int) int {
	if k <= 0 {
		return a.Arr[0]
	}
	right := MaxRight(0, func(r int) bool { return a.Arr[r-1]*r-a.Presum[r] < k }, len(a.Arr))
	filled := a.Arr[right-1]*right - a.Presum[right]
	remain := k - filled
	return a.Arr[right-1] + remain/right
}

// 每次选取最矮的矩形中编号最小的，并把它的高度+1
// 返回操作后的数组.
func (a *IncreasingArray) IncreaseForArray(k int) []int {
	if k <= 0 {
		return a.Arr
	}
	right := MaxRight(0, func(r int) bool { return a.Arr[r-1]*r-a.Presum[r] < k }, len(a.Arr))
	filled := a.Arr[right-1]*right - a.Presum[right]
	remain := k - filled
	div, mod := remain/right, remain%right
	base := a.Arr[right-1] + div
	res := append(a.Arr[:0:0], a.Arr...)
	for i := 0; i < right; i++ {
		res[i] = base
		if i < mod {
			res[i]++
		}
	}
	return res
}

// 每次选取最高的矩形中编号最小的，并把它的高度-1
// 返回第K次选取的矩形的高度和编号.
// k>=1.
func (a *IncreasingArray) Decrease(k int) (value, pos int) {
	if k <= 0 {
		panic("k should be >=1")
	}
	left := MinLeft(
		len(a.Arr),
		func(l int) bool {
			return (a.Presum[len(a.Arr)]-a.Presum[l])-a.Arr[l]*(len(a.Arr)-l) < k
		},
		0,
	)
	filled := a.Presum[len(a.Arr)] - a.Presum[left] - a.Arr[left]*(len(a.Arr)-left)
	remain := k - filled
	div, mod := remain/(len(a.Arr)-left), remain%(len(a.Arr)-left)
	value = a.Arr[left] - div
	if mod > 0 {
		value--
	}
	pos = mod - 1
	if pos < 0 {
		pos += len(a.Arr) - left
	}
	pos += left
	return
}

func (a *IncreasingArray) DecreaseForMax(k int) int {
	if k <= 0 {
		return a.Arr[len(a.Arr)-1]
	}
	left := MinLeft(
		len(a.Arr),
		func(l int) bool {
			return (a.Presum[len(a.Arr)]-a.Presum[l])-a.Arr[l]*(len(a.Arr)-l) < k
		},
		0,
	)
	filled := a.Presum[len(a.Arr)] - a.Presum[left] - a.Arr[left]*(len(a.Arr)-left)
	remain := k - filled
	return a.Arr[left] - remain/(len(a.Arr)-left)
}

// 每次选取最高的矩形中编号最小的，并把它的高度-1
// 返回操作后的数组.
func (a *IncreasingArray) DecreaseForArray(k int) []int {
	if k <= 0 {
		return a.Arr
	}
	left := MinLeft(
		len(a.Arr),
		func(l int) bool {
			return (a.Presum[len(a.Arr)]-a.Presum[l])-a.Arr[l]*(len(a.Arr)-l) < k
		},
		0,
	)
	filled := a.Presum[len(a.Arr)] - a.Presum[left] - a.Arr[left]*(len(a.Arr)-left)
	remain := k - filled
	div, mod := remain/(len(a.Arr)-left), remain%(len(a.Arr)-left)
	base := a.Arr[left] - div
	res := append(a.Arr[:0:0], a.Arr...)
	for i := left; i < len(a.Arr); i++ {
		res[i] = base
		if i-left < mod {
			res[i]--
		}
	}
	return res
}

// 求所有数与v取min的和.
func (a *IncreasingArray) SumWithUpClamp(v int) int {
	pos := sort.SearchInts(a.Arr, v)
	return a.Presum[pos] + (len(a.Arr)-pos)*v
}

func (a *IncreasingArray) SumWithUpClampRange(v int, start, end int) int {
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

// 求所有数与v取max的和.
func (a *IncreasingArray) SumWithLowClamp(v int) int {
	pos := sort.SearchInts(a.Arr, v)
	return pos*v + a.Presum[len(a.Arr)] - a.Presum[pos]
}

func (a *IncreasingArray) SumWithLowClampRange(v int, start, end int) int {
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

// 求所有数与v的绝对值差的和.
func (a *IncreasingArray) DiffSum(v int) int {
	pos := sort.SearchInts(a.Arr, v)
	n := len(a.Arr)
	leftSum := v*pos - a.Presum[pos]
	rightSum := a.Presum[n] - a.Presum[pos] - v*(n-pos)
	return leftSum + rightSum
}

func (a *IncreasingArray) DiffSumRange(v int, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0
	}
	presum := a.Presum
	pos := sort.SearchInts(a.Arr, v)
	if pos <= start {
		return (presum[end] - presum[start]) - v*(end-start)
	}
	if pos >= end {
		return v*(end-start) - (presum[end] - presum[start])
	}
	leftSum := v*(pos-start) - (presum[pos] - presum[start])
	rightSum := presum[end] - presum[pos] - v*(end-pos)
	return leftSum + rightSum
}

// !WaveletMatrixLike Api

// [start,end) x [y1,y2) 中的数的个数.
func (a *IncreasingArray) CountRange(start, end int, y1, y2 int) int {
	count, _ := a.CountAndSumRange(start, end, y1, y2)
	return count
}

// [start,end) x [y1,y2) 中的数的和.
func (a *IncreasingArray) SumRange(start, end int, y1, y2 int) int {
	_, sum := a.CountAndSumRange(start, end, y1, y2)
	return sum
}

// [start,end) x [y1,y2) 中的数的个数、和.
func (a *IncreasingArray) CountAndSumRange(start, end int, y1, y2 int) (int, int) {
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}
