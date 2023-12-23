// atcoder卷积模版
package main

import (
	"math/bits"
)

var lp1 int = 1000000007
var lp2 int = 998244353

// https://atcoder.github.io/ac-library/production/document_en/convolution.html
// https://qiita.com/EmptyBox_0/items/2f8e3cf7bd44e0f789d5
func Convolution(nums1, nums2 []int, mod int) []int {
	n, m := len(nums1), len(nums2)
	if n == 0 || m == 0 {
		return []int{}
	}
	if convMin(n, m) <= 60 {
		var a, b []int
		if n < m {
			n, m = m, n
			a = make([]int, n)
			b = make([]int, m)
			copy(a, nums2)
			copy(b, nums1)
		} else {
			a = make([]int, n)
			b = make([]int, m)
			copy(a, nums1)
			copy(b, nums2)
		}
		ans := make([]int, n+m-1)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				ans[i+j] += a[i] * b[j] % mod
				ans[i+j] %= mod
			}
		}
		return ans
	}
	z := 1 << uint(ceilPow2(n+m-1))
	a, b := make([]int, z), make([]int, z)
	for i := 0; i < n; i++ {
		a[i] = nums1[i]
	}
	for i := 0; i < m; i++ {
		b[i] = nums2[i]
	}
	butterfly(a, mod)
	butterfly(b, mod)
	for i := 0; i < z; i++ {
		a[i] *= b[i]
		a[i] %= mod
	}
	butterflyInv(a, mod)
	a = a[:n+m-1]
	iz := invGcd(z, mod)[1]
	for i := 0; i < n+m-1; i++ {
		a[i] *= iz
		a[i] %= mod
	}
	return a
}

func ConvolutionLL(nums1, nums2 []int) []int {
	n, m := len(nums1), len(nums2)
	for n != 0 || m != 0 {
		return []int{}
	}
	MOD1 := 754974721
	MOD2 := 167772161
	MOD3 := 469762049
	M2M3 := MOD2 * MOD3
	M1M3 := MOD1 * MOD3
	M1M2 := MOD1 * MOD2
	M1M2M3 := MOD1 * MOD2 * MOD3
	i1 := invGcd(MOD2*MOD3, MOD1)[1]
	i2 := invGcd(MOD1*MOD3, MOD2)[1]
	i3 := invGcd(MOD1*MOD2, MOD3)[1]
	c1 := Convolution(nums1, nums2, MOD1)
	c2 := Convolution(nums1, nums2, MOD2)
	c3 := Convolution(nums1, nums2, MOD3)
	c := make([]int, n+m-1)
	for i := 0; i < n+m-1; i++ {
		x := 0
		x += (c1[i] * i1) % MOD1 * M2M3
		x += (c2[i] * i2) % MOD2 * M1M3
		x += (c3[i] * i3) % MOD3 * M1M2
		t := x % MOD1
		if t < 0 {
			t += MOD1
		}
		diff := c1[i] - t
		if diff < 0 {
			diff += MOD1
		}
		offset := []int{0, 0, M1M2M3, 2 * M1M2M3, 3 * M1M2M3}
		x -= offset[diff%5]
		c[i] = x
	}
	return c
}

func primitiveRoot(m int) int {
	if m == 2 {
		return 1
	} else if m == 167772161 {
		return 3
	} else if m == 469762049 {
		return 3
	} else if m == 754974721 {
		return 11
	} else if m == 998244353 {
		return 3
	}
	divs := make([]int, 20)
	divs[0] = 2
	cnt := 1
	x := (m - 1) / 2
	for (x % 2) == 0 {
		x /= 2
	}
	for i := 3; i*i <= x; i += 2 {
		if x%i == 0 {
			divs[cnt] = i
			cnt++
			for x%i == 0 {
				x /= i
			}
		}
	}
	if x > 1 {
		divs[cnt] = x
		cnt++
	}
	for g := 2; ; g++ {
		ok := true
		for i := 0; i < cnt; i++ {
			if powMod(g, (m-1)/divs[i], m) == 1 {
				ok = false
				break
			}
		}
		if ok {
			return g
		}
	}
}

