/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number}
 *
 * @description 数组中的每个整数前添加 '+' 或 '-'使其等于targte
 * @summary 注意到x-y=target,x+y=Sum,等价于的01背包问题
 * 问题转换为：只使用 +运算符凑出(sum + target) / 2
 */
var findTargetSumWays = function (nums, target) {
  const sum = nums.reduce((pre, cur) => pre + cur, 0)
  const volume = (sum + target) / 2
  if (!Number.isInteger(volume)) return 0
  const dp = Array(volume + 1).fill(0)
  // 我们定义 dp[0] = 1 表示只有当不选取任何元素时，元素之和才为 0，因此只有 1 种方案。
  dp[0] = 1

  for (let i = 0; i < nums.length; i++) {
    const num = nums[i]
    for (let j = volume; j >= num; j--) {
      // 上一次jp[j]的加这一次的dp[j-num]
      dp[j] += dp[j - num]
    }
  }

  return dp[volume]
}

console.log(findTargetSumWays([1, 1, 1, 1, 1], 3))
