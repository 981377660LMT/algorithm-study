/**
 * 线段树区间叠加最大值RMQ
 *
 * !叠加更新可以省去isLazy数组
 */
class MaxSegmentTree {
  private readonly _tree: number[]
  private readonly _lazyValue: number[]
  private readonly _size: number

  /**
   * @param sizeOrNums 数组长度或数组
   */
  constructor(sizeOrNums: number | number[]) {
    this._size = typeof sizeOrNums === 'number' ? sizeOrNums : sizeOrNums.length
    this._tree = Array(this._size << 2).fill(0)
    this._lazyValue = Array(this._size << 2).fill(0)
    if (Array.isArray(sizeOrNums)) {
      this._build(1, 1, this._size, sizeOrNums)
    }
  }

  query(l: number, r: number): number {
    // this._checkRange(l, r)
    return this._query(1, l, r, 1, this._size)
  }

  add(l: number, r: number, delta: number): void {
    // this._checkRange(l, r)
    this._add(1, l, r, 1, this._size, delta)
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
    let res = -Infinity
    if (L <= mid) res = Math.max(res, this._query(rt << 1, L, R, l, mid))
    if (mid < R) res = Math.max(res, this._query((rt << 1) | 1, L, R, mid + 1, r))

    return res
  }

  private _add(rt: number, L: number, R: number, l: number, r: number, delta: number): void {
    if (L <= l && r <= R) {
      this._lazyValue[rt] += delta
      this._tree[rt] += delta
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    if (L <= mid) this._add(rt << 1, L, R, l, mid, delta)
    if (mid < R) this._add((rt << 1) | 1, L, R, mid + 1, r, delta)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this._tree[rt] = Math.max(this._tree[rt << 1], this._tree[(rt << 1) | 1])
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this._lazyValue[rt]) {
      const delta = this._lazyValue[rt]

      this._lazyValue[rt << 1] += delta
      this._lazyValue[(rt << 1) | 1] += delta
      this._tree[rt << 1] += delta
      this._tree[(rt << 1) | 1] += delta

      this._lazyValue[rt] = 0
    }
  }

  private _checkRange(l: number, r: number): void {
    if (l < 1 || r > this._size) {
      throw new RangeError(`[${l}, ${r}] out of range: [1, ${this._size}]`)
    }
  }
}

if (require.main === module) {
  const tree = new MaxSegmentTree([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
  console.log(tree.query(1, 3))
  tree.add(1, 3, 2)
  console.log(tree.query(1, 3))
  console.log(tree.query(1, 1))
  console.log(tree.query(2, 2))
  console.log(tree.queryAll())
  tree.add(1, 10, -20)
  console.log(tree.queryAll())
}

export { MaxSegmentTree }
