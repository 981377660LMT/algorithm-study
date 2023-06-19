import { SegmentTreeRangeUpdateRangeQuery } from '../../template/atcoder_segtree/SegmentTreeRangeUpdateRangeQuery'

class BookMyShow {
  private readonly _row: number
  private readonly _col: number
  // 维护每行的剩余座位数和最大剩余座位数
  private readonly _tree: SegmentTreeRangeUpdateRangeQuery<[max: number, sum: number], number>

  constructor(n: number, m: number) {
    const leaves = Array(n)
    for (let i = 0; i < n; i++) leaves[i] = [m, m]
    this._row = n
    this._col = m
    this._tree = new SegmentTreeRangeUpdateRangeQuery<[max: number, sum: number], number>(leaves, {
      e() {
        return [0, 0]
      },
      id() {
        return 0
      },
      op(left, right) {
        return [Math.max(left[0], right[0]), left[1] + right[1]]
      },
      mapping(lazy, data) {
        return [data[0] + lazy, data[1] + lazy]
      },
      composition(parentLazy, childLazy) {
        return parentLazy + childLazy
      }
    })
  }

  gather(k: number, maxRow: number): number[] {
    const first = this._tree.maxRight(0, e => e[0] < k) // !找到第一个空座位>=k的行
    if (first > maxRow) return []
    const used = this._col - this._tree.query(first, first + 1)[1]
    this._tree.update(first, first + 1, -k)
    return [first, used]
  }

  scatter(k: number, maxRow: number): boolean {
    const remain = this._tree.query(0, maxRow + 1)[1]
    if (remain < k) return false

    let first = this._tree.maxRight(0, e => e[1] === 0) // !找到第一个未坐满的行
    while (k > 0) {
      const remain = this._tree.query(first, first + 1)[1]
      const min_ = Math.min(k, remain)
      this._tree.update(first, first + 1, -min_)
      k -= min_
      first++
    }

    return true
  }
}
