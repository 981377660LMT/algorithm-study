// api:
//  1.Set(index int32, value E) -> O(sqrt(n))
//  2.Query(start, end int32) E -> O(sqrt(n))
//  3.QueryAll() E -> O(sqrt(n))
//  !4.Get(index int32) E -> O(1)
//  5.GetAll() []E -> O(n)
//  6.MaxRight(start int32, predicate func(end int32, sum E) bool) int32
//  7.MinLeft(end int32, predicate func(start int32, sum E) bool) int32

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	// demo()
	// test()
	testTime()
}

func demo() {
	type E = int
	e, op := func() int { return 0 }, func(a, b int) int { return a + b }
	seg := NewSegmentTreeSqrtDecomposition(e, op, 10, func(i int32) int { return int(i) }, -1)
	fmt.Println(seg.GetAll())
	seg.Set(3, 5)
	seg.Set(4, 6)
	fmt.Println(seg.GetAll())
	fmt.Println(seg.Query(3, 5)) // 11
	fmt.Println("----------------")
	fmt.Println(seg.MaxRight(3, func(end int32, sum E) bool { return sum <= 10 }))     // 4
	fmt.Println(seg.MaxRight(3, func(end int32, sum E) bool { return sum <= 11 }))     // 5
	fmt.Println(seg.MaxRight(5, func(start int32, sum E) bool { return sum <= 1000 })) // 10
	fmt.Println(seg.MinLeft(5, func(start int32, sum E) bool { return sum <= 10 }))    // 4
	fmt.Println(seg.MinLeft(5, func(start int32, sum E) bool { return sum <= 11 }))    // 3
	fmt.Println(seg.MinLeft(5, func(start int32, sum E) bool { return sum <= 1001 }))  // 0
}

type SegmentTreeSqrtDecomposition[E comparable] struct {
	e           func() E
	op          func(a, b E) E
	n           int32
	bucketSize  int32
	bucketCount int32
	buckets     [][]E
	bucketSums  []E
}

// bucketSize 为 -1 时，使用默认值 sqrt(n).
func NewSegmentTreeSqrtDecomposition[E comparable](
	e func() E, op func(a, b E) E,
	n int32, f func(i int32) E, bucketSize int32,
) *SegmentTreeSqrtDecomposition[E] {
	if bucketSize == -1 {
		bucketSize = int32(math.Sqrt(float64(n))) + 1
	}
	if bucketSize < 100 {
		bucketSize = 100 // 防止 blockSize 过小
	}
	bucketCount := (n + bucketSize - 1) / bucketSize
	res := &SegmentTreeSqrtDecomposition[E]{e: e, op: op, n: n, bucketSize: bucketSize, bucketCount: bucketCount}
	buckets, bucketSum := make([][]E, bucketCount), make([]E, bucketCount)
	for bid := int32(0); bid < bucketCount; bid++ {
		start, end := bid*bucketSize, (bid+1)*bucketSize
		if end > n {
			end = n
		}
		bucket := make([]E, end-start)
		sum := res.e()
		for i := start; i < end; i++ {
			bucket[i-start] = f(i)
			sum = res.op(sum, bucket[i-start])
		}
		buckets[bid], bucketSum[bid] = bucket, sum
	}
	res.buckets, res.bucketSums = buckets, bucketSum
	return res
}

func (st *SegmentTreeSqrtDecomposition[E]) Set(index int32, value E) {
	bid := index / st.bucketSize
	pos := index - bid*st.bucketSize
	if st.buckets[bid][pos] == value {
		return
	}
	st.buckets[bid][pos] = value
	newSum := st.e()
	for _, v := range st.buckets[bid] {
		newSum = st.op(newSum, v)
	}
	st.bucketSums[bid] = newSum
}

