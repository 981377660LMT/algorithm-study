import { Tree } from '../../6_tree/重链剖分/Tree'

export {}

const INF = 2e15
// 现有一棵由 n 个节点组成的无向树，节点按从 0 到 n - 1 编号。给你一个整数 n 和一个长度为 n - 1 的二维整数数组 edges ，其中 edges[i] = [ui, vi, wi] 表示树中存在一条位于节点 ui 和节点 vi 之间、权重为 wi 的边。

// 另给你一个长度为 m 的二维整数数组 queries ，其中 queries[i] = [ai, bi] 。对于每条查询，请你找出使从 ai 到 bi 路径上每条边的权重相等所需的 最小操作次数 。在一次操作中，你可以选择树上的任意一条边，并将其权重更改为任意值。

// 注意：

// 查询之间 相互独立 的，这意味着每条新的查询时，树都会回到 初始状态 。
// 从 ai 到 bi的路径是一个由 不同 节点组成的序列，从节点 ai 开始，到节点 bi 结束，且序列中相邻的两个节点在树中共享一条边。
// 返回一个长度为 m 的数组 answer ，其中 answer[i] 是第 i 条查询的答案。
function minOperationsQueries(n: number, edges: number[][], queries: number[][]): number[] {
  const adjList: [number, number][][] = Array(n)
  for (let i = 0; i < n; i++) adjList[i] = []
  const weights: Map<number, number>[] = Array(n)
  for (let i = 0; i < n; i++) weights[i] = new Map()
  edges.forEach(([u, v, w]) => {
    adjList[u].push([v, w])
    adjList[v].push([u, w])
    weights[u].set(v, w)
    weights[v].set(u, w)
  })

  // 对每个查询求出边权，答案为边数减去最大频率
  const res = Array(queries.length).fill(0)
  for (let i = 0; i < queries.length; i++) {
    const [u, v] = queries[i]
    const weightCounter = new Uint32Array(30)
    getPath(u, -1, v, weightCounter)
    let sum = 0
    let maxFreq = 0
    for (let i = 0; i < weightCounter.length; i++) {
      sum += weightCounter[i]
      maxFreq = Math.max(maxFreq, weightCounter[i])
    }
    res[i] = sum - maxFreq
  }

  return res

  function getPath(cur: number, pre: number, to: number, weightCounter: Uint32Array): boolean {
    if (cur === to) return true
    for (const [next, w] of adjList[cur]) {
      if (next === pre) continue
      weightCounter[w]++
      if (getPath(next, cur, to, weightCounter)) return true
      weightCounter[w]--
    }
    return false
  }
}
