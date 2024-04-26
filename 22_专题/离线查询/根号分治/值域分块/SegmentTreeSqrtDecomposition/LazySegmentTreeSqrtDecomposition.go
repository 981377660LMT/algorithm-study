// api:
//  1.Set(index int32, value E) -> O(sqrt(n))
//  !2.Get(index int32) E -> O(1)
//  3.Query(start, end int32) E
//  4.QueryAll() E
//  6.Update(start, end int32, lazy Id)
//  5.GetAll() []E

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
	seg := NewLazySegmentTreeSqrtDecomposition(10, func(i int32) int { return int(i) }, -1)
	fmt.Println(seg.GetAll())
	seg.Update(1, 5, 1)

	fmt.Println(seg.GetAll())
	seg.Update(3, 6, 1)

	fmt.Println(seg.GetAll())
	fmt.Println(seg.Query(3, 5)) // 6

	fmt.Println(seg.GetAll())
}

type E = int
type Id = int32

// RangeAddRangeMax

func (*LazySegmentTreeSqrtDecomposition) e() E        { return 0 }
func (*LazySegmentTreeSqrtDecomposition) id() Id      { return 0 }
func (*LazySegmentTreeSqrtDecomposition) op(a, b E) E { return max(a, b) }
func (*LazySegmentTreeSqrtDecomposition) mapping(f Id, a E) E {
	return a + E(f)
}
func (*LazySegmentTreeSqrtDecomposition) composition(f, g Id) Id {
	return f + g
}

type LazySegmentTreeSqrtDecomposition struct {
	n           int32
	bucketSize  int32
	bucketCount int32
	data        [][]E
	sum         []E
	lazy        []Id
}

// bucketSize 为 -1 时，使用默认值 sqrt(n).
func NewLazySegmentTreeSqrtDecomposition(n int32, f func(i int32) E, bucketSize int32) *LazySegmentTreeSqrtDecomposition {
	if bucketSize == -1 {
		bucketSize = int32(math.Sqrt(float64(n))) + 1
	}
	bucketCount := (n + bucketSize - 1) / bucketSize
	res := &LazySegmentTreeSqrtDecomposition{n: n, bucketSize: bucketSize, bucketCount: bucketCount}
	buckets, bucketSum := make([][]E, bucketCount), make([]E, bucketCount)
	bucketLazys := make([]Id, bucketCount)
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
		bucketLazys[bid] = res.id()
	}
	res.data, res.sum, res.lazy = buckets, bucketSum, bucketLazys
	return res
}

func (st *LazySegmentTreeSqrtDecomposition) Set(index int32, value E) {
	bid := index / st.bucketSize
	st._propagate(bid)
	st.data[bid][index-bid*st.bucketSize] = value
	newSum := st.e()
	for _, v := range st.data[bid] {
		newSum = st.op(newSum, v)
	}
	st.sum[bid] = newSum
}

func (st *LazySegmentTreeSqrtDecomposition) Get(index int32) E {
	bid := index / st.bucketSize
	pos := index - bid*st.bucketSize
	if st.lazy[bid] == st.id() {
		return st.data[bid][pos]
	} else {
		return st.mapping(st.lazy[bid], st.data[bid][pos])
	}
}

func (st *LazySegmentTreeSqrtDecomposition) Query(start, end int32) E {
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
	res := st.e()
	if bid1 == bid2 {
		block := st.data[bid1]
		for i := start; i < end; i++ {
			res = st.op(res, block[i])
		}
		if st.lazy[bid1] != st.id() {
			res = st.mapping(st.lazy[bid1], res)
		}
	} else {
		if start < int32(len(st.data[bid1])) {
			block1 := st.data[bid1]
			for i := start; i < int32(len(block1)); i++ {
				res = st.op(res, block1[i])
			}
			if st.lazy[bid1] != st.id() {
				res = st.mapping(st.lazy[bid1], res)
			}
		}

		for i := bid1 + 1; i < bid2; i++ {
			res = st.op(res, st.sum[i])
		}

		if bid2 < st.bucketCount && end > 0 {
			block2 := st.data[bid2]
			tmp := st.e()
			for i := int32(0); i < end; i++ {
				tmp = st.op(tmp, block2[i])
			}
			if st.lazy[bid2] != st.id() {
				tmp = st.mapping(st.lazy[bid2], tmp)
			}
			res = st.op(res, tmp)
		}
	}
	return res
}

func (st *LazySegmentTreeSqrtDecomposition) QueryAll() E {
	res := st.e()
	for bid := int32(0); bid < st.bucketCount; bid++ {
		res = st.op(res, st.sum[bid])
	}
	return res
}

