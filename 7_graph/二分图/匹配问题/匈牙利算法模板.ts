function hungarian(adjList: number[][]) {
  let maxMatching = 0
  let visited: boolean[]
  const matching = Array<number>(adjList.length).fill(-1)

  const colors = bisect(adjList)
  for (let i = 0; i < adjList.length; i++) {
    visited = Array<boolean>(adjList.length).fill(false)
    if (colors[i] === 0 && matching[i] === -1) {
      if (dfs(i)) maxMatching++
    }
  }

  return maxMatching

  function bisect(adjList: number[][]): number[] {
    const colors = Array<number>(adjList.length).fill(-1)

    const dfs = (cur: number, color: number) => {
      colors[cur] = color

      for (const next of adjList[cur]) {
        if (colors[next] === -1) {
          dfs(next, color ^ 1)
        } else {
          if (colors[cur] === colors[next]) throw new Error('不是二分图')
        }
      }
    }

    for (let i = 0; i < adjList.length; i++) {
      if (colors[i] === -1) dfs(i, 0)
    }

    return colors
  }

  function dfs(cur: number): boolean {
    if (visited[cur]) return false
    visited[cur] = true

    for (const next of adjList[cur]) {
      if (matching[next] === -1 || dfs(matching[next])) {
        matching[cur] = next
        matching[next] = cur
        return true
      }
    }

    return false
  }
}

if (require.main === module) {
  console.assert(hungarian([[1, 3], [0, 2], [1], [0], [], []]) === 2)
}
