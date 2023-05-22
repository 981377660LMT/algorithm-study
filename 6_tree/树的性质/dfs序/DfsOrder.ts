//            0 [0,5)
//           /       \
//          /         \
//        1 [0,3)      2 [3,4)
//       /    \
//      /      \
//    3 [0,1)   4[1,2)

class DfsOrder {
  private readonly _tree: number[][]
  private _dfsId = 0
  ins: Uint32Array
  outs: Uint32Array

  /**
   * @param n 树节点个数
   * @param tree 无向图邻接表
   * @param root 树根节点
   */
  constructor(n: number, tree: number[][], root = 0) {
    this._tree = tree
    this.ins = new Uint32Array(n) // 子树中最小的结点序号
    this.outs = new Uint32Array(n) // 子树中最大的结点序号，即自己的id
    this._dfs(root, -1)
  }

  /**
   * @param curRoot 求root所在子树映射到的左闭右开区间
   * @returns [start, end). `0 <= start < end <= n`
   */
  queryRange(curRoot: number): [number, number] {
    return [this.ins[curRoot], this.outs[curRoot] + 1]
  }

  /**
   * @param curRoot 求root自身的dfsId
   * @returns `0 <= dfsId < n`
   */
  queryId(curRoot: number): number {
    return this.outs[curRoot]
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
  isAncestor(root1: number, root2: number): boolean {
    const left1 = this.ins[root1]
    const right1 = this.outs[root1]
    const left2 = this.ins[root2]
    const right2 = this.outs[root2]
    return left1 <= left2 && right2 <= right1
  }

  private _dfs(cur: number, pre: number): void {
    this.ins[cur] = this._dfsId
    this._tree[cur].forEach(next => {
      if (next !== pre) this._dfs(next, cur)
    })
    this.outs[cur] = this._dfsId
    this._dfsId++
  }
}

export { DfsOrder }

if (require.main === module) {
  const n = 5
  const edges = [
    [0, 1],
    [0, 2],
    [1, 3],
    [1, 4]
  ]
  const tree: number[][] = Array(n)
  for (let i = 0; i < n; i++) {
    tree[i] = []
  }
  edges.forEach(([pre, cur]) => {
    tree[pre].push(cur)
    tree[cur].push(pre)
  })

  const order = new DfsOrder(n, tree, 0)
  for (let i = 0; i < n; i++) {
    console.log(i, order.queryRange(i), order.queryId(i))
  }
}
