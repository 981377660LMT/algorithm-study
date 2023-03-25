// st表, 查询区间最大值以及对应的下标(多个最大值时取最小的下标).

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	leaves := []S{{1, 0}, {2, 1}, {3, 2}, {4, 3}, {5, 4}, {6, 5}, {7, 6}, {8, 7}, {9, 8}, {10, 9}}
	query := NewSparseTable(leaves,
		func(s1, s2 S) S {
			if s1.max > s2.max {
				return s1
			}
			if s1.max < s2.max {
				return s2
			}
			return S{max: s1.max, index: min(s1.index, s2.index)}
		})

	fmt.Println(query(0, 9))
	fmt.Println(query(0, 8))
}

// RangeMaxWIthIndex
type S = struct{ max, index int }

// SparseTable 稀疏表: st[j][i] 表示区间 [i, i+2^j-1] 的贡献值
//   query: 查询 [`left`,`right`] 闭区间的贡献值
//     0 <= left <= right < len(nums)
func NewSparseTable(nums []S, op func(S, S) S) (query func(int, int) S) {
	n := len(nums)
	size := bits.Len(uint(n))
	dp := make([][]S, size)
	for i := range dp {
		dp[i] = make([]S, n)
	}

	for i := 0; i < n; i++ {
		dp[0][i] = nums[i]
	}

	for i := 1; i < size; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			dp[i][j] = op(dp[i-1][j], dp[i-1][j+(1<<(i-1))])
		}
	}

	query = func(left, right int) S {
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
