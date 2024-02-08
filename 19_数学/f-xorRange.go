// 枚举异或范围(enumerateXorRange)
// lo <= (x ^ num) < hi となる x の区間 [L, R) の列挙
// 枚举所有的区间[L,R)使得在这个区间里的x满足
// !floor <= x^num < higher
// !注意每段[L,R)的长度都是2的幂次,这样的区间最多有O(log(higher-floor))个

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"sort"
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

// E. Beautiful Subarrays 子数组异或>=k的个数
// https://www.luogu.com.cn/problem/CF665E
// 给定长度为n的数组nums和整数k.
// !求有多少个nums的子数组满足子数组的异或>=k.
// k<=1e9, n<=1e5, nums[i]<=1e9.
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	k := io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}

	preXor := make([]int, n+1)
	for i := 0; i < n; i++ {
		preXor[i+1] = preXor[i] ^ nums[i]
	}
	sort.Ints(preXor)

	count := func(start, end int) int {
		return sort.SearchInts(preXor, end) - sort.SearchInts(preXor, start)
	}

	res := 0
	for _, v := range preXor {
		EnumerateXorRange(v, k, 1<<32, func(start, end int) {
			res += count(start, end)
		})
	}
	io.Println(res / 2)
}

func EnumerateXorRange(num int, floor, higher int, f func(start, end int)) {
	bit := 0
	for floor < higher {
		if floor&1 != 0 {
			f((floor^num)<<bit, ((floor^num)+1)<<bit)
			floor++
		}
		if higher&1 != 0 {
			higher--
			f((higher^num)<<bit, ((higher^num)+1)<<bit)
		}
		floor >>= 1
		higher >>= 1
		num >>= 1
		bit++
	}
}
