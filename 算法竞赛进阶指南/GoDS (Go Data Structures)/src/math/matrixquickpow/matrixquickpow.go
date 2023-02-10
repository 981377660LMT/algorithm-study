// 矩阵快速幂
// !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/math_matrix.go#L70-L117

package main

type M = [][]int

func NewMatrix(row, col int) M {
	res := make(M, row)
	for i := range res {
		res[i] = make([]int, col)
	}
	return res
}

func NewIdentityMatrix(n int) M {
	res := make(M, n)
	for i := range res {
		res[i] = make([]int, n)
		res[i][i] = 1
	}
	return res
}

func Mul(m1, m2 M, mod int) M {
	res := NewMatrix(len(m1), len(m2[0]))
	for i, row := range m1 {
		for j := range m2[0] {
			for k, v := range row {
				res[i][j] = (res[i][j] + v*m2[k][j]) % mod
			}
			if res[i][j] < 0 {
				res[i][j] += mod
			}
		}
	}
	return res
}

func Pow(m M, exp, mod int) M {
	res := NewIdentityMatrix(len(m))
	for ; exp > 0; exp >>= 1 {
		if exp&1 > 0 {
			res = Mul(res, m, mod)
		}
		m = Mul(m, m, mod)
	}
	return res
}

// 矩阵求逆、行列式、高斯消元
// https://github.dev/EndlessCheng/codeforces-go/blob/3dd70515200872705893d52dc5dad174f2c3b5f3/copypasta/math_matrix.go#L200
