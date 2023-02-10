// st表

package cmnx

import (
	"fmt"
	"math/bits"
)

func main() {
	query := NewSparseTable([]S{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, max)
	fmt.Println(query(0, 9))
	fmt.Println(query(0, 8))
}

type S = int

// SparseTable 稀疏表: st[j][i] 表示区间 [i, i+2^j-1] 的贡献值
//   query: 查询 [`left`,`right`] 闭区间的贡献值
//     0 <= left <= right < len(nums)
func NewSparseTable(nums []S, op func(S, S) S) (query func(int, int) S) {
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
