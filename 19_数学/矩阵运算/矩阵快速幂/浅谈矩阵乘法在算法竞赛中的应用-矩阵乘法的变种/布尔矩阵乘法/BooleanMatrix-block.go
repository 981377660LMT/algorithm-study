// !布尔矩阵乘法(Boolean Matrix Multiplication, BMM)
// 输入和输出矩阵的元素均为布尔值。
// 按矩阵乘法的公式运算时，可以把“乘”看成and，把“加”看成or
// 对矩阵乘法 C[i][j] |= A[i][k] & B[k][j], 它的一个直观意义是把A的行和B的列看成集合，
// A的第i行包含元素k当且仅当A[i][k]=1。
// B的第j列包含元素k当且仅当B[k][j]=1。
// !那么C[i][j]代表A的第i行和B的第j列是否包含公共元素。
//
// 一个应用是传递闭包(Transitive Closure)的加速计算。
//
// https://zhuanlan.zhihu.com/p/631804105
// https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/math/matrix_bool.cc#L4
//
//
// Boolean matrix
//
// Description:
//   This admits very fast operations for boolean matrices.
//
// Algorithm:
//   Block matrix decomposition technique:
//   For a matrix A of size n x n, we split A as the block matrix
//   each block is of size n/W x n/W. Here, computation of each
//   W x W block is performed by bit operations;
//
// Complexity: (in practice)
//   50--60 times faster than the naive implementation.
//
// 这里是分块实现(四毛子方法), 参考 https://www.cnblogs.com/whx1003/p/13996517.html

package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	mat := NewBooleanMatrix(3, 3)
	mat.Set(0, 0, true)
	mat.Set(0, 1, true)
	mat.Set(1, 2, true)
	mat.Set(1, 0, true)
	fmt.Println(mat)
	fmt.Println(Transpose(mat))
	fmt.Println(Mul(Mul(mat, mat), mat), Pow(mat, 3))
	fmt.Println(Eye(8))

	// 5000*5000的矩阵乘法 => 643.2133ms
	time1 := time.Now()
	eye := Eye(5000)
	Mul(eye, eye)
	time2 := time.Now()
	fmt.Println(time2.Sub(time1))

	// 2000*2000的传递闭包 => 620.574ms
	time3 := time.Now()
	n := 2000
	newEye := Eye(n)
	TransitiveClosure(newEye)
	time4 := time.Now()
	fmt.Println(time4.Sub(time3))
}

func NewBooleanMatrix(row, col int) *BooleanMatrix {
	return &BooleanMatrix{ROW: row, COL: col, x: make([]uint64, (1+row>>3)*(1+col>>3))}
}

func Eye(n int) *BooleanMatrix {
	res := NewBooleanMatrix(n, n)
	size := 1 + n>>3
	for i := 0; i < size; i++ {
		res.x[i*size+i] = 0x8040201008040201
	}
	return res
}

// 返回一个新的矩阵,值为mat1*mat2.
func Mul(mat1, mat2 *BooleanMatrix) *BooleanMatrix {
	r1, c1, r2, c2 := 1+mat1.ROW>>3, 1+mat1.COL>>3, 1+mat2.ROW>>3, 1+mat2.COL>>3
	res := NewBooleanMatrix(mat1.ROW, mat2.COL)
	for i := 0; i < r1; i++ {
		for k := 0; k < r2; k++ {
			for j := 0; j < c2; j++ {
				res.x[i*c2+j] |= _mul(mat1.x[i*c1+k], mat2.x[k*c2+j])
			}
		}
	}
	return res
}

// 求方阵mat的k次幂.
func Pow(mat *BooleanMatrix, k int) *BooleanMatrix {
	res := Eye(mat.ROW)
	for ; k > 0; k >>= 1 {
		if k&1 != 0 {
			res = Mul(res, mat)
		}
		mat = Mul(mat, mat)
	}
	return res
}

func Transpose(mat *BooleanMatrix) *BooleanMatrix {
	row, col := 1+mat.ROW>>3, 1+mat.COL>>3
	res := NewBooleanMatrix(mat.COL, mat.ROW)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			res.x[j*row+i] = _transpose(mat.x[i*col+j])
		}
	}
	return res
}

// 返回一个新的矩阵,值为mat1|mat2.
func Add(mat1, mat2 *BooleanMatrix) *BooleanMatrix {
	res := mat1.Copy()
	row, col := 1+mat1.ROW>>3, 1+mat1.COL>>3
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			pos := i*col + j
			res.x[pos] |= mat2.x[pos]
		}
	}
	return res
}

// (A + I)^n 是传递闭包.
func TransitiveClosure(mat *BooleanMatrix) *BooleanMatrix {
	n := mat.ROW
	return Pow(Add(mat, Eye(n)), n)
}

//
//
type BooleanMatrix struct {
	ROW, COL int
	x        []uint64
}

func (mat *BooleanMatrix) Get(i, j int) bool {
	return mat.x[(i>>3)<<3|j>>3]&(1<<((i&7)<<3|j&7)) != 0 // % 8 => & 7
}

func (mat *BooleanMatrix) Set(i, j int, b bool) {
	if b {
		mat.x[(i>>3)<<3|j>>3] |= (1 << ((i&7)<<3 | j&7))
	} else {
		mat.x[(i>>3)<<3|j>>3] &= ^(1 << ((i&7)<<3 | j&7))
	}
}

func (mat *BooleanMatrix) Copy() *BooleanMatrix {
	return &BooleanMatrix{ROW: mat.ROW, COL: mat.COL, x: append(mat.x[:0:0], mat.x...)}
}

// To 2D grid.
func (mat *BooleanMatrix) String() string {
	grid := make([][]int, mat.ROW)
	for i := 0; i < mat.ROW; i++ {
		grid[i] = make([]int, mat.COL)
		for j := 0; j < mat.COL; j++ {
			if mat.Get(i, j) {
				grid[i][j] = 1
			} else {
				grid[i][j] = 0
			}
		}
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("BooleanMatrix(%d,%d)\n", mat.ROW, mat.COL))
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			sb.WriteString(fmt.Sprintf("%d ", grid[i][j]))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// C[i][j] |= A[i][k] & B[k][j]
func _mul(a, b uint64) uint64 {
	u := uint64(0xff)
	v := uint64(0x101010101010101)
	c := uint64(0)
	for a != 0 && b != 0 {
		c |= (((a & v) * u) & ((b & u) * v))
		a >>= 1
		b >>= 8
	}
	return c
}

func _transpose(a uint64) uint64 {
	t := (a ^ (a >> 7)) & 0x00aa00aa00aa00aa
	a = a ^ t ^ (t << 7)
	t = (a ^ (a >> 14)) & 0x0000cccc0000cccc
	a = a ^ t ^ (t << 14)
	t = (a ^ (a >> 28)) & 0x00000000f0f0f0f0
	a = a ^ t ^ (t << 28)
	return a
}
