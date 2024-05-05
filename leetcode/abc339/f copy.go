// from math import log2
// import sys

// sys.setrecursionlimit(int(1e6))
// input = lambda: sys.stdin.readline().rstrip("\r\n")

// if __name__ == "__main__":
//     n = int(input())
//     nums = [log2(int(input())) for _ in range(n)]
//     mp = dict()
//     for num in nums:
//         mp[num] = mp.get(num, 0) + 1

//     res = 0
//     for i in range(n):
//         for j in range(n):
//             res += mp.get(nums[i] + nums[j], 0)
//     print(res)

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/big"
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

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	bigInts := make([]big.Int, n)
	for i := 0; i < n; i++ {
		bigInts[i].SetString(io.Text(), 10)
	}
	mp := make(map[string]int)
	for i := 0; i < n; i++ {
		mp[bigInts[i].String()]++
	}

	res := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var tmp big.Int
			tmp.Mul(&bigInts[i], &bigInts[j])
			res += mp[tmp.String()]
		}
	}
	io.Println(res)

}
