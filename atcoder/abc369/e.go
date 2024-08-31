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

// N 個の島と、
// 2 つの島の間を双方向に結ぶ
// M 本の橋があり、 それぞれ島
// 1, 島
// 2,
// …, 島
// N および 橋
// 1, 橋
// 2,
// …, 橋
// M と番号づけられています。
// 橋
// i は島
// U
// i
// ​
//   と島
// V
// i
// ​
//   を相互に結んでおり、どちらの方向に移動するにも
// T
// i
// ​
//   だけ時間がかかります。
// ここで、橋の両端が同一の島であるような橋は存在しませんが、ある
// 2 つの島の間が
// 2 本以上の橋で直接繋がれている可能性はあります。
// また、どの
// 2 つの島の間もいくつかの橋をわたって移動することができます。

// Q 個の問題が与えられるので、各問題に対する答えを求めてください。
// i 番目の問題は次のようなものです。

// 相異なる
// K
// i
// ​
//   本の橋、橋
// B
// i,1
// ​
//  , 橋
// B
// i,2
// ​
//  ,
// …, 橋
// B
// i,K
// i
// ​

// ​
//
//	が与えられます。
//
// これらの橋をすべて
// 1 回以上わたり、島
// 1 から島
// N まで移動するために必要な時間の最小値を求めてください。
// ただし、島
// 1 から島
// N までの移動において、橋をわたって島の間を移動するのにかかる時間以外は無視できるものとします。
// また、与えられた橋はどの順で、またどの向きにわたってもかまいません。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

}
