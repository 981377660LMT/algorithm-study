package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

// from https://atcoder.jp/users/ccppjsrb
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

// # 1 から
// # N までの番号のついた
// # N 個の区間が与えられます。 区間
// # i は
// # [L
// # i
// # ​
// #  ,R
// # i
// # ​
// #  ] です。

// # 区間
// # [l
// # a
// # ​
// #  ,r
// # a
// # ​
// #  ] と区間
// # [l
// # b
// # ​
// #  ,r
// # b
// # ​
// #  ] は
// # (l
// # a
// # ​
// #  <l
// # b
// # ​
// #  <r
// # a
// # ​
// #  <r
// # b
// # ​
// #  ) または
// # (l
// # b
// # ​
// #  <l
// # a
// # ​
// #  <r
// # b
// # ​
// #  <r
// # a
// # ​
// #  ) を満たすとき、交差するといいます。

// # f(l,r) を
// # 1≤i≤N を満たし、区間
// # [l,r] と区間
// # i が交差する
// # i の個数と定義します。

// # 0≤l<r≤10
// # 9
// #   を満たす整数の組
// # (l,r) において、
// # f(l,r) の最大値を達成する
// # (l,r) の組のうち
// # l が最小のものを答えてください。そのような組が複数存在する場合はさらにそのうちで
// # r が最小のものを答えてください (
// # 0≤l<r より、 答えるべき
// # (l,r) の組は一意に定まります)。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N := io.NextInt()
	L := make([]int, N)
	R := make([]int, N)
	for i := 0; i < N; i++ {
		L[i] = io.NextInt()
		R[i] = io.NextInt()
	}

}
