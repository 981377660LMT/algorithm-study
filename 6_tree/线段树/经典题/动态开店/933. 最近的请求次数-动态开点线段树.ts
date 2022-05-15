// 写一个 RecentCounter 类来计算特定时间范围内最近的请求

// 如果时间戳递增 直接用Queue
// 如果时间戳不递增 直接用BIT/SegmentTree

class SegmentTreeNode {
  left!: SegmentTreeNode
  right!: SegmentTreeNode
  isLazy = false // 因为是单点修改，这里不要
  lazyValue = 0 // 因为是单点修改，这里不要
  value = 0
}

// 动态开点线段树
class SegmentTree {
  private readonly root = new SegmentTreeNode()
  private readonly size: number

  constructor(size: number) {
    this.size = size
  }

  // 1 <= t <= 1e9
  update(left: number, right: number, delta: number): void {
    this._update(left, right, 1, this.size, this.root, delta)
  }

  query(left: number, right: number): number {
    return this._query(left, right, 1, this.size, this.root)
  }

  private _update(
    L: number,
    R: number,
    left: number,
    right: number,
    root: SegmentTreeNode,
    delta: number
  ): void {
    if (L <= left && right <= R) {
      root.value += (right - left + 1) * delta
      return
    }

    this.pushDown(root)

    const mid = Math.floor((left + right) / 2)
    if (L <= mid) {
      this._update(L, R, left, mid, root.left, delta)
    }
    if (R >= mid + 1) {
      this._update(L, R, mid + 1, right, root.right, delta)
    }

    this.pushUp(root)
  }

  private _query(L: number, R: number, left: number, right: number, root: SegmentTreeNode): number {
    if (L <= left && right <= R) {
      return root.value
    }

    this.pushDown(root)
    let res = 0

    const mid = Math.floor((left + right) / 2)
    if (L <= mid) {
      res += this._query(L, R, left, mid, root.left)
    }
    if (R >= mid + 1) {
      res += this._query(L, R, mid + 1, right, root.right)
    }

    return res
  }

  private pushUp(root: SegmentTreeNode): void {
    root.value = root.left.value + root.right.value
  }

  private pushDown(root: SegmentTreeNode): void {
    if (root.left == undefined) root.left = new SegmentTreeNode()
    if (root.right == undefined) root.right = new SegmentTreeNode()
  }
}

class RecentCounter {
  private readonly tree = new SegmentTree(1e9 + 10)

  ping(t: number): number {
    this.tree.update(t, t, 1)
    return this.tree.query(t - 3000, t)
  }
}

export {}
