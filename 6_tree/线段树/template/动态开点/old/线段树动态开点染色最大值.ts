/* eslint-disable class-methods-use-this */
/* eslint-disable no-inner-declarations */
class SegmentTreeNodeWithLazy {
  left!: SegmentTreeNodeWithLazy
  right!: SegmentTreeNodeWithLazy
  isLazy = false
  lazyValue = 0
  value = 0
}

/**
 * 线段树区间染色最大值RMQ
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

  update(left: number, right: number, target: number): void {
    if (left < this._lower) left = this._lower
    if (right > this._upper) right = this._upper
    if (left > right) return
    this._update(left, right, this._lower, this._upper, this._root, target)
  }

  query(left: number, right: number): number {
    if (left < this._lower) left = this._lower
    if (right > this._upper) right = this._upper
    if (left > right) return 0
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
    target: number
  ): void {
    if (L <= l && r <= R) {
      root.value = Math.max(root.value, target)
      root.lazyValue = Math.max(root.lazyValue, target)
      root.isLazy = true
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(l, mid, r, root)
    if (L <= mid) this._update(L, R, l, mid, root.left, target)
    if (mid < R) this._update(L, R, mid + 1, r, root.right, target)
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
    let res = 0
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
      root.left.lazyValue = Math.max(root.left.lazyValue, root.lazyValue)
      root.right.lazyValue = Math.max(root.right.lazyValue, root.lazyValue)
      root.left.value = Math.max(root.left.value, root.lazyValue)
      root.right.value = Math.max(root.right.value, root.lazyValue)

      root.isLazy = false
      root.lazyValue = 0
    }
  }

  private pushUp(root: SegmentTreeNodeWithLazy): void {
    root.value = Math.max(root.left.value, root.right.value)
  }
}

if (require.main === module) {
  function fallingSquares(positions: number[][]): number[] {
    const res = Array<number>(positions.length).fill(0)
    const tree = new MaxSegmentTree(1, 2e8)
    for (const [i, [left, size]] of positions.entries()) {
      const right = left + size - 1

      const preHeihgt = tree.query(left, right)
      tree.update(left, right, preHeihgt + size)
      res[i] = tree.queryAll()
    }

    return res
  }

  console.log(
    fallingSquares([
      [1, 5],
      [2, 2],
      [7, 5]
    ])
  )
}

export {}
