package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	abc242_g()
}

func abc242_g() {
	// G - Range Pairing Query
	// https://atcoder.jp/contests/abc242/tasks/abc242_g

	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	// 读一个整数，支持负数
	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}
	_ = NextInt

	const N int32 = 1e5 + 10
	n := int32(NextInt())
	nums := make([]int32, n)
	for i := range nums {
		nums[i] = int32(NextInt())
	}

	q := int32(NextInt())
	mo := NewMoV2(n, q)
	for i := int32(0); i < q; i++ {
		l, r := int32(NextInt()), int32(NextInt())
		l--
		mo.AddQuery(l, r)
	}

	pair := 0
	counter := [N + 1]int{}
	res := make([]int, q)
	add := func(i int32) {
		v := nums[i]
		pair -= counter[v] >> 1
		counter[v]++
		pair += counter[v] >> 1
	}
	remove := func(i int32) {
		v := nums[i]
		pair -= counter[v] >> 1
		counter[v]--
		pair += counter[v] >> 1
	}
	query := func(qid int32) { res[qid] = pair }

	mo.Run(add, add, remove, remove, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type MoV2 struct {
	lefts, rights []int32
}

func NewMoV2(n, q int32) *MoV2 {
	res := &MoV2{
		lefts:  make([]int32, 0, q),
		rights: make([]int32, 0, q),
	}
	return res
}

func (m *MoV2) AddQuery(start, end int32) {
	m.lefts = append(m.lefts, start)
	m.rights = append(m.rights, end)
}

func (m *MoV2) Run(
	addL, addR func(i int32),
	removeL, removeR func(i int32),
	query func(qid int32),
) {
	order := getMoOrder(m.lefts, m.rights)
	l, r := int32(0), int32(0)
	for _, idx := range order {
		left, right := m.lefts[idx], m.rights[idx]
		for l > left {
			l--
			addL(l)
		}
		for r < right {
			addR(r)
			r++
		}
		for l < left {
			removeL(l)
			l++
		}
		for r > right {
			r--
			removeR(r)
		}
		query(idx)
	}
}

func getMoOrder(lefts, rights []int32) []int32 {
	n := int32(1)
	for i := 0; i < len(lefts); i++ {
		n = max32(n, lefts[i])
		n = max32(n, rights[i])
	}
	q := len(lefts)
	if q == 0 {
		return []int32{}
	}
	bs := int32(math.Sqrt(3) * float64(n) / math.Sqrt(2*float64(q)))
	bs = max32(bs, 1)
	order := make([]int32, q)
	for i := 0; i < q; i++ {
		order[i] = int32(i)
	}
	belong := make([]int32, q)
	for i := 0; i < q; i++ {
		belong[i] = lefts[i] / bs
	}
	sort.Slice(order, func(a, b int) bool {
		oa, ob := order[a], order[b]
		bida, bidb := belong[oa], belong[ob]
		if bida != bidb {
			return bida < bidb
		}
		if bida&1 == 1 {
			return rights[oa] > rights[ob]
		}
		return rights[oa] < rights[ob]
	})

	cost := func(a, b int32) int32 {
		oa, ob := order[a], order[b]
		return abs32(lefts[oa]-lefts[ob]) + abs32(rights[oa]-rights[ob])
	}
	for k := int32(0); k < int32(q-5); k++ {
		if cost(k, k+2)+cost(k+1, k+3) < cost(k, k+1)+cost(k+2, k+3) {
			order[k+1], order[k+2] = order[k+2], order[k+1]
		}
		if cost(k, k+3)+cost(k+1, k+4) < cost(k, k+1)+cost(k+3, k+4) {
			order[k+1], order[k+3] = order[k+3], order[k+1]
		}
	}

	return order
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
