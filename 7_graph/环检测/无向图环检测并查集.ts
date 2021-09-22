import { useUnionFindArray } from '../../14_并查集/推荐使用并查集精简版'

const hasCycle = (n: number, edges: [number, number][]) => {
  const uf = useUnionFindArray(n)

  for (const [u, w] of edges) {
    if (uf.isConnected(u, w)) return true
    uf.union(u, w)
  }

  return false
}

export {}
