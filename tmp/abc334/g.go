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

const MOD int = 998244353

// H 行
// W 列のグリッドがあり、グリッドの各マスは赤色あるいは緑色に塗られています。

// グリッドの上から
// i 行目、左から
// j 列目のマスをマス
// (i,j) と表記します。

// マス
// (i,j) の色は文字
// S
// i,j
// ​
//   で表され、
// S
// i,j
// ​
//  = . のときマス
// (i,j) は赤色、
// S
// i,j
// ​
//  = # のときマス
// (i,j) は緑色に塗られています。

// グリッドにおいて、緑色に塗られたマスを頂点集合、隣り合った
// 2 つの緑色のマスを結ぶ辺全体を辺集合としたグラフにおける連結成分の個数を 緑の連結成分数 と呼びます。ただし、
// 2 つのマス
// (x,y) と
// (x
// ′
//  ,y
// ′
//  ) が隣り合っているとは、
// ∣x−x
// ′
//  ∣+∣y−y
// ′
//  ∣=1 であることを指します。

// 緑色に塗られたマスを一様ランダムに
// 1 つ選び、赤色に塗り替えたとき、塗り替え後のグリッドの緑の連結成分数の期待値を
// mod 998244353 で出力してください。

const GREEN byte = '#'
const RED byte = '.'

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	ROW, COL := io.NextInt(), io.NextInt()
	grid := make([][]byte, ROW)
	for i := range grid {
		grid[i] = []byte(io.Text())
	}

	uf := NewUnionFindArray(ROW * COL)
	greenCount := 0
	adjList := make([][]Edge, ROW*COL)
	deg := make([]int, ROW*COL)
	eid := 0
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == GREEN {
				greenCount++
				if i > 0 && grid[i-1][j] == GREEN {
					adjList[i*COL+j] = append(adjList[i*COL+j], Edge{i*COL + j, (i-1)*COL + j, eid})
					eid++
					adjList[(i-1)*COL+j] = append(adjList[(i-1)*COL+j], Edge{(i-1)*COL + j, i*COL + j, eid})
					eid++
					uf.Union(i*COL+j, (i-1)*COL+j)
					deg[i*COL+j]++
					deg[(i-1)*COL+j]++
				}
				if j > 0 && grid[i][j-1] == GREEN {
					adjList[i*COL+j] = append(adjList[i*COL+j], Edge{i*COL + j, i*COL + j - 1, eid})
					eid++
					adjList[i*COL+j-1] = append(adjList[i*COL+j-1], Edge{i*COL + j - 1, i*COL + j, eid})
					eid++
					uf.Union(i*COL+j, i*COL+j-1)
					deg[i*COL+j]++
					deg[i*COL+j-1]++
				}
			}
		}
	}
	redCount := ROW*COL - greenCount

	lowLink := NewLowLink(adjList)
	lowLink.Build()

	isCut := make([]bool, ROW*COL)
	for _, v := range lowLink.Articulation {
		isCut[v] = true
	}

	res := 0
	curPart := uf.Part - redCount

	countNeighbor := func(r, c int) int {
		cnt := 0
		if r > 0 && grid[r-1][c] == GREEN {
			cnt++
		}
		if r < ROW-1 && grid[r+1][c] == GREEN {
			cnt++
		}
		if c > 0 && grid[r][c-1] == GREEN {
			cnt++
		}
		if c < COL-1 && grid[r][c+1] == GREEN {
			cnt++
		}
		return cnt
	}

	countNeighborWithDeg1 := func(r, c int) int {
		cnt := 0
		if r > 0 && grid[r-1][c] == GREEN && deg[(r-1)*COL+c] == 1 {
			cnt++
		}
		if r < ROW-1 && grid[r+1][c] == GREEN && deg[(r+1)*COL+c] == 1 {
			cnt++
		}
		if c > 0 && grid[r][c-1] == GREEN && deg[r*COL+c-1] == 1 {
			cnt++
		}
		if c < COL-1 && grid[r][c+1] == GREEN && deg[r*COL+c+1] == 1 {
			cnt++
		}
		return cnt
	}

	// 孤立点、割点、其他情形
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == GREEN {
				cur := 0
				neigh := countNeighbor(i, j)
				if neigh == 0 {
					cur = curPart - 1
				} else if isCut[i*COL+j] {
					deg1 := countNeighborWithDeg1(i, j)
					cur = curPart + deg1
					if neigh == deg1 {
						cur--
					}
				} else {
					cur = curPart
				}
				res += cur
				res %= MOD

			}
		}
	}

	inv := Pow(greenCount, MOD-2, MOD)
	res = res * inv % MOD
	io.Println(res)

}

