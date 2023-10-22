package main

func main() {

}

// 二元一次不定方程（线性丢番图方程中的一种）https://en.wikipedia.org/wiki/Diophantine_equation
// exgcd solve equation ax+by=gcd(a,b)
// 特解满足 |x|<=|b|, |y|<=|a|
// https://cp-algorithms.com/algebra/extended-euclid-algorithm.html
func exgcd(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

// 任意非零模数逆元 ax ≡ 1 (mod m)，要求 |gcd(a,m)| = 1     注：不要求 m 为质数
// 返回最小正整数解
// 模板题 https://www.luogu.com.cn/problem/P1082
// https://codeforces.com/problemset/problem/772/C
func modInv(a, m int) int {
	g, x, _ := exgcd(a, m)
	if g != 1 && g != -1 {
		return -1
	}
	res := x % m
	if res < 0 {
		res += m
	}
	return res
}

// ax ≡ b (mod m)，要求 gcd(a,m) | b       注：不要求 m 为质数
// 或者，ax-b 是 m 的倍数，求最小非负整数 x
// 或者，求 ax-km = b 的一个最小非负整数解
// 示例代码 https://codeforces.com/contest/1748/submission/205834351
func modInv2(a, b, m int) int {
	g, x, _ := exgcd(a, m)
	if b%g != 0 {
		return -1
	}
	x *= b / g
	m /= g
	return (x%m + m) % m
}

// a*x + b*y = c 的通解为
// x = (c/g)*x0 + (b/g)*k
// y = (c/g)*y0 - (a/g)*k
// 其中 g = gcd(a,b) 且需要满足 g|c
// x0 和 y0 是 ax+by=g 的一组特解（即 exgcd(a,b) 的返回值）
//
// 为方便讨论，这里要求输入的 a b c 必须为正整数
// 注意若输入超过 1e9 可能要用高精
// 返回：正整数解的个数（无解时为 -1，无正整数解时为 0）
//
//	x 取最小正整数时的解 x1 y1，此时 y1 是最大正整数解
//	y 取最小正整数时的解 x2 y2，此时 x2 是最大正整数解
//
// 相关论文 THE NUMBER OF SOLUTIONS TO ax + by = n http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.376.403
// 相关题目 https://www.luogu.com.cn/problem/P5656
// 使非负解 x+y 尽量小 https://codeforces.com/problemset/problem/1244/C
//
//	最简单的做法就是 min(x1+y1, x2+y2)
//
// 需要转换一下符号 https://atcoder.jp/contests/abc186/tasks/abc186_e
// https://codeforces.com/problemset/problem/1748/D
func solveLinearDiophantineEquations(a, b, c int) (n, x1, y1, x2, y2 int) {
	g, x0, y0 := exgcd(a, b)

	// 无解
	if c%g != 0 {
		n = -1
		return
	}

	a /= g
	b /= g
	c /= g
	x0 *= c
	y0 *= c

	x1 = x0 % b
	if x1 <= 0 { // 若要求的是非负整数解，去掉等号
		x1 += b
	}
	k1 := (x1 - x0) / b
	y1 = y0 - k1*a

	y2 = y0 % a
	if y2 <= 0 { // 若要求的是非负整数解，去掉等号
		y2 += a
	}
	k2 := (y0 - y2) / a
	x2 = x0 + k2*b

	// 无正整数解
	if y1 <= 0 {
		return
	}

	// k 越大 x 越大
	n = k2 - k1 + 1
	return
}

// 费马小定理求质数逆元
// ax ≡ 1 (mod p)
// x^-1 ≡ a^(p-2) (mod p)
// 滑窗 https://codeforces.com/contest/1833/problem/F
func invP(a, p int) int {
	if a <= 0 {
		panic(-1)
	}
	return powM(a, p-2, p)
}

// 有理数取模
// 模板题 https://www.luogu.com.cn/problem/P2613
func divM(a, b, m int) int { return a * invM(b, m) % m }
func divP(a, b, p int) int { return a * invP(b, p) % p }

// // 线性求逆元·其一
// // 求 1^-1, 2^-1, ..., (p−1)^-1 mod p
// // http://blog.miskcoo.com/2014/09/linear-find-all-invert
// // https://www.zhihu.com/question/59033693
// // 模板题 https://www.luogu.com.cn/problem/P3811
// {
// 	const mod = 998244353
// 	const mx int = 1e6
// 	inv := [mx + 1]int{}
// 	inv[1] = 1
// 	for i := 2; i <= mx; i++ {
// 		inv[i] = (mod - mod/i) * inv[mod%i] % mod
// 	}
// }

// // 线性求逆元·其二（离线逆元）
// // 求 a1, a2, ..., an mod p 的逆元
// // 根据 ai^-1 ≡ Πai/ai * (Πai)^-1 (mod p)，求出 Πai 的前缀积和后缀积可以得到 Πai/ai，从而求出 ai^-1 mod p
// // https://zhuanlan.zhihu.com/p/86561431
// // 模板题 https://www.luogu.com.cn/problem/P5431
// calcAllInv := func(a []int, p int) []int {
// 	n := len(a)
// 	pre := make([]int, n+1)
// 	pre[0] = 1
// 	for i, v := range a {
// 		pre[i+1] = pre[i] * v % p
// 	}
// 	invMulAll := invP(pre[n], p)
// 	suf := make([]int, n+1)
// 	suf[n] = 1
// 	for i := len(a) - 1; i > 0; i-- { // i=0 不用求
// 		suf[i] = suf[i+1] * a[i] % p
// 	}
// 	inv := make([]int, n)
// 	for i, pm := range pre[:n] {
// 		inv[i] = pm * suf[i+1] % p * invMulAll % p
// 	}
// 	return inv
// }

// // 模数两两互质的线性同余方程组 - 中国剩余定理 (CRT)
// // x ≡ bi (mod mi), bi 与 mi 互质且 Πmi <= 1e18
// // bi 可以是负数
// // https://blog.csdn.net/synapse7/article/details/9946013
// // https://codeforces.com/blog/entry/61290
// // 模板题 https://www.luogu.com.cn/problem/P1495
// crt := func(b, m []int) (x int) {
// 	M := 1
// 	for _, mi := range m {
// 		M *= mi
// 	}
// 	for i, mi := range m {
// 		Mi := M / mi
// 		_, inv, _ := exgcd(Mi, mi)
// 		x = (x + b[i]*Mi*inv) % M
// 	}
// 	x = (x + M) % M // 调整为非负
// 	return
// }

// // 扩展中国剩余定理 (EXCRT)
// // ai * x ≡ bi (mod mi)
// // 解为 x ≡ b (mod m)
// // 有解时返回 (b, m)，无解时返回 (0, -1)
// // 推导过程见《挑战程序设计竞赛》P292
// // 注意乘法溢出的可能
// // 推荐 https://blog.csdn.net/niiick/article/details/80229217
// // todo 模板题 https://www.luogu.com.cn/problem/P4777 https://www.luogu.com.cn/problem/P4774
// // https://codeforces.com/contest/1500/problem/B
// excrt := func(A, B, M []int) (x, m int) {
// 	m = 1
// 	for i, mi := range M {
// 		a, b := A[i]*m, B[i]-A[i]*x
// 		d := gcd(a, mi)
// 		if b%d != 0 {
// 			return 0, -1
// 		}
// 		t := divM(b/d, a/d, mi/d)
// 		x += m * t
// 		m *= mi / d
// 	}
// 	x = (x%m + m) % m
// 	return
// }
//
// O(n) 预处理阶乘及其逆元，O(1) 求组合数
// {
// 	const mx int = 2e6
// 	F := [mx + 1]int{1}
// 	for i := 1; i <= mx; i++ {
// 		F[i] = F[i-1] * i % mod
// 	}
// 	invF := [...]int{mx: pow(F[mx], mod-2)}
// 	for i := mx; i > 0; i-- {
// 		invF[i-1] = invF[i] * i % mod
// 	}
// 	C := func(n, k int) int {
// 		if k < 0 || k > n {
// 			return 0
// 		}
// 		return F[n] * invF[k] % mod * invF[n-k] % mod
// 	}
// 	P := func(n, k int) int {
// 		if k < 0 || k > n {
// 			return 0
// 		}
// 		return F[n] * invF[n-k] % mod
// 	}

// 	// EXTRA: 卢卡斯定理 https://en.wikipedia.org/wiki/Lucas%27s_theorem
// 	// https://yangty.blog.luogu.org/lucas-theorem-note
// 	// C(n,m)%p = C(n%p,m%p) * C(n/p,m/p) % p
// 	// 注意初始化 F 和 invF 时 mx 取 mod-1
// 	// 推论：n&m==m 时 C(n,m) 为奇数，否则为偶数 https://en.wikipedia.org/wiki/Lucas%27s_theorem#Consequences
// 	// - https://www.zhihu.com/question/64270942
// 	// - https://atcoder.jp/contests/agc043/tasks/agc043_b
// 	// https://www.luogu.com.cn/problem/P3807
// 	// https://www.luogu.com.cn/problem/P7386
// 	var lucas func(int, int) int
// 	lucas = func(n, k int) int {
// 		if k == 0 {
// 			return 1
// 		}
// 		return C(n%mod, k%mod) * lucas(n/mod, k/mod) % mod
// 	}
// // 适用于 n 巨大但 k 或 n-k 较小的情况
// comb := func(n, k int) int {
// 	if k > n-k {
// 		k = n - k
// 	}
// 	a, b := 1, 1
// 	for i := 1; i <= k; i++ {
// 		a = a * n % mod
// 		n--
// 		b = b * i % mod
// 	}
// 	return divP(a, b, mod)
// }
