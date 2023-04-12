// # https://natsugiri.hatenablog.com/entry/2016/10/10/035445

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/double_ended_priority_queue
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	pq := NewPriorityQueue(
		func(a, b PqItem) bool { return a < b },
		nums,
	)

	for i := 0; i < q; i++ {
		var t, x int
		fmt.Fscan(in, &t)
		if t == 0 {
			fmt.Fscan(in, &x)
			pq.Push(x)
		} else if t == 1 {

			fmt.Fprintln(out, pq.PopMin())
		} else {
			fmt.Fprintln(out, pq.PopMax())
		}
	}
}

type PqItem = int

type PriorityQueue struct {
	less func(a, b PqItem) bool
	data []PqItem
}

func NewPriorityQueue(less func(a, b PqItem) bool, data []PqItem) *PriorityQueue {
	res := &PriorityQueue{less: less, data: append(data[:0:0], data...)}
	if len(data) > 0 {
		res.heapify()
	}
	return res
}

func (pq *PriorityQueue) Push(x PqItem) {
	k := len(pq.data)
	pq.data = append(pq.data, x)
	pq.pushUp(k, 1)
}

func (pq *PriorityQueue) PopMin() (res PqItem) {
	res = pq.Min()
	if len(pq.data) < 3 {
		pq.data = pq.data[:len(pq.data)-1]
		return
	}
	pq.data[1] = pq.data[len(pq.data)-1]
	pq.data = pq.data[:len(pq.data)-1]
	k := pq.pushDown(1)
	pq.pushUp(k, 1)
	return
}

func (pq *PriorityQueue) PopMax() (res PqItem) {
	res = pq.Max()
	if len(pq.data) < 2 {
		pq.data = pq.data[:len(pq.data)-1]
		return
	}
	pq.data[0] = pq.data[len(pq.data)-1]
	pq.data = pq.data[:len(pq.data)-1]
	k := pq.pushDown(0)
	pq.pushUp(k, 1)
	return
}

func (pq *PriorityQueue) Min() PqItem {
	if len(pq.data) < 2 {
		return pq.data[0]
	}
	return pq.data[1]
}

func (pq *PriorityQueue) Max() PqItem {
	return pq.data[0]
}

func (pq *PriorityQueue) Len() int {
	return len(pq.data)
}

func (pq *PriorityQueue) Empty() bool {
	return len(pq.data) == 0
}

func (pq *PriorityQueue) heapify() {
	for i := pq.Len() - 1; i >= 0; i-- {
		if i&1 != 0 && pq.less(pq.data[i-1], pq.data[i]) {
			pq.data[i-1], pq.data[i] = pq.data[i], pq.data[i-1]
		}
		k := pq.pushDown(i)
		pq.pushUp(k, i)
	}
}

func (pq *PriorityQueue) pushDown(k int) int {
	n := pq.Len()
	if k&1 != 0 { // min heap
		for k<<1|1 < n {
			c := 2*k + 3
			if n <= c || pq.less(pq.data[c-2], pq.data[c]) {
				c -= 2
			}
			if c < n && pq.less(pq.data[c], pq.data[k]) {
				pq.data[k], pq.data[c] = pq.data[c], pq.data[k]
				k = c
			} else {
				break
			}
		}
	} else { // max heap
		for 2*k+2 < n {
			c := 2*k + 4
			if n <= c || pq.less(pq.data[c], pq.data[c-2]) {
				c -= 2
			}
			if c < n && pq.less(pq.data[k], pq.data[c]) {
				pq.data[k], pq.data[c] = pq.data[c], pq.data[k]
				k = c
			} else {
				break
			}
		}
	}
	return k
}

func (pq *PriorityQueue) pushUp(k int, root int) int {
	if a, b := k&^1, k|1; b < pq.Len() && pq.less(pq.data[a], pq.data[b]) {
		pq.data[a], pq.data[b] = pq.data[b], pq.data[a]
		k ^= 1
	}
	p := 0
	for root < k {
		p = ((k >> 1) - 1) &^ 1 // parent(k)
		if !pq.less(pq.data[p], pq.data[k]) {
			break
		}
		pq.data[p], pq.data[k] = pq.data[k], pq.data[p]
		k = p
	}
	for root < k {
		p = ((k>>1)-1)&^1 | 1
		if !pq.less(pq.data[k], pq.data[p]) {
			break
		}
		pq.data[p], pq.data[k] = pq.data[k], pq.data[p]
		k = p
	}
	return k
}
