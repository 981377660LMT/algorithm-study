package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
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

type G = int
type S = int

// 求离散对数 a^n = b (mod m) 的最小非负整数解 n.
//
//	如果不存在则返回-1.
func ModLog(a, b, mod int) int {
	a %= mod
	b %= mod
	return discreteLogActed(
		func() G { return 1 % mod },
		func(g1, g2 G) G { return (g1 * g2) % mod },
		func(s, g G) S { return s * g % mod },
		a,
		1,
		b,
		0,
		mod,
	)
}

func discreteLogActed(
	e func() G,
	op func(g1, g2 G) G,

	act func(s S, g G) S,

	x G,
	s S,
	t S,
	lower int,
	higher int,
) int {
	if lower >= higher {
		return -1
	}
	UNIT := e()
	set := make(map[S]struct{})
	xpow := func(n int) G {
		p := x
		res := UNIT
		for n > 0 {
			if n&1 == 1 {
				res = op(res, p)
			}
			p = op(p, p)
			n /= 2
		}
		return res
	}

	ht := t
	s = act(s, xpow(lower))
	LIM := higher - lower
	K := int(math.Sqrt(float64(LIM))) + 1
	for k := 0; k < K; k++ {
		t = act(t, x)
		set[t] = struct{}{}
	}

	y := xpow(K)
	failed := false
	for k := 0; k <= K; k++ {
		s1 := act(s, y)
		if _, ok := set[s1]; ok {
			for i := 0; i < K; i++ {
				if s == ht {
					cand := k*K + i + lower
					if cand >= higher {
						return -1
					}
					return cand
				}
				s = act(s, x)
			}
			if failed {
				return -1
			}
			failed = true
		}
		s = s1
	}
	return -1
}

// N 個の整数
// A
// 1
// ​
//  ,…,A
// N
// ​
//   と素数
// P が与えられます。 次の条件をともに満たす整数の組
// (i,j) の個数を求めてください。

// 1≤i,j≤N
// ある正整数
// k が存在し、
// A
// i
// k
// ​
//
//	≡A
//
// j
// ​
//
//	modP
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N, P := io.NextInt(), io.NextInt()
	A := make([]int, N)
	for i := range A {
		A[i] = io.NextInt()
	}

	// low := make([]int, N)
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			io.Println(ModLog(A[i], A[j], P))
		}
	}

}
