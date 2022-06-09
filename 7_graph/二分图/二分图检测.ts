// 二分图检测最好用dfs

import { useUnionFindMap } from '../../14_并查集/useUnionFind'

/**
 * @description 二分图检测
 */
function isBipartite<V = number>(adjMap: Map<V, Set<V>>): boolean {
  const colorMap = new Map<V, number>()

  for (const cur of adjMap.keys()) {
    if (colorMap.has(cur)) continue
    if (!dfs(cur, 0)) return false
  }

  return true

  function dfs(cur: V, color: number): boolean {
    colorMap.set(cur, color)

    for (const next of adjMap.get(cur) ?? []) {
      if (!colorMap.has(next)) {
        if (!dfs(next, color ^ 1)) return false
      } else if (colorMap.get(next) === color) return false
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
