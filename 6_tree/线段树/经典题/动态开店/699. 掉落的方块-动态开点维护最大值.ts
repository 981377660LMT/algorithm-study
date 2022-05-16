// 第 i 个掉落的方块（positions[i] = (left, side_length)）是正方形
// 1 <= positions.length <= 1000.
// 1 <= positions[i][0] <= 10^8.
// 1 <= positions[i][1] <= 10^6.

// 来源：力扣（LeetCode）
// 链接：https://leetcode.cn/problems/falling-squares
// 著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
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

  queryAll(): number {
    return this.root.value
  }

  private _update(
    L: number,
    R: number,
    left: number,
    right: number,
    root: SegmentTreeNode,
    target: number
  ): void {
    if (L <= left && right <= R) {
      root.isLazy = true
      root.value = Math.max(root.value, target)
      root.lazyValue = Math.max(root.lazyValue, target)
      return
    }

    this.pushDown(root)

    const mid = Math.floor((left + right) / 2)
    if (L <= mid) {
      this._update(L, R, left, mid, root.left, target)
    }
    if (R >= mid + 1) {
      this._update(L, R, mid + 1, right, root.right, target)
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
      res = Math.max(res, this._query(L, R, left, mid, root.left))
    }
    if (R >= mid + 1) {
      res = Math.max(res, this._query(L, R, mid + 1, right, root.right))
    }

    return res
  }

  private pushUp(root: SegmentTreeNode): void {
    root.value = Math.max(root.left.value, root.right.value)
  }

  private pushDown(root: SegmentTreeNode): void {
    if (root.left == undefined) root.left = new SegmentTreeNode()
    if (root.right == undefined) root.right = new SegmentTreeNode()
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
}

function fallingSquares(positions: number[][]): number[] {
  const res = Array<number>(positions.length).fill(0)
  const tree = new SegmentTree(2e8 + 10)
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
    [1, 2],
    [2, 3],
    [6, 1],
  ])
)
console.log(
  fallingSquares([
    [100, 100],
    [200, 100],
  ])
)
export {}
