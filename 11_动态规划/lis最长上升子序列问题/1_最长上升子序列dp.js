// https://leetcode.cn/problems/longest-increasing-subsequence/
// 300. 最长递增子序列

/**
 * @param {number[]} nums
 * @return {number}
 * @description LIS(i)表示以第i个数字为结尾的最长上升子序列的长度
 * 即在范围0到i中选择nums[i]可以获得最长上升子序列的长度\
 * 时间复杂度O(n^2)
 */
function lengthOfLIS(nums) {
  if (nums.length === 0) return 0
  const dp = Array(nums.length).fill(1)

  for (let i = 1; i < nums.length; i++) {
    for (let j = 0; j < i; j++) {
      if (nums[j] < nums[i]) dp[i] = Math.max(dp[i], dp[j] + 1)
    }
  }

  return Math.max(...dp)
}

console.log(lengthOfLIS([10, 9, 2, 5, 3, 7, 101, 18]))
// 输出：4
// 解释：最长递增子序列是 [2,3,7,101]，因此长度为 4 。
