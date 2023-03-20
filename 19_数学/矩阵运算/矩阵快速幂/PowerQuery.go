// 幂运算预处理

package main

import "fmt"

const MOD int = 1e9 + 7

func main() {
	// !快速幂预处理, 用于多次查询.
	pow := NewPowerQuery(2, func() E { return 1 }, func(a, b E) E { return a.(int) * b.(int) % MOD }, 32)
	fmt.Println(pow.Pow(1e18))

	// !矩阵快速幂预处理, 用于多次查询.
	T := [][]int{{1, 1}, {1, 0}}
	matpow := NewPowerQuery(
		T,
		func() E { return eye(len(T)) },
		func(a, b E) E { return matMul(a.([][]int), b.([][]int), MOD) },
		32,
	)
	fmt.Println(matpow.Pow(10)) // [[89 55] [55 34]]

	_ = []interface{}{pow, matpow}
}

type E = interface{}

type PowerQuery struct {
	data    [][]E
	logBase int
	e       func() E
	mul     func(a, b E) E
}

// 幂运算预处理.
//  base: 幂运算的基.
//  e: Monoid 的单位元.
//  mul: Monoid 的乘法.
//  logBase: 幂运算的 log 底数(一般取32).
func NewPowerQuery(base E, e func() E, mul func(a, b E) E, logBase int) *PowerQuery {
	res := &PowerQuery{logBase: logBase, e: e, mul: mul}
	res.data = [][]E{res.makePow(base)}
	return res
}

func (pq *PowerQuery) Pow(n int) E {
	res := pq.e()
	k := 0
	for n > 0 {
		mod := n % pq.logBase
		n /= pq.logBase
		if len(pq.data) == k {
			pq.data = append(pq.data, pq.makePow(pq.data[k-1][len(pq.data[k-1])-1]))
		}
		res = pq.mul(res, pq.data[k][mod])
		k++
	}
	return res
}

func (pq *PowerQuery) makePow(e E) []E {
	res := []E{pq.e()}
	for i := 0; i < pq.logBase; i++ {
		res = append(res, pq.mul(res[len(res)-1], e))
	}
	return res
}

//
//
func newMatrix(row, col int) [][]int {
	res := make([][]int, row)
	for i := range res {
		res[i] = make([]int, col)
	}
	return res
}

func eye(n int) [][]int {
	res := make([][]int, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
		res[i][i] = 1
	}
	return res
}

func matMul(a, b [][]int, mod int) [][]int {
	res := newMatrix(len(a), len(b[0]))
	for i := 0; i < len(a); i++ {
		for k := 0; k < len(b); k++ {
			for j := 0; j < len(b[0]); j++ {
				res[i][j] = (res[i][j] + a[i][k]*b[k][j]) % mod
				if res[i][j] < 0 {
					res[i][j] += mod
				}
			}
		}
	}
	return res
}
