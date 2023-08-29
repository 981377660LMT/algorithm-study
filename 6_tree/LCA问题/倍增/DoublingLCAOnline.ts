/* eslint-disable max-len */
/* eslint-disable no-param-reassign */

// !不预先给出整棵树,而是动态添加叶子节点,在线维护树节点的LCA和k级祖先.

class DoublingLCAOnline {
  readonly depth: Uint32Array
  readonly depthWeighted: number[]
  private readonly _n: number
  private readonly _root: number
  private readonly _bitLen: number
  private readonly _dp: Int32Array

  /**
   * 不预先给出整棵树,而是动态添加叶子节点,维护树节点的LCA和k级祖先.
   * @param n 树节点 0~n-1.
   * @param root 根节点.初始时只有一个根节点root,默认为0.
   */
  constructor(n: number, root = 0) {
    n++ // 防止越界
    this._n = n
    this._root = root
    this._bitLen = 32 - Math.clz32(n)
    this._dp = new Int32Array(this._bitLen * n).fill(-1)
    this.depth = new Uint32Array(n)
    this.depthWeighted = Array(n).fill(0)
  }

  /**
   * 在树中添加一条从parent到child的边.要求parent已经存在于树中.
   */
  addDirectedEdge(parent: number, child: number, weight = 1): void {
    if (parent !== this._root && !this.depth[parent]) {
      throw new Error(`parent ${parent} not exists`)
    }
    this.depth[child] = this.depth[parent] + 1
    this.depthWeighted[child] = this.depthWeighted[parent] + weight
    this._dp[child] = parent
    for (let i = 0; i < this._bitLen - 1; i++) {
      const pre = this._dp[i * this._n + child]
      if (pre === -1) break
      this._dp[(i + 1) * this._n + child] = this._dp[i * this._n + pre]
    }
  }

  /**
   * 查询节点node的第k个祖先(向上跳k步).如果不存在,返回-1.
   */
  kthAncestor(node: number, k: number): number {
    if (k > this.depth[node]) return -1
    let bit = 0
    while (k > 0) {
      if (k & 1) {
        node = this._dp[bit * this._n + node]
        if (node === -1) return -1
      }
      bit++
      k >>>= 1
    }
    return node
  }

  /**
   * 从 root 开始向上跳到指定深度 toDepth,toDepth<=depth[v],返回跳到的节点.
   */
  upToDepth(root: number, toDepth: number): number {
    if (toDepth >= this.depth[root]) return root
    for (let i = this._bitLen - 1; ~i; i--) {
      if ((this.depth[root] - toDepth) & (1 << i)) {
        root = this._dp[i * this._n + root]
      }
    }
    return root
  }

  lca(root1: number, root2: number): number {
    if (this.depth[root1] < this.depth[root2]) {
      root1 ^= root2
      root2 ^= root1
      root1 ^= root2
    }
    root1 = this.upToDepth(root1, this.depth[root2])
    if (root1 === root2) return root1
    for (let i = this._bitLen - 1; ~i; i--) {
      if (this._dp[i * this._n + root1] !== this._dp[i * this._n + root2]) {
        root1 = this._dp[i * this._n + root1]
        root2 = this._dp[i * this._n + root2]
      }
    }
    return this._dp[root1]
  }

  /**
   * 从start节点跳向target节点,跳过step个节点(0-indexed)
   * 返回跳到的节点,如果不存在这样的节点,返回-1
   */
  jump(start: number, target: number, step: number): number {
    const lca = this.lca(start, target)
    const dep1 = this.depth[start]
    const dep2 = this.depth[target]
    const deplca = this.depth[lca]
    const dist = dep1 + dep2 - 2 * deplca
    if (step > dist) return -1
    if (step <= dep1 - deplca) return this.kthAncestor(start, step)
    return this.kthAncestor(target, dist - step)
  }

  dist(root1: number, root2: number, weighted = false): number {
    if (weighted) {
      return this.depthWeighted[root1] + this.depthWeighted[root2] - 2 * this.depthWeighted[this.lca(root1, root2)]
    }
    return this.depth[root1] + this.depth[root2] - 2 * this.depth[this.lca(root1, root2)]
  }
}

export { DoublingLCAOnline }

if (require.main === module) {
  // https://leetcode.cn/problems/kth-ancestor-of-a-tree-node/

  class TreeAncestor {
    private readonly _lca: DoublingLCAOnline
    constructor(n: number, parent: number[]) {
      this._lca = new DoublingLCAOnline(n)
      const adjList: number[][] = Array(n)
      for (let i = 0; i < n; i++) adjList[i] = []
      for (let i = 1; i < n; i++) {
        const pre = parent[i]
        const cur = i
        adjList[pre].push(cur)
      }
      const dfs = (cur: number, pre: number) => {
        adjList[cur].forEach(next => {
          if (next !== pre) {
            this._lca.addDirectedEdge(cur, next)
            dfs(next, cur)
          }
        })
      }

      dfs(0, -1)
    }

    getKthAncestor(node: number, k: number): number {
      return this._lca.kthAncestor(node, k)
    }
  }

  /**
   * Your TreeAncestor object will be instantiated and called as such:
   * var obj = new TreeAncestor(n, parent)
   * var param_1 = obj.getKthAncestor(node,k)
   */
}
