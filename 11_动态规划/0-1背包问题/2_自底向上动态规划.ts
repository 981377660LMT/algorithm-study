// 建立一个二维的dp数组(放前几个物品(从1开始)，背包0到总容量j)
const dpKnapsack = (
  stuffVolume: number[],
  stuffValue: number[],
  knapsackVolume: number
): number => {
  const dpRow = stuffVolume.length
  if (dpRow === 0) return 0
  const dpCol = knapsackVolume + 1
  const dp = Array.from({ length: dpRow }, () => Array(dpCol).fill(0))

  // 初始化第一行
  for (let j = 0; j < dpCol; j++) {
    dp[0][j] = j >= stuffVolume[0] ? stuffValue[0] : 0
  }

  for (let i = 1; i < dpRow; i++) {
    for (let j = 0; j < dpCol; j++) {
      // 默认不选这个物品
      dp[i][j] = dp[i - 1][j]
      if (j >= stuffVolume[i]) {
        dp[i][j] = Math.max(dp[i][j], stuffValue[i] + dp[i - 1][j - stuffVolume[i]])
      }
    }
  }

  console.table(dp)
  return dp[dpRow - 1][dpCol - 1]
}

console.dir(dpKnapsack([1, 2, 3], [6, 10, 12], 5), { depth: null })

export {}
