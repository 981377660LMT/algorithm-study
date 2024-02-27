package main

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"math/rand"
	"os"
	"sort"
)

func main() {
	fmt.Println(SegmentedSieve(0, 100))
	fmt.Println(IsPrimeMillerRabin(4))
	fmt.Println(PollardRhoPrimeFactor(100))
	fmt.Println(GetPrimeFactorsBig(1e18 + 9))
	EnumerateFactors(100, func(factor int) bool {
		fmt.Println(factor)
		return false
	})
	fmt.Println(GetFactors(10), SumFactors2(10))

	for i := 1; i <= 1000; i++ {
		if SumFactors(GetPrimeFactors(i)) != SumFactors2(i) {
			panic("err")
		}
		if CountFactors(GetPrimeFactors(i)) != CountFactors2(i) {
			panic("err")
		}
	}

	fmt.Println(MaxDivisorNum(100))
	for i := 0; i <= 10; i++ {
		fmt.Println(MaxDivisorNumWithLimit(i))
	}

	fmt.Println(MaxDivisorNum(5e17))
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
			if f*f < n {
				big = append(big, n/f)
			}
		}
	}
	for i := len(big) - 1; i >= 0; i-- {
		small = append(small, big[i])
	}
	return small
}

// 空间复杂度为O(1)的枚举因子.枚举顺序为从小到大.
func EnumerateFactors(n int, f func(factor int) (shouldBreak bool)) {
	if n <= 0 {
		return
	}
	i := 1
	upper := int(math.Sqrt(float64(n)))
	for ; i <= upper; i++ {
		if n%i == 0 {
			if f(i) {
				return
			}
		}
	}
	i--
	if i*i == n {
		i--
	}
	for ; i > 0; i-- {
		if n%i == 0 {
			if f(n / i) {
				return
			}
		}
	}
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
	if n <= 1 {
		return res
	}

	count2 := 0
	for n%2 == 0 {
		n /= 2
		count2++
	}
	if count2 > 0 {
		res[2] = count2
	}

	for i := 3; i*i <= n; i += 2 {
		count := 0
		for n%i == 0 {
			n /= i
			count++
		}
		if count > 0 {
			res[i] = count
		}
	}

	if n > 1 {
		res[n] = 1
	}

	return res
}

func GetPrimeFactorsBig(n int) map[int]int {
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
// 如果`primeFactors`为空,返回1.
func CountFactors(primeFactors map[int]int) int {
	res := 1
	for _, count := range primeFactors {
		res *= count + 1
	}
	return res
}

func CountFactors2(x int) int {
	if x <= 0 {
		return 0
	}
	res := 1
	if x&1 == 0 {
		e := 2
		x >>= 1
		for x&1 == 0 {
			x >>= 1
			e++
		}
		res *= e
	}
	for f := 3; f*f <= x; f += 2 {
		if x%f == 0 {
			e := 2
			x /= f
			for x%f == 0 {
				x /= f
				e++
			}
			res *= e
		}
	}
	if x > 1 {
		res *= 2
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
// 如果`primeFactors`为空,返回1.
func SumFactors(primeFactors map[int]int) int {
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

func SumFactors2(n int) int {
	if n <= 0 {
		return 0
	}
	res := 1
	if n&1 == 0 {
		cur := 1
		for n&1 == 0 {
			n >>= 1
			cur = cur*2 + 1
		}
		res *= cur
	}
	for f := 3; f*f <= n; f += 2 {
		if n%f == 0 {
			cur := 1
			for n%f == 0 {
				n /= f
				cur = cur*f + 1
			}
			res *= cur
		}
	}
	if n > 1 {
		res *= n + 1
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

// n 以内的最多约数个数，以及对应的最小数字
// n <= 1e9
// https://www.luogu.com.cn/problem/P1221
func MaxDivisorNum(n int) (count, res int) {
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	var dfs func(i, maxExp, curCount, curRes int)
	dfs = func(i, maxExp, curCount, curRes int) {
		if curCount > count || (curCount == count && curRes < res) {
			count, res = curCount, curRes
		}
		for e := 1; e <= maxExp; e++ {
			curRes *= primes[i]
			if curRes > n {
				break
			}
			dfs(i+1, e, curCount*(e+1), curRes)
		}
	}
	dfs(0, bits.Len(uint(n)), 1, 1)
	return
}

// 在有 最大约数个数限制 的前提下，maxCount 最大是多少，以及对应的最小数字.
func MaxDivisorNumWithLimit(maxCount int) (count, res int) {
	if maxCount == 0 {
		return
	}
	num := sort.Search(1e9, func(n int) bool {
		c, _ := MaxDivisorNum(n + 1)
		return c > maxCount
	})
	return MaxDivisorNum(num)
}

// [min,max]以内的最多约数个数，以及对应的最小数字.
// 1<=min<=max<=1e9
// dfs+剪枝
// https://www.luogu.com.cn/problem/P1221
func MaxDivisorNumInInterval(min, max int) (count, res int) {
	if max-min <= 100000 {
		for i := min; i <= max; i++ {
			curCount := CountFactors2(i)
			if curCount > count {
				count, res = curCount, i
			}
		}
		return
	}
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	var dfs func(i, maxExp, curCount, curRes int)
	dfs = func(i, maxExp, curCount, curRes int) {
		if curRes >= min && (curCount > count || (curCount == count && curRes < res)) {
			count, res = curCount, curRes
		}
		for e := 1; e <= maxExp; e++ {
			curRes *= primes[i]
			if curRes > max {
				break
			}
			dfs(i+1, e, curCount*(e+1), curRes)
		}
	}
	dfs(0, bits.Len(uint(max)), 1, 1)
	return
}

// 整数n拆分成若干连续整数的方法数/奇约数个数
// Number of odd divisors of n https://oeis.org/A001227
// e.g. 9 = 2+3+4 or 4+5 or 9 so a(9)=3
// LC829 连续整数求和
// https://leetcode.cn/problems/consecutive-numbers-sum/
func OddDivisorsNum(n int) int {
	res := 0
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			if i&1 == 1 {
				res++
			}
			if i*i < n && n/i&1 == 1 {
				res++
			}
		}
	}
	return res
}

// 因子的中位数（偶数个因子时取小的那个）
// Lower central (median) divisor of n https://oeis.org/A060775
// EXTRA: Largest divisor of n <= sqrt(n) https://oeis.org/A033676
func MedianDivisor(n int) int {
	for d := int(math.Sqrt(float64(n))); ; d-- {
		if n%d == 0 {
			return d
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
