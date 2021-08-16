/**
 * @param {number[]} nums
 * @return {number}
 */
const findNumberOfLIS = function (nums: number[]): number {
  let res = 0
  // dp[i]：到nums[i]为止的最长递增子序列长度
  const dp = Array(nums.length).fill(1)
  // count[i]：到nums[i]为止的最长递增子序列个数
  const count = Array(nums.length).fill(1)

  for (let i = 1; i < dp.length; i++) {
    for (let j = 0; j < i; j++) {
      if (nums[i] > nums[j]) {
        // 新增长度
        if (dp[j] + 1 > dp[i]) {
          dp[i] = dp[j] + 1
          count[i] = count[j]
          // 不新增长度
        } else if (dp[j] + 1 === dp[i]) {
          count[i] += count[j]
        }
      }
    }
  }

  const max = Math.max.apply(null, dp)
  for (let i = 0; i < nums.length; i++) {
    if (dp[i] === max) res += count[i]
  }
  return res
}

console.log(findNumberOfLIS([1, 3, 5, 4, 7]))
console.log(findNumberOfLIS([2, 2, 2, 2, 2]))
export default 1
