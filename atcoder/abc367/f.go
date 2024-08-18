package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N, Q := int32(io.NextInt()), int32(io.NextInt())
	A := make([]int, N)
	B := make([]int, N)
	for i := int32(0); i < N; i++ {
		A[i] = io.NextInt()
	}
	for i := int32(0); i < N; i++ {
		B[i] = io.NextInt()
	}

	seg := NewSegmentTree(N, func(i int32) E { return FromElement(int32(A[i])) })
	seg2 := NewSegmentTree(N, func(i int32) E { return FromElement(int32(B[i])) })

	query := func(l1, r1, l2, r2 int32) bool {
		e1 := seg.Query(l1, r1)
		e2 := seg2.Query(l2, r2)
		return e1.hash == e2.hash && e1.hash2 == e2.hash2 && e1.min_ == e2.min_ && e1.max_ == e2.max_
	}

	for i := int32(0); i < Q; i++ {
		l1, r1, l2, r2 := int32(io.NextInt()), int32(io.NextInt()), int32(io.NextInt()), int32(io.NextInt())
		l1--
		l2--
		if r1-l1 != r2-l2 {
			io.Println("No")
			continue
		}

		if query(l1, r1, l2, r2) {
			io.Println("Yes")
		} else {
			io.Println("No")
		}
	}
}

const INF int = 1e18
const INF32 int32 = 1e9 + 10
const MOD int = 999999929
const BASE int = 13331
const MOD2 int = 999999937
const BASE2 int = 131

func qpow(a, b int, mod int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

type E = struct {
	min_  int32
	hash  int
	max_  int32
	hash2 int
}

func FromElement(v int32) E { return E{v, qpow(BASE, int(v), MOD), v, qpow(BASE2, int(v), MOD2)} }
func (*SegmentTree) e() E   { return E{INF32, 0, 0, 0} }
func (*SegmentTree) op(a, b E) E {
	newMin := min32(a.min_, b.min_)
	newMax := max32(a.max_, b.max_)
	newHash := (a.hash + b.hash) % MOD
	newHash2 := (a.hash2 + b.hash2) % MOD2
	return E{newMin, newHash, newMax, newHash2}
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

type SegmentTree struct {
	n, size int32
	seg     []E
}

func NewSegmentTree(n int32, f func(int32) E) *SegmentTree {
	res := &SegmentTree{}
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}

func (st *SegmentTree) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

func (st *SegmentTree) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTree) QueryAll() E { return st.seg[1] }
func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}
