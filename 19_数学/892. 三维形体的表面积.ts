/**
 * @param {number[][]} grid
 * @return {number}
 */
const surfaceArea = function (grid: number[][]): number {
  let res = 0
  const width = grid.length
  const length = grid[0].length

  for (let i = 0; i < width; i++) {
    for (let j = 0; j < length; j++) {
      const height = grid[i][j]
      if (height > 0) {
        // 2个底面 + 所有的正方体都贡献了4个侧表面积
        res += 4 * height + 2
        // 减掉 i 与 i-1 相贴的两份表面积
        i > 0 && (res -= 2 * Math.min(height, grid[i - 1][j]))
        // 减掉 j 与 j-1 相贴的两份表面积
        j > 0 && (res -= 2 * Math.min(height, grid[i][j - 1]))
      }
    }
  }

  return res
}

console.log(
  surfaceArea([
    [1, 2],
    [3, 4],
  ])
)

export default 1
