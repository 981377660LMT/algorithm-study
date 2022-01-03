function hasCycle(adjList: number[][]) {
  const n = adjList.length
  const visited = Array<boolean>(n).fill(false)
  const onPath = Array<boolean>(n).fill(false)

  for (let i = 0; i < n; i++) {
    if (visited[i]) continue
    if (dfs(i)) return true
  }

  return false

  /**
   * @description
   * 与无向图的区别:`已经遍历过不代表形成环`
   * 需要添加标记，那些顶点在搜索路径上，回溯需要清除标记
   */
  function dfs(cur: number): boolean {
    if (onPath[cur]) return true
    if (visited[cur]) return false
    visited[cur] = true
    onPath[cur] = true

    for (const next of adjList[cur]) {
      if (dfs(next)) return true
    }

    onPath[cur] = false
    return false
  }
}

if (require.main === module) {
  console.log(hasCycle([[1], [0]]))

  // // 无环
  console.log(hasCycle([[1], [2, 3], [4], [2], []]))
  // 有环
  console.log(hasCycle([[1], [2], [0]]))
}

export { hasCycle }
