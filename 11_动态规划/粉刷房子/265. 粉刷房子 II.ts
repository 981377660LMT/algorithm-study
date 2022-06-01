/**
 *
 * @param costs 每个房子粉刷成不同颜色的花费是以一个 n x k 的正整数矩阵 costs 来表示的
 * 每个房子可以被粉刷成红色、蓝色或者绿色这三种颜色中的一种，
 * 你需要粉刷所有的房子并且使其相邻的两个房子颜色不能相同。
 *
 */
function minCostII(costs: number[][]): number {
  const [row, col] = [costs.length, costs[0].length]
  let dp = new Uint16Array(costs[0])

  for (let r = 1; r < row; r++) {
    const ndp = new Uint16Array(col).fill(-1) // 相当于65535 (1<<16)-1

    for (let preC = 0; preC < col; preC++) {
      for (let curC = 0; curC < col; curC++) {
        if (preC === curC) continue
        ndp[curC] = Math.min(ndp[curC], dp[preC] + costs[r][curC])
      }
    }

    dp = ndp
  }

  return Math.min(...dp)
}

export {}
