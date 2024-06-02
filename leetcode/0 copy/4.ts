export {}

const INF = 2e9 // !超过int32使用2e15

// 给你一个正整数 days，表示员工可工作的总天数（从第 1 天开始）。另给你一个二维数组 meetings，长度为 n，其中 meetings[i] = [start_i, end_i] 表示第 i 次会议的开始和结束天数（包含首尾）。

// 返回员工可工作且没有安排会议的天数。

// 注意：会议时间可能会有重叠。
function countDays(days: number, meetings: number[][]): number {
  const seg = new SegmentTree()
  meetings.forEach(([l, r]) => seg.update(l, r, 1))
  return days - seg.queryAll()
}

class SegmentTreeNode {
  value = 0
  left?: SegmentTreeNode
  right?: SegmentTreeNode
}

/**
 * @description 区间染色线段树，动态开点
 */
class SegmentTree {
  private readonly _root: SegmentTreeNode = new SegmentTreeNode()
  private readonly _lower: number
  private readonly _upper: number

  constructor(lower = 0, upper = 1e9 + 10) {
    this._lower = lower
    this._upper = upper
  }

  update(left: number, right: number, target: 0 | 1): void {
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
}
