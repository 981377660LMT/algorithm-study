// https://leetcode.cn/problems/find-the-maximum-number-of-marked-indices/description/
// https://atcoder.jp/contests/abc388/tasks/abc388_g

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)
	M := NewMaxNumOfPair(nums)
	for i := 0; i < q; i++ {
		var start, end int32
		fmt.Fscan(in, &start, &end)
		start--
		fmt.Fprintln(out, M.Query(start, end))
	}
}

type MaxNumOfPair struct {
	n   int32
	seg *SparseTableFast
}

func NewMaxNumOfPair(sortedNums []int) *MaxNumOfPair {
	n := int32(len(sortedNums))
	nexts := make([]int32, n) // !每个元素距离右侧第一个大于等于它两倍的元素的下标.
	ptr := int32(0)
	for i := int32(0); i < n; i++ {
		ptr = max32(ptr, i+1)
		for ptr < n && sortedNums[i]*2 > sortedNums[ptr] {
			ptr++
		}
		nexts[i] = ptr
	}

	seg := NewSparseTableFast(
		n,
		func(i int32) int32 { return nexts[i] - i },
		func() int32 { return 0 },
		func(a, b int32) int32 { return max32(a, b) },
	)
	return &MaxNumOfPair{n: n, seg: seg}
}

func (m *MaxNumOfPair) Query(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > m.n {
		end = m.n
	}
	if start >= end {
		return 0
	}

	right := MaxRight32(
		0,
		// 能否配k对.
		func(k int32) bool {
			k--
			if 2*k > end-start {
				return false
			}

			maxDist := m.seg.Query(start, start+k)
			return maxDist <= end-start-k
		},
		m.n/2+1,
	)

	return right - 1
}

type S = int32

// Static RMQ, O(n)预处理, O(1)查询.
type SparseTableFast struct {
	n        int32
	leaves   []S
	pre, suf []S
	st       *SparseTable
	data     []int32
	e        func() S
	op       func(S, S) S
}

func NewSparseTableFast(n int32, f func(int32) S, e func() S, op func(S, S) S) *SparseTableFast {
	res := &SparseTableFast{}
	bNum := n >> 4
	leaves := make([]S, n)
	for i := int32(0); i < n; i++ {
		leaves[i] = f(i)
	}
	pre, suf := append(leaves[:0:0], leaves...), append(leaves[:0:0], leaves...)
	for i := int32(1); i < n; i++ {
		if i&15 != 0 {
			pre[i] = op(pre[i-1], leaves[i])
		}
	}
	for i := n - 1; i > 0; i-- {
		if i&15 != 0 {
			suf[i-1] = op(leaves[i-1], suf[i])
		}
	}
	st := NewSparseTable(bNum, func(i int32) S { return suf[i<<4] }, e, op)

	// 处理长度小于或等于16的查询
	// 在区间 [i, i+16) 内，如果 i+j 的位置上的值是 [i, i+j] 这个子区间的最小值，那么就将 j-th 位设置为1
	data := make([]int32, n)
	stack := int32(0)
	for i := n - 1; i >= 0; i-- {
		stack = (stack << 1) & 65535
		for stack > 0 {
			k := lowbit(stack)
			tmp := op(leaves[i], leaves[i+k])
			if tmp != leaves[i] {
				break
			}
			stack &= ^(1 << k)
		}
		stack |= 1
		data[i] = stack
	}
	res.n = n
	res.leaves = leaves
	res.pre, res.suf = pre, suf
	res.st = st
	res.data = data
	res.e = e
	res.op = op
	return res
}

func NewSparseTableFastFrom(leaves []S, e func() S, op func(S, S) S) *SparseTableFast {
	return NewSparseTableFast(int32(len(leaves)), func(i int32) S { return leaves[i] }, e, op)
}

func (st *SparseTableFast) Query(start, end int32) S {
	if start >= end {
		return st.e()
	}
	if end-start <= 16 {
		d := st.data[start] & ((1 << (end - start)) - 1)
		return st.leaves[start+topbit(d)]
	}
	end--
	a, b := start>>4, end>>4
	x := st.st.Query(a+1, b)
	x = st.op(st.suf[start], x)
	x = st.op(x, st.pre[end])
	return x
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *SparseTableFast) MaxRight(left int32, check func(e S) bool) int32 {
	if left == ds.n {
		return ds.n
	}
	ok, ng := left, ds.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(ds.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (ds *SparseTableFast) MinLeft(right int32, check func(e S) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(ds.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// SparseTable 稀疏表: st[j][i] 表示区间 [i, i+2^j) 的贡献值.
type SparseTable struct {
	st [][]S
	e  func() S
	op func(S, S) S
	n  int32
}

func NewSparseTable(n int32, f func(int32) S, e func() S, op func(S, S) S) *SparseTable {
	res := &SparseTable{}

	b := bits.Len32(uint32(n))
	st := make([][]S, b)
	for i := range st {
		st[i] = make([]S, n)
	}
	for i := int32(0); i < n; i++ {
		st[0][i] = f(i)
	}
	for i := 1; i < b; i++ {
		for j := int32(0); j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	res.st = st
	res.e = e
	res.op = op
	res.n = n
	return res
}

func NewSparseTableFrom(leaves []S, e func() S, op func(S, S) S) *SparseTable {
	return NewSparseTable(int32(len(leaves)), func(i int32) S { return leaves[i] }, e, op)
}

// 查询区间 [start, end) 的贡献值.
func (st *SparseTable) Query(start, end int32) S {
	if start >= end {
		return st.e()
	}
	b := bits.Len(uint(end-start)) - 1 // log2
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *SparseTable) MaxRight(left int32, check func(e S) bool) int32 {
	if left == ds.n {
		return ds.n
	}
	ok, ng := left, ds.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(ds.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (ds *SparseTable) MinLeft(right int32, check func(e S) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(ds.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func topbit(x int32) int32 {
	if x == 0 {
		return -1
	}
	return int32(bits.Len32(uint32(x)) - 1)
}

func lowbit(x int32) int32 {
	if x == 0 {
		return -1
	}
	return int32(bits.TrailingZeros32(uint32(x)))
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

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含,使用时需要right-1.
// right<=upper.
func MaxRight32(left int32, check func(right int32) bool, upper int32) int32 {
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
