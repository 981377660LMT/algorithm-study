package main

import "fmt"

func main() {
	LS := NewLatinSquare([][]int{{0, 1, 2}, {1, 2, 0}, {2, 0, 1}})
	fmt.Println(LS.ToMatrix())
	LS.Up(1)
	fmt.Println(LS.ToMatrix())
	LS.Left(1)
	fmt.Println(LS.ToMatrix())
	LS.Down(1)
	fmt.Println(LS.ToMatrix())
	LS.Plus(1)
	fmt.Println(LS.ToMatrix())
	LS.Down(1)
	fmt.Println(LS.ToMatrix())

	LS.Transpose()
	fmt.Println(LS.ToMatrix())
	fmt.Println("t")
	LS.RowInv()
	fmt.Println(LS.ToMatrix())
	LS.ColInv()
	fmt.Println(LS.ToMatrix())
}

// 拉丁方是一个n^2的矩阵，矩阵的每个元素都是0到n-1之间的整数，且每一行每一列都是一个置换。
// 支持O(1)的左右上下滚动，行逆，列逆，转置操作
// 支持O(1)单个单元格查询操作
type LatinSquare struct {
	ij   [][]int
	iv   [][]int
	jv   [][]int
	n    int
	cast []*elementModifyWrapper
}

func NewLatinSquare(matrix [][]int) *LatinSquare {
	return newLatinSquare(matrix, true)
}

func newLatinSquare(matrix [][]int, copy bool) *LatinSquare {
	n := len(matrix)
	ij := matrix
	if copy {
		newMatrix := make([][]int, n)
		for i := 0; i < n; i++ {
			newMatrix[i] = append(matrix[i][:0:0], matrix[i]...)
		}
		ij = newMatrix
	}
	iv := make([][]int, n)
	jv := make([][]int, n)
	for i := 0; i < n; i++ {
		iv[i] = make([]int, n)
		jv[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := matrix[i][j]
			ij[i][j] = v
			iv[i][v] = j
			jv[j][v] = i
		}
	}
	cast := make([]*elementModifyWrapper, 3)
	cast[0] = newElementModifyWrapper(
		newElement(
			func(i, j, v int) int { return i },
			func() int { return 0 },
		),
	)
	cast[1] = newElementModifyWrapper(
		newElement(
			func(i, j, v int) int { return j },
			func() int { return 1 },
		),
	)
	cast[2] = newElementModifyWrapper(
		newElement(
			func(i, j, v int) int { return v },
			func() int { return 2 },
		),
	)
	return &LatinSquare{ij: ij, iv: iv, jv: jv, n: n, cast: cast}
}

// O(1)
func (ls *LatinSquare) Get(i, j int) int {
	kind := ls.cast[2].Index()
	var table [][]int
	var inv int
	if kind == 0 {
		table = ls.jv
		if ls.cast[0].Index() == 0 {
			inv = 0
		} else {
			inv = 1
		}
	} else if kind == 1 {
		table = ls.iv
		if ls.cast[0].Index() == 0 {
			inv = 0
		} else {
			inv = 1
		}
	} else {
		table = ls.ij
		if ls.cast[0].Index() == 0 {
			inv = 0
		} else {
			inv = 1
		}
	}
	if inv == 1 {
		i, j = j, i
	}
	mi := calMod(i-ls.cast[0^inv].mod, ls.n)
	mj := calMod(j-ls.cast[1^inv].mod, ls.n)
	return calMod(table[mi][mj]+ls.cast[2].mod, ls.n)
}

func (ls *LatinSquare) Up(x int) {
	ls.Down(-x)
}

func (ls *LatinSquare) Left(x int) {
	ls.Right(-x)
}

func (ls *LatinSquare) Down(x int) {
	ls.cast[0].Modify(x, ls.n)
}

func (ls *LatinSquare) Right(x int) {
	ls.cast[1].Modify(x, ls.n)
}

func (ls *LatinSquare) RowInv() {
	ls.cast[1], ls.cast[2] = ls.cast[2], ls.cast[1]
}

func (ls *LatinSquare) ColInv() {
	ls.cast[0], ls.cast[2] = ls.cast[2], ls.cast[0]
}

func (ls *LatinSquare) Plus(x int) {
	ls.cast[2].Modify(x, ls.n)
}

func (ls *LatinSquare) Transpose() {
	ls.cast[0], ls.cast[1] = ls.cast[1], ls.cast[0]
}

func (ls *LatinSquare) ToMatrix() [][]int {
	matrix := make([][]int, ls.n)
	for i := 0; i < ls.n; i++ {
		matrix[i] = make([]int, ls.n)
		for j := 0; j < ls.n; j++ {
			matrix[i][j] = ls.Get(i, j)
		}
	}
	return matrix
}

type element struct {
	apply func(i, j, v int) int
	index func() int
}

func newElement(apply func(i, j, v int) int, index func() int) *element {
	return &element{apply: apply, index: index}
}

type elementModifyWrapper struct {
	e   *element
	mod int
}

func newElementModifyWrapper(e *element) *elementModifyWrapper {
	return &elementModifyWrapper{e: e}
}

func (wrapper *elementModifyWrapper) Modify(x int, n int) {
	wrapper.mod = calMod(wrapper.mod+x, n)
}

func (wrapper *elementModifyWrapper) Apply(i, j, v int) int {
	return wrapper.e.apply(i, j, v) + wrapper.mod
}

func (wrapper *elementModifyWrapper) Index() int {
	return wrapper.e.index()
}

func calMod(x int, mod int) int {
	if x < -mod || x >= mod {
		x %= mod
	}
	if x < 0 {
		x += mod
	}
	return x
}
