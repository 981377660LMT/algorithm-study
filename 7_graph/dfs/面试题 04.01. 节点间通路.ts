function findWhetherExistsPath(
  n: number,
  graph: number[][],
  start: number,
  target: number
): boolean {
  const adjList = Array.from<number, number[]>({ length: n + 1 }, () => [])
  for (const [u, v] of graph) {
    adjList[u].push(v)
  }

  const dfs = (cur: number, visited: Set<number>): boolean => {
    if (cur === target) return true
    for (const next of adjList[cur]) {
      if (visited.has(next)) continue
      visited.add(next)
      if (dfs(next, visited)) return true
    }
    return false
  }

  return dfs(start, new Set([start]))
}
