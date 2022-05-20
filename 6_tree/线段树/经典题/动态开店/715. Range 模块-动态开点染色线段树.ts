class SegmentTreeNode {
  left!: SegmentTreeNode
  right!: SegmentTreeNode
  isLazy = false
  lazyValue = 0
  value = 0
}

/**
 * @description 区间染色线段树，动态开点
 */
class SegmentTree {
  private readonly root: SegmentTreeNode = new SegmentTreeNode()
  private readonly lower: number
  private readonly upper: number

  constructor(lower = 0, upper = 1e9 + 10) {
    this.lower = lower
    this.upper = upper
  }

  update(left: number, right: number, target: 0 | 1): void {
    this._update(left, right, this.lower, this.upper, this.root, target)
  }

  query(left: number, right: number): number {
    return this._query(left, right, this.lower, this.upper, this.root)
  }

  queryAll(): number {
    return this.root.value
  }

  private _update(
    L: number,
    R: number,
    l: number,
    r: number,
    root: SegmentTreeNode,
    target: 0 | 1
  ): void {
    if (L <= l && r <= R) {
      root.value = (r - l + 1) * target
      root.lazyValue = target
      root.isLazy = true
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(l, mid, r, root)
    if (L <= mid) this._update(L, R, l, mid, root.left, target)
    if (mid < R) this._update(L, R, mid + 1, r, root.right, target)
    this.pushUp(root)
  }

  private _query(L: number, R: number, l: number, r: number, root: SegmentTreeNode): number {
    if (L <= l && r <= R) {
      return root.value
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(l, mid, r, root)
    let res = 0
    if (L <= mid) res += this._query(L, R, l, mid, root.left)
    if (mid < R) res += this._query(L, R, mid + 1, r, root.right)
    return res
  }

  private pushDown(left: number, mid: number, right: number, root: SegmentTreeNode): void {
    !root.left && (root.left = new SegmentTreeNode())
    !root.right && (root.right = new SegmentTreeNode())
    if (root.isLazy) {
      root.left.isLazy = true
      root.right.isLazy = true
      root.left.lazyValue = root.lazyValue
      root.right.lazyValue = root.lazyValue
      root.left.value = root.lazyValue * (mid - left + 1)
      root.right.value = root.lazyValue * (right - mid)

      root.isLazy = false
      root.lazyValue = 0
    }
  }

  private pushUp(root: SegmentTreeNode): void {
    root.value = root.left.value + root.right.value
  }
}

class RangeModule {
  private readonly tree: SegmentTree = new SegmentTree(0, 1e9 + 10)

  /**
   * @description 添加 半开区间 [left, right)
   */
  addRange(left: number, right: number): void {
    this.tree.update(left, right - 1, 1)
  }

  /**
   * @description  只有在当前正在跟踪区间 [left, right) 中的每一个实数时，才返回 true
   */
  queryRange(left: number, right: number): boolean {
    return this.tree.query(left, right - 1) === right - left
  }

  /**
   * @description 停止跟踪 半开区间 [left, right)
   */
  removeRange(left: number, right: number): void {
    this.tree.update(left, right - 1, 0)
  }
}

export {}
