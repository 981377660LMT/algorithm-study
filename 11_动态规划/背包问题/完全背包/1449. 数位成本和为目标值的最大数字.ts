// 属于完全背包问题
/**
 * @param {number[]} cost
 * @param {number} target
 * @return {string}
 * @description 给当前结果添加一个数位（i + 1）的成本为 cost[i]
 * 总成本必须恰好等于 target
 * 添加的数位中没有数字 0 。
 * @summary 完全背包，选出背包，求最后组成的数的最大值
 */
const largestNumber = function (cost: number[], target: number): string {
  const dp = Array(target + 1).fill(-Infinity)
  dp[0] = 0

  // 遍历nums
  for (let i = 9; i >= 0; i--) {
    // 正序遍历容量
    for (let j = cost[i - 1]; j <= target; j++) {
      // 这里是从左向右插入，因此nums是从大到小选的
      dp[j] = Math.max(dp[j], dp[j - cost[i - 1]] * 10 + i)
    }
  }
  console.log(dp)
  return dp[target] === -Infinity ? 0 : dp[target]
}

console.log(largestNumber([7, 6, 5, 5, 5, 6, 8, 7, 8], 12))
