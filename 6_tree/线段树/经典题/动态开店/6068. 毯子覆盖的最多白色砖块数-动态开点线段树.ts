class SegmentTreeNode {
  left!: SegmentTreeNode
  right!: SegmentTreeNode
  isLazy = false
  lazyValue = 0
  value = 0
}

// 动态开点线段树，区间叠加更新
class SegmentTree {
  private readonly root = new SegmentTreeNode()
  private readonly size: number

  constructor(size: number) {
    this.size = size
  }

  update(left: number, right: number, delta: number): void {
    this._update(left, right, 1, this.size, this.root, delta)
  }

  query(left: number, right: number): number {
    return this._query(left, right, 1, this.size, this.root)
  }

  queryAll(): number {
    return this.root.value
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
      root.lazyValue += delta
      root.isLazy = true
      return
    }

    const mid = Math.floor((left + right) / 2)
    this.pushDown(left, mid, right, root)
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

    let res = 0

    const mid = Math.floor((left + right) / 2)
    this.pushDown(left, mid, right, root)
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

  private pushDown(left: number, mid: number, right: number, root: SegmentTreeNode): void {
    if (root.left == undefined) root.left = new SegmentTreeNode()
    if (root.right == undefined) root.right = new SegmentTreeNode()
    if (root.isLazy) {
      root.left.isLazy = true
      root.right.isLazy = true
      root.left.lazyValue += root.lazyValue
      root.right.lazyValue += root.lazyValue
      root.left.value += root.lazyValue * (mid - left + 1)
      root.right.value += root.lazyValue * (right - mid)

      root.isLazy = false
      root.lazyValue = 0
    }
  }
}

function maximumWhiteTiles(tiles: number[][], carpetLen: number): number {
  // 注意js不要直接开1e9 容易MLE
  const tree = new SegmentTree(Math.max(...tiles.flat()) + 10)
  for (const [left, right] of tiles) tree.update(left, right, 1)
  let res = 0
  for (const [left] of tiles) res = Math.max(res, tree.query(left, left + carpetLen - 1))
  return res
}

export {}
