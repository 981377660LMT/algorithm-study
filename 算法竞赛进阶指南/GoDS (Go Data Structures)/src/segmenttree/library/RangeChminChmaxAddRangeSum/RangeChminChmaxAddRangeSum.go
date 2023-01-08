// https://judge.yosupo.jp/problem/range_chmin_chmax_add_range_sum
// 区间min/max/加法/区间和
// 0 left right x => 区间[left,right)的每个元素变为min(x,原值)
// 1 left right x => 区间[left,right)的每个元素变为max(x,原值)
// 2 left right x => 区间[left,right)的每个元素都加上x
// 3 left right => 求区间[left,right)的和
// !0<=left<right<=n
// 普通的线段树不能用了
// 需要用 SegmentTreeBeats

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

	var n, q int
	fmt.Fscan(in, &n, &q)

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	tree := NewSegmentTreeBeats(n)
	tree.Build(nums)

	for i := 0; i < q; i++ {
		var op, left, right, x int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &left, &right, &x)
			tree.RangeChmin(left, right, x)
		} else if op == 1 {
			fmt.Fscan(in, &left, &right, &x)
			tree.RangeChmax(left, right, x)
		} else if op == 2 {
			fmt.Fscan(in, &left, &right, &x)
			tree.RangeAdd(left, right, x)
		} else {
			fmt.Fscan(in, &left, &right)
			fmt.Fprintln(out, tree.GetSum(left, right))
		}
	}

}

const INF int = 1e18

type SegmentTreeBeats struct {
	n, log, size                                      int
	fmax, fmin, smax, smin, maxc, minc, add, upd, sum []int
	up, down, lt, rt                                  []int
}

func NewSegmentTreeBeats(n int) *SegmentTreeBeats {
	log := bits.Len(uint(n - 1))
	size := 1 << uint(log)
	fmax := make([]int, 2*size)
	fmin := make([]int, 2*size)
	smax := make([]int, 2*size)
	smin := make([]int, 2*size)
	maxc := make([]int, 2*size)
	minc := make([]int, 2*size)
	add := make([]int, 2*size)
	upd := make([]int, 2*size)
	sum := make([]int, 2*size)
	for i := range fmax {
		fmax[i] = -INF
		fmin[i] = INF
		smax[i] = -INF
		smin[i] = INF
		upd[i] = INF
	}
	up := make([]int, 0, 2*size)
	down := make([]int, 0, 2*size)
	lt := make([]int, 2*size)
	rt := make([]int, 2*size)
	for i := 0; i < size; i++ {
		lt[i+size] = i
		rt[i+size] = i + 1
	}
	for i := size - 1; i > -1; i-- {
		lt[i] = lt[i<<1]
		rt[i] = rt[i<<1|1]
	}

	return &SegmentTreeBeats{
		n:    n,
		log:  log,
		size: size,
		fmax: fmax,
		fmin: fmin,
		smax: smax,
		smin: smin,
		maxc: maxc,
		minc: minc,
		add:  add,
		upd:  upd,
		sum:  sum,
		up:   up,
		down: down,
		lt:   lt,
		rt:   rt,
	}
}

func (stb *SegmentTreeBeats) Build(arr []int) {
	for i, a := range arr {
		stb.fmax[i+stb.size] = a
		stb.fmin[i+stb.size] = a
		stb.sum[i+stb.size] = a
		stb.maxc[i+stb.size] = 1
		stb.minc[i+stb.size] = 1
	}

	for i := stb.size - 1; i > 0; i-- {
		stb.merge(i)
	}
}

// RangeAdd updates all elements in [left, right) by x.
//  0 <= left < right <= n
func (stb *SegmentTreeBeats) RangeChmax(left, right, x int) {
	stb.down = append(stb.down, 1)
	for len(stb.down) > 0 {
		k := stb.down[len(stb.down)-1]
		stb.down = stb.down[:len(stb.down)-1]
		if right <= stb.lt[k] || stb.rt[k] <= left || x <= stb.fmin[k] {
			continue
		}
		if left <= stb.lt[k] && stb.rt[k] <= right && x < stb.smin[k] {
			stb.chmin(k, x)
			continue
		}
		stb.downPropagate(k)
		stb.up = append(stb.up, k)
	}
	stb.upMerge()
}

