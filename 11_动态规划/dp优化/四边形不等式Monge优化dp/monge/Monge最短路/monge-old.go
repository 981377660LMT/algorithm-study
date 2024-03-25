// https://maspypy.github.io/library/convex/monge.hpp

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF868F()
}

// Yet Another Minimization Problem (决策单调性+莫队)
// https://www.luogu.com.cn/problem/CF868F
// 有一个长度为 n 的序列，要求将其分成 k 个子段，每个子段的花费是子段内相同元素的对数，求最小花费。
// dp[k][i] 表示前 i 个元素分成 k 个子段的最小花费。
func CF868F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	D := NewDictionary()
	for i, v := range nums {
		nums[i] = D.Id(v)
	}

	counter := make([]int, D.Size())
	left, right, cost := 0, 0, 0
	add := func(i int) {
		cost += counter[nums[i]]
		counter[nums[i]]++
	}
	remove := func(i int) {
		counter[nums[i]]--
		cost -= counter[nums[i]]
	}
	f := func(l, r int) int {
		for left > l {
			left--
			add(left)
		}
		for right < r {
			add(right)
			right++
		}
		for left < l {
			remove(left)
			left++
		}
		for right > r {
			right--
			remove(right)
		}
		return cost
	}

	res := MongeShortestPathDEdge(n, k, 1e18, f)
	fmt.Fprintln(out, res)

}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

const INF int = 1e18

