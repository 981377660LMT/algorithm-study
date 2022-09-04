1. 准备好所有边 [u,v,weight] 按 wweight 升序排列
2. 并查集构造生成树

```JS
type U = number
type V = number
type Weight = number

function minimumCost(n: number, connections: number[][]): number {
  const edges: [U, V, Weight][] = [] // 图中两两点间的权值
  for (let i = 0; i < connections.length; i++) {
    for (let j = i + 1; j < connections.length; j++) {
      const [x1, y1] = connections[i]
      const [x2, y2] = connections[j]
      const weight = Math.abs(x1 - x2) + Math.abs(y1 - y2)
      edges.push([i, j, weight])
    }
  }
  edges.sort((a, b) => a[2] - b[2])

  return useUnionFind(n, edges)

  function useUnionFind(size: number, edges: [U, V, Weight][]) {
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
}

```
