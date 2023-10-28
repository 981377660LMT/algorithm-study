package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
)

func demo() {
	F := NewFastGcd(1000000)
	fmt.Println(F.Gcd(2, 4))
	fmt.Println(F.Gcd(29, 58))
}

const MOD int = 998244353

// https://www.luogu.com.cn/problem/P5435
// P5435 基于值域预处理的快速 GCD
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums1 := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums1[i])
	}
	nums2 := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums2[i])
	}

	for i := 1; i <= n; i++ {
		res := 0
		for j, now := 1, i; j <= n; j, now = j+1, now*i%MOD {
			res = (res + now*binaryGcd(nums1[i-1], nums2[j-1])) % MOD
		}
		fmt.Fprintln(out, res)
	}
}

// binary binaryGcd
func binaryGcd(a, b int) int {
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

type FastGcd struct {
	max     int
	sqrt    int
	preGcd  [][]int
	fac     [][3]int
	isPrime []bool
	primes  []int
	total   int
}

// O(值域)时间预处理，O(1)时间查询
func NewFastGcd(max int) *FastGcd {
	max++
	res := &FastGcd{}
	sqrt := int(math.Sqrt(float64(max)))

	preGcd := make([][]int, sqrt+1)
	for i := range preGcd {
		preGcd[i] = make([]int, sqrt+1)
	}
	fac := make([][3]int, max+1)
	isPrime := make([]bool, max+1)
	primes := make([]int, max+1)

	res.max = max
	res.sqrt = sqrt
	res.preGcd = preGcd
	res.fac = fac
	res.isPrime = isPrime
	res.primes = primes

	res._build()
	return res
}

// alias of `Query`.
func (fg *FastGcd) Gcd(a, b int) int {
	return fg.Query(a, b)
}

func (fg *FastGcd) Query(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	res := 1
	for i := 0; i < 3; i++ {
		var tmp int
		if f := fg.fac[a][i]; f > fg.sqrt {
			if b%f != 0 {
				tmp = 1
			} else {
				tmp = f
			}
		} else {
			tmp = fg.preGcd[f][b%f]
		}
		b /= tmp
		res *= tmp
	}
	return res
}

func (fg *FastGcd) _build() {
	fg.fac[1][0] = 1
	fg.fac[1][1] = 1
	fg.fac[1][2] = 1
	for i := 2; i <= fg.max; i++ {
		if !fg.isPrime[i] {
			fg.fac[i][0] = 1
			fg.fac[i][1] = 1
			fg.fac[i][2] = i
			fg.total++
			fg.primes[fg.total] = i
		}
		for j := 1; fg.primes[j]*i <= fg.max; j++ {
			tmp := fg.primes[j] * i
			fg.isPrime[tmp] = true
			fg.fac[tmp][0] = fg.fac[i][0] * fg.primes[j]
			fg.fac[tmp][1] = fg.fac[i][1]
			fg.fac[tmp][2] = fg.fac[i][2]
			if fg.fac[tmp][0] > fg.fac[tmp][1] {
				fg.fac[tmp][0] ^= fg.fac[tmp][1]
				fg.fac[tmp][1] ^= fg.fac[tmp][0]
				fg.fac[tmp][0] ^= fg.fac[tmp][1]
			}
			if fg.fac[tmp][1] > fg.fac[tmp][2] {
				fg.fac[tmp][1] ^= fg.fac[tmp][2]
				fg.fac[tmp][2] ^= fg.fac[tmp][1]
				fg.fac[tmp][1] ^= fg.fac[tmp][2]
			}
			if i%fg.primes[j] == 0 {
				break
			}
		}
	}
	for i := 0; i <= fg.sqrt; i++ {
		fg.preGcd[0][i] = i
		fg.preGcd[i][0] = i
	}
	for i := 1; i <= fg.sqrt; i++ {
		for j := 1; j <= i; j++ {
			fg.preGcd[j][i] = fg.preGcd[j][i%j]
			fg.preGcd[i][j] = fg.preGcd[j][i]
		}
	}
}
