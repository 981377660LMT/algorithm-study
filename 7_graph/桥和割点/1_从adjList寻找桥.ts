// 寻找无向图的桥 tarjan算法
// 1. order 数组记录 dfs 的访问顺序
// 2. low 数组记录能访问到的**最小的 orderIndex**
// 3. 对 cur 的 low 取 min
// 4. 回退的最后阶段，如果low[next]>order[cur]，那么[cur,next]是桥
function findBridge(n: number, adjList: number[][]): [number, number][] {
  const res: [cur: number, next: number][] = []

  let orderIndex = 0
  const visited = Array<boolean>(n).fill(false)
  const order = Array<number>(n).fill(Infinity)
  const low = Array<number>(n).fill(Infinity)

  for (let i = 0; i < n; i++) {
    if (visited[i]) continue
    dfs(i, Infinity)
  }

  return res

  function dfs(cur: number, parent: number): void {
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

console.log(findBridge(4, [[1, 2], [0, 2, 3], [1, 0], [1]]))
