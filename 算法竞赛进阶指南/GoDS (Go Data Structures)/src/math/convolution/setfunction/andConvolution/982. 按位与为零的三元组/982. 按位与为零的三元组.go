// 982. 按位与为零的三元组

package main

import "math/bits"

const MOD int = 1e18 // 不取模

// 按位与三元组 是由下标 (i, j, k) 组成的三元组，并满足下述全部条件：
// 0<=i,j,k<=n
// nums[i]&nums[j]&nums[k]==0
// 1 <= nums.length <= 1000
// 0 <= nums[i] < 2**16
func countTriplets(nums []int) int {
	max_ := maxs(nums...)
	bit := bits.Len(uint(max_))
	V := make([]int, 1<<bit) // 值域数组
	for _, num := range nums {
		V[num]++
	}
	res1 := AndConvolution(V, V) // 两个数的与为0,1,...,2**16-1的组数
	res2 := AndConvolution(res1, V)
	return res2[0]
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

func AndConvolution(a, b []int) []int {
	a = append(a[:0:0], a...)
	b = append(b[:0:0], b...)
	supersetZetaTransform(a)
	supersetZetaTransform(b)
	for i := range a {
		a[i] = (a[i]*b[i]%MOD + MOD) % MOD
	}
	supersetMoebiusTransform(a)
	return a
}

func supersetZetaTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j] = (f[j] + f[j|i]) % MOD
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
				f[j] = ((f[j]-f[j|i])%MOD + MOD) % MOD
			}
		}
	}
}
