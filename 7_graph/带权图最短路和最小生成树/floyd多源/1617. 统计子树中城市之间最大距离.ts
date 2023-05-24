// 2 <= n <= 15
// 请你返回一个大小为 n-1 的数组，其中第 d 个元素（下标从 1 开始）是城市间 最大距离 恰好等于 d 的`子树`数目。
// 1.求每个点到所有点的最短距离--多源最短路径算法 floyd
// 2.二进制枚举所有状态（子集）

import { Floyd } from './Floyd'

function countSubgraphsForEachDiameter(n: number, edges: number[][]): number[] {
  const F = new Floyd(n)
  edges.forEach(([u, v]) => F.addEdge(u - 1, v - 1, 1))
  const res: number[] = Array(n).fill(0)

  for (let state = 0; state < 1 << n; state++) {
    let maxDist = 0
    let edge = 0
    const vertex: number[] = []
    for (let i = 0; i < n; i++) {
      if ((state >> i) & 1) vertex.push(i)
    }

    for (let i = 0; i < vertex.length; i++) {
      for (let j = i + 1; j < vertex.length; j++) {
        maxDist = Math.max(maxDist, F.dist(vertex[i], vertex[j]))
        edge += +(F.dist(vertex[i], vertex[j]) === 1)
      }
    }

    const isTree = vertex.length === edge + 1
    if (isTree) res[maxDist]++
  }

  return res.slice(1)
}

console.log(
  countSubgraphsForEachDiameter(4, [
    [1, 2],
    [2, 3],
    [2, 4]
  ])
)
// 输出：[3,4,0]