func (st *SegmentTreeSqrtDecomposition[E]) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	bid1, bid2 := start/st.bucketSize, end/st.bucketSize
	start, end = start-bid1*st.bucketSize, end-bid2*st.bucketSize
	if bid1 == bid2 {
		res := st.e()
		bucket := st.buckets[bid1]
		for i := start; i < end; i++ {
			res = st.op(res, bucket[i])
		}
		return res
	}
	res := st.e()
	bucket1, bucket2 := st.buckets[bid1], st.buckets[bid2]
	for i := start; i < int32(len(bucket1)); i++ {
		res = st.op(res, bucket1[i])
	}
	for i := bid1 + 1; i < bid2; i++ {
		res = st.op(res, st.bucketSums[i])
	}
	for i := int32(0); i < end; i++ {
		res = st.op(res, bucket2[i])
	}
	return res
}

func (st *SegmentTreeSqrtDecomposition[E]) QueryAll() E {
	res := st.e()
	for _, v := range st.bucketSums {
		res = st.op(res, v)
	}
	return res
}

func (st *SegmentTreeSqrtDecomposition[E]) Get(index int32) E {
	bid := index / st.bucketSize
	pos := index - bid*st.bucketSize
	return st.buckets[bid][pos]
}

func (st *SegmentTreeSqrtDecomposition[E]) GetAll() []E {
	res := make([]E, 0, st.n)
	for _, bucket := range st.buckets {
		for _, v := range bucket {
			res = append(res, v)
		}
	}
	return res
}

// 查询最大的 end 使得切片 [start:end] 内的值满足 predicate.
func (st *SegmentTreeSqrtDecomposition[E]) MaxRight(start int32, predicate func(end int32, sum E) bool) int32 {
	if start >= st.n {
		return st.n
	}

	curSum := st.e()
	res := start
	bid := start / st.bucketSize

	// 散块内
	{
		pos := start - bid*st.bucketSize
		block := st.buckets[bid]
		m := int32(len(block))
		for ; pos < m; pos++ {
			nextRes, nextSum := res+1, st.op(curSum, block[pos])
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				return res
			}
		}
	}
	bid++

	// 整块跳跃
	{
		m := st.bucketCount
		for ; bid < m; bid++ {
			nextRes := res + int32(len(st.buckets[bid]))
			nextSum := st.op(curSum, st.bucketSums[bid])
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				// 答案在这个块内
				block := st.buckets[bid]
				for _, v := range block {
					nextRes, nextSum = res+1, st.op(curSum, v)
					if predicate(nextRes, nextSum) {
						res, curSum = nextRes, nextSum
					} else {
						return res
					}
				}
			}
		}
	}

	return res
}

// 查询最小的 start 使得切片 [start:end] 内的值满足 predicate.
func (st *SegmentTreeSqrtDecomposition[E]) MinLeft(end int32, predicate func(start int32, sum E) bool) int32 {
	if end <= 0 {
		return 0
	}

	curSum := st.e()
	res := end
	bid := (end - 1) / st.bucketSize

	// 散块内
	{
		pos := (end - 1) - bid*st.bucketSize
		block := st.buckets[bid]
		for ; pos >= 0; pos-- {
			nextRes, nextSum := res-1, st.op(block[pos], curSum)
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				return res
			}
		}
	}
	bid--

	// 整块跳跃
	{
		for ; bid >= 0; bid-- {
			nextRes := res - int32(len(st.buckets[bid]))
			nextSum := st.op(st.bucketSums[bid], curSum)
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				// 答案在这个块内
				block := st.buckets[bid]
				for i := int32(len(block)) - 1; i >= 0; i-- {
					nextRes, nextSum = res-1, st.op(block[i], curSum)
					if predicate(nextRes, nextSum) {
						res, curSum = nextRes, nextSum
					} else {
						return res
					}
				}
			}
		}
	}
	return res
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

