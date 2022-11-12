// 矩阵快速幂

package matrixquickpow

// !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/math_matrix.go#L70-L117
type Matrix [][]int

func NewMatrix(row, col int) Matrix {
	res := make(Matrix, row)
	for i := range res {
		res[i] = make([]int, col)
	}
	return res
}

func NewIdentityMatrix(n int) Matrix {
	res := make(Matrix, n)
	for i := range res {
		res[i] = make([]int, n)
		res[i][i] = 1
	}
	return res
}

func (matrix Matrix) Mul(other Matrix, mod int) Matrix {
	res := NewMatrix(len(matrix), len(other[0]))
	for i, row := range matrix {
		for j := range other[0] {
			for k, v := range row {
				res[i][j] = (res[i][j] + v*other[k][j]) % mod // 注：此处不能化简
			}
			if res[i][j] < 0 {
				res[i][j] += mod
			}
		}
	}
	return res
}

func (matrix Matrix) Pow(exp int, mod int) Matrix {
	res := NewIdentityMatrix(len(matrix))
	for ; exp > 0; exp >>= 1 {
		if exp&1 > 0 {
			res = res.Mul(matrix, mod)
		}
		matrix = matrix.Mul(matrix, mod)
	}
	return res
}
