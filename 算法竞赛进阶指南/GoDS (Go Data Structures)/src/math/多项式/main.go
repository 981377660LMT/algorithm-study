// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/math_ntt.go
// 多项式库

// Usage:
// polyA:= make(poly,n)

// PolyConv()  // 多个多项式卷积
// poly.Conv   // 两个多项式卷积
// poly.Reverse() // 逆序
// poly.Reversed() // 逆序(返回新多项式)
// poly.Neg() // 取负
// poly.Add() // 加法
// poly.Sub() // 减法
// poly.Mul() // 乘法
// poly.Div() // 除法
// poly.Lsh() // 左移
// poly.Rsh() // 右移
// poly.Derivative() // 求导
// poly.Integral() // 求积分
// poly.Inv() // 求逆
// poly.Sqrt() // 开根
// poly.Ln() // 多项式对数函数
// poly.Exp() // 多项式指数函数
// poly.Pow() // 多项式幂函数
// poly.Mod() // 求模
// poly.Sincos() // 多项式三角函数
// poly.Asin() // 多项式反三角函数

package main

import (
	"math/big"
	"math/bits"
)

func main() {

}

const MOD = 998244353

var omega, omegaInv [31]int64 // 多开一点空间

func init() {
	const g, invG = 3, 332748118
	for i := 1; i < len(omega); i++ {
		omega[i] = _pow(g, (MOD-1)/(1<<i))
		omegaInv[i] = _pow(invG, (MOD-1)/(1<<i))
	}
}

type ntt struct {
	n    int
	invN int64
}

func newNTT(n int) ntt { return ntt{n, _pow(int64(n), MOD-2)} }

// 注：下面 swap 的代码，另一种写法是初始化每个 i 对应的 j https://blog.csdn.net/Flag_z/article/details/99163939
// 由于不是性能瓶颈，实测对性能影响不大
func (t ntt) transform(a, omega []int64) {
	for i, j := 0, 0; i < t.n; i++ {
		if i > j {
			a[i], a[j] = a[j], a[i]
		}
		for l := t.n >> 1; ; l >>= 1 {
			j ^= l
			if j >= l {
				break
			}
		}
	}
	for l, li := 2, 1; l <= t.n; l <<= 1 {
		m := l >> 1
		wn := omega[li]
		li++
		for st := 0; st < t.n; st += l {
			b := a[st:]
			for i, w := 0, int64(1); i < m; i++ {
				d := b[m+i] * w % MOD
				b[m+i] = (b[i] - d + MOD) % MOD
				b[i] = (b[i] + d) % MOD
				w = w * wn % MOD
			}
		}
	}
}

func (t ntt) dft(a []int64) {
	t.transform(a, omega[:])
}

func (t ntt) idft(a []int64) {
	t.transform(a, omegaInv[:])
	for i, v := range a {
		a[i] = v * t.invN % MOD
	}
}

type poly []int64

func (a poly) resize(n int) poly {
	b := make(poly, n)
	copy(b, a)
	return b
}

// 计算 A(x) 和 B(x) 的卷积 (convolution)
// c[i] = ∑a[k]*b[i-k], k=0..i
// 入参出参都是次项从低到高的系数
// 模板题 https://www.luogu.com.cn/problem/P3803 https://www.luogu.com.cn/problem/P1919 https://atcoder.jp/contests/practice2/tasks/practice2_f
func (a poly) Conv(b poly) poly {
	n, m := len(a), len(b)
	limit := 1 << bits.Len(uint(n+m-1))
	A := a.resize(limit)
	B := b.resize(limit)
	t := newNTT(limit)
	t.dft(A)
	t.dft(B)
	for i, v := range A {
		A[i] = v * B[i] % MOD
	}
	t.idft(A)
	return A[:n+m-1]
}

// 计算多个多项式的卷积
// 入参出参都是次项从低到高的系数
func PolyConvNTTs(coefs []poly) poly {
	n := len(coefs)
	if n == 1 {
		return coefs[0]
	}
	return PolyConvNTTs(coefs[:n/2]).Conv(PolyConvNTTs(coefs[n/2:]))
}

func (a poly) Reverse() poly {
	for i, n := 0, len(a); i < n/2; i++ {
		a[i], a[n-1-i] = a[n-1-i], a[i]
	}
	return a
}

func (a poly) Reversed() poly {
	n := len(a)
	b := make(poly, n)
	for i, v := range a {
		b[n-1-i] = v
	}
	return b
}

