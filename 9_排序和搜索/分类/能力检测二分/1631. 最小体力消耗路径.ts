/**
 * @param {number[][]} heights  1 <= heights[i][j] <= 106
 * @return {number}
 * 请你返回从左上角走到右下角的最小 体力消耗值
 * 一条路径耗费的 体力值 是路径上相邻格子之间 高度差绝对值 的 最大值 决定的。
 * @summary
 * 关于可不可以，我们可以使用 DFS 来做，由于只需要找到一条满足条件的，或者找到一个不满足的提前退出，
 */
const minimumEffortPath = function (heights: number[][]): number {
  const m = heights.length
  const n = heights[0].length

  // function* dfs(mid: number, x: number, y: number, visited: Set<string>): Generator<boolean> {
  //   if (x === m - 1 && y === n - 1) yield true

  //   const next = [
  //     [x - 1, y],
  //     [x + 1, y],
  //     [x, y - 1],
  //     [x, y + 1],
  //   ]

  //   for (const [nextX, nextY] of next) {
  //     if (
  //       nextX >= 0 &&
  //       nextX < m &&
  //       nextY >= 0 &&
  //       nextY < n &&
  //       !visited.has(`${nextX}#${nextY}`) &&
  //       Math.abs(heights[nextX][nextY] - heights[x][y]) <= mid
  //     ) {
  //       visited.add(`${nextX}#${nextY}`)
  //       yield* dfs(mid, nextX, nextY, visited)
  //     }
  //   }
  // }
  const dfs = (mid: number, x: number, y: number, visited: Set<string>): boolean => {
    if (x === m - 1 && y === n - 1) return true

    const next = [
      [x - 1, y],
      [x + 1, y],
      [x, y - 1],
      [x, y + 1],
    ]

    for (const [nextX, nextY] of next) {
      if (
        nextX >= 0 &&
        nextX < m &&
        nextY >= 0 &&
        nextY < n &&
        !visited.has(`${nextX}#${nextY}`) &&
        Math.abs(heights[nextX][nextY] - heights[x][y]) <= mid
      ) {
        visited.add(`${nextX}#${nextY}`)
        if (dfs(mid, nextX, nextY, visited)) return true
      }
    }

    return false
  }

  let l = 0
  let r = 10 ** 6 - 1
  while (l <= r) {
    const mid = (l + r) >> 1
    if (dfs(mid, 0, 0, new Set('0#0'))) r = mid - 1
    else l = mid + 1
  }

  return l
}

console.log(
  minimumEffortPath([
    [1, 2, 2],
    [3, 8, 2],
    [5, 3, 5],
  ])
)

export {}