// RangeChmin updates all elements in [left, right) to min(x, a[i])
//  0 <= left < right <= n
func (stb *SegmentTreeBeats) RangeChmin(left, right, x int) {
	stb.down = append(stb.down, 1)
	for len(stb.down) > 0 {
		k := stb.down[len(stb.down)-1]
		stb.down = stb.down[:len(stb.down)-1]
		if right <= stb.lt[k] || stb.rt[k] <= left || stb.fmax[k] <= x {
			continue
		}
		if left <= stb.lt[k] && stb.rt[k] <= right && stb.smax[k] < x {
			stb.chmax(k, x)
			continue
		}
		stb.downPropagate(k)
		stb.up = append(stb.up, k)
	}
	stb.upMerge()
}

// RangeAdd adds x to all elements in [left, right)
//  0 <= left < right <= n
func (stb *SegmentTreeBeats) RangeAdd(left, right, x int) {
	stb.down = append(stb.down, 1)
	for len(stb.down) > 0 {
		k := stb.down[len(stb.down)-1]
		stb.down = stb.down[:len(stb.down)-1]
		if right <= stb.lt[k] || stb.rt[k] <= left {
			continue
		}
		if left <= stb.lt[k] && stb.rt[k] <= right {
			stb.add_(k, x)
			continue
		}
		stb.downPropagate(k)
		stb.up = append(stb.up, k)
	}
	stb.upMerge()
}

// RangeUpdate updates all elements in [left, right) to x.
//  0 <= left < right <= n
func (stb *SegmentTreeBeats) RangeUpdate(left, right, x int) {
	stb.down = append(stb.down, 1)
	for len(stb.down) > 0 {
		k := stb.down[len(stb.down)-1]
		stb.down = stb.down[:len(stb.down)-1]
		if right <= stb.lt[k] || stb.rt[k] <= left {
			continue
		}
		if left <= stb.lt[k] && stb.rt[k] <= right {
			stb.update(k, x)
			continue
		}
		stb.downPropagate(k)
		stb.up = append(stb.up, k)
	}
	stb.upMerge()
}

//  0 <= left < right <= n
func (stb *SegmentTreeBeats) GetMax(left, right int) int {
	stb.down = append(stb.down, 1)
	v := -INF
	for len(stb.down) > 0 {
		k := stb.down[len(stb.down)-1]
		stb.down = stb.down[:len(stb.down)-1]
		if right <= stb.lt[k] || stb.rt[k] <= left {
			continue
		}
		if left <= stb.lt[k] && stb.rt[k] <= right {
			v = max(v, stb.fmax[k])
			continue
		}
		stb.downPropagate(k)
	}
	return v
}

//  0 <= left < right <= n
func (stb *SegmentTreeBeats) GetMin(left, right int) int {
	stb.down = append(stb.down, 1)
	v := INF
	for len(stb.down) > 0 {
		k := stb.down[len(stb.down)-1]
		stb.down = stb.down[:len(stb.down)-1]
		if right <= stb.lt[k] || stb.rt[k] <= left {
			continue
		}
		if left <= stb.lt[k] && stb.rt[k] <= right {
			v = min(v, stb.fmin[k])
			continue
		}
		stb.downPropagate(k)
	}
	return v
}

//  0 <= left < right <= n
func (stb *SegmentTreeBeats) GetSum(left, right int) int {
	stb.down = append(stb.down, 1)
	var v int
	for len(stb.down) > 0 {
		k := stb.down[len(stb.down)-1]
		stb.down = stb.down[:len(stb.down)-1]
		if right <= stb.lt[k] || stb.rt[k] <= left {
			continue
		}
		if left <= stb.lt[k] && stb.rt[k] <= right {
			v += stb.sum[k]
			continue
		}
		stb.downPropagate(k)
	}
	return v
}

