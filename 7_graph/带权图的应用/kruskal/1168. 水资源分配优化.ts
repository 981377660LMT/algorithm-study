type U = number
type V = number
type Weight = number

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
  const edges: [U, V, Weight][] = []
  wells.forEach((weight, index) => edges.push([0, index + 1, weight]))
  pipes.forEach(([u, v, w]) => edges.push([u, v, w]))
  edges.sort((a, b) => a[2] - b[2])
  return useUnionFind(n + 1, edges)

  function useUnionFind(size: number, edges: [U, V, Weight][]) {
    let cost = 0
    const parent = Array.from<number, number>({ length: size }, (_, i) => i)

    const find = (key: number) => {
      while (parent[key] !== undefined && parent[key] !== key) {
        parent[key] = parent[parent[key]]
        key = parent[key]
      }
      return key
    }

    for (const [u, v, w] of edges) {
      const root1 = find(u)
      const root2 = find(v)
      if (root1 !== root2) {
        cost += w
        parent[Math.max(root1, root2)] = Math.min(root1, root2)
      }
    }

    return cost
  }
}

console.log(
  minCostToSupplyWater(
    3,
    [1, 2, 2],
    [
      [1, 2, 1],
      [2, 3, 1],
    ]
  )
)
