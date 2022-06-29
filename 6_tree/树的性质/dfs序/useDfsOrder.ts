/**
 * @param n 树节点个数
 * @param tree 无向图邻接表
 * @param root 树根节点
 */
function useDfsOrder(n: number, tree: Iterable<number>[], root = 0) {
  const ins = new Uint32Array(n + 1) // 子树中最小的结点序号
  const outs = new Uint32Array(n + 1) // 子树中最大的结点序号，即自己的id
  let dfsId = 1
  dfs(root, -1)

  // 求dfs序
  function dfs(cur: number, pre: number): void {
    ins[cur] = dfsId
    for (const next of tree[cur]) {
      if (next !== pre) dfs(next, cur)
    }
    outs[cur] = dfsId
    dfsId++
  }

  /**
   * @param root 求root所在子树映射到的区间
   * @returns [start, end] 1 <= start <= end <= n
   */
  function queryRange(root: number): [left: number, right: number] {
    return [ins[root], outs[root]]
  }

  /**
   * @param root 求root自身的dfsId
   * @returns dfsId 1 <= dfsId <= n
   */
  function queryId(root: number): number {
    return outs[root]
  }

  /**
   * @returns root是否是child的祖先
   * @description 应用:枚举边时给树的边定向
   *
   * ```js
   *  if (!isAncestor(e[0], e[1])) {
   *    [e[0], e[1]] = [e[1], e[0]]
   *  }
   * ```
   */
  function isAncestor(root: number, child: number): boolean {
    const [left1, right1] = [ins[root], outs[root]]
    const [left2, right2] = [ins[child], outs[child]]
    return left1 <= left2 && right2 <= right1
  }

  return { queryRange, queryId, isAncestor }
}

export { useDfsOrder }
