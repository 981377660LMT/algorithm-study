// https://hitonanode.github.io/cplib-cpp/number/zeta_moebius_transform.hpp
// !约数/倍数 zeta/μ变换   (约数/倍数的前缀和)
// 整除関係に基づくゼータ変換・メビウス変換．
// またこれらを使用した添字 GCD・添字 LCM 畳み込み．計算量は O(N loglog N)．
// !なお，引数の vector<> について第 0 要素（f[0]）の値は無視される．

package main

func main() {

}

// 倍数的前缀和变换 O(N loglog N)
//  g[n] = f[n] + f[2n] + f[3n] + ... + f[dn]
//  使用例 https://yukicoder.me/submissions/385043
func MultipleZeta(f []int) {
	n := len(f) - 1
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	for p := 2; p <= n; p++ {
		if isPrime[p] {
			for q := p * 2; q <= n; q += p {
				isPrime[q] = false
			}
			for j := n / p; j > 0; j-- {
				f[j] += f[j*p] // op
			}
		}
	}
}

// 倍数的前缀和逆变换 O(N loglog N)
//  使用例 https://yukicoder.me/submissions/385120
func MultipleMoebius(f []int) {
	n := len(f) - 1
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	for p := 2; p <= n; p++ {
		if isPrime[p] {
			for q := p * 2; q <= n; q += p {
				isPrime[q] = false
			}
			for j := 1; j*p <= n; j++ {
				f[j] -= f[j*p] // op
			}
		}
	}
}

// 約数的前缀和变换 O(N loglog N)
//  g[n] = f[fac1] + f[fac2] + f[fac3] + ... + f[facm]
func DivisorZeta(f []int) {
	n := len(f) - 1
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	for p := 2; p <= n; p++ {
		if isPrime[p] {
			for q := p * 2; q <= n; q += p {
				isPrime[q] = false
			}
			for j := 1; j*p <= n; j++ {
				f[j*p] += f[j] // op
			}
		}
	}
}

// 約数的前缀和逆变换 O(N loglog N)
//   codeforces.com/contest/1630/problem/E
func DivisorMoebius(f []int) {
	n := len(f) - 1
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	for p := 2; p <= n; p++ {
		if isPrime[p] {
			for q := p * 2; q <= n; q += p {
				isPrime[q] = false
			}
			for j := n / p; j > 0; j-- {
				f[j*p] -= f[j] // op
			}
		}
	}
}

// gcd卷积
// ret[k] = sum(f[i] * g[j] | gcd(i, j) = k)
func GcdConv(f, g []int) []int {
	n := len(f)
	if n != len(g) {
		panic("len(f) != len(g)")
	}
	MultipleZeta(f)
	MultipleZeta(g)
	for i := range f {
		f[i] *= g[i]
	}
	MultipleMoebius(f)
	return f
}

// lcm卷积
// ret[k] = sum(f[i] * g[j] | lcm(i, j) = k)
func LcmConv(f, g []int) []int {
	n := len(f)
	if n != len(g) {
		panic("len(f) != len(g)")
	}
	DivisorZeta(f)
	DivisorZeta(g)
	for i := range f {
		f[i] *= g[i]
	}
	DivisorMoebius(f)
	return f
}
