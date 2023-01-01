// gcd卷积/gcdConvolution
// n<=1e6
// 0<=nums[i]<998244353
package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

// https://judge.yosupo.jp/problem/gcd_convolution
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	nums1 := make([]int, n)
	for i := range nums1 {
		fmt.Fscan(in, &nums1[i])
	}
	nums2 := make([]int, n)
	for i := range nums2 {
		fmt.Fscan(in, &nums2[i])
	}

	res := gcdConvolution(nums1, nums2)
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// c[k] = ∑a[i]*b[j] mod MOD, gcd(i,j)=k
func gcdConvolution(nums1, nums2 []int) []int {
	n := len(nums1)
	pf := make([]int, n+1)
	copy1, copy2 := append([]int{0}, nums1...), append([]int{0}, nums2...)

	for i := 2; i < n+1; i++ {
		if pf[i] == 0 {
			for j := n / i; j > 0; j-- {
				pf[j*i] = 1
				copy1[j] = (copy1[j] + copy1[j*i]) % MOD
				copy2[j] = (copy2[j] + copy2[j*i]) % MOD
			}
			pf[i] = 0
		}
	}

	res := make([]int, n+1)
	for i := 0; i < n+1; i++ {
		res[i] = copy1[i] * copy2[i] % MOD
	}

	for i := 2; i < n+1; i++ {
		if pf[i] == 0 {
			for j := 1; j < n/i+1; j++ {
				res[j] = (res[j] - res[j*i] + MOD) % MOD
			}
		}
	}

	return res[1:]
}
