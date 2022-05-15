// 线段树解法
class SegmentTreeNode {
  left!: SegmentTreeNode
  right!: SegmentTreeNode
  isLazy = false
  lazyValue = false
  value = false
}

/**
 * @description 线段树懒更新模板
 */
class SegmentTree {
  private readonly root: SegmentTreeNode = new SegmentTreeNode()
  private readonly size: number

  constructor(size: number) {
    this.size = size
  }

  update(left: number, right: number, delta: boolean): void {
    this._update(left, right, 0, this.size, this.root, delta)
  }

  query(left: number, right: number): boolean {
    return this._query(left, right, 0, this.size, this.root)
  }

  queryAll(): boolean {
    return this.root.value
  }

  private _update(
    L: number,
    R: number,
    l: number,
    r: number,
    root: SegmentTreeNode,
    delta: boolean
  ): void {
    if (L <= l && r <= R) {
      root.value = delta
      root.lazyValue = delta
      root.isLazy = true
      return
    }

    this.pushDown(root)
    const mid = (l + r) >> 1
    if (L <= mid) this._update(L, R, l, mid, root.left, delta)
    if (mid < R) this._update(L, R, mid + 1, r, root.right, delta)
    this.pushUp(root)
  }

  private _query(L: number, R: number, l: number, r: number, root: SegmentTreeNode): boolean {
    if (L <= l && r <= R) {
      return root.value
    }

    this.pushDown(root)
    let res = true
    const mid = (l + r) >> 1
    if (L <= mid) res = res && this._query(L, R, l, mid, root.left)
    if (mid < R) res = res && this._query(L, R, mid + 1, r, root.right)
    return res
  }

  /**
   * @param root 向下传递懒标记和懒更新的值 `isLazy`, `lazyValue`，并用 `lazyValue` 更新子区间的值
   */
  private pushDown(root: SegmentTreeNode): void {
    !root.left && (root.left = new SegmentTreeNode())
    !root.right && (root.right = new SegmentTreeNode())
    if (root.isLazy) {
      root.left.isLazy = true
      root.right.isLazy = true
      root.left.lazyValue = root.lazyValue
      root.right.lazyValue = root.lazyValue
      root.left.value = root.value
      root.right.value = root.value

      root.isLazy = false
      root.lazyValue = false
    }
  }

  /**
   * @param root 用子节点更新父节点的值
   */
  private pushUp(root: SegmentTreeNode): void {
    root.value = root.left.value && root.right.value
  }
}

class RangeModule {
  private readonly tree: SegmentTree = new SegmentTree(1e9 + 10)

  /**
   * @description 添加 半开区间 [left, right)
   */
  addRange(left: number, right: number): void {
    this.tree.update(left, right - 1, true)
  }

  /**
   * @description  只有在当前正在跟踪区间 [left, right) 中的每一个实数时，才返回 true
   */
  queryRange(left: number, right: number): boolean {
    return this.tree.query(left, right - 1)
  }

  /**
   * @description 停止跟踪 半开区间 [left, right)
   */
  removeRange(left: number, right: number): void {
    this.tree.update(left, right - 1, false)
  }
}

export {}
