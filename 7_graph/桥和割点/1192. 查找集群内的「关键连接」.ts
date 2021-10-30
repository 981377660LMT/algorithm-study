// 1. order 数组记录 dfs 的访问顺序
// 2. low 数组记录能访问到的**最小的 orderIndex**
// 3. 前进时对 cur 的 low 取 min
// 4. 回退的最后阶段，如果low[next]>order[cur]，那么[cur,next]是桥
function criticalConnections(n: number, connections: number[][]): number[][] {
  const res: [cur: number, next: number][] = []

  const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  connections.forEach(([cur, next]) => {
    adjList[cur].push(next)
    adjList[next].push(cur)
  })

  let orderIndex = 0
  const visited = Array<boolean>(n).fill(false)
  const order = Array<number>(n).fill(Infinity)
  const low = Array<number>(n).fill(Infinity)

  dfs(0, Infinity)
  return res

  function dfs(cur: number, parent: number) {
    if (visited[cur]) return
    visited[cur] = true

    order[cur] = orderIndex
    low[cur] = orderIndex
    orderIndex++

    for (const next of adjList[cur]) {
      if (next === parent) continue
      if (!visited[next]) {
        dfs(next, cur)
        // 回退阶段
        low[cur] = Math.min(low[cur], low[next])
        if (low[next] > order[cur]) res.push([cur, next])
      } else {
        low[cur] = Math.min(low[cur], low[next])
      }
    }
  }
}

console.log(
  criticalConnections(4, [
    [0, 1],
    [1, 2],
    [2, 0],
    [1, 3],
  ])
)

// connections[i][0] != connections[i][1]
// 不存在重复的连接

// 没有自环边
// 没有平行边
// 连通图
// ps:没有自环边并且没有平行边的图称为**简单图**
