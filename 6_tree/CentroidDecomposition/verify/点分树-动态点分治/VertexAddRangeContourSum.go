// https://judge.yosupo.jp/problem/vertex_add_range_contour_sum_on_tree
// 给定q个操作，操作有两种：
// 0 root x : 将root节点的值加上x (点权加)
// 1 root floor higher: 求出距离root节点距离在[floor,higher)之间的所有节点的值的和 (区间点权和)
// n<=1e5 q<=2e5

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	values := make([]AbleGroup, n)
	for i := range values {
		fmt.Fscan(in, &values[i])
	}
	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}
	cs := NewContourSum(tree, values)
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 0 {
			var root, x int
			fmt.Fscan(in, &root, &x)
			cs.Add(root, x)
		} else {
			var root, floor, higher int
			fmt.Fscan(in, &root, &floor, &higher)
			fmt.Fprintln(out, cs.Sum(root, floor, higher))
		}
	}
}

type AbleGroup = int

func e() AbleGroup                { return 0 }
func op(a, b AbleGroup) AbleGroup { return a + b }
func inv(a AbleGroup) AbleGroup   { return -a }

type ContourSum struct {
	n        int
	g        [][]int
	bit      *_BT
	bitRange [][]int    // 每个重心的每个方向的bit范围
	data     [][][3]int // 方向、离重心的距离、bit的index
}

func NewContourSum(tree [][]int, values []AbleGroup) *ContourSum {
	res := &ContourSum{n: len(tree), g: tree}
	res.build(values)
	return res
}

// root的点权加上value.O(logn).
func (cs *ContourSum) Add(root int, value AbleGroup) {
	for _, d := range cs.data[root] {
		i := d[2]
		cs.bit.Add(i, value)
	}
}

// 查询距离root的距离在[lower, higher)的点的和.O(logn^2).
func (cs *ContourSum) Sum(root int, floor, higher int) AbleGroup {
	res := e()
	for _, d := range cs.data[root] {
		k, x := d[0], d[1]
		lo, hi := floor-x, higher-x
		p := k
		if k < 0 {
			lo -= 2
			hi -= 2
			p = ^k
		}
		n := len(cs.bitRange[p]) - 2
		lo = max(lo, 0)
		hi = min(hi, n+1)
		if lo >= hi {
			continue
		}
		a, b := cs.bitRange[p][lo], cs.bitRange[p][hi]
		val := cs.bit.Query(a, b)
		if k < 0 {
			val = inv(val)
		}
		res = op(res, val)
	}
	return res
}

func (cs *ContourSum) build(values []AbleGroup) {
	N := cs.n
	nextBitIdx := 0
	done := make([]bool, N)
	sz := make([]int, N)
	par := make([]int, N)
	dist := make([]int, N)
	for i := range par {
		par[i] = -1
		dist[i] = -1
	}
	st := [][2]int{{0, N}}
	cs.bitRange = make([][]int, N)
	cs.data = make([][][3]int, N)

	for len(st) > 0 {
		v0, n := st[len(st)-1][0], st[len(st)-1][1]
		st = st[:len(st)-1]
		c := -1
		{
			var dfs func(v int) int
			dfs = func(v int) int {
				sz[v] = 1
				for _, to := range cs.g[v] {
					if to != par[v] && !done[to] {
						par[to] = v
						sz[v] += dfs(to)
					}
				}
				if c == -1 && n-sz[v] <= n/2 {
					c = v
				}
				return sz[v]
			}
			dfs(v0)
		}

		// 从重心开始bfs
		done[c] = true
		{
			off := nextBitIdx
			que := []int{}
			add := func(v, d, p int) {
				if dist[v] != -1 {
					return
				}
				sz[v] = 1
				dist[v] = d
				par[v] = p
				que = append(que, v)
			}

			p := 0
			add(c, 0, -1)
			for p < len(que) {
				v := que[p]
				p++
				for _, to := range cs.g[v] {
					if done[to] {
						continue
					}
					add(to, dist[v]+1, v)
				}
			}

			for i := len(que) - 1; i >= 1; i-- {
				v := que[i]
				sz[par[v]] += sz[v]
			}

			maxD := dist[que[len(que)-1]]
			count := make([]int, maxD+1)
			for _, v := range que {
				cs.data[v] = append(cs.data[v], [3]int{c, dist[v], nextBitIdx})
				nextBitIdx++
				count[dist[v]]++
				par[v] = -1
				dist[v] = -1
			}
			preSum := make([]int, len(count)+1)
			for i := 0; i < len(count); i++ {
				preSum[i+1] = preSum[i] + count[i]
			}
			cs.bitRange[c] = preSum
			for i := range cs.bitRange[c] {
				cs.bitRange[c][i] += off
			}
		}

		// 每个方向bfs
		for _, to := range cs.g[c] {

			off := nextBitIdx
			nbd := to
			if done[nbd] {
				continue
			}
			K := len(cs.bitRange)
			que := []int{}
			add := func(v, d int) {
				if dist[v] != -1 || v == c {
					return
				}
				dist[v] = d
				que = append(que, v)
			}
			p := 0
			add(nbd, 0)
			for p < len(que) {
				v := que[p]
				p++
				for _, to := range cs.g[v] {
					if done[to] {
						continue
					}
					add(to, dist[v]+1)
				}
			}

			maxD := dist[que[len(que)-1]]
			count := make([]int, maxD+1)
			for _, v := range que {
				cs.data[v] = append(cs.data[v], [3]int{^K, dist[v], nextBitIdx})
				nextBitIdx++
				count[dist[v]]++
				par[v] = -1
				dist[v] = -1
			}
			preSum := make([]int, len(count)+1)
			for i := 0; i < len(count); i++ {
				preSum[i+1] = preSum[i] + count[i]
			}
			cs.bitRange = append(cs.bitRange, preSum)
			for i := range cs.bitRange[K] {
				cs.bitRange[K][i] += off
			}
			st = append(st, [2]int{nbd, sz[nbd]})
		}
	}

	// build bit
	bitRaw := make([]AbleGroup, nextBitIdx)
	for i := 0; i < N; i++ {
		for _, d := range cs.data[i] {
			bitRaw[d[2]] = values[i]
		}
	}
	cs.bit = &_BT{}
	cs.bit.Build(bitRaw)
}

type _BT struct {
	n    int
	data []AbleGroup
}

func (bit *_BT) Add(i int, x AbleGroup) {
	for i++; i <= bit.n; i += i & -i {
		bit.data[i-1] = op(bit.data[i-1], x)
	}
}

func (bit *_BT) Query(start, end int) AbleGroup {
	if start < 0 {
		start = 0
	}
	if end > bit.n {
		end = bit.n
	}
	if start == 0 {
		return bit.queryPrefix(end)
	}
	pos, neg := e(), e()
	for start < end {
		pos = op(pos, bit.data[end-1])
		end -= end & -end
	}
	for end < start {
		neg = op(neg, bit.data[start-1])
		start -= start & -start
	}
	return op(pos, inv(neg))
}

func (bit *_BT) Build(values []AbleGroup) {
	n := len(values)
	data := append(values[:0:0], values...)
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = op(data[i-1], data[j-1])
		}
	}
	bit.n = n
	bit.data = data
}

func (bit *_BT) queryPrefix(right int) AbleGroup {
	if right > bit.n {
		right = bit.n
	}
	res := e()
	for ; right > 0; right -= right & -right {
		res = op(res, bit.data[right-1])
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
