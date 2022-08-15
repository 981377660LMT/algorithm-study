// 相邻两个1组成一条边，每条边都要去掉一个端点，
// !其实是找最小点覆盖，即求二分图的最大匹配，跑匈牙利算法
// 2123. 使矩阵中的 1 互不相邻的最小操作数

import { useHungarian } from './匈牙利算法'

const DIR4 = [
  [-1, 0],
  [0, 1],
  [1, 0],
  [0, -1]
]

function minimumOperations(grid: number[][]): number {
  const [ROW, COL] = [grid.length, grid[0].length]
  const H = useHungarian(ROW * COL, ROW * COL) // 男生:偶数格子, 女生:奇数格子

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (grid[r][c] === 0 || (r + c) & 1) continue

      const cur = r * COL + c
      for (const [dr, dc] of DIR4) {
        const nr = r + dr
        const nc = c + dc
        if (nr < 0 || nr >= ROW || nc < 0 || nc >= COL) continue
        if (grid[nr][nc] === 1) {
          H.addEdge(cur, nr * COL + nc)
        }
      }
    }
  }

  return H.work()
}

console.log(
  minimumOperations([
    [1, 1, 0],
    [0, 1, 1],
    [1, 1, 1]
  ])
)

export {}
