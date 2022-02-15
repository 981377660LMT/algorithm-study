import { hungarian } from './匈牙利算法'

/**
 * @param {number} n  1 <= n <= 8
 * @param {number} m  1 <= m <= 8
 * @param {number[][]} broken  棋盘上每一个坏掉的格子的位置
 * @return {number}
 * @description 有无穷块大小为1 * 2的多米诺骨牌,把这些骨牌不重叠地覆盖在完好的格子上，请找出你最多能在棋盘上放多少块骨牌
 * @description 棋盘格间隔染色变成二分图，注意1*2的骨牌覆盖两个格子相当于二分图匹配
 * @description 等价于求二分图最大匹配
 */
const domino = (n: number, m: number, broken: number[][]): number => {
  const board = Array.from<unknown, number[]>({ length: n }, () => Array(m).fill(0))
  for (const [brokenX, brokenY] of broken) {
    board[brokenX][brokenY] = 1
  }

  const adjList = Array.from<unknown, number[]>({ length: n * m }, () => [])
  // 建图 因为是无向图 所以只需下方和右侧相连避免重复看
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < m; j++) {
      if (board[i][j] === 1) continue
      const cur = i * m + j

      if (j + 1 < m && board[i][j + 1] === 0) {
        const next = i * m + j + 1
        adjList[cur].push(next)
        adjList[next].push(cur)
      }

      if (i + 1 < n && board[i + 1][j] === 0) {
        const next = (i + 1) * m + j
        adjList[cur].push(next)
        adjList[next].push(cur)
      }
    }
  }

  return hungarian(adjList)
}

console.log(
  domino(2, 3, [
    [1, 0],
    [1, 1],
  ])
)
// 输出：2
// 解释：我们最多可以放两块骨牌：[[0, 0], [0, 1]]以及[[0, 2], [1, 2]]。（见下图）
export {}
console.log(
  domino(4, 4, [
    [1, 1],
    [2, 2],
  ])
)
// 预期结果：
// 6
console.log(
  domino(2, 3, [
    [1, 1],
    [1, 2],
  ])
)
// 预期结果：
// 2

// 1.通过题目条件建立无向图adjList，每个点可以连接上下左右。因为是无向图，所以建图时每个点只需考虑下方和右侧相连，避免重复连接。
// 2.对无向图进行匈牙利算法获取最大匹配数。
//     其中匈牙利算法的步骤:
//     1.预处理，将图二分染色，得到colors数组。
//     2.从二分图的左侧还没有匹配到的点出发，dfs寻找增广路径。如果找到，最大匹配数就加一。
