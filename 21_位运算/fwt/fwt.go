package main

import (
	"math/bits"
	"slices"
)

// 982. 按位与为零的三元组
// https://leetcode.cn/problems/triples-with-bitwise-and-equal-to-zero/description/
// 按位与三元组 是由下标 (i, j, k) 组成的三元组，并满足下述全部条件：
// 0<=i,j,k<=n
// nums[i]&nums[j]&nums[k]==0
// 1 <= nums.length <= 1000
// 0 <= nums[i] < 2**16
func countTriplets(nums []int) int {
	U := 1 << (bits.Len(uint(slices.Max(nums))))
	counter := make([]int, U)
	for _, num := range nums {
		counter[num]++
	}

	FwtAnd(counter, false)
	for i, v := range counter {
		counter[i] *= v * v
	}
	FwtAnd(counter, true)

	return counter[0]
}

// 3514. 不同 XOR 三元组的数目 II
//
// https://leetcode.cn/problems/number-of-unique-xor-triplets-ii/
func uniqueXorTriplets(nums []int) int {
	U := 1 << (bits.Len(uint(slices.Max(nums))))
	counter := make([]int, U)
	for _, num := range nums {
		counter[num]++
	}

	FwtXor(counter, false)
	for i, v := range counter {
		counter[i] *= v * v
	}
	FwtXor(counter, true)

	res := 0
	for _, v := range counter {
		if v > 0 {
			res++
		}
	}
	return res
}

func FwtXor(a []int, inv bool) {
	if !inv {
		n := len(a)
		if n&(n-1) != 0 {
			panic("n is not a power of 2")
		}
		for l, k := 2, 1; l <= n; l, k = l<<1, k<<1 {
			for i := 0; i < n; i += l {
				for j := 0; j < k; j++ {
					a[i+j], a[i+j+k] = add(a[i+j], a[i+j+k]), sub(a[i+j], a[i+j+k])
				}
			}
		}
	} else {
		n := len(a)
		if n&(n-1) != 0 {
			panic("n is not a power of 2")
		}
		for l, k := 2, 1; l <= n; l, k = l<<1, k<<1 {
			for i := 0; i < n; i += l {
				for j := 0; j < k; j++ {
					a[i+j], a[i+j+k] = div2(add(a[i+j], a[i+j+k])), div2(sub(a[i+j], a[i+j+k]))
				}
			}
		}
	}
}

func FwtOr(a []int, inv bool) {
	if !inv {
		subsetZetaTransform(a)
	} else {
		subsetMoebiusTransform(a)
	}
}

func FwtAnd(a []int, inv bool) {
	if !inv {
		supersetZetaTransform(a)
	} else {
		supersetMoebiusTransform(a)
	}
}

func supersetZetaTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j] = add(f[j], f[j|i])
			}
		}
	}
}

func supersetMoebiusTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j] = sub(f[j], f[j|i])
			}
		}
	}
}

func subsetZetaTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j|i] = add(f[j|i], f[j])
			}
		}
	}
}

func subsetMoebiusTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j|i] = sub(f[j|i], f[j])
			}
		}
	}
}

// Shallow clone.
func clone[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
}

func add(a, b int) int {
	return (a + b)
}

func sub(a, b int) int {
	return a - b
}

func div2(a int) int {
	return a >> 1
}
