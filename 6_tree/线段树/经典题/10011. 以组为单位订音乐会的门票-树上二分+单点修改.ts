// 用一个线段树维护区间最值、区间和
// 需要找第一个行>=k，类似主席树找mex，检查区间max线段树上二分即可(到某个节点检查是否成立)

// https://leetcode.cn/problems/booking-concert-tickets-in-groups/submissions/
// 0. 线段树维护区间的 sum 和区间的 max
// 1. gather 操作里需要找到第一个 `剩余座位>=k` 的行，可以通过查询左右区间里的 `max` 来确定向左子树递归还是向右子树递归，比较简单的写法是 `在入口判断区间max是否不小于k`
// 2. scatter 操作里需要按行号从小到大依次填充k个位置，这个过程就是二叉树的后序遍历，先遍历左子树再遍历右子树找到最左边的结点，先填充完小的行再填充大的行
// 3. gather 操作需要返回 `[第一个行，第一个列]`，比较方便的做法是在update操作里多传一个引用类型的参数 `resRef`，在递归过程中直接修改返回值
// 4. 由于每次调用api都需要单点修改线段树叶子值，所以不需要懒标记和 pushDown 操作

// !都是到叶子节点的单点更新 不用加懒标记
class SegmentTree {
  private readonly sums: number[] // 每行剩余座位数
  private readonly maxs: Uint32Array // 区间最大值
  // private readonly isLazy: Uint8Array // 这里没有区间更新 所以可以不用加懒标记
  // private readonly lazyValue: number[] // 这里没有区间更新 所以可以不用加懒标记
  private readonly rowSize: number
  private readonly colSize: number

  /**
   * @param rowSize 区间右边界
   */
  constructor(rowSize: number, colSize: number) {
    this.rowSize = rowSize
    this.colSize = colSize
    this.sums = Array(rowSize << 2).fill(0)
    this.maxs = new Uint32Array(rowSize << 2)
    // this.lazyValue = Array(size << 2).fill(0)
    // this.isLazy = new Uint8Array(size << 2)
    this._build(1, 1, rowSize, colSize)
  }

  queryMax(l: number, r: number): number {
    this._checkRange(l, r)
    return this._queryMax(1, l, r, 1, this.rowSize)
  }

  querySum(l: number, r: number): number {
    this._checkRange(l, r)
    return this._querySum(1, l, r, 1, this.rowSize)
  }

  // 找到最小的行满足sum>=k的位置
  updateOneRow(
    l: number,
    r: number,
    delta: number,
    resRef: { resRow: number; resCol: number }
  ): void {
    this._checkRange(l, r)
    this._updateOneRow(1, l, r, 1, this.rowSize, delta, resRef)
  }

  // 从最小的行开始填充delta
  updateManyRows(l: number, r: number, deltaRef: { value: number }): void {
    this._checkRange(l, r)
    this._updateManyRows(1, l, r, 1, this.rowSize, deltaRef)
  }

  // queryAll(): number {
  //   return this.tree[1]
  // }

  private _build(rt: number, l: number, r: number, colSize: number): void {
    if (l === r) {
      this.sums[rt] = colSize
      this.maxs[rt] = colSize
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._build(rt << 1, l, mid, colSize)
    this._build((rt << 1) | 1, mid + 1, r, colSize)
    this._pushUp(rt)
  }

  private _queryMax(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this.maxs[rt]

    const mid = Math.floor((l + r) / 2)
    // this._pushDown(rt, l, r, mid)
    let res = 0
    if (L <= mid) res = Math.max(res, this._queryMax(rt << 1, L, R, l, mid))
    if (mid < R) res = Math.max(res, this._queryMax((rt << 1) | 1, L, R, mid + 1, r))

    return res
  }

  private _querySum(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this.sums[rt]

    const mid = Math.floor((l + r) / 2)
    // this._pushDown(rt, l, r, mid)
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
    delta: number,
    resRef: { resRow: number; resCol: number } // 返回值
  ): void {
    if (this.maxs[rt] < delta) return // 二分看max是否满足 简单的写法就是在入口判断

    // 单点修改 找到答案了
    if (l === r) {
      resRef.resRow = l - 1
      resRef.resCol = this.colSize - this.sums[rt]
      this.sums[rt] -= delta
      this.maxs[rt] -= delta
      return
    }

    const mid = Math.floor((l + r) / 2)
    // this._pushDown(rt, l, r, mid)
    if (resRef.resRow === -1 && L <= mid) {
      this._updateOneRow(rt << 1, L, R, l, mid, delta, resRef)
    }
    if (resRef.resRow === -1 && mid < R) {
      this._updateOneRow((rt << 1) | 1, L, R, mid + 1, r, delta, resRef)
    }
    this._pushUp(rt)
  }

  // scatter 尽量往左填充delta 即二叉树前序遍历(先遍历左子树再遍历右子树)
  private _updateManyRows(
    rt: number,
    L: number,
    R: number,
    l: number,
    r: number,
    deltaRef: { value: number }
  ): void {
    // 单点修改
    if (l === r) {
      const remain = Math.min(deltaRef.value, this.sums[rt])
      this.sums[rt] -= remain
      this.maxs[rt] -= remain
      deltaRef.value -= remain
      return
    }

    const mid = Math.floor((l + r) / 2)
    // this._pushDown(rt, l, r, mid)
    if (deltaRef.value > 0 && L <= mid) {
      this._updateManyRows(rt << 1, L, R, l, mid, deltaRef)
    }
    if (deltaRef.value > 0 && mid < R) {
      this._updateManyRows((rt << 1) | 1, L, R, mid + 1, r, deltaRef)
    }
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this.sums[rt] = this.sums[rt << 1] + this.sums[(rt << 1) | 1]
    this.maxs[rt] = Math.max(this.maxs[rt << 1], this.maxs[(rt << 1) | 1])
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    // if (this.isLazy[rt]) {
    //   const target = this.sum[rt]
    //   this.sum[rt << 1] = Math.max(this.sum[rt << 1], target)
    //   this.sum[(rt << 1) | 1] = Math.max(this.sum[(rt << 1) | 1], target)
    //   this.isLazy[rt << 1] = 1
    //   this.tree[rt << 1] = Math.max(this.tree[rt << 1], target)
    //   this.tree[(rt << 1) | 1] = Math.max(this.tree[(rt << 1) | 1], target)
    //   this.isLazy[(rt << 1) | 1] = 1
    //   this.sum[rt] = Infinity
    //   this.isLazy[rt] = 0
    // }
  }

  private _checkRange(l: number, r: number): void {
    if (l < 1 || r > this.rowSize) {
      throw new RangeError(`[${l}, ${r}] out of range: [1, ${this.rowSize}]`)
    }
  }
}

class BookMyShow {
  private readonly tree: SegmentTree
  private readonly row: number
  private readonly col: number

  constructor(row: number, col: number) {
    this.row = row
    this.col = col
    this.tree = new SegmentTree(row, col) // 行 列
  }

  gather(k: number, maxRow: number): number[] {
    if (k > this.col) return []
    maxRow++
    if (this.tree.queryMax(1, maxRow) < k) return []
    const res = { resRow: -1, resCol: -1 }
    this.tree.updateOneRow(1, maxRow, k, res)
    return [res.resRow, res.resCol]
  }

  scatter(k: number, maxRow: number): boolean {
    maxRow++
    if (this.tree.querySum(1, maxRow) < k) return false
    this.tree.updateManyRows(1, maxRow, { value: k })
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
