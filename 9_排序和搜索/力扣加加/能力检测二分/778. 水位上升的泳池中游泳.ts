/**
 * @param {number[][]} grid
 * @return {number}
 * @description
 * 在一个 N x N 的坐标方格 grid 中，每一个方格的值 grid[i][j] 表示在位置 (i,j) 的平台高度。
 * 当时间为 t 时，此时雨水导致水池中任意位置的水位为 t 。
 * 你可以从一个平台游向四周相邻的任意一个平台，但是前提是此时水位必须同时淹没这两个平台。
 * 你从坐标方格的左上平台 (0，0) 出发。最少耗时多久你才能到达坐标方格的右下平台 (N-1, N-1)？
 */
const swimInWater = function (grid: number[][]): number {
  const n = grid.length

  // 能力检测
  // 这里使用yield 碰到结果就暂停后面的函数运行(大海捞针,yield可以冲出递归调用栈)
  function* dfs(mid: number, x: number, y: number, visited: Set<string>): Generator<boolean> {
    if (x === n - 1 && y === n - 1) yield true
    const next = [
      [x - 1, y],
      [x + 1, y],
      [x, y - 1],
      [x, y + 1],
    ]

    for (const [nextX, nextY] of next) {
      if (
        nextX >= 0 &&
        nextX < n &&
        nextY >= 0 &&
        nextY < n &&
        grid[nextX][nextY] <= mid &&
        !visited.has(`${nextX}#${nextY}`)
      ) {
        visited.add(`${nextX}#${nextY}`)
        yield* dfs(mid, nextX, nextY, visited)
      }
    }

    // visited.delete(`${x}#${y}`)
  }

  // 解空间
  let l = 0
  let r = Math.max.apply(
    null,
    grid.map(row => Math.max.apply(null, row))
  )

  while (l <= r) {
    const mid = Math.floor((l + r) / 2)
    // 最左二分
    if (dfs(mid, 0, 0, new Set<string>(['0#0'])).next().value) {
      r = mid - 1
    } else {
      l = mid + 1
    }
  }

  return l
}

console.log(
  swimInWater([
    [0, 1, 2, 3, 4],
    [24, 23, 22, 21, 5],
    [12, 13, 14, 15, 16],
    [11, 17, 18, 19, 20],
    [10, 9, 8, 7, 6],
  ])
)
// 输出: 16
export {}
