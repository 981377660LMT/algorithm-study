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

// 2 次元平面上に
// N 個の石が置かれています。
// i 番目の石は座標
// (X
// i
// ​
//  ,Y
// i
// ​
//  ) にあります。石は全て第一象限(軸上含む)の格子点にあります。

// 石の置かれていない格子点
// (x,y) であって、上下左右のいずれかに
// 1 移動することを繰り返すことで、石の置かれている座標を通らずに
// (−1,−1) に到達することができないものの個数を求めてください。

// より正確には、石の置かれていない格子点
// (x,y) であって、以下の
// 4 条件を全て満たすような整数の組の有限列
// (x
// 0
// ​
//  ,y
// 0
// ​
//  ),…,(x
// k
// ​
//  ,y
// k
// ​
//  ) が存在しないものの個数を求めてください。

// (x
// 0
// ​
//
//	,y
//
// 0
// ​
//
//	)=(x,y)
//
// (x
// k
// ​
//
//	,y
//
// k
// ​
//
//	)=(−1,−1)
//
// 全ての
// 0≤i<k で
// ∣x
// i
// ​
//
//	−x
//
// i+1
// ​
//
//	∣+∣y
//
// i
// ​
//
//	−y
//
// i+1
// ​
//
//	∣=1
//
// どの
// 0≤i≤k でも、
// (x
// i
// ​
//
//	,y
//
// i
// ​
//
//	) に石はない
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

}