func (stb *SegmentTreeBeats) merge(k int) {
	stb.sum[k] = stb.sum[k<<1] + stb.sum[k<<1|1]

	if stb.fmax[k<<1] < stb.fmax[k<<1|1] {
		stb.fmax[k] = stb.fmax[k<<1|1]
		stb.maxc[k] = stb.maxc[k<<1|1]
		stb.smax[k] = max(stb.fmax[k<<1], stb.smax[k<<1|1])
	} else if stb.fmax[k<<1] > stb.fmax[k<<1|1] {
		stb.fmax[k] = stb.fmax[k<<1]
		stb.maxc[k] = stb.maxc[k<<1]
		stb.smax[k] = max(stb.smax[k<<1], stb.fmax[k<<1|1])
	} else {
		stb.fmax[k] = stb.fmax[k<<1]
		stb.maxc[k] = stb.maxc[k<<1] + stb.maxc[k<<1|1]
		stb.smax[k] = max(stb.smax[k<<1], stb.smax[k<<1|1])
	}

	if stb.fmin[k<<1] > stb.fmin[k<<1|1] {
		stb.fmin[k] = stb.fmin[k<<1|1]
		stb.minc[k] = stb.minc[k<<1|1]
		stb.smin[k] = min(stb.fmin[k<<1], stb.smin[k<<1|1])
	} else if stb.fmin[k<<1] < stb.fmin[k<<1|1] {
		stb.fmin[k] = stb.fmin[k<<1]
		stb.minc[k] = stb.minc[k<<1]
		stb.smin[k] = min(stb.smin[k<<1], stb.fmin[k<<1|1])
	} else {
		stb.fmin[k] = stb.fmin[k<<1]
		stb.minc[k] = stb.minc[k<<1] + stb.minc[k<<1|1]
		stb.smin[k] = min(stb.smin[k<<1], stb.smin[k<<1|1])
	}
}

func (stb *SegmentTreeBeats) propagate(k int) {
	if k >= stb.size {
		return
	}

	if stb.upd[k] != INF {
		stb.update(k<<1, stb.upd[k])
		stb.update(k<<1|1, stb.upd[k])
		stb.upd[k] = INF
		return
	}

	if stb.add[k] != 0 {
		stb.add_(k<<1, stb.add[k])
		stb.add_(k<<1|1, stb.add[k])
		stb.add[k] = 0
	}

	if stb.fmax[k] < stb.fmax[k<<1] {
		stb.chmax(k<<1, stb.fmax[k])
	}
	if stb.fmin[k<<1] < stb.fmin[k] {
		stb.chmin(k<<1, stb.fmin[k])
	}
	if stb.fmax[k] < stb.fmax[k<<1|1] {
		stb.chmax(k<<1|1, stb.fmax[k])
	}
	if stb.fmin[k<<1|1] < stb.fmin[k] {
		stb.chmin(k<<1|1, stb.fmin[k])
	}
}

func (stb *SegmentTreeBeats) upMerge() {
	for len(stb.up) > 0 {
		popped := stb.up[len(stb.up)-1]
		stb.up = stb.up[:len(stb.up)-1]
		stb.merge(popped)
	}
}

func (stb *SegmentTreeBeats) downPropagate(k int) {
	stb.propagate(k)
	stb.down = append(stb.down, k<<1, k<<1|1)
}

func (stb *SegmentTreeBeats) update(k, x int) {
	stb.fmax[k] = x
	stb.fmin[k] = x
	stb.smax[k] = -INF
	stb.smin[k] = INF
	stb.maxc[k] = stb.rt[k] - stb.lt[k]
	stb.minc[k] = stb.rt[k] - stb.lt[k]
	stb.upd[k] = x
	stb.add[k] = 0
	stb.sum[k] = x * (stb.rt[k] - stb.lt[k])
}

func (stb *SegmentTreeBeats) add_(k, x int) {
	stb.fmax[k] += x
	if stb.smax[k] != -INF {
		stb.smax[k] += x
	}
	stb.fmin[k] += x
	if stb.smin[k] != INF {
		stb.smin[k] += x
	}
	stb.sum[k] += x * (stb.rt[k] - stb.lt[k])
	if stb.upd[k] != INF {
		stb.upd[k] += x
	} else {
		stb.add[k] += x
	}
}

func (stb *SegmentTreeBeats) chmax(k, x int) {
	stb.sum[k] += (x - stb.fmax[k]) * stb.maxc[k]
	if stb.fmax[k] == stb.fmin[k] {
		stb.fmax[k] = x
		stb.fmin[k] = x
	} else if stb.fmax[k] == stb.smin[k] {
		stb.fmax[k] = x
		stb.smin[k] = x
	} else {
		stb.fmax[k] = x
	}
	if stb.upd[k] != INF && x < stb.upd[k] {
		stb.upd[k] = x
	}
}

func (stb *SegmentTreeBeats) chmin(k, x int) {
	stb.sum[k] += (x - stb.fmin[k]) * stb.minc[k]
	if stb.fmin[k] == stb.fmax[k] {
		stb.fmin[k] = x
		stb.fmax[k] = x
	} else if stb.fmin[k] == stb.smax[k] {
		stb.fmin[k] = x
		stb.smax[k] = x
	} else {
		stb.fmin[k] = x
	}
	if stb.upd[k] != INF && stb.upd[k] < x {
		stb.upd[k] = x
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
