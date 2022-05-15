export {}

class SegmentTreeNode {
  left!: SegmentTreeNode
  right!: SegmentTreeNode
  isLazy = false
  lazyValue = 0
  value = 0
}

class SegmentTree {
  private readonly root: SegmentTreeNode = new SegmentTreeNode()
  private readonly size: number

  constructor(size: number) {
    this.size = size
  }

  update(left: number, right: number, delta: boolean): void {
    this._update(left, right, 0, this.size, this.root, delta)
  }

  query(left: number, right: number): number {
    return this._query(left, right, 0, this.size, this.root)
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
    delta: boolean
  ): void {
    if (L <= l && r <= R) {
      root.value = r - l + 1
      root.lazyValue = 1
      root.isLazy = true
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(root, l, r, mid)
    if (L <= mid) this._update(L, R, l, mid, root.left, delta)
    if (mid < R) this._update(L, R, mid + 1, r, root.right, delta)
    this.pushUp(root)
  }

  private _query(L: number, R: number, l: number, r: number, root: SegmentTreeNode): number {
    if (L <= l && r <= R) {
      return root.value
    }

    let res = 0
    const mid = Math.floor((l + r) / 2)
    this.pushDown(root, l, r, mid)
    if (L <= mid) res += this._query(L, R, l, mid, root.left)
    if (mid < R) res += this._query(L, R, mid + 1, r, root.right)
    return res
  }

  private pushDown(root: SegmentTreeNode, l: number, r: number, mid: number): void {
    !root.left && (root.left = new SegmentTreeNode())
    !root.right && (root.right = new SegmentTreeNode())
    if (root.isLazy) {
      root.left.isLazy = true
      root.right.isLazy = true
      root.left.lazyValue = root.lazyValue
      root.right.lazyValue = root.lazyValue
      root.left.value = mid - l + 1
      root.right.value = r - mid

      root.isLazy = false
      root.lazyValue = 0
    }
  }

  private pushUp(root: SegmentTreeNode): void {
    root.value = root.left.value + root.right.value
  }
}

class CountIntervals {
  private readonly tree = new SegmentTree(1e9 + 10)

  add(left: number, right: number): void {
    this.tree.update(left, right, true)
  }

  count(): number {
    return this.tree.queryAll()
  }
}
