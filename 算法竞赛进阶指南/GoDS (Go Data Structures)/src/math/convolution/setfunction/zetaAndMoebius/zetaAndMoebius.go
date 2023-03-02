// https://nyaannyaan.github.io/library/set-function/zeta-mobius-transform.hpp
// Zeta变换 与 Mobius变换 互为逆变换

// !这种技巧 用 一个状态就可以表示自己(每个mask)的所有子集(超集)的状态了
// SubsetZetaTransform: 求出每个mask的所有子集的贡献 (SoS-DP)
// SupersetZetaTransform: 求出每个mask的所有超集的贡献
// SubsetMoebiusTransform: 从每个mask的所有子集的贡献还原出原数组
// SupersetMoebiusTransform: 从每个mask的所有超集的贡献还原出原数组

package main

import "fmt"

func main() {
	// 0b000 0b001 0b010 0b011 0b100 0b101 0b110 0b111
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7}
	SubsetZetaTransform(nums)
	fmt.Println(nums)
	nums = []int{0, 1, 2, 3, 4, 5, 6, 7}
	SupersetZetaTransform(nums)
	fmt.Println(nums)
}

// f是长度为2^d的数组
//  g[mask] = sum(f[subset])
//  g[0] = f[0]
//  g[1] = f[1]+f[0]
//  g[2] = f[2]+f[0]
//  g[3] = f[3]+f[2]+f[1]+f[0]
//  g[4] = f[4]+f[0]
//  g[5] = f[5]+f[4]+f[1]+f[0]
//  g[6] = f[6]+f[4]+f[2]+f[0]
//  g[7] = f[7]+f[6]+f[5]+f[4]+f[3]+f[2]+f[1]+f[0]
func SubsetZetaTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j|i] += f[j]
			}
		}
	}
}

func SubsetMoebiusTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j|i] -= f[j]
			}
		}
	}
}

// f是长度为2^d的数组
//  g[mask] = sum(f[superset])
//  g[0] = f[0]+f[1]+f[2]+f[3]+f[4]+f[5]+f[6]+f[7]
//  g[1] = f[1]+f[3]+f[5]+f[7]
//  g[2] = f[2]+f[3]+f[6]+f[7]
//  g[3] = f[3]+f[7]
//  g[4] = f[4]+f[5]+f[6]+f[7]
//  g[5] = f[5]+f[7]
//  g[6] = f[6]+f[7]
//  g[7] = f[7]
func SupersetZetaTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j] += f[j|i]
			}
		}
	}
}

func SupersetMoebiusTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j] -= f[j|i]
			}
		}
	}
}
