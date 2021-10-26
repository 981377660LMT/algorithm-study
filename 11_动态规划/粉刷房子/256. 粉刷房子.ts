/**
 *
 * @param costs 每个房子粉刷成不同颜色的花费是以一个 n x 3 的正整数矩阵 costs 来表示的
 * 每个房子可以被粉刷成红色、蓝色或者绿色这三种颜色中的一种，
 * 你需要粉刷所有的房子并且使其相邻的两个房子颜色不能相同。
 *
 */
function minCost(costs: number[][]): number {
  const len = costs.length
  const dp = Array.from<number, number[]>({ length: len }, () => Array(3).fill(0))
  dp[0] = costs[0]

  for (let i = 1; i < len; i++) {
    for (let j = 0; j < 3; j++) {
      dp[i][j] =
        Math.min(
          Math.min.apply(null, dp[i - 1].slice(0, j)),
          Math.min.apply(null, dp[i - 1].slice(j + 1))
        ) + costs[i][j]
    }
  }

  return Math.min.apply(null, dp[len - 1])
}

export {}

console.log(
  minCost([
    [3, 5, 3],
    [6, 17, 6],
    [7, 13, 18],
    [9, 10, 18],
  ])
)
