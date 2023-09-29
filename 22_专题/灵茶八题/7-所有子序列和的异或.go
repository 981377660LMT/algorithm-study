package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	fmt.Fprintln(out, xorOfSubsetSum(nums))
}

func xorOfSubsetSum(nums []int) int {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	dp := make([]int, sum+1)
	dp[0] = 1

	for _, num := range nums {
		for i := sum; i >= num; i-- {
			dp[i] ^= dp[i-num]
		}
	}

	res := 0
	for i, v := range dp {
		if v != 0 {
			res ^= i
		}
	}
	return res
}
