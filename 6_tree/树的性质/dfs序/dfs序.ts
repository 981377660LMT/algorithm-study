/**
 * @param n 树节点个数
 * @param adjList 无向图邻接表
 */
function useDfsOrder(n: number, adjList: number[][]) {
  const starts = Array<number>(n + 1).fill(0) // 子树中最小的结点序号
  const ends = Array<number>(n + 1).fill(0) // 子树中最大的结点序号，即自己的id
  let dfsId = 1
  dfs(0, -1)

  // 求dfs序
  function dfs(cur: number, pre: number): void {
    starts[cur] = dfsId
    for (const next of adjList[cur]) {
      if (next !== pre) dfs(next, cur)
    }
    ends[cur] = dfsId
    dfsId++
  }

  /**
   * @param root 求root所在子树映射到的区间
   * @returns [start, end] 1 <= start <= end <= n
   */
  function queryRange(root: number): [left: number, right: number] {
    return [starts[root], ends[root]]
  }

  /**
   * @param root 求root自身的dfsId
   * @returns dfsId 1 <= dfsId <= n
   */
  function queryId(root: number): number {
    return ends[root]
  }

  return { queryRange, queryId }
}

export { useDfsOrder }