func (st *LazySegmentTreeSqrtDecomposition) Update(start, end int32, lazy Id) {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return
	}

	changeData := func(k, l, r int32) {
		st._propagate(k)
		d := st.data[k]
		for i := l; i < r; i++ {
			d[i] = st.mapping(lazy, d[i])
		}
		e := st.e()
		for _, v := range d {
			e = st.op(e, v)
		}
		st.sum[k] = e
	}

	bid1, bid2 := start/st.bucketSize, end/st.bucketSize
	start, end = start-bid1*st.bucketSize, end-bid2*st.bucketSize
	if bid1 == bid2 {
		if bid1 < st.bucketCount {
			changeData(bid1, start, end)
		}
	} else {
		if bid1 < st.bucketCount {
			if start == 0 {
				if st.lazy[bid1] == st.id() {
					st.lazy[bid1] = lazy
				} else {
					st.lazy[bid1] = st.composition(lazy, st.lazy[bid1])
				}
				st.sum[bid1] = st.mapping(lazy, st.sum[bid1])
			} else {
				changeData(bid1, start, int32(len(st.data[bid1])))
			}
		}

		for i := bid1 + 1; i < bid2; i++ {
			if st.lazy[i] == st.id() {
				st.lazy[i] = lazy
			} else {
				st.lazy[i] = st.composition(lazy, st.lazy[i])
			}
			st.sum[i] = st.mapping(lazy, st.sum[i])
		}

		if bid2 < st.bucketCount {
			if end == int32(len(st.data[bid2])) {
				if st.lazy[bid2] == st.id() {
					st.lazy[bid2] = lazy
				} else {
					st.lazy[bid2] = st.composition(lazy, st.lazy[bid2])
				}
				st.sum[bid2] = st.mapping(lazy, st.sum[bid2])
			} else {
				changeData(bid2, 0, end)
			}
		}
	}
}

func (st *LazySegmentTreeSqrtDecomposition) GetAll() []E {
	st._propagateAll()
	res := make([]E, 0, st.n)
	for _, bucket := range st.data {
		res = append(res, bucket...)
	}
	return res
}

func (st *LazySegmentTreeSqrtDecomposition) _propagate(k int32) {
	if st.lazy[k] == st.id() {
		return
	}
	f := st.lazy[k]
	dk := st.data[k]
	for i := int32(0); i < int32(len(dk)); i++ {
		dk[i] = st.mapping(f, dk[i])
	}
	st.lazy[k] = st.id()
}

func (st *LazySegmentTreeSqrtDecomposition) _propagateAll() {
	for k := int32(0); k < st.bucketCount; k++ {
		st._propagate(k)
	}
	for i := int32(0); i < st.bucketCount; i++ {
		st.lazy[i] = st.id()
	}
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func test() {
	for i := 0; i < 10; i++ {
		n := rand.Int31n(50) + 1
		nums := make([]int, n)
		for i := 0; i < int(n); i++ {
			nums[i] = rand.Intn(5)
		}
		seg := NewLazySegmentTreeSqrtDecomposition(n, func(i int32) int { return nums[i] }, -1)

		for j := 0; j < 3; j++ {

			index := rand.Int31n(n)
			// Get
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
			res := E(0)
			for i := start; i < end; i++ {
				if nums[i] > res {
					res = E(nums[i])
				}
			}
			if seg.Query(start, end) != res {
				fmt.Println("Query Error", seg.Query(start, end), res)
				panic("Query Error")
			}

			// QueryAll
			res = E(0)
			for _, v := range nums {
				if v > res {
					res = E(v)
				}
			}
			fmt.Println("QueryAll", seg.QueryAll(), res, seg.GetAll(), nums)
			if seg.QueryAll() != res {
				fmt.Println("QueryAll Error", seg.QueryAll(), res)
				fmt.Println(seg.GetAll())
				panic("QueryAll Error")
			}

			// Update
			start, end = rand.Int31n(n), rand.Int31n(n)
			if start > end {
				start, end = end, start
			}
			lazy := rand.Intn(100)
			for i := start; i < end; i++ {
				nums[i] += lazy
			}
			seg.Update(start, end, Id(lazy))
			for i := 0; i < int(n); i++ {
				if seg.Get(int32(i)) != E(nums[i]) {
					fmt.Println("Update Error")
					panic("Update Error")
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
	seg := NewLazySegmentTreeSqrtDecomposition(n, func(i int32) int { return nums[i] }, -1)

	for i := int32(0); i < n; i++ {
		seg.Get(i)
		seg.Set(i, int(E(i)))
		seg.Query(i, n)
		seg.QueryAll()
		seg.Update(i, n, Id(i))
	}
	fmt.Println("Time1", time.Since(time1)) // Time1 276.134709ms
}
