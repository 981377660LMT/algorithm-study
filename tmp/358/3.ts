export {}

const INF = 2e15
// # 给你一个下标从 0 开始、大小为 n x n 的二维矩阵 grid ，其中 (r, c) 表示：

// # 如果 grid[r][c] = 1 ，则表示一个存在小偷的单元格
// # 如果 grid[r][c] = 0 ，则表示一个空单元格
// # 你最开始位于单元格 (0, 0) 。在一步移动中，你可以移动到矩阵中的任一相邻单元格，包括存在小偷的单元格。

// # 矩阵中路径的 安全系数 定义为：从路径中任一单元格到矩阵中任一小偷所在单元格的 最小 曼哈顿距离。

// # 返回所有通向单元格 (n - 1, n - 1) 的路径中的 最大安全系数 。

// # 单元格 (r, c) 的某个 相邻 单元格，是指在矩阵中存在的 (r, c + 1)、(r, c - 1)、(r + 1, c) 和 (r - 1, c) 之一。

// # 两个单元格 (a, b) 和 (x, y) 之间的 曼哈顿距离 等于 | a - x | + | b - y | ，其中 |val| 表示 val 的绝对值。

// # 预处理出每个格子的安全系数
function maximumSafenessFactor(grid: number[][]): number {
  const ROW = grid.length
  const COL = grid[0].length
  // 每个格子到小偷的最小曼哈顿距离
  const distToThief = Array.from({ length: ROW }, () => Array.from({ length: COL }, () => INF))
}
