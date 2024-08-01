package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
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
	mo := NewMoFast(n, q)
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

	mo.Run(add, remove, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// MoAlgo.
type MoFast struct {
	n, q       int32
	bucketSize int32
	bit        int32
	mask       int32
	bucket     [][]int
	id         int32
}

// !max(n,q) <=2e6.
func NewMoFast(n, q int32) *MoFast {
	var bucketSize int32
	if q > 0 {
		bucketSize = int32(math.Ceil(math.Sqrt(3) * float64(n) / math.Sqrt(2*float64(q))))
	} else {
		bucketSize = n
	}
	if bucketSize == 0 {
		bucketSize = 1
	}
	bit := int32(bits.Len32(uint32(max32(n, q))))
	mask := int32((1 << bit) - 1)
	bucket := make([][]int, n/bucketSize+1)
	return &MoFast{n: n, q: q, bucketSize: bucketSize, bit: bit, mask: mask, bucket: bucket}
}

func (m *MoFast) AddQuery(start, end int32) {
	bid := start / m.bucketSize
	s64, e64, id64 := int(start), int(end), int(m.id)
	m.bucket[bid] = append(m.bucket[bid], (((e64<<m.bit)|s64)<<m.bit)|id64)
	m.id++
}

func (m *MoFast) Run(add, remove, query func(index int32)) {
	bucket, bit, mask := m.bucket, m.bit, int(m.mask)

	for i, b := range bucket {
		if i&1 == 1 {
			sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
		} else {
			sort.Ints(b)
		}
	}
	nl, nr := int32(0), int32(0)
	for _, b := range bucket {
		for _, rli := range b {
			r, l := int32(rli>>bit>>bit), int32(rli>>bit&mask)
			for nl > l {
				nl--
				add(nl)
			}
			for nr < r {
				add(nr)
				nr++
			}
			for nl < l {
				remove(nl)
				nl++
			}
			for nr > r {
				nr--
				remove(nr)
			}
			query(int32(rli & mask))
		}
	}
}

func (m *MoFast) RunVerbose(
	addLeft, addRight, removeLeft, removeRight func(i int32),
	query func(qi int32),
) {
	bucket, bit, mask := m.bucket, m.bit, int(m.mask)
	for i, b := range bucket {
		if i&1 == 1 {
			sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
		} else {
			sort.Ints(b)
		}
	}
	nl, nr := int32(0), int32(0)
	for _, b := range bucket {
		for _, rli := range b {
			r, l := int32(rli>>bit>>bit), int32(rli>>bit&mask)
			for nl > l {
				nl--
				addLeft(nl)
			}
			for nr < r {
				addRight(nr)
				nr++
			}
			for nl < l {
				removeLeft(nl)
				nl++
			}
			for nr > r {
				nr--
				removeRight(nr)
			}
			query(int32(rli & mask))
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

func max32(a, b int32) int32 {
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
