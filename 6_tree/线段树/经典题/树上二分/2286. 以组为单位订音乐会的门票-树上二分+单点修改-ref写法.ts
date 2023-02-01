/* eslint-disable max-len */
// 用一个线段树维护区间最值、区间和
// 需要找第一个行>=k，类似主席树找mex，检查区间max线段树上二分即可(到某个节点检查是否成立)

// https://leetcode.cn/problems/booking-concert-tickets-in-groups/submissions/
// 0. 线段树维护区间的 sum 和区间的 max
// 1. gather 操作里需要找到第一个 `剩余座位>=k` 的行，
//    可以通过查询左右区间里的 `max` 来确定向左子树递归还是向右子树递归，比较简单的写法是 `在入口判断区间max是否不小于k`
// 2. scatter 操作里需要按行号从小到大依次填充k个位置，
//    这个过程就是二叉树的后序遍历，先遍历左子树再遍历右子树找到最左边的结点，先填充完小的行再填充大的行
// 3. gather 操作需要返回 `[第一个行，第一个列]`，
//    比较方便的做法是在update操作里多传一个引用类型的参数 `resRef`，在递归过程中直接修改返回值
// 4. 由于每次调用api都是单点修改单点查询，所以不需要懒标记和 pushDown 操作

// !都是到叶子节点的单点更新 不用加懒标记
class SegmentTree {
  private readonly _row: number
  private readonly _col: number
  private readonly _sums: number[] // 每行剩余座位数
  private readonly _maxs: Uint32Array // 区间最大值

  constructor(rowSize: number, colSize: number) {
    this._row = rowSize
    this._col = colSize
    this._sums = Array(rowSize << 2).fill(0)
    this._maxs = new Uint32Array(rowSize << 2)
    this._build(1, 1, rowSize, colSize)
  }

  queryMax(l: number, r: number): number {
    return this._queryMax(1, l, r, 1, this._row)
  }

  querySum(l: number, r: number): number {
    return this._querySum(1, l, r, 1, this._row)
  }

  // !找到最小的行满足sum>=k的位置
  updateOneRow(
    l: number,
    r: number,
    delta: number,
    resRef: { resRow: number; resCol: number }
  ): void {
    this._updateOneRow(1, l, r, 1, this._row, delta, resRef)
  }

  // !从最小的行开始填充delta
  updateManyRows(l: number, r: number, remainRef: { value: number }): void {
    this._updateManyRows(1, l, r, 1, this._row, remainRef)
  }

  private _build(rt: number, l: number, r: number, colSize: number): void {
    if (l === r) {
      this._sums[rt] = colSize
      this._maxs[rt] = colSize
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._build(rt << 1, l, mid, colSize)
    this._build((rt << 1) | 1, mid + 1, r, colSize)
    this._pushUp(rt)
  }

  private _queryMax(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._maxs[rt]
    const mid = Math.floor((l + r) / 2)
    let res = 0
    if (L <= mid) res = Math.max(res, this._queryMax(rt << 1, L, R, l, mid))
    if (mid < R) res = Math.max(res, this._queryMax((rt << 1) | 1, L, R, mid + 1, r))
    return res
  }

  private _querySum(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._sums[rt]
    const mid = Math.floor((l + r) / 2)
    let res = 0
    if (L <= mid) res += this._querySum(rt << 1, L, R, l, mid)
    if (mid < R) res += this._querySum((rt << 1) | 1, L, R, mid + 1, r)
    return res
  }

  // 找到第一个行空座位>=k的位置 即二叉树前序遍历(先遍历左子树再遍历右子树)
  private _updateOneRow(
    rt: number,
    L: number,
    R: number,
    l: number,
    r: number,
    remain: number,
    resRef: { resRow: number; resCol: number } // 返回值
  ): void {
    if (this._maxs[rt] < remain) return // 二分看max是否满足 简单的写法就是在入口判断

    // 单点修改 找到答案了
    if (l === r) {
      resRef.resRow = l - 1
      resRef.resCol = this._col - this._sums[rt]
      this._sums[rt] -= remain
      this._maxs[rt] -= remain
      return
    }

    const mid = Math.floor((l + r) / 2)
    if (resRef.resRow === -1 && L <= mid) {
      this._updateOneRow(rt << 1, L, R, l, mid, remain, resRef)
    }
    if (resRef.resRow === -1 && mid < R) {
      this._updateOneRow((rt << 1) | 1, L, R, mid + 1, r, remain, resRef)
    }
    this._pushUp(rt)
  }

  // scatter 尽量往左填充 即二叉树前序遍历(先遍历左子树再遍历右子树)
  private _updateManyRows(
    rt: number,
    L: number,
    R: number,
    l: number,
    r: number,
    remainRef: { value: number }
  ): void {
    // 单点修改
    if (l === r) {
      const remain = Math.min(remainRef.value, this._sums[rt])
      this._sums[rt] -= remain
      this._maxs[rt] -= remain
      remainRef.value -= remain
      return
    }

    const mid = Math.floor((l + r) / 2)
    if (remainRef.value > 0 && L <= mid) {
      this._updateManyRows(rt << 1, L, R, l, mid, remainRef)
    }
    if (remainRef.value > 0 && mid < R) {
      this._updateManyRows((rt << 1) | 1, L, R, mid + 1, r, remainRef)
    }
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this._sums[rt] = this._sums[rt << 1] + this._sums[(rt << 1) | 1]
    this._maxs[rt] = Math.max(this._maxs[rt << 1], this._maxs[(rt << 1) | 1])
  }
}

class BookMyShow {
  private readonly _tree: SegmentTree
  private readonly _row: number
  private readonly _col: number

  constructor(row: number, col: number) {
    this._row = row
    this._col = col
    this._tree = new SegmentTree(row, col) // 行 列
  }

  gather(k: number, maxRow: number): number[] {
    if (k > this._col) return []
    maxRow++
    if (this._tree.queryMax(1, maxRow) < k) return []
    const res = { resRow: -1, resCol: -1 }
    this._tree.updateOneRow(1, maxRow, k, res)
    return [res.resRow, res.resCol]
  }

  scatter(k: number, maxRow: number): boolean {
    maxRow++
    if (this._tree.querySum(1, maxRow) < k) return false
    this._tree.updateManyRows(1, maxRow, { value: k })
    return true
  }
}

/**
 * Your BookMyShow object will be instantiated and called as such:
 * var obj = new BookMyShow(n, m)
 * var param_1 = obj.gather(k,maxRow)
 * var param_2 = obj.scatter(k,maxRow)
 */
if (require.main === module) {
  const bookMyShow = new BookMyShow(2, 5)
  console.log(bookMyShow.gather(4, 0))
  console.log(bookMyShow.gather(2, 0))
  console.log(bookMyShow.scatter(5, 1))
  console.log(bookMyShow.scatter(5, 1))
}

export {}
