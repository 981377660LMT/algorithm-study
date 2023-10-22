package main

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
)

func demo() {
	fmt.Println(SegmentedSieve(0, 100))
	fmt.Println(IsPrimeMillerRabin(4))
	fmt.Println(PollardRhoPrimeFactor(100))
	fmt.Println(GetPrimeFactorsFast(1e18 + 9))
}

func Luogu4718() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for i := 0; i < T; i++ {
		var n int
		fmt.Fscan(in, &n)
		if IsPrimeMillerRabin(n) {
			fmt.Fprintln(out, "Prime")
		} else {
			maxPf := PollardRhoPrimeFactor(n)
			fmt.Fprintln(out, maxPf)
		}
	}
}

// 埃氏筛
type EratosthenesSieve struct {
	minPrime []int
}

func NewEratosthenesSieve(maxN int) *EratosthenesSieve {
	minPrime := make([]int, maxN+1)
	for i := range minPrime {
		minPrime[i] = i
	}
	upper := int(math.Sqrt(float64(maxN))) + 1
	for i := 2; i < upper; i++ {
		if minPrime[i] < i {
			continue
		}
		for j := i * i; j <= maxN; j += i {
			if minPrime[j] == j {
				minPrime[j] = i
			}
		}
	}
	return &EratosthenesSieve{minPrime}
}

func (es *EratosthenesSieve) IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	return es.minPrime[n] == n
}

func (es *EratosthenesSieve) GetPrimeFactors(n int) map[int]int {
	res := make(map[int]int)
	for n > 1 {
		m := es.minPrime[n]
		res[m]++
		n /= m
	}
	return res
}

func (es *EratosthenesSieve) GetPrimes() []int {
	res := []int{}
	for i, x := range es.minPrime {
		if i >= 2 && i == x {
			res = append(res, x)
		}
	}
	return res
}

// 返回 n 的所有因子. O(n^0.5).
func GetFactors(n int) []int {
	if n <= 0 {
		return nil
	}
	small := []int{}
	big := []int{}
	upper := int(math.Sqrt(float64(n)))
	for f := 1; f <= upper; f++ {
		if n%f == 0 {
			small = append(small, f)
			big = append(big, n/f)
		}
	}
	if small[len(small)-1] == big[len(big)-1] {
		big = big[:len(big)-1]
	}
	for i, j := 0, len(big)-1; i < j; i, j = i+1, j-1 {
		big[i], big[j] = big[j], big[i]
	}
	res := append(small, big...)
	return res
}

// 返回区间 `[0, upper]` 内所有数的约数.
func GetFactorsOfAll(upper int) [][]int {
	res := make([][]int, upper+1)
	for i := 1; i <= upper; i++ {
		for j := i; j <= upper; j += i {
			res[j] = append(res[j], i)
		}
	}
	return res
}

