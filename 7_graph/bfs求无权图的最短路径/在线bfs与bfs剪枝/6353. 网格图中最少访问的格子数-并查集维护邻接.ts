import { onlineBfs } from '../../../22_专题/implicit_graph/OnlineBfs-在线bfs'
import { NextFinder } from '../../../22_专题/implicit_graph/RangeFinder/Finder-线性序列并查集(NextFinder).'

// 给你一个下标从 0 开始的 m x n 整数矩阵 grid 。
// 你一开始的位置在 左上角 格子 (0, 0) 。

// 当你在格子 (i, j) 的时候，你可以移动到以下格子之一：

// 满足 j < k <= grid[i][j] + j 的格子 (i, k) （向右移动），或者
// 满足 i < k <= grid[i][j] + i 的格子 (k, j) （向下移动）。
// 请你返回到达 右下角 格子 (m - 1, n - 1) 需要经过的最少移动格子数，
// 如果无法到达右下角格子，请你返回 -1 。
// https://leetcode.cn/problems/minimum-number-of-visited-cells-in-a-grid/solution/typescript-zai-xian-bfs-jie-jue-bian-shu-9awi/

const INF = 2e15

function minimumVisitedCells(grid: number[][]): number {
  const ROW = grid.length
  const COL = grid[0].length
  const rowVisited: NextFinder[] = Array(ROW).fill(0)
  const colVisited: NextFinder[] = Array(COL).fill(0)
  for (let i = 0; i < ROW; i++) {
    rowVisited[i] = new NextFinder(COL)
  }
  for (let j = 0; j < COL; j++) {
    colVisited[j] = new NextFinder(ROW)
  }

  const dist = onlineBfs(
    ROW * COL,
    0,
    cur => {
      const r = ~~(cur / COL)
      const c = cur % COL
      rowVisited[r].erase(c)
      colVisited[c].erase(r)
    },
    cur => {
      const r = ~~(cur / COL)
      const c = cur % COL
      if (grid[r][c] === 0) return null
      const rightFirst = rowVisited[r].next(c)
      if (rightFirst && rightFirst <= c + grid[r][c]) {
        return r * COL + rightFirst
      }
      const downFirst = colVisited[c].next(r)
      if (downFirst && downFirst <= r + grid[r][c]) {
        return downFirst * COL + c
      }
      return null
    }
  )[0]

  return dist[ROW * COL - 1] < INF ? 1 + dist[ROW * COL - 1] : -1
}
