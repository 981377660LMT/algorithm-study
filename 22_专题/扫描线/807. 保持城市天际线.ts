/**
 * @param {number[][]} grid
 * @return {number}
 * grid[i][j]代表位于某处的建筑物的高度
 * 我们被允许增加任何数量（不同建筑物的数量可能不同）的建筑物的高度
 * 从新数组的所有四个方向（即顶部，底部，左侧和右侧）观看的“天际线”必须与原始数组的天际线相同
 * 建筑物高度可以增加的最大总和是多少？
 */
var maxIncreaseKeepingSkyline = function (grid: number[][]): number {
  const m = grid.length
  const n = grid[0].length
  const xMax = Array<number>(m).fill(-Infinity)
  const yMax = Array<number>(n).fill(-Infinity)
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      xMax[i] = Math.max(xMax[i], grid[i][j])
      yMax[j] = Math.max(yMax[j], grid[i][j])
    }
  }

  let res = 0
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      res += Math.min(xMax[i], yMax[j]) - grid[i][j]
    }
  }

  return res
}

console.log(
  maxIncreaseKeepingSkyline([
    [3, 0, 8, 4],
    [2, 4, 5, 7],
    [9, 2, 6, 3],
    [0, 3, 1, 0],
  ])
)

export {}
