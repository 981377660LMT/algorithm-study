/* eslint-disable @typescript-eslint/no-non-null-assertion */
// # 给定一个 m x n 的整数矩阵 grid ，其中 grid[i][j] 等于 0 或 1 。
// # 若 grid[i][j] == 1 ，则表示第 i 个男孩可以邀请第 j 个女孩参加派对。
// # 一个男孩最多可以邀请一个女孩，一个女孩最多可以接受一个男孩的一个邀请。
// # 返回可能的最多邀请的个数。
// 男孩是行,女孩是列

import { useHungarian } from './匈牙利算法'

function maximumInvitations(grid: number[][]): number {
  const [ROW, COL] = [grid.length, grid[0].length]
  const H = useHungarian(ROW, COL) // 男生:行, 女生:列
  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (grid[r][c] === 1) H.addEdge(r, c)
    }
  }
  return H.work().length
}

console.log(
  maximumInvitations([
    [1, 1, 1],
    [1, 0, 1],
    [0, 0, 1]
  ])
)

// # 输出: 3
// # 解释: 按下列方式邀请：
// # - 第 1 个男孩邀请第 3 个女孩。
// # - 第 2 个男孩邀请第 1 个女孩。
// # - 第 3 个男孩未邀请任何人。
// # - 第 4 个男孩邀请第 2 个女孩。
export {}
