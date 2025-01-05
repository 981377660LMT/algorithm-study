package main

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

// 防溢出的lcm.
// 如果lcm(a,b) > clamp, 返回clamp.
func lcmClamped(a, b int, clamp int) int {
	if a == 0 || b == 0 {
		return 0
	}
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a >= clamp || b >= clamp {
		return clamp
	}
	gcd_ := gcd(a, b)
	a /= gcd_
	if a >= (clamp+b-1)/b {
		return clamp
	}
	return a * b
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
	res := x % m
	if res < 0 {
		res += m
	}
	return res
}

// Chicken McNugget Theorem
// 麦乐鸡定理
// https://artofproblemsolving.com/wiki/index.php/Chicken_McNugget_Theorem
// 给定两个互质正整数，求最大不能用这两个数的线性组合表示的数
// https://leetcode.cn/problems/most-expensive-item-that-can-not-be-bought/description/
func mostExpensiveItem(primeOne int, primeTwo int) int {
	return primeOne*primeTwo - primeOne - primeTwo
}
