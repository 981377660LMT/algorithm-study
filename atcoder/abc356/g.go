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

// 高橋くんは
// N 種類の泳ぎ方ができます。
// 高橋くんが
// i 種類目の泳ぎ方で泳ぐと、
// 1 秒あたり体力を
// A
// i
// ​
//   消費して
// B
// i
// ​
//   [m] 進みます。

// Q 個のクエリに答えてください。そのうち
// i 個目は次の通りです。

// 消費する体力の合計を
// C
// i
// ​
//   以下にして
// D
// i
// ​
//   [m] 進むことができるか判定し、進める場合は必要な最小の秒数を求めよ。
// ただし、高橋くんは泳ぎ方を自由に組み合わせることができ、泳ぎ方を変える時間は無視できます。
// 具体的には、次の手順で泳ぐことができます。

// 正整数
// m 、全ての要素が正である長さ
// m の実数列
// t=(t
// 1
// ​
//
//	,t
//
// 2
// ​
//
//	,…,t
//
// m
// ​
//
//	) 、全ての要素が
//
// 1 以上
// N 以下の長さ
// m の整数列
// x=(x
// 1
// ​
//
//	,x
//
// 2
// ​
//
//	,…,x
//
// m
// ​
//
//	) を選択する。
//
// その後、
// i=1,2,…,m の順で、
// x
// i
// ​
//
//	種類目の泳ぎ方で
//
// t
// i
// ​
//
//	秒間泳ぐ。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	A, B := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		A[i], B[i] = io.NextInt(), io.NextInt()
	}
	q := io.NextInt()
	C, D := make([]int, q), make([]int, q)
	for i := 0; i < q; i++ {
		C[i], D[i] = io.NextInt(), io.NextInt()
	}

}
