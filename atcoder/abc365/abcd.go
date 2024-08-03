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

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N := io.NextInt()
	S := io.Text()
	mp := map[string]int{"R": 0, "P": 1, "S": 2}
	nums := make([]int, N)
	for i, s := range S {
		nums[i] = mp[string(s)]
	}

	winBy := []int{1, 2, 0}
	winBys := make([]int, N)
	for i, num := range nums {
		winBys[i] = winBy[num]
	}

	memo := make([][]int, N+10)
	for i := 0; i < N; i++ {
		memo[i] = make([]int, 4)
		for j := 0; j < 4; j++ {
			memo[i][j] = -1
		}
	}
	var dfs func(int, int) int
	dfs = func(index, preType int) int {
		if index == N {
			return 0
		}
		if memo[index][preType+1] != -1 {
			return memo[index][preType+1]
		}
		res := 0
		for i := 0; i < 3; i++ {
			if i != preType && winBy[i] != nums[index] {
				if winBys[index] == i {
					res = max(res, dfs(index+1, i)+1)
				} else {
					res = max(res, dfs(index+1, i))
				}
			}
		}
		memo[index][preType+1] = res
		return res
	}
	res := dfs(0, -1)
	io.Println(res)

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
