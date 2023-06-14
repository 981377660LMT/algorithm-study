class SegmentTreeNodeWithLazy {
  left!: SegmentTreeNodeWithLazy
  right!: SegmentTreeNodeWithLazy
  isLazy = false
  lazyValue = 0
  value = 0 // 结点最大值
}

/**
 * 线段树区间叠加最大值RMQ
 *
 * 超出范围返回0
 */
class MaxSegmentTree {
  private readonly _root: SegmentTreeNodeWithLazy = new SegmentTreeNodeWithLazy()
  private readonly _lower: number
  private readonly _upper: number

  constructor(lower = 0, upper = 1e9 + 10) {
    this._lower = lower
    this._upper = upper
  }

  update(left: number, right: number, delta: number): void {
    if (left < this._lower) left = this._lower
    if (right > this._upper) right = this._upper
    if (left > right) return
    this._update(left, right, this._lower, this._upper, this._root, delta)
  }

  query(left: number, right: number): number {
    if (left < this._lower) left = this._lower
    if (right > this._upper) right = this._upper
    if (left > right) return 0 // !超出范围返回0
    return this._query(left, right, this._lower, this._upper, this._root)
  }

  queryAll(): number {
    return this._root.value
  }

  private _update(
    L: number,
    R: number,
    l: number,
    r: number,
    root: SegmentTreeNodeWithLazy,
    delta: number
  ): void {
    if (L <= l && r <= R) {
      root.value += delta
      root.lazyValue += delta
      root.isLazy = true
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(l, mid, r, root)
    if (L <= mid) this._update(L, R, l, mid, root.left, delta)
    if (mid < R) this._update(L, R, mid + 1, r, root.right, delta)
    this.pushUp(root)
  }

  private _query(
    L: number,
    R: number,
    l: number,
    r: number,
    root: SegmentTreeNodeWithLazy
  ): number {
    if (L <= l && r <= R) {
      return root.value
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(l, mid, r, root)
    let res = 0 // !默认返回0
    if (L <= mid) res = Math.max(res, this._query(L, R, l, mid, root.left))
    if (mid < R) res = Math.max(res, this._query(L, R, mid + 1, r, root.right))
    return res
  }

  private pushDown(left: number, mid: number, right: number, root: SegmentTreeNodeWithLazy): void {
    !root.left && (root.left = new SegmentTreeNodeWithLazy())
    !root.right && (root.right = new SegmentTreeNodeWithLazy())
    if (root.isLazy) {
      root.left.isLazy = true
      root.right.isLazy = true
      root.left.lazyValue += root.lazyValue
      root.right.lazyValue += root.lazyValue
      root.left.value += root.lazyValue
      root.right.value += root.lazyValue

      root.isLazy = false
      root.lazyValue = 0
    }
  }

  private pushUp(root: SegmentTreeNodeWithLazy): void {
    root.value = Math.max(root.left.value, root.right.value)
  }
}

if (require.main === module) {
  class MyCalendarThree {
    private readonly tree = new MaxSegmentTree(0, 1e9)

    book(start: number, end: number): number {
      this.tree.update(start, end - 1, 1) // 注意这里
      return this.tree.queryAll()
    }
  }
}

export {}
