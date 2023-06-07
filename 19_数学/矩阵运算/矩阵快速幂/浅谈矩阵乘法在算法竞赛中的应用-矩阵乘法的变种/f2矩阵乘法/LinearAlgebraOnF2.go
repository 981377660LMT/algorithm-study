// f2上的线性代数
// https://hitonanode.github.io/cplib-cpp/linear_algebra_matrix/linalg_bitset.hpp

package main

import "fmt"

func main() {
	demo := func() {
		ROW, COL := 3, 3
		mat := make([]BitSet64, ROW)
		for i := range mat {
			mat[i] = NewBitset(COL)
		}

		mat[0].Set(0)
		mat[0].Set(1)
		mat[1].Set(1)
		mat[1].Set(2)
		mat[2].Set(0)
		mat[2].Set(2)
		fmt.Println(mat)
		mat = GaussJordan(mat, COL)
		fmt.Println(mat)
	}
	demo()
}

// 高斯-约当消元法.将矩阵原地修改为最简行阶梯矩阵,即对角线上全是1，其他地方的值都为0的对角阵(diagonal matrix).
//  即将Ax=b变为Ux=c.
//  col: 矩阵的列数
//  时间复杂度: O(HW + HW rank(M) / 64)
//  Verified: abc276_h (2000 x 8000)
func GaussJordan(mat []BitSet64, col int) []BitSet64 {
	row := len(mat)
	c := 0
	for h := 0; h < row && c < col; h, c = h+1, c+1 {
		pivot := -1
		for j := h; j < row; j++ {
			if mat[j].Has(c) {
				pivot = j
				break
			}
		}
		if pivot == -1 {
			h--
			continue
		}
		mat[pivot], mat[h] = mat[h], mat[pivot]
		for hh := 0; hh < row; hh++ {
			if hh != h && mat[hh].Has(c) {
				mat[hh].IXOr(mat[h])
			}
		}
	}
	return mat
}

// 01矩阵的秩.
func CalRank(mat []BitSet64, col int) int {
	for h := len(mat) - 1; h >= 0; h-- {
		j := 0
		for j < col && !mat[h].Has(j) {
			j++
		}
		if j < col {
			return h + 1
		}
	}
	return 0
}

// F2矩阵乘法.
func MatMul(mat1, mat2 []BitSet64) []BitSet64 {
	row, col := len(mat1), len(mat2)
	res := make([]BitSet64, row)
	for i := 0; i < row; i++ {
		res[i] = NewBitset(col)
		for j := 0; j < col; j++ {
			if mat1[i].Has(j) {
				res[i].IXOr(mat2[j])
			}
		}
	}
	return res
}

// F2矩阵快速幂.返回一个新的矩阵.
func MatPow(mat []BitSet64, k int) []BitSet64 {
	n := len(mat)
	res := make([]BitSet64, n)
	for i := 0; i < n; i++ {
		res[i] = NewBitset(n)
		res[i].Set(i)
	}
	matCopy := make([]BitSet64, n)
	for i := 0; i < n; i++ {
		matCopy[i] = mat[i].Copy()
	}
	mat = matCopy
	for ; k > 0; k >>= 1 {
		if k&1 == 1 {
			res = MatMul(res, mat)
		}
		mat = MatMul(mat, mat)
	}
	return res
}

// 求解f2上的线性方程组Ax=b.
//  返回值: (一组解, 方程组的自由变量, 是否有解)
//  时间复杂度: O(HW + HW rank(A) / 64 + W^2 len(freedoms))
func SystemOfLinearEquations(A []BitSet64, b []bool, col int) (solution0 BitSet64, freedoms []BitSet64, ok bool) {
	row := len(A)
	if row != len(b) {
		panic("len(A) != len(b)")
	}

	M := make([]BitSet64, row)
	for i := 0; i < row; i++ {
		bs := NewBitset(col + 1)
		for j := 0; j < col; j++ {
			if A[i].Has(j) {
				bs.Set(j)
			}
		}
		if b[i] {
			bs.Set(col)
		}
		M[i] = bs
	}

	M = GaussJordan(M, col+1)
	ss := make([]int, col)
	for i := range ss {
		ss[i] = -1
	}
	var ssNonnegJs []int
	for i := 0; i < row; i++ {
		j := 0
		for j <= col && !M[i].Has(j) {
			j++
		}
		if j == col {
			return
		}
		if j < col {
			ssNonnegJs = append(ssNonnegJs, j)
			ss[j] = i
		}
	}

	ok = true
	solution0 = NewBitset(col)
	for j := 0; j < col; j++ {
		if ss[j] == -1 {
			// This part may require W^2 space complexity in output
			d := NewBitset(col)
			d.Set(j)
			for _, jj := range ssNonnegJs {
				if M[ss[jj]].Has(j) {
					d.Set(jj)
				}
			}
			freedoms = append(freedoms, d)
		} else {
			if M[ss[j]].Has(col) {
				solution0.Set(j)
			}
		}
	}

	return
}

type BitSet64 []uint64

func NewBitset(n int) BitSet64 { return make(BitSet64, n>>6+1) } // (n+_w-1)>>6

func (b BitSet64) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 }
func (b BitSet64) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b BitSet64) Set(p int)      { b[p>>6] |= 1 << (p & 63) }
func (b BitSet64) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) }

func (b BitSet64) Copy() BitSet64 {
	res := make(BitSet64, len(b))
	copy(res, b)
	return res
}

// !f2上的加法
func (b BitSet64) IXOr(c BitSet64) {
	for i, v := range c {
		b[i] ^= v
	}
}

func Xor(a, b BitSet64) BitSet64 {
	res := make(BitSet64, len(a))
	for i, v := range a {
		res[i] = v ^ b[i]
	}
	return res
}
