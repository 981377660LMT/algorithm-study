import assert from 'assert'

class SegmentTree {
  private readonly tree: Int32Array
  private readonly lazyValue: Int32Array
  private readonly isLazy: Uint8Array
  private readonly size: number

  /**
   * @param size 值域间右边界
   * 2 <= A.length <= 50000
     0 <= A[i] <= 50000
   */
  constructor(size: number) {
    this.size = size
    this.tree = new Int32Array(size << 2).fill(-1) // -1 表示不存在 存储最大索引值
    this.lazyValue = new Int32Array(size << 2).fill(-1)
    this.isLazy = new Uint8Array(size << 2)
  }

  query(l: number, r: number): number {
    this._checkBoundsBeginEnd(l, r)
    return this._query(1, l, r, 1, this.size)
  }

  update(l: number, r: number, maxIndex: number): void {
    this._checkBoundsBeginEnd(l, r)
    this._update(1, l, r, 1, this.size, maxIndex)
  }

  queryAll(): number {
    return this.tree[1]
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this.tree[rt]

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    let res = -1
    if (L <= mid) res = Math.max(res, this._query(rt << 1, L, R, l, mid))
    if (mid < R) res = Math.max(res, this._query((rt << 1) | 1, L, R, mid + 1, r))

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, maxIndex: number): void {
    if (L <= l && r <= R) {
      this.lazyValue[rt] = Math.max(this.lazyValue[rt], maxIndex)
      this.tree[rt] = Math.max(this.tree[rt], maxIndex)
      this.isLazy[rt] = 1
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, maxIndex)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, maxIndex)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this.tree[rt] = Math.max(this.tree[rt << 1], this.tree[(rt << 1) | 1])
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this.isLazy[rt]) {
      const target = this.lazyValue[rt]
      this.lazyValue[rt << 1] = Math.max(this.lazyValue[rt << 1], target)
      this.lazyValue[(rt << 1) | 1] = Math.max(this.lazyValue[(rt << 1) | 1], target)
      this.isLazy[rt << 1] = 1
      this.tree[rt << 1] = Math.max(this.tree[rt << 1], target)
      this.tree[(rt << 1) | 1] = Math.max(this.tree[(rt << 1) | 1], target)
      this.isLazy[(rt << 1) | 1] = 1

      this.lazyValue[rt] = -1
      this.isLazy[rt] = 0
    }
  }

  private _checkBoundsBeginEnd(begin: number, end: number): void {
    if (begin < 1 || begin > end || end > this.size) {
      throw new RangeError(`[${begin}, ${end}] out of range: [1, ${this.size}]`)
    }
  }
}

/**
 * 对每个数，寻找右侧最后一个比自己大的数；注意线段树要偏移
 * 值域线段树 维护每个数的最大索引
 *
 * !@param nums 0 <= nums[i] <= 50000  2 <= nums.length <= 50000
 * @returns 每个元素右侧最后一个比自己大的数的索引
 */
function findLastLarge(nums: number[]): number[] {
  const n = nums.length
  const res = Array<number>(n).fill(-1)
  const max = Math.max(...nums)

  const tree = new SegmentTree(max + 10)
  for (let i = n - 1; i >= 0; i--) {
    res[i] = tree.query(nums[i] + 1, nums[i] + 1) // 偏移量为1
    tree.update(1, nums[i] + 1, i)
  }

  return res
}

if (require.main === module) {
  assert.deepStrictEqual(
    findLastLarge([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]),
    [9, 9, 9, 9, 9, 9, 9, 9, 9, -1]
  )

  assert.deepStrictEqual(
    findLastLarge([9, 8, 1, 0, 1, 9, 4, 0, 4, 1]),
    [5, 5, 9, 9, 9, -1, 8, 9, -1, -1]
  )
}

export { findLastLarge }
