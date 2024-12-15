package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Fenw struct {
	n  int
	fw []int
}

func NewFenw(n int) *Fenw {
	return &Fenw{n: n, fw: make([]int, n+1)}
}

func (f *Fenw) Add(pos, val int) {
	for pos <= f.n {
		f.fw[pos] += val
		pos += pos & -pos
	}
}

func (f *Fenw) Sum(pos int) int {
	res := 0
	for pos > 0 {
		res += f.fw[pos]
		pos -= pos & -pos
	}
	return res
}

func (f *Fenw) RangeSum(l, r int) int {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

func main() {
	input := bufio.NewReader(os.Stdin)

	var N int
	fmt.Fscan(input, &N)
	A := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(input, &A[i])
	}
	B := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(input, &B[i])
	}

	var K int
	fmt.Fscan(input, &K)
	type Query struct {
		X, Y, id int
	}
	queries := make([]Query, K)
	for i := 0; i < K; i++ {
		fmt.Fscan(input, &queries[i].X, &queries[i].Y)
		queries[i].id = i
	}

	sort.Slice(queries, func(i, j int) bool {
		if queries[i].X == queries[j].X {
			return queries[i].Y < queries[j].Y
		}
		return queries[i].X < queries[j].X
	})

	prefixA := make([]int, N+1)
	for i := 1; i <= N; i++ {
		prefixA[i] = prefixA[i-1] + A[i-1]
	}

	Avals := append([]int(nil), A...)
	sort.Ints(Avals)

	fenwCount := NewFenw(N)
	fenwSum := NewFenw(N)

	rankVal := func(x int) int {
		l, r := 0, N
		for l < r {
			mid := (l + r) >> 1
			if Avals[mid] <= x {
				l = mid + 1
			} else {
				r = mid
			}
		}
		return r
	}

	lastResultForX := make([]int, N+1)
	res := make([]int, K)

	xPtr := 0

	currYVal := 0

	for _, q := range queries {
		for xPtr < q.X {
			xPtr++
			val := A[xPtr-1]
			r := rankVal(val)
			fenwCount.Add(r, 1)
			fenwSum.Add(r, val)
		}

		currYVal = B[q.Y-1]
		r := rankVal(currYVal)
		count_left := fenwCount.Sum(r)
		sum_left := fenwSum.Sum(r)
		X := q.X
		fXY := prefixA[X] - 2*sum_left + (2*count_left-X)*currYVal

		PXY := lastResultForX[X] + fXY
		res[q.id] = PXY
		lastResultForX[X] = PXY
	}

	wr := bufio.NewWriter(os.Stdout)
	for i := 0; i < K; i++ {
		fmt.Fprintln(wr, res[i])
	}
	wr.Flush()
}
