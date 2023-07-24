package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5}
	res := AllInv(nums)
	for i := 0; i < len(nums); i++ {
		fmt.Println(res[i])
	}

}

const MOD int = 1e9 + 7

// 线性求数组中所有数的逆元。nums 中不能包含 0。
func AllInv(nums []int) []int {
	n := len(nums)
	res := make([]int, n+1)
	res[0] = 1
	for i, v := range nums {
		res[i+1] = res[i] * v % MOD
	}
	inv := Pow(res[n], MOD-2, MOD)
	res = res[:n]
	for i := n - 1; i >= 0; i-- {
		res[i] = res[i] * inv % MOD
		inv = inv * nums[i] % MOD
	}
	return res
}

func Pow(base, exp, mod int) int {
	base %= mod
	res := 1
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}