// dp[j]=min(dp[i]+f(i,j)) (0<=i<j<=n)
//
//	!f(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
//	f函数需要满足Monge性质.
func MongeShortestPath(N int, f func(i, j int) int) []int {
	dp := make([]int, N+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	larsch := _NewLARSCH(N, func(i, j int) int {
		i++
		if i <= j {
			return INF
		}
		return dp[j] + f(j, i)
	})
	for r := 1; r <= N; r++ {
		l := larsch.GetArgmin()
		dp[r] = dp[l] + f(l, r)
	}
	return dp
}

// O(nlogn)求从0到n恰好经过d条边的最短路,其中边权f(from,to)满足Monge性质.
//
//	d: 边数, fLimit: 边权的最大值, f: 边权函数
//	AliensDp可以归结为这个问题.
//	https://noshi91.github.io/algorithm-encyclopedia/d-edge-shortest-path-monge
//
// 网格图负权最短路 https://zhuanlan.zhihu.com/p/33808530
func MongeShortestPathDEdge(N, d, fLimit int, f func(i, j int) int) int {
	if d > N {
		panic("d > N")
	}
	calcL := func(lambda int) int {
		cost := func(from, to int) int {
			return f(from, to) + lambda
		}
		dp := MongeShortestPath(N, cost)
		return dp[N] - lambda*d
	}
	_, fx := _FibonacciSearch(calcL, -fLimit, fLimit+1, false)
	return fx
}

// 检验[0,n]范围内的f是否满足Monge性质.
func CheckMonge(n int, f func(i, j int) int) bool {
	for l := 0; l <= n; l++ {
		for k := 0; k < l; k++ {
			for j := 0; j < k; j++ {
				for i := 0; i < j; i++ {
					lhs := f(i, l) + f(j, k)
					rhs := f(i, k) + f(j, l)
					if lhs < rhs {
						return false
					}
				}
			}
		}
	}
	return true
}

// ndp[j] = min(dp[i] + f(i,j))
func MongeDpUpdate(n int, dp []int, f func(i, j int) int) []int {
	if len(dp) != n+1 {
		panic("len(dp) != n+1")
	}
	choose := func(i, j, k int) int {
		if i <= k {
			return j
		}
		if dp[j]+f(j, i) > dp[k]+f(k, i) {
			return k
		}
		return j
	}
	I := _SMAWK(n+1, n+1, choose)
	ndp := make([]int, n+1)
	for i := range ndp {
		ndp[i] = INF
	}
	for j := range ndp {
		i := I[j]
		ndp[j] = dp[i] + f(i, j)
	}
	return ndp
}

// choose: func(i, j, k int) int 选择(i,j)和(i,k)中的哪一个(j or k)
//
//	返回值: minArg[i] 表示第i行的最小值的列号
func _SMAWK(H, W int, choose func(i, j, k int) int) (minArg []int) {
	var dfs func(X, Y []int) []int
	dfs = func(X, Y []int) []int {
		n := len(X)
		if n == 0 {
			return nil
		}
		YY := []int{}
		for _, y := range Y {
			for len(YY) > 0 {
				py := YY[len(YY)-1]
				x := X[len(YY)-1]
				if choose(x, py, y) == py {
					break
				}
				YY = YY[:len(YY)-1]
			}
			if len(YY) < len(X) {
				YY = append(YY, y)
			}
		}
		XX := []int{}
		for i := 1; i < len(X); i += 2 {
			XX = append(XX, X[i])
		}
		II := dfs(XX, YY)
		I := make([]int, n)
		for i, v := range II {
			I[i+i+1] = v
		}
		p := 0
		for i := 0; i < n; i += 2 {
			var lim int
			if i+1 == n {
				lim = Y[len(Y)-1]
			} else {
				lim = I[i+1]
			}
			best := Y[p]
			for Y[p] < lim {
				p++
				best = choose(X[i], best, Y[p])
			}
			I[i] = best
		}

		return I
	}

	X, Y := make([]int, H), make([]int, W)
	for i := range X {
		X[i] = i
	}
	for i := range Y {
		Y[i] = i
	}
	return dfs(X, Y)
}

type _LARSCH struct {
	base *_reduceRow
}

func _NewLARSCH(n int, f func(i, j int) int) *_LARSCH {
	res := &_LARSCH{base: _newReduceRow(n)}
	res.base.setF(f)
	return res
}

func (l *_LARSCH) GetArgmin() int {
	return l.base.getArgmin()
}

type _reduceRow struct {
	n      int
	f      func(i, j int) int
	curRow int
	state  int
	rec    *_reduceCol
}

func _newReduceRow(n int) *_reduceRow {
	res := &_reduceRow{n: n}
	m := n / 2
	if m != 0 {
		res.rec = _newReduceCol(m)
	}
	return res
}

func (r *_reduceRow) setF(f func(i, j int) int) {
	r.f = f
	if r.rec != nil {
		r.rec.setF(func(i, j int) int {
			return f(2*i+1, j)
		})
	}
}

func (r *_reduceRow) getArgmin() int {
	curRow := r.curRow
	r.curRow += 1
	if curRow%2 == 0 {
		prevArgmin := r.state
		var nextArgmin int
		if curRow+1 == r.n {
			nextArgmin = r.n - 1
		} else {
			nextArgmin = r.rec.getArgmin()
		}
		r.state = nextArgmin
		ret := prevArgmin
		for j := prevArgmin + 1; j <= nextArgmin; j += 1 {
			if r.f(curRow, ret) > r.f(curRow, j) {
				ret = j
			}
		}
		return ret
	}

	if r.f(curRow, r.state) <= r.f(curRow, curRow) {
		return r.state
	}
	return curRow
}

type _reduceCol struct {
	n      int
	f      func(i, j int) int
	curRow int
	cols   []int
	rec    *_reduceRow
}

func _newReduceCol(n int) *_reduceCol {
	return &_reduceCol{n: n, rec: _newReduceRow(n)}
}

func (c *_reduceCol) setF(f func(i, j int) int) {
	c.f = f
	c.rec.setF(func(i, j int) int {
		return f(i, c.cols[j])
	})
}

func (r *_reduceCol) getArgmin() int {
	curRow := r.curRow
	r.curRow += 1
	var cs []int
	if curRow == 0 {
		cs = []int{0}
	} else {
		cs = []int{2*curRow - 1, 2 * curRow}
	}

	for _, j := range cs {
		for {
			size := len(r.cols)
			flag := size != curRow && r.f(size-1, r.cols[size-1]) > r.f(size-1, j)
			if !flag {
				break
			}
			r.cols = r.cols[:size-1]
		}
		if len(r.cols) != r.n {
			r.cols = append(r.cols, j)
		}
	}
	return r.cols[r.rec.getArgmin()]
}

// 寻找[start,end)中的一个极值点,不要求单峰性质.
// GoldenSectionSearch, 黄金比搜索.
//
//	返回值: (极值点,极值)
func _FibonacciSearch(f func(x int) int, start, end int, minimize bool) (int, int) {
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
	_valueToId map[V]int32
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int32{},
	}
}

func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return int(res)
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = int32(id)
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
