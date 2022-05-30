function amountPainted(paint: number[][]): number[] {
  const res: number[] = []
  const tree = new SegmentTree(Math.max(...paint.flat()))
  for (let [start, end] of paint) {
    ;[start, end] = [start + 1, end + 1]
    const cur = tree.query(start, end - 1)
    res.push(end - start - cur)
    tree.update(start, end - 1, 1)
  }

  return res
}

class SegmentTree {
  private readonly tree: Uint32Array
  private readonly lazyValue: Uint8Array
  private readonly isLazy: Uint8Array
  private readonly size: number

  /**
   *
   * @param size 区间右边界
   */
  constructor(size: number) {
    this.size = size
    this.tree = new Uint32Array(size << 2)
    this.lazyValue = new Uint8Array(size << 2)
    this.isLazy = new Uint8Array(size << 2)
  }

  query(l: number, r: number): number {
    this.checkRange(l, r)
    return this._query(1, l, r, 1, this.size)
  }

  update(l: number, r: number, target: number): void {
    this.checkRange(l, r)
    this._update(1, l, r, 1, this.size, target)
  }

  queryAll(): number {
    return this.tree[1]
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this.tree[rt]

    const mid = Math.floor((l + r) / 2)
    this.pushDown(rt, l, r, mid)
    let res = 0
    if (L <= mid) res += this._query(rt << 1, L, R, l, mid)
    if (mid < R) res += this._query((rt << 1) | 1, L, R, mid + 1, r)

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, target: number): void {
    if (L <= l && r <= R) {
      this.isLazy[rt] = 1
      this.lazyValue[rt] = target
      this.tree[rt] = target === 1 ? r - l + 1 : 0
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, target)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, target)
    this.pushUp(rt)
  }

  private pushUp(rt: number): void {
    this.tree[rt] = this.tree[rt << 1] + this.tree[(rt << 1) | 1]
  }

  private pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this.isLazy[rt]) {
      const target = this.lazyValue[rt]
      this.lazyValue[rt << 1] = target
      this.lazyValue[(rt << 1) | 1] = target
      this.tree[rt << 1] = target === 1 ? mid - l + 1 : 0
      this.tree[(rt << 1) | 1] = target === 1 ? r - mid : 0
      this.isLazy[rt << 1] = 1
      this.isLazy[(rt << 1) | 1] = 1

      this.lazyValue[rt] = 0
      this.isLazy[rt] = 0
    }
  }

  private checkRange(l: number, r: number): void {
    if (l < 1 || r > this.size) throw new RangeError(`[${l}, ${r}] out of range: [1, ${this.size}]`)
  }
}

if (require.main === module) {
  console.log(
    amountPainted([
      [6, 17],
      [3, 6],
      [7, 17],
      [16, 20],
      [2, 20],
    ])
  )
}

export {}
