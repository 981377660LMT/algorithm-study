// SegmentTreePointSetRangeContinuousCount-区间连续1的个数

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type SegmentTreePointSetRangeContinuousCount struct {
	data []bool
	seg  *SegmentTree
}

func NewSegmentTreePointSetRangeContinuousCount(n int, f func(i int) bool) *SegmentTreePointSetRangeContinuousCount {
	data := make([]bool, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
	}
	seg := NewSegmentTree(n, func(i int) E { return of(data[i]) })
	return &SegmentTreePointSetRangeContinuousCount{data: data, seg: seg}
}

func (st *SegmentTreePointSetRangeContinuousCount) Get(index int) bool {
	return st.data[index]
}

func (st *SegmentTreePointSetRangeContinuousCount) Set(index int, value bool) {
	if st.data[index] == value {
		return
	}
	st.data[index] = value
	st.seg.Set(index, of(value))
}

func (st *SegmentTreePointSetRangeContinuousCount) Query(start, end int, target bool) int {
	res := st.seg.Query(start, end)
	if target {
		return res.onesCount
	}
	return res.zerosCount
}

const INF int = 1e18

// SegmentTreePointSetRangeContinuousCount

type E = struct {
	leftOne, rightOne     bool
	zerosCount, onesCount int
	len                   int
}

func of(value bool) E {
	if value {
		return E{leftOne: true, rightOne: true, onesCount: 1, len: 1}
	}
	return E{zerosCount: 1, len: 1}
}

func (*SegmentTree) e() E { return E{} }
func (*SegmentTree) op(a, b E) E {
	if a.len == 0 {
		return b
	}
	if b.len == 0 {
		return a
	}
	res := E{
		leftOne:    a.leftOne,
		rightOne:   b.rightOne,
		onesCount:  a.onesCount + b.onesCount,
		zerosCount: a.zerosCount + b.zerosCount,
		len:        a.len + b.len,
	}
	if a.rightOne && b.leftOne {
		res.onesCount--
	}
	if !a.rightOne && !b.leftOne {
		res.zerosCount--
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

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(n int, f func(int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int, value E) {
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
func (st *SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
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
func (st *SegmentTree) QueryAll() E { return st.seg[1] }
func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

func bruteOnes(arr []bool, l, r int) int {
	if l >= r {
		return 0
	}
	cnt := 0
	i := l
	for i < r {
		if arr[i] {

			cnt++
			for i < r && arr[i] {
				i++
			}
		} else {
			i++
		}
	}
	return cnt
}

func main() {
	rand.Seed(time.Now().UnixNano())
	const rounds = 5000
	for tc := 1; tc <= rounds; tc++ {
		n := rand.Intn(100) + 1  // 1..100
		q := rand.Intn(2000) + 1 // 1..2000
		initArr := make([]bool, n)
		for i := 0; i < n; i++ {
			initArr[i] = rand.Intn(2) == 1
		}

		seg := NewSegmentTreePointSetRangeContinuousCount(n, func(i int) bool { return initArr[i] })
		naive := append([]bool(nil), initArr...)

		for step := 0; step < q; step++ {
			if rand.Intn(2) == 0 { // 点修改
				idx := rand.Intn(n)
				val := rand.Intn(2) == 1
				seg.Set(idx, val)
				naive[idx] = val
			} else { // 区间查询
				l := rand.Intn(n)
				r := rand.Intn(n-l) + l + 1
				got := seg.Query(l, r, true)
				exp := bruteOnes(naive, l, r)
				if got != exp {
					fmt.Println("Mismatch detected!")
					fmt.Printf("n=%d  step=%d  l=%d r=%d\n", n, step, l, r)
					fmt.Printf("expected=%d  got=%d\n", exp, got)
					fmt.Println("array=", naive)
					return
				}
			}
		}
	}
	fmt.Println("All tests passed!")
}
