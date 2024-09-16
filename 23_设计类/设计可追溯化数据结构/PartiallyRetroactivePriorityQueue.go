package main

import (
	"bufio"
	"fmt"
	"os"
)

// G - Dynamic Scheduling
// https://atcoder.jp/contests/abc363/tasks/abc363_g
// 给定两个长度为 n 的序列D和P。
// 有 Q 个操作，每个操作形如i x y，将 D[i] 改为 x ，将 P[i] 改为 y 。
// 每次修改完成之后，回答如下问题的答案：
// !有 n 个任务，每天可以做一个任务，如果一个任务在 D[i] 天前被解决，则可以获得奖励 P[i] 。
// 问可以得到的最大奖励。
// n,q<=1e5, 1<=D[i]<=n, 1<=P[i]<=1e9.
//
// 如果不带修改，那么只需要维护一个优先队列即可，弹出过期的任务，每天完成P最大的任务即可。
// 带修改，可以使用部分可追溯优先队列。
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, Q int32
	fmt.Fscan(in, &N, &Q)
	D := make([]int32, N)
	P := make([]int32, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &D[i])
		D[i]-- // 天数从0开始
	}
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &P[i])
	}

	type modify struct{ i, d, p int32 }
	ms := make([]modify, Q)
	for i := int32(0); i < Q; i++ {
		fmt.Fscan(in, &ms[i].i, &ms[i].d, &ms[i].p)
		ms[i].i--
		ms[i].d--
	}

	counter := make([]int32, N)
	for i := range counter {
		counter[i] = 1
	}
	for _, d := range D {
		counter[d]++
	}
	for _, q := range ms {
		counter[q.d]++
	}
	for i := N - 1; i > 0; i-- {
		counter[i-1] += counter[i]
	}

	pq := NewPartiallyRetroactivePriorityQueue(counter[0], func(a, b int32) bool { return a < b }, 0, 1e9+10)
	for i := int32(0); i < N; i++ {
		counter[i]--
		pq.SetPop(counter[i])
	}

	res := 0
	update := func(op ReturnType[int32]) {
		for _, e := range op.Insert {
			res -= int(e)
		}
		for _, e := range op.Erase {
			res += int(e)
		}
	}

	idx := make([]int32, N)
	for i := int32(0); i < N; i++ {
		counter[D[i]]--
		idx[i] = counter[D[i]]
		update(pq.SetPush(idx[i], P[i]))
		res += int(P[i])
	}

	for _, q := range ms {
		c, x, y := q.i, q.d, q.p
		update(pq.SetNoOp(idx[c]))
		res -= int(P[c])
		counter[x]--
		idx[c] = counter[x]
		P[c] = y
		update(pq.SetPush(idx[c], y))
		res += int(P[c])
		fmt.Fprintln(out, res)
	}
}

type ReturnType[T any] struct {
	Insert []T
	Erase  []T
}

type treeNode[T any] struct {
	ssum, smin int32
	qmax, dmin T
}

// PartiallyRetroactivePriorityQueue/部分可追溯优先队列.
// 可以更新过去所有版本，但是只能查询最新版本的变化。
// https://atcoder.jp/contests/abc363/editorial/10496
type PartiallyRetroactivePriorityQueue[T any] struct {
	n, m int32
	less func(a, b T) bool
	tree []treeNode[T]
}

func NewPartiallyRetroactivePriorityQueue[T any](maxTime int32, less func(a, b T) bool, min, max T) *PartiallyRetroactivePriorityQueue[T] {
	m := int32(1)
	for m < maxTime+1 {
		m <<= 1
	}
	tree := make([]treeNode[T], 2*m)
	for i := range tree {
		tree[i] = treeNode[T]{qmax: min, dmin: max}
	}
	tree[0].qmax = min
	tree[0].dmin = max
	return &PartiallyRetroactivePriorityQueue[T]{n: maxTime, m: m, less: less, tree: tree}
}

// 在第 i 次操作后插入将元素 x 入堆的操作.
func (q *PartiallyRetroactivePriorityQueue[T]) SetPush(i int32, x T) ReturnType[T] {
	res := q.SetNoOp(i)
	i += q.m
	q.tree[i].dmin = x
	q.updateD(i)
	q.incrementalUpdate(i, &res)
	return res
}

// 在第 i 次操作后插入弹出最小元素的操作。
func (q *PartiallyRetroactivePriorityQueue[T]) SetPop(i int32) ReturnType[T] {
	res := q.SetNoOp(i)
	i += q.m
	q.decrementalUpdate(i, &res)
	return res
}

