// https://hitonanode.github.io/cplib-cpp/combinatorial_opt/convex_sum.hpp
// 自变量和为定值的条件下求凸函数的最小值
// https://codeforces.com/contest/1344/problem/D
// https://yukicoder.me/problems/no/1495

// 现在有n个变量 x1,x2,...,xn,每个变量有一个取值范围[loweri,upperi]
// 给定自变量限制 sum(x1,x2,...,xn) = c
// 求y的最小值,其中y = sum(f1(x1),f2(x2),...,fm(xn)) (f是关于xi的凸函数)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// 采购帽子
	// 班级决定采购帽子.帽子一共有N个种类,老师决定从中购买K个帽子.
	// 同学们给出了M个建议,每个建议是一个(xi,yi)对,表示希望购买第xi种帽子yi个.
	// 为了迎合同学们的需求,老师决定购买第i种帽子的数量为Bi个.
	// 为了接近同学们的意见,要最小化: sum(Bxi-yi)^2 (1<=i<=M) 之和
	// 即求 ∑(Bxi-yi)^2 (1<=i<=M) 的最小值
	// N,M<=2e5 K<=1e9 1<=Xi<=n
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for t := 0; t < T; t++ {
		var N, M, K int
		fmt.Fscan(in, &N, &M, &K)
		// 这里的函数 f[i] = sum(Bxi-yi)^2 (1<=i<=M)
		A, B, C := make([]int, N), make([]int, N), make([]int, N)
		for i := 0; i < M; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x--
			A[x]++
			B[x] -= y * 2
			C[x] += y * y
		}
		funcs := make([]SlopeFunc, N)
		for i := 0; i < N; i++ {
			funcs[i] = NewQuadratic(A[i], B[i], C[i], 0, K)
		}
		coeff := make([]int, N)
		for i := 0; i < N; i++ {
			coeff[i] = 1
		}
		res, _, _ := MinConvexSumUnderLinearConstraint(coeff, funcs, K)
		fmt.Fprintln(out, res)
	}
}

const INF int = 1e18

// !minimize sum(fi(xij) for j in range(1, ki+1) for i in range(1, n+1))
//  k: coefficient of each variable
//  f: convex function
//  c: constraint (sum of all variables)
//  return: (y, [[(x_i, # of such x_i), ... ], ...])
func MinConvexSumUnderLinearConstraint(k []int, f []SlopeFunc, c int) (minimum int, res [][][2]int, ok bool) {
	if len(k) != len(f) {
		panic("len(k) != len(f)")
	}
	lowerSum, upperSum := 0, 0
	for _, func_ := range f {
		lowerSum += func_.getLower()
		upperSum += func_.getUpper()
	}
	if lowerSum > c || upperSum < c {
		return
	}

	n := len(k)
	few, enough := -INF, INF
	for enough-few > 1 {
		slope := few + (enough-few)/2
		cnt := 0
		for i := 0; i < n; i++ {
			tmp := f[i].Slope(slope)
			cnt += tmp * k[i]
			if cnt >= c {
				break
			}
		}
		if cnt >= c {
			enough = slope
		} else {
			few = slope
		}
	}

	res = make([][][2]int, n)
	additional := []int{}
	ctmp := 0
	for i := 0; i < n; i++ {
		xLower := f[i].Slope(few)
		xUpper := f[i].Slope(few + 1)
		ctmp += k[i] * xLower
		res[i] = append(res[i], [2]int{xLower, k[i]})
		if xLower < xUpper {
			additional = append(additional, i)
		}
		minimum += k[i] * f[i].Eval(xLower)
	}

	minimum += (c - ctmp) * (few + 1)
	for len(additional) > 0 {
		i := additional[len(additional)-1]
		additional = additional[:len(additional)-1]
		add := 0
		if c-ctmp > k[i] {
			add = k[i]
		} else {
			add = c - ctmp
		}
		x := res[i][0][0]
		if add != 0 {
			res[i][0][1] -= add
			if res[i][0][1] == 0 {
				res[i] = res[i][:len(res[i])-1]
			}
			res[i] = append(res[i], [2]int{x + 1, add})
			ctmp += add
		}
	}

	ok = true
	return
}

type SlopeFunc interface {
	Slope(s int) int
	Eval(x int) int
	getLower() int
	getUpper() int
}

// ax^2 + bx + c (convex), lower <= x <= upper
type Quadratic struct{ a, b, c, lower, upper int }

func NewQuadratic(a, b, c, lower, upper int) *Quadratic { return &Quadratic{a, b, c, lower, upper} }

func (q *Quadratic) Slope(s int) int {
	if q.a == 0 {
		if q.b <= s {
			return q.upper
		}
		return q.lower
	}
	res := (s + q.a - q.b) / (q.a * 2)
	if res > q.upper {
		return q.upper
	}
	if res < q.lower {
		return q.lower
	}
	return res
}

func (q *Quadratic) Eval(x int) int { return (q.a*x+q.b)*x + q.c }

// f(x) - f(x - 1)
func (q *Quadratic) nextCost(x int) int { return 2*q.a*x - q.a + q.b }
func (q *Quadratic) getLower() int      { return q.lower }
func (q *Quadratic) getUpper() int      { return q.upper }

// x^3 - ax, x >= 0 (convex)
type Cubic struct {
	a, lower, upper int
}

func NewCubic(a, upper int) *Cubic { return &Cubic{a, 0, upper} }
func (c *Cubic) Slope(s int) int {
	lo, hi := c.lower, c.upper+1
	for hi-lo > 1 {
		mid := (lo + hi) / 2
		if c.nextCost(mid) <= s {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo
}

func (c *Cubic) Eval(x int) int { return (x*x - c.a) * x }

// f(x) - f(x - 1)
func (c *Cubic) nextCost(x int) int { return 3*x*x - 3*x + 1 - c.a }
func (q *Cubic) getLower() int      { return q.lower }
func (q *Cubic) getUpper() int      { return q.upper }
