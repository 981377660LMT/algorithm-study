/**
 * @description 线段树区间最大值RMQ
 */
class MaxSegmentTree {
  private readonly tree: number[]
  private readonly lazyValue: number[]
  private readonly isLazy: Uint8Array
  private readonly size: number

  /**
   *
   * @param size 区间右边界
   */
  constructor(size: number) {
    this.size = size
    this.tree = Array(size << 2).fill(0)
    this.lazyValue = Array(size << 2).fill(0)
    this.isLazy = new Uint8Array(size << 2)
  }

  query(l: number, r: number): number {
    this.checkRange(l, r)
    return this._query(1, l, r, 1, this.size)
  }

  update(l: number, r: number, delta: number): void {
    this.checkRange(l, r)
    this._update(1, l, r, 1, this.size, delta)
  }

  queryAll(): number {
    return this.tree[1]
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this.tree[rt]

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    let res = -Infinity
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

  private checkRange(l: number, r: number): void {
    if (l < 1 || r > this.size) throw new RangeError(`[${l}, ${r}] out of range: [1, ${this.size}]`)
  }
}

if (require.main === module) {
}

export {}