func (a poly) Neg() poly {
	b := make(poly, len(a))
	for i, v := range a {
		if v > 0 {
			b[i] = MOD - v
		}
	}
	return b
}

func (a poly) Add(b poly) poly {
	c := make(poly, len(a))
	for i, v := range a {
		c[i] = (v + b[i]) % MOD
	}
	return c
}

func (a poly) Sub(b poly) poly {
	c := make(poly, len(a))
	for i, v := range a {
		c[i] = (v - b[i] + MOD) % MOD
	}
	return c
}

func (a poly) Mul(k int64) poly {
	k %= MOD
	b := make(poly, len(a))
	for i, v := range a {
		b[i] = v * k % MOD
	}
	return b
}

func (a poly) Lsh(k int) poly {
	b := make(poly, len(a))
	if k > len(a) {
		return b
	}
	copy(b[k:], a)
	return b
}

func (a poly) Rsh(k int) poly {
	b := make(poly, len(a))
	if k > len(a) {
		return b
	}
	copy(b, a[k:])
	return b
}

func (a poly) Derivative() poly {
	n := len(a)
	d := make(poly, n)
	for i := 1; i < n; i++ {
		d[i-1] = a[i] * int64(i) % MOD
	}
	return d
}

func (a poly) Integral() poly {
	n := len(a)
	s := make(poly, n)
	s[0] = 0 // C
	// 线性求逆元，详见 math.go 中的 initAllInv
	inv := make([]int64, n)
	inv[1] = 1
	for i := 2; i < n; i++ {
		inv[i] = int64(MOD-MOD/i) * inv[MOD%i] % MOD
	}
	for i := 1; i < n; i++ {
		s[i] = a[i-1] * inv[i] % MOD
	}
	return s
}

// 多项式乘法逆 (mod x^n, 下同)
// 参考 https://blog.orzsiyuan.com/archives/Polynomial-Inversion/
// https://oi-wiki.org/math/poly/Inv/
// 模板题 https://www.luogu.com.cn/problem/P4238
func (a poly) Inv() poly {
	n := len(a)
	m := 1 << bits.Len(uint(n))
	A := a.resize(m)
	invA := make(poly, m)
	invA[0] = _pow(A[0], MOD-2)
	for l := 2; l <= m; l <<= 1 {
		ll := l << 1
		b := A[:l].resize(ll)
		iv := invA[:l].resize(ll)
		t := newNTT(ll)
		t.dft(b)
		t.dft(iv)
		for i, v := range iv {
			b[i] = v * (2 - v*b[i]%MOD + MOD) % MOD
		}
		t.idft(b)
		copy(invA, b[:l])
	}
	return invA[:n]
}

// 多项式除法
// https://blog.orzsiyuan.com/archives/Polynomial-Division-and-Modulo/
// https://oi-wiki.org/math/poly/Div-mod/
// 模板题 https://www.luogu.com.cn/problem/P4512
func (a poly) Div(b poly) poly {
	k := len(a) - len(b) + 1
	if k <= 0 {
		return make(poly, 1)
	}
	A := a.Reversed().resize(k)
	B := b.Reversed().resize(k)
	return A.Conv(B.Inv())[:k].Reverse()
}

// 多项式取模
func (a poly) Mod(b poly) poly {
	m := len(b)
	return a[:m-1].Sub(a.Div(b).Conv(b)[:m-1])
}

func (a poly) Divmod(b poly) (quo, rem poly) {
	m := len(b)
	quo = a.Div(b)
	rem = a[:m-1].Sub(quo.Conv(b)[:m-1])
	return
}

// 多项式开根
// 参考 https://blog.orzsiyuan.com/archives/Polynomial-Square-Root/
// https://oi-wiki.org/math/poly/Sqrt/
// 模板题 https://www.luogu.com.cn/problem/P5205
// 模板题（二次剩余）https://www.luogu.com.cn/problem/P5277
func (a poly) Sqrt() poly {
	const inv2 = (MOD + 1) / 2
	n := len(a)
	m := 1 << bits.Len(uint(n))
	A := a.resize(m)
	rt := make(poly, m)
	rt[0] = 1
	if a[0] != 1 {
		rt[0] = new(big.Int).ModSqrt(big.NewInt(a[0]), big.NewInt(MOD)).Int64()
		//if 2*rt[0] > P { // P5277 需要
		//	rt[0] = P - rt[0]
		//}
	}
	for l := 2; l <= m; l <<= 1 {
		ll := l << 1
		b := A[:l].resize(ll)
		r := rt[:l].resize(ll)
		ir := rt[:l].Inv().resize(ll)
		t := newNTT(ll)
		t.dft(b)
		t.dft(r)
		t.dft(ir)
		for i, v := range r {
			b[i] = (b[i] + v*v%MOD) * inv2 % MOD * ir[i] % MOD
		}
		t.idft(b)
		copy(rt, b[:l])
	}
	return rt[:n]
}

