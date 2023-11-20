// Reference:
//  - https://atcoder.jp/contests/arc028/submissions/19916627
//  - https://judge.yosupo.jp/submission/35057

// Usage:
// trie := NewBinaryTrie(n, maxLog, allowMultipleElements)
// trie.Add(x)
// trie.Discard(x)
// trie.XorAll(x)
// trie.Kth(k)
// trie.At(i)
// trie.Min()
// trie.Max()
// trie.Count(x)
// trie.Erase(-1)
// trie.Index(x)
// trie.Find(x)
// trie.Size()
// trie.Has()
// trie.BisectLeft(x)
// trie.BisectRight(x)
// trie.ForEach(func(x int) bool { return true })

package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"sort"
)

func main() {
	xorTrie := NewBinaryTrie(1e9, 1e5, true)
	xorTrie.Add(2)
	xorTrie.Add(1)
	xorTrie.Add(3)
	fmt.Println(xorTrie)
}

// https://leetcode.cn/problems/maximum-xor-of-two-numbers-in-an-array/
func findMaximumXOR(nums []int) int {
	max_ := 1
	for _, num := range nums {
		max_ = max(max_, num)
	}
	tree := NewBinaryTrie(max_, len(nums), true)

	maxXor := 0
	for _, num := range nums {
		tree.Add(num)
		tree.XorAll(num)
		maxXor = max(maxXor, tree.Max())
		tree.XorAll(num)
	}
	return maxXor

}

// 1803. 统计异或值在范围内的数对有多少
// https://leetcode.cn/problems/count-pairs-with-xor-in-a-range/description/
func countPairs(nums []int, low int, high int) int {
	n := len(nums)
	bt := NewBinaryTrie(max(nums...), n, true)
	for _, v := range nums {
		bt.Add(v)
	}
	res := 0
	for _, v := range nums {
		bt.XorAll(v)
		res += bt.BisectRight(high) - bt.BisectLeft(low)
		bt.XorAll(v)
	}
	return res / 2
}

// 2935. 找出强数对的最大异或值 II
// https://leetcode.cn/problems/maximum-strong-pair-xor-ii/description/
func maximumStrongPairXor(nums []int) int {
	sort.Ints(nums)
	res, left, n := 0, 0, len(nums)
	trie := NewBinaryTrie(max(nums...), n, true)
	for right, cur := range nums {
		trie.Add(cur)
		for left <= right && cur > 2*nums[left] {
			trie.Discard(nums[left])
			left++
		}
		trie.XorAll(cur)
		res = max(res, trie.Max())
		trie.XorAll(cur)
	}
	return res
}

// XorTrie.
type BinaryTrie struct {
	_multiset                        bool
	_maxLog, _xEnd, _maxV, _lazy     int
	_vList, _edges, _size, _endCount []int
}

// max: max of x
// addLimit: max number of add and query operations
// allowMultipleElements: allow multiple elements with the same value (SortedList or SortedSet)
func NewBinaryTrie(max, addLimit int, allowMultipleElements bool) *BinaryTrie {
	maxLog := bits.Len(uint(max))
	n := maxLog*addLimit + 1
	edges := make([]int, 2*n)
	for i := range edges {
		edges[i] = -1
	}

	return &BinaryTrie{
		_multiset: allowMultipleElements,
		_maxLog:   maxLog,
		_xEnd:     1 << maxLog,
		_vList:    make([]int, maxLog+1),
		_edges:    edges,
		_size:     make([]int, n),
		_endCount: make([]int, n),
	}
}

func (bt *BinaryTrie) Add(x int) {
	if x < 0 || x >= bt._xEnd {
		return
	}
	x ^= bt._lazy
	v := 0
	for i := bt._maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		if bt._edges[2*v+d] == -1 {
			bt._maxV++
			bt._edges[2*v+d] = bt._maxV
		}
		v = bt._edges[2*v+d]
		bt._vList[i] = v
	}

	if bt._multiset || bt._endCount[v] == 0 {
		bt._endCount[v]++
		for _, v := range bt._vList {
			bt._size[v]++
		}
	}
}

func (bt *BinaryTrie) Discard(x int) {
	if x < 0 || x >= bt._xEnd {
		return
	}
	x ^= bt._lazy
	v := 0
	for i := bt._maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		if bt._edges[2*v+d] == -1 {
			return
		}
		v = bt._edges[2*v+d]
		bt._vList[i] = v
	}
	if bt._endCount[v] > 0 {
		bt._endCount[v]--
		for _, v := range bt._vList {
			bt._size[v]--
		}
	}
}

