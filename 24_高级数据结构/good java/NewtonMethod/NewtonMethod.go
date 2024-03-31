// 牛顿迭代法求方程的近似解 f(x) = 0 (牛顿求根法解方程).

package main

import "fmt"

func main() {
	eps := 1e-6
	newton := NewNewtonMethod(eps)
	f := func(x float64) float64 { return x*x - 2 }
	df := func(x float64) float64 { return 2 * x }
	x0 := 1.0
	fmt.Println(newton.Search(f, df, x0))
}

// NewtonMethod is a template for solving equation f(x) = 0 using Newton's method.
type NewtonMethod struct {
	absError float64
}

func NewNewtonMethod(absError float64) *NewtonMethod {
	return &NewtonMethod{absError: absError}
}

// Search finds the root of f(x) = 0 using Newton's method.
// f is the function to solve.
// df is the derivative of f.
// x0 is the initial guess.
func (n *NewtonMethod) Search(f, df func(x float64) float64, x0 float64) float64 {
	for abs64(f(x0)) > n.absError {
		x0 = x0 - f(x0)/df(x0)
	}
	return x0
}

func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
