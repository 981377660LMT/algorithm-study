/**
 * @param {number[]} cost
 * @return {number}
 * 请你找出达到楼层顶部的最低花费
 * @summary 到达第i个台阶所花费的最少体力为dp[i]。第一步是要花费体力的，最后一步不用花费体力了
 */
const minCostClimbingStairs = function (cost: number[]): number {
  // 第一步是要花费体力的，最后一步不用
  const dp = Array<number>(cost.length).fill(Infinity)
  dp[0] = cost[0]
  dp[1] = cost[1]
  for (let i = 2; i < cost.length; i++) {
    dp[i] = Math.min(dp[i - 1] + cost[i], dp[i - 2] + cost[i])
  }
  console.log(dp)
  // 最后一步不用花费体力
  return Math.min(dp[cost.length - 1], dp[cost.length - 2])
}

console.log(minCostClimbingStairs([1, 100, 1, 1, 1, 100, 1, 1, 100, 1]))