func test() {
	for i := int32(0); i < 100; i++ {
		n := rand.Int31n(10000) + 1000
		nums := make([]int, n)
		for i := int32(0); i < n; i++ {
			nums[i] = rand.Intn(100)
		}

		type E = int
		e, op := func() int { return 0 }, func(a, b int) int { return a + b }
		seg := NewSegmentTreeSqrtDecomposition(e, op, n, func(i int32) E { return E(nums[i]) }, -1)

		for j := 0; j < 1000; j++ {
			// Get
			index := rand.Int31n(n)
			if seg.Get(index) != E(nums[index]) {
				fmt.Println("Get Error")
				panic("Get Error")
			}

			// Set
			index = rand.Int31n(n)
			value := rand.Intn(100)
			nums[index] = value
			seg.Set(index, E(value))
			if seg.Get(index) != E(value) {
				fmt.Println("Set Error")
				panic("Set Error")
			}

			// Query
			start, end := rand.Int31n(n), rand.Int31n(n)
			if start > end {
				start, end = end, start
			}
			sum_ := E(0)
			for i := start; i < end; i++ {
				sum_ += E(nums[i])
			}
			if seg.Query(start, end) != sum_ {
				fmt.Println("Query Error")
				panic("Query Error")
			}

			// QueryAll
			sum_ = E(0)
			for _, v := range nums {
				sum_ += E(v)
			}
			if seg.QueryAll() != sum_ {
				fmt.Println("QueryAll Error")
				panic("QueryAll Error")
			}

			// GetAll
			all := seg.GetAll()
			for i, v := range all {
				if v != E(nums[i]) {
					fmt.Println("GetAll Error")
					panic("GetAll Error")
				}
			}

			// MaxRight
			maxRightBf := func(start int32, predicate func(end int32, sum E) bool) (res int32) {
				res = start
				curSum := seg.e()
				for i := start; i < n; i++ {
					curSum = seg.op(curSum, E(nums[i]))
					if !predicate(i+1, curSum) {
						return
					}
					res = i + 1
				}
				return
			}

			minLeftBf := func(end int32, predicate func(start int32, sum E) bool) (res int32) {
				res = end
				curSum := seg.e()
				for i := end - 1; i >= 0; i-- {
					curSum = seg.op(curSum, E(nums[i]))
					if !predicate(i, curSum) {
						return
					}
					res = i
				}
				return
			}

			{
				start := rand.Int31n(n)
				upper := rand.Intn(100)
				res1 := seg.MaxRight(start, func(end int32, sum E) bool { return sum <= E(upper) })
				res2 := maxRightBf(start, func(end int32, sum E) bool { return sum <= E(upper) })
				if res1 != res2 {
					fmt.Println("MaxRight Error")
					panic("MaxRight Error")
				}
				res3 := seg.MinLeft(start, func(start int32, sum E) bool { return sum <= E(upper) })
				res4 := minLeftBf(start, func(start int32, sum E) bool { return sum <= E(upper) })
				if res3 != res4 {
					fmt.Println("MinLeft Error")
					panic("MinLeft Error")
				}
			}
		}

	}
	fmt.Println("Pass")
}

func testTime() {
	// 1e5
	n := int32(2e5)
	nums := make([]int, n)
	for i := 0; i < int(n); i++ {
		nums[i] = rand.Intn(100)
	}

	time1 := time.Now()

	type E = int
	e, op := func() int { return 0 }, func(a, b int) int { return a + b }
	seg := NewSegmentTreeSqrtDecomposition(e, op, n, func(i int32) int { return nums[i] }, -1)

	for i := int32(0); i < n; i++ {
		seg.Query(i, n)
		seg.QueryAll()
		seg.Get(i)
		seg.Set(i, int(E(i)))
		seg.MaxRight(i, func(end int32, sum E) bool { return sum <= nums[i] })
		seg.MinLeft(i, func(start int32, sum E) bool { return sum <= nums[i] })
	}
	fmt.Println("Time1", time.Since(time1)) // Time1 49.79925ms
}
