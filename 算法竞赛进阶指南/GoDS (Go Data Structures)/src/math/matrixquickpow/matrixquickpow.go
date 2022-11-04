// 矩阵快速幂
package matrixquickpow

// https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/math_matrix.go#L70-L117

type matrix [][]int64

func newMatrix(n, m int) matrix {
	a := make(matrix, n)
	for i := range a {
		a[i] = make([]int64, m)
	}
	return a
}

func newIdentityMatrix(n int) matrix {
	a := make(matrix, n)
	for i := range a {
		a[i] = make([]int64, n)
		a[i][i] = 1
	}
	return a
}

func (a matrix) mul(b matrix) matrix {
	const mod int64 = 1e9 + 7 // 998244353
	c := newMatrix(len(a), len(b[0]))
	for i, row := range a {
		for j := range b[0] {
			for k, v := range row {
				c[i][j] = (c[i][j] + v*b[k][j]) % mod // 注：此处不能化简
			}
			if c[i][j] < 0 {
				c[i][j] += mod
			}
		}
	}
	return c
}

func (a matrix) pow(n int64) matrix {
	res := newIdentityMatrix(len(a))
	for ; n > 0; n >>= 1 {
		if n&1 > 0 {
			res = res.mul(a)
		}
		a = a.mul(a)
	}
	return res
}
