// cf-Centroids-删边加边成为重心的点的个数
// https://codeforces.com/contest/708/problem/C
// 给定一颗树，你有一次将树改造的机会，改造的意思是删去一条边，再加入一条边，保证改造后还是一棵树。
// 请问有多少点可以通过改造，成为这颗树的重心？（如果以某个点为根，每个子树的大小都不大于 n/2，则称某个点为重心）
// n<=4e5

// 1.对于一个节点 x ，如果本身就是重心，就无需改造，否则就必然有且只有一个儿子的子树大于 n/2 。
// !2.如果可以从这个儿子的子树中拿出一个子树接到节点x 上，
//   使得接上的子树大小不超过 n/2 且剩下的子树大小也不超过 n/2 ，那么节点
//   x 就是可以改造为重心的。

// !以每个结点v作为根考虑
// v的子树中存在一个子树T的大小 > n/2 (重儿子)
// 如果T中存在一个子树S的大小满足size(T)-size(S)<=n/2且size(S)<=n/2
// !贪心，这样的S需要是不超过n//2的size最大的子树
// !那么就可以把子树S切出来，接到根v的子树中,使得v成为重心
// 总之,需要维护子树内的最大子树和次大子树的大小以及对应的编号

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"sort"
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
	R := NewRerootingSubTree(n)
	for i := 0; i < n-1; i++ {
		u, v := io.NextInt()-1, io.NextInt()-1
		R.AddEdge2(u, v, 1, 1)
	}

	R.ReRooting(
		func(root int) E { return E{} },
		func(dp1, dp2 E) E {
			return E{
				size: dp1.size + dp2.size,
				max_: max(dp1.max_, dp2.max_),
			}
		},
		func(dp E, e Edge) E {
			size, max_ := dp.size, dp.max_
			cand := 0
			if size+e <= n/2 { // 如果整个子树大小不超过n/2,那么就可以作为最大子树大小
				cand = size + e
			}
			return E{size + 1, max(max_, cand)}
		},
	)

	// 检查每个点为根时，是否存在这样的子树
	tree := R.G
	for root := 0; root < n; root++ {
		ok := true
		for _, e := range tree[root] {
			to := e.to
			subDp := R.SubTree(root, to)
			size, max_ := subDp.size, subDp.max_
			if size <= n/2 { // 已经是重心
				continue
			}
			if !(size-max_ <= n/2 && max_ <= n/2) {
				ok = false
				break
			}
		}

		if ok {
			io.Println(1)
		} else {
			io.Println(0)
		}
	}
}

// 整颗子树大小, 子树中不超过n/2的最大的子树的大小
type E = struct{ size, max_ int }

type Edge = int // Edge为边权

type ReRootingSubTree struct {
	G           [][]Node
	ld, rd      [][]E
	lp, rp      []int
	e           func(root int) E
	op          func(dp1, dp2 E) E
	composition func(dp E, e Edge) E
}

type Node struct {
	to, rev int
	edge    Edge
}

func NewRerootingSubTree(n int) *ReRootingSubTree {
	res := &ReRootingSubTree{
		G:  make([][]Node, n),
		ld: make([][]E, n),
		rd: make([][]E, n),
		lp: make([]int, n),
		rp: make([]int, n),
	}
	return res
}

func (rr *ReRootingSubTree) AddEdge(u, v int, e Edge) {
	rr.AddEdge2(u, v, e, e)
}

func (rr *ReRootingSubTree) AddEdge2(u, v int, e, revE Edge) {
	rr.G[u] = append(rr.G[u], Node{to: v, edge: e})
	rr.G[v] = append(rr.G[v], Node{to: u, edge: revE})
}

func (rr *ReRootingSubTree) ReRooting(
	e func(root int) E,
	op func(dp1, dp2 E) E,
	compositionEdge func(dp E, e Edge) E,
) []E {
	rr.e = e
	rr.op = op
	rr.composition = compositionEdge
	n := len(rr.G)
	for i := 0; i < n; i++ {
		sort.Slice(rr.G[i], func(j, k int) bool {
			return rr.G[i][j].to < rr.G[i][k].to
		})
		rr.ld[i] = make([]E, len(rr.G[i])+1)
		rr.rd[i] = make([]E, len(rr.G[i])+1)
		for j := range rr.ld[i] {
			rr.ld[i][j] = e(i)
			rr.rd[i][j] = e(i)
		}
		rr.lp[i] = 0
		rr.rp[i] = len(rr.G[i]) - 1
	}

	for i := 0; i < n; i++ {
		for j := range rr.G[i] {
			rr.G[i][j].rev = rr.search(rr.G[rr.G[i][j].to], i)
		}
	}

	res := make([]E, n)
	for i := 0; i < n; i++ {
		res[i] = rr.dfs(i, -1)
	}
	return res
}

// !root 作为根节点时, 子树 v 的 dp 值
func (rr *ReRootingSubTree) SubTree(root, v int) E {
	k := rr.search(rr.G[root], v)
	return rr.composition(rr.dfs(v, rr.G[root][k].rev), rr.G[root][k].edge)
}

func (rr *ReRootingSubTree) dfs(root, eid int) E {
	for rr.lp[root] != eid && rr.lp[root] < len(rr.G[root]) {
		e := rr.G[root][rr.lp[root]]
		rr.ld[root][rr.lp[root]+1] = rr.op(rr.ld[root][rr.lp[root]], rr.composition(rr.dfs(e.to, e.rev), e.edge))
		rr.lp[root]++
	}
	for rr.rp[root] != eid && rr.rp[root] >= 0 {
		e := rr.G[root][rr.rp[root]]
		rr.rd[root][rr.rp[root]] = rr.op(rr.rd[root][rr.rp[root]+1], rr.composition(rr.dfs(e.to, e.rev), e.edge))
		rr.rp[root]--
	}
	if eid < 0 {
		return rr.rd[root][0]
	}
	return rr.op(rr.ld[root][eid], rr.rd[root][eid+1])
}

func (rr *ReRootingSubTree) search(vs []Node, idx int) int {
	return sort.Search(len(vs), func(i int) bool { return vs[i].to >= idx })
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
