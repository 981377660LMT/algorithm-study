package main

import (
	"fmt"
	"math"
)

func main() {
	// in := bufio.NewReader(os.Stdin)
	// out := bufio.NewWriter(os.Stdout)
	// defer out.Flush()

	// var n int
	// fmt.Fscan(in, &n)
	// for i := 0; i < n; i++ {
	// 	var p int
	// 	fmt.Fscan(in, &p)
	// 	if IsPrimeFast(p) {
	// 		fmt.Fprintln(out, "Yes")
	// 	} else {
	// 		fmt.Fprintln(out, "No")
	// 	}
	// }
	fmt.Println(SegmentedSieve(0, 100))

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

// O(n^0.5) 判断 n 是否为素数.
func isPrime(n int) bool {
	if n < 2 {
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

func Pow(base, exp, mod int) int {
	base %= mod
	res := 1
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