// 多项式对数函数
// https://blog.orzsiyuan.com/archives/Polynomial-Natural-Logarithm/
// https://oi-wiki.org/math/poly/Ln-exp/
// 模板题 https://www.luogu.com.cn/problem/P4725
func (a poly) Ln() poly {
	if a[0] != 1 {
		panic(a[0])
	}
	return a.Derivative().Conv(a.Inv())[:len(a)].Integral()
}

// 多项式指数函数
// https://blog.orzsiyuan.com/archives/Polynomial-Exponential/
// https://oi-wiki.org/math/poly/ln-Exp/
// 模板题 https://www.luogu.com.cn/problem/P4726
func (a poly) Exp() poly {
	if a[0] != 0 {
		panic(a[0])
	}
	n := len(a)
	m := 1 << bits.Len(uint(n))
	A := a.resize(m)
	e := make(poly, m)
	e[0] = 1
	for l := 2; l <= m; l <<= 1 {
		b := e[:l].Ln()
		b[0]--
		for i, v := range b {
			b[i] = (A[i] - v + MOD) % MOD
		}
		copy(e, b.Conv(e[:l])[:l])
	}
	return e[:n]
}

// 多项式幂函数
// https://blog.orzsiyuan.com/archives/Polynomial-Power/
// https://oi-wiki.org/math/poly/ln-exp/#_5
// 模板题 https://www.luogu.com.cn/problem/P5245
// 模板题（a[0] != 1）https://www.luogu.com.cn/problem/P5273
func (a poly) Pow(k int64) poly {
	n := len(a)
	if k >= int64(n) && a[0] == 0 {
		return make(poly, n)
	}
	k1 := k % (MOD - 1)
	k %= MOD
	if a[0] == 1 {
		return a.Ln().Mul(k).Exp()
	}
	shift := 0
	for ; shift < n && a[shift] == 0; shift++ {
	}
	if int64(shift)*k >= int64(n) {
		return make(poly, n)
	}
	a = a.Rsh(shift)         // a[0] != 0
	a.Mul(_pow(a[0], MOD-2)) // a[0] == 1
	return a.Ln().Mul(k).Exp().Mul(_pow(a[0], int(k1))).Lsh(shift * int(k))
}

// 多项式三角函数
// 模意义下的单位根 i = w4 = g^((P-1)/4), 其中 g 为 P 的原根
// https://blog.orzsiyuan.com/archives/Polynomial-Trigonometric-Function/
// https://oi-wiki.org/math/poly/tri-func/
// 模板题 https://www.luogu.com.cn/problem/P5264
func (a poly) Sincos() (sin, cos poly) {
	if a[0] != 0 {
		panic(a[0])
	}
	const i = 911660635    // pow(g, (P-1)/4)
	const inv2i = 43291859 // pow(2*i, P-2)
	const inv2 = (MOD + 1) / 2
	e := a.Mul(i).Exp()
	invE := e.Inv()
	sin = e.Sub(invE).Mul(inv2i)
	cos = e.Add(invE).Mul(inv2)
	return
}

func (a poly) Tan() poly {
	sin, cos := a.Sincos()
	return sin.Conv(cos.Inv())
}

// 多项式反三角函数
// https://oi-wiki.org/math/poly/inv-tri-func/
// 模板题 https://www.luogu.com.cn/problem/P5265
func (a poly) Asin() poly {
	if a[0] != 0 {
		panic(a[0])
	}
	n := len(a)
	b := a.Conv(a)[:n].Neg()
	b[0] = 1
	return a.Derivative().Conv(b.Sqrt().Inv())[:n].Integral()
}

func (a poly) Acos() poly {
	return a.Asin().Neg()
}

func (a poly) Atan() poly {
	if a[0] != 0 {
		panic(a[0])
	}
	n := len(a)
	b := a.Conv(a)[:n]
	b[0] = 1
	return a.Derivative().Conv(b.Inv())[:n].Integral()
}

func _pow(x int64, n int) (res int64) {
	res = 1
	for ; n > 0; n >>= 1 {
		if n&1 == 1 {
			res = res * x % MOD
		}
		x = x * x % MOD
	}
	return
}
