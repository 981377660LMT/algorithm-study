import { kruskal1 } from './Kruskal'

// 该最小成本应该是所用全部连接代价的综合。如果根据已知条件无法完成该项任务，则请你返回 -1。
function minimumCost(n: number, connections: number[][]): number {
  return kruskal1(n, connections as [u: number, v: number, w: number][])
}

console.log(
  minimumCost(3, [
    [1, 2, 5],
    [1, 3, 6],
    [2, 3, 1]
  ])
)

export {}
