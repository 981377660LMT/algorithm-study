// api:
//  1.Set(index int32, value E) -> O(sqrt(n))
//  2.Query(start, end int32) E -> O(sqrt(n))
//  3.QueryAll() E -> O(sqrt(n))
//  !4.Get(index int32) E -> O(1)
//  5.GetAll() []E -> O(n)

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	test()
	testTime()
}

func demo() {
	seg := NewSegmentTreeSqrtDecompositionDynamic(10, func(i int32) int { return int(i) }, -1)
	fmt.Println(seg.GetAll())
	seg.Set(3, 5)
	seg.Set(4, 6)
	fmt.Println(seg.GetAll())
	fmt.Println(seg.Query(3, 5)) // 11
}

type E = int

func (*SegmentTreeSqrtDecompositionDynamic) e() E        { return 0 }
func (*SegmentTreeSqrtDecompositionDynamic) op(a, b E) E { return a + b }

type SegmentTreeSqrtDecompositionDynamic struct {
	n           int32
	bucketSize  int32
	bucketCount int32
	buckets     [][]E
	bucketSums  []E
}

// bucketSize 为 -1 时，使用默认值 sqrt(n).
func NewSegmentTreeSqrtDecompositionDynamic(n int32, f func(i int32) E, bucketSize int32) *SegmentTreeSqrtDecompositionDynamic {
	if bucketSize == -1 {
		bucketSize = int32(math.Sqrt(float64(n))) + 1
	}
	bucketCount := (n + bucketSize - 1) / bucketSize
	res := &SegmentTreeSqrtDecompositionDynamic{n: n, bucketSize: bucketSize, bucketCount: bucketCount}
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

func (st *SegmentTreeSqrtDecompositionDynamic) Set(index int32, value E) {
	bid := index / st.bucketSize
	pos := index - bid*st.bucketSize
	st.buckets[bid][pos] = value
	newSum := st.e()
	for _, v := range st.buckets[bid] {
		newSum = st.op(newSum, v)
	}
	st.bucketSums[bid] = newSum
}

func (st *SegmentTreeSqrtDecompositionDynamic) Query(start, end int32) E {
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

func (st *SegmentTreeSqrtDecompositionDynamic) QueryAll() E {
	res := st.e()
	for _, v := range st.bucketSums {
		res = st.op(res, v)
	}
	return res
}

func (st *SegmentTreeSqrtDecompositionDynamic) Get(index int32) E {
	bid := index / st.bucketSize
	pos := index - bid*st.bucketSize
	return st.buckets[bid][pos]
}

func (st *SegmentTreeSqrtDecompositionDynamic) GetAll() []E {
	res := make([]E, 0, st.n)
	for _, bucket := range st.buckets {
		for _, v := range bucket {
			res = append(res, v)
		}
	}
	return res
}

func test() {
	for i := int32(0); i < 100; i++ {
		n := rand.Int31n(10000) + 1000
		nums := make([]int, n)
		for i := int32(0); i < n; i++ {
			nums[i] = rand.Intn(100)
		}
		seg := NewSegmentTreeSqrtDecompositionDynamic(n, func(i int32) E { return E(nums[i]) }, -1)

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

		}

	}
	fmt.Println("Pass")
}

func testTime() {
	// 2e5
	n := int32(2e5)
	nums := make([]int, n)
	for i := 0; i < int(n); i++ {
		nums[i] = rand.Intn(5)
	}

	time1 := time.Now()
	seg := NewSegmentTreeSqrtDecompositionDynamic(n, func(i int32) int { return nums[i] }, -1)

	for i := int32(0); i < n; i++ {
		seg.Get(i)
		seg.Set(i, int(E(i)))
		seg.Query(i, n)
		seg.QueryAll()
	}
	fmt.Println("Time1", time.Since(time1)) // Time1 128.573042ms
}
