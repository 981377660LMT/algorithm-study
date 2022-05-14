import { Queue } from 'datastructures-js'

/**
 * @param {number[][]} mat
 * @return {number[][]}
 * 请输出一个大小相同的矩阵，其中每一个格子是 mat 中对应位置元素到最近的 0 的距离。
 * @summary
 * 跟之前求海洋到陆地的最大距离的最小值那个做法一样，多源 BFS
 * 「1162.地图分析」
 */
function updateMatrix(mat: number[][]): number[][] {
  const row = mat.length
  const col = mat[0].length
  const res = Array.from({ length: row }, () => Array(col).fill(Infinity))
  const visited = new Set<number>()

  // 从陆地向海洋多源bfs
  const queue = new Queue<[row: number, col: number]>()
  for (let r = 0; r < row; r++) {
    for (let c = 0; c < col; c++) {
      if (mat[r][c] === 0) {
        res[r][c] = 0
        queue.enqueue([r, c])
        visited.add(r * col + c)
      }
    }
  }

  while (queue.size()) {
    const [row, col] = queue.dequeue()
    ;[
      [row - 1, col],
      [row + 1, col],
      [row, col - 1],
      [row, col + 1],
    ].forEach(([nextRow, nextColumn]) => {
      if (
        nextRow >= 0 &&
        nextRow < row &&
        nextColumn >= 0 &&
        nextColumn < col &&
        !visited.has(nextRow * col + nextColumn)
      ) {
        res[nextRow][nextColumn] = res[row][col] + 1
        queue.enqueue([nextRow, nextColumn])
        visited.add(nextRow * col + nextColumn)
      }
    })
  }

  return res
}

console.log(
  updateMatrix([
    [0, 0, 0],
    [0, 1, 0],
    [0, 0, 0],
  ])
)
// 输出：[[0,0,0],[0,1,0],[0,0,0]]
