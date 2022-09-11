/* eslint-disable no-inner-declarations */

/**
 * @description 线段树区间最大值RMQ
 *
 * 如果查询超出范围 返回0
 */
class RMQSegmentTree {
  private readonly _tree: number[]
  private readonly _lazyValue: number[]
  private readonly _isLazy: Uint8Array
  private readonly _size: number

  /**
   * @param size 区间右边界
   */
  constructor(nOrNums: number | number[]) {
    this._size = Array.isArray(nOrNums) ? nOrNums.length : nOrNums
    this._tree = Array(this._size << 2).fill(0)
    this._lazyValue = Array(this._size << 2).fill(0)
    this._isLazy = new Uint8Array(this._size << 2)
    if (Array.isArray(nOrNums)) this._build(1, 1, this._size, nOrNums)
  }

  query(l: number, r: number): number {
    if (l < 1) l = 1
    if (r > this._size) r = this._size
    if (l > r) return 0 // !超出范围返回0
    return this._query(1, l, r, 1, this._size)
  }

  update(l: number, r: number, target: number): void {
    if (l < 1) l = 1
    if (r > this._size) r = this._size
    if (l > r) return
    this._update(1, l, r, 1, this._size, target)
  }

  queryAll(): number {
    return this._tree[1]
  }

  private _build(rt: number, l: number, r: number, nums: number[]): void {
    if (l === r) {
      this._tree[rt] = nums[l - 1]
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._build(rt << 1, l, mid, nums)
    this._build((rt << 1) | 1, mid + 1, r, nums)
    this._pushUp(rt)
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._tree[rt]

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    let res = 0 // !默认的最小值为0
    if (L <= mid) res = Math.max(res, this._query(rt << 1, L, R, l, mid))
    if (mid < R) res = Math.max(res, this._query((rt << 1) | 1, L, R, mid + 1, r))

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, target: number): void {
    if (L <= l && r <= R) {
      this._lazyValue[rt] = Math.max(this._lazyValue[rt], target)
      this._tree[rt] = Math.max(this._tree[rt], target)
      this._isLazy[rt] = 1
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, target)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, target)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this._tree[rt] = Math.max(this._tree[rt << 1], this._tree[(rt << 1) | 1])
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this._isLazy[rt]) {
      const target = this._lazyValue[rt]
      this._lazyValue[rt << 1] = Math.max(this._lazyValue[rt << 1], target)
      this._lazyValue[(rt << 1) | 1] = Math.max(this._lazyValue[(rt << 1) | 1], target)
      this._isLazy[rt << 1] = 1

      this._tree[rt << 1] = Math.max(this._tree[rt << 1], target)
      this._tree[(rt << 1) | 1] = Math.max(this._tree[(rt << 1) | 1], target)
      this._isLazy[(rt << 1) | 1] = 1

      this._lazyValue[rt] = 0
      this._isLazy[rt] = 0
    }
  }
}

if (require.main === module) {
  function fallingSquares(positions: number[][]): number[] {
    const tmpSet = new Set<number>()
    for (const [left, size] of positions) {
      tmpSet.add(left)
      tmpSet.add(left + size - 1)
    }
    const allNums = [...tmpSet].sort((a, b) => a - b)

    const mapping = new Map<number, number>()
    for (const [index, num] of allNums.entries()) mapping.set(num, index + 1)

    const res = Array<number>(positions.length).fill(0)
    const tree = new RMQSegmentTree(mapping.size + 10)
    for (const [i, [left, size]] of positions.entries()) {
      const right = left + size - 1

      const preHeihgt = tree.query(mapping.get(left)!, mapping.get(right)!)
      tree.update(mapping.get(left)!, mapping.get(right)!, preHeihgt + size)
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
