// 快速幂
package quickpow

func main() {
	pow := func(x, n, mod int64) int64 {
		x %= mod
		res := int64(1) % mod
		for ; n > 0; n >>= 1 {
			if n&1 == 1 {
				res = res * x % mod
			}
			x = x * x % mod
		}
		return res
	}
}
