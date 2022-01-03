import { useUnionFindArray } from './推荐使用并查集精简版'

// n-1条边且无环
const validTree = (n: number, edges: number[][]) => {
  if (edges.length !== n - 1) return false
  const uf = useUnionFindArray(n)
  for (const [u, v] of edges) {
    if (uf.isConnected(u, v)) return false
    uf.union(u, v)
  }
  return true
}
// 判断图是不是树，满足两点

// 是连通图
// 不存在环
export { validTree }
