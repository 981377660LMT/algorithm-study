/**
 * @param {number[]} stones
 * @return {number}
 * 1 <= stones.length <= 30
   1 <= stones[i] <= 100
 * 01背包问题 最后一块石头的重量
 * @description
 * 每一回合，从中选出任意两块石头，然后将它们一起粉碎
 * 最后，最多只会剩下一块 石头。返回此石头 最小的可能重量 。如果没有石头剩下，就返回 0。
 */
function lastStoneWeightII(stones: number[]): number {
  const n = stones.length
  const sum = stones.reduce((pre, cur) => pre + cur, 0)
  const target = sum >>> 1
  let dp = new Uint8Array(target + 1) // !dp[i] 表示若干块石头中能否选出一些组成重量和为 i
  dp[0] = 1

  for (let i = 0; i < n; i++) {
    const ndp = dp.slice()
    for (let pre = 0; pre + stones[i] <= target; pre++) {
      ndp[pre + stones[i]] |= dp[pre]
    }
    dp = ndp
  }

  const maxHalf = dp.lastIndexOf(1)
  return sum - 2 * maxHalf
}

export {}

if (require.main === module) {
  console.log(lastStoneWeightII([2, 7, 4, 1, 8, 1]))
  console.log(lastStoneWeightII([31, 26, 33, 21, 40]))
}
