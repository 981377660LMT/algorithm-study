import { ArrayDeque } from '../../2_queue/Deque'

/**
 * @param {number} n  有 n 名玩家，所有玩家编号分别为 0 ～ n-1，其中小朋友 A 的编号为 0
 * @param {number[][]} relation  传信息的关系是单向的
 * @param {number} k
 * @return {number}
 * 每轮信息必须需要传递给另一个人，且信息可重复经过同一个人
 * 返回信息从小 A (编号 0 ) 经过 k 轮传递到编号为 n-1 的小伙伴处的方案数；若不能到达，返回 0
 * @summary
 * 无权图统计有限步数的到达某个节点的方案数，最常见的方式是使用 BFS 或 DFS
 */
var numWays = function (n: number, relation: number[][], k: number): number {
  const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  relation.forEach(([u, v]) => adjList[u].push(v))
  let res = 0
  const queue = new ArrayDeque(10 ** 4)
  queue.push(0)

  while (queue.length && k) {
    const len = queue.length

    for (let i = 0; i < len; i++) {
      const cur = queue.shift()!
      for (const next of adjList[cur]) {
        if (next === n - 1 && k === 1) {
          res++
        }
        queue.push(next)
      }
    }

    k--
  }

  return res
}

console.log(
  numWays(
    5,
    [
      [0, 2],
      [2, 1],
      [3, 4],
      [2, 3],
      [1, 4],
      [2, 0],
      [0, 4],
    ],
    3
  )
)
console.log(1e1)
