interface ITreePath {
  readonly from: number
  readonly to: number
  readonly lca: number
  readonly length: number

  /** 从路径的起点开始，第k个节点(0-indexed).*/
  kthNodeOnPath(k: number): number
  onPath(node: number): boolean
  hasIntersection(other: ITreePath): boolean
  /** 求两条路径的交, 返回相交线段的两个端点.无交点则返回undefined. */
  getIntersection(other: ITreePath): { p1: number; p2: number } | undefined
  countIntersection(other: ITreePath): number
}

class TreePath implements ITreePath {
  private readonly _depth: ArrayLike<number>
  private readonly _kthAncestorFn: (node: number, k: number) => number
  private readonly _lcaFn: (node1: number, node2: number) => number
  private readonly _from: number
  private readonly _to: number
  private readonly _lca: number

  constructor(
    from: number,
    to: number,
    treeProps: {
      depth: ArrayLike<number>
      kthAncestorFn: (node: number, k: number) => number
      lcaFn: (node1: number, node2: number) => number
    }
  ) {
    this._depth = treeProps.depth
    this._kthAncestorFn = treeProps.kthAncestorFn
    this._lcaFn = treeProps.lcaFn
    this._from = from
    this._to = to
    this._lca = this._lcaFn(from, to)
  }

  kthNodeOnPath(k: number): number {
    if (k <= this._depth[this._from] - this._depth[this._lca]) {
      return this._kthAncestorFn(this._from, k)
    }
    return this._kthAncestorFn(this._to, this.length - k)
  }

  onPath(node: number): boolean {
    const lcaFn = this._lcaFn
    return (
      lcaFn(node, this._lca) === this._lca &&
      (lcaFn(node, this._from) === node || lcaFn(node, this._to) === node)
    )
  }

  hasIntersection(other: ITreePath): boolean {
    return this.onPath(other.lca) || other.onPath(this.lca)
  }

  /**
   * 给出两条路径`(a, b)`和`(c, d)`.
   * 四个点两两求lca，得到 `x1 = lca(a, c), x2 = lca(a, d), x3 = lca(b, c), x4 = lca(b, d)`.
   * 再从这四个点中找到深度最大的两个点`p1`和`p2`.
   * - 若`p1!==p2`，则两条路径相交，交点为`p1`和`p2`.
   * - 若`p1===p2`，且`p1`的深度小于`lca(a,b)`或小于`lca(c,d)`，则两条路径无交点.否则交点为`p1`.
   */
  getIntersection(other: ITreePath): { p1: number; p2: number } | undefined {
    const a = this._from
    const b = this._to
    const c = other.from
    const d = other.to
    const lcaFn = this._lcaFn
    const depth = this._depth
    const x1 = lcaFn(a, c)
    const x2 = lcaFn(a, d)
    const x3 = lcaFn(b, c)
    const x4 = lcaFn(b, d)
    let p1 = x1
    let p2 = x2
    if (depth[x2] > depth[p1]) {
      p2 = p1
      p1 = x2
    }
    const update = (x: number): void => {
      const curDepth = depth[x]
      if (curDepth > depth[p1]) {
        p2 = p1
        p1 = x
      } else if (curDepth > depth[p2]) {
        p2 = x
      }
    }
    update(x3)
    update(x4)
    const lca1 = this.lca
    const lca2 = other.lca
    if (p1 !== p2) return { p1, p2 }
    if (depth[p1] < depth[lca1] || depth[p1] < depth[lca2]) return undefined
    return { p1, p2 }
  }

  countIntersection(other: ITreePath): number {
    const res = this.getIntersection(other)
    if (res === undefined) return 0
    const { p1, p2 } = res
    if (p1 === p2) return 1
    return this._depth[p1] + this._depth[p2] - 2 * this._depth[this._lca] + 1
  }

  get from(): number {
    return this._from
  }

  get to(): number {
    return this._to
  }

  get lca(): number {
    return this._lca
  }

  get length(): number {
    return this._depth[this._from] + this._depth[this._to] - 2 * this._depth[this._lca]
  }
}

export { ITreePath, TreePath }
