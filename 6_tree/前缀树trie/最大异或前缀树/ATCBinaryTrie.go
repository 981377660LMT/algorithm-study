package main

import (
	"bytes"
	"fmt"
	"math/bits"
)

func findMaximumXOR(nums []int) int {
	n := len(nums)
	maxLog := bits.Len(uint(max(nums...))) // 1+int(log2(x))
	bt := NewATCBinaryTrie(n, maxLog, true)
	for _, v := range nums {
		bt.Add(v)
	}
	res := 0
	for _, v := range nums {
		bt.XorAll(v)
		res = max(res, bt.Max())
		bt.XorAll(v)
	}
	return res
}

// Reference:
//  - https://atcoder.jp/contests/arc028/submissions/19916627
//  - https://judge.yosupo.jp/submission/35057
type ATCBinaryTrie struct {
	multiset                                bool
	maxLog, xEnd, addQueryLimit, maxV, lazy int
	vList, edges, size, isEnd               []int
}

//  addQueryLimit: max number of add and query operations
//  maxLog: max log of x
//  allowMultipleElements: allow multiple elements with the same value (SortedList or SortedSet)
//
// Example of `Count Pairs With XOR in a Range`:
//  n := len(nums)
//  maxLog := bits.Len(uint(max(nums...)))
//  bt := NewATCBinaryTrie(n, maxLog, true)
//  for _, v := range nums {
//  	bt.Add(v)
//  }
//  res := 0
//  for _, v := range nums {
//  	bt.XorAll(v)
//  	res += bt.BisectRight(high) - bt.BisectLeft(low)
//  	bt.XorAll(v)
//  }
//	return res / 2
func NewATCBinaryTrie(addQueryLimit, maxLog int, allowMultipleElements bool) *ATCBinaryTrie {
	n := maxLog*addQueryLimit + 1
	edges := make([]int, 2*n)
	for i := range edges {
		edges[i] = -1
	}

	return &ATCBinaryTrie{
		maxLog:        maxLog,
		xEnd:          1 << maxLog,
		multiset:      allowMultipleElements,
		addQueryLimit: addQueryLimit,
		edges:         edges,
		size:          make([]int, n),
		isEnd:         make([]int, n),
		vList:         make([]int, maxLog+1),
	}
}

func (bt *ATCBinaryTrie) Add(x int) {
	x ^= bt.lazy
	v := 0
	for i := bt.maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		if bt.edges[2*v+d] == -1 {
			bt.maxV++
			bt.edges[2*v+d] = bt.maxV
		}
		v = bt.edges[2*v+d]
		bt.vList[i] = v
	}

	if bt.multiset || bt.isEnd[v] == 0 {
		bt.isEnd[v]++
		for _, v := range bt.vList {
			bt.size[v]++
		}
	}
}

func (bt *ATCBinaryTrie) Discard(x int) {
	if x < 0 || x >= bt.xEnd {
		return
	}
	x ^= bt.lazy
	v := 0
	for i := bt.maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		if bt.edges[2*v+d] == -1 {
			return
		}
		v = bt.edges[2*v+d]
		bt.vList[i] = v
	}
	if bt.isEnd[v] > 0 {
		bt.isEnd[v]--
		for _, v := range bt.vList {
			bt.size[v]--
		}
	}
}

// 删除count个x count=-1表示删除所有x
func (bt *ATCBinaryTrie) Erase(x int, count int) {
	if x < 0 || x >= bt.xEnd {
		return
	}
	x ^= bt.lazy
	v := 0
	for i := bt.maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		if bt.edges[2*v+d] == -1 {
			return
		}
		v = bt.edges[2*v+d]
		bt.vList[i] = v
	}
	if count == -1 || bt.isEnd[v] < count {
		count = bt.isEnd[v]
	}
	if bt.isEnd[v] > 0 {
		bt.isEnd[v] -= count
		for _, v := range bt.vList {
			bt.size[v] -= count
		}
	}
}

