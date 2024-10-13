// TODO: Faster https://www.luogu.com.cn/problem/P5824 十二重计数法

package main

import "fmt"

const N int = 1e5
const MOD int = 1e9 + 7

var B = NewBallBoxProblem(N, MOD)

// 写像十二相.
// https://qiita.com/drken/items/f2ea4b58b0d21621bd51.
type BallBoxProblem struct {
	mod  int
	fac  []int
	ifac []int
	inv  []int
}

func NewBallBoxProblem(size, mod int) *BallBoxProblem {
	bb := &BallBoxProblem{mod: mod}
	bb.fac = []int{1}
	bb.ifac = []int{1}
	bb.inv = []int{1}
	bb.expand(size)
	return bb
}

// n个球放入k个盒子的方案数.
//
//	isBallDistinct: 球是否有区别.
//	isBoxDistinct: 盒子是否有区别.
//	atMostOneBallPerBox: 每个盒子最多放一个球.
//	noLimitWithBox: 每个盒子可以放任意个球.
//	atLeastOneBallPerBox: 每个盒子至少放一个球.
func (bb *BallBoxProblem) Query(
	n, k int,
	isBallDistinct, isBoxDistinct bool,
	atMostOneBallPerBox, noLimitWithBox, atLeastOneBallPerBox bool,
) int {
	limitCount := func() int {
		c := 0
		if atMostOneBallPerBox {
			c++
		}
		if noLimitWithBox {
			c++
		}
		if atLeastOneBallPerBox {
			c++
		}
		return c
	}()
	if limitCount != 1 {
		panic("Must have one limit and only one limit with box.")
	}

	if isBallDistinct && isBoxDistinct {
		if atMostOneBallPerBox {
			return bb.solve1(n, k)
		}
		if noLimitWithBox {
			return bb.solve2(n, k)
		}
		if atLeastOneBallPerBox {
			return bb.solve3(n, k)
		}
	}
	if !isBallDistinct && isBoxDistinct {
		if atMostOneBallPerBox {
			return bb.solve4(n, k)
		}
		if noLimitWithBox {
			return bb.solve5(n, k)
		}
		if atLeastOneBallPerBox {
			return bb.solve6(n, k)
		}
	}
	if isBallDistinct && !isBoxDistinct {
		if atMostOneBallPerBox {
			return bb.solve7(n, k)
		}
		if noLimitWithBox {
			return bb.solve8(n, k)
		}
		if atLeastOneBallPerBox {
			return bb.solve9(n, k)
		}
	}
	if !isBallDistinct && !isBoxDistinct {
		if atMostOneBallPerBox {
			return bb.solve10(n, k)
		}
		if noLimitWithBox {
			return bb.solve11(n, k)
		}
		if atLeastOneBallPerBox {
			return bb.solve12(n, k)
		}
	}
	panic("Unreachable")
}

// 有区别的球放入有区别的盒子(每个盒子最多放一个球).
func (bb *BallBoxProblem) solve1(n, k int) int {
	return bb.P(n, k)
}

// 有区别的球放入有区别的盒子(每个盒子可以放任意个球).
func (bb *BallBoxProblem) solve2(n, k int) int {
	return pow(k, n, bb.mod)
}

// 有区别的球放入有区别的盒子(每个盒子至少放一个球).
// O(k*log(n)).
func (bb *BallBoxProblem) solve3(n, k int) int {
	mod := bb.mod
	res := 0
	for i := 0; i <= k; i++ {
		if (k-i)&1 == 1 {
			res -= bb.C(k, i) * pow(i, n, mod)
		} else {
			res += bb.C(k, i) * pow(i, n, mod)
		}
		res %= mod
	}
	if res < 0 {
		res += mod
	}
	return res
}

// 无区别的球放入有区别的盒子(每个盒子最多放一个球).
func (bb *BallBoxProblem) solve4(n, k int) int {
	return bb.C(n, k)
}

// 无区别的球放入有区别的盒子(每个盒子可以放任意个球).
func (bb *BallBoxProblem) solve5(n, k int) int {
	return bb.C(n+k-1, n)
}

// 无区别的球放入有区别的盒子(每个盒子至少放一个球).
func (bb *BallBoxProblem) solve6(n, k int) int {
	return bb.C(n-1, k-1)
}

// 有区别的球放入无区别的盒子(每个盒子最多放一个球).
func (bb *BallBoxProblem) solve7(n, k int) int {
	if n > k {
		return 0
	}
	return 1
}

// 有区别的球放入无区别的盒子(每个盒子可以放任意个球).
// 贝尔数B(n,k).
// O(min(n,k)*logn).
func (bb *BallBoxProblem) solve8(n, k int) int {
	return bb.bell(n, k)
}

