// abc429-F - Shortest Path Query-线段树维护矩阵乘法(动态dp)
// https://atcoder.jp/contests/abc429/tasks/abc429_f
// 给定一个 3×N 的栅格，. 为可通行，# 为墙。你要处理 Q 次单点翻转（. 与 # 互换，不会翻到起点 (1,1) 与终点 (3,N)）。
// 每次翻转后，判断能否从 (1,1) 走到 (3,N)（四联通），若能输出最短步数，否则输出 -1。
//
// 思路：矩阵链乘法 + 线段树维护矩阵乘法。
// 将整张 3×N 图拆为 N−1 个“相邻两列”的小块，每个小块 j 对应 3×3 的传递矩阵 Tj，Tj[a][b] 表示在仅使用第 j 与 j+1 列内的格子，从 (a, j) 走到 (b, j+1) 的最短步数（不可达为 INF）。
// 计算 Tj：对 6 个点（两列×3 行）建小图，边为上下/左右（需两端为 .），对每个起点 (a,j) 在小图中 BFS，得到到 (b,j+1) 的距离。
// 整体最短路即为矩阵链的 min-plus 乘积 T1⊗T2⊗…⊗T_{N-1} 的 [0][2] 元素。维护一个线段树存放矩阵并支持区间乘积（这里就是整段），min-plus 乘法是 C[i][k]=min_j(A[i][j]+B[j][k])。
// 单点翻转 (r,c) 只影响 T_{c-1} 与 T_c（若存在），重算这两个块的矩阵并在段树中点更新。查询时取整段乘积的 [0][2]，若为 INF 输出 -1。
// 复杂度：每次更新重算至多 2 个 3×3 矩阵（每个做 3 次 6 点 BFS，常数极小）+ 段树合并 O(log N) 次，每次 3×3 min-plus 乘法常数 27；总计 O((N+Q) log N)。

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int
	fmt.Fscan(in, &N)

	S := make([][]byte, 3)
	open := make([][]bool, 3)
	for i := 0; i < 3; i++ {
		S[i] = make([]byte, N)
		fmt.Fscan(in, &S[i])
		open[i] = make([]bool, N)
		for j := 0; j < N; j++ {
			open[i][j] = S[i][j] == '.'
		}
	}

	st := NewSegmentTree(max(0, N-1), func(i int) E {
		return calcMat(i, open)
	})

	var Q int
	fmt.Fscan(in, &Q)
	for ; Q > 0; Q-- {
		var r, c int
		fmt.Fscan(in, &r, &c)
		r--
		c--

		// 翻转
		open[r][c] = !open[r][c]

		// 受影响的块：c-1 与 c
		if c-1 >= 0 {
			st.Set(c-1, calcMat(c-1, open))
		}
		if c < N-1 {
			st.Set(c, calcMat(c, open))
		}

		res := st.QueryAll()[0][2]
		if res >= INF/2 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, res)
		}
	}
}

type Mat [3][3]int

func idMat() Mat {
	var m Mat
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == j {
				m[i][j] = 0
			} else {
				m[i][j] = INF
			}
		}
	}
	return m
}

// min-plus 乘法
func mul(a, b Mat) Mat {
	var c Mat
	for i := 0; i < 3; i++ {
		for k := 0; k < 3; k++ {
			c[i][k] = INF
			for j := 0; j < 3; j++ {
				if a[i][j] == INF || b[j][k] == INF {
					continue
				}
				c[i][k] = min(c[i][k], a[i][j]+b[j][k])
			}
		}
	}
	return c
}

// 计算第 j 块（列 j 与 j+1）3×3 传递矩阵
func calcMat(j int, open [][]bool) Mat {
	var res Mat
	for i := 0; i < 3; i++ {
		for k := 0; k < 3; k++ {
			res[i][k] = INF
		}
	}
	// 小图：两列×三行，共 6 节点。编号: col*3 + row (col=0→j, col=1→j+1)
	isOpen := func(col, row int) bool {
		if col == 0 {
			return open[row][j]
		}
		return open[row][j+1]
	}
	neighbors := func(u int, buf *[]int) {
		*buf = (*buf)[:0]
		col, row := u/3, u%3
		// 上
		if row > 0 && isOpen(col, row) && isOpen(col, row-1) {
			*buf = append(*buf, col*3+(row-1))
		}
		// 下
		if row < 2 && isOpen(col, row) && isOpen(col, row+1) {
			*buf = append(*buf, col*3+(row+1))
		}
		// 左右（两列之间的水平）
		oc := col ^ 1
		if isOpen(col, row) && isOpen(oc, row) {
			*buf = append(*buf, oc*3+row)
		}
	}

	for sr := 0; sr < 3; sr++ {
		if !open[sr][j] { // 起点关着则不可达
			continue
		}
		// BFS on 6 nodes
		dist := [6]int{INF, INF, INF, INF, INF, INF}
		q := [6]int{}
		h, t := 0, 0
		push := func(x int) {
			q[t] = x
			t++
		}
		pop := func() int {
			x := q[h]
			h++
			return x
		}
		start := 0*3 + sr
		dist[start] = 0
		push(start)
		nb := make([]int, 0, 4)
		for h < t {
			u := pop()
			neighbors(u, &nb)
			for _, v := range nb {
				if dist[v] != INF {
					continue
				}
				dist[v] = dist[u] + 1
				push(v)
			}
		}
		for tr := 0; tr < 3; tr++ {
			res[sr][tr] = dist[1*3+tr]
		}
	}
	return res
}

type E = Mat
type SegmentTree struct {
	n, size int
	seg     []E
}

func (*SegmentTree) e() E        { return idMat() }
func (*SegmentTree) op(a, b E) E { return mul(a, b) }
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
func NewSegmentTree(n int, f func(int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func NewSegmentTreeFrom(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTree) QueryAll() E { return st.seg[1] }
func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}
