// 快速幂
package quickpow

func Pow(base, exp, mod int64) int64 {
	base %= mod
	res := int64(1)
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}

	return res
}
