// https://leetcode.cn/contest/ccbft-2021fall/problems/lSjqMF/
// # 实验目标要求同学们用导线连接所有「目标插孔」，
// # 即从任意一个「目标插孔」沿导线可以到达其他任意「目标插孔」
// # 一条导线可连接相邻两列的且行间距不超过 1 的两个插孔
// # 每一列插孔中最多使用其中一个插孔（包括「目标插孔」）
// # 若实验目标可达成，请返回使用导线数量最少的连接所有目标插孔的方案数；否则请返回 0。

// # 1 <= row <= 20
// # 3 <= col <= 10^9
// # 1 < position.length <= 1000

// # O(row^3*log(col)*position.length) = 20^3*log(10^9)*1000

package main

import (
	"sort"
)

const MOD int = 1e9 + 7

func electricityExperiment(row int, col int, position [][]int) int {
	n := len(position)
	sort.Slice(position, func(i, j int) bool {
		return position[i][1] < position[j][1]
	})

	T := make([][]int, row) // 转移矩阵
	for i := 0; i < row; i++ {
		T[i] = make([]int, row)
		T[i][i] = 1
		if i > 0 {
			T[i][i-1] = 1
		}
		if i < row-1 {
			T[i][i+1] = 1
		}
	}

	mp := NewMatPow(T, MOD, 4)
	res := 1
	for i := 0; i < n-1; i++ {
		r1, c1, r2, c2 := position[i][0], position[i][1], position[i+1][0], position[i+1][1]
		colDiff := c2 - c1
		pow_ := mp.Pow(colDiff)
		res = res * pow_[r1][r2] % MOD
	}
	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type M = [][]int

type MatPow struct {
	n          int
	mod        int
	base       []int
	cacheLevel int
	useCache   bool
	cache      [][][]int
}

// 带缓存的矩阵快速幂,适合多次查询
//
//	base: 转移矩阵,必须是方阵;
//	mod: 快速幂的模;
//	cacheLevel: 矩阵快速幂的log底数.启用缓存时一般设置为`4`;
//	当调用 pow 次数很多时,可设置为 `16`;
//	小于 `2` 时不启用缓存.
func NewMatPow(base M, mod, cacheLevel int) *MatPow {
	res := &MatPow{}
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
func (mp *MatPow) Pow(exp int) M {
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

func (mp *MatPow) mul(mat1, mat2 []int) []int {
	res := make([]int, mp.n*mp.n)
	for i := 0; i < mp.n; i++ {
		for k := 0; k < mp.n; k++ {
			for j := 0; j < mp.n; j++ {
				res[i*mp.n+j] = (res[i*mp.n+j] + mat1[i*mp.n+k]*mat2[k*mp.n+j]) % mp.mod
				if res[i*mp.n+j] < 0 {
					res[i*mp.n+j] += mp.mod
				}
			}
		}
	}
	return res
}

func (mp *MatPow) powWithOutCache(exp int) M {
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

func (mp *MatPow) eye(n int) []int {
	res := make([]int, n*n)
	for i := 0; i < n; i++ {
		res[i*n+i] = 1
	}
	return res
}

func (mp *MatPow) to2D(mat []int) M {
	res := make(M, mp.n)
	for i := 0; i < mp.n; i++ {
		res[i] = make([]int, mp.n)
		copy(res[i], mat[i*mp.n:(i+1)*mp.n])
	}
	return res
}
