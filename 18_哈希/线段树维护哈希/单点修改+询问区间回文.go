package main

import "fmt"

func main() {
	// abcda
	// 5
	// palindrome? 1 5
	// palindrome? 1 1
	// change 4 b
	// palindrome? 1 5
	// palindrome? 2 4
	s := "abcda"
	nums := make([]int, len(s))
	for i := range s {
		nums[i] = int(s[i])
	}
	operations := [][3]int{
		{1, 1, 5},
		{1, 1, 1},
		{2, 4, int('b')},
		{1, 1, 5},
		{1, 2, 4},
	}

	fmt.Println(Subpalindromes(nums, operations))
}

// nums[i]>=0
// 1：palindrome? x y  :询问[x, y]是否为回文串
// 2：change 4 b       :把s[4]修改为b
func Subpalindromes(nums []int, operations [][3]int) []bool {
	n := len(nums)
	leaves0, leaves1 := make([]E, n), make([]E, n)
	for i := 0; i < n; i++ {
		leaves0[i] = CreateE(uint(nums[i]))
		leaves1[i] = CreateE(uint(nums[n-1-i]))
	}
	seg0, seg1 := NewSegmentTree(leaves0), NewSegmentTree(leaves1)

	// 0<=start<=end<=n
	isPalindrome := func(start, end int) bool {
		h1 := seg0.Query(start, end)
		h2 := seg1.Query(start, end)
		return h1.hash1 == h2.hash1 && h1.hash2 == h2.hash2
	}

	res := []bool{}
	for _, op := range operations {
		if op[0] == 1 {
			start, end := op[1]-1, op[2]
			res = append(res, isPalindrome(start, end))
		} else {
			i, v := op[1]-1, uint(op[2])
			seg0.Set(i, CreateE(v))
			seg1.Set(n-1-i, CreateE(v))
		}
	}
	return res
}

const N int = 1e5 + 10

var BASEPOW0 [N]uint
var BASEPOW1 [N]uint

func init() {
	BASEPOW0[0] = 1
	BASEPOW1[0] = 1
	for i := 1; i < N; i++ {
		BASEPOW0[i] = BASEPOW0[i-1] * BASE0
		BASEPOW1[i] = BASEPOW1[i-1] * BASE1
	}
}

// PointSetRangeHash
// 131/13331/1713302033171(回文素数)
const BASE0 = 131
const BASE1 = 13331

type E = struct {
	len          int
	hash1, hash2 uint
}

func CreateE(c uint) E {
	return E{len: 1, hash1: c, hash2: c}
}
func (*SegmentTreeHash) e() E { return E{} }
func (*SegmentTreeHash) op(a, b E) E {
	return E{
		len:   a.len + b.len,
		hash1: a.hash1*BASEPOW0[b.len] + b.hash1,
		hash2: a.hash2*BASEPOW1[b.len] + b.hash2,
	}
}

type SegmentTreeHash struct {
	n, size int
	seg     []E
	unit    E
}

func NewSegmentTree(leaves []E) *SegmentTreeHash {
	res := &SegmentTreeHash{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	res.unit = res.e()
	return res
}
func (st *SegmentTreeHash) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.unit
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeHash) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeHash) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTreeHash) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.unit
	}
	leftRes, rightRes := st.unit, st.unit
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}

func (st *SegmentTreeHash) QueryAll() E { return st.seg[1] }

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
