// 基环树:树加一条边使之成环
// 两种情况:1.所有的二元环连树;2.唯一的最长环
function maximumInvitations(favorite: number[]): number {
  const n = favorite.length
  const visited = Array<boolean>(n).fill(false)
  const maxTopoLevels = Array<number>(n).fill(0)
  const adjList = favorite.slice()
  const indegrees = Array<number>(n).fill(0)
  favorite.forEach(f => indegrees[f]++)
  topoSort(adjList, indegrees)

  let res = 0

  // 1. 二元环连树
  for (let i = 0; i < n; i++) {
    if (adjList[adjList[i]] === i) res += 1 + maxTopoLevels[i]
  }

  // 2. 最长环
  for (let i = 0; i < n; i++) {
    if (visited[i]) continue
    dfs(i, 1)
  }

  return res

  // 计算每个点在拓扑排序中的最大深度
  function topoSort(adjList: number[], indegrees: number[]): void {
    let level = 0
    let queue: number[] = []
    indegrees.forEach((degree, id) => degree === 0 && queue.push(id))

    while (queue.length > 0) {
      const len = queue.length
      const nextQueue: number[] = []
      level++

      for (let _ = 0; _ < len; _++) {
        const cur = queue.pop()!
        visited[cur] = true
        const next = adjList[cur]
        indegrees[next]--
        if (indegrees[next] === 0) nextQueue.push(next)
        maxTopoLevels[next] = level
      }

      queue = nextQueue
    }
  }

  // 在各个环中寻找最长环
  function dfs(cur: number, count: number): void {
    if (visited[cur]) return
    visited[cur] = true
    res = Math.max(res, count)
    dfs(adjList[cur], count + 1)
  }
}

console.log(maximumInvitations([2, 2, 1, 2]))
console.log(maximumInvitations([3, 0, 1, 4, 1]))

// 数组用s结尾比较好
// 非纯函数返回void比较好
