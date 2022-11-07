import { useHungarian } from './匈牙利算法'

const DIR4 = [
  [-1, 0],
  [1, 0],
  [0, -1],
  [0, 1]
]

/**
 * 有无穷块大小为1 * 2的多米诺骨牌,把这些骨牌不重叠地覆盖在完好的格子上，
 * 请找出你最多能在棋盘上放多少块骨牌
 *
 * @param {number} n  1 <= n <= 8
 * @param {number} m  1 <= m <= 8
 * @param {number[][]} broken  棋盘上每一个坏掉的格子的位置
 *
 * 棋盘格间隔染色变成二分图，注意1*2的骨牌覆盖两个格子相当于二分图匹配
 * 等价于求二分图最大匹配
 * https://leetcode.cn/problems/broken-board-dominoes/solution/suan-fa-xiao-ai-cong-ling-dao-yi-jiao-hu-8b4k/
 */
function domino(n: number, m: number, broken: number[][]): number {
  const board = Array.from({ length: n }, () => new Uint8Array(m))
  broken.forEach(([r, c]) => {
    board[r][c] = 1
  })

  const H = useHungarian(n * m, n * m) // 男生:偶数格子, 女生:奇数格子
  for (let r = 0; r < n; r++) {
    for (let c = 0; c < m; c++) {
      if (board[r][c] === 1 || (r + c) & 1) continue // !从男生连边到女生

      const cur = r * m + c
      for (const [dr, dc] of DIR4) {
        const nr = r + dr
        const nc = c + dc
        if (nr < 0 || nr >= n || nc < 0 || nc >= m) continue
        if (board[nr][nc] === 0) {
          H.addEdge(cur, nr * m + nc)
        }
      }
    }
  }

  return H.work()
}
export {}

if (require.main === module) {
  console.log(
    domino(2, 3, [
      [1, 0],
      [1, 1]
    ])
  )
  // 输出：2
  // 解释：我们最多可以放两块骨牌：[[0, 0], [0, 1]]以及[[0, 2], [1, 2]]。（见下图）
  console.log(
    domino(4, 4, [
      [1, 1],
      [2, 2]
    ])
  )
  // 预期结果：
  // 6
  console.log(
    domino(2, 3, [
      [1, 1],
      [1, 2]
    ])
  )
  // 预期结果：
  // 2

  // 1.通过题目条件建立无向图adjList，每个点可以连接上下左右。因为是无向图，所以建图时每个点只需考虑下方和右侧相连，避免重复连接。
  // 2.对无向图进行匈牙利算法获取最大匹配数。
  //     其中匈牙利算法的步骤:
  //     1.预处理，将图二分染色，得到colors数组。
  //     2.从二分图的左侧还没有匹配到的点出发，dfs寻找增广路径。如果找到，最大匹配数就加一。
}
