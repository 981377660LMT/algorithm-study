type U = number
type V = number
type Weight = number

/**
 * @param {number[][]} points
 * @return {number}
 * 连接点 [xi, yi] 和点 [xj, yj] 的费用为它们之间的 曼哈顿距
 *
 */
const minCostConnectPoints = function (points: number[][]): number {
  const useUnionFind = (size: number, edges: [U, V, Weight][]) => {
    let res = 0
    const parent = Array.from<number, number>({ length: size }, (_, i) => i)

    const find = (key: number) => {
      while (parent[key] && parent[key] !== key) {
        parent[key] = parent[parent[key]]
        key = parent[key]
      }
      return key
    }

    for (const [u, v, w] of edges) {
      const root1 = find(u)
      const root2 = find(v)
      // 不连通
      if (root1 !== root2) {
        res += w
        parent[Math.max(root1, root2)] = Math.min(root1, root2)
      }
    }

    return res
  }

  const edges: [U, V, Weight][] = [] // 图中两两点间的权值
  for (let i = 0; i < points.length; i++) {
    for (let j = i + 1; j < points.length; j++) {
      const [x1, y1] = points[i]
      const [x2, y2] = points[j]
      const weight = Math.abs(x1 - x2) + Math.abs(y1 - y2)
      edges.push([i, j, weight])
    }
  }
  edges.sort((a, b) => a[2] - b[2])

  return useUnionFind(points.length, edges)
}

console.log(
  minCostConnectPoints([
    [0, 0],
    [2, 2],
    [3, 10],
    [5, 2],
    [7, 0],
  ])
)

export default 1

1
