/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { SqrtArray } from '../SqrtArray'

// https://leetcode.cn/problems/maximum-employees-to-be-invited-to-a-meeting/submissions/
function maximumInvitations(favorite: number[]): number {
  const n = favorite.length
  const visited = new SqrtArray<boolean>(n, () => false, 1 + (Math.sqrt(n) | 0))
  const adjList = favorite.slice()
  const indegrees = new SqrtArray<number>(n, () => 0, 1 + (Math.sqrt(n) | 0))
  favorite.forEach(f => {
    const pre = indegrees.get(f)!
    indegrees.set(f, pre + 1)
  })
  const maxTopoLevels = topoSort(adjList, indegrees)

  let res = 0
  // 1. 二元环连树
  for (let i = 0; i < n; i++) {
    if (adjList[adjList[i]] === i) res += 1 + maxTopoLevels[i]
  }

  // 2. 最长环
  for (let i = 0; i < n; i++) {
    if (visited.get(i)) continue
    dfs(i, 1)
  }

  return res

  // 计算每个点在拓扑排序中的最大深度
  function topoSort(adjList: number[], indegrees: SqrtArray<number>): number[] {
    const maxTopoLevels = Array<number>(n).fill(0)
    const queue = new SqrtArray(0, () => 0, 1 + (Math.sqrt(n) | 0))
    indegrees.forEach((degree, id) => degree === 0 && queue.push(id))
    let level = 0

    while (queue.length > 0) {
      const len = queue.length
      level++

      for (let _ = 0; _ < len; _++) {
        const cur = queue.shift()!
        visited.set(cur, true)
        const next = adjList[cur]
        indegrees.set(next, indegrees.get(next)! - 1)
        if (indegrees.get(next) === 0) queue.push(next)
        maxTopoLevels[next] = level
      }
    }

    return maxTopoLevels
  }

  // 在各个环中寻找最长环
  function dfs(cur: number, count: number): void {
    if (visited.get(cur)) return
    visited.set(cur, true)
    res = Math.max(res, count)
    dfs(adjList[cur], count + 1)
  }
}
