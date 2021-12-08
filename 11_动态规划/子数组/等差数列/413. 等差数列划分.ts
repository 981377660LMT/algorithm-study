/**
 *
 * @param nums  如果一个数列 至少有三个元素 ，并且任意两个相邻元素之差相同，则称该数列为等差数列
 * 返回数组 nums 中所有为等差数组的 子数组 个数。
 */
function numberOfArithmeticSlices(nums: number[]): number {
  // dp[i] 是以 A[i] 为终点的等差数列的个数。
  const n = nums.length
  const dp = Array(n + 1).fill(0)

  for (let i = 2; i < n + 1; i++) {
    if (nums[i] - nums[i - 1] === nums[i - 1] - nums[i - 2]) dp[i] = dp[i - 1] + 1
  }

  return dp.reduce((pre, cur) => pre + cur, 0)
}

console.log(numberOfArithmeticSlices([1, 2, 3, 4]))
// 输入：nums = [1,2,3,4]
// 输出：3
// 解释：nums 中有三个子等差数组：[1, 2, 3]、[2, 3, 4] 和 [1,2,3,4] 自身。
