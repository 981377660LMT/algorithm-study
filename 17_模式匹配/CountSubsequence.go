package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yuki1493()
}

const MOD int = 1e9 + 7

// https://yukicoder.me/problems/no/1493
// 给定一个长度为n的数组，每次可以将相邻的两个数换成xor
// 问可以得到的数组的个数模1e9+7
func yuki1493() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	for i := 0; i < n-1; i++ {
		nums[i+1] ^= nums[i]
	}
	nums = nums[:n-1]
	res := CountSubSequence(nums, MOD)
	res = (res + 1) % MOD // 空集
	fmt.Fprintln(out, res)
}

func CountSubSequence(seq []int, mod int) int {
	n := len(seq)
	dp := make([]int, n+1)
	dp[0] = 1
	last := make(map[int]int32)
	for i, c := range seq {
		dp[i+1] = 2 * dp[i] % mod
		if v, ok := last[c]; ok {
			dp[i+1] -= dp[v]
			if dp[i+1] < 0 {
				dp[i+1] += mod
			}
		}
		last[c] = int32(i)
	}
	res := (dp[n] - 1) % mod
	if res < 0 {
		res += mod
	}
	return res
}

func CountSubSequenceString(seq string, mod int) int {
	n := len(seq)
	dp := make([]int, n+1)
	dp[0] = 1
	last := make(map[rune]int32)
	for i, c := range seq {
		dp[i+1] = 2 * dp[i] % mod
		if v, ok := last[c]; ok {
			dp[i+1] -= dp[v]
			if dp[i+1] < 0 {
				dp[i+1] += mod
			}
		}
		last[c] = int32(i)
	}
	res := (dp[n] - 1) % mod
	if res < 0 {
		res += mod
	}
	return res
}