// 删除第i次操作.
func (q *PartiallyRetroactivePriorityQueue[T]) SetNoOp(i int32) ReturnType[T] {
	res := ReturnType[T]{}
	i += q.m
	if q.tree[i].ssum == -1 {
		q.incrementalUpdate(i, &res)
	} else if q.less(q.tree[0].qmax, q.tree[i].qmax) {
		res.Erase = append(res.Erase, q.tree[i].qmax)
		q.tree[i].qmax = q.tree[0].qmax
		q.updateQ(i)
	} else if q.less(q.tree[i].dmin, q.tree[0].dmin) {
		q.tree[i].dmin = q.tree[0].dmin
		q.updateD(i)
		q.decrementalUpdate(i, &res)
	}
	return res
}

func (q *PartiallyRetroactivePriorityQueue[T]) updateS(i int32) {
	for i >>= 1; i > 0; i >>= 1 {
		q.tree[i].ssum = q.tree[i<<1].ssum + q.tree[i<<1|1].ssum
		q.tree[i].smin = min32(q.tree[i<<1].smin, q.tree[i<<1].ssum+q.tree[i<<1|1].smin)
	}
}

func (q *PartiallyRetroactivePriorityQueue[T]) updateQ(i int32) {
	for i >>= 1; i > 0; i >>= 1 {
		q.tree[i].qmax = q.maxOp(q.tree[i<<1].qmax, q.tree[i<<1|1].qmax)
	}
}

func (q *PartiallyRetroactivePriorityQueue[T]) updateD(i int32) {
	for i >>= 1; i > 0; i >>= 1 {
		q.tree[i].dmin = q.minOp(q.tree[i<<1].dmin, q.tree[i<<1|1].dmin)
	}
}

func (q *PartiallyRetroactivePriorityQueue[T]) incrementalUpdate(i int32, ret *ReturnType[T]) {
	q.tree[i].ssum++
	q.updateS(i)
	s := q.tree[1].ssum - 1
	k := int32(1)
	for k < q.m {
		k <<= 1
		if q.tree[k].ssum+q.tree[k|1].smin == s {
			s -= q.tree[k].ssum
			k |= 1
		}
	}
	if k == q.m {
		return
	}
	c := int32(0)
	r := int32(len(q.tree))
	for k < r {
		if k&1 == 1 {
			if q.less(q.tree[k].dmin, q.tree[c].dmin) {
				c = k
			}
			k++
		}
		k >>= 1
		r >>= 1
	}
	if c == 0 {
		panic("bug")
	}
	for c < q.m {
		c <<= 1
		if q.less(q.tree[c|1].dmin, q.tree[c].dmin) {
			c |= 1
		}
	}
	ret.Insert = append(ret.Insert, q.tree[c].dmin)
	q.tree[c].ssum = 0
	q.tree[c].qmax = q.tree[c].dmin
	q.tree[c].dmin = q.tree[0].dmin
	q.updateS(c)
	q.updateQ(c)
	q.updateD(c)
}

func (q *PartiallyRetroactivePriorityQueue[T]) decrementalUpdate(i int32, ret *ReturnType[T]) {
	q.tree[i].ssum--
	q.updateS(i)
	s := q.tree[1].ssum
	k := int32(1)
	for k < q.m {
		k <<= 1
		if s != q.tree[k].smin {
			s -= q.tree[k].ssum
			k |= 1
		}
	}
	c := int32(0)
	for k > 0 {
		if k&1 == 1 {
			if q.less(q.tree[c].qmax, q.tree[k-1].qmax) {
				c = k - 1
			}
		}
		k >>= 1
	}
	if c == 0 {
		return
	}
	for c < q.m {
		c <<= 1
		if q.less(q.tree[c].qmax, q.tree[c|1].qmax) {
			c |= 1
		}
	}
	ret.Erase = append(ret.Erase, q.tree[c].qmax)
	q.tree[c].ssum = 1
	q.tree[c].dmin = q.tree[c].qmax
	q.tree[c].qmax = q.tree[0].qmax
	q.updateS(c)
	q.updateQ(c)
	q.updateD(c)
}

func (q *PartiallyRetroactivePriorityQueue[T]) minOp(x, y T) T {
	if q.less(x, y) {
		return x
	}
	return y
}

func (q *PartiallyRetroactivePriorityQueue[T]) maxOp(x, y T) T {
	if q.less(x, y) {
		return y
	}
	return x
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
