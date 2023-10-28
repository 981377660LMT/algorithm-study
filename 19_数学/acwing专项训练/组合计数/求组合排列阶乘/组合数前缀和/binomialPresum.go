package main

import (
	"math"
	"sort"
)

// https://leetcode.cn/contest/hust_1024_2023/problems/yH1vqC/
// 华科大-04. 美丽字符串
func beautifulString(s string) int {
	curOne, curZero := 0, 0
	todo := make([][2]int, 0, len(s))
	res := 0
	for i, c := range s {
		if c == '1' {
			curOne++
		} else {
			curZero++
		}

		atLeastSelect := i + 1
		if c == '1' {
			atLeastSelect -= curOne
		} else {
			atLeastSelect -= curZero
		}
		res += Pow(2, i, MOD)
		res %= MOD
		todo = append(todo, [2]int{i, atLeastSelect - 1})
	}

	preSum := BinominalPresum(todo, E)
	for i := 0; i < len(preSum); i++ {
		res -= preSum[i]
		res %= MOD
	}

	if res < 0 {
		res += MOD
	}
	return res
}

const MOD int = 998244353

var E = NewEnumeration(2e5+10, MOD)

type IEnumeration interface {
	Inv(v int) int
	C(n, k int) int
}

// 莫队求组合数前缀和.
//
//	queries[i] = [n, k] 表示组合数 C(n, k).
//	返回数组第i项为组合数前缀和 `C(ni,0) + C(ni,1) + ... + C(ni,ki)`.
func BinominalPresum(queries [][2]int, enumeration IEnumeration) []int {
	maxN := 2
	for _, q := range queries {
		maxN = max(maxN, q[0])
	}

	q := len(queries)
	mo := NewMoAlgo(maxN+1, q)
	for _, q := range queries {
		mo.AddQuery(q[1], q[0])
	}

	res := make([]int, q)
	inv2 := enumeration.Inv(2)
	cur := 1
	curN, curK := 0, 0
	addLeft := func(_ int) {
		cur -= enumeration.C(curN, curK)
		cur %= MOD
		curK--
	}
	addRight := func(_ int) {
		cur += cur - enumeration.C(curN, curK)
		cur %= MOD
		curN++
	}
	removeLeft := func(_ int) {
		curK++
		cur += enumeration.C(curN, curK)
		cur %= MOD
	}
	removeRight := func(_ int) {
		curN--
		cur = (cur + enumeration.C(curN, curK)) * inv2 % MOD
	}
	query := func(qid int) {
		if cur < 0 {
			cur += MOD
		}
		res[qid] = cur
	}
	mo.Run(addLeft, addRight, removeLeft, removeRight, query)
	return res
}

type Enumeration struct {
	fac, ifac, inv []int
	mod            int
}

// 模数为质数时的组合数计算.
func NewEnumeration(initSize, mod int) *Enumeration {
	res := &Enumeration{
		fac:  make([]int, 1, initSize+1),
		ifac: make([]int, 1, initSize+1),
		inv:  make([]int, 1, initSize+1),
		mod:  mod,
	}
	res.fac[0] = 1
	res.ifac[0] = 1
	res.inv[0] = 1
	res.expand(initSize)
	return res
}

// 阶乘.
func (e *Enumeration) Fac(k int) int {
	e.expand(k)
	return e.fac[k]
}

// 阶乘逆元.
func (e *Enumeration) Ifac(k int) int {
	e.expand(k)
	return e.ifac[k]
}

// 模逆元.
func (e *Enumeration) Inv(k int) int {
	e.expand(k)
	return e.inv[k]
}

// 组合数.
func (e *Enumeration) C(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	mod := e.mod
	return e.Fac(n) * e.Ifac(k) % mod * e.Ifac(n-k) % mod
}

// 排列数.
func (e *Enumeration) P(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	mod := e.mod
	return e.Fac(n) * e.Ifac(n-k) % mod
}

// 可重复选取元素的组合数.
func (e *Enumeration) H(n, k int) int {
	if n == 0 {
		if k == 0 {
			return 1
		}
		return 0
	}
	return e.C(n+k-1, k)
}

// n个相同的球放入k个不同的盒子(盒子可放任意个球)的方法数.
func (e *Enumeration) Put(n, k int) int {
	return e.C(n+k-1, n)
}

// 卡特兰数.
func (e *Enumeration) Catalan(n int) int {
	return e.C(2*n, n) * e.Inv(n+1) % e.mod
}

// lucas定理求解组合数.适合模数较小的情况.
func (e *Enumeration) Lucas(n, k int) int {
	if k == 0 {
		return 1
	}
	mod := e.mod
	return e.C(n%mod, k%mod) * e.Lucas(n/mod, k/mod) % mod
}

func (e *Enumeration) expand(size int) {
	if upper := e.mod - 1; size > upper {
		size = upper
	}
	if len(e.fac) < size+1 {
		mod := e.mod
		preSize := len(e.fac)
		diff := size + 1 - preSize
		e.fac = append(e.fac, make([]int, diff)...)
		e.ifac = append(e.ifac, make([]int, diff)...)
		e.inv = append(e.inv, make([]int, diff)...)
		for i := preSize; i < size+1; i++ {
			e.fac[i] = e.fac[i-1] * i % mod
		}
		e.ifac[size] = Pow(e.fac[size], mod-2, mod) // !modInv
		for i := size - 1; i >= preSize; i-- {
			e.ifac[i] = e.ifac[i+1] * (i + 1) % mod
		}
		for i := preSize; i < size+1; i++ {
			e.inv[i] = e.ifac[i] * e.fac[i-1] % mod
		}
	}
}

func Pow(base, exp, mod int) int {
	base %= mod
	res := 1 % mod
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}

type MoAlgo struct {
	queryOrder int
	chunkSize  int
	buckets    [][]query
}

type query struct{ qi, left, right int }

func NewMoAlgo(n, q int) *MoAlgo {
	chunkSize := max(1, n/max(1, int(math.Sqrt(float64(q*2/3)))))
	buckets := make([][]query, n/chunkSize+1)
	return &MoAlgo{chunkSize: chunkSize, buckets: buckets}
}

// 添加一个查询，查询范围为`左闭右开区间` [left, right).
//
//	0 <= left <= right <= n
func (mo *MoAlgo) AddQuery(left, right int) {
	index := left / mo.chunkSize
	mo.buckets[index] = append(mo.buckets[index], query{mo.queryOrder, left, right})
	mo.queryOrder++
}

func (mo *MoAlgo) Run(
	addLeft func(index int),
	addRight func(index int),
	removeLeft func(index int),
	removeRight func(index int),
	query func(qid int),
) {
	left, right := 0, 0

	for i, bucket := range mo.buckets {
		if i&1 == 1 {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right < bucket[j].right })
		} else {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right > bucket[j].right })
		}

		for _, q := range bucket {
			// !窗口扩张
			for left > q.left {
				left--
				addLeft(left)
			}
			for right < q.right {
				addRight(right)
				right++
			}

			// !窗口收缩
			for left < q.left {
				removeLeft(left)
				left++
			}
			for right > q.right {
				right--
				removeRight(right)
			}

			query(q.qi)
		}
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
