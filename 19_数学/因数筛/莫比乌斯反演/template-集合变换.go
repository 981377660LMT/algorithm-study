// 集合上的zeta/mobius变换：高维前缀和(SOSDp)

package main

// 超集的前缀和变换.
//
// c2[v] = c1[s1] + c1[s2] + ... + c1[s_n] (si & v == v)
func SupersetZeta(c1 []int) {
	n := len(c1)
	if n&(n-1) != 0 {
		panic("n must be a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if j&i == 0 {
				c1[j] += c1[j|i]
			}
		}
	}
}

// 超集的前缀和逆变换.
func SupersetMoebius(c2 []int) {
	n := len(c2)
	if n&(n-1) != 0 {
		panic("n must be a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if j&i == 0 {
				c2[j] -= c2[j|i]
			}
		}
	}
}

// 子集的前缀和变换.
//
// c2[v] = c1[s1] + c1[s2] + ... + c1[v] (v & si == si)
func SubsetZeta(c1 []int) {
	n := len(c1)
	if n&(n-1) != 0 {
		panic("n must be a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if j&i == 0 {
				c1[j|i] += c1[j]
			}
		}
	}
}

// 子集的前缀和逆变换.
func SubsetMoebius(c2 []int) {
	n := len(c2)
	if n&(n-1) != 0 {
		panic("n must be a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if j&i == 0 {
				c2[j|i] -= c2[j]
			}
		}
	}
}
