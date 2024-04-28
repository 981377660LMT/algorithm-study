package main

import "fmt"

func main() {
	arr2 := NewArray2D(3, 4, func() int { return 0 })
	arr2.Set(1, 2, 3)
	fmt.Println(arr2.ToList())
}

type Array2D[E any] struct {
	n, m int32
	e    func() E
	data []E
}

func NewArray2D[E any](n, m int32, e func() E) *Array2D[E] {
	data := make([]E, n*m)
	for i := int32(0); i < n*m; i++ {
		data[i] = e()
	}
	return &Array2D[E]{n: n, m: m, e: e, data: data}
}

func (a *Array2D[E]) Set(i, j int32, val E) {
	a.data[i*a.m+j] = val
}

func (a *Array2D[E]) Get(i, j int32) E {
	return a.data[i*a.m+j]
}

func (a *Array2D[E]) ToList() [][]E {
	res := make([][]E, a.n)
	for i := int32(0); i < a.n; i++ {
		res[i] = make([]E, a.m)
		for j := int32(0); j < a.m; j++ {
			res[i][j] = a.Get(i, j)
		}
	}
	return res
}
