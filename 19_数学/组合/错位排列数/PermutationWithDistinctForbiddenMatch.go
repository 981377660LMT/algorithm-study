// 带禁止位置的排列数.
// O(n^2).
package main

func main() {

}

type PermutationWithDistinctForbiddenMatch struct {
	dp  [][]int32
	dp2 [][]int32
}

func NewPermutationWithDistinctForbiddenMatch(n, mod int32) *PermutationWithDistinctForbiddenMatch {
	dp := make([][]int32, n+1)
	for i := int32(0); i <= n; i++ {
		dp[i] = make([]int32, n+1)
	}
	dp[0][0] = 1
	mod64 := int64(mod)
	for i := int32(1); i <= n; i++ {
		for j := int32(0); j <= i; j++ {
			dp[i][j] = int32((int64(dp[i-1][j]) * int64(j+j) % mod64))
			if j+1 <= n {
				dp[i][j] += dp[i-1][j+1]
				if dp[i][j] >= mod {
					dp[i][j] -= mod
				}
			}
			if j > 0 {
				dp[i][j] = int32((int64(dp[i][j]) + int64(dp[i-1][j-1])*int64(j)%mod64*int64(j)) % mod64)
			}
		}
	}
	dp2 := make([][]int32, n+1)
	for i := int32(0); i <= n; i++ {
		dp2[i] = make([]int32, n+1)
	}
	for i := int32(0); i <= n; i++ {
		for j := int32(0); j <= n; j++ {
			if j == 0 {
				dp2[i][j] = dp[i][j]
				continue
			}
			if j > 0 {
				dp2[i][j] = int32((int64(dp2[i][j]) + int64(dp2[i][j-1])*int64(j)) % mod64)
			}
			if i > 0 {
				dp2[i][j] = int32((int64(dp2[i][j]) + int64(dp2[i-1][j])*int64(i)) % mod64)
			}
		}
	}
	return &PermutationWithDistinctForbiddenMatch{dp: dp, dp2: dp2}
}

// i+j个数(i,j<=n)，其中对于k<=i，第k个数不允许为k。
func (p *PermutationWithDistinctForbiddenMatch) Get(i, j int32) int32 {
	return p.dp2[i][j]
}
