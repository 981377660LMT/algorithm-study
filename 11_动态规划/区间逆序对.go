// 区间逆序对

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	nums := make([]int, 100)
	for i := range nums {
		nums[i] = rand.Intn(100)
	}
	inv := RangeInv(nums)
	fmt.Println(inv(0, 10))

	bruteForce := func(start, end int) int {
		res := 0
		for i := start; i < end; i++ {
			for j := i + 1; j < end; j++ {
				if nums[i] > nums[j] {
					res++
				}
			}
		}
		return res
	}

	for i := 0; i < 10; i++ {
		for j := i + 1; j <= 10; j++ {
			if inv(i, j) != bruteForce(i, j) {
				fmt.Println(i, j, inv(i, j), bruteForce(i, j))
				panic("")
			}
		}
	}
	fmt.Println("done")
}

// 区间逆序对.
// 时空复杂度 O(n^2) 预处理, O(1) 查询.
func RangeInv(nums []int) func(start, end int) int {
	n := len(nums)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for left := n; left >= 0; left-- {
		tmp1 := dp[left]
		var tmp2 []int
		if left+1 <= n {
			tmp2 = dp[left+1]
		}
		for right := left; right <= n; right++ {
			if right-left <= 1 {
				continue
			}
			tmp1[right] = tmp2[right] + tmp1[right-1] - tmp2[right-1]
			if nums[left] > nums[right-1] {
				tmp1[right]++
			}
		}
	}

	return func(start, end int) int {
		if start >= end {
			return 0
		}
		return dp[start][end]
	}
}
