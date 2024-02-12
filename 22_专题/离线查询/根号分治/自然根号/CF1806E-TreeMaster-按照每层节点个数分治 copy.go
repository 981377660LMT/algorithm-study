// CF1806E-TreeMaster
// https://www.luogu.com.cn/problem/CF1806E
// 给定一颗以 0 为根的树，每个节点有一个权值，
// 每次询问给出两个深度相同的点，求他们到根的这条路径上所有深度相同的点权值的积的和。
//
// 记忆化dfs
// !节点数 x <=sqrt(n) 的层，最坏情况被调用x^2次，最多有n/x个这样的层，这部分复杂度为O(n/x*x^2)=O(n*x),这部分复杂度为O(n*sqrt(n)).
// !节点数 x > sqrt(n) 时层，不走缓存，直接继续调用dfs，这样的层小于sqrt(n)个，每一层最多被调用q次，这部分复杂度为O(q*sqrt(n)).

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

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, q := io.NextInt(), io.NextInt()
	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = io.NextInt()
	}
	parents := make([]int, n)
	parents[0] = -1
	for i := 1; i < n; i++ {
		parents[i] = io.NextInt() - 1
	}

	depth := make([]int, n)
	levelCount := make([]int, n) // 每一层的节点数
	levelId := make([]int, n)    // 每个节点在每一层的id
	resToRoot := make([]int, n)  // 每个节点到根的路径上的节点权值的积的和
	resToRoot[0] = values[0] * values[0]
	for i := 1; i < n; i++ {
		depth[i] = depth[parents[i]] + 1
		levelId[i] = levelCount[depth[i]]
		levelCount[depth[i]]++
		resToRoot[i] = resToRoot[parents[i]] + values[i]*values[i]
	}

	const SQRT = 100
	memo := make(map[int]int)

	var dfs func(x, y int) int
	dfs = func(x, y int) int {
		if x == y {
			return resToRoot[x]
		}
		if levelCount[depth[x]] < SQRT {
			xid := levelId[x]
			hash := xid*n + y
			if v, ok := memo[hash]; ok {
				return v
			}
			res := dfs(parents[x], parents[y]) + values[x]*values[y]
			memo[hash] = res
			return res
		} else {
			return dfs(parents[x], parents[y]) + values[x]*values[y]
		}
	}

	for i := 0; i < q; i++ {
		x, y := io.NextInt()-1, io.NextInt()-1
		io.Println(dfs(x, y))
	}
}
