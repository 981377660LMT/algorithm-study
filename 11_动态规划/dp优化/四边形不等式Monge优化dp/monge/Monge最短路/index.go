// 王钦石二分(wqs二分)/Monge图d边最短路/aliens dp

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
	CF321E()

	// Yuki952()
	// Yuki705()
}

// P1484 种树
// https://www.luogu.com.cn/problem/P1484
func P1484() {}

// P2619 [国家集训队] Tree I
// https://www.luogu.com.cn/problem/P2619
func P2619() {}

// MST Company
// https://www.luogu.com.cn/problem/CF125E
func CF125E() {}

// Ciel and Gondolas
// https://www.luogu.com.cn/problem/CF321E
//
// 转移代价为二维前缀和.
func CF321E() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, k := io.NextInt(), io.NextInt()
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
		for j := 0; j < n; j++ {
			grid[i][j] = io.NextInt()
		}
	}

	preSum2d := NewPreSum2DFrom(grid)

	f := func(i, j int) int {
		res := preSum2d.QueryRange(i, i, j-1, j-1)
		return res
	}

	res := MongeShortestPathDEdge(n, k, 1e9, f)
	res /= 2
	io.Println(res)
}

// Gosha is hunting
// https://www.luogu.com.cn/problem/CF739E
func CF739E() {}

// 危险的火药库
// 有n个火药库，每个火药库有一个危险程度A[i].
// 现在需要打开某些火药库.
// 总危险程度为 `被打开的火药库连续段的危险程度之和`的平方和.
// 对k=1,...,n,求打开k个火药库时的最小总危险程度.
// n<=3000, A[i]<=1e6
// https://yukicoder.me/problems/no/952
func Yuki952() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	preSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + nums[i]
	}

	// 0 -> n+1 的路径
	// i -> j 的边权为: (sum_{k=i...j-2} nums[k])^2
	cost := func(i, j int) int {
		if j-i <= 1 {
			return 0
		}
		tmp := preSum[j-1] - preSum[i]
		return tmp * tmp
	}

	res := EnumerateMongeDEdgeShortestPath(n+1, cost)
	for k := 1; k <= n; k++ {
		fmt.Fprintln(out, res[n+1-k])
	}

	// for i := 1; i < n+1; i++ {
	// 	res2 := MongeShortestPathDEdge(n+1, i, 1e18, cost)
	// 	if res2 != res[i] {
	// 		panic("res2 != res[i]")
	// 	}
	// }
}

// 清理垃圾
// https://yukicoder.me/problems/no/705
func Yuki705() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	A := make([]int, n)
	X := make([]int, n)
	Y := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &X[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &Y[i])
	}

	cubic := func(x int) int {
		return abs(x * x * x)
	}
	// 人i,垃圾j
	dist := func(i, j int) int {
		return cubic(X[j]-A[i]) + cubic(Y[j])
	}
	f := func(i, j int) int {
		return dist(j-1, i)
	}
	fmt.Fprintln(out, MongeShortestPath(n, f)[n])
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

const INF int = 1e18

// Monge图最短路.求从0点出发到各个点的最短路.
// O(nlogn).
// f(i,j) : 边权函数, 即从i到j的边的权值(0<=i<j<=n).
// https://noshi91.hatenablog.com/entry/2023/02/18/005856
func MongeShortestPath(n int, f func(i, j int) int) []int {
	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	x := make([]int, n+1)

	check := func(from, to int) {
		if from >= to {
			return
		}
		cost := f(from, to)
		if tmp := dp[from] + cost; tmp < dp[to] {
			dp[to] = tmp
			x[to] = from
		}
	}

	var dfs func(int, int)
	dfs = func(l, r int) {
		if l+1 >= r {
			return
		}
		m := (l + r) / 2
		for i := x[l]; i <= x[r]; i++ {
			check(i, m)
		}
		dfs(l, m)
		for i := l + 1; i <= m; i++ {
			check(i, r)
		}
		dfs(m, r)
	}
	dp[0] = 0
	check(0, n)
	dfs(0, n)
	return dp
}

// Monge图d边最短路.
// O(nlogn)求从0到n恰好经过d条边的最短路,其中边权f(from,to)满足Monge性质.
// n: 点数.
// d: 边数.
// f(i,j) : 边权函数, 即从i到j的边的权值(0<=i<j<=n).
// maxWeight: 边权的最大值.
// 返回值为从0到N恰好经过d条边的最短路.
//
//	AliensDp可以归结为这个问题.
//	https://noshi91.github.io/algorithm-encyclopedia/d-edge-shortest-path-monge
//
// 网格图负权最短路 https://zhuanlan.zhihu.com/p/33808530
func MongeShortestPathDEdge(n, d, maxWeight int, f func(i, j int) int) int {
	if d > n {
		panic("d > N")
	}
	cal := func(x int) int {
		g := func(frm, to int) int {
			return f(frm, to) + x
		}
		cost := MongeShortestPath(n, g)[n]
		return cost - x*d
	}
	_, res := FibonacciSearch(cal, -maxWeight, maxWeight+1, false)
	return res
}

