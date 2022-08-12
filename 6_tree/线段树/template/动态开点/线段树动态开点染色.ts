// !两种颜色的区间染色不需要用懒标记记录更新状态 直接用节点值判断
class SegmentTreeNode {
  value = 0
  left?: SegmentTreeNode
  right?: SegmentTreeNode
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
    this.checkRange(left, right)
    this._update(left, right, this.lower, this.upper, this.root, target)
  }

  query(left: number, right: number): number {
    this.checkRange(left, right)
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
    // if (root.value === r - l + 1) return  // 如果线段树只加区间不减区间 可以这样优化
    if (L <= l && r <= R) {
      root.value = (r - l + 1) * target
      return
    }

    const mid = Math.floor((l + r) / 2)
    SegmentTree.pushDown(l, mid, r, root)
    if (L <= mid) this._update(L, R, l, mid, root.left!, target)
    if (mid < R) this._update(L, R, mid + 1, r, root.right!, target)
    SegmentTree.pushUp(root)
  }

  private _query(L: number, R: number, l: number, r: number, root: SegmentTreeNode): number {
    if (L <= l && r <= R) {
      return root.value
    }

    const mid = Math.floor((l + r) / 2)
    SegmentTree.pushDown(l, mid, r, root)
    let res = 0
    if (L <= mid) res += this._query(L, R, l, mid, root.left!)
    if (mid < R) res += this._query(L, R, mid + 1, r, root.right!)
    return res
  }

  private static pushDown(left: number, mid: number, right: number, root: SegmentTreeNode): void {
    !root.left && (root.left = new SegmentTreeNode())
    !root.right && (root.right = new SegmentTreeNode())
    if (root.value === right - left + 1) {
      root.left.value = mid - left + 1
      root.right.value = right - mid
    } else if (root.value === 0) {
      root.left.value = 0
      root.right.value = 0
    }
  }

  private static pushUp(root: SegmentTreeNode): void {
    root.value = root.left!.value + root.right!.value
  }

  private checkRange(l: number, r: number): void {
    if (l < this.lower || r > this.upper) {
      throw new RangeError(`[${l}, ${r}] out of range: [${this.lower}, ${this.upper}]`)
    }
  }
}

//! /////////////////////////////////////////////////////////////////////////////////////////
// 染色种类大于两种则需要懒标记
class SegmentTreeNodeWithLazy {
  left!: SegmentTreeNodeWithLazy
  right!: SegmentTreeNodeWithLazy
  isLazy = false
  lazyValue = 0
  value = 0
}

/**
 * @description 区间染色线段树，动态开点
 */
class SegmentTree2 {
  private readonly root: SegmentTreeNodeWithLazy = new SegmentTreeNodeWithLazy()
  private readonly lower: number
  private readonly upper: number

  constructor(lower = 0, upper = 1e9 + 10) {
    this.lower = lower
    this.upper = upper
  }

  update(left: number, right: number, target: number): void {
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
    root: SegmentTreeNodeWithLazy,
    target: number
  ): void {
    if (L <= l && r <= R) {
      root.value = (r - l + 1) * target
      root.lazyValue = target
      root.isLazy = true
      return
    }

    const mid = Math.floor((l + r) / 2)
    SegmentTree2.pushDown(l, mid, r, root)
    if (L <= mid) this._update(L, R, l, mid, root.left, target)
    if (mid < R) this._update(L, R, mid + 1, r, root.right, target)
    SegmentTree2.pushUp(root)
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
    SegmentTree2.pushDown(l, mid, r, root)
    let res = 0
    if (L <= mid) res += this._query(L, R, l, mid, root.left)
    if (mid < R) res += this._query(L, R, mid + 1, r, root.right)
    return res
  }

  private static pushDown(
    left: number,
    mid: number,
    right: number,
    root: SegmentTreeNodeWithLazy
  ): void {
    !root.left && (root.left = new SegmentTreeNodeWithLazy())
    !root.right && (root.right = new SegmentTreeNodeWithLazy())
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

  private static pushUp(root: SegmentTreeNodeWithLazy): void {
    root.value = root.left.value + root.right.value
  }
}

if (require.main === module) {
  class CountIntervals {
    private readonly tree = new SegmentTree(1, 1e9)

    // Adds the interval [left, right] to the set of intervals.
    add(left: number, right: number): void {
      this.tree.update(left, right, 1)
    }

    // Returns the number of integers that are present in at least one interval.
    count(): number {
      return this.tree.queryAll()
    }
  }
}

export {}
