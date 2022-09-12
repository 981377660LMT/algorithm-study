// 如果只是用区间叠加 可以用树状数组代替

class SegmentTree {
  private readonly _tree: number[]
  private readonly _lazyValue: number[]
  private readonly _size: number

  /**
   *
   * @param size 区间右边界
   */
  constructor(nOrNums: number | number[]) {
    this._size = Array.isArray(nOrNums) ? nOrNums.length : nOrNums
    this._tree = Array(this._size << 2).fill(0)
    this._lazyValue = Array(this._size << 2).fill(0)
    if (Array.isArray(nOrNums)) this._build(1, 1, this._size, nOrNums)
  }

  query(left: number, right: number): number {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return 0 // !超出范围返回0
    return this._query(1, left, right, 1, this._size)
  }

  update(left: number, right: number, delta: number): void {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return
    this._update(1, left, right, 1, this._size, delta)
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
    let res = 0
    if (L <= mid) res += this._query(rt << 1, L, R, l, mid)
    if (mid < R) res += this._query((rt << 1) | 1, L, R, mid + 1, r)

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, delta: number): void {
    if (L <= l && r <= R) {
      this._lazyValue[rt] += delta
      this._tree[rt] += delta * (r - l + 1)
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, delta)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, delta)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this._tree[rt] = this._tree[rt << 1] + this._tree[(rt << 1) | 1]
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this._lazyValue[rt]) {
      const delta = this._lazyValue[rt]
      this._lazyValue[rt << 1] += delta
      this._lazyValue[(rt << 1) | 1] += delta
      this._tree[rt << 1] += delta * (mid - l + 1)
      this._tree[(rt << 1) | 1] += delta * (r - mid)
      this._lazyValue[rt] = 0
    }
  }
}

if (require.main === module) {
  const sg = new SegmentTree(10)
  sg.update(2, 3, 2)
  console.log(sg.query(1, 8))
  console.log(sg.query(1, 1))
  console.log(sg.queryAll())
}

export {}
