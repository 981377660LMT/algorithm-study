/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number}
 * @description 考虑排列顺序的完全背包问题
 */
var combinationSum4 = function (nums, target) {
  const dp = Array(target + 1).fill(0)
  dp[0] = 1
  for (let i = 1; i <= target; i++) {
    for (let num of nums) {
      if (i - num >= 0) dp[i] += dp[i - num]
    }
  }

  return dp[target]
}

console.log(combinationSum4([1, 2, 3], 4))
