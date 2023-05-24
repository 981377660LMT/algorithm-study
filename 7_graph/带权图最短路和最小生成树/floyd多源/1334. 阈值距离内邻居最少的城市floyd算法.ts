import { floyd } from './Floyd'

const INF = 2e15

/**
 *
 * @param n  2 <= n <= 100
 * @param edges
 * @param distanceThreshold
 * 返回能通过某些路径到达其他城市数目最少、
 * 且路径距离 最大 为 distanceThreshold 的城市。如果有多个这样的城市，则返回编号最大的城市。
 * @summary
 * 要求出求所有点对最短路径,使用floyd算法
 */
function findTheCity(n: number, edges: number[][], distanceThreshold: number): number {
  const dist = floyd(n, edges, false)
  let res = 0
  let minNeighbor = Infinity
  for (let i = 0; i < n; i++) {
    let count = 0
    for (const dis of dist[i]) {
      dis <= distanceThreshold && count++
    }

    if (count <= minNeighbor) {
      res = i
      minNeighbor = count
    }
  }

  return res
}

console.log(
  findTheCity(
    4,
    [
      [0, 1, 3],
      [1, 2, 1],
      [1, 3, 4],
      [2, 3, 1]
    ],
    4
  )
)

export {}
