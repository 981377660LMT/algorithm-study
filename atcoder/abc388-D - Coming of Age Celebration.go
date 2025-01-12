// D - Coming of Age Celebration
// https://atcoder.jp/contests/abc388/tasks/abc388_d
//
// 在一个星球上， 有 N 名外星人，且全员未成年。
// 第 i 名外星人目前拥有 A_i 个石头，并将在恰好 i 年后成年。
// 在这个星球上，每当有人成年时，所有已成年的且持有 1 个以上石头的外星人，都会向这位即将成年的人赠送 1 个石头，以示祝贺。
// 请计算在 N 年后，每位外星人持有的石头数量。
// 请注意，未来不会再有新的外星人出生。
//
// !等价于，按照成年顺序排序后，每个人需要给后面的人赠送石头。

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	seg := NewSegmentTreeDual32(n, func(i int32) Id { return 0 })
	for i := int32(0); i < n; i++ {
		nums[i] += seg.Get(i)
		cur := nums[i]
		k := min(cur, int(n-1-i))
		seg.Update(i+1, i+int32(k)+1, 1)
		nums[i] -= k
	}

	for i := int32(0); i < n; i++ {
		fmt.Fprint(out, nums[i], " ")
	}
}

// RangeAssignPointGet

type Id = int

const COMMUTATIVE = true

func (*SegmentTreeDual32) id() Id                 { return 0 }
func (*SegmentTreeDual32) composition(f, g Id) Id { return f + g }

type SegmentTreeDual32 struct {
	n         int32
	size, log int32
	lazy      []Id
	unit      Id
}

func NewSegmentTreeDual32(n int32, f func(i int32) Id) *SegmentTreeDual32 {
	res := &SegmentTreeDual32{}
	log := int32(1)
	for 1<<log < n {
		log++
	}
	size := int32(1 << log)
	lazy := make([]Id, 2*size)
	unit := res.id()
	for i := int32(0); i < size; i++ {
		lazy[i] = unit
	}
	for i := int32(0); i < n; i++ {
		lazy[size+i] = f(i)
	}
	res.n = n
	res.size = size
	res.log = log
	res.lazy = lazy
	res.unit = unit
	return res
}
func (seg *SegmentTreeDual32) Get(index int32) Id {
	index += seg.size
	for i := seg.log; i > 0; i-- {
		seg.propagate(index >> i)
	}
	return seg.lazy[index]
}
func (seg *SegmentTreeDual32) Set(index int32, value Id) {
	index += seg.size
	for i := seg.log; i > 0; i-- {
		seg.propagate(index >> i)
	}
	seg.lazy[index] = value
}
func (seg *SegmentTreeDual32) GetAll() []Id {
	for i := int32(0); i < seg.size; i++ {
		seg.propagate(i)
	}
	res := make([]Id, seg.n)
	copy(res, seg.lazy[seg.size:seg.size+seg.n])
	return res
}
func (seg *SegmentTreeDual32) Update(left, right int32, value Id) {
	if left < 0 {
		left = 0
	}
	if right > seg.n {
		right = seg.n
	}
	if left >= right {
		return
	}
	left += seg.size
	right += seg.size
	if !COMMUTATIVE {
		for i := seg.log; i > 0; i-- {
			if (left>>i)<<i != left {
				seg.propagate(left >> i)
			}
			if (right>>i)<<i != right {
				seg.propagate((right - 1) >> i)
			}
		}
	}
	for left < right {
		if left&1 > 0 {
			seg.lazy[left] = seg.composition(value, seg.lazy[left])
			left++
		}
		if right&1 > 0 {
			right--
			seg.lazy[right] = seg.composition(value, seg.lazy[right])
		}
		left >>= 1
		right >>= 1
	}
}
func (seg *SegmentTreeDual32) propagate(k int32) {
	if seg.lazy[k] != seg.unit {
		seg.lazy[k<<1] = seg.composition(seg.lazy[k], seg.lazy[k<<1])
		seg.lazy[k<<1|1] = seg.composition(seg.lazy[k], seg.lazy[k<<1|1])
		seg.lazy[k] = seg.unit
	}
}
func (st *SegmentTreeDual32) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := int32(0); i < st.n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(fmt.Sprint(st.Get(i)))
	}
	buf.WriteByte(']')
	return buf.String()
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
