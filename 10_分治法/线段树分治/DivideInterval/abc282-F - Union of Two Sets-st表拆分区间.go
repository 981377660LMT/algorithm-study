// F - Union of Two Sets
// https://atcoder.jp/contests/abc282/tasks/abc282_f
// 交互题

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

	var n int32
	fmt.Fscan(in, &n)
	D := NewDivideIntervalByBinaryLift(n)
	mp := make(map[[2]int32]int32, D.Size())
	intervals := make([][2]int32, 0, D.Size())
	D.EnumerateJump(func(level, index int32) {
		start, end := D.JumpToRange(level, index)
		intervals = append(intervals, [2]int32{start, end})
		mp[[2]int32{start, end}] = int32(len(intervals) - 1)
	})

	fmt.Fprintln(out, len(mp))
	for _, interval := range intervals {
		fmt.Fprintln(out, interval[0]+1, interval[1])
	}
	out.Flush()

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var l, r int32
		fmt.Fscan(in, &l, &r)
		l--

		D.EnumerateRangeDangerously(l, r, func(level, index int32) {
			start, end := D.JumpToRange(level, index)
			id := mp[[2]int32{start, end}]
			fmt.Fprint(out, id+1, " ")
		})
		fmt.Fprintln(out)
		out.Flush()
	}
}

// DivideIntervalBinaryLift/反向st表
// 倍增拆分序列上的区间 `[start,end)`
// 一共拆分成[0,log]层，每层有n个元素.
// !jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n).
type DivideIntervalByBinaryLift struct {
	n, log int32
	size   int32
}

func NewDivideIntervalByBinaryLift(n int32) *DivideIntervalByBinaryLift {
	log := int32(bits.Len32(uint32(n))) - 1
	size := n * (log + 1)
	return &DivideIntervalByBinaryLift{n: n, log: log, size: size}
}

// O(logn)遍历[start,end)区间内的所有jump.
func (d *DivideIntervalByBinaryLift) EnumerateRange(start, end int32, f func(level, index int32)) {
	if start >= end {
		return
	}
	cur := start
	log := d.log
	len := end - start
	for k := log; k >= 0; k-- {
		if len&(1<<k) != 0 {
			f(k, cur)
			cur += 1 << k
			if cur >= end {
				return
			}
		}
	}
	f(0, cur)
}

// st表拆分区间.
// O(1)遍历[start,end)区间内的所有jump.
// !要求运算幂等(idempotent).
func (d *DivideIntervalByBinaryLift) EnumerateRangeDangerously(start, end int32, f func(level, index int32)) {
	if start >= end {
		return
	}
	k := int32(bits.Len32(uint32(end-start))) - 1
	f(k, start)
	f(k, end-(1<<k))
}

// 从高的jump开始下推信息，更新底部jump的答案.
// O(n*log(n)).
func (d *DivideIntervalByBinaryLift) PushDown(
	f func(parentLevel, parentIndex int32, childLevel, childIndex1, childIndex2 int32),
) {
	n, log := d.n, d.log
	for k := log - 1; k >= 0; k-- {
		for i := int32(0); i < n-(1<<k); i++ {
			f(k+1, i, k, i, i+(1<<k))
		}
	}
}

func (d *DivideIntervalByBinaryLift) Size() int32 { return d.size }
func (d *DivideIntervalByBinaryLift) Log() int32  { return d.log }

// 遍历所有jump结点.
func (d *DivideIntervalByBinaryLift) EnumerateJump(f func(level, index int32)) {
	n, log := d.n, d.log
	for k := log; k >= 0; k-- {
		for i := int32(0); i <= n-(1<<k); i++ {
			f(k, i)
		}
	}
}

func (d *DivideIntervalByBinaryLift) JumpToRange(level, index int32) (start, end int32) {
	start = index
	end = index + 1<<level
	return
}
