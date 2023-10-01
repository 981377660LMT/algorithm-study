/* eslint-disable no-inner-declarations */

// TODO: 速度较慢(与python差不多)，建议使用python
// !比较慢的地方是数组的push和pop，可以用静态数组+ptr来优化

/**
 * Tarjan 求有向图的强联通分量.用于缩点成拓扑图.
 * @param n 图的节点数. 节点编号从 0 到 n - 1.
 * @param graph 图的邻接表表示.
 * @returns
 * [groups, belong].
 * groups 是每个 scc 组里包含的点.**`每个 group 之间按照拓扑序排列`**.
 * belong 是每个节点所属的强联通分量的编号(从0开始).
 */
function findSCC(
  n: number,
  graph: ArrayLike<ArrayLike<number>>
): [groups: number[][], belong: Uint32Array] {
  const dfn = new Uint32Array(n)
  let time = 0
  const stack = new Uint32Array(n)
  let stackTop = 0
  const inStack = new Uint8Array(n)
  const groups: number[][] = []

  const dfs = (cur: number): number => {
    time++
    dfn[cur] = time
    let curLow = time
    stack[stackTop++] = cur
    inStack[cur] = 1
    const nexts = graph[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (!dfn[next]) {
        curLow = Math.min(curLow, dfs(next))
      } else if (inStack[next] && dfn[next] < curLow) {
        curLow = dfn[next]
      }
    }
    if (dfn[cur] === curLow) {
      const group: number[] = []
      while (true) {
        const top = stack[--stackTop]
        group.push(top)
        inStack[top] = 0
        if (top === cur) break
      }
      groups.push(group)
    }
    return curLow
  }

  for (let i = 0; i < n; i++) {
    if (!dfn[i]) dfs(i)
  }

  groups.reverse()
  const belong = new Uint32Array(n)
  for (let i = 0; i < groups.length; i++) {
    const group = groups[i]
    for (let j = 0; j < group.length; j++) {
      belong[group[j]] = i
    }
  }

  return [groups, belong]
}

/**
 * 有向图的强联通分量缩点成拓扑图.
 * 缩点后得到了一张 DAG，点的编号范围为 `0 ~ {@link groups.length} - 1`.
 * @param graph 图的邻接表表示.
 * @param groups 每个 scc 组里包含的点.
 * @param belong 每个节点所属的强联通分量的编号(从0开始).
 * @param f 回调函数，入参为 `(from, fromSccId, to, toSccId)`.
 * @returns
 * [dag, indeg].
 * dag 是缩点后的拓扑图的邻接表表示.
 * indeg 是每个点的入度.
 */
function toDAG(
  graph: ArrayLike<ArrayLike<number>>,
  groups: ArrayLike<ArrayLike<number>>,
  belong: ArrayLike<number>,
  f?: (from: number, fromSccId: number, to: number, toSccId: number) => void
): [dag: number[][], indeg: Uint32Array] {
  const m = groups.length
  const dag: number[][] = Array(m)
  for (let i = 0; i < m; i++) dag[i] = []
  const indeg = new Uint32Array(m)
  const visitedEdge = new Set<number>() // !去除重边

  for (let cur = 0; cur < graph.length; cur++) {
    const curId = belong[cur]
    const nexts = graph[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      const nextId = belong[next]
      if (curId !== nextId) {
        const hash = curId * m + nextId
        if (visitedEdge.has(hash)) continue
        visitedEdge.add(hash)
        dag[curId].push(nextId)
        indeg[nextId]++
      }
      f && f(cur, curId, next, nextId)
    }
  }

  return [dag, indeg]
}

export { findSCC, findSCC as getSCC, toDAG }

if (require.main === module) {
  // 100075. 有向图访问计数
  // 给定一个有向图，对每个结点 0 <= i < n，统计从 i 出发可以访问到的结点数量。
  // https://leetcode.cn/problems/count-visited-nodes-in-a-directed-graph/

  function countVisitedNodes(edges: number[]): number[] {
    const n = edges.length
    const adjList: number[][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    for (let i = 0; i < n; i++) adjList[i].push(edges[i])

    const [groups, belong] = findSCC(n, adjList)
    const [dag, indeg] = toDAG(adjList, groups, belong)

    const cache = new Int32Array(dag.length).fill(-1)
    const dfs = (curId: number): number => {
      if (cache[curId] !== -1) return cache[curId]
      let res = groups[curId].length
      const nexts = dag[curId]
      for (let i = 0; i < nexts.length; i++) {
        res += dfs(nexts[i])
      }
      cache[curId] = res
      return res
    }

    const res: number[] = Array(n)
    for (let i = 0; i < n; i++) res[i] = dfs(belong[i])
    return res
  }
}
