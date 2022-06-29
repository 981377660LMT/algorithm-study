/**
 * @description 线段树区间最大值RMQ
 */
class MaxSegmentTree {
  private static readonly MIN = -Infinity // !注意是0还是-Infinity
  private readonly tree: number[]
  private readonly lazyValue: number[]
  private readonly isLazy: Uint8Array
  private readonly size: number

  /**
   * @param sizeOrNums 数组长度或数组
   */
  constructor(sizeOrNums: number | number[]) {
    this.size = typeof sizeOrNums === 'number' ? sizeOrNums : sizeOrNums.length
    this.tree = Array(this.size << 2).fill(MaxSegmentTree.MIN)
    this.lazyValue = Array(this.size << 2).fill(0)
    this.isLazy = new Uint8Array(this.size << 2)
    if (Array.isArray(sizeOrNums)) {
      this._build(1, 1, this.size, sizeOrNums)
    }
  }

  query(l: number, r: number): number {
    // this._checkRange(l, r)
    return this._query(1, l, r, 1, this.size)
  }

  update(l: number, r: number, delta: number): void {
    // this._checkRange(l, r)
    this._update(1, l, r, 1, this.size, delta)
  }

  queryAll(): number {
    return this.tree[1]
  }

  private _build(rt: number, l: number, r: number, nums: number[]): void {
    if (l === r) {
      this.tree[rt] = nums[l - 1]
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._build(rt << 1, l, mid, nums)
    this._build((rt << 1) | 1, mid + 1, r, nums)
    this._pushUp(rt)
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this.tree[rt]

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    let res = MaxSegmentTree.MIN
    if (L <= mid) res = Math.max(res, this._query(rt << 1, L, R, l, mid))
    if (mid < R) res = Math.max(res, this._query((rt << 1) | 1, L, R, mid + 1, r))

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, delta: number): void {
    if (L <= l && r <= R) {
      this.lazyValue[rt] += delta
      this.tree[rt] += delta
      this.isLazy[rt] = 1
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, delta)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, delta)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this.tree[rt] = Math.max(this.tree[rt << 1], this.tree[(rt << 1) | 1])
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this.isLazy[rt]) {
      const delta = this.lazyValue[rt]
      this.isLazy[rt << 1] = 1
      this.isLazy[(rt << 1) | 1] = 1
      this.lazyValue[rt << 1] += delta
      this.lazyValue[(rt << 1) | 1] += delta
      this.tree[rt << 1] += delta
      this.tree[(rt << 1) | 1] += delta

      this.lazyValue[rt] = 0
      this.isLazy[rt] = 0
    }
  }

  private _checkRange(l: number, r: number): void {
    if (l < 1 || r > this.size) throw new RangeError(`[${l}, ${r}] out of range: [1, ${this.size}]`)
  }
}

if (require.main === module) {
}

export { MaxSegmentTree }
