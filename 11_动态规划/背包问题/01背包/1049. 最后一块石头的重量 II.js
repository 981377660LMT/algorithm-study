/**
 * @param {number[]} stones
 * @return {number}
 * 1 <= stones.length <= 30
   1 <= stones[i] <= 100
 * 01背包问题
 * @description
 * 每一回合，从中选出`任意`(不是最大的了)两块石头，然后将它们一起粉碎
 * 最后，最多只会剩下一块 石头。返回此石头 最小的可能重量 。如果没有石头剩下，就返回 0。
 */
function lastStoneWeightII(stones) {
  const sum = stones.reduce((pre, cur) => pre + cur, 0)
  const volumn = sum >> 1
  // dp[i] 表示若干块石头中能否选出一些组成重量和为 i
  const dp = Array(volumn + 1).fill(false)
  dp[0] = true

  for (let i = 0; i < stones.length; i++) {
    for (let j = volumn; j >= 0; j--) {
      j >= stones[i] && (dp[j] = dp[j] || dp[j - stones[i]])
    }
  }

  const maxWeight = dp.lastIndexOf(true)
  return sum - 2 * maxWeight
}
