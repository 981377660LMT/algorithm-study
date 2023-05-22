// 动态区间频率查询
// 区间加
// 查询区间某元素出现次数

package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

func main() {
	nums := []int{1, 2, 2, 4, 5, 6, 7, 8, 9, 10}
	rf := NewRangeFreqQueryDynamic(nums)
	rf.Add(0, 10, 1)
	fmt.Println(rf.RangeFreq(0, 10, 5))
	rf.Add(0, 10, 2)
	fmt.Println(rf.RangeFreq(0, 10, 5))
	fmt.Println(rf.RangeFreqWithFloor(0, 10, 5))

	nums = make([]int, 1e5)
	rf = NewRangeFreqQueryDynamic(nums)
	time1 := time.Now()
	for i := 0; i < 1e5; i++ {
		rf.Add(i, i+1, i)
		rf.RangeFreq(0, i+1, i)
	}
	fmt.Println(time.Since(time1))
}

type RangeFreqQueryDynamic struct {
	sqrt *SqrtDecomposition
}

func NewRangeFreqQueryDynamic(nums []int) *RangeFreqQueryDynamic {
	n := len(nums)
	return &RangeFreqQueryDynamic{sqrt: NewSqrtDecomposition(nums, 1+int(math.Sqrt(float64(n))))}
}

// 左闭右开区间[start,end)的每个数加上add.
func (rf *RangeFreqQueryDynamic) Add(start, end, add int) {
	rf.sqrt.Update(start, end, add)
}

// 左闭右开区间[start,end)内元素k的出现次数.
func (rf *RangeFreqQueryDynamic) RangeFreq(start, end int, k int) int {
	res := 0
	rf.sqrt.Query(start, end, func(blockRes int) { res += blockRes }, k, true)
	return res
}

// 左闭右开区间[start,end)内大于等于floor的元素个数.
func (rf *RangeFreqQueryDynamic) RangeFreqWithFloor(start, end int, floor int) int {
	res := 0
	rf.sqrt.Query(start, end, func(blockRes int) { res += blockRes }, floor, false)
	return res
}

type E = int
type Id = int

type Block struct {
	id, start, end int
	nums           []E

	sorted  []E
	lazyAdd Id
}

func (b *Block) Created() { b.Updated() }
func (b *Block) Updated() {
	b.sorted = append(b.sorted[:0:0], b.nums...)
	sort.Ints(b.sorted)
}

// !区间加.
func (b *Block) UpdateAll(lazy Id) { b.lazyAdd += lazy }
func (b *Block) UpdatePart(start, end int, lazy Id) {
	for i := start; i < end; i++ {
		b.nums[i] += lazy
	}
}

// !查询区间等于x的元素个数/大于等于x的元素个数.
func (b *Block) QueryAll(x int, same bool) E {
	if same {
		pos2 := sort.SearchInts(b.sorted, x-b.lazyAdd+1)
		pos1 := sort.SearchInts(b.sorted, x-b.lazyAdd)
		return pos2 - pos1
	}

	lower := sort.SearchInts(b.sorted, x-b.lazyAdd)
	return len(b.sorted) - lower
}

func (b *Block) QueryPart(start, end int, x int, same bool) E {
	if same {
		res := 0
		for i := start; i < end; i++ {
			if b.nums[i]+b.lazyAdd == x {
				res++
			}
		}
		return res
	}

	res := 0
	for i := start; i < end; i++ {
		if b.nums[i]+b.lazyAdd >= x {
			res++
		}
	}
	return res
}

type SqrtDecomposition struct {
	n   int
	bs  int
	bls []Block
}

// 指定维护的序列和分块大小初始化.
//  blockSize:分块大小,一般取根号n(300)
func NewSqrtDecomposition(nums []E, blockSize int) *SqrtDecomposition {
	nums = append(nums[:0:0], nums...)
	res := &SqrtDecomposition{n: len(nums), bs: blockSize, bls: make([]Block, len(nums)/blockSize+1)}
	for i := range res.bls {
		res.bls[i].id = i
		res.bls[i].start = i * blockSize
		res.bls[i].end = min((i+1)*blockSize, len(nums))
		res.bls[i].nums = nums[res.bls[i].start:res.bls[i].end]
		res.bls[i].Created()
	}
	return res
}

// 更新左闭右开区间[start,end)的值.
//  0<=start<=end<=n
func (s *SqrtDecomposition) Update(start, end int, lazy Id) {
	if start >= end {
		return
	}
	id1, id2 := start/s.bs, end/s.bs
	pos1, pos2 := start%s.bs, end%s.bs
	if id1 == id2 {
		s.bls[id1].UpdatePart(pos1, pos2, lazy)
		s.bls[id1].Updated()
	} else {
		s.bls[id1].UpdatePart(pos1, s.bs, lazy)
		s.bls[id1].Updated()
		for i := id1 + 1; i < id2; i++ {
			s.bls[i].UpdateAll(lazy)
		}
		s.bls[id2].UpdatePart(0, pos2, lazy)
		s.bls[id2].Updated()
	}
}

// 查询左闭右开区间[start,end)的值.
//  0<=start<=end<=n
func (s *SqrtDecomposition) Query(start, end int, cb func(blockRes E), k int, same bool) {
	if start >= end {
		return
	}
	id1, id2 := start/s.bs, end/s.bs
	pos1, pos2 := start%s.bs, end%s.bs
	if id1 == id2 {
		cb(s.bls[id1].QueryPart(pos1, pos2, k, same))
		return
	}
	cb(s.bls[id1].QueryPart(pos1, s.bs, k, same))
	for i := id1 + 1; i < id2; i++ {
		cb(s.bls[i].QueryAll(k, same))
	}
	cb(s.bls[id2].QueryPart(0, pos2, k, same))
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
