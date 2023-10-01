/* eslint-disable no-inner-declarations */

import { bfs } from '../../bfs求无权图的最短路径/bfs模板'

/**
 * 返回基环树森林的环分组信息(环的大小>=2)以及每个点在拓扑排序中的最大深度.
 * @param n 图的节点数.
 * @param graph 图的邻接表表示.
 * @param directed 图是否有向.
 */
function cyclePartition(
  n: number,
  graph: ArrayLike<ArrayLike<number>>,
  directed: boolean
): {
  /** 环分组.每个环的大小>=2. */
  groups: number[][]
  /** 每个点是否在环中. */
  inCycle: Uint8Array
  /** 每个点所在的环的编号.如果不在环中,则为-1. */
  belong: Uint32Array
  /** 每个点在拓扑排序中的最大深度.最外层的点深度为0. */
  depth: Uint32Array
} {
  const deg = new Uint32Array(n)
  if (directed) {
    for (let u = 0; u < n; ++u) {
      const nexts = graph[u]
      for (let i = 0; i < nexts.length; ++i) {
        const v = nexts[i]
        deg[v]++
      }
    }
  } else {
    for (let u = 0; u < n; ++u) {
      deg[u] = graph[u].length
    }
  }

  const startDeg = directed ? 0 : 1
  const visited = new Uint8Array(n)
  const depth = new Uint32Array(n)
  const queue = new Uint32Array(n)
  let head = 0
  let tail = 0
  for (let i = 0; i < n; ++i) {
    if (deg[i] === startDeg) {
      queue[tail++] = i
    }
  }

  while (head < tail) {
    const cur = queue[head++]
    visited[cur] = 1
    const nexts = graph[cur]
    for (let i = 0; i < nexts.length; ++i) {
      const next = nexts[i]
      depth[next] = Math.max(depth[next], depth[cur] + 1)
      deg[next]--
      if (deg[next] === startDeg) {
        queue[tail++] = next
      }
    }
  }

  const dfs = (cur: number, path: number[]) => {
    if (visited[cur]) return
    visited[cur] = 1
    path.push(cur)
    const nexts = graph[cur]
    for (let i = 0; i < nexts.length; ++i) {
      dfs(nexts[i], path)
    }
  }

  const groups: number[][] = []
  for (let i = 0; i < n; ++i) {
    if (visited[i]) continue
    const path: number[] = []
    dfs(i, path)
    groups.push(path)
  }

  const inCycle = new Uint8Array(n)
  const belong = new Uint32Array(n)
  for (let gid = 0; gid < groups.length; ++gid) {
    const group = groups[gid]
    for (let i = 0; i < group.length; ++i) {
      const node = group[i]
      inCycle[node] = 1
      belong[node] = gid
    }
  }

  return { groups, inCycle, belong, depth }
}

export { cyclePartition }

if (require.main === module) {
  // 100075. 有向图访问计数
  // https://leetcode.cn/problems/count-visited-nodes-in-a-directed-graph/description/
  function countVisitedNodes(edges: number[]): number[] {
    const n = edges.length
    const adjList: number[][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    for (let i = 0; i < n; i++) adjList[i].push(edges[i])
    const { groups, inCycle, belong } = cyclePartition(n, adjList, true)

    const cache = new Int32Array(n).fill(-1)
    const dfs = (cur: number): number => {
      if (cache[cur] !== -1) return cache[cur]
      if (inCycle[cur]) return groups[belong[cur]].length
      const next = edges[cur]
      const res = dfs(next) + 1
      cache[cur] = res
      return res
    }

    const res: number[] = Array(n)
    for (let i = 0; i < n; i++) res[i] = dfs(i)
    return res
  }

  // 2204. 无向图中到环的距离
  // https://leetcode.cn/problems/distance-to-a-cycle-in-undirected-graph/
  function distanceToCycle(n: number, edges: number[][]): number[] {
    const adjList: number[][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    edges.forEach(([u, v]) => {
      adjList[u].push(v)
      adjList[v].push(u)
    })

    const { groups } = cyclePartition(n, adjList, false)
    return bfs(n, adjList, groups[0])
  }
}
