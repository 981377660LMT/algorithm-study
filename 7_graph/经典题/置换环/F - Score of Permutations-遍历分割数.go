// 遍历分割数/划分数/分割方案
// https://maspypy.github.io/library/enumerate/partition.hpp

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int = 998244353
const N int = 100

var fac, ifac [N + 1]int
var modInv [N + 1]int

func init() {
	fac[0] = 1
	ifac[0] = 1
	for i := 1; i < N+1; i++ {
		fac[i] = fac[i-1] * i % MOD
	}
	ifac[N] = Pow(fac[N], MOD-2, MOD)
	for i := N; i > 0; i-- {
		ifac[i-1] = ifac[i] * i % MOD
	}

	modInv[0] = 1
	for i := 1; i < N+1; i++ {
		modInv[i] = Pow(i, MOD-2, MOD)
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

func main() {
	// https://atcoder.jp/contests/abc226/tasks/abc226_f
	// 给定一个长度为 N 的排列 P=(p1,p2,...,pn)
	// 我们定义 P 的分数 S(P) 为所有置换环的长度的lcm.
	// 一共有n!种排列，求出所有排列的分数的k次方的和对998244353取模的结果。
	// n<=50 k<=1e4
	// 分割数与排列的关系
	// !partition: [3,3,2,2,1] => 两个大小为3的环，两个大小为2的环，一个大小为1的环有多少种取法?
	// => 11!/((3*3*2*2*1)*(2!*2!*1!))

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	res := 0
	EnumeratePartition(n, func(partition []int) {
		count := fac[n]
		for _, size := range partition {
			count = count * modInv[size] % MOD // 每个环内部循环位移
		}
		cycle := make(map[int]int)
		for _, size := range partition {
			cycle[size]++
		}
		for _, v := range cycle {
			count = count * ifac[v] % MOD // 大小相同的环重复
		}
		lcm := 1
		for size := range cycle {
			lcm = lcm * size / gcd(size, lcm) % MOD
		}
		res = (res + Pow(lcm, k, MOD)*count) % MOD
	}, -1, -1)

	fmt.Fprintln(out, res)
}

// https://maspypy.github.io/library/enumerate/partition.hpp
//
//	按照字典序降序遍历给定整数 n 的所有可能的正整数分割(和等于 n 的正整数组合).
//	n：要进行整数划分的目标值。
//	・n = 50（204226）：20 ms
//	・n = 60（966467）：80 ms
//	・n = 70（4087968）：360 ms
//	・n = 80（15796476）：1500 ms
//	f ：在每次找到有效的整数划分时被调用，参数是一个表示划分的整数数组。
//	lenLimit：限制整数划分的最大长度。-1表示没有限制。
//	valLimit：限制整数划分中最大整数的值。-1表示没有限制。
func EnumeratePartition(n int, f func(partition []int), lenLimit, valLimit int) {
	var dfs func([]int, int)
	dfs = func(partition []int, sum int) {
		if sum == n {
			f(partition)
			return
		}

		if lenLimit != -1 && len(partition) == lenLimit {
			return
		}

		next := n
		if len(partition) > 0 {
			next = partition[len(partition)-1]
		}
		if valLimit != -1 && next > valLimit {
			next = valLimit
		}
		if tmp := n - sum; next > tmp {
			next = tmp
		}

		partition = append(partition, 0)
		for x := next; x >= 1; x-- {
			partition[len(partition)-1] = x
			dfs(partition, sum+x)
		}
		partition = partition[:len(partition)-1]
	}

	dfs(make([]int, 0, n), 0) // 最多n个1
}

func gcd(a, b int) int {
	// 取绝对值
	x, y := a, b
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	if x == 0 || y == 0 {
		return x + y
	}
	n := bits.TrailingZeros(uint(x))
	m := bits.TrailingZeros(uint(y))
	x >>= n
	y >>= m
	for x != y {
		d := bits.TrailingZeros(uint(x - y))
		f := x > y
		var c int
		if f {
			c = x
		} else {
			c = y
		}
		if !f {
			y = x
		}
		x = (c - y) >> d
	}
	return x << min(n, m)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