// 删除count个x count=-1表示删除所有x
func (bt *BinaryTrie) Erase(x int, count int) {
	if x < 0 || x >= bt._xEnd {
		return
	}
	x ^= bt._lazy
	v := 0
	for i := bt._maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		if bt._edges[2*v+d] == -1 {
			return
		}
		v = bt._edges[2*v+d]
		bt._vList[i] = v
	}
	if count == -1 || bt._endCount[v] < count {
		count = bt._endCount[v]
	}
	if bt._endCount[v] > 0 {
		bt._endCount[v] -= count
		for _, v := range bt._vList {
			bt._size[v] -= count
		}
	}
}

func (bt *BinaryTrie) Count(x int) int {
	if x < 0 || x >= bt._xEnd {
		return 0
	}
	x ^= bt._lazy
	v := 0
	for i := bt._maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		if bt._edges[2*v+d] == -1 {
			return 0
		}
		v = bt._edges[2*v+d]
	}
	return bt._endCount[v]
}

func (bt *BinaryTrie) BisectLeft(x int) int {
	if x < 0 {
		return 0
	}
	if bt._xEnd <= x {
		return bt.Size()
	}
	v := 0
	res := 0
	for i := bt._maxLog - 1; i > -1; i-- {
		d := (x >> i) & 1
		l := (bt._lazy >> i) & 1
		lc := bt._edges[2*v]
		rc := bt._edges[2*v+1]
		if l == 1 {
			lc, rc = rc, lc
		}
		if d != 0 {
			if lc != -1 {
				res += bt._size[lc]
			}
			if rc == -1 {
				return res
			}
			v = rc
		} else {
			if lc == -1 {
				return res
			}
			v = lc
		}
	}
	return res
}

func (bt *BinaryTrie) BisectRight(x int) int {
	return bt.BisectLeft(x + 1)
}

func (bt *BinaryTrie) Index(x int) int {
	if bt.Count(x) == 0 {
		panic(fmt.Sprintf("%d is not in BinaryTrie", x))
	}
	return bt.BisectLeft(x)
}

func (bt *BinaryTrie) Find(x int) int {
	if bt.Count(x) == 0 {
		return -1
	}
	return bt.BisectLeft(x)
}

// 0<=k<bt.Size()
//
//	support negative index
func (bt *BinaryTrie) At(k int) int {
	if k < 0 {
		k += bt._size[0]
	}
	v := 0
	res := 0
	for i := bt._maxLog - 1; i > -1; i-- {
		l := (bt._lazy >> i) & 1
		lc := bt._edges[2*v]
		rc := bt._edges[2*v+1]
		if l == 1 {
			lc, rc = rc, lc
		}
		if lc == -1 {
			v = rc
			res |= 1 << i
			continue
		}
		if bt._size[lc] <= k {
			k -= bt._size[lc]
			v = rc
			res |= 1 << i
		} else {
			v = lc
		}
	}
	return res
}

func (bt *BinaryTrie) Min() int {
	return bt.At(0)
}

func (bt *BinaryTrie) Max() int {
	return bt.At(-1)
}

func (bt *BinaryTrie) XorAll(x int) {
	bt._lazy ^= x
}

func (bt *BinaryTrie) Has(x int) bool {
	return bt.Count(x) > 0
}

func (bt *BinaryTrie) Size() int {
	return bt._size[0]
}

func (bt *BinaryTrie) ForEach(f func(value, index int)) {
	queue := [][2]int{{0, 0}}
	for i := bt._maxLog - 1; i > -1; i-- {
		l := (bt._lazy >> i) & 1
		nextQueue := [][2]int{}
		for _, v := range queue {
			lc := bt._edges[2*v[0]]
			rc := bt._edges[2*v[0]+1]
			if l == 1 {
				lc, rc = rc, lc
			}
			if lc != -1 {
				nextQueue = append(nextQueue, [2]int{lc, 2 * v[1]})
			}
			if rc != -1 {
				nextQueue = append(nextQueue, [2]int{rc, 2*v[1] + 1})
			}
		}
		queue = nextQueue
	}

	i := 0
	for _, item := range queue {
		v, x := item[0], item[1]
		for j := 0; j < bt._endCount[v]; j++ {
			f(x, i)
			i++
		}
	}
}

func (bt *BinaryTrie) String() string {
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
