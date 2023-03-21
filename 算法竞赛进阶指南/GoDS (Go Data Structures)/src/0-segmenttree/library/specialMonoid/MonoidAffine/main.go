// 仿射变换群

package main

import "fmt"

func main() {
	M := MonoidAffine{}
	e1 := M.e()
	e2 := M.op(e1, M.op(M.add(1), M.add(2)))
	e3 := M.op(e2, M.op(M.mul(3), M.add(4)))
	fmt.Println(e3) // {3 10
}

const MOD int = 1e9 + 7

type E = struct{ mul, add int }

// 仿射变换群
type MonoidAffine struct{}

func (*MonoidAffine) e() E          { return E{1, 0} }
func (*MonoidAffine) op(e1, e2 E) E { return E{e1.mul * e2.mul % MOD, (e1.add*e2.mul + e2.add) % MOD} }
func (*MonoidAffine) inv(e E) E {
	mul, add := e.mul, e.add
	mul = pow(mul, MOD-2, MOD) // modInv of mul
	return E{mul, mul * (MOD - add) % MOD}
}
func (*MonoidAffine) mul(x int) E { return E{(x%MOD + MOD) % MOD, 0} }
func (*MonoidAffine) add(x int) E { return E{1, (x%MOD + MOD) % MOD} }

//
//
func pow(base, exp, mod int) int {
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

func exgcd(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

// 注意模为1时不存在逆元
func modInv(a, mod int) int {
	gcd, x, _ := exgcd(a, mod)
	if gcd != 1 {
		panic(fmt.Sprintf("no inverse element for %d", a))
	}
	return (x%mod + mod) % mod
}
