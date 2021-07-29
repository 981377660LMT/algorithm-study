import { AdjMap } from '../../chapter02图的基本表示/图的基本表示/2_邻接表'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'
import { MaxFlow } from '../../chapter14有向有权图最大流算法/网络流/MaxFlow'
import { BipartiteMaxGrow } from './最大流解决二分匹配'

/**
 * @param {number} n
 * @param {number} m
 * @param {number[][]} broken
 * @return {number}
 * @description 有无穷块大小为1 * 2的多米诺骨牌,把这些骨牌不重叠地覆盖在完好的格子上，请找出你最多能在棋盘上放多少块骨牌
 * @description 棋盘格间隔染色变成二分图，注意1*2的骨牌覆盖两个格子相当于二分图匹配
 * @description 等价于求二分图最大匹配
 */
const domino = (n: number, m: number, broken: number[][]): number => {
  const board = Array.from<number, number[]>({ length: n }, () => Array(m).fill(0))
  for (const [brokenX, brokenY] of broken) {
    board[brokenX][brokenY] = 1
  }
  const adjMap = new Map<number, Set<number>>()
  const delta = [
    [-1, 0],
    [1, 0],
    [0, -1],
    [0, 1],
  ]

  for (let i = 0; i < n; i++) {
    for (let j = 0; j < m; j++) {
      for (const [deltaX, deltaY] of delta) {
        const nextX = i + deltaX
        const nextY = j + deltaY
        if (
          nextX < n &&
          nextX >= 0 &&
          nextY < m &&
          nextY >= 0 &&
          board[i][j] === 0 &&
          board[nextX][nextY] === 0
        ) {
          adjMap.set(
            i * m + j,
            adjMap.get(i * m + j)?.add(nextX * m + nextY) || new Set([nextX * m + nextY])
          )
        }
      }
    }
  }

  // 我们并不关心多少条边，所以设Infinity
  const netWork = new AdjMap(n * m, Infinity, adjMap, false, [], [])

  const dfs = new DFS(netWork)
  return new BipartiteMaxGrow(dfs).maxMatching
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
