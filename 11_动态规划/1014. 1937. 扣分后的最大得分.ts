// print(Solution().maxPoints(points=[[1, 2, 3], [1, 5, 1], [3, 1, 1]]))
//  输出：9
//  解释：
//  蓝色格子是最优方案选中的格子，坐标分别为 (0, 2)，(1, 1) 和 (2, 0) 。
//  你的总得分增加 3 + 5 + 3 = 11 。
//  但是你的总得分需要扣除 abs(2 - 1) + abs(1 - 0) = 2 。
//  你的最终得分为 11 - 2 = 9 。

// 相邻行 r 和 r + 1 （其中 0 <= r < m - 1），
// 选中坐标为 (r, c1) 和 (r + 1, c2) 的格子，你的总得分 减少 abs(c1 - c2) 。

// 对于所有上一行左边的数，都是加上它的坐标减去自己的坐标；
// 对于上一行所有右边的数，都是减去它的坐标加上自己的坐标；
// 分别求最大值即可。
function maxPoints(points: number[][]): number {
  const [row, col] = [points.length, points[0].length]
  const dp = Array.from<unknown, number[]>({ length: row }, () => Array(col).fill(0))

  for (let c = 0; c < col; c++) {
    dp[0][c] = points[0][c]
  }

  for (let r = 1; r < row; r++) {
    // 看左
    let leftMax = -Infinity
    for (let c = 0; c < col; c++) {
      leftMax = Math.max(leftMax, dp[r - 1][c] + c)
      dp[r][c] = Math.max(dp[r][c], leftMax + points[r][c] - c)
    }

    // 看右
    let rightMax = -Infinity
    for (let c = col - 1; ~c; c--) {
      rightMax = Math.max(rightMax, dp[r - 1][c] - c)
      dp[r][c] = Math.max(dp[r][c], rightMax + points[r][c] + c)
    }
  }

  return Math.max(...dp[row - 1])
}
console.log(
  maxPoints([
    [1, 2, 3],
    [1, 5, 1],
    [3, 1, 1],
  ])
)
