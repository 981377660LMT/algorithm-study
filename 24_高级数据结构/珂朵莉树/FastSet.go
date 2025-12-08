// 又叫做 64-ary tree
// !时间复杂度:O(log64n)
// https://zhuanlan.zhihu.com/p/107238627
// https://www.luogu.com.cn/blog/RuntimeErrror/ni-suo-fou-zhi-dao-di-shuo-ju-jie-gou-van-emde-boas-shu
// 使用场景:
// 1. 在存储IP地址的时候， 需要快速查找某个IP地址（2 ^32大小)是否在访问的列表中，
//    或者需要找到比这个IP地址大一点或者小一点的IP作为重新分配的IP。
// 2. 一条路上开了很多商店，用int来表示商店的位置（假设位置为1-256之间的数），
//    不断插入，删除商店，同时需要找到离某个商店最近的商店在哪里。
// 3. 代替链表(写起来比链表简单)

// !Insert/Erase/Prev/Next/Has/Enumerate
// https://maspypy.github.io/library/ds/fastset.hpp

// ! 注意频繁查找普通属性耗时(res.B), 把B写在代码里面或者定义成const,
// ! 会快很多(200ms->30ms), 编译器会优化

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// fs := NewFastSetFrom(66, func(i int) bool { return true })
	// fmt.Println(fs)
	P7912()
	// CF899E()
}

// P7912 [CSP-J 2021] 小熊的果篮
// https://www.luogu.com.cn/problem/P7912
// 连续排在一起的同一种水果称为一个“块”.
// 每次都把每一个“块”中最左边的水果同时挑出，组成一个果篮。重复这一操作，直至水果用完。
// !注意：取出水果后，如果左右两边的“块”中的水果种类相同，那么这两个“块”将会合并成一个“块”。
//
// set维护每个区间的长度,数组维护每个区间的编号，合并时启发式合并.
func P7912() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	fs := NewFastSetFrom(n, func(i int) bool { return true })
	indexes := make([]*Deque, n)
	for i := 0; i < n; i++ {
		indexes[i] = NewDeque2(2)
		indexes[i].Append(i)
	}

	merge := func(leftDeq, rightDeq **Deque) {
		if (*leftDeq).Size() < (*rightDeq).Size() {
			*leftDeq, *rightDeq = *rightDeq, *leftDeq
			for (*rightDeq).Size() > 0 {
				(*leftDeq).AppendLeft((*rightDeq).Pop())
			}
		} else {
			for (*rightDeq).Size() > 0 {
				(*leftDeq).Append((*rightDeq).PopLeft())
			}
		}
	}

	// 向左合并，将区间信息保存到区间起点left上.
	tryMergeToLeft := func(right int) {
		left := fs.Prev(right - 1)
		if left == -1 || nums[left] != nums[right] {
			return
		}
		fs.Erase(right)
		merge(&indexes[left], &indexes[right])
	}
	for i := 0; i < n; i++ {
		tryMergeToLeft(i)
	}

	var res [][]int
	for fs.Size() > 0 {
		curGroup := []int{}
		toErase := []int{}
		fs.Enumerate(0, n, func(left int) {
			curGroup = append(curGroup, indexes[left].PopLeft())
			if indexes[left].Size() == 0 {
				toErase = append(toErase, left)
			}
		})
		for _, v := range toErase {
			fs.Erase(v)
			right := fs.Next(v)
			if right < n {
				tryMergeToLeft(right) // !注意删除后，左右区间可能合并
			}
		}

		res = append(res, curGroup)
	}

	for _, group := range res {
		for _, v := range group {
			fmt.Fprint(out, v+1, " ")
		}
		fmt.Fprintln(out)
	}
}

// Segments Removal
// https://www.luogu.com.cn/problem/CF899E
// 每次删除长度最长的最左边的一段连续的数字区间，求删除次数
// 优先队列维护区间长度，fastset维护区间(也可以链表维护)
func CF899E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	size := make([]int, n)
	for i := 0; i < n; i++ {
		size[i] = 1
	}
	fs := NewFastSetFrom(n, func(i int) bool { return true })
	pq := NewHeap(
		func(a, b H) bool {
			if a.length == b.length {
				return a.start < b.start
			}
			return a.length > b.length
		},
		nil,
	)
	for i := 0; i < n; i++ {
		pq.Push(H{length: 1, start: i}) // 长度为1，起始位置为i
	}

	// 向左合并，将区间信息保存到区间起点left上.
	tryMergeToLeft := func(right int) {
		left := fs.Prev(right - 1)
		if left == -1 || nums[left] != nums[right] {
			return
		}
		fs.Erase(right)
		size[left] += size[right]
		pq.Push(H{length: size[left], start: left})
	}
	for i := 0; i < n; i++ {
		tryMergeToLeft(i)
	}

	res := 0
	for pq.Len() > 0 {
		item := pq.Pop()
		start := item.start
		if !fs.Has(start) {
			continue
		}
		res++
		fs.Erase(start)
		right := fs.Next(start)
		if right < n {
			tryMergeToLeft(right) // !注意删除后，左右区间可能合并
		}
	}

	fmt.Println(res)
}

func demo() {
	demo := func() {
		n := int(1e7)
		fs := NewFastSet(n)
		time1 := time.Now()
		for i := 0; i < n; i++ {
			fs.Insert(i)
			fs.Next(i)
			fs.Prev(i)
			fs.Has(i)
			fs.Erase(i)
			fs.Insert(i)
		}
		fmt.Println(time.Since(time1))
	}
	_ = demo

	// https://judge.yosupo.jp/problem/predecessor_problem
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	const N = 1e7 + 10
	set := NewFastSet(N)
	var s string
	fmt.Fscan(in, &s)
	for i, v := range s {
		if v == '1' {
			set.Insert(i)
		}
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		switch op {
		case 0:
			var k int
			fmt.Fscan(in, &k)
			set.Insert(k)
		case 1:
			var k int
			fmt.Fscan(in, &k)
			set.Erase(k)
		case 2:
			var k int
			fmt.Fscan(in, &k)
			if set.Has(k) {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
		case 3:
			var k int
			fmt.Fscan(in, &k)
			ceiling := set.Next(k)
			if ceiling < N {
				fmt.Fprintln(out, ceiling)
			} else {
				fmt.Fprintln(out, -1)
			}
		case 4:
			var k int
			fmt.Fscan(in, &k)
			floor := set.Prev(k)
			fmt.Fprintln(out, floor)

		}
	}
}

type FastSet struct {
	n, lg int
	seg   [][]int
	size  int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet) Clear() {
	fs.Enumerate(0, fs.n, func(i int) {
		fs.Erase(i)
	})
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.size
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}

//
//

type H = struct{ length, start int }

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}

		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}

type D = int
type Deque struct{ l, r []D }

func NewDeque2(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
