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

// 東西に続く道があり、道の上には
// N 人の高橋くんがいます。 道は原点と呼ばれる点から東西に十分長く続いています。

// i 番目
// (1≤i≤N) の高橋くんは、はじめ原点から東に
// X
// i
// ​
//   メートル進んだところにいます。

// 高橋くんたちは道の上を東西に動くことができます。 具体的には、次の移動を好きなだけ行うことができます。

// 高橋くんを一人選ぶ。移動する先に他の高橋くんがいない場合、選んだ高橋くんを
// 1 メートル東に、もしくは西に移動させる。
// 高橋くんたちには合計
// Q 個の用事があり、
// i 個目
// (1≤i≤Q) の用事は次の形式で表されます。

// T
// i
// ​
//
//	番目の高橋くんが座標
//
// G
// i
// ​
//
//	に到着する。
//
// Q 個の用事を先頭から順にすべて完了するために必要な移動回数の最小値を求めてください。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

}
