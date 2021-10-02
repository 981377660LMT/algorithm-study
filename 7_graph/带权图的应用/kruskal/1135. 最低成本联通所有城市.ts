type U = number
type V = number
type Weight = number

// 该最小成本应该是所用全部连接代价的综合。如果根据已知条件无法完成该项任务，则请你返回 -1。
function minimumCost(n: number, connections: number[][]): number {
  return useUnionFind(n, connections.sort((a, b) => a[2] - b[2]) as [U, V, Weight][])

  function useUnionFind(size: number, edges: [U, V, Weight][]) {
    let res = 0
    let count = size
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
        count--
      }
    }

    return count === 1 ? res : -1
  }
}

console.log(
  minimumCost(3, [
    [1, 2, 5],
    [1, 3, 6],
    [2, 3, 1],
  ])
)

export {}
