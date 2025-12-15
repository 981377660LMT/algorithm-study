package main

import (
	"fmt"
	"strings"
)

type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int, f func(i int) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITArray) QueryAll() int {
	return b.total
}

func (b *BITArray) MaxRight(check func(index, preSum int) bool) int {
	i := 0
	s := 0
	k := 1
	for 2*k <= b.n {
		k *= 2
	}
	for k > 0 {
		if i+k-1 < b.n {
			t := s + b.data[i+k-1]
			if check(i+k, t) {
				i += k
				s = t
			}
		}
		k >>= 1
	}
	return i
}

func (b *BITArray) Kth(k int) int {
	return b.MaxRight(func(index, preSum int) bool { return preSum <= k })
}

func (b *BITArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}

func minDeletions(s string, queries [][]int) []int {
	bytes := []byte(s)
	n := len(bytes)
	bit := NewBitArray(n)

	check := func(i int) int {
		if i < 0 || i >= n-1 {
			return 0
		}
		if bytes[i] == bytes[i+1] {
			return 1
		}
		return 0
	}

	for i := 0; i < n-1; i++ {
		if check(i) == 1 {
			bit.Add(i, 1)
		}
	}

	res := make([]int, 0, len(queries))
	for _, q := range queries {
		if q[0] == 1 {
			pos := q[1]
			if pos > 0 {
				bit.Add(pos-1, -check(pos-1))
			}
			if pos < n-1 {
				bit.Add(pos, -check(pos))
			}

			if bytes[pos] == 'A' {
				bytes[pos] = 'B'
			} else {
				bytes[pos] = 'A'
			}

			if pos > 0 {
				bit.Add(pos-1, check(pos-1))
			}
			if pos < n-1 {
				bit.Add(pos, check(pos))
			}
		} else {
			l, r := q[1], q[2]
			res = append(res, bit.QueryRange(l, r))
		}
	}

	return res
}
