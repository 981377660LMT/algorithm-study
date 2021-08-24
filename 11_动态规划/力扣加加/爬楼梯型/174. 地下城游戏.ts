// 编写一个函数来计算确保骑士能够拯救到公主所需的最低初始健康点数。
// 反向DP
// 状态定义：dp[i][j]表示从[i,j]到终点需要的最小血量，dp[0][0]就是最小初始血量
// 状态转移：1. 如果dungeon[i][j] == 0，那么，dp[i][j] = min(dp[i+1][j], dp[i][j+1])
//          2. 如果dungeon[i][j] < 0，那么，dp[i][j] = min(dp[i+1][j], dp[i][j+1]) - dungeon[i][j]
//          3. 如果dungeon[i][j] > 0，那么，dp[i][j] = max(1, min(dp[i+1][j], dp[i][j+1]) - dungeon[i][j])
// 所以，三种情况可以统一成一种dp[i][j] = max(1, min(dp[i+1][j], dp[i][j+1]) - dungeon[i][j])
// 处理边界：dp[m-1][n-1] = max(1, 1-dungeon[m-1][n-1])，右边和下边，相临的元素只有一个，特殊处理一下。

function calculateMinimumHP(dungeon: number[][]): number {
  const m = dungeon.length
  const n = dungeon[0].length
  const dp = Array.from<number, number[]>({ length: m }, () => Array(n).fill(Infinity))

  dp[m - 1][n - 1] = Math.max(1, 1 - dungeon[m - 1][n - 1])
  for (let i = m - 2; i >= 0; i--) {
    dp[i][n - 1] = Math.max(1, dp[i + 1][n - 1] - dungeon[i][n - 1])
  }
  for (let j = n - 2; j >= 0; j--) {
    dp[m - 1][j] = Math.max(1, dp[m - 1][j + 1] - dungeon[m - 1][j])
  }

  for (let i = m - 2; i >= 0; i--) {
    for (let j = n - 2; j >= 0; j--) {
      dp[i][j] = Math.max(1, Math.min(dp[i + 1][j], dp[i][j + 1]) - dungeon[i][j])
    }
  }

  // console.table(dp)

  return dp[0][0]
}

console.log(
  calculateMinimumHP([
    [-2, -3, 3],
    [-5, -10, 1],
    [10, 30, -5],
  ])
)

export {}
