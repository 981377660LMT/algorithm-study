/**
 * @param {number[]} nums
 * @return {number[]}
 * 有点像LIS问题
 */
const largestDivisibleSubset = function (nums: number[]): number[] {
  nums.sort((a, b) => a - b)

  // dp表示以nums[i]结尾的数组的最大整除子集
  const dp = Array.from(nums, v => [v])
  let res: number[] = []
  for (let i = 0; i < nums.length; i++) {
    for (let j = 0; j < i; j++) {
      if (nums[i] % nums[j] === 0 && dp[j].length >= dp[i].length) {
        //  合并ij
        dp[i] = [...dp[j], nums[i]]
      }
    }
    if (dp[i].length > res.length) res = dp[i]
  }

  return res
}

console.log(largestDivisibleSubset([1, 2, 3]))
// 输出：[1,2]
// 解释：[1,3] 也会被视为正确答案。

// answer[i] % answer[j] == 0 ，或
// answer[j] % answer[i] == 0