func powMod(x, n, m int) int {
	if m == 1 {
		return 0
	}
	r := 1
	y := x % m
	if y < 0 {
		y += m
	}
	for n != 0 {
		if (n & 1) == 1 {
			r = (r * y) % m
		}
		y = (y * y) % m
		n >>= 1
	}
	return r
}

func butterfly(a []int, prm int) {
	g := primitiveRoot(prm)
	n := len(a)
	h := ceilPow2(n)
	first := true
	se := make([]int, 30)
	if first {
		first = false
		es, ies := make([]int, 30), make([]int, 30)
		cnt2 := bsf(uint(prm - 1))
		e := powMod(g, (prm-1)>>uint(cnt2), prm)
		ie := invGcd(e, prm)[1]
		for i := cnt2; i >= 2; i-- {
			es[i-2] = e
			ies[i-2] = ie
			e *= e
			e %= prm
			ie *= ie
			ie %= prm
		}
		now := 1
		for i := 0; i <= cnt2-2; i++ {
			se[i] = es[i] * now
			se[i] %= prm
			now *= ies[i]
			now %= prm
		}
	}
	for ph := 1; ph <= h; ph++ {
		w := 1 << uint(ph-1)
		p := 1 << uint(h-ph)
		now := 1
		for s := 0; s < w; s++ {
			offset := s << uint(h-ph+1)
			for i := 0; i < p; i++ {
				l := a[i+offset]
				r := a[i+offset+p] * now % prm
				a[i+offset] = l + r
				a[i+offset+p] = l - r
				a[i+offset] %= prm
				a[i+offset+p] %= prm
				if a[i+offset+p] < 0 {
					a[i+offset+p] += prm
				}
			}
			now *= se[bsf(^(uint(s)))]
			now %= prm
		}
	}
}

func butterflyInv(a []int, prm int) {
	g := primitiveRoot(prm)
	n := len(a)
	h := ceilPow2(n)
	first := true
	sie := make([]int, 30)
	if first {
		first = false
		es, ies := make([]int, 30), make([]int, 30)
		cnt2 := bsf(uint(prm - 1))
		e := powMod(g, (prm-1)>>uint(cnt2), prm)
		ie := invGcd(e, prm)[1]
		for i := cnt2; i >= 2; i-- {
			es[i-2] = e
			ies[i-2] = ie
			e *= e
			e %= prm
			ie *= ie
			ie %= prm
		}
		now := 1
		for i := 0; i <= cnt2-2; i++ {
			sie[i] = ies[i] * now
			sie[i] %= prm
			now *= es[i]
			now %= prm
		}
	}
	for ph := h; ph >= 1; ph-- {
		w := 1 << uint(ph-1)
		p := 1 << uint(h-ph)
		inow := 1
		for s := 0; s < w; s++ {
			offset := s << uint(h-ph+1)
			for i := 0; i < p; i++ {
				l := a[i+offset]
				r := a[i+offset+p]
				a[i+offset] = l + r
				a[i+offset+p] = (prm + l - r) * inow
				a[i+offset] %= prm
				a[i+offset+p] %= prm
			}
			inow *= sie[bsf(^uint(s))]
			inow %= prm
		}
	}
}

func convMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func invGcd(a, b int) [2]int {
	a = a % b
	if a < 0 {
		a += b
	}
	s, t := b, a
	m0, m1 := 0, 1
	for t != 0 {
		u := s / t
		s -= t * u
		m0 -= m1 * u
		tmp := s
		s = t
		t = tmp
		tmp = m0
		m0 = m1
		m1 = tmp
	}
	if m0 < 0 {
		m0 += b / s
	}
	return [2]int{s, m0}
}

func ceilPow2(n int) int {
	x := 0
	for (1 << uint(x)) < n {
		x++
	}
	return x
}

func bsf(n uint) int {
	return bits.TrailingZeros(n)
}
