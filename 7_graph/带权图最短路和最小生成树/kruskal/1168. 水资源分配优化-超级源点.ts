import { kruskal1 } from './Kruskal'

type Edge<V = unknown> = [u: V, v: V, weight: number]

/**
 *
 * @param n 1到n
 * @param wells  一种是直接在房子内建造水井，成本为 wells[i] (即与虚拟节点水源0相连边的权重为wells[i])
 * @param pipes  另一种是从另一口井铺设管道引水  pipes[i] = [house1, house2, cost] 代表用管道将 house1 和 house2 连接在一起的成本
 * 建造水井和铺设管道来为所有房子供水。
 * 请你帮忙计算为所有房子都供水的最低总成本。
 * @summary
 * 加一个虚拟节点就和1135一样
 */
function minCostToSupplyWater(n: number, wells: number[], pipes: number[][]): number {
  const edges: Edge<number>[] = []
  wells.forEach((weight, index) => edges.push([0, index + 1, weight]))
  pipes.forEach(([u, v, w]) => edges.push([u, v, w]))
  return kruskal1(n + 1, edges)
}

console.log(
  minCostToSupplyWater(
    3,
    [1, 2, 2],
    [
      [1, 2, 1],
      [2, 3, 1]
    ]
  )
)

export {}
