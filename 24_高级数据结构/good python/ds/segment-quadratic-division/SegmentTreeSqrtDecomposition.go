package main

import (
	"fmt"
	"math"
)

func main() {
	seg := NewSegmentTreeSqrtDecomposition(10, func(i int32) int { return int(i) }, -1)
	fmt.Println(seg.GetAll())
	seg.Set(3, 5)
	seg.Set(4, 6)
	fmt.Println(seg.GetAll())
	fmt.Println(seg.Query(3, 5)) // 11
}

type E = int

func (*SegmentTreeSqrtDecomposition) e() E        { return 0 }
func (*SegmentTreeSqrtDecomposition) op(a, b E) E { return a + b }

type SegmentTreeSqrtDecomposition struct {
	n           int32
	bucketSize  int32
	bucketCount int32
	buckets     [][]E
	bucketSums  []E
}

// bucketSize 为 -1 时，使用默认值 sqrt(n).
func NewSegmentTreeSqrtDecomposition(n int32, f func(i int32) E, bucketSize int32) *SegmentTreeSqrtDecomposition {
	if bucketSize == -1 {
		bucketSize = int32(math.Sqrt(float64(n))) + 1
	}
	bucketCount := (n + bucketSize - 1) / bucketSize
	res := &SegmentTreeSqrtDecomposition{n: n, bucketSize: bucketSize, bucketCount: bucketCount}
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

func (st *SegmentTreeSqrtDecomposition) Set(index int32, value E) {
	bid := index / st.bucketSize
	pos := index - bid*st.bucketSize
	st.buckets[bid][pos] = value
	newSum := st.e()
	for _, v := range st.buckets[bid] {
		newSum = st.op(newSum, v)
	}
	st.bucketSums[bid] = newSum
}

func (st *SegmentTreeSqrtDecomposition) Query(start, end int32) E {
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

func (st *SegmentTreeSqrtDecomposition) QueryAll() E {
	res := st.e()
	for _, v := range st.bucketSums {
		res = st.op(res, v)
	}
	return res
}

func (st *SegmentTreeSqrtDecomposition) Get(index int32) E {
	bid := index / st.bucketSize
	pos := index - bid*st.bucketSize
	return st.buckets[bid][pos]
}

func (st *SegmentTreeSqrtDecomposition) GetAll() []E {
	res := make([]E, 0, st.n)
	for _, bucket := range st.buckets {
		for _, v := range bucket {
			res = append(res, v)
		}
	}
	return res
}
