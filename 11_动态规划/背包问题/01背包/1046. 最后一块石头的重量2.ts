/**
 * @param {number[]} stones
 * @return {number}
 * @description 每一回合，从中选出任意两块石头，然后将它们一起粉碎。
   返回此石头 最小的可能重量 。
   @summary 01 背包问题
   本题其实就是选一个子集让和接近一半
 */
const lastStoneWeight = function (stones: number[]) {
  const sum = stones.reduce((pre, cur) => pre + cur, 0)
  const volumn = sum >> 1
  // dp[i][j] 保存前 i 块石头中能否选出一些组成重量和为 j 的可能
  const dp = Array<boolean>(volumn + 1).fill(false)
  dp[0] = true

  for (let i = 0; i < stones.length; i++) {
    for (let j = volumn; j >= 0; j--) {
      dp[j] = dp[j] || dp[j - stones[i]]
    }
  }

  const maxWeight = dp.lastIndexOf(true)
  return sum - 2 * maxWeight
}

console.log(lastStoneWeight([2, 7, 4, 1, 8, 1]))

export default 1
