// 矩阵快速幂
// NewMatPow(base,mod,cacheLevel): 带缓存的矩阵快速幂,适合多次查询;
// Pow(m1,exp,mod): 普通的矩阵快速幂;
// Mul(m1,m2,mod): 矩阵乘法.

package main

type M = [][]int

type MatPowWithCache struct {
	n          int
	mod        int
	base       []int
	cacheLevel int
	useCache   bool
	cache      [][][]int
}

// 带缓存的矩阵快速幂,适合多次查询
//  base: 转移矩阵,必须是方阵;
//  mod: 快速幂的模;
//  cacheLevel: 矩阵快速幂的log底数.启用缓存时一般设置为`4`;
//  当调用 pow 次数很多时,可设置为 `16`;
//  小于 `2` 时不启用缓存.
func NewMatPow(base M, mod, cacheLevel int) *MatPowWithCache {
	res := &MatPowWithCache{}
	n := len(base)
	base_ := make([]int, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			base_[i*n+j] = base[i][j]
		}
	}
	useCache := cacheLevel >= 2

	res.n = n
	res.mod = mod
	res.base = base_
	res.cacheLevel = cacheLevel
	res.useCache = useCache
	if useCache {
		cache_ := make([][][]int, cacheLevel-1)
		for i := 0; i < cacheLevel-1; i++ {
			cache_[i] = make([][]int, 0, 64)
		}
		res.cache = cache_
	}

	return res
}

// 时间复杂度: O(n^3*log(exp))
func (mp *MatPowWithCache) Pow(exp int) M {
	if !mp.useCache {
		return mp.powWithOutCache(exp)
	}

	if len(mp.cache[0]) == 0 {
		mp.cache[0] = append(mp.cache[0], mp.base)
		for i := 1; i < mp.cacheLevel-1; i++ {
			mp.cache[i] = append(mp.cache[i], mp.mul(mp.cache[i-1][0], mp.base))
		}
	}

	e := mp.eye(mp.n)
	div := 0
	for exp > 0 {
		if div == len(mp.cache[0]) {
			mp.cache[0] = append(mp.cache[0], mp.mul(mp.cache[mp.cacheLevel-2][div-1], mp.cache[0][div-1]))
			for i := 1; i < mp.cacheLevel-1; i++ {
				mp.cache[i] = append(mp.cache[i], mp.mul(mp.cache[i-1][div], mp.cache[0][div]))
			}
		}

		mod := exp % mp.cacheLevel
		if mod > 0 {
			e = mp.mul(e, mp.cache[mod-1][div])
		}
		exp /= mp.cacheLevel
		div++
	}

	return mp.to2D(e)
}

func (mp *MatPowWithCache) mul(mat1, mat2 []int) []int {
	n := mp.n
	res := make([]int, n*n)
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			for j := 0; j < n; j++ {
				res[i*n+j] = (res[i*n+j] + mat1[i*n+k]*mat2[k*n+j]) % mp.mod
				if res[i*n+j] < 0 {
					res[i*n+j] += mp.mod
				}
			}
		}
	}
	return res
}

func (mp *MatPowWithCache) powWithOutCache(exp int) M {
	e := mp.eye(mp.n)
	b := append(mp.base[:0:0], mp.base...)
	for exp > 0 {
		if exp&1 == 1 {
			e = mp.mul(e, b)
		}
		b = mp.mul(b, b)
		exp >>= 1
	}

	return mp.to2D(e)
}

func (mp *MatPowWithCache) eye(n int) []int {
	res := make([]int, n*n)
	for i := 0; i < n; i++ {
		res[i*n+i] = 1
	}
	return res
}

func (mp *MatPowWithCache) to2D(mat []int) M {
	res := make(M, mp.n)
	for i := 0; i < mp.n; i++ {
		res[i] = make([]int, mp.n)
		copy(res[i], mat[i*mp.n:(i+1)*mp.n])
	}
	return res
}

//
//
//
func NewIdentityMatrix(n int) M {
	res := make(M, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
		res[i][i] = 1
	}
	return res
}

func NewMatrix(row, col int) M {
	res := make(M, row)
	for i := range res {
		res[i] = make([]int, col)
	}
	return res
}

func MatMul(m1, m2 M, mod int) M {
	res := NewMatrix(len(m1), len(m2[0]))
	for i := 0; i < len(m1); i++ {
		for k := 0; k < len(m2); k++ {
			for j := 0; j < len(m2[0]); j++ {
				res[i][j] = (res[i][j] + m1[i][k]*m2[k][j]) % mod
				if res[i][j] < 0 {
					res[i][j] += mod
				}
			}
		}
	}
	return res
}

// matPow/matqpow
func MatPow(m1 M, exp, mod int) M {
	n := len(m1)
	e := NewIdentityMatrix(n)
	b := make(M, n)
	for i := 0; i < n; i++ {
		b[i] = make([]int, n)
		copy(b[i], m1[i])
	}
	for exp > 0 {
		if exp&1 == 1 {
			e = MatMul(e, b, mod)
		}
		b = MatMul(b, b, mod)
		exp >>= 1
	}
	return e
}
