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

// 長さ
// N の数列
// P=(P
// 1
// ​
//  ,P
// 2
// ​
//  ,…,P
// N
// ​
//  ) が与えられます。高橋君と青木君が、数列
// P を使ってゲームを行います。

// まず、高橋君が、
// 1,2,…,N から
// K 個の相異なる整数
// x
// 1
// ​
//  ,x
// 2
// ​
//  ,…,x
// K
// ​
//   を選びます。

// 次に、青木君が、
// 1,2,…,N から
// 1 つの整数
// y を
// P
// y
// ​
//   に比例する確率で選びます。すなわち、整数
// y が選ばれる確率は
// ∑
// y
// ′
//  =1
// N
// ​
//  P
// y
// ′

// ​

// P
// y
// ​

// ​
//   です。そして、青木君が
// i=1,2,…,K
// min
// ​
//  ∣x
// i
// ​
//  −y∣ のスコアを得ます。

// 高橋君は、青木君が得るスコアの期待値を最小化したいです。高橋君が適切に
// x
// 1
// ​
//  ,x
// 2
// ​
//  ,…,x
// K
// ​
//   を選んだときに、青木君が得るスコアの期待値の最小値を
// ∑
// y
// ′
//  =1
// N
// ​
//  P
// y
// ′

// ​
//
//	倍した値を求めてください。なお、出力すべき値は整数になることが証明できます。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

}
