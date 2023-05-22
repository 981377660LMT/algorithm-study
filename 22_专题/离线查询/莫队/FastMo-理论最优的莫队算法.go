package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	// G - Range Pairing Query
	// https://atcoder.jp/contests/abc242/tasks/abc242_g
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const N int = 1e5 + 10
	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)
	mo := NewFastMo(n, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		mo.AddQuery(l, r)
	}

	pair := 0
	counter := [N + 1]int{}
	res := make([]int, q)
	add := func(i int) {
		v := nums[i]
		pair -= counter[v] / 2
		counter[v]++
		pair += counter[v] / 2
	}
	remove := func(i int) {
		v := nums[i]
		pair -= counter[v] / 2
		counter[v]--
		pair += counter[v] / 2
	}
	query := func(qid int) { res[qid] = pair }

	mo.Run(add, add, remove, remove, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type FastMo struct {
	n, q, width        int
	left, right, order []int
	isBuild            bool
}

func NewFastMo(n, q int) *FastMo {
	width := max(1, n/max(1, int(math.Sqrt(float64(q)/2))))
	order := make([]int, q)
	for i := 0; i < q; i++ {
		order[i] = i
	}
	return &FastMo{n: n, q: q, width: width, order: order}
}

// [start, end).
func (mo *FastMo) AddQuery(start, end int) {
	mo.left = append(mo.left, start)
	mo.right = append(mo.right, end)
}

func (mo *FastMo) Build() {
	mo.sort()
	mo.climb(3, 5)
	mo.isBuild = true
}

func (mo *FastMo) Run(addLeft, addRight, deleteLeft, deleteRight, query func(i int)) {
	if !mo.isBuild {
		mo.Build()
	}
	nl, nr := 0, 0
	for _, qi := range mo.order {
		for nl > mo.left[qi] {
			addLeft(nl - 1)
			nl--
		}
		for nr < mo.right[qi] {
			addRight(nr)
			nr++
		}
		for nl < mo.left[qi] {
			deleteLeft(nl)
			nl++
		}
		for nr > mo.right[qi] {
			deleteRight(nr - 1)
			nr--
		}
		query(qi)
	}
}

func (mo *FastMo) sort() {
	n, q := mo.n, mo.q
	cnt := make([]int, n+1)
	buf := make([]int, q)
	for i := 0; i < q; i++ {
		cnt[mo.right[i]]++
	}
	for i := 1; i < len(cnt); i++ {
		cnt[i] += cnt[i-1]
	}
	for i := 0; i < q; i++ {
		cnt[mo.right[i]]--
		buf[cnt[mo.right[i]]] = i
	}
	b := make([]int, q)
	for i := 0; i < q; i++ {
		b[i] = mo.left[i] / mo.width
	}
	newCnt := make([]int, n/mo.width+1)
	for i := 0; i < q; i++ {
		newCnt[b[i]]++
	}
	for i := 1; i < len(newCnt); i++ {
		newCnt[i] += newCnt[i-1]
	}
	for i := 0; i < q; i++ {
		newCnt[b[buf[i]]]--
		mo.order[newCnt[b[buf[i]]]] = buf[i]
	}
	for i, j := 0, 0; i < q; i = j {
		bi := b[mo.order[i]]
		j = i + 1
		for j != q && bi == b[mo.order[j]] {
			j++
		}
		if bi&1 == 0 {
			for p1, p2 := i, j-1; p1 < p2; p1, p2 = p1+1, p2-1 {
				mo.order[p1], mo.order[p2] = mo.order[p2], mo.order[p1]
			}
		}
	}
}

func (mo *FastMo) dist(i, j int) int {
	return abs(mo.left[i]-mo.left[j]) + abs(mo.right[i]-mo.right[j])
}

func (mo *FastMo) climb(iter, interval int) {
	q := mo.q
	d := make([]int, q-1)
	for i := 0; i < q-1; i++ {
		d[i] = mo.dist(i, i+1)
	}
	for iter > 0 {
		iter--
		for i := 1; i < q; i++ {
			pre1 := d[i-1]
			js := i + 1
			je := min(i+interval, q-1)
			for j := je - 1; j >= js; j-- {
				pre2 := d[j]
				now1 := mo.dist(i-1, j)
				now2 := mo.dist(i, j+1)
				if now1+now2 < pre1+pre2 {
					for p1, p2 := i, j; p1 < p2; p1, p2 = p1+1, p2-1 {
						mo.order[p1], mo.order[p2] = mo.order[p2], mo.order[p1]
					}
					for p1, p2 := i-1, j; p1 < p2; p1, p2 = p1+1, p2-1 {
						d[p1], d[p2] = d[p2], d[p1]
					}
					d[i-1] = now1
					d[j] = now2
				}
			}
		}
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
