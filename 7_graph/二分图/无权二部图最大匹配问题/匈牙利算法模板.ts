// 复杂度O(NM)
function hungarian(adjList: number[][]) {
  let maxMatching = 0
  const matched = Array<number>(adjList.length).fill(-1)

  const colors = genColor(adjList)
  for (let i = 0; i < adjList.length; i++) {
    const visited = Array<boolean>(adjList.length).fill(false)
    if (colors[i] === 0 && matched[i] === -1) {
      if (dfs(i, visited)) maxMatching++
    }
  }

  return maxMatching

  function dfs(boy: number, visited: boolean[]): boolean {
    if (visited[boy]) return false
    visited[boy] = true

    for (const girl of adjList[boy]) {
      // 女孩还没配对,或则可以找到增广路径
      if (matched[girl] === -1 || dfs(matched[girl], visited)) {
        matched[boy] = girl
        matched[girl] = boy
        return true
      }
    }

    return false
  }

  function genColor(adjList: number[][]): number[] {
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
}

if (require.main === module) {
  console.assert(hungarian([[1, 3], [0, 2], [1], [0], [], []]) === 2)
}
