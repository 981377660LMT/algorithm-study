// https://www.luogu.com.cn/problem/CF1790F
// https://zhuanlan.zhihu.com/p/601326343

// CF1790F-Timofey and Black-White Tree-树上点对距离
// !给你一个 n 个节点的树，初始时root节点为黑色，其余节点为白色。
// 每一次染黑一个点，每次操作后问你任选两个黑点之后这两个黑点的最短距离。
// n<=2e5
//
// 记录每个点到最近黑点的距离.
// !在树上随机撒 sqrt(n) 个点，这些点之间的最短距离不超过 sqrt(n)。其在链上是成立的，而树比链要更紧。
// 因此，经过 sqrt(n) 次操作后，dist 只有 sqrt(n) 以下的取值有意义，如果一个点更新上的值超过这个值我们就可以选择不去更新。

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

	T := io.NextInt()
	for i := 0; i < T; i++ {
		n := io.NextInt()
		root := io.NextInt() - 1
		seq := make([]int, n-1)
		for i := 0; i < n-1; i++ {
			seq[i] = io.NextInt() - 1
		}
		tree := make([][]int, n)
		for i := 0; i < n-1; i++ {
			u := io.NextInt() - 1
			v := io.NextInt() - 1
			tree[u] = append(tree[u], v)
			tree[v] = append(tree[v], u)
		}

		curMin := n
		dist := make([]int, n)
		for i := 0; i < n; i++ {
			dist[i] = n
		}
		queue := make([]int, n)
		head, tail := 0, 0
		updateDist := func(start int) {
			head, tail = 0, 0
			dist[start] = 0
			queue[tail] = start
			tail++
			for head < tail {
				cur := queue[head]
				head++
				for _, next := range tree[cur] {
					cand := dist[cur] + 1
					if cand < dist[next] && cand < curMin { // !关键
						dist[next] = cand
						queue[tail] = next
						tail++
					}
				}
			}
		}
		updateDist(root)
		for _, v := range seq {
			curMin = min(curMin, dist[v])
			io.Println(curMin)
			updateDist(v)
		}
	}
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
