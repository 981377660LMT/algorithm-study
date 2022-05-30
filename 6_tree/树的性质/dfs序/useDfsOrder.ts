/**
 * @param n 树节点个数
 * @param tree 无向图邻接表
 * @param root 树根节点
 */
function useDfsOrder(n: number, tree: Iterable<number>[], root = 0) {
  const starts = new Uint32Array(n + 1) // 子树中最小的结点序号
  const ends = new Uint32Array(n + 1) // 子树中最大的结点序号，即自己的id
  let dfsId = 1
  dfs(root, -1)

  // 求dfs序
  function dfs(cur: number, pre: number): void {
    starts[cur] = dfsId
    for (const next of tree[cur]) {
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
