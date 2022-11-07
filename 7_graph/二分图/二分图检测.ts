// 二分图检测最好用dfs

import { useUnionFindMap } from '../../14_并查集/useUnionFind'

/**
 * 二分图检测
 */
function isBipartite(n: number, adjList: number[][]): boolean {
  const colors = new Int8Array(n).fill(-1)

  for (let i = 0; i < n; i++) {
    if (colors[i] === -1) {
      if (!dfs(i, 0)) return false
    }
  }

  return true

  function dfs(cur: number, color: number): boolean {
    colors[cur] = color
    for (const next of adjList[cur]) {
      if (colors[next] === -1) {
        if (!dfs(next, color ^ 1)) return false
      } else if (colors[next] === color) {
        return false
      }
    }

    return true
  }
}

// 扩展域并查集
function isBipartite2(adjMap: Map<number, Set<number>>): boolean {
  const uf = useUnionFindMap()
  const OFFSET = 1e9

  for (const [cur, nexts] of adjMap) {
    for (const next of nexts) {
      uf.union(cur, next + OFFSET)
      uf.union(next, cur + OFFSET)
      if (uf.isConnected(cur, next)) return false
    }
  }

  return true
}

export { isBipartite }
