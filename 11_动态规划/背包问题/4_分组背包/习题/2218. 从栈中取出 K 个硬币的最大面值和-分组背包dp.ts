// 2218. 从栈中取出K个硬币的最大面值和
// https://leetcode.cn/problems/maximum-value-of-k-coins-from-piles/
// 1 <= k <= sum(piles[i].length) <= 2000
// 1 <= n <= 1000
// 这道题时间复杂度为O(k*sum(piles[i].length))
// 分组背包，每组只能取一个前缀

function maxValueOfCoins(piles: number[][], k: number): number {
  const n = piles.length
  const dp = new Uint32Array(k + 1)
  let curCount = 0
  for (let i = 0; i < n; i++) {
    const pile = piles[i]
    const m = pile.length
    // pile 前缀和
    for (let j = 1; j < m; j++) pile[j] += pile[j - 1]
    curCount = Math.min(curCount + m, k) // 优化：j 从前 i 个栈的大小之和开始枚举（不超过 k）
    for (let j = curCount; j > 0; j--) {
      let max = 0
      for (let w = 0; w < Math.min(m, j); w++) {
        max = Math.max(max, dp[j - w - 1] + pile[w]) // w 从 0 开始，物品体积为 w+1
      }
      dp[j] = Math.max(dp[j], max)
    }
  }

  return dp[k]
}

export {}
