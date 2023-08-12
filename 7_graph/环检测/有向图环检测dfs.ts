function hasCycle(n: number, adjList: number[][]) {
  const visited = new Uint8Array(n)
  const onPath = new Uint8Array(n)

  for (let i = 0; i < n; i++) {
    if (visited[i]) continue
    if (dfs(i)) return true
  }

  return false

  /**
   * 与无向图的区别:`已经遍历过不代表形成环`
   * 需要添加标记，那些顶点在搜索路径上，回溯需要清除标记
   */
  function dfs(cur: number): boolean {
    if (onPath[cur]) return true
    if (visited[cur]) return false
    visited[cur] = 1
    onPath[cur] = 1

    for (const next of adjList[cur]) {
      if (dfs(next)) return true
    }

    onPath[cur] = 0
    return false
  }
}

export { hasCycle }
