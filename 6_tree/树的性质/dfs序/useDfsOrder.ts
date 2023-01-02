/**
 * @param n 树节点个数
 * @param tree 无向图邻接表
 * @param root 树根节点
 */
function useDfsOrder(n: number, tree: number[][], root = 0) {
  const ins = new Uint32Array(n + 10) // 子树中最小的结点序号
  const outs = new Uint32Array(n + 10) // 子树中最大的结点序号，即自己的id
  let dfsId = 1
  dfs(root, -1)

  // 求dfs序
  function dfs(cur: number, pre: number): void {
    ins[cur] = dfsId
    tree[cur].forEach(next => {
      if (next !== pre) dfs(next, cur)
    })
    outs[cur] = dfsId
    dfsId++
  }

  /**
   * @param curRoot 求root所在子树映射到的区间
   * @returns [start, end] 1 <= start <= end <= n
   */
  function queryRange(curRoot: number): [left: number, right: number] {
    return [ins[curRoot], outs[curRoot]]
  }

  /**
   * @param curRoot 求root自身的dfsId
   * @returns dfsId 1 <= dfsId <= n
   */
  function queryId(curRoot: number): number {
    return outs[curRoot]
  }

  /**
   * @returns root1是否是root2的祖先
   * @description 应用:枚举边时给树的边定向
   *
   * ```ts
   *  if (!isAncestor(e[0], e[1])) {
   *    [e[0], e[1]] = [e[1], e[0]]
   *  }
   * ```
   */
  function isAncestor(root1: number, root2: number): boolean {
    const left1 = ins[root1]
    const right1 = outs[root1]
    const left2 = ins[root2]
    const right2 = outs[root2]
    return left1 <= left2 && right2 <= right1
  }

  return { queryRange, queryId, isAncestor }
}

export { useDfsOrder }
