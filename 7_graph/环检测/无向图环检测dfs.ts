const hasCycle = (n: number, edges: [number, number][]) => {
  const visited = Array<boolean>(n).fill(false)
  const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  for (const [u, v] of edges) {
    adjList[v].push(u)
    adjList[u].push(v)
  }

  /**
   *
   * @param cur
   * @param pre
   * @returns
   * 有无环
   */
  const dfs = (cur: number, pre: number): boolean => {
    visited[cur] = true

    for (const next of adjList[cur]) {
      if (!visited[next]) {
        if (dfs(next, cur)) return true
      } else {
        // 走回了之前走过的非 pre 的节点
        if (next !== pre) {
          return true
        }
      }
    }

    return false
  }

  for (let i = 0; i < n; i++) {
    if (visited[i]) continue
    if (dfs(i, i)) return true
  }

  return false
}