// O(n^0.5) 判断 n 是否为素数.
func IsPrime(n int) bool {
	if n < 2 || (n >= 4 && n%2 == 0) {
		return false
	}
	upper := int(math.Sqrt(float64(n)))
	for i := 2; i < upper+1; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func IsPrimeMillerRabin(n int) bool {
	return big.NewInt(int64(n)).ProbablyPrime(0)
}

// 返回 n 的所有质数因子，键为质数，值为因子的指数.O(n^0.5).
func GetPrimeFactors(n int) map[int]int {
	res := make(map[int]int)
	upper := int(math.Sqrt(float64(n)))
	for f := 2; f <= upper; f++ {
		count := 0
		for n%f == 0 {
			n /= f
			count++
		}
		if count > 0 {
			res[f] = count
		}
	}
	if n > 1 {
		res[n] = 1
	}
	return res
}

func GetPrimeFactorsFast(n int) map[int]int {
	res := make(map[int]int)
	for n > 1 {
		p := PollardRhoPrimeFactor(n)
		for n%p == 0 {
			n /= p
			res[p]++
		}
	}
	return res
}

// Pollard-Rho 算法求出一个因子 O(n^1/4)
func PollardRhoFactor(n int) int {
	if n == 4 {
		return 2
	}
	if IsPrimeMillerRabin(n) {
		return n
	}

	gcd := func(a, b int) int {
		for a != 0 {
			a, b = b%a, a
		}
		return b
	}

	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	mul := func(a, b int) (res int) {
		for ; b > 0; b >>= 1 {
			if b&1 == 1 {
				res = (res + a) % n
			}
			a = (a + a) % n
		}
		return
	}

	for {
		c := 1 + rand.Intn(n-1)
		f := func(x int) int { return (mul(x, x) + c) % n }
		for t, r := f(0), f(f(0)); t != r; t, r = f(t), f(f(r)) {
			if d := gcd(abs(t-r), n); d > 1 {
				return d
			}
		}
	}
}

// 判断质数+求最大质因子
// 先用 Pollard-Rho 算法求出一个因子，然后递归求最大质因子
// https://zhuanlan.zhihu.com/p/267884783
// https://www.luogu.com.cn/problem/P4718
func PollardRhoPrimeFactor(n int) int {
	if n == 4 {
		return 2
	}
	if IsPrimeMillerRabin(n) {
		return n
	}

	cache := map[int]int{}
	var getPrimeFactor func(int) int
	getPrimeFactor = func(x int) (res int) {
		if cache[x] > 0 {
			return cache[x]
		}
		p := PollardRhoFactor(x)
		if p == x {
			cache[x] = x
			return p
		}
		res = max(getPrimeFactor(p), getPrimeFactor(x/p))
		cache[x] = res
		return
	}

	return getPrimeFactor(n)
}

// 区间质数个数.
// [floor, ceiling]内的质数个数.
// !1<=floor<=ceiling<=1e12,ceiling-floor<=5e5
func CountPrime(floor, ceiling int) int {
	if floor > ceiling {
		return 0
	}
	isPrime := make([]bool, ceiling-floor+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	if floor == 1 {
		isPrime[0] = false
	}
	last := int(math.Sqrt(float64(ceiling)))
	for fac := 2; fac <= last; fac++ {
		start := fac*max(2, (floor+fac-1)/fac) - floor
		for start < len(isPrime) {
			isPrime[start] = false
			start += fac
		}
	}
	res := 0
	for _, v := range isPrime {
		if v {
			res++
		}
	}
	return res
}

// 区间筛/分段筛求 [floor,higher) 中的每个数是否为质数.
// 1<=floor<=higher<=1e12,higher-floor<=5e5.
func SegmentedSieve(floor, higher int) []bool {
	root := 1
	for (root+1)*(root+1) < higher {
		root++
	}

	isPrime := make([]bool, root+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	isPrime[0] = false
	isPrime[1] = false

	res := make([]bool, higher-floor)
	for i := range res {
		res[i] = true
	}
	for i := 0; i < 2-floor; i++ {
		res[i] = false
	}

	for i := 2; i <= root; i++ {
		if isPrime[i] {
			for j := i * i; j <= root; j += i {
				isPrime[j] = false
			}
			for j := max(2, (floor+i-1)/i) * i; j < higher; j += i {
				res[j-floor] = false
			}
		}
	}

	return res
}

// 返回约数个数.`primeFactors`为这个数的所有质数因子分解.
// 如果`primeFactors`为空,返回0.
func CountFactors(primeFactors map[int]int) int {
	if len(primeFactors) == 0 {
		return 0
	}
	res := 1
	for _, count := range primeFactors {
		res *= count + 1
	}
	return res
}

// 返回[0,upper]的所有数的约数个数.
func CountFactorsOfAll(upper int) []int {
	res := make([]int, upper+1)
	for i := 1; i <= upper; i++ {
		for j := i; j <= upper; j += i {
			res[j]++
		}
	}
	return res
}

// 返回约数之和.`primeFactors`为这个数的所有质数因子分解.
// 如果`primeFactors`为空,返回0.
func SumFactors(primeFactors map[int]int) int {
	if len(primeFactors) == 0 {
		return 0
	}
	res := 1
	for prime, count := range primeFactors {
		cur := 1
		for i := 0; i < count; i++ {
			cur = cur*prime + 1
		}
		res *= cur
	}
	return res
}

// 返回[0,upper]的所有数的约数之和.
func SumFactorsOfAll(upper int) []int {
	res := make([]int, upper+1)
	for i := 1; i <= upper; i++ {
		for j := i; j <= upper; j += i {
			res[j] += i
		}
	}
	return res
}

// // n 以内的最多约数个数，以及对应的最小数字
// 	// n <= 1e9
// 	// https://www.luogu.com.cn/problem/P1221
// 	maxDivisorNum := func(n int) (mxc, ans int) {
// 		primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29} // 多取一个质数，让乘法超出 n
// 		var dfs func(int, int, int, int)
// 		dfs = func(i, mxE, c, v int) {
// 			if c > mxc || c == mxc && v < ans {
// 				mxc, ans = c, v
// 			}
// 			for e := 1; e <= mxE; e++ {
// 				v *= primes[i]
// 				if v > n {
// 					break
// 				}
// 				dfs(i+1, e, c*(e+1), v)
// 			}
// 		}
// 		dfs(0, 30, 1, 1)
// 		return
// 	}

// // 在有 mxcLimit 的前提下（限制约数个数），mxc 最大是多少
//
//	maxDivisorNumWithLimit := func(mxcLimit int) (mxc, ans int) {
//		rawAns := sort.Search(1e9, func(n int) bool {
//			c, _ := maxDivisorNum(n + 1)
//			return c > mxcLimit
//		})
//		return maxDivisorNum(rawAns)
//	}

// // Number of odd divisors of n https://oeis.org/A001227
// // a(n) = d(2*n) - d(n)
// // 亦为整数 n 分拆成若干连续整数的方法数
// // Number of partitions of n into consecutive positive integers including the trivial partition of length 1
// // e.g. 9 = 2+3+4 or 4+5 or 9 so a(9)=3
// // 相关题目 LC829 https://leetcode.cn/problems/consecutive-numbers-sum/
// // Kick Start 2021 Round C Alien Generator https://codingcompetitions.withgoogle.com/kickstart/round/0000000000435c44/00000000007ec1cb
// oddDivisorsNum := func(n int) (ans int) {
// 	for i := 1; i*i <= n; i++ {
// 		if n%i == 0 {
// 			if i&1 == 1 {
// 				ans++
// 			}
// 			if i*i < n && n/i&1 == 1 {
// 				ans++
// 			}
// 		}
// 	}
// 	return
// }

// // 因子的中位数（偶数个因子时取小的那个）
// // Lower central (median) divisor of n https://oeis.org/A060775
// // EXTRA: Largest divisor of n <= sqrt(n) https://oeis.org/A033676
// maxSqrtDivisor := func(n int) int {
// 	for d := int(math.Sqrt(float64(n))); ; d-- {
// 		if n%d == 0 {
// 			return d
// 		}
// 	}
// }

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
