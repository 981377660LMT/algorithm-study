function hasCycle(adjList: number[][]) {
  const n = adjList.length
  const visited = Array<boolean>(n).fill(false)

  for (let i = 0; i < n; i++) {
    if (visited[i]) continue
    if (dfs(i, -1)) return true
  }

  return false

  function dfs(cur: number, pre: number): boolean {
    if (visited[cur]) return true
    visited[cur] = true

    for (const next of adjList[cur]) {
      if (next === pre) continue
      if (dfs(next, cur)) return true
    }

    return false
  }
}

// 0->1->2->0
console.log(
  hasCycle([
    [1, 2],
    [0, 2],
    [0, 1],
  ])
)

// 0->1 只有一条边，不存在环
console.log(hasCycle([[1], [0]]))
