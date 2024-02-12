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

// 正整数
// N,M,K と、長さ
// N の正整数列
// (C
// 1
// ​
//  ,C
// 2
// ​
//  ,…,C
// N
// ​
//  ) が与えられるので、
// r=0,1,2,…,N−1 の場合それぞれについて、下記の問題の答えを出力してください。

// 色がついた
// N 個のボールからなる列があり、
// i=1,2,…,N について、列の先頭から
// i 番目にあるボールの色は
// C
// i
// ​
//   です。 また、
// 1 から
// M の番号がつけられた
// M 個の空の箱があります。

// 下記の手順を行った後に箱に入っているボールの総数を求めてください。

// まず、下記の操作を
// r 回行う。

// 列の先頭のボール
// 1 個を列の最後尾に移動する。
// その後、列にボールが
// 1 個以上残っている限り、下記の操作を繰り返す。

// 列の先頭のボールと同じ色のボールが既に
// 1 個以上
// K 個未満入っている箱が存在する場合、その箱に列の先頭のボールを入れる。
// そのような箱が存在しない場合、
// 空の箱が存在するなら、そのうち番号が最小のものに、列の先頭のボールを入れる。
// 空の箱が存在しない場合、列の先頭のボールをどの箱にも入れず、食べる。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N, M, K := io.NextInt(), io.NextInt(), io.NextInt()

}
