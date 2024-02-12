// some dp can be solved using this trick
// Problem:
//     Lets have a look at the LIS problem. it has following solution-
//     for(int i = 1; i <= n; i++) {
//       for(int j = 1; j < i; j++) if(a[j] < a[i])
//         dp[i] = max(dp[i], dp[j] + 1)
//     }
//     it can modeled like update-query problem
//     for(int i = 1; i <= n; i++) {
//       dp[i] = query(a[i]);
//       insert({a[i], dp[i]});
//
// Solution:
//     solve(l, r) {
//       m = (l + r) / 2
//       solve(l, m)
//       make ds with a[l..m] and dp[l..m]
//       update dp[m+1..r] using the ds
//       solve(m+1, r)
//     }

package main

func main() {

}

// https://leetcode.cn/problems/longest-increasing-subsequence/
// 最长上升子序列的问题则可以表示成，当满足约束条件：
// 1. i<j
// 2. ai<aj
// 时，可以进行状态转移 dp[j]=max(dp[i]+1,...)
func lengthOfLIS(nums []int) int {

}
