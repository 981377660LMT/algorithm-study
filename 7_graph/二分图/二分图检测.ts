function isBipartite(n: number, adjList: number[][]): [colors: Int8Array, ok: boolean] {
  const colors = new Int8Array(n).fill(-1)
  for (let i = 0; i < n; i++) {
    if (colors[i] === -1 && !dfs(i, 0)) {
      return [colors, false]
    }
  }
  return [colors, true]

  function dfs(cur: number, color: number): boolean {
    colors[cur] = color
    for (let i = 0; i < adjList[cur].length; i++) {
      const next = adjList[cur][i]
      if (colors[next] === -1) {
        if (!dfs(next, color ^ 1)) {
          return false
        }
      } else if (colors[next] === color) {
        return false
      }
    }
    return true
  }
}

export { isBipartite }