func (bt *ATCBinaryTrie) Count(x int) int {
	if x < 0 || x >= bt.xEnd {
		return 0
	}
	x ^= bt.lazy
	v := 0
	for i := bt.maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		if bt.edges[2*v+d] == -1 {
			return 0
		}
		v = bt.edges[2*v+d]
	}
	return bt.isEnd[v]
}

func (bt *ATCBinaryTrie) BisectLeft(x int) int {
	if x < 0 {
		return 0
	}
	if bt.xEnd <= x {
		return bt.Size()
	}
	v := 0
	ret := 0
	for i := bt.maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		l := (bt.lazy >> i) & 1
		lc := bt.edges[2*v]
		rc := bt.edges[2*v+1]
		if l == 1 {
			lc, rc = rc, lc
		}
		if d != 0 {
			if lc != -1 {
				ret += bt.size[lc]
			}
			if rc == -1 {
				return ret
			}
			v = rc
		} else {
			if lc == -1 {
				return ret
			}
			v = lc
		}
	}
	return ret
}

func (bt *ATCBinaryTrie) BisectRight(x int) int {
	return bt.BisectLeft(x + 1)
}

func (bt *ATCBinaryTrie) Index(x int) int {
	if bt.Count(x) == 0 {
		panic(fmt.Sprintf("%d is not in BinaryTrie", x))
	}
	return bt.BisectLeft(x)
}

func (bt *ATCBinaryTrie) Find(x int) int {
	if bt.Count(x) == 0 {
		return -1
	}
	return bt.BisectLeft(x)
}

// 0<=k<bt.Size()
//  support negative index
func (bt *ATCBinaryTrie) At(k int) int {
	if k < 0 {
		k += bt.size[0]
	}
	v := 0
	res := 0
	for i := bt.maxLog - 1; i > -1; i-- {
		l := (bt.lazy >> i) & 1
		lc := bt.edges[2*v]
		rc := bt.edges[2*v+1]
		if l == 1 {
			lc, rc = rc, lc
		}
		if lc == -1 {
			v = rc
			res |= 1 << i
			continue
		}
		if bt.size[lc] <= k {
			k -= bt.size[lc]
			v = rc
			res |= 1 << i
		} else {
			v = lc
		}
	}
	return res
}

func (bt *ATCBinaryTrie) Min() int {
	return bt.At(0)
}

func (bt *ATCBinaryTrie) Max() int {
	return bt.At(-1)
}

func (bt *ATCBinaryTrie) XorAll(x int) {
	bt.lazy ^= x
}

func (bt *ATCBinaryTrie) Has(x int) bool {
	return bt.Count(x) > 0
}

func (bt *ATCBinaryTrie) Size() int {
	return bt.size[0]
}

func (bt *ATCBinaryTrie) ForEach(callbackfn func(value, index int)) {
	q := make([][2]int, 0, 16)
	q = append(q, [2]int{0, 0})
	for i := bt.maxLog - 1; i > -1; i-- {
		l := (bt.lazy >> i) & 1
		nq := make([][2]int, 0, 16)
		for _, v := range q {
			lc := bt.edges[2*v[0]]
			rc := bt.edges[2*v[0]+1]
			if l == 1 {
				lc, rc = rc, lc
			}
			if lc != -1 {
				nq = append(nq, [2]int{lc, 2 * v[1]})
			}
			if rc != -1 {
				nq = append(nq, [2]int{rc, 2*v[1] + 1})
			}
		}
		q = nq
	}

	i := 0
	for _, item := range q {
		v, x := item[0], item[1]
		for j := 0; j < bt.isEnd[v]; j++ {
			callbackfn(x, i)
			i++
		}
	}
}

func (bt *ATCBinaryTrie) String() string {
	var buf bytes.Buffer
	buf.WriteString("BinaryTrie{")
	bt.ForEach(func(x, i int) {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%d", x))
	})
	buf.WriteString("}")
	return buf.String()
}

func max(nums ...int) int {
	if len(nums) == 0 {
		panic("max: no arguments")
	}
	res := nums[0]
	for _, num := range nums[1:] {
		if num > res {
			res = num
		}
	}
	return res
}
