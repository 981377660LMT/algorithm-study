/* eslint-disable no-inner-declarations */
// !两种颜色的区间染色不需要用懒标记记录更新状态 直接用节点值判断
class SegmentTree {
  private readonly _tree: Uint32Array
  private readonly _size: number

  /**
   * @param size 区间右边界
   */
  constructor(size: number) {
    this._size = size
    this._tree = new Uint32Array(size << 2)
  }

  /**
   * 1 <= left <= right <= size
   */
  query(left: number, right: number): number {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return 0
    return this._query(1, left, right, 1, this._size)
  }

  /**
   * 1 <= left <= right <= size
   */
  update(left: number, right: number, target: 0 | 1): void {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return
    this._update(1, left, right, 1, this._size, target)
  }

  queryAll(): number {
    return this._tree[1]
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._tree[rt]

    const mid = Math.floor((l + r) / 2)
    this.pushDown(rt, l, r, mid)
    let res = 0
    if (L <= mid) res += this._query(rt << 1, L, R, l, mid)
    if (mid < R) res += this._query((rt << 1) | 1, L, R, mid + 1, r)

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, target: 0 | 1): void {
    if (L <= l && r <= R) {
      this._tree[rt] = target === 1 ? r - l + 1 : 0
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, target)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, target)
    this.pushUp(rt)
  }

  private pushUp(rt: number): void {
    this._tree[rt] = this._tree[rt << 1] + this._tree[(rt << 1) | 1]
  }

  private pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this._tree[rt] === r - l + 1) {
      this._tree[rt << 1] = mid - l + 1
      this._tree[(rt << 1) | 1] = r - mid
    } else if (this._tree[rt] === 0) {
      this._tree[rt << 1] = 0
      this._tree[(rt << 1) | 1] = 0
    }
  }
}

//! ///////////////////////////////////////////////////
// !写isLazy的版本
class SegmentTree2 {
  private readonly _tree: Uint32Array
  private readonly _lazyValue: Uint8Array
  private readonly _isLazy: Uint8Array
  private readonly _size: number

  /**
   *
   * @param size 区间右边界
   */
  constructor(size: number) {
    this._size = size
    this._tree = new Uint32Array(size << 2)
    this._lazyValue = new Uint8Array(size << 2)
    this._isLazy = new Uint8Array(size << 2)
  }

  query(left: number, right: number): number {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return 0
    return this._query(1, left, right, 1, this._size)
  }

  update(left: number, right: number, target: 0 | 1): void {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return
    this._update(1, left, right, 1, this._size, target)
  }

  queryAll(): number {
    return this._tree[1]
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._tree[rt]

    const mid = Math.floor((l + r) / 2)
    this.pushDown(rt, l, r, mid)
    let res = 0
    if (L <= mid) res += this._query(rt << 1, L, R, l, mid)
    if (mid < R) res += this._query((rt << 1) | 1, L, R, mid + 1, r)

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, target: 0 | 1): void {
    if (L <= l && r <= R) {
      this._isLazy[rt] = 1
      this._lazyValue[rt] = target
      this._tree[rt] = target * (r - l + 1)
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, target)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, target)
    this.pushUp(rt)
  }

  private pushUp(rt: number): void {
    this._tree[rt] = this._tree[rt << 1] + this._tree[(rt << 1) | 1]
  }

  private pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this._isLazy[rt]) {
      const target = this._lazyValue[rt]
      this._lazyValue[rt << 1] = target
      this._lazyValue[(rt << 1) | 1] = target
      this._tree[rt << 1] = target * (mid - l + 1)
      this._tree[(rt << 1) | 1] = target * (r - mid)
      this._isLazy[rt << 1] = 1
      this._isLazy[(rt << 1) | 1] = 1

      this._lazyValue[rt] = 0
      this._isLazy[rt] = 0
    }
  }
}

if (require.main === module) {
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

  console.log(
    amountPainted([
      [6, 17],
      [3, 6],
      [7, 17],
      [16, 20],
      [2, 20]
    ])
  )
}

export {}
