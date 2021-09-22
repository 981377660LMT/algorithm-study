const hasCycle = (n: number, prerequisites: [number, number][]) => {
  const visited = Array<boolean>(n).fill(false)
  const path = new Set<number>()
  const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  for (const [cur, pre] of prerequisites) {
    adjList[pre].push(cur)
  }

  /**
   *
   * @param cur
   * @param visited
   * @param path 用于检测有向图的环,回溯需要删除
   * @returns
   */
  const dfs = (cur: number, visited: boolean[], path: Set<number>): boolean => {
    visited[cur] = true
    path.add(cur)

    for (const next of adjList[cur]) {
      if (!visited[next]) {
        if (dfs(next, visited, path)) return true
      } else {
        if (path.has(next)) {
          return true
        }
      }
    }

    path.delete(cur)
    return false
  }

  for (let i = 0; i < n; i++) {
    if (visited[i]) continue
    if (dfs(i, visited, path)) return true
  }

  return false
}

export {}

console.log(
  hasCycle(2, [
    [0, 1],
    [1, 0],
  ])
)