func Pow(base, exp, mod int) int {
	base %= mod
	res := 1 % mod
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}

type BiConnectedComponents struct {
	BCC     [][]Edge // 每个边双连通分量中的边
	g       [][]Edge
	lowLink *LowLink
	used    []bool
	tmp     []Edge
}

func NewBiConnectedComponents(g [][]Edge) *BiConnectedComponents {
	return &BiConnectedComponents{
		g:       g,
		lowLink: NewLowLink(g),
	}
}

func (bcc *BiConnectedComponents) Build() {
	bcc.lowLink.Build()
	bcc.used = make([]bool, len(bcc.g))
	for i := 0; i < len(bcc.used); i++ {
		if !bcc.used[i] {
			bcc.dfs(i, -1)
		}
	}
}

func (bcc *BiConnectedComponents) dfs(idx, par int) {
	bcc.used[idx] = true
	beet := false
	for _, next := range bcc.g[idx] {
		if next.to == par {
			b := beet
			beet = true
			if !b {
				continue
			}
		}

		if !bcc.used[next.to] || bcc.lowLink.ord[next.to] < bcc.lowLink.ord[idx] {
			bcc.tmp = append(bcc.tmp, next)
		}

		if !bcc.used[next.to] {
			bcc.dfs(next.to, idx)
			if bcc.lowLink.low[next.to] >= bcc.lowLink.ord[idx] {
				bcc.BCC = append(bcc.BCC, []Edge{})
				for {
					e := bcc.tmp[len(bcc.tmp)-1]
					bcc.BCC[len(bcc.BCC)-1] = append(bcc.BCC[len(bcc.BCC)-1], e)
					bcc.tmp = bcc.tmp[:len(bcc.tmp)-1]
					if e.index == next.index {
						break
					}
				}
			}
		}
	}
}

type Edge = struct{ from, to, index int }
type LowLink struct {
	Articulation []int  // 関節点
	Bridge       []Edge // 橋
	g            [][]Edge
	ord, low     []int
	used         []bool
}

func NewLowLink(g [][]Edge) *LowLink {
	return &LowLink{g: g}
}

func (ll *LowLink) Build() {
	ll.used = make([]bool, len(ll.g))
	ll.ord = make([]int, len(ll.g))
	ll.low = make([]int, len(ll.g))
	k := 0
	for i := 0; i < len(ll.g); i++ {
		if !ll.used[i] {
			k = ll.dfs(i, k, -1)
		}
	}
}

func (ll *LowLink) dfs(idx, k, par int) int {
	ll.used[idx] = true
	ll.ord[idx] = k
	k++
	ll.low[idx] = ll.ord[idx]
	isArticulation := false
	beet := false
	cnt := 0
	for _, e := range ll.g[idx] {
		if e.to == par {
			tmp := beet
			beet = true
			if !tmp {
				continue
			}
		}
		if !ll.used[e.to] {
			cnt++
			k = ll.dfs(e.to, k, idx)
			ll.low[idx] = min(ll.low[idx], ll.low[e.to])
			if par >= 0 && ll.low[e.to] >= ll.ord[idx] {
				isArticulation = true
			}
			if ll.ord[idx] < ll.low[e.to] {
				ll.Bridge = append(ll.Bridge, e)
			}
		} else {
			ll.low[idx] = min(ll.low[idx], ll.ord[e.to])
		}
	}

	if par == -1 && cnt > 1 {
		isArticulation = true
	}
	if isArticulation {
		ll.Articulation = append(ll.Articulation, idx)
	}
	return k
}

func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		size:   n,
		rank:   rank,
		parent: parent,
	}
}

type _UnionFindArray struct {
	size   int
	Part   int
	rank   []int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