// 有区别的球放入无区别的盒子(每个盒子至少放一个球).
// 第二类斯特林数S(n,k).
// O(k*logn).
func (bb *BallBoxProblem) solve9(n, k int) int {
	return bb.stirling2(n, k)
}

// 无区别的球放入无区别的盒子(每个盒子最多放一个球).
func (bb *BallBoxProblem) solve10(n, k int) int {
	if n > k {
		return 0
	}
	return 1
}

// 无区别的球放入无区别的盒子(每个盒子可以放任意个球).
func (bb *BallBoxProblem) solve11(n, k int) int {
	return bb.partition(n, k)
}

// 无区别的球放入无区别的盒子(每个盒子至少放一个球).
// 分割数P(n-k,k).
func (bb *BallBoxProblem) solve12(n, k int) int {
	if n < k {
		return 0
	}
	return bb.partition(n-k, k)
}

func (bb *BallBoxProblem) facF(k int) int {
	bb.expand(k)
	return bb.fac[k]
}

func (bb *BallBoxProblem) ifacF(k int) int {
	bb.expand(k)
	return bb.ifac[k]
}

func (bb *BallBoxProblem) invF(k int) int {
	bb.expand(k)
	return bb.inv[k]
}

func (bb *BallBoxProblem) C(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	return bb.facF(n) * bb.ifacF(k) % bb.mod * bb.ifacF(n-k) % bb.mod
}

func (bb *BallBoxProblem) P(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	return bb.facF(n) * bb.ifacF(n-k) % bb.mod
}

// 可重复选取元素的组合数.
func (bb *BallBoxProblem) H(n, k int) int {
	if n == 0 {
		if k == 0 {
			return 1
		}
		return 0
	}
	return bb.C(n+k-1, k)
}

// n个相同的球放入k个不同的盒子(盒子可放任意个球)的方案数.
func (bb *BallBoxProblem) Put(n, k int) int {
	return bb.C(n+k-1, n)
}

func (bb *BallBoxProblem) partition(n, k int) int {
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, k+1)
	}
	dp[0][0] = 1
	for i := 0; i <= n; i++ {
		for j := 1; j <= k; j++ {
			if i >= j {
				dp[i][j] = dp[i][j-1] + dp[i-j][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	return dp[n][k]
}

func (bb *BallBoxProblem) bell(n, k int) int {
	if k > n {
		k = n
	}
	mod := bb.mod
	jsum := make([]int, k+2)
	for j := 0; j <= k; j++ {
		add := bb.ifacF(j)
		if j&1 == 1 {
			jsum[j+1] = (jsum[j] - add) % mod
		} else {
			jsum[j+1] = (jsum[j] + add) % mod
		}
	}
	res := 0
	for i := 0; i <= k; i++ {
		res += pow(i, n, mod) * bb.ifacF(i) % mod * jsum[k-i+1]
		res %= mod
	}
	if res < 0 {
		res += mod
	}
	return res
}

func (bb *BallBoxProblem) stirling2(n, k int) int {
	mod := bb.mod
	res := 0
	for i := 0; i <= k; i++ {
		if (k-i)&1 == 1 {
			res -= bb.C(k, i) * pow(i, n, mod)
		} else {
			res += bb.C(k, i) * pow(i, n, mod)
		}
		res %= mod
	}
	if res < 0 {
		res += mod
	}
	return res * bb.ifacF(k) % mod
}

func (bb *BallBoxProblem) expand(size int) {
	if len(bb.fac) < size+1 {
		preSize := len(bb.fac)
		diff := size + 1 - preSize
		bb.fac = append(bb.fac, make([]int, diff)...)
		bb.ifac = append(bb.ifac, make([]int, diff)...)
		bb.inv = append(bb.inv, make([]int, diff)...)
		for i := preSize; i <= size; i++ {
			bb.fac[i] = bb.fac[i-1] * i % bb.mod
		}
		bb.ifac[size] = pow(bb.fac[size], bb.mod-2, bb.mod)
		for i := size - 1; i >= preSize; i-- {
			bb.ifac[i] = bb.ifac[i+1] * (i + 1) % bb.mod
		}
		for i := preSize; i <= size; i++ {
			bb.inv[i] = bb.ifac[i] * bb.fac[i-1] % bb.mod
		}
	}
}

func pow(base, exp, mod int) int {
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

// 100450. 安排活动的方案数
// https://leetcode.cn/contest/biweekly-contest-141/problems/find-the-number-of-possible-ways-for-an-event/
func numberOfWays(n int, x int, y int) int {
	res := 0
	for k := 1; k <= x; k++ {
		v1 := B.Query(
			n, k, true, true, false, false, true,
		) * B.C(x, k) % MOD
		v2 := pow(y, k, MOD)
		res = (res + v1*v2) % MOD
	}
	return res
}

func main() {
	fmt.Println(numberOfWays(76, 31, 194))
}
