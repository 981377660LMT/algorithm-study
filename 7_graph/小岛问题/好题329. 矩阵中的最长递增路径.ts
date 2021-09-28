/**
 * @param {number[][]} matrix 1 <= m, n <= 200
 * @return {number}
 * @summary 需要带记忆化搜索
 */
var longestIncreasingPath = function (matrix: number[][]): number {
  const m = matrix.length
  const n = matrix[0].length
  const memo = Array.from({ length: m }, () => Array(n).fill(0))

  const dfs = (x: number, y: number): number => {
    if (memo[x][y]) return memo[x][y]
    let res = 1
    for (const [dx, dy] of [
      [-1, 0],
      [1, 0],
      [0, -1],
      [0, 1],
    ]) {
      const [nextX, nextY] = [x + dx, y + dy]
      if (
        nextX >= 0 &&
        nextX < m &&
        nextY >= 0 &&
        nextY < n &&
        matrix[nextX][nextY] > matrix[x][y]
      ) {
        res = Math.max(res, 1 + dfs(nextX, nextY))
      }
    }

    memo[x][y] = res
    return res
  }

  let res = 1
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      res = Math.max(res, dfs(i, j))
    }
  }

  return res
}

console.log(
  longestIncreasingPath([
    [9, 9, 4],
    [6, 6, 8],
    [2, 1, 1],
  ])
)
// 4
// 朴素深度优先搜索的时间复杂度过高的原因是进行了大量的重复计算，
// 同一个单元格会被访问多次，每次访问都要重新计算。
// 由于同一个单元格对应的最长递增路径的长度是固定不变的，
// 因此可以使用记忆化的方法进行优化。
// 用矩阵 \textit{memo}memo 作为缓存矩阵，
// 已经计算过的单元格的结果存储到缓存矩阵中
