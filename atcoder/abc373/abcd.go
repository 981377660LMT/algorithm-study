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

	N, M := int32(io.NextInt()), int32(io.NextInt())
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	uf := NewPotentializedUnionFind(N, e, op, inv)

	for i := int32(0); i < M; i++ {
		u, v, w := int32(io.NextInt()), int32(io.NextInt()), io.NextInt()
		u, v = u-1, v-1
		uf.Union(v, u, w)
	}

	groups := make(map[int32][]int32)
	for i := int32(0); i < N; i++ {
		root, _ := uf.Find(i)
		groups[root] = append(groups[root], i)
	}
	res := make([]int, N)
	for _, group := range groups {
		for _, v := range group {
			diff, _ := uf.Diff(v, group[0])
			res[v] = diff
		}
	}

	for i := int32(0); i < N; i++ {
		io.Print(res[i], " ")
	}
}

type PotentializedUnionFind[E any] struct {
	n, Part int32
	parents []int32
	sizes   []int32
	values  []E
	e       func() E
	op      func(E, E) E
	inv     func(E) E
}

func NewPotentializedUnionFind[E any](n int32, e func() E, op func(E, E) E, inv func(E) E) *PotentializedUnionFind[E] {
	values, parents, sizes := make([]E, n), make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		parents[i] = i
		sizes[i] = 1
		values[i] = e()
	}
	return &PotentializedUnionFind[E]{n: n, Part: n, parents: parents, sizes: sizes, values: values, e: e, op: op, inv: inv}
}

func (uf *PotentializedUnionFind[E]) Union(a, b int32, x E) bool {
	v1, x1 := uf.Find(b)
	v2, x2 := uf.Find(a)
	if v1 == v2 {
		return false
	}
	if uf.sizes[v1] < uf.sizes[v2] {
		v1, v2 = v2, v1
		x1, x2 = x2, x1
		x = uf.inv(x)
	}
	x = uf.op(x1, x)
	x = uf.op(x, uf.inv(x2))
	uf.values[v2] = x
	uf.parents[v2] = v1
	uf.sizes[v1] += uf.sizes[v2]
	uf.Part--
	return true
}

func (uf *PotentializedUnionFind[E]) Find(v int32) (root int32, diff E) {
	diff = uf.e()
	vs, ps := uf.values, uf.parents
	for v != ps[v] {
		diff = uf.op(vs[v], diff)
		diff = uf.op(vs[ps[v]], diff)
		vs[v] = uf.op(vs[ps[v]], vs[v])
		ps[v] = ps[ps[v]]
		v = ps[v]
	}
	root = v
	return
}

func (uf *PotentializedUnionFind[E]) Diff(a, b int32) (E, bool) {
	ru, xu := uf.Find(b)
	rv, xv := uf.Find(a)
	if ru != rv {
		return uf.e(), false
	}
	return uf.op(uf.inv(xu), xv), true
}

func (uf *PotentializedUnionFind[E]) Union2(a, b int32, x E, beforeUnion func(big, small int32)) bool {
	v1, x1 := uf.Find(b)
	v2, x2 := uf.Find(a)
	if v1 == v2 {
		return false
	}
	if uf.sizes[v1] < uf.sizes[v2] {
		v1, v2 = v2, v1
		x1, x2 = x2, x1
		x = uf.inv(x)
	}
	if beforeUnion != nil {
		beforeUnion(v1, v2)
	}
	x = uf.op(x1, x)
	x = uf.op(x, uf.inv(x2))
	uf.values[v2] = x
	uf.parents[v2] = v1
	uf.sizes[v1] += uf.sizes[v2]
	uf.Part--
	return true
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
