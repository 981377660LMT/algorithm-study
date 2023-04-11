import { UnionFindArray } from '../../14_并查集/UnionFind'

const hasCycle = (n: number, edges: [number, number][]) => {
  const uf = new UnionFindArray(n)

  for (const [u, w] of edges) {
    if (uf.isConnected(u, w)) return true
    uf.union(u, w)
  }

  return false
}

export {}
