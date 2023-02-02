// st表

package cmnx

import (
	"fmt"
	"math/bits"
)

func main() {
	query := NewSparseTable([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, max)
	fmt.Println(query(0, 9))
	fmt.Println(query(0, 8))
}

// SparseTable 稀疏表: st[j][i] 表示区间 [i, i+2^j-1] 的贡献值
//   query: 查询 [`left`,`right`] 闭区间的贡献值
//     0 <= left <= right < len(nums)
func NewSparseTable(nums []int, op func(int, int) int) (query func(int, int) int) {
	n := len(nums)
	size := bits.Len(uint(n))
	dp := make([][]int, size)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		dp[0][i] = nums[i]
	}

	for i := 1; i < size; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			dp[i][j] = op(dp[i-1][j], dp[i-1][j+(1<<(i-1))])
		}
	}

	query = func(left, right int) int {
		k := bits.Len(uint(right-left+1)) - 1
		return op(dp[k][left], dp[k][right-(1<<k)+1])
	}

	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