// monge图从0到n，经过d边最短路的d=1,...,N的答案.
// 无法到达的点的距离为INF.
// f(i,j) : 边权函数, 即从i到j的边的权值(0<=i<j<=n).
// O(n^2logn)
func EnumerateMongeDEdgeShortestPath(n int, f func(i, j int) int) []int {
	res := make([]int, n+1)
	for i := range res {
		res[i] = INF
	}
	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for d := 1; d <= n; d++ {
		midx := MonotoneMinima2(n+1, n+1, func(j, i int) int {
			if i < j {
				return dp[i] + f(i, j)
			}
			return INF
		})
		for i := n; i >= d; i-- {
			dp[i] = dp[midx[i]] + f(midx[i], i)
		}
		res[d] = dp[n]
	}
	return res
}

// 给定一个 n 行 m 列的矩阵.
// minI=argmin_j(A[i][j]) 单调递增时, 返回 minI.
// f(i, j, k) :
// 比较 A[i][j] 和 A[i][k] (保证 j < k)
// 当 A[i][j] <= A[i][k] 时返回 true.
func MonotoneMinima(row, col int, f func(i, j, k int) bool) []int {
	res := make([]int, row)
	var dfs func(int, int, int, int)
	dfs = func(is, ie, l, r int) {
		if is == ie {
			return
		}
		i := (is + ie) / 2
		m := l
		for k := l + 1; k < r; k++ {
			if !f(i, m, k) {
				m = k
			}
		}
		res[i] = m
		dfs(is, i, l, m+1)
		dfs(i+1, ie, m, r)
	}
	dfs(0, row, 0, col)
	return res
}

// 给定一个 n 行 m 列的矩阵.
// minI=argmin_j(A[i][j]) 单调递增时, 返回 minI.
// get(i, j) : 返回 A[i][j] 的函数
func MonotoneMinima2(row, col int, get func(i, j int) int) []int {
	f := func(i, j, k int) bool {
		return get(i, j) <= get(i, k)
	}
	return MonotoneMinima(row, col, f)
}

// 寻找[start,end)中的一个极值点,不要求单峰性质.
//
//	返回值: (极值点,极值)
func FibonacciSearch(f func(x int) int, start, end int, minimize bool) (int, int) {
	end--
	a, b, c, d := start, start+1, start+2, start+3
	n := 0
	for d < end {
		b = c
		c = d
		d = b + c - a
		n++
	}

	get := func(i int) int {
		if end < i {
			return INF
		}
		if minimize {
			return f(i)
		}
		return -f(i)
	}

	ya, yb, yc, yd := get(a), get(b), get(c), get(d)
	for i := 0; i < n; i++ {
		if yb < yc {
			d = c
			c = b
			b = a + d - c
			yd = yc
			yc = yb
			yb = get(b)
		} else {
			a = b
			b = c
			c = a + d - b
			ya = yb
			yb = yc
			yc = get(c)
		}
	}

	x := a
	y := ya
	if yb < y {
		x = b
		y = yb
	}
	if yc < y {
		x = c
		y = yc
	}
	if yd < y {
		x = d
		y = yd
	}

	if minimize {
		return x, y
	}
	return x, -y

}

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
}

type PreSum2D struct {
	preSum [][]int
}

func NewPreSum2D(row, col int, f func(int, int) int) *PreSum2D {
	preSum := make([][]int, row+1)
	for i := range preSum {
		preSum[i] = make([]int, col+1)
	}
	for r := 0; r < row; r++ {
		for c := 0; c < col; c++ {
			preSum[r+1][c+1] = f(r, c) + preSum[r][c+1] + preSum[r+1][c] - preSum[r][c]
		}
	}
	return &PreSum2D{preSum}
}

func NewPreSum2DFrom(mat [][]int) *PreSum2D {
	return NewPreSum2D(len(mat), len(mat[0]), func(r, c int) int { return mat[r][c] })
}

// 查询sum(A[r1:r2+1, c1:c2+1])的值.
// 0 <= r1 <= r2 < row, 0 <= c1 <= c2 < col.
func (ps *PreSum2D) QueryRange(row1, col1, row2, col2 int) int {
	return ps.preSum[row2+1][col2+1] - ps.preSum[row2+1][col1] - ps.preSum[row1][col2+1] + ps.preSum[row1][col1]
}
