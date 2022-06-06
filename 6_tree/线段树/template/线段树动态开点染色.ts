// !染色不用带懒标记 因为信息存在区间值里
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
    this.pushDown(l, mid, r, root)
    if (L <= mid) this._update(L, R, l, mid, root.left!, target)
    if (mid < R) this._update(L, R, mid + 1, r, root.right!, target)
    this.pushUp(root)
  }

  private _query(L: number, R: number, l: number, r: number, root: SegmentTreeNode): number {
    if (L <= l && r <= R) {
      return root.value
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(l, mid, r, root)
    let res = 0
    if (L <= mid) res += this._query(L, R, l, mid, root.left!)
    if (mid < R) res += this._query(L, R, mid + 1, r, root.right!)
    return res
  }

  private pushDown(left: number, mid: number, right: number, root: SegmentTreeNode): void {
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

  private pushUp(root: SegmentTreeNode): void {
    root.value = root.left!.value + root.right!.value
  }

  private checkRange(l: number, r: number): void {
    if (l < this.lower || r > this.upper)
      throw new RangeError(`[${l}, ${r}] out of range: [${this.lower}, ${this.upper}]`)
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
