package main

// 求组合数.适用于 n 巨大但 k 或 n-k 较小的情况.
func BigModSmallKComb(n, k, mod int) (res int) {
	if k > n-k {
		k = n - k
	}
	a, b := 1, 1
	for i := 1; i <= k; i++ {
		a = a * n % mod
		n--
		b = b * i % mod
	}
	return a * pow(b, mod-2, mod) % mod
}

func pow(x, n, mod int) (res int) {
	res = 1
	for ; n > 0; n >>= 1 {
		if n&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
	}
	return
}
