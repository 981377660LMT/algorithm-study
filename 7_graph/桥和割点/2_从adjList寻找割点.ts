function findCutvertex(n: number, adjList: number[][]): number[] {
  const res = new Set<number>()

  let orderIndex = 0
  const visited = Array<boolean>(n).fill(false)
  const order = Array<number>(n).fill(Infinity)
  const low = Array<number>(n).fill(Infinity)

  dfs(0, Infinity)
  return [...res]

  function dfs(cur: number, parent: number) {
    if (visited[cur]) return
    visited[cur] = true

    order[cur] = orderIndex
    low[cur] = orderIndex
    orderIndex++

    let dfsChild = 0
    for (const next of adjList[cur]) {
      if (next === parent) continue

      if (!visited[next]) {
        dfsChild++
        dfs(next, cur)

        // 回退阶段
        low[cur] = Math.min(low[cur], low[next])
        // 根节点单独讨论
        if (parent !== Infinity && low[next] >= order[cur]) res.add(cur)
        else if (parent === Infinity && dfsChild > 1) res.add(cur)
      } else {
        low[cur] = Math.min(low[cur], low[next])
      }
    }
  }
}

console.log(findCutvertex(4, [[1, 2], [0, 2, 3], [1, 0], [1]]))
