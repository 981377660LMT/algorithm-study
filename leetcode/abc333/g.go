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

// sys.setrecursionlimit(int(1e6))
// input = lambda: sys.stdin.readline().rstrip("\r\n")
// MOD = 998244353
// INF = int(4e18)
// # N 人の人が一列に並んでおり、人
// # i は先頭から
// # i 番目に並んでいます。

// # 以下の操作を、列に並んでいる人が
// # 1 人になるまで繰り返します。

// # 先頭に並んでいる人を
// # 2
// # 1
// # ​
// #   の確率で列から取り除き、そうでない場合は列の末尾に移す。
// # 人
// # i=1,2,…,N それぞれについて、人
// # i が最後まで列に並んでいる
// # 1 人になる確率を
// # mod 998244353 で求めて下さい。(取り除くかどうかの選択はすべてランダムかつ独立です。)

const MOD int = 998244353
const INV2 int = (MOD + 1) / 2

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()

	memo := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		memo[i] = make([]int, n+1)
		for j := 0; j <= n; j++ {
			memo[i][j] = -1
		}
	}

	var dfs func(int, int) int
	dfs = func(round int, shift int) int {
		return 1
	}

	for i := 0; i < n; i++ {
		io.Print(dfs(0, i), " ")
	}

}
